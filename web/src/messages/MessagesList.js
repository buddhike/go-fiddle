import './Messages.css';
import React, { Component } from 'react';
import MessageRow from './MessageRow';

class MessagesList extends Component {
  constructor(props) {
    super(props);

    this.handleSelect = this.handleSelect.bind(this);
    this.handleKeyDown = this.handleKeyDown.bind(this);
  }

  handleSelect(message) {
    if (this.props.onSelect) {
      this.props.onSelect(message);
    }
  }

  handleKeyDown(e) {
    const { messages, activeMessageId } = this.props;

    if (!messages || !messages.length) return;
    if (!['ArrowUp', 'ArrowDown', 'Home', 'End', 'PageUp', 'PageDown'].includes(e.key)) return;

    const height = this.refs.rows.clientHeight - this.refs.header.clientHeight;
    const rowHeight = height / messages.length;
    const pageSize = Math.floor((this.refs.container.parentElement.clientHeight - this.refs.header.clientHeight) / rowHeight);

    console.log('pageSize', pageSize, 'header', this.refs.header, this.refs.header.clientHeight, 'rows', this.refs.rows, this.refs.rows.clientHeight);

    const selectedIndex = messages.findIndex(m => m.id === activeMessageId);
    let newIndex = selectedIndex;

    if (e.key === 'ArrowUp') {
      newIndex = Math.max(selectedIndex - 1, 0);
    } else if (e.key === 'ArrowDown') {
      newIndex = Math.min(selectedIndex + 1, messages.length - 1);
    } else if (e.key === 'PageUp') {
      newIndex = Math.max(selectedIndex - pageSize, 0);
    } else if (e.key === 'PageDown') {
      newIndex = Math.min(selectedIndex + pageSize, messages.length - 1);
    } else if (e.key === 'Home') {
      newIndex = 0;
    } else if (e.key === 'End') {
      newIndex = messages.length - 1;
    }

    e.preventDefault();
    e.stopPropagation();

    const message = messages[newIndex];
    this.handleSelect(message);

    console.log('Key down', e, e.keyCode, e.key);
  }

  render() {
    const { activeMessageId } = this.props;

    return (
      <div ref="container" className="MessageList">
        <table ref="header" className="head" cellSpacing="0" cellPadding="0">
          <thead>
            <tr>
              <th className="col-time">Time</th>
              <th className="col-method">Method</th>
              <th className="col-uri">Uri</th>
              <th className="col-status">Status</th>
            </tr>
          </thead>
        </table>
        <table ref="rows" className="body" cellSpacing="0" cellPadding="0" onKeyDown={this.handleKeyDown}>
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
