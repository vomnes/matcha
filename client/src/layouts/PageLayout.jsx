import React, { Component } from 'react';
import './PageLayout.css';

class GeneralLayout extends Component {
  render() {
    const {
      children,
    } = this.props;

    return (
      <div className="general-layout">
        <div className="nav-top">
        </div>
        <div className="content">
          { children }
        </div>
      </div>
    );
  }
}
