import React, { Component } from 'react';
import './Chat.css';
import SendButton from '../../../design/icons/send-button.svg';
import api from '../../../library/api'
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
  return (
    <div>
      <div className="match-element" id={props.username} style={props.selectedProfile === props.username ? selectedStyle : null } onClick={() => props.updateState("selectedProfile", props.username)}>
        {props.isOnline ? (<span className="online-dot" title={`${props.name} is online`}>&bull;</span>) : null}
        <div className="picture-list">
          <a href={`/profile/${props.username}` + (props.optionsBase64 ? '/' + props.optionsBase64 : '')} title={`Click to see ${props.name}'s profile`}>
            <div className="picture-list-background" style={{ backgroundImage: "url(" + props.picture + ")" }}></div>
          </a>
        </div>
        <span className="match-element-list">{props.name}</span>
        {props.selectedProfile === props.username ? null : (
          <div>
            <span className="match-element-time" title={`Last message sent at ${props.time}`}>{props.time}</span>
            <span className="match-element-lastmsg">{props.lastmsg}</span>
          </div>
        )}
      </div>
      <div className="limit" style={{ width: "94.25%", margin: "2.5px 0px 2.5px" }}></div>
    </div>
  )
}

const MsgItem = (props) => {
  var msgStyle = null;
  var pictureStyle = null;
  var msgContentStyle = null;
  var msgHeaderStyle = null;
  if (props.side === "left") {
    msgStyle = { marginLeft: "5%" };
    pictureStyle = { left: "-22px" };
    msgContentStyle = { marginLeft: "3%" };
    msgHeaderStyle = { left: "25px" };
  } else {
    msgStyle = { marginLeft: "10%" };
    pictureStyle = { right: "-22px" };
    msgHeaderStyle = { right: "30px" };
  }
  return (
    <div className="msg-element">
      <div className="msg" style={msgStyle}>
        <div className="picture-msg-list" style={pictureStyle}>
          <div className="picture-list-background" style={{ backgroundImage: "url(" + props.picture + ")" }}></div>
        </div>
        <span className="msg-header" style={msgHeaderStyle}>{props.firstname} {props.lastname} - {props.received_at}</span>
        <div className="msg-content" style={msgContentStyle}>
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
        picture={message.picture}
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
  } else {
    dataMessage = (
      <div id="no-msg">
        <a href={`/profile/${props.username}` + (props.optionsBase64 ? '/' + props.optionsBase64 : '')} title={`Click to see ${props.name}'s profile`}>
          <div id="picture-no-msg" style={null}>
            <div className="picture-list-background" style={{ backgroundImage: "url(https://images.unsplash.com/photo-1522234811749-abc512463137?h=1000&q=10)" }}></div>
          </div>
        </a>
        <div id="text-no-msg">
          <span id="fullname-no-msg">Pamela Ross</span>
          <br/>
          <span style={{ fontSize: "12.5px" }}>Start the discussion by typing a message</span>
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
        updateState={props.updateState}
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
      throw new Error("Bad response from server - GetMe has failed");
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

class Chat extends Component {
  constructor (props) {
    super(props);
    this.state = {
      selectedProfile: '',
      messages: [
        {
          username: "vomnes",
          firstname: "Valentin",
          lastname: "Omnes",
          picture: "https://images.unsplash.com/photo-1471943068829-26c1a4ac5bfa?h=1000&q=10",
          content: "laborum explicabo est autem voluptatum esse. debitis quis natus sequi vero velit eos. sit ratione doloremque necessitatibus. et omnis et sed veritatis! laudantium sit quas enim explicabo. omnis et sed veritatis! laudantium. laborum explicabo est autem voluptatum esse. debitis quis natus sequi vero velit eos. sit ratione doloremque necessitatibus. et omnis et sed veritatis! laudantium sit quas enim explicabo.",
          received_at: "10:10",
        }
      ],
      selectedProfileData: [],
      listMatches: [],
    }
    this.updateState = this.updateState.bind(this);
  }
  updateState(key, content) {
    this.setState({
      [key]: content,
    });
  }
  componentDidMount() {
    GetListMatches(this.updateState);
  }
  render () {
    return (
      <div>
        <div id="matches-list" style={{ height: (window.innerHeight - 75) + "px" }}>
          <div id="matches-list-top">
            <span className="matches-list-title">Your matches</span>
          </div>
          <div className="limit" style={{ width: "50%", margin: "0px", marginBottom: "2.5px" }}></div>
          <ListMatches listMatches={this.state.listMatches} updateState={this.updateState} selectedProfile={this.state.selectedProfile}/>
        </div>
        <div id="chat-area">
          <div id="header-msg-area">
            <span id="main-title">Pamela Ross is typing ...</span>
            {/*  */}
          </div>
          <MsgArea listMsg={this.state.messages} myusername={"vomnes"}/>
          <div id="new-msg-area">
            <form id="form">
                <input id="new-msg-input" type="text" placeholder="Type something to send ..." size="64"/>
                <button id="new-msg-submit" title="Send message"><img alt="Submit message" src={SendButton} style={{position: "absolute", top: "25%", left: "0px", width: "100%"}}/></button>
            </form>
          </div>
        </div>
      </div>
    )
  }
}

export default Chat;
