import './App.css';
import './Resizer.css';

import React, { Component } from 'react';
import Header from './header/Header';
import MessageDetails from './messages/MessageDetails';
import MessageList from './messages/MessagesList';
import Sockette from 'sockette';
import SplitPane from 'react-split-pane';
import StatusPanel from './StatusPanel';
import config from './config';
import fetch from 'isomorphic-fetch';

const MIN_WIDTH = 640;

class App extends Component {
  constructor(props) {
    super(props);

    this.state = {
      messages: [],
      selectedMessageId: null,
      selectedMessage: null,
      gridWidth: MIN_WIDTH,
    };

    this.handleMessageSelect = this.handleMessageSelect.bind(this);
    this.handleData = this.handleData.bind(this);
    this.handleError = this.handleError.bind(this);
    this.handleStatusClose = this.handleStatusClose.bind(this);
    this.handleResize = this.handleResize.bind(this);
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

  handleError(err) {
    console.error(err);
    this.setState({
      status: {
        type: 'error',
        message: 'Oops, something went wrong',
      },
    });
  }

  handleStatusClose(err) {
    this.setState({
      status: null,
    });
  }

  async handleMessageSelect(message) {
    this.setState({
      selectedMessageId: message.id,
      selectedMessage: null,
    });

    try {
      const response = await fetch(`${config.restApi}messages/${message.id}`);
      const messageDetails = await response.json();

      this.setState({
        selectedMessage: messageDetails,
      });
    } catch (err) {
      this.handleError(err);
    }
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

  handleResize(size) {
    this.setState({gridWidth: size});
  }

  render() {
    return (
      <div className="App">
        <Header />
        <div className="container">
          <SplitPane split="vertical" minSize={MIN_WIDTH} defaultSize={this.state.gridWidth} onChange={this.handleResize}>
            <div className="list-panel">
              <MessageList messages={this.state.messages} activeMessageId={this.state.selectedMessageId} onSelect={this.handleMessageSelect} width={this.state.gridWidth} />
            </div>
            <div className="details-panel">
              <MessageDetails message={this.state.selectedMessage} />
            </div>
          </SplitPane>
        </div>
        { this.state.status ?
          <StatusPanel type={this.state.status.type} onDismiss={this.handleStatusClose}>{this.state.status.message}</StatusPanel> :
          null
        }
      </div>
    );
  }
}

export default App;
