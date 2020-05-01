window.onload = function () {
    let ws = new WebSocket("wss://" + document.location.host + "/entry");
    let chatbox = document.getElementById("chatbox");

    ws.onerror = function () {
        alert("WEBSOCKET SERVER DOESN'T WORK!");
    };

    ws.onmessage = function (e) {
        let msg = JSON.parse(e.data);
        chatbox.innerHTML += msg.author + ": " + msg.body + "<br>";
        chatbox.scrollTop = 9999;
    };

    //отправка сообщений на вебсокет
    let form = document.querySelector('form');
    form.onsubmit = function () {
        if (form[0].value !== '') {
            ws.send(JSON.stringify({body: form[0].value}));
        }
        form[0].value = '';
        return false;
    };
};
