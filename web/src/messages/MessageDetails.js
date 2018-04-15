import React, { Component } from 'react';
import moment from 'moment';

const DATE_FORMAT = 'dddd D MMMM YYYY HH:mm:ss.SSS';

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
      <div className="MessageDetails">
        <div className="request-details">
          <dt>Request</dt>
          <pre>{rawRequest}</pre>
          <div className="time">{message && message.request ? moment(message.request.timestamp / 1000000).format(DATE_FORMAT) : ''}</div>
        </div>
        <div className="response-details">
          <dt>Response</dt>
          <pre>{rawResponse}</pre>
          <div className="time">{message && message.response ? moment(message.response.timestamp / 1000000).format(DATE_FORMAT) : ''}</div>
        </div>
      </div>
    );
  }
}

export default MessagesDetails;
