const btn = document.getElementById("downloadBtn");
const aboutBtn = document.getElementById("aboutBtn");
const about = document.getElementById("about");
const aboutClose = document.getElementById("aboutClose");

aboutBtn.addEventListener("click", () => {
    about.hidden = false;
    aboutBtn.setAttribute("aria-expanded", "true");
    aboutClose.focus();
});

function closeAbout() {
    about.hidden = true;
    aboutBtn.setAttribute("aria-expanded", "false");
}

aboutClose.addEventListener("click", closeAbout);
about.addEventListener("click", (e) => { if (e.target === about) closeAbout(); });
document.addEventListener("keydown", (e) => { if (e.key === "Escape" && !about.hidden) closeAbout(); });


async function generateRandomNumbers() {
    const stems = Array.from({length: 8}, (_, i) => i + 1);

    // Fisher-Yates shuffle, pick first 3
    for (let i = stems.length - 1; i > 0; i--) {
        const j = Math.floor(Math.random() * (i + 1));
        [stems[i], stems[j]] = [stems[j], stems[i]];
    }
    const selectedStems = stems.slice(0, 3);

    const volumes = Array.from({length: 32}, () => Math.random());

    btn.disabled = true;
    btn.textContent = "Generating…";

    try {
        const response = await fetch("/mixdown", {
            method: "POST",
            headers: {"Content-Type": "application/json"},
            body: JSON.stringify({de_values: {stems: selectedStems, volumes}}),
        });

        if (!response.ok) throw new Error(`Server error: ${response.status}`);

        const blob = await response.blob();
        const url = URL.createObjectURL(blob);
        const a = document.createElement("a");
        a.href = url;
        a.download = "mixdown.mp3";
        a.click();
        URL.revokeObjectURL(url);
    } catch (err) {
        alert(`Something went wrong: ${err.message}`);
    } finally {
        btn.disabled = false;
        btn.textContent = "Download";
    }
}

btn.addEventListener("click", generateRandomNumbers);
