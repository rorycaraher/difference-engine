function generateRandomNumbers() {

    const stems = Array.from({length: 8}, (_, i) => i + 1);
    const volumes = [];

    // Shuffle the array to randomize the order
    for (let i = stems.length - 1; i > 0; i--) {
        const j = Math.floor(Math.random() * (i + 1));
        [stems[i], stems[j]] = [stems[j], stems[i]];
    }
    const selectedStems = stems.slice(0, 3);

    // generate a bunch of random values
    for (let i = 32 - 1; i > 0; i--) {
        volumes.push(Math.random())
    }

    // print the value
    const lvg_values = {
        stems: selectedStems,
        volumes: volumes
    }

    console.log(lvg_values);

    const xhr = new XMLHttpRequest();
    xhr.open("POST", "http://localhost:5000/test_mixdown", true);
    xhr.setRequestHeader("Content-Type", "application/json;charset=UTF-8");
    xhr.send(JSON.stringify({ lvg_values: lvg_values }));
    xhr.onload = function() {
        if (xhr.status === 200) {
            alert("Worked!");
        } else {
            alert("Something wrong!");
        }
    };
}


function testDatabase() {

    const stems = Array.from({length: 8}, (_, i) => i + 1);
    const volumes = [];

    // Shuffle the array to randomize the order
    for (let i = stems.length - 1; i > 0; i--) {
        const j = Math.floor(Math.random() * (i + 1));
        [stems[i], stems[j]] = [stems[j], stems[i]];
    }
    const selectedStems = stems.slice(0, 3);

    // generate a bunch of random values
    for (let i = 32 - 1; i > 0; i--) {
        volumes.push(Math.random())
    }

    // print the value
    const lvg_values = {
        stems: selectedStems,
        volumes: volumes
    }

    console.log(lvg_values);

    const xhr = new XMLHttpRequest();
    xhr.open("POST", "http://localhost:5000/test_database", true);
    xhr.setRequestHeader("Content-Type", "application/json;charset=UTF-8");
    xhr.send(JSON.stringify({ lvg_values: lvg_values }));
    xhr.onload = function() {
        if (xhr.status === 200) {
            alert("Worked!");
        } else {
            alert("Something wrong!");
        }
    };
}
