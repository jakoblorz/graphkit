class App extends React.Component {

  static SVG = function (props) {
    const { innerHTML } = props;
    return <section className="svg" dangerouslySetInnerHTML={{
      __html: innerHTML,
    }}></section>
  };

  constructor(props) {
    super(props);
    this.state = {
      innerHTML: ""
    };
  }

  configureSocket(s) {
    const that = this;
    const onMessageReceived = function(e) {
      that.setState({
        innerHTML: e.data,
      })
    };
    const onSocketClosed = function() {
      that.componentDidMount();
    };

    s.addEventListener("message", onMessageReceived);
    s.addEventListener("close", onSocketClosed);
  }

  componentDidMount() {
    setTimeout(() => {
      const s = getSocket();
      if (s instanceof Promise) {
        s.then((s) => this.configureSocket(s));
      } else {
        this.configureSocket(s);
      }
    });
  }

  render() {
    const { innerHTML } = this.state;
    return <React.Fragment>
      <App.SVG innerHTML={innerHTML} />
    </React.Fragment>
  }
}

ReactDOM.render(<App />, document.getElementById("body"));