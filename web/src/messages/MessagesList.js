import './Messages.css';
import React, { Component } from 'react';
import MessageRow from './MessageRow';

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
    const { activeMessageId } = this.props;

    return (
      <div className="MessageList">
        <table className="head" cellSpacing="0" cellPadding="0">
          <thead>
            <tr>
              <th className="col-time">Time</th>
              <th className="col-method">Method</th>
              <th className="col-uri">Uri</th>
              <th className="col-status">Status</th>
            </tr>
          </thead>
        </table>
        <table className="body" cellSpacing="0" cellPadding="0">
          <tbody>
            {this.props.messages.map(m => (
              <MessageRow key={m.id}
                message={m}
                active={m.id === activeMessageId}
                onClick={this.handleSelect}
              />
            ))}
          </tbody>
        </table>
      </div>
    );
  }
}

export default MessagesList;
