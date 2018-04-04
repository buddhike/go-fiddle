import React, { Component } from 'react';

class MessagesDetails extends Component {
  render() {
    const { message } = this.props;

    const rawRequest = message ? [
      `${message.request.method} ${message.request.uri} ${message.request.version}`,
      ...message.request.headers.map(h => `${h.name}: ${h.value}`),
      '',
      message.request.body,
    ].join('\r\n') : '';
    const rawResponse = message && message.response ? [
      `${message.response.version} ${message.response.statuscode} ${message.response.statustext}`,
      ...message.response.headers.map(h => `${h.name}: ${h.value}`),
      '',
      message.response.body,
    ].join('\r\n') : '';

    return (
      <div>
        <div>
          <div>Request:</div>
          <pre>{rawRequest}</pre>
        </div>
        <div>
          <div>Response:</div>
          <pre>{rawResponse}</pre>
        </div>
      </div>
    );
  }
}

export default MessagesDetails;
