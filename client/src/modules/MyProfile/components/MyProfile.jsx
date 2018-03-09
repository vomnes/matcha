import React, { Component } from 'react';
import './MyProfile.css';
import Pictures  from './Pictures.jsx'
import Tags  from './Tags.jsx'
import Location  from './Location.jsx'
import ConfirmModal from '../../../components/ConfirmModal'
import utils from '../../../library/utils/array.js'

class MyProfile extends Component {
  constructor(props) {
    super(props);
    this.state = {
      firstname: '',
      lastname: '',
      email: '',
      biography: '',
      genre: 'male',
      preference: 'female',
      password: '',
      rePassword: '',
      birthday: '2018-03-09',
      tags: [
        "hello",
        "bonjour",
        "play",
        "tennis",
        "hello",
        "bonjour",
        "play",
        "tennis",
        "hello",
        "bonjour",
        "play",
        "tennis",
        "yes"
      ],
      newTag: '',
      variableModal: '',
    }
    this.handleUserInput = this.handleUserInput.bind(this);
    this.clickDeletePicture = this.clickDeletePicture.bind(this);
    this.cancelAction = this.cancelAction.bind(this);
    this.confirmDeletePicture = this.confirmDeletePicture.bind(this);
    this.deleteTag = this.deleteTag.bind(this);
  }
  handleUserInput = (e) => {
    const fieldName = e.target.name;
    var data = e.target.value;
    this.setState({
      [fieldName]: data,
    });
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
  render() {
    var profilePictures = [
      require('../../../design/pictures/Register-robson-hatsukami-morgan-250757-unsplash.jpg'),
      require('../../../design/pictures/Profile-molly-belle-73279-unsplash.jpg'),
      require('../../../design/pictures/Login-sorasak-217807-unsplash.jpg'),
    ];
    return (
      <div>
        <h1 className="top-title">Your profile</h1>
        <Pictures
          profilePictures={profilePictures}
          clickDeletePicture={this.clickDeletePicture}
        />
        <div className="myprofile-data">
          <div className="myprofile-id">
            <div className="fields-area">
              <div className="limit" style={{ width: "10%" }}></div>
              <div className="field-title">
                Edit your personal settings<br />
                <span className="profile-username">@vomnes</span><br />
              </div>
              <form className="profile-personal-data">
                <span className="field-name">Firstname</span><br />
                <input className="field-input" placeholder="Valentin" type="text" name="firstname"
                  value={this.state.firstname} onChange={this.handleUserInput}/><br />
                <span className="field-name">Lastname</span><br />
                <input className="field-input" placeholder="Omnes" type="text" name="lastname"
                  value={this.state.lastname} onChange={this.handleUserInput}/><br />
                <span className="field-name">Email address</span><br />
                <input className="field-input" placeholder="valentin.omnes@gmail.com" type="text" name="email"
                  value={this.state.email} onChange={this.handleUserInput}/><br />
                <span className="field-name">Biography</span><br />
                <input className="field-input" placeholder="Greatly hearted has who believe..." type="text" name="biography"
                  value={this.state.biography} onChange={this.handleUserInput}/><br />
                <span className="field-name">Birthday</span><br />
                <input className="field-input" type="date" name="birthday"
                  value={this.state.birthday} onChange={this.handleUserInput}/><br />
                <input className="submit-profile" type="submit" value="Save"/>
              </form>
              <div className="limit" style={{ width: "10%" }}></div>
              <form className="profile-personal-data">
                <span className="field-name">Genre</span><br />
                <select name="genre" value={this.state.genre} onChange={this.handleUserInput}>
                  <option value="female">Female</option>
                  <option value="male">Male</option>
                </select><br />
                <span className="field-name">Interesting in</span><br />
                <select name="preference" value={this.state.preference} onChange={this.handleUserInput}>
                  <option value="female">Female</option>
                  <option value="male">Male</option>
                  <option value="bisexual">Bisexual</option>
                </select><br />
                <input className="submit-profile" type="submit" value="Save"/>
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
                <input className="submit-profile" type="submit" value="Update password"/>
              </form>
              <div className="limit" style={{ width: "10%" }}></div>
              <div className="field-title">Set your location</div>
              <Location />
              <div className="limit" style={{ width: "10%" }}></div>
              <div className="field-title">Update your tags<br />
                <Tags tags={this.state.tags} deleteTag={this.deleteTag} newTag={this.state.newTag} appendTag={this.appendTag} handleUserInput={this.handleUserInput}/>
              </div>
            </div>
          </div>
        </div>
        <ConfirmModal type="deletePicture" number={this.state.variableModal}
          cancelAction={this.cancelAction}
          confirmAction={this.confirmDeletePicture}
        />
      </div>
    )
  }
}

export default MyProfile
