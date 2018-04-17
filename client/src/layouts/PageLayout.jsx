import React, { Component } from 'react';
import './PageLayout.css';
import ViewIcon from '../design/icons/eye.svg'
import HeartIcon from '../design/icons/like.svg'
import BrokenHeartIcon from '../design/icons/broken-heart-divided-in-two-parts.svg'
import MatchIcon from '../design/icons/flame.svg'
import MessageIcon from '../design/icons/conversation-speech-bubbles.svg'
import ViewIconRed from '../design/icons/eye-red.svg'
import HeartIconRed from '../design/icons/like-red.svg'
import BrokenHeartIconRed from '../design/icons/broken-heart-divided-in-two-parts-red.svg'
import MatchIconRed from '../design/icons/flame-red.svg'
import MessageIconRed from '../design/icons/conversation-speech-bubbles-red.svg'
import SearchIcon from '../design/icons/magnifier.svg'
import NotificationIcon from '../design/icons/notifications-button.svg'
import ProfileIcon from '../design/icons/man-user.svg'
import LogoutIcon from '../design/icons/logout.svg'
import api from '../library/api'

var moment = require('moment');

const NotifElement = (props) => {
  var icon, content;
  switch (props.type) {
    case "view":
        if (props.new) {
          icon = ViewIconRed;
        } else {
          icon = ViewIcon;
        }
        content = `${props.fullname} has viewed your profile`;
        break;
    case "like":
        if (props.new) {
          icon = HeartIconRed;
        } else {
          icon = HeartIcon;
        }
        content = `${props.fullname} has liked your profile`;
        break;
    case "unmatch":
        if (props.new) {
          icon = BrokenHeartIconRed;
        } else {
          icon = BrokenHeartIcon;
        }
        content = `${props.fullname} has broken your match`;
        break;
    case "match":
        if (props.new) {
          icon = MatchIconRed;
        } else {
          icon = MatchIcon;
        }
        content = `You have a match with ${props.fullname}`;
        break;
    case "message":
        if (props.new) {
          icon = MessageIconRed;
        } else {
          icon = MessageIcon;
        }
        content = `${props.fullname} has sent you a message`;
        break;
    default:
        icon = ViewIcon;
        content = '';
  }
  return (
    <div className="notif-element">
      <div className="picture-notif-background" style={{ backgroundImage: `url(${props.picture})` }}></div>
      <div className="white-notif-background"></div>
      <div className="notif-logo">
        <img alt="View profile" src={icon} style={{ fill: "#434343", width: "100%", opacity: "0.7" }}/>
      </div>
      <span className="notif-text">{content}<br/>{moment(props.date).format("M/D/YYYY - HH:mm")}</span>
    </div>
  )
}

const NotifList = (props) => {
  var listNotifications = [];
  var index = 0;
  props.notifications.forEach((notification) => {
    listNotifications.push(
      <NotifElement
        type={notification.type}
        fullname={`${notification.firstname} ${notification.lastname}`}
        picture={notification.user_picture_url}
        date={notification.date}
        key={index}
        new={notification.new}
      />
    );
    index += 1;
  });
  if (props.isOpen) {
    return (
      <div id="notif">
        {listNotifications}
      </div>
    )
  }
  return null;
}

const GetListNotifications = async (updateState) => {
  let res = await api.getNotifications();
  if (res) {
    const response = await res.json();
    if (res.status >= 500) {
      throw new Error("Bad response from server - GetListNotifications has failed");
    } else if (res.status >= 400) {
      console.log(response.error);
    } else {
      if (response.data === "No notifications") {
        updateState("notifications", []);
      } else {
        console.log(response);
        updateState("notifications", response);
      }
    }
  }
}

class PageLayout extends Component {
  constructor(props) {
    super(props);
    this.state = {
      notificationsOpen: false,
      newNotification: true,
      loggedProfileData: {},
      notifications: [],
      isMobile: window.innerWidth > 505,
    }
    this.updateState = this.updateState.bind(this);
  }
  updateState(key, content) {
    this.setState({
      [key]: content,
    });
  }
  handleNotifications(updateState) {
    if (!this.state.notificationsOpen) {
      GetListNotifications(this.updateState);
    }
    this.updateState("notificationsOpen", !this.state.notificationsOpen);
    var profile = this.state.loggedProfileData;
    profile["total_new_notifications"] = 0;
    this.updateState("loggedProfileData", profile);
  }
  updateHeader = () => {
    this.updateState("isMobile", window.innerWidth > 505);
  }
  componentDidMount() {
    window.addEventListener("resize", this.updateHeader);
    this.props.myProfileData.then((data) => {
      this.updateState("loggedProfileData", data);
    });
  }
  handleWebsocket = (msg) => {
    var profile = this.state.loggedProfileData;
    if (msg.event === "message") {
      profile["total_new_messages"] += 1;
    } else if (msg.event !== "isTyping" && msg.event !== "login" && msg.event !== "logout") {
      profile["total_new_notifications"] += 1;
      this.updateState("notificationsOpen", false);
    }
    this.updateState("loggedProfileData", profile);
  }
  render() {
    const {
      children,
    } = this.props;
    this.props.wsConn.onmessage = (e) => {
      try {
        var msg = JSON.parse(e.data);
      } catch (e) {
        return;
      }
      this.handleWebsocket(msg);
    }
    return (
      <div className="general-layout">
        <div className="header">
          {this.state.isMobile ? (
            <div className="nav-top">
              <a href='/home' className="logo"><span>Matcha</span></a>
              <div className="header-left-side">
                <a href='/browsing'><span>Browsing</span></a>
                <a href='/matches'><span>Matches</span></a>
                {this.state.loggedProfileData && this.state.loggedProfileData.total_new_messages && !window.location.href.endsWith("matches") ? (<span className="top-notif red-cercle-notif" id="matches-notif">{this.state.loggedProfileData && this.state.loggedProfileData.total_new_messages}</span>) : null}
                <a href='/profile'><span>My Profile</span></a>
                <span id="notif-btn" onClick={() => this.handleNotifications(this.updateState)}>Notifications</span>
                {this.state.loggedProfileData && this.state.loggedProfileData.total_new_notifications ? (<span className="top-notif red-cercle-notif" id="true-notif">{this.state.loggedProfileData && this.state.loggedProfileData.total_new_notifications}</span>) : null}
                <a href='/logout'><span>Logout</span></a>
              </div>
              {this.state.notificationsOpen ? (<div id="notif-arrow"></div>) : null}
            </div>
          ) : (
            <div className="nav-top-responsive">
              <a href='/home' className="logo"><span>Matcha</span></a>
              <a href='/browsing'><img alt="Browsing button" title="browsing" src={SearchIcon} className="header-button" style={{ left: "25vw" }}/></a>
              <a href='/matches'><span><img alt="Matches button" title="matches" src={MatchIcon} className="header-button" style={{ left: "40vw" }}/></span></a>
              {this.state.loggedProfileData && this.state.loggedProfileData.total_new_messages && !window.location.href.endsWith("matches") ? (<span className="top-notif red-cercle-notif" id="matches-notif-responsive">{this.state.loggedProfileData && this.state.loggedProfileData.total_new_messages}</span>) : null}
              <a href='/profile'><span><img alt="My profile button" title="profile" src={ProfileIcon} className="header-button" style={{ left: "55vw" }}/></span></a>
              <img alt="Notification list button" title="notifications" src={NotificationIcon} className="header-button" onClick={() => this.handleNotifications(this.updateState)} style={{ left: "70vw" }}/>
              {this.state.loggedProfileData && this.state.loggedProfileData.total_new_notifications ? (<span className="top-notif red-cercle-notif" id="true-notif-responsive">{this.state.loggedProfileData && this.state.loggedProfileData.total_new_notifications}</span>) : null}
              <a href='/logout'><img alt="Logout button" title="logout" src={LogoutIcon} className="header-button" style={{ right: "5vw" }}/></a>
              {this.state.notificationsOpen ? (<div id="notif-arrow"></div>) : null}
            </div>
          )}
        </div>
        <div className="content">
          <NotifList isOpen={this.state.notificationsOpen} notifications={this.state.notifications}/>
          { children }
        </div>
      </div>
    );
  }
}

export default PageLayout;
