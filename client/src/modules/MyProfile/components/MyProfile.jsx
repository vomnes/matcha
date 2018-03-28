import React, { Component } from 'react';
import './MyProfile.css';
import Pictures  from './Pictures.jsx'
import Tags  from './Tags.jsx'
import Location  from './Location.jsx'
import ConfirmModal from '../../../components/ConfirmModal'
import Modal from '../../../components/Modal'
import utils from '../../../library/utils/array.js'
import api from '../../../library/api'

const getProfileData = async (ip, updateState, setStateGenre) => {
  let res = await api.getprofile(ip);
  if (res) {
    const response = await res.json();
    if (res.status >= 500) {
      throw new Error("Bad response from server - GetProfileData has failed");
    } else if (res.status >= 400) {
      console.log(response.error);
    } else {
      console.log(response);
      updateState("data", response);
      setStateGenre();
      return;
    }
  }
}

const getIP = async (getData, updateState, setStateGenre) => {
  let res;
  res = await fetch("https://freegeoip.net/json/", {
    method: 'GET',
    mod: 'no-cors',
  });
  if (res) {
    const json = await res.json();
    getData(json.ip, updateState, setStateGenre);
    return;
  }
}

const editProfile = async (args, originalData, updateState, updatePersonal, updateData) => {
  let res = await api.editprofile(args);
  if (res && res.status >= 400) {
    const response = await res.json();
    if (res.status >= 500) {
      throw new Error("Bad response from server - EditProfile has failed - ", response.error);
    } else if (res.status >= 400) {
      updateState('newError', response.error);
    }
  } else {
    if (args.firstname) {
      updateData('firstname', args.firstname);
      updatePersonal('firstname', '');
    }
    if (args.lastname) {
      updateData('lastname', args.lastname);
      updatePersonal('lastname', '');
    }
    if (args.email) {
      updateData('email', args.email);
      updatePersonal('email', '');
    }
    if (args.biography) {
      updateData('firstname', args.firstname);
      updatePersonal('firstname', '');
    }
    if (args.birthday) {
      updateData('birthday', args.birthday);
      updatePersonal('birthday', '');
    }
    if (args.genre !== originalData.genre) {
      updateData('genre', args.genre);
    }
    if (args.interesting_in !== originalData.interesting_in) {
      updateData('interesting_in', args.interesting_in);
    }
  }
}

const isEmpty = (value) => {
  return value !== undefined && value !== ''
}

class MyProfile extends Component {
  constructor(props) {
    super(props);
    this.state = {
      data: {},
      personal: {},
      password: '',
      rePassword: '',
      newTag: '',
      variableModal: '',
      newError: '',
      newSuccess: '',
    }
    this.handleUserInput = this.handleUserInput.bind(this);
    this.handleUserInputPersonal = this.handleUserInputPersonal.bind(this);
    this.handleUserInputData = this.handleUserInputData.bind(this);
    this.updateStatePersonal = this.updateStatePersonal.bind(this);
    this.clickDeletePicture = this.clickDeletePicture.bind(this);
    this.cancelAction = this.cancelAction.bind(this);
    this.confirmDeletePicture = this.confirmDeletePicture.bind(this);
    this.deleteTag = this.deleteTag.bind(this);
    this.updateState = this.updateState.bind(this);
    this.setStateGenre = this.setStateGenre.bind(this);
    this.closeModal = this.closeModal.bind(this);
    this.handleSubmitPersonal = this.handleSubmitPersonal.bind(this);
  }
  handleUserInput = (e) => {
    const fieldName = e.target.name;
    var data = e.target.value;
    this.setState({
      [fieldName]: data,
    });
  }
  handleUserInputPersonal = (e) => {
    var data = this.state.personal;
    data[e.target.name] = e.target.value;
    this.setState({
      personal: data,
    });
  }
  handleUserInputData = (field, value) => {
    var data = this.state.data;
    data[field] = value;
    console.log("data:", data);
    this.setState({
      data: data,
    });
  }
  updateStatePersonal = (field, value) => {
    var personal = this.state.personal;
    personal[field] = value;
    console.log("personal:", personal);
    this.setState({
      personal: personal,
    });
  }
  updateState(key, content) {
    this.setState({
      [key]: content,
    });
  }
  setStateGenre() {
    this.setState({
      personal: {
        genre: this.state.data.genre,
        interesting_in: this.state.data.interesting_in,
      }
    })
  }
  cancelAction = () => {
    this.setState({
      variableModal: '',
    })
  }
  clickDeletePicture = (id) => {
    this.setState({
      variableModal: id,
    })
  }
  confirmDeletePicture = () => {
    console.log("Picture " + this.state.variableModal + " deleted");
    this.setState({
      variableModal: '',
    })
  }
  deleteTag = (name) => {
    var newTags = utils.removeItemByValue(this.state.tags, name);
    this.setState({
      tags: newTags,
    })
  }
  appendTag = (name) => {
    var newTags = this.state.tags.concat(name);
    this.setState({
      tags: newTags,
      newTag: '',
    })
  }
  closeModal(event) {
    this.setState({
      newError: '',
      newSuccess: ''
    });
    event.preventDefault();
  }
  handleSubmitPersonal = (event) => {
    editProfile(this.state.personal, this.state.data, this.updateState, this.updateStatePersonal, this.handleUserInputData);
    event.preventDefault();
  }

  componentDidMount() {
    getIP(getProfileData, this.updateState, this.setStateGenre);
  }

  render() {
    let userData = Object.assign({}, this.state.data);
    var profilePictures = [
      userData.picture_url_1,
      userData.picture_url_2,
      userData.picture_url_3,
      userData.picture_url_4,
      userData.picture_url_5,
    ];
    var updatePersonalDataBtn;
    if (isEmpty(this.state.personal.firstname) || isEmpty(this.state.personal.lastname) || isEmpty(this.state.personal.email) ||
  isEmpty(this.state.personal.biography) || isEmpty(this.state.personal.birthday) ||
  this.state.personal.genre !== this.state.data.genre || this.state.personal.interesting_in !== this.state.data.interesting_in) {
        updatePersonalDataBtn = <input className="submit-profile" type="submit" value="Save" title="Save personal data"/>
    }
    var updatePasswordBtn;
    if (isEmpty(this.state.password) && isEmpty(this.state.rePassword)) {
        updatePasswordBtn = <input className="submit-profile" type="submit" value="Update" title="Update password"/>
    }
    return (
      <div>
        <h1 className="top-title">Your profile</h1>
        <Pictures
          profilePictures={profilePictures}
          clickDeletePicture={this.clickDeletePicture}
          updatePicture={this.handleUserInputData}
          updateState={this.updateState}
        />
        <div className="myprofile-data">
          <div className="myprofile-id">
            <div className="fields-area">
              <div className="limit" style={{ width: "10%" }}></div>
              <div className="field-title">
                Edit your personal settings<br />
                <span className="profile-username">@{userData.username || ''}</span><br />
              </div>
              <form className="profile-personal-data" onSubmit={this.handleSubmitPersonal}>
                <span className="field-name">Firstname</span><br />
                <input className="field-input" placeholder={this.state.data.firstname || ''} type="text" name="firstname"
                  value={this.state.personal.firstname || ''} onChange={this.handleUserInputPersonal}/><br />
                <span className="field-name">Lastname</span><br />
                <input className="field-input" placeholder={userData.lastname || ''} type="text" name="lastname"
                  value={this.state.personal.lastname || ''} onChange={this.handleUserInputPersonal}/><br />
                <span className="field-name">Email address</span><br />
                <input className="field-input" placeholder={userData.email || ''} type="text" name="email"
                  value={this.state.personal.email || ''} onChange={this.handleUserInputPersonal}/><br />
                <span className="field-name">Biography</span><br />
                <input className="field-input" placeholder={userData.biography || ''} type="text" name="biography"
                  value={this.state.personal.biography || ''} onChange={this.handleUserInputPersonal}/><br />
                <span className="field-name">Birthday</span><br />
                <input className="field-input" placeholder={userData.birthday || ''} type="text" name="birthday"
                  value={this.state.personal.birthday || ''} onChange={this.handleUserInputPersonal}/><br />
                <div className="limit" style={{ width: "10%" }}></div>
                <span className="field-name">Genre</span><br />
                <select className="field-input" name="genre" value={this.state.personal.genre} onChange={this.handleUserInputPersonal}>
                  <option value="female">Female</option>
                  <option value="male">Male</option>
                </select><br />
                <span className="field-name">Interesting in</span><br />
                <select className="field-input" name="interesting_in" value={this.state.personal.interesting_in} onChange={this.handleUserInputPersonal}>
                  <option value="female">Female</option>
                  <option value="male">Male</option>
                  <option value="bisexual">Bisexual</option>
                </select><br />
                {updatePersonalDataBtn}
              </form>
              <div className="limit" style={{ width: "10%" }}></div>
              <div className="field-title">Set your password</div>
              <form className="profile-personal-data">
                <span className="field-name">New password</span><br />
                <input className="field-input" type="password" name="password"
                  value={this.state.password} onChange={this.handleUserInput}/><br />
                <span className="field-name">Type it again</span><br />
                  <input className="field-input" type="password" name="rePassword"
                    value={this.state.rePassword} onChange={this.handleUserInput}/><br />
                {updatePasswordBtn}
              </form>
              <div className="limit" style={{ width: "10%" }}></div>
              <div className="field-title">Update your location</div>
              <Location
                lat={this.state.data.latitude}
                lng={this.state.data.longitude}
                geolocalisation_allowed={this.state.data.geolocalisation_allowed}
                updateState={this.updateState}
              />
              <div className="limit" style={{ width: "10%" }}></div>
              <div className="field-title">Update your tags<br />
                <Tags tags={userData.tags || []} deleteTag={this.deleteTag} newTag={this.state.newTag} appendTag={this.appendTag} handleUserInput={this.handleUserInput}/>
              </div>
            </div>
          </div>
        </div>
        <ConfirmModal type="deletePicture" number={this.state.variableModal}
          cancelAction={this.cancelAction}
          confirmAction={this.confirmDeletePicture}
        />
        <Modal online={true} type="error" content={this.state.newError} onClose={this.closeModal}/>
        <Modal online={true} type="success" content={this.state.newSuccess} onClose={this.closeModal}/>
      </div>
    )
  }
}

export default MyProfile
