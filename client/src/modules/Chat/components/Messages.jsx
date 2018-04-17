import React from 'react';
import './Messages.css';
import utils from '../../../library/utils/pictures.js'

var moment = require('moment');

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
          <div className="picture-list-background" style={{ backgroundImage: "url(" + utils.pictureURLFormated(props.picture) + ")" }}></div>
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
            <div className="picture-list-background" style={{ backgroundImage: `url(${utils.pictureURLFormated(props.selectedProfileData.picture_url)})` }}></div>
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

export default MsgArea;
