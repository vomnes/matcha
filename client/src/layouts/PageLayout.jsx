import React, { Component } from 'react';
import './PageLayout.css';

class PageLayout extends Component {
  render() {
    const {
      children,
    } = this.props;

    return (
      <div className="general-layout">
        <div className="header">
          <div className="nav-top">
            <a href='/home' className="logo"><span>Matcha</span></a>
            <div className="header-left-side">
              <a href='/' className="logout"><span>Browsing</span></a>
              <a href='/' className="logout"><span>Matches</span></a>
              <a href='/myprofile' className="logout"><span>My Profile</span></a>
              <a href='/logout' className="logout"><span>Logout</span></a>
            </div>
          </div>
        </div>
        <div className="content">
          { children }
        </div>
      </div>
    );
  }
}

export default PageLayout;
