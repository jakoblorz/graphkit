function configureSocket(s) {
    const onMessageReceived = function(e) {
        document.getElementById("body").innerHTML = Viz(e.data, "svg");
    };
    const onSocketClosed = function() {
        componentDidMount();
    };

    s.addEventListener("message", onMessageReceived);
    s.addEventListener("close", onSocketClosed);
}

function componentDidMount() {
    setTimeout(() => {
        const s = getSocket();
        if (s instanceof Promise) {
            s.then((s) => configureSocket(s));
        } else {
            configureSocket(s);
        }
    })
};

(function() {
    componentDidMount();
})();