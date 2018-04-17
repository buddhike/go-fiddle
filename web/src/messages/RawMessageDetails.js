import React, { Component } from 'react';
import Expander from '../expander/Expander';

class RawMessageDetails extends Component {
  render() {
    let { message } = this.props;
    if (!message) message = {};

    const rawRequest = message.request ? [
      `${message.request.method} ${message.request.uri} ${message.request.version}`,
      ...message.request.headers.map(h => `${h.name}: ${h.value}`),
      '',
      message.request.body,
    ].join('\r\n') : '';
    const rawResponse = message.response ? [
      `${message.response.version} ${message.response.statuscode} ${message.response.statustext}`,
      ...message.response.headers.map(h => `${h.name}: ${h.value}`),
      '',
      message.response.body,
    ].join('\r\n') : '';

    return (
      <div>
        <Expander title="Request">
          <pre className="raw">{rawRequest}</pre>
        </Expander>
        <Expander title="Response">
          <pre className="raw">{rawResponse}</pre>
        </Expander>
      </div>
    );
  }
}

export default RawMessageDetails;
