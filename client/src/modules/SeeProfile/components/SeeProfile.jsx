import React, { Component } from 'react';
import './SeeProfile.css';
import './SideProfiles.css';
import Modal from '../../../components/Modal';
import PictureArea from './PictureArea.jsx'
import DataArea from './DataArea.jsx'
import api from '../../../library/api'
import { Redirect } from 'react-router-dom';

const getUserData = async (username, updateState) => {
  let res = await api.getuser(username);
  if (res) {
    const response = await res.json();
    if (res.status >= 500) {
      updateState("userExist", false);
      throw new Error("Bad response from server - GetUserData has failed");
    } else if (res.status >= 400) {
      updateState("userExist", false);
      console.log(response.error);
    } else {
      console.log(response);
      updateState("data", response);
      return;
    }
  }
}
class SeeProfile extends Component {
  constructor(props) {
    super(props);
    this.state = {
      userExist: true,
      data: {},
      indexProfilePictures: 0,
      lengthProfilePictures: 0,
      newSuccess: '',
    }
    this.changePicture = this.changePicture.bind(this);
    this.updateState = this.updateState.bind(this);
    this.closeModal = this.closeModal.bind(this);
    this.updateStateData = this.updateStateData.bind(this);
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
  updateState(key, successContent) {
    this.setState({
      [key]: !this.state[key],
      newSuccess: successContent
    });
  }
  updateStateData(key, content) {
    this.setState({
      [key]: content,
    });
  }
  closeModal(event) {
    this.setState({
      newError: '',
      newSuccess: ''
    });
    event.preventDefault();
  }

  componentDidMount() {
    getUserData(this.props.match.params.userId, this.updateStateData);
  }

  render() {
    let userData = Object.assign({}, this.state.data);
    var leftPicture = require('../../../design/pictures/Profile-molly-belle-73279-unsplash.jpg');
    var rightPicture = require('../../../design/pictures/Side-Picture-aziz-acharki-253909-unsplash.jpg');
    if (this.state.userExist) {
      return (
        <div>
          <div className="left-side-profile">
            <div style={{ backgroundImage: "url(" + leftPicture + ")" }}></div>
            <span className="side-username">Valentin Omnès</span>
            <span className="side-see-profile" style={{ left: "25px" }}>&larr;</span>
          </div>
          <PictureArea
            picture={(userData.pictures && userData.pictures[this.state.indexProfilePictures]) || null}
            changePicture={this.changePicture}
            pictureArrayLength={(userData.pictures && userData.pictures.length) || 0}
            index={this.state.indexProfilePictures}
            liked={userData.liked}
            updateState={this.updateStateData}
            reportedAsFakeAccount={userData.reportedAsFakeAccount}
            usersAreConnected={userData.usersAreConnected}
            firstname={userData.firstname}
            username={userData.username}
            isMe={userData.isMe}
          />
          <div className="right-side-profile">
            <div style={{ backgroundImage: "url(" + rightPicture + ")" }}></div>
            <span className="side-username">Emma Thaero</span>
            <span className="side-see-profile" style={{ right: "25px" }}>&rarr;</span>
          </div>
          <DataArea username={this.props.match.params.userId} data={this.state.data}/>
          <Modal type="success" online="true" content={this.state.newSuccess} onClose={this.closeModal}/>
        </div>
      )
    } else {
      return <Redirect to='/home'/>;
    }
  }
}

export default SeeProfile
