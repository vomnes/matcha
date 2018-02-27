import React, { Component } from 'react';
import './SeeProfile.css';
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
        <PictureArea
          picture={profilePictures[this.state.indexProfilePictures]}
          changePicture={this.changePicture}
          pictureArrayLength={profilePictures.length}
          index={this.state.indexProfilePictures}
          liked={this.state.liked}
          likeUser={this.likeUser}
          updateState={this.updateState}
          reportedAsFakeAccount={this.state.reportedAsFakeAccount}
        />
        <DataArea
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
