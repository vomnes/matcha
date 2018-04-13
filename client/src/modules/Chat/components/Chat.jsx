import React, { Component } from 'react';
import './Chat.css';
import SendButton from '../../../design/icons/send-button.svg';
import api from '../../../library/api'

var moment = require('moment');
// import encode from '../../../library/utils/encode.js'

// var conn;
// var msg = document.getElementById("msg");
// document.getElementById("form").onsubmit = function () {
//     if (!conn) {
//         return false;
//     }
//     if (!msg.value) {
//         return false;
//     }
//     conn.send(msg.value);
//     msg.value = "";
//     return false;
// };
// if (window["WebSocket"]) {
//     var room = encode.objectToBase64({username1: "voluptas_atque_etpawSY", username2: "Elsa"});
//     conn = new WebSocket(`ws://localhost:8081/ws/chat/${room}`, localStorage.getItem('matcha_token'));
//     conn.onclose = function (evt) {
//         var item = document.createElement("div");
//         item.innerHTML = "<b>Connection closed.</b>";
//         document.getElementById("log").appendChild(item);
//     };
//     conn.onmessage = function (evt) {
//         var messages = evt.data.split('\n');
//         for (var i = 0; i < messages.length; i++) {
//             var item = document.createElement("div");
//             item.innerText = messages[i];
//             document.getElementById("log").appendChild(item);
//         }
//     };
// }
/*<span>Nice chat isn't it ?</span>
<div id="log"></div>*/

const MatchItem = (props) => {
  const selectedStyle = {
    top: "-5px",
    left: "-5px",
    border: "none",
    boxShadow: "0 3px 6px rgba(0,0,0,0.16), 0 3px 6px rgba(0,0,0,0.23)",
    zIndex: "2",
  }
  if (props && props.time) {
    var date = moment(props.time)
    var formatedDate;
    if (moment().diff(date, 'days') > 1) {
      if (moment().diff(date, 'years') > 1) {
        formatedDate = date.format("DD MMM, YYYY")
      } else {
        formatedDate = date.format("DD MMM")
      }
    } else {
      formatedDate = date.startOf("minute").fromNow();
    }
  }
  return (
    <div>
      <div className="match-element" id={props.username} style={props.selectedProfile === props.username ? selectedStyle : null } onClick={() => props.updateSelectedProfile(props.username)}>
        {props.isOnline ? (<span className="online-dot" title={`${props.name} is online`}>&bull;</span>) : null}
        <div className="picture-list">
          <a href={`/profile/${props.username}` + (props.optionsBase64 ? '/' + props.optionsBase64 : '')} title={`Click to see ${props.name}'s profile`}>
            <div className="picture-list-background" style={{ backgroundImage: "url(" + props.picture + ")" }}></div>
          </a>
        </div>
        <span className="match-element-list">{props.name}</span>
        {props.selectedProfile === props.username ? null : (
          <div>
            <span className="match-element-time" title={`Last message sent - ${formatedDate}`}>{formatedDate}</span>
            <span className="match-element-lastmsg">{props.lastmsg}</span>
          </div>
        )}
      </div>
      <div className="limit" style={{ width: "94.25%", margin: "2.5px 0px 2.5px" }}></div>
    </div>
  )
}

const MsgItem = (props) => {
  var pictureStyle = null;
  var msgHeaderStyle = null;
  if (props.side === "left") {
    pictureStyle = { left: "-22px" };
    msgHeaderStyle = { left: "25px" };
  } else {
    pictureStyle = { right: "-22px" };
    msgHeaderStyle = { right: "30px" };
  }
  return (
    <div className="msg-element">
      <div className="msg">
        <div className="picture-msg-list" style={pictureStyle}>
          <div className="picture-list-background" style={{ backgroundImage: "url(" + props.picture + ")" }}></div>
        </div>
        <span className="msg-header" style={msgHeaderStyle}>{props.firstname} {props.lastname} - {moment(props.received_at).format('LT')}</span>
        <div className="msg-content" style={{ textAlign: props.side }}>
          <span>{props.content}</span>
        </div>
      </div>
    </div>
  )
}

const MsgArea = (props) => {
  var dataMessage;
  var msgScroll;
  var messages = [];
  var index = 0;
  props.listMsg.forEach((message) => {
    messages.push(
      <MsgItem
        key={index} index={index}
        side={props.myusername === message.username ? "right" : "left"}
        firstname={message.firstname}
        lastname={message.lastname}
        picture={message.picture_url}
        content={message.content}
        received_at={message.received_at}
      />
    );
    index += 1;
  });
  if (props.listMsg.length > 0) {
    msgScroll = { overflowY: "scroll", overflowX: "hidden" };
    dataMessage = (
      <div>
        {/* {window.onload = () => {document.getElementById("list-msg-area").scrollTop = document.getElementById("list-msg-area").scrollHeight;}} */}
        <div id="msg-array">
          {/* <div className="day">
            <span>Monday 15 December</span>
          </div> */}

          {/* <div className="day">
            <span>Today</span>
          </div> */}
          {messages}
        </div>
      </div>
    )
  } else if (props.selectedProfileData.username !== undefined) {
    dataMessage = (
      <div id="no-msg">
        <a href={`/profile/${props.selectedProfileData.username}`} title={`Click to see ${props.selectedProfileData.firstname}'s profile`}>
          <div id="picture-no-msg" style={null}>
            <div className="picture-list-background" style={{ backgroundImage: `url(${props.selectedProfileData.picture_url})` }}></div>
          </div>
        </a>
        <div id="text-no-msg">
          <span id="fullname-no-msg">{props.selectedProfileData.firstname} {props.selectedProfileData.lastname}</span>
          <br/>
          <span style={{ fontSize: "12.5px" }}>Start the discussion by sending a message</span>
        </div>
      </div>
    )
  }
  return (
    <div id="list-msg-area" style={msgScroll}>
      {dataMessage}
    </div>
  )
}

const ListMatches = (props) => {
  var matches = [];
  var index = 0;
  props.listMatches.forEach((profile) => {
    matches.push(
      <MatchItem
        selectedProfile={props.selectedProfile}
        picture={profile.picture_url}
        username={profile.username}
        name={`${profile.firstname} ${profile.lastname}`}
        lastmsg={profile.last_message_content}
        time={profile.last_message_date}
        isOnline={profile.online}
        updateSelectedProfile={props.updateSelectedProfile}
        key={index}
      />
    );
    index += 1;
  });
  return (
    <div id="matches-list-main">
      {matches}
    </div>
  )
}

const GetListMatches = async (updateState) => {
  let res = await api.listMatches();
  if (res) {
    const response = await res.json();
    if (res.status >= 500) {
      throw new Error("Bad response from server - ListMatches has failed");
    } else if (res.status >= 400) {
      console.log(response.error);
    } else {
      if (response.data === "No matches") {
        updateState("listMatches", []);
      } else {
        updateState("listMatches", response);
      }
      return;
    }
  }
}

const GetListMessages = async (username, updateState) => {
  let res = await api.getMessages(username);
  if (res) {
    const response = await res.json();
    if (res.status >= 500) {
      throw new Error("Bad response from server - GetMe has failed");
    } else if (res.status >= 400) {
      console.log(response.error);
    } else {
      if (response.data === "No messages") {
        updateState("messages", []);
      } else {
        updateState("messages", response);
        console.log(response);
      }
      return;
    }
  }
}

class Chat extends Component {
  constructor (props) {
    super(props);
    this.state = {
      selectedProfile: '',
      messages: [],
      selectedProfileData: {},
      listMatches: [],
      message: '',
    }
    this.updateState = this.updateState.bind(this);
    this.handleChange = this.handleChange.bind(this);
    this.handleSubmit = this.handleSubmit.bind(this);
  }
  updateState(key, content) {
    this.setState({
      [key]: content,
    });
  }
  updateSelectedProfile = (username) => {
    this.state.listMatches.forEach((profile) => {
      if (profile.username === username) {
        this.updateState("selectedProfileData", profile);
        return;
      }
    })
    this.updateState("selectedProfile", username);
    GetListMessages(username, this.updateState);
  }
  componentDidMount() {
    GetListMatches(this.updateState);
  }
  handleSubmit = (e) => {
    console.log(this.state.message);
    this.updateState("message", "");
    e.preventDefault();
  }
  handleChange (e) {
    this.setState({
      [e.target.name]: e.target.value,
    });
  }
  render () {
    return (
      <div>
        <div id="matches-list" style={{ height: (window.innerHeight - 75) + "px" }}>
          <div id="matches-list-top">
            <span className="matches-list-title">Your matches</span>
          </div>
          <div className="limit" style={{ width: "50%", margin: "0px", marginBottom: "2.5px" }}></div>
          <ListMatches listMatches={this.state.listMatches} updateSelectedProfile={this.updateSelectedProfile} selectedProfile={this.state.selectedProfile}/>
        </div>
        {this.state.selectedProfile ? (
          <div id="chat-area">
            <div id="header-msg-area">
              <span id="main-title">{this.state.selectedProfileData.firstname} {this.state.selectedProfileData.lastname}</span>
              {/* is typing ... */}
            </div>
            <MsgArea listMsg={this.state.messages} myusername={"LawrenceHall2Hx24"} selectedProfileData={this.state.selectedProfileData}/>
            <div id="new-msg-area">
              {this.state.selectedProfileData.username !== undefined ? (
                <form id="form" onSubmit={this.handleSubmit}>
                    <input id="new-msg-input" type="text" placeholder="Type something to send ..." name="message" value={this.state.message} onChange={this.handleChange} size="64"/>
                    <button id="new-msg-submit" title="Send message"><img alt="Submit message" src={SendButton} style={{position: "absolute", top: "25%", left: "0px", width: "100%"}}/></button>
                </form>
              ) : null}
            </div>
          </div>
        ) : null}
      </div>
    )
  }
}

export default Chat;
