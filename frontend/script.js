const TRACK = "first-principles";
const PICK = 3;

function shuffle(arr) {
    for (let i = arr.length - 1; i > 0; i--) {
        const j = Math.floor(Math.random() * (i + 1));
        [arr[i], arr[j]] = [arr[j], arr[i]];
    }
    return arr;
}

const btn = document.getElementById("playBtn");
const progressFill = document.getElementById("progress-fill");
const refreshBtn = document.getElementById("refreshBtn");
let players = [];
let selectedStems = [];
let volumes = [];
let playing = false;
let endCount = 0;
let rafId = null;

function updateProgress() {
    const p = players[0];
    if (p && p.duration) {
        progressFill.style.width = (p.currentTime / p.duration * 100) + "%";
    }
    if (playing) rafId = requestAnimationFrame(updateProgress);
}

async function togglePlay() {
    if (!playing) {
        btn.disabled = true;
        btn.classList.add("loading");
        try {
            await Promise.all(players.map(p => p.play()));
            playing = true;
            btn.classList.add("playing");
            rafId = requestAnimationFrame(updateProgress);
        } catch (err) {
            console.error("playback error:", err);
        } finally {
            btn.disabled = false;
            btn.classList.remove("loading");
        }
    } else {
        players.forEach(p => p.pause());
        playing = false;
        btn.classList.remove("playing");
        cancelAnimationFrame(rafId);
    }
}

async function init() {
    const res = await fetch(`/stems/${TRACK}/count`);
    if (!res.ok) throw new Error(`stems list failed: ${res.status}`);
    const { stems } = await res.json();

    selectedStems = shuffle([...stems]).slice(0, PICK);
    volumes = selectedStems.map(() => Math.random() * 0.5 + 0.5);

    players = selectedStems.map((stem, i) => {
        const audio = new Audio(`/stems/${TRACK}/${stem}`);
        audio.volume = Math.min(1, volumes[i]);
        audio.preload = "auto";
        return audio;
    });

    players.forEach(p => {
        p.addEventListener("ended", () => {
            endCount++;
            if (endCount === players.length) {
                playing = false;
                endCount = 0;
                btn.classList.remove("playing");
                cancelAnimationFrame(rafId);
                progressFill.style.width = "0%";
                players.forEach(p => { p.currentTime = 0; });
            }
        });
    });

    btn.disabled = false;
}

btn.disabled = true;
init().catch(err => {
    console.error("init error:", err);
    btn.disabled = false;
});

btn.addEventListener("click", togglePlay);

document.addEventListener("keydown", e => {
    if (e.code === "Space" && e.target === document.body) {
        e.preventDefault();
        togglePlay();
    }
    if (e.key === "Escape") closeAbout();
});

const aboutModal = document.getElementById("about-modal");
const aboutLink = document.getElementById("aboutLink");
const aboutClose = document.getElementById("about-close");

function openAbout(e) {
    e.preventDefault();
    aboutModal.classList.add("open");
    aboutClose.focus();
}

function closeAbout() {
    aboutModal.classList.remove("open");
}

aboutLink.addEventListener("click", openAbout);
aboutClose.addEventListener("click", closeAbout);
aboutModal.addEventListener("click", e => {
    if (e.target === aboutModal) closeAbout();
});

const dlBtn = document.getElementById("downloadBtn");

refreshBtn.addEventListener("click", async () => {
    if (playing) {
        cancelAnimationFrame(rafId);
        playing = false;
        btn.classList.remove("playing");
    }
    players.forEach(p => { p.pause(); p.src = ""; });
    players = [];
    endCount = 0;
    progressFill.style.width = "0%";

    btn.disabled = true;
    dlBtn.disabled = true;
    refreshBtn.disabled = true;
    refreshBtn.classList.add("loading");
    try {
        await init();
    } catch (err) {
        console.error("refresh error:", err);
    } finally {
        refreshBtn.classList.remove("loading");
        refreshBtn.disabled = false;
        dlBtn.disabled = false;
    }
});

dlBtn.addEventListener("click", async () => {
    dlBtn.disabled = true;
    dlBtn.classList.add("loading");
    try {
        const response = await fetch("/mixdown", {
            method: "POST",
            headers: {"Content-Type": "application/json"},
            body: JSON.stringify({de_values: {track: TRACK, stems: selectedStems, volumes}}),
        });
        if (!response.ok) throw new Error(`Server error: ${response.status}`);

        const disposition = response.headers.get("Content-Disposition") ?? "";
        const match = disposition.match(/filename="([^"]+)"/);
        const filename = match ? match[1] : "mixdown.mp3";

        const blob = await response.blob();
        const url = URL.createObjectURL(blob);
        const a = document.createElement("a");
        a.href = url;
        a.download = filename;
        a.click();
        URL.revokeObjectURL(url);
    } catch (err) {
        console.error("download error:", err);
    } finally {
        dlBtn.disabled = false;
        dlBtn.classList.remove("loading");
    }
});
