import React, { Component } from 'react';
import './SeeProfile.css'
import Modal from '../../../components/Modal'
import ConnectedLogo from '../../../design/icons/connected-128.png'
import GenderLogo from '../../../design/icons/gender-128.png'
import HeartLogo from '../../../design/icons/heart-128.png'
import LocationLogo from '../../../design/icons/location-128.png'
import Star from '../../../design/icons/star-128.png'
import Linked from '../../../design/icons/link-128.png'

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

class ProfilePicture extends Component {
  constructor(props) {
    super(props);
    this.state = {
      moreInformationOpen: false,
    }
    this.openInformation = this.openInformation.bind(this);
  }
  openInformation() {
    this.setState({
      moreInformationOpen: !this.state.moreInformationOpen,
    });
  }
  render() {
    const index = this.props.index;
    return (
      <div className="picture-area">
        <div className="picture-user-background" style={{ backgroundImage: "url(" + this.props.picture + ")" }}></div>
        <div className="more-information">
          <span className="more" onClick={this.openInformation}>{this.state.moreInformationOpen ? '-' : '+'}</span>
        </div>
        <div className="information-area" style={{ visibility: this.state.moreInformationOpen ? "visible" :  "hidden" } }>
            {this.props.reportedAsFakeAccount ?
              <span onClick={() => this.props.updateState("reportedAsFakeAccount", "You have just invalide you fake account report.")}>Invalidate fake account report</span> :
              <span onClick={() => this.props.updateState("reportedAsFakeAccount", "This profile has been declared as fake account.")}>Report as a fake account</span>
            }<br />
            {this.props.liked ? <span onClick={() => this.props.updateState("liked", "You have just unliked this profile")}>Unlike profile</span> : null}
        </div>
        <IndexPictures pictureArrayLength={this.props.pictureArrayLength} index={index}/>
        <div id="picture-previous" style={{ visibility: (index === 0) ? "hidden" :  "visible" }}>
          <span className="arrow" onClick={() => this.props.changePicture(0, this.props.pictureArrayLength)}>&#x21A9;</span>
        </div>
        <div id="picture-next" style={{ visibility: (index === (this.props.pictureArrayLength - 1)) ? "hidden" :  "visible" } }>
          <span className="arrow" onClick={() => this.props.changePicture(1, this.props.pictureArrayLength)}>&#x21AA;</span>
        </div>
        <div title="Like profile" className="btn-like"
          onClick={() => this.props.likeUser()}
          style={{
            background: (this.props.liked ? "white" :  "#F80759"),
            color: (this.props.liked ? "#F80759" :  "white"),
            cursor: (this.props.liked ? "default" :  "pointer") }}
          >
          <span>&#9829;</span>
        </div>
      </div>
    )
  }
}

const ShowTags = (props) => {
  var index = 0;
  var elements = [];
  props.tags.forEach(function (tag) {
    elements.push(<span key={index} className="picture-tag">#{tag}</span>);
    index += 1;
  });
  return (
    <div id="picture-tags">
      {elements}
    </div>
  )
}

class DateArea extends Component {
  constructor(props) {
    super(props);
    this.state = {
    }
  }
  render() {
    return (
      <div className="data-area">
        <div className="profile-id">
          <span className="profile-title">{this.props.firstname + ' ' + this.props.lastname}</span><br />
          <span className="profile-username"> @{this.props.username}</span><br />
          <div className="profile-bio">
            <span>Greatly hearted has who believe...</span>
          </div>
          <div className="profile-data-list">
            <img alt="Connected status logo" title="Connected status" src={ConnectedLogo}/>
            <span>{this.props.online ? "Online" : "Offline - Last connection 60 minutes ago"}</span><br />
            <img alt="Gender logo" title="Gender" src={GenderLogo}/><span>Male</span><br />
            <img alt="Preferred gender logo" title="Preferred gender" src={HeartLogo}/><span>Female</span><br />
            <img alt="Location logo" title="Location" src={LocationLogo}/><span>Paris, France</span><br />
            <img alt="Rating logo" title="Rating" src={Star}/><span>95/100</span><br />
          </div>
          <ShowTags tags={this.props.tags}/>
        </div>
        {this.props.usersAreConnected ? (
          <div className="profiles-linked">
            <span title={'You are connected with ' + this.props.firstname + ' - Click here to take contact ;)' }>&#x1F517;</span>
          </div>
        ) : null}
      </div>
    )
  }
}

class SeeProfile extends Component {
  constructor(props) {
    super(props);
    this.state = {
      indexProfilePictures: 0,
      lengthProfilePictures: 0,
      liked: false,
      reportedAsFakeAccount: false,
      newSuccess: '',
      online: false,
      usersAreConnected: true
    }
    this.changePicture = this.changePicture.bind(this);
    this.updateState = this.updateState.bind(this);
    this.closeModal = this.closeModal.bind(this);
    this.likeUser = this.likeUser.bind(this);
  }
  changePicture(value, len) {
    var change;
    if (value === 1) {
      change = 1;
    } else {
      change = -1;
    }
    const update = this.state.indexProfilePictures + change;
    if (update === -1 || update === len) {
      return
    }
    this.setState({
      indexProfilePictures: update
    });
  }
  likeUser() {
    if (this.state.liked) {
      return
    }
    this.setState({
      liked: true,
      newSuccess: 'You have just liked xxx\'s profile'
    });
  }
  updateState(key, successContent) {
    this.setState({
      [key]: !this.state[key],
      newSuccess: successContent
    });
  }
  closeModal(event) {
    this.setState({
      newError: '',
      newSuccess: ''
    });
    event.preventDefault();
  }
  render() {
    var profilePictures = [
      require('../../../design/pictures/Register-robson-hatsukami-morgan-250757-unsplash.jpg'),
      require('../../../design/pictures/Profile-molly-belle-73279-unsplash.jpg'),
      require('../../../design/pictures/Login-sorasak-217807-unsplash.jpg'),
    ];
    var tags = [
      "hello",
      "bonjour",
      "play",
      "tennis",
      "bye",
      "tag",
      "myTag",
      "awesome",
      "bye",
      "tag",
      "myTag",
      "awesome",
      "yes"
    ]
    return (
      <div>
        <ProfilePicture
          picture={profilePictures[this.state.indexProfilePictures]}
          changePicture={this.changePicture}
          pictureArrayLength={profilePictures.length}
          index={this.state.indexProfilePictures}
          liked={this.state.liked}
          likeUser={this.likeUser}
          updateState={this.updateState}
          reportedAsFakeAccount={this.state.reportedAsFakeAccount}
        />
        <DateArea
          firstname="Valentin"
          lastname="Omnes"
          username="vomnes"
          tags={tags}
          online={this.state.online}
          usersAreConnected={this.state.usersAreConnected}
        />
        <Modal type="success" online="true" content={this.state.newSuccess} onClose={this.closeModal}/>
      </div>
    )
  }
}

export default SeeProfile
