import React, { Component } from 'react';
import './Messages.css';

class MessagesRow extends Component {
  constructor(props) {
    super(props);

    this.handleClick = this.handleClick.bind(this);
  }

  handleClick() {
    if (this.props.onClick) {
      this.props.onClick(this.props.message);
    }
  }

  render() {
    const { message, active } = this.props;
    const { uri, method, statuscode } = message;

    return (
      <tr className={active ? 'active' : ''} onClick={this.handleClick}>
        <td className="col-method">{method}</td>
        <td className="col-uri" title={uri}>{uri}</td>
        <td className="col-status">{statuscode || '-'}</td>
      </tr>
    );
  }
}

export default MessagesRow;
