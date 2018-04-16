import React, { Component } from 'react';
import './Chat.css';
import SendButton from '../../../design/icons/send-button.svg';
import api from '../../../library/api'
import ListMatches from './Matches.jsx'
import MsgArea from './Messages.jsx'

var moment = require('moment');

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
        console.log(response);
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
      }
      return;
    }
  }
}

const wsSend = (conn, object) => {
  conn.send(JSON.stringify(object));
}

class Chat extends Component {
  constructor (props) {
    super(props);
    this.state = {
      selectedProfile: '',
      messages: [],
      selectedProfileData: {},
      loggedProfileData: {},
      listMatches: [],
      message: '',
      isTyping: false,
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
    this.updateProfile({total_unread_messages: 0, username: username});
    GetListMessages(username, this.updateState);
  }
  componentDidMount() {
    GetListMatches(this.updateState);
    var wsConn = new WebSocket(`ws://localhost:8081/ws/${localStorage.getItem(`matcha_token`)}`);
    this.props.myProfileData.then((data) => {
      this.updateState("loggedProfileData", data);
    });
    wsConn.onmessage = (e) => {
      try {
        var msg = JSON.parse(e.data);
      } catch (e) {
        return;
      }
      this.handleWebsocket(msg);
    }
    this.updateState('wsConn', wsConn);
  }
  handleSubmit = (e) => {
    if (this.state.message !== "") {
      wsSend(this.state.wsConn, {
        "event": "message",
        "target": this.state.selectedProfile,
        "data": this.state.message,
      });
      const data = this.state.loggedProfileData;
      this.appendMessage(
        {
          content: this.state.message,
          from: data.username
        }, {
          firstname: data.firstname,
          lastname: data.lastname,
          picture_url: data.profile_picture
        },
      );
      this.updateState("message", "");
      this.updateState('isTyping', false);
    }
    e.preventDefault();
  }
  handleChange (e) {
    if (e.target.value.length > 255) {
      return;
    }
    wsSend(this.state.wsConn, {
      "event": "isTyping",
      "target": this.state.selectedProfile,
    });
    this.setState({
      [e.target.name]: e.target.value,
    });
  }
  appendMessage(message, {firstname, lastname, picture_url}) {
    this.state.messages.push({
      content: message.content,
      username: message.from,
      received_at: moment(),
      firstname,
      lastname,
      picture_url,
    })
    this.updateState('messages', this.state.messages);
  }
  updateProfile({total_unread_messages, online, username, content}) {
    var data = this.state.listMatches;
    var len = data.length;
    for (var i = 0; i < len; i++) {
      if (data[i].username === username) {
        if (online !== undefined) {
          data[i].online = online;
        }
        if (total_unread_messages !== undefined) {
          if (total_unread_messages === 0) {
            data[i].total_unread_messages = 0;
          } else {
            data[i].total_unread_messages += 1;
          }
        }
        if (content !== undefined) {
          data[i].last_message_content = content;
          data[i].last_message_date = moment();
        }
      }
    }
    this.updateState('listMatches', data);
  }
  handleWebsocket = (msg) => {
    if (msg.event === "message") {
      if (msg.data.from === this.state.selectedProfile) {
        const data = this.state.selectedProfileData;
        this.appendMessage(msg.data, {
          firstname: data.firstname,
          lastname: data.lastname,
          picture_url: data.picture_url,
        });
      } else {
        this.updateProfile({total_unread_messages: 1, username: msg.data.from, content: msg.data.content});
      }
      console.log(`New message from ${msg.data.from} - Content: ${msg.data.content}`);
    }
    if (msg.event === "isTyping") {
      if (msg.from === this.state.selectedProfile && !this.state.isTyping) {
        this.updateState('isTyping', true);
        setTimeout(async () => {
          this.updateState('isTyping', false);
        }, 2000);
      }
    }
    if (msg.event === "login" || msg.event === "logout") {
      this.updateProfile({online: msg.event === "login" ? true : false, username: msg.username});
    }
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
              <span id="main-title">{this.state.selectedProfileData.firstname} {this.state.selectedProfileData.lastname} {this.state.isTyping ? ("is typing ...") : null}</span>
            </div>
            <MsgArea listMsg={this.state.messages} myusername={this.state.loggedProfileData.username} selectedProfileData={this.state.selectedProfileData}/>
            <div id="new-msg-area">
              {this.state.selectedProfileData.username !== undefined ? (
                <form id="form" onSubmit={this.handleSubmit}>
                    <input id="new-msg-input" type="text" placeholder="Type something to send ..." maxLength="255" name="message" value={this.state.message} onChange={this.handleChange} size="64"/>
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
