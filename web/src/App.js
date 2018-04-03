import React, { Component } from 'react';
import fetch from 'isomorphic-fetch';
import MessageList from './messages/MessagesList';
import './App.css';

class App extends Component {
  constructor(props) {
    super(props)
    this.state = {
      messages: []
    };
  }

  async componentDidMount() {
    const response = await fetch('http://localhost:8888/messages');
    const messages = await response.json();

    this.setState({messages});
  }

  render() {
    return (
      <div className="App">
        <MessageList messages={this.state.messages} />
      </div>
    );
  }
}

export default App;
