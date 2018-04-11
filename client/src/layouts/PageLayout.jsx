import React, { Component } from 'react';
import './PageLayout.css';

const NotifList = (props) => {
  return (
    <div id="notif">
      <div id="notif-arrow"></div>
      <div className="notif-element">
        <div className="picture-notif-background" style={{ backgroundImage: "url(https://images.unsplash.com/photo-1520512202623-51c5c53957df?h=1000&q=10)" }}></div>
        <div className="white-notif-background"></div>
      </div>
      <div className="notif-element">
        <div className="picture-notif-background" style={{ backgroundImage: "url(https://images.unsplash.com/photo-1522234811749-abc512463137?h=1000&q=10)" }}></div>
        <div className="white-notif-background"></div>
      </div>
    </div>
  )
}

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
              <a href='/browsing' className="logout"><span>Browsing</span></a>
              <a href='/matches' className="logout"><span>Matches</span></a>
              <a href='/profile' className="logout"><span>My Profile</span></a>
              <span className="logout" style={{ cursor: "pointer" }}>Notifications</span>
              <NotifList />
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
