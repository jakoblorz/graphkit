let socket;

function setSocket(s) {
    socket = s;
}

function getSocket() {
    return socket;
}

(function connect() {
    const url = window.location.href.replace(
      window.location.protocol,
      "ws:",
    );

    return new Promise((resolve) => {
        const s = new WebSocket(`${url}${
        window.location.pathname == "/" ? "ws" : "/ws"
        }`);

        s.addEventListener("close", function() {
            setSocket(connect());
        });
        setSocket(s);
        resolve(s);
    });
})()