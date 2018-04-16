import React, { Component } from 'react';
import './SeeProfile.css';
import './SideProfiles.css';
import Modal from '../../../components/Modal';
import PictureArea from './PictureArea.jsx'
import DataArea from './DataArea.jsx'
import api from '../../../library/api'
import { Redirect } from 'react-router-dom';

var moment = require('moment');

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
      updateState("data", response);
      return;
    }
  }
}

const targetedMatch = async (optionsBase64, username, updateState) => {
  let res = await api.targetedmatch(optionsBase64, username);
  if (res) {
    const response = await res.json();
    if (res.status >= 500) {
      throw new Error("Bad response from server - GetMatch has failed");
    } else if (res.status >= 400) {
      console.log(response.error);
    } else {
      if (response.previous) {
        updateState("previousProfile", response.previous);
      }
      if (response.next) {
        updateState("nextProfile", response.next);
      }
    }
  }
}

const SideProfile = (props) => {
  return (
    <div className={`${props.order}-profile`} title={`See ${props.firstname} ${props.lastname}'s profile`} onClick={() => props.getSideProfile(props.username)}>
      <div className="picture-user-background" style={{ backgroundImage: "url(" + props.picture_url + ")" }}>
        <span className="side-fullname">{`${props.firstname} ${props.lastname}`}</span>
      </div>
    </div>
  );
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
      previousProfile: {},
      nextProfile: {},
      searchparameters: this.props.match.params.searchparameters,
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
  getSideProfile = (username) => {
    this.updateStateData("previousProfile", {});
    this.updateStateData("nextProfile", {});
    getUserData(username, this.updateStateData);
    targetedMatch(this.state.searchparameters, username, this.updateStateData);
    this.props.match.params.username = username;
  }
  handleWebsocket = (msg) => {
    if ((msg.event === "login" || msg.event === "logout") && this.state.data.username === msg.username) {
      var profile = this.state.data
      profile['online'] = msg.event === "login" ? true : false;
      profile['logout_at'] = msg.event === "login" ? '' : moment();
      this.updateStateData('data', profile);
    }
  }

  componentDidMount() {
    getUserData(this.props.match.params.username, this.updateStateData);
    targetedMatch(this.props.match.params.searchparameters, this.props.match.params.username, this.updateStateData);
  }

  render() {
    let userData = Object.assign({}, this.state.data);
    this.props.wsConn.onmessage = (e) => {
      try {
        var msg = JSON.parse(e.data);
      } catch (e) {
        return;
      }
      this.handleWebsocket(msg);
    }
    console.log(userData);
    var left;
    var right;
    if (this.state.previousProfile.picture_url) {
      left = <SideProfile getSideProfile={this.getSideProfile} order="previous" username={this.state.previousProfile.username} picture_url={this.state.previousProfile.picture_url} firstname={this.state.previousProfile.firstname} lastname={this.state.previousProfile.lastname}/>;
    } else {
      left = null;
    }
    if (this.state.nextProfile.picture_url) {
      right = <SideProfile getSideProfile={this.getSideProfile} order="next" username={this.state.nextProfile.username} picture_url={this.state.nextProfile.picture_url} firstname={this.state.nextProfile.firstname} lastname={this.state.nextProfile.lastname}/>;
    } else {
      right = null;
    }

    if (this.state.userExist) {
      return (
        <div>
          {left}
          <PictureArea
            picture={(userData.pictures && userData.pictures[this.state.indexProfilePictures]) || null}
            changePicture={this.changePicture}
            pictureArrayLength={(userData.pictures && userData.pictures.length) || 0}
            index={this.state.indexProfilePictures}
            liked={userData.liked}
            updateStateData={this.updateStateData}
            updateState={this.updateState}
            reportedAsFakeAccount={userData.reported_as_fake}
            usersAreConnected={userData.users_linked}
            firstname={userData.firstname}
            username={userData.username}
            isMe={userData.isMe}
          />
          {right}
          <DataArea username={this.props.match.params.username} data={this.state.data}/>
          <Modal type="success" online="true" content={this.state.newSuccess} onClose={this.closeModal}/>
        </div>
      )
    } else {
      return <Redirect to='/home'/>;
    }
  }
}

export default SeeProfile
