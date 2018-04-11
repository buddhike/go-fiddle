import React, { Component } from 'react';
import fetch from 'isomorphic-fetch';
import Sockette from 'sockette';
import MessageList from './messages/MessagesList';
import MessageDetails from './messages/MessageDetails';
import config from './config';

import './App.css';

class App extends Component {
  constructor(props) {
    super(props)

    this.state = {
      messages: [],
      selectedMessageId: null,
      selectedMessage: null,
    };

    this.handleMessageSelect = this.handleMessageSelect.bind(this);
    this.handleData = this.handleData.bind(this);
  }

  componentDidMount() {
    this.websocket = new Sockette(config.websocket, {
      timeout: 5000,
      maxAttempts: 10,
      onmessage: e => {
        this.handleData(e);
      },
    });

    return this.refreshData();
  }

  async refreshData() {
    const response = await fetch(`${config.restApi}messages`);
    const messages = await response.json();

    this.setState({messages});
  }

  async handleMessageSelect(message) {
    this.setState({
      selectedMessageId: message.id,
      selectedMessage: null,
    });
    const response = await fetch(`${config.restApi}messages/${message.id}`);
    const messageDetails = await response.json();

    this.setState({
      selectedMessage: messageDetails,
    });
  }

  handleData(e) {
    const data = JSON.parse(e.data);
    const messages = this.state.messages.slice();
    const index = messages.findIndex(m => m.id === data.id);

    if (index === -1) {
      messages.push(data);
    } else {
      messages[index] = data;
    }

    this.setState({messages});
  }

  render() {
    return (
      <div className="App">
        <div className="list-panel">
          <MessageList messages={this.state.messages} activeMessageId={this.state.selectedMessageId} onSelect={this.handleMessageSelect} />
        </div>
        <div className="details-panel">
          <MessageDetails message={this.state.selectedMessage} />
        </div>
      </div>
    );
  }
}

export default App;
