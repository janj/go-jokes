class App extends React.Component {
    constructor(props) {
        super(props);
        console.log("WIN:", window.reactApp);
        if(window.reactApp) {
            this.state = { jokeId: window.reactApp.jokeId, jokes: [] };
        }
        else {
            this.state = { jokeId: undefined, jokes: [] };
        }
        this.serverRequest = this.serverRequest.bind(this);
    }

    render() {
        return (
            <div className="container">
                {this.state.jokes.map(function(joke, i) {
                    return <Joke key={i} joke={joke} />;
                })}
            </div>
        );
    }

    componentDidMount() {
        this.serverRequest();
    }

    serverRequest() {
        const withId = this.state.jokeId ? `/${this.state.jokeId}` : "";
        $.get(`http://localhost:8080/api/jokes${withId}`, res => {
            console.log("RESP:", res);
            this.setState({
                jokes: res
            });
        });
    }
}

class Joke extends React.Component {
    constructor(props) {
        super(props);
        console.log("Joke:", props);
    }

    render() {
        return (
            <div className="col-xs-6">
                <div className="panel panel-default">
                    <div className="panel-heading">
                        #{this.props.joke.Id}{" "}
                    </div>
                    <div className="panel-body">{this.props.joke.Joke}</div>
                </div>
            </div>
        )
    }
}

ReactDOM.render(<App />, document.getElementById('app'));
