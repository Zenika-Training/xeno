function sendAliens() {
    var alienNumber = 10;
    $.get("/sendAliens?alien=" + alienNumber).done(function () {
        var msg = `⚠️ Sent ${alienNumber} 👽 aliens in !`;
        $.toast({
            text: msg,
            hideAfter: 1000
        })
    });
}

function sendMarines() {
    var marineNumber = 10;
    $.get("/sendMarines?marine=" + marineNumber).done(function () {
        var msg = `⚠️ Sent ${marineNumber} 👮 marines to the rescue !`;
        console.log(msg)
        $.toast({
            text: msg,
            hideAfter: 1000
        })
    });
}

function resetSimulation() {
    $.get("/resetSimulation").done(function () {
        var msg = `⚠️ Simulation reseted`;
        console.log(msg)
        $.toast({
            text: msg,
            hideAfter: 1000
        })
    });
}