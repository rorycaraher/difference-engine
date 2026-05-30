const TRACK = "first-principles";
const TOTAL_STEMS = 8;
const PICK = 3;

function pickStems() {
    const all = Array.from({length: TOTAL_STEMS}, (_, i) => i + 1);
    for (let i = all.length - 1; i > 0; i--) {
        const j = Math.floor(Math.random() * (i + 1));
        [all[i], all[j]] = [all[j], all[i]];
    }
    return all.slice(0, PICK);
}

const selectedStems = pickStems();
const volumes = selectedStems.map(() => Math.random() * 0.5 + 0.5);

const players = selectedStems.map((stem, i) => {
    const audio = new Audio(`/stems/${TRACK}/${stem}`);
    audio.volume = Math.min(1, volumes[i]);
    audio.preload = "auto";
    return audio;
});

const btn = document.getElementById("playBtn");
const progressFill = document.getElementById("progress-fill");
let playing = false;
let endCount = 0;
let rafId = null;

function updateProgress() {
    const p = players[0];
    if (p.duration) {
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

btn.addEventListener("click", togglePlay);

document.addEventListener("keydown", e => {
    if (e.code === "Space" && e.target === document.body) {
        e.preventDefault();
        togglePlay();
    }
});

const dlBtn = document.getElementById("downloadBtn");

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
