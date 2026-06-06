package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/joho/godotenv"
	"github.com/rs/cors"

	"difference-engine/mixer"
	"difference-engine/store"
)

type appConfig struct {
	port    string
	siteDir string
	dbPath  string
	mixer   mixer.Config
}

func loadConfig() appConfig {
	_ = godotenv.Load()
	return appConfig{
		port:    getEnv("PORT", "8000"),
		siteDir: getEnv("SITE_DIR", "frontend"),
		dbPath:  getEnv("DB_PATH", "de.db"),
		mixer: mixer.Config{
			StemsDir:          getEnv("STEMS_DIR", ""),
			OutputDir:         getEnv("OUTPUT_DIR", "./output"),
			R2AccountID:       os.Getenv("R2_ACCOUNT_ID"),
			R2AccessKeyID:     os.Getenv("R2_ACCESS_KEY_ID"),
			R2SecretAccessKey: os.Getenv("R2_SECRET_ACCESS_KEY"),
			R2StemsBucket:     os.Getenv("R2_STEMS_BUCKET"),
		},
	}
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

type server struct {
	mixer   *mixer.Mixer
	siteDir string
	store   *store.Store
}

type statusRecorder struct {
	http.ResponseWriter
	status int
}

func (r *statusRecorder) WriteHeader(code int) {
	r.status = code
	r.ResponseWriter.WriteHeader(code)
}

func logRequests(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		rec := &statusRecorder{ResponseWriter: w, status: http.StatusOK}
		next.ServeHTTP(rec, r)
		log.Printf("%s %s %s %d %v",
			r.RemoteAddr, r.Method, r.URL.Path, rec.status, time.Since(start).Round(time.Millisecond),
		)
	})
}

var safeSegment = regexp.MustCompile(`^[A-Za-z0-9_-]+$`)

func (s *server) routes() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /stems/{track}/count", s.handleStemCount)
	mux.HandleFunc("GET /stems/{track}/{stem}", s.handleStem)
	mux.HandleFunc("POST /mixdown", s.handleMixdown)
	mux.HandleFunc("GET /mixdown/{id}", s.handleRecall)
	// SITE_DIR is only set in local dev; in production Caddy serves the static site.
	if s.siteDir != "" {
		mux.Handle("/", http.FileServer(http.Dir(s.siteDir)))
	}

	return cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost", "http://localhost:8000"},
		AllowCredentials: true,
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"*"},
		ExposedHeaders:   []string{"Content-Disposition", "X-Mixdown-ID"},
	}).Handler(mux)
}

type mixdownRequest struct {
	DeValues struct {
		Track   string        `json:"track"`
		Stems   []json.Number `json:"stems"`
		Volumes []float64     `json:"volumes"`
	} `json:"de_values"`
}

func (s *server) handleStemCount(w http.ResponseWriter, r *http.Request) {
	track := r.PathValue("track")
	if !safeSegment.MatchString(track) {
		http.Error(w, "invalid track", http.StatusBadRequest)
		return
	}
	count, err := s.mixer.CountStems(track)
	if err != nil {
		log.Printf("count stems %s: %v", track, err)
		http.Error(w, "failed to count stems", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	fmt.Fprintf(w, `{"count":%d}`, count)
}

func (s *server) handleStem(w http.ResponseWriter, r *http.Request) {
	track := r.PathValue("track")
	stem := r.PathValue("stem")
	if !safeSegment.MatchString(track) || !safeSegment.MatchString(stem) {
		http.Error(w, "invalid path", http.StatusBadRequest)
		return
	}

	paths, err := s.mixer.GetStems(track, []string{stem})
	if err != nil {
		log.Printf("get stem %s/%s: %v", track, stem, err)
		http.NotFound(w, r)
		return
	}

	// R2 mode: presigned URL — redirect the browser directly to it
	if strings.HasPrefix(paths[0], "http") {
		http.Redirect(w, r, paths[0], http.StatusTemporaryRedirect)
		return
	}

	http.ServeFile(w, r, paths[0])
}

func (s *server) handleMixdown(w http.ResponseWriter, r *http.Request) {
	var req mixdownRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, fmt.Sprintf("invalid request: %v", err), http.StatusBadRequest)
		return
	}

	dv := req.DeValues
	if dv.Track == "" {
		http.Error(w, "track is required", http.StatusBadRequest)
		return
	}
	if len(dv.Stems) == 0 {
		http.Error(w, "no stems provided", http.StatusBadRequest)
		return
	}

	// json.Number.String() matches Python's f"{stem}" for integer stems like 3 → "3"
	stems := make([]string, len(dv.Stems))
	for i, n := range dv.Stems {
		stems[i] = n.String()
	}

	vols := dv.Volumes
	if len(vols) > len(stems) {
		vols = vols[:len(stems)]
	}
	rounded := make([]float64, len(vols))
	for i, v := range vols {
		rounded[i] = math.Round(v*1000) / 1000
	}

	id, err := s.store.RecordRequest(dv.Track, stems, rounded)
	if err != nil {
		log.Printf("record request: %v", err)
	}

	s.serveMixdown(w, r, dv.Track, stems, rounded, id)
}

func (s *server) handleRecall(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}

	req, err := s.store.GetRequest(id)
	if err != nil {
		log.Printf("get request %d: %v", id, err)
		http.Error(w, "database error", http.StatusInternalServerError)
		return
	}
	if req == nil {
		http.NotFound(w, r)
		return
	}

	s.serveMixdown(w, r, req.Track, req.Stems, req.Volumes, req.ID)
}

// serveMixdown generates the audio from the given parameters and writes it as a response.
func (s *server) serveMixdown(w http.ResponseWriter, r *http.Request, track string, stems []string, volumes []float64, id int64) {
	inputFiles, err := s.mixer.GetStems(track, stems)
	if err != nil {
		log.Printf("get stems: %v", err)
		http.Error(w, "failed to get stems", http.StatusInternalServerError)
		return
	}

	filePath, err := s.mixer.CreateMixdown(track, inputFiles, volumes)
	if err != nil {
		log.Printf("create mixdown: %v", err)
		http.Error(w, "failed to create mixdown", http.StatusInternalServerError)
		return
	}

	if id > 0 {
		w.Header().Set("X-Mixdown-ID", strconv.FormatInt(id, 10))
	}
	filename := fmt.Sprintf("%s_%s.mp3", track, time.Now().UTC().Format("20060102_150405"))
	w.Header().Set("Content-Disposition", fmt.Sprintf(`attachment; filename="%s"`, filename))
	w.Header().Set("Content-Type", "audio/mpeg")
	http.ServeFile(w, r, filePath)
}

func main() {
	cfg := loadConfig()

	st, err := store.Open(cfg.dbPath)
	if err != nil {
		log.Fatalf("open store: %v", err)
	}
	defer st.Close()

	srv := &server{
		mixer:   mixer.New(cfg.mixer),
		siteDir: cfg.siteDir,
		store:   st,
	}

	log.Printf("listening on :%s", cfg.port)
	log.Fatal(http.ListenAndServe(":"+cfg.port, logRequests(srv.routes())))
}
