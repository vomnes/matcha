import React from 'react';
import './Matches.css';
import utils from '../../../library/utils/pictures.js'

var moment = require('moment');

const MatchItem = (props) => {
  if (props && props.time) {
    var date = moment(props.time)
    var formatedDate;
    if (moment().diff(date, 'years') > 100) {
      formatedDate = '';
    } else {
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
  }
  var matchStyle = {
    top: "-5px",
    left: "2.5px",
    border: "none",
    boxShadow: "0 3px 6px rgba(0,0,0,0.16), 0 3px 6px rgba(0,0,0,0.23)",
    zIndex: "2",
  };
  if (props.selectedProfile !== props.username) {
    matchStyle = {}
  }
  if (props.total_unread_messages) {
    matchStyle["background"] = "#eaeaea";
  }
  return (
    <div>
      <div className="match-element" id={props.username} style={matchStyle} onClick={() => props.updateSelectedProfile(props.username)}>
        {props.isOnline ? (<span className="online-dot" title={`${props.name} is online`}>&bull;</span>) : null}
        <div className="picture-list center">
          <a href={`/profile/${props.username}` + (props.optionsBase64 ? '/' + props.optionsBase64 : '')} title={`Click to see ${props.name}'s profile`}>
            <div className="picture-list-background" style={{ backgroundImage: "url(" + utils.pictureURLFormated(props.picture) + ")" }}></div>
          </a>
        </div>
        <span className="match-element-list">{props.name}</span>
        {props.selectedProfile === props.username ? null : (
          <div>
            {<span className="match-element-time" title={`Last message sent - ${formatedDate}`}>{formatedDate}</span>}
            <span className="match-element-lastmsg">{props.lastmsg}</span>
            {props.total_unread_messages ? (<span className="match-element-new-message-count red-cercle-notif">{props.total_unread_messages}</span>) : null}
          </div>
        )}
      </div>
      <div className="limit" style={{ width: "94.25%", margin: "2.5px 0px 2.5px" }}></div>
    </div>
  )
}

const ListMatches = (props) => {
  var matches = [];
  var index = 0;
  props.listMatches
  .sort((a, b) => {
    return moment(a.last_message_date).isBefore(moment(b.last_message_date));
  })
  .sort((a, b) => {
    return a.total_unread_messages < b.total_unread_messages;
  })
  .forEach((profile) => {
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
        total_unread_messages={profile.total_unread_messages}
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

export default ListMatches;
