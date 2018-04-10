import React, { Component } from 'react';
import './Chat.css';
import SendButton from '../../../design/icons/send-button.svg';
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
        <div className="picture-list">
          <a href={`/profile/${props.username}` + (props.optionsBase64 ? '/' + props.optionsBase64 : '')} title={`Click to see ${props.name}'s profile`}>
            <div className="picture-list-background" style={{ backgroundImage: "url(" + props.picture + ")" }}></div>
          </a>
        </div>
        <span className="match-element-list">{props.name}</span>
        {props.selectedProfile === props.username ? null : (
          <div>
            <span className="match-element-time">{props.time}</span>
            <span className="match-element-lastmsg">{props.lastmsg}</span>
          </div>
        )}
      </div>
      <div className="limit" style={{ width: "94.25%", margin: "2.5px 0px 2.5px" }}></div>
    </div>
  )
}

class Chat extends Component {
  constructor (props) {
    super(props);
    this.state = {
      selectedProfile: '',
    }
    this.updateState = this.updateState.bind(this);
  }
  updateState(key, content) {
    this.setState({
      [key]: content,
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
          <div id="matches-list-main">
            <MatchItem selectedProfile={this.state.selectedProfile} picture="https://images.unsplash.com/photo-1522234811749-abc512463137?h=1000&q=10" username="a" name="Carolyn Wells" lastmsg="Good newsGood newsGood newsGood newsGood newsGood newsGood news" time={"12:30"} updateState={this.updateState}/>
            <MatchItem selectedProfile={this.state.selectedProfile} picture="https://images.unsplash.com/photo-1520512202623-51c5c53957df?h=1000&q=10" username="v" name="Pamela Ross" lastmsg="Good news" time={"12:30"} updateState={this.updateState}/>
            <MatchItem selectedProfile={this.state.selectedProfile} picture="https://images.unsplash.com/photo-1425009294879-3f15dd0b4ed5?h=1000&q=10" username="c" name="Pamela Ross" lastmsg="Good news" time={"12:30"} updateState={this.updateState}/>
            <MatchItem selectedProfile={this.state.selectedProfile} picture="https://images.unsplash.com/photo-1522234811749-abc512463137?h=1000&q=10" username="d" name="Louise Walker" lastmsg="Good news" time={"12:30"} updateState={this.updateState}/>
            <MatchItem selectedProfile={this.state.selectedProfile} picture="https://images.unsplash.com/photo-1522234811749-abc512463137?h=1000&q=10" username="1" name="Louise Walker" lastmsg="Good news" time={"12:30"} updateState={this.updateState}/>
            <MatchItem selectedProfile={this.state.selectedProfile} picture="https://images.unsplash.com/photo-1522234811749-abc512463137?h=1000&q=10" username="e" name="Louise Walker" lastmsg="Good news" time={"12:30"} updateState={this.updateState}/>
            <MatchItem selectedProfile={this.state.selectedProfile} picture="https://images.unsplash.com/photo-1522234811749-abc512463137?h=1000&q=10" username="f" name="Louise Walker" lastmsg="Good news" time={"12:30"} updateState={this.updateState}/>
            <MatchItem selectedProfile={this.state.selectedProfile} picture="https://images.unsplash.com/photo-1522234811749-abc512463137?h=1000&q=10" username="g" name="Louise Walker" lastmsg="Good news" time={"12:30"} updateState={this.updateState}/>
            <MatchItem selectedProfile={this.state.selectedProfile} picture="https://images.unsplash.com/photo-1522234811749-abc512463137?h=1000&q=10" username="h" name="Louise Walker" lastmsg="Good news" time={"12:30"} updateState={this.updateState}/>
            <MatchItem selectedProfile={this.state.selectedProfile} picture="https://images.unsplash.com/photo-1522234811749-abc512463137?h=1000&q=10" username="8" name="Louise Walker" lastmsg="Good news" time={"12:30"} updateState={this.updateState}/>
            <MatchItem selectedProfile={this.state.selectedProfile} picture="https://images.unsplash.com/photo-1522234811749-abc512463137?h=1000&q=10" username="j" name="Louise Walker" lastmsg="Good news" time={"12:30"} updateState={this.updateState}/>
            <MatchItem selectedProfile={this.state.selectedProfile} picture="https://images.unsplash.com/photo-1522234811749-abc512463137?h=1000&q=10" username="k" name="Louise Walker" lastmsg="Good news" time={"12:30"} updateState={this.updateState}/>
            <MatchItem selectedProfile={this.state.selectedProfile} picture="https://images.unsplash.com/photo-1522234811749-abc512463137?h=1000&q=10" username="l" name="Louise Walker" lastmsg="Good news" time={"12:30"} updateState={this.updateState}/>
            <MatchItem selectedProfile={this.state.selectedProfile} picture="https://images.unsplash.com/photo-1522234811749-abc512463137?h=1000&q=10" username="m" name="Louise Walker" lastmsg="Good news" time={"12:30"} updateState={this.updateState}/>
            <MatchItem selectedProfile={this.state.selectedProfile} picture="https://images.unsplash.com/photo-1522234811749-abc512463137?h=1000&q=10" username="n" name="Louise Walker" lastmsg="Good news" time={"12:30"} updateState={this.updateState}/>
            <MatchItem selectedProfile={this.state.selectedProfile} picture="https://images.unsplash.com/photo-1522234811749-abc512463137?h=1000&q=10" username="o" name="Louise Walker" lastmsg="Good news" time={"12:30"} updateState={this.updateState}/>
            <MatchItem selectedProfile={this.state.selectedProfile} picture="https://images.unsplash.com/photo-1522234811749-abc512463137?h=1000&q=10" username="p" name="Louise Walker" lastmsg="Good news" time={"12:30"} updateState={this.updateState}/>
            <MatchItem selectedProfile={this.state.selectedProfile} picture="https://images.unsplash.com/photo-1522234811749-abc512463137?h=1000&q=10" username="q" name="Louise Walker" lastmsg="Good news" time={"12:30"} updateState={this.updateState}/>
            <MatchItem selectedProfile={this.state.selectedProfile} picture="https://images.unsplash.com/photo-1522234811749-abc512463137?h=1000&q=10" username="r" name="Louise Walker" lastmsg="Good news" time={"12:30"} updateState={this.updateState}/>
            <MatchItem selectedProfile={this.state.selectedProfile} picture="https://images.unsplash.com/photo-1522234811749-abc512463137?h=1000&q=10" username="s" name="Louise Walker" lastmsg="Good news" time={"12:30"} updateState={this.updateState}/>
            <MatchItem selectedProfile={this.state.selectedProfile} picture="https://images.unsplash.com/photo-1522234811749-abc512463137?h=1000&q=10" username="t" name="Louise Walker" lastmsg="Good news" time={"12:30"} updateState={this.updateState}/>
            <MatchItem selectedProfile={this.state.selectedProfile} picture="https://images.unsplash.com/photo-1522234811749-abc512463137?h=1000&q=10" username="u" name="Louise Walker" lastmsg="Good news" time={"12:30"} updateState={this.updateState}/>
          </div>
        </div>
        <div id="chat-area">
          <div id="header-msg-area">
            <span id="main-title">Pamela Ross</span>
          </div>
          {/* <div className="picture-msg-list" style={{ right: "-22px" }}>
            <div className="picture-list-background" style={{ backgroundImage: "url(" + "https://images.unsplash.com/photo-1471943068829-26c1a4ac5bfa?h=1000&q=10" + ")" }}></div>
          </div> */}
          <div id="list-msg-area">
            <div id="msg-array">
              <div className="msg-element">
                <div className="msg" style={{ marginLeft: "5%" }}>
                  <div className="picture-msg-list" style={{ left: "-22px" }}>
                    <div className="picture-list-background" style={{ backgroundImage: "url(" + "https://images.unsplash.com/photo-1471943068829-26c1a4ac5bfa?h=1000&q=10" + ")" }}></div>
                  </div>
                  <div className="msg-content" style={{ marginLeft: "3%" }}>
                    <span>
                      laborum explicabo est autem voluptatum esse. debitis quis natus sequi vero velit eos. sit ratione doloremque necessitatibus. et omnis et sed veritatis! laudantium sit quas enim explicabo.
                      omnis et sed veritatis! laudantium.
                      laborum explicabo est autem voluptatum esse. debitis quis natus sequi vero velit eos. sit ratione doloremque necessitatibus. et omnis et sed veritatis! laudantium sit quas enim explicabo.
                    </span>
                  </div>
                </div>
              </div>
              <div className="msg-element">
                <div className="msg" style={{ marginLeft: "5%" }}>
                  <div className="picture-msg-list" style={{ left: "-22px" }}>
                    <div className="picture-list-background" style={{ backgroundImage: "url(" + "https://images.unsplash.com/photo-1471943068829-26c1a4ac5bfa?h=1000&q=10" + ")" }}></div>
                  </div>
                  <div className="msg-content" style={{ marginLeft: "3%" }}>
                    <span>
                      laborum explicabo est autem voluptatum esse. debitis quis natus sequi vero velit eos. sit ratione doloremque necessitatibus. et omnis et sed veritatis! laudantium sit quas enim explicabo.
                      omnis et sed veritatis! laudantium.
                      laborum explicabo est autem voluptatum esse. debitis quis natus sequi vero velit eos. sit ratione doloremque necessitatibus. et omnis et sed veritatis! laudantium sit quas enim explicabo.
                    </span>
                  </div>
                </div>
              </div>
              <div className="msg-element">
                <div className="msg" style={{ marginLeft: "5%" }}>
                  <div className="picture-msg-list" style={{ left: "-22px" }}>
                    <div className="picture-list-background" style={{ backgroundImage: "url(" + "https://images.unsplash.com/photo-1471943068829-26c1a4ac5bfa?h=1000&q=10" + ")" }}></div>
                  </div>
                  <div className="msg-content" style={{ marginLeft: "3%" }}>
                    <span>
                      laborum explicabo est autem voluptatum esse. debitis quis natus sequi vero velit eos. sit ratione doloremque necessitatibus. et omnis et sed veritatis! laudantium sit quas enim explicabo.
                      omnis et sed veritatis! laudantium.
                      laborum explicabo est autem voluptatum esse. debitis quis natus sequi vero velit eos. sit ratione doloremque necessitatibus. et omnis et sed veritatis! laudantium sit quas enim explicabo.
                    </span>
                  </div>
                </div>
              </div>
              <div className="msg-element">
                <div className="msg" style={{ marginLeft: "5%" }}>
                  <div className="picture-msg-list" style={{ left: "-22px" }}>
                    <div className="picture-list-background" style={{ backgroundImage: "url(" + "https://images.unsplash.com/photo-1471943068829-26c1a4ac5bfa?h=1000&q=10" + ")" }}></div>
                  </div>
                  <div className="msg-content" style={{ marginLeft: "3%" }}>
                    <span>
                      laborum explicabo est autem voluptatum esse. debitis quis natus sequi vero velit eos. sit ratione doloremque necessitatibus. et omnis et sed veritatis! laudantium sit quas enim explicabo.
                      omnis et sed veritatis! laudantium.
                      laborum explicabo est autem voluptatum esse. debitis quis natus sequi vero velit eos. sit ratione doloremque necessitatibus. et omnis et sed veritatis! laudantium sit quas enim explicabo.
                    </span>
                  </div>
                </div>
              </div>
              <div className="msg-element">
                <div className="msg" style={{ marginLeft: "10%" }}>
                  <div className="picture-msg-list" style={{ right: "-22px" }}>
                    <div className="picture-list-background" style={{ backgroundImage: "url(" + "https://images.unsplash.com/photo-1522234811749-abc512463137?h=1000&q=10" + ")" }}></div>
                  </div>
                  <div className="msg-content">
                    <span>
                      laborum explicabo est autem voluptatum esse. debitis quis natus sequi vero velit eos. sit ratione doloremque necessitatibus. et omnis et sed veritatis! laudantium sit quas enim explicabo.
                      omnis et sed veritatis! laudantium.
                    </span>
                  </div>
                </div>
              </div>
              <div className="msg-element">
                <div className="msg" style={{ marginLeft: "10%" }}>
                  <div className="picture-msg-list" style={{ right: "-22px" }}>
                    <div className="picture-list-background" style={{ backgroundImage: "url(" + "https://images.unsplash.com/photo-1522234811749-abc512463137?h=1000&q=10" + ")" }}></div>
                  </div>
                  <div className="msg-content">
                    <span>
                      laborum explicabo est autem voluptatum esse. debitis quis natus sequi vero velit eos. sit ratione doloremque necessitatibus. et omnis et sed veritatis! laudantium sit quas enim explicabo.
                      omnis et sed veritatis! laudantium.
                    </span>
                  </div>
                </div>
              </div>
              <div className="msg-element">
                <div className="msg" style={{ marginLeft: "5%" }}>
                  <div className="picture-msg-list" style={{ left: "-22px" }}>
                    <div className="picture-list-background" style={{ backgroundImage: "url(" + "https://images.unsplash.com/photo-1471943068829-26c1a4ac5bfa?h=1000&q=10" + ")" }}></div>
                  </div>
                  <div className="msg-content" style={{ marginLeft: "3%" }}>
                    <span>
                      laborum explicabo est autem voluptatum esse. debitis quis natus sequi vero velit eos.
                    </span>
                  </div>
                </div>
              </div>
            </div>
          </div>
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
