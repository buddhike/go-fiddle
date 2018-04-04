import React, { Component } from 'react';
import './Messages.css';

class MessagesRow extends Component {
  constructor(props) {
    super(props);

    this.handleClick = this.handleClick.bind(this);
  }

  handleClick(e) {
    if (this.props.onClick) {
      this.props.onClick(this.props.message);
    }
  }

  render() {
    const { message } = this.props;
    const host = message.request.headers.filter(h => /^host$/i.test(h.name))[0].value;
    const { uri, method } = message.request;
    const { statuscode } = message.response;
    return (
      <tr onClick={this.handleClick}>
        <td className="ColumnMethod">{method}</td>
        <td className="ColumnHost">{host}</td>
        <td className="ColumnUri">{uri}</td>
        <td className="ColumnStatus">{statuscode}</td>
      </tr>
    );
  }
}

export default MessagesRow;
