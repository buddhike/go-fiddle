import React, { Component } from 'react';
import MessageRow from './MessageRow';
import './Messages.css';

class MessagesList extends Component {
  render() {
    return (
      <table className="MessageList">
        <thead>
          <tr>
            <th>method</th>
            <th>uri</th>
            <th>host</th>
            <th>status</th>
          </tr>
        </thead>
        <tbody>
          {this.props.messages.map(m => <MessageRow key={m.id} message={m} />)}
        </tbody>
      </table>
    );
  }
}

export default MessagesList;
