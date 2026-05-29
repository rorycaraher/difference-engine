package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math"
	"net/http"
	"os"

	"github.com/rs/cors"

	"difference-engine/mixer"
)

type mixdownRequest struct {
	DeValues struct {
		Stems   []json.Number `json:"stems"`
		Volumes []float64     `json:"volumes"`
	} `json:"de_values"`
}

var m = &mixer.Mixer{}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("POST /mixdown", handleMixdown)
	mux.Handle("/", http.FileServer(http.Dir("site")))

	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost", "http://localhost:8000"},
		AllowCredentials: true,
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"*"},
	})

	port := os.Getenv("PORT")
	if port == "" {
		port = "8000"
	}

	log.Printf("listening on :%s", port)
	log.Fatal(http.ListenAndServe(":"+port, c.Handler(mux)))
}

func handleMixdown(w http.ResponseWriter, r *http.Request) {
	var req mixdownRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, fmt.Sprintf("invalid request: %v", err), http.StatusBadRequest)
		return
	}

	dv := req.DeValues
	if len(dv.Stems) == 0 {
		http.Error(w, "no stems provided", http.StatusBadRequest)
		return
	}

	// Convert json.Number stems to strings, matching Python f"{stem}" behaviour
	stems := make([]string, len(dv.Stems))
	for i, n := range dv.Stems {
		stems[i] = n.String()
	}

	// Take only as many volumes as there are stems, round to 3 decimal places
	vols := dv.Volumes
	if len(vols) > len(stems) {
		vols = vols[:len(stems)]
	}
	rounded := make([]float64, len(vols))
	for i, v := range vols {
		rounded[i] = math.Round(v*1000) / 1000
	}

	inputFiles, err := m.GetStems(mixer.StemsDir, stems)
	if err != nil {
		log.Printf("get stems: %v", err)
		http.Error(w, "failed to get stems", http.StatusInternalServerError)
		return
	}

	filePath, err := m.CreateMixdown(inputFiles, rounded)
	if err != nil {
		log.Printf("create mixdown: %v", err)
		http.Error(w, "failed to create mixdown", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Disposition", `attachment; filename="mixdown.mp3"`)
	w.Header().Set("Content-Type", "audio/mpeg")
	http.ServeFile(w, r, filePath)
}
