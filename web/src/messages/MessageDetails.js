import React, { Component } from 'react';
import { Tab, TabList, TabPanel, Tabs } from 'react-tabs';
import Expander from '../expander/Expander';
import moment from 'moment';

const DATE_FORMAT = 'dddd D MMMM YYYY HH:mm:ss.SSS';

function createDefinition(name, value) {
  if (value === null) return null;

  return (
    <div className="property">
      <dt>{name}</dt>
      <dd>{value}</dd>
    </div>
  );
}

function mapHeaders(headers) {
  return headers.map((h, i) => (
    <div className="property" key={i}>
      <dt>{h.name}</dt>
      <dd>{h.value}</dd>
    </div>
  ));
}

function getUri(message) {
  if (!message || !message.request) return null;

  if (/^https?:\/\//i.test(message.request.uri)) {
    return message.request.uri;
  }

  const hostHeader = message.request.headers.filter(h => /^host$/i.test(h.name))[0];
  if (hostHeader) {
    return `https://${hostHeader.value}${message.request.uri}`;
  }

  return null;
}

class MessagesDetails extends Component {
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
      <div className="MessageDetails">
        <Tabs>
          <TabList>
            <Tab>Headers</Tab>
            <Tab>Raw</Tab>
          </TabList>

          <TabPanel>
            <div style={{display: this.props.message ? 'block' : 'none'}}>
              <Expander title="General">
                <dl className="properties">
                  {createDefinition('URL', getUri(message))}
                  {createDefinition('Method', message.request ? message.request.method : null)}
                  {createDefinition('Status Code', message.response ? `${message.response.statuscode} ${message.response.statustext}` : null)}
                  {createDefinition('Time', message.request ? moment(message.request.timestamp / 1000000).format(DATE_FORMAT) : null)}
                  {createDefinition('Duration', message.request && message.request ? `${Math.round((message.response.timestamp - message.request.timestamp) / 1000000)}ms` : null)}
                </dl>
              </Expander>
              <Expander title="Request">
                <dl className="properties">
                  {mapHeaders(message.request ? message.request.headers : [])}
                </dl>
              </Expander>
              <Expander title="Response">
                <dl className="properties">
                  {mapHeaders(message.response ? message.response.headers : [])}
                </dl>
              </Expander>
            </div>
          </TabPanel>
          <TabPanel>
            <Expander title="Request">
              <pre className="raw">{rawRequest}</pre>
            </Expander>
            <Expander title="Response">
              <pre className="raw">{rawResponse}</pre>
            </Expander>
          </TabPanel>
        </Tabs>
      </div>
    );
  }
}

export default MessagesDetails;
