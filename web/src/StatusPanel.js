import './StatusPanel.css';

import React, { Component } from 'react';

class StatusPanel extends Component {
  constructor(props) {
    super(props);

    this.handleHide = this.handleHide.bind(this);
  }

  handleHide() {
    if (this.props.onDismiss) {
      this.props.onDismiss();
    }
  }

  render() {
    const { type } = this.props;
    return (
      <div className={['StatusPanel', type].join(' ')}>
        <button className="close" onClick={this.handleHide}>
          <svg width={20} height={20}>
            <path className="path" fill="none" strokeWidth={2} d="M4,4 L16,16 M4,16 L16,4" />
          </svg>
        </button>
        <div className="content">
          {this.props.children}
        </div>
      </div>
    );
  }
}

export default StatusPanel;
