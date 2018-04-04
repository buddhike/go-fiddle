import React, { Component } from 'react';
import fetch from 'isomorphic-fetch';
import MessageList from './messages/MessagesList';
import MessageDetails from './messages/MessageDetails';
import './App.css';

class App extends Component {
  constructor(props) {
    super(props)

    this.state = {
      messages: [],
      selectedMessage: null,
    };

    this.handleMessageSelect = this.handleMessageSelect.bind(this);
  }

  async componentDidMount() {
    const response = await fetch('http://localhost:8888/messages');
    const messages = await response.json();

    this.setState({messages});
  }

  handleMessageSelect(message) {
    this.setState({
      selectedMessage: message,
    });
  }

  render() {
    return (
      <div className="App">
        <div className="list-panel">
          <MessageList messages={this.state.messages} activeMessage={this.state.selectedMessage} onSelect={this.handleMessageSelect} />
        </div>
        <div className="details-panel">
          <MessageDetails message={this.state.selectedMessage} />
        </div>
      </div>
    );
  }
}

export default App;
