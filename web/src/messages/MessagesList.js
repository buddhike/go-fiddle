import React, { Component } from 'react';
import MessageRow from './MessageRow';
import './Messages.css';

class MessagesList extends Component {
  constructor(props) {
    super(props);

    this.handleSelect = this.handleSelect.bind(this);
  }

  handleSelect(message) {
    if (this.props.onSelect) {
      this.props.onSelect(message);
    }
  }

  render() {
    return (
      <table className="MessageList" cellspacing="0" cellpadding="0">
        <thead>
          <tr>
            <th className="ColumnMethod">Method</th>
            <th className="ColumnHost">Host</th>
            <th className="ColumnUri">Uri</th>
            <th className="ColumnStatus">Status</th>
          </tr>
        </thead>
        <tbody>
          {this.props.messages.map(m => <MessageRow key={m.id} message={m} onClick={this.handleSelect} />)}
        </tbody>
      </table>
    );
  }
}

export default MessagesList;
