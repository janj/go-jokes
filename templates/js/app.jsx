class App extends React.Component {
    constructor(props) {
        super(props);
        this.state = {
            jokes: []
        };

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
        $.get("http://localhost:8080/api/jokes", res => {
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
        console.log("PROPS:", props);
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
