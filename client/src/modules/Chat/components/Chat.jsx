import React, { Component } from 'react';
import './Chat.css';
import encode from '../../../library/utils/encode.js'

class Chat extends Component {
  constructor (props) {
    super(props);
    console.log("create Chat");
  }
  componentDidMount() {
    var conn;
    var msg = document.getElementById("msg");
    document.getElementById("form").onsubmit = function () {
        if (!conn) {
            return false;
        }
        if (!msg.value) {
            return false;
        }
        conn.send(msg.value);
        msg.value = "";
        return false;
    };
    if (window["WebSocket"]) {
        var room = encode.objectToBase64({username1: "voluptas_atque_etpawSY", username2: "Elsa"});
        conn = new WebSocket(`ws://localhost:8081/ws/chat/${room}`, localStorage.getItem('matcha_token'));
        conn.onclose = function (evt) {
            var item = document.createElement("div");
            item.innerHTML = "<b>Connection closed.</b>";
            document.getElementById("log").appendChild(item);
        };
        conn.onmessage = function (evt) {
            var messages = evt.data.split('\n');
            for (var i = 0; i < messages.length; i++) {
                var item = document.createElement("div");
                item.innerText = messages[i];
                document.getElementById("log").appendChild(item);
            }
        };
    }
  }
  render () {
    return (
      <div>
        <span>Nice chat isn't it ?</span>
        <div id="log"></div>
        <form id="form">
            <input type="submit" value="Send" />
            <input type="text" id="msg" size="64"/>
        </form>
      </div>
    )
  }
}

export default Chat;
