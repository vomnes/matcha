import React, { Component } from 'react';
import './SeeProfile.css';
import './SideProfiles.css';
import Modal from '../../../components/Modal';
import PictureArea from './PictureArea.jsx'
import DataArea from './DataArea.jsx'

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
    var matchedTags = [
      "hello",
      "bonjour",
      "play",
      "tennis",
      "yes"
    ]
    var otherTags = [
      "bye",
      "tag",
      "myTag",
      "awesome",
      "bye",
      "tag",
      "myTag",
      "awesome"
    ]
    var leftPicture = require('../../../design/pictures/Profile-molly-belle-73279-unsplash.jpg');
    var rightPicture = require('../../../design/pictures/Side-Picture-aziz-acharki-253909-unsplash.jpg');
    return (
      <div>
        <div className="left-side-profile">
          <div style={{ backgroundImage: "url(" + leftPicture + ")" }}></div>
          <span className="side-username">Valentin Omnès</span>
          <span className="side-see-profile" style={{ left: "40px" }}>&larr;</span>
        </div>
        <PictureArea
          picture={profilePictures[this.state.indexProfilePictures]}
          changePicture={this.changePicture}
          pictureArrayLength={profilePictures.length}
          index={this.state.indexProfilePictures}
          liked={this.state.liked}
          likeUser={this.likeUser}
          updateState={this.updateState}
          reportedAsFakeAccount={this.state.reportedAsFakeAccount}
          usersAreConnected={this.state.usersAreConnected}
          firstname="Valentin"
        />
        <div className="right-side-profile">
          <div style={{ backgroundImage: "url(" + rightPicture + ")" }}></div>
          <span className="side-username">Emma Thaero</span>
          <span className="side-see-profile" style={{ right: "40px" }}>&rarr;</span>
        </div>
        <DataArea
          firstname="Valentin"
          lastname="Omnes"
          username="vomnes"
          online={this.state.online}
          matchedTags={matchedTags}
          otherTags={otherTags}
        />
        <Modal type="success" online="true" content={this.state.newSuccess} onClose={this.closeModal}/>
      </div>
    )
  }
}

export default SeeProfile
