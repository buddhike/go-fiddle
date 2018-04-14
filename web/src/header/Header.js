import React, { Component } from 'react';
import './Header.css';
import config from '../config';

class Header extends Component {
  constructor(props) {
    super(props);

    this.state = {
      showCertificatesPanel: false,
    };

    this.handleSelect = this.handleSelect.bind(this);
  }

  handleSelect(message) {
    if (this.props.onSelect) {
      this.props.onSelect(message);
    }
  }

  render() {
    const { activeMessageId } = this.props;

    return (
      <div className="Header">
        <div className="brand">GoFiddle</div>
        <div className="tools">
          <a href={`${config.restApi}certificate`}>Download certificate</a>
        </div>
      </div>
    );
  }
}

export default Header;
