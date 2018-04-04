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
    const { message, active } = this.props;
    const host = message.request.headers.filter(h => /^host$/i.test(h.name))[0].value;
    const { uri, method } = message.request;
    const { statuscode } = message.response;
    return (
      <tr className={active ? 'active' : ''} onClick={this.handleClick}>
        <td className="col-method">{method}</td>
        <td className="col-host">{host}</td>
        <td className="col-uri">{uri}</td>
        <td className="col-status">{statuscode}</td>
      </tr>
    );
  }
}

export default MessagesRow;
