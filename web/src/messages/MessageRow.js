import React, { Component } from 'react';
import './Messages.css';

const hostFromUriExpression = /^(\w+):\/\/([^\/]+)/i;
function getHeaderValue(headers, headerName) {
  const expression = new RegExp(`^${headerName}$`, "i");
  const header = headers.filter(h => expression.test(h.name))[0];
  if (header) {
    return header.value;
  }
  return null;
}

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
    const { uri, method } = message.request;
    const host = getHeaderValue(message.request.headers, "host") || (hostFromUriExpression.exec(uri) || [])[2];
    const { statuscode } = message.response || {};
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
