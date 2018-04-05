import React, { Component } from 'react';
import fetch from 'isomorphic-fetch';
import Websocket from 'react-websocket';
import MessageList from './messages/MessagesList';
import MessageDetails from './messages/MessageDetails';
import config from './config';

import './App.css';

class App extends Component {
  constructor(props) {
    super(props)

    this.state = {
      messages: [],
      selectedMessage: null,
    };

    this.handleMessageSelect = this.handleMessageSelect.bind(this);
    this.handleData = this.handleData.bind(this);
  }

  componentDidMount() {
    return this.refreshData();
  }

  async refreshData() {
    const response = await fetch(`${config.restApi}messages`);
    const messages = await response.json();

    this.setState({messages});
  }

  async handleMessageSelect(message) {
    const response = await fetch(`${config.restApi}messages/${message.id}`);
    const messageDetails = await response.json();
    this.setState({
      selectedMessage: messageDetails,
    });
  }

  handleData(data) {
    data = JSON.parse(data);
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
          <MessageList messages={this.state.messages} activeMessage={this.state.selectedMessage} onSelect={this.handleMessageSelect} />
        </div>
        <div className="details-panel">
          <MessageDetails message={this.state.selectedMessage} />
        </div>
        <Websocket url={config.websocket} onMessage={this.handleData} />
      </div>
    );
  }
}

export default App;
