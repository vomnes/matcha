import React, { Component } from 'react';
import './PageLayout.css';
import ViewIcon from '../design/icons/eye.svg'
import HeartIcon from '../design/icons/like.svg'
import BrokenHeartIcon from '../design/icons/broken-heart-divided-in-two-parts.svg'
import MatchIcon from '../design/icons/flame.svg'
import MessageIcon from '../design/icons/conversation-speech-bubbles.svg'

const NotifElement = (props) => {
  var icon, content;
  switch (props.type) {
    case "view":
        icon = ViewIcon;
        content = `${props.fullname} has viewed your profile`;
        break;
    case "like":
        icon = HeartIcon;
        content = `${props.fullname} has liked your profile`;
        break;
    case "unlike":
        icon = BrokenHeartIcon;
        content = `${props.fullname} has unliked your profile`;
        break;
    case "match":
        icon = MatchIcon;
        content = `You have a match with ${props.fullname}`;
        break;
    case "message":
        icon = MessageIcon;
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
      <div className="notif-logo" title="See more">
        <img alt="View profile" src={icon} style={{ width: "100%", opacity: "0.7" }}/>
      </div>
      <span className="notif-text">{content}<br/>{props.date}</span>
    </div>
  )
}

const NotifList = (props) => {
  if (props.isOpen) {
    return (
      <div id="notif">
        <NotifElement type="unlike" fullname="Valentin Omnes" picture="https://images.unsplash.com/photo-1520512202623-51c5c53957df?h=1000&q=10" date={"12/12/12 - 12:42"}/>
        <NotifElement type="message" fullname="Valentin Omnes" picture="https://images.unsplash.com/photo-1520512202623-51c5c53957df?h=1000&q=10" date={"12/12/12 - 12:42"}/>
        <NotifElement type="match" fullname="Valentin Omnes" picture="https://images.unsplash.com/photo-1520512202623-51c5c53957df?h=1000&q=10" date={"12/12/12 - 12:42"}/>
        <NotifElement type="like" fullname="Valentin Omnes" picture="https://images.unsplash.com/photo-1520512202623-51c5c53957df?h=1000&q=10" date={"12/12/12 - 12:42"}/>
        <NotifElement type="view" fullname="Valentin Omnes" picture="https://images.unsplash.com/photo-1520512202623-51c5c53957df?h=1000&q=10" date={"12/12/12 - 12:42"}/>
      </div>
    )
  }
  return null;
}

// -> visit
// -> like
// -> unlike
// -> match
// -> message

class PageLayout extends Component {
  constructor() {
    super();
    this.state = {
      notificationsOpen: false,
      newNotification: true,
    }
  }
  render() {
    const {
      children,
    } = this.props;
    var notifStyle = { cursor: "pointer" };
    notifStyle["border"] = this.state.newNotification ? "2px solid #e63946" : "2px solid white";
    return (
      <div className="general-layout">
        <div className="header">
          <div className="nav-top">
            <a href='/home' className="logo"><span>Matcha</span></a>
            <div className="header-left-side">
              <a href='/browsing' className="logout"><span>Browsing</span></a>
              <a href='/matches' className="logout"><span>Matches</span></a>
              <a href='/profile' className="logout"><span>My Profile</span></a>
              <span id="notif-btn" className="logout" style={notifStyle}
                onClick={() => this.setState({notificationsOpen: !this.state.notificationsOpen})}
              >Notifications</span>
              <NotifList isOpen={this.state.notificationsOpen}/>
              <a href='/logout' className="logout"><span>Logout</span></a>
            </div>
            {this.state.notificationsOpen ? (<div id="notif-arrow"></div>) : null}
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
