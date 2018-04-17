import React, { Component } from 'react';
import './PictureArea.css';
import api from '../../../library/api'

const like = async (isLiked, method, username, updateStateHere, updateState, handleWebsocketSend) => {
  if (method === `POST` && isLiked) {
    return;
  }
  let res = await api.like(method, username);
  const response = await res.json();
  if (res.status >= 500) {
    throw new Error("Bad response from server - Like has failed - ", response.error);
  } else if (res.status >= 400) {
    updateState('newError', response.error);
    return;
  } else {
    var type = "liked"
    if (method === `DELETE`) {
      type = "unliked"
      updateStateHere("liked", false);
      updateStateHere("usersAreConnected", false);
      if (response.users_were_linked) {
        handleWebsocketSend("unmatch", username);
      }
    } else {
      updateStateHere("liked", true);
      updateStateHere("usersAreConnected", response.users_linked);
      if (response.users_linked) {
        handleWebsocketSend("match", username);
      } else {
        handleWebsocketSend("like", username);
      }
    }
    updateState('newSuccess', `You have just ${type} ${username}'s profile`);
  }
}

const fake = async (method, username, updateStateHere, updateState) => {
  let res = await api.fake(method, username);
  if (res.status >= 400) {
    const response = await res.json();
    if (res.status >= 500) {
      throw new Error("Bad response from server - Fake has failed - ", response.error);
    } else if (res.status >= 400) {
      updateState('newError', response.error);
      return;
    }
  }
  if (method === `DELETE`) {
    updateStateHere("reportedAsFakeAccount", false);
    updateState('newError', `You have just invalide you fake account report.`);
  } else {
    updateStateHere("reportedAsFakeAccount", true);
    updateState('newError', `You have just declared this profile as fake account.`);
  }
}

const IndexPictures = (props) => {
  var elements = [];
  for (var i = 0; i < props.pictureArrayLength; i++) {
    let color = props.index === i ? "black" : "grey";
    elements.push(<span key={i} style={{ color: color }}>&#8226;</span>);
  }
  return (
    <div id="picture-index">
      {elements}
    </div>
  )
}

class PictureArea extends Component {
  constructor(props) {
    super(props);
    this.state = {
      moreInformationOpen: false,
      index: 0,
      pictureArrayLength: 0,
      picture: '',
      reportedAsFakeAccount: null,
      liked: null,
    }
    this.openInformation = this.openInformation.bind(this);
    this.updateState = this.updateState.bind(this);
  }
  openInformation() {
    this.setState({
      moreInformationOpen: !this.state.moreInformationOpen,
    });
  }
  updateState(key, content) {
    this.setState({
      [key]: content,
    });
  }
  componentWillReceiveProps(nextProps) {
    if (this.state.username !== nextProps.username) {
      this.updateState('pictureArrayLength', nextProps.pictureArrayLength);
      this.updateState('index', nextProps.index);
      this.updateState('picture', nextProps.picture);
      this.updateState('liked', nextProps.liked);
      this.updateState('reportedAsFakeAccount', nextProps.reportedAsFakeAccount);
      this.updateState('usersAreConnected', nextProps.usersAreConnected);
      this.updateState('firstname', nextProps.firstname);
      this.updateState('username', nextProps.username);
      this.updateState('isMe', nextProps.isMe);
    }
  }
  render() {
    const index =  this.state.index;
    let pictureURL = this.state.picture ? `http://localhost:8080${this.state.picture}` : null
    if (this.state.picture && this.state.picture.includes('images.unsplash.com/photo-')) {
      pictureURL = this.state.picture;
    }
    var length = this.state.pictureArrayLength;
    return (
      <div className="picture-area">
        <div className="picture-user-background" style={{ backgroundImage: "url(" + pictureURL + ")" }}></div>
        {!this.state.isMe ? (
          <div>
            <div className="more-information" title={ 'See ' + (this.state.moreInformationOpen ? "less" :  "more") + ' options' }>
              <span className="more" onClick={this.openInformation}>{this.state.moreInformationOpen ? '-' : '+'}</span>
            </div>
            <div className="information-area" style={{ visibility: this.state.moreInformationOpen ? "visible" :  "hidden" } }>
                {this.state.reportedAsFakeAccount ?
                  <span onClick={() => fake(`DELETE`, this.state.username, this.updateState, this.props.updateState)}>Invalidate fake account report</span> :
                  <span onClick={() => fake(`POST`, this.state.username, this.updateState, this.props.updateState)}>Report as a fake account</span>
                }<br />
                {this.state.liked ? <span onClick={() => like(null, `DELETE`, this.state.username, this.updateState, this.props.updateState, this.props.handleWebsocketSend)}>Unlike profile</span> : null}
            </div>
          </div>
        ) : null }
        <IndexPictures pictureArrayLength={length} index={index}/>
        <div id="picture-previous" style={{ visibility: (index === 0) ? "hidden" :  "visible" }}>
          <span className="arrow" onClick={() => this.props.changePicture(0, length)}>&#x21A9;</span>
        </div>
        <div id="picture-next" style={{ visibility: (!length || index === (length - 1)) ? "hidden" :  "visible" } }>
          <span className="arrow" onClick={() => this.props.changePicture(1, length)}>&#x21AA;</span>
        </div>
        {!this.state.isMe ? (
          !this.state.usersAreConnected ? (
          <div title={(this.state.liked ? "You like this profile" :  "Like profile")} className="btn-like"
            onClick={() => like(this.state.liked, `POST`, this.state.username, this.updateState, this.props.updateState, this.props.handleWebsocketSend)}
            style={{
              background: (this.state.liked ? "white" :  "#F80759"),
              color: (this.state.liked ? "#F80759" :  "white"),
              cursor: (this.state.liked ? "default" :  "pointer") }}
            >
            <span>&#9829;</span>
          </div>
          ) : (
            <a href={`/matches`}>
              <div className="profiles-linked-picture">
                <span
                  role="img" aria-labelledby="Connected with"
                  title={'You are connected with ' + this.state.firstname + ' - Click here to take contact ;)' }
                >&#x1f525;</span>
              </div>
            </a>
          )
        ) : null}
      </div>
    )
  }
}

export default PictureArea;
