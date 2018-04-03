import React, { Component } from 'react';

class MessagesRow extends Component {
  render() {
    const { message } = this.props;
    const host = message.request.headers.filter(h => /^host$/i.test(h.name))[0].value;
    const { uri, method } = message.request;
    const { statuscode } = message.response;
    return (
      <tr>
        <td>{method}</td>
        <td>{uri}</td>
        <td>{host}</td>
        <td>{statuscode}</td>
      </tr>
    );
  }
}

export default MessagesRow;
