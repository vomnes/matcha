import React, { Component } from 'react';
import './MyProfile.css';

class EditPicture extends Component {
  constructor(props) {
    super(props);
    this.state = {
      pictureVisible: false,
    }
    this.showEditPicture = this.showEditPicture.bind(this);
    this.uploadPicture = this.uploadPicture.bind(this);
  }
  showEditPicture() {
    this.setState({
      pictureVisible: true,
    })
  }
  uploadPicture() {
    this.inputPicture.click();
  }
  render() {
    var url = this.props.urlPicture
    return (
      <div className={this.props.className} onMouseOver={this.showEditPicture} onMouseOut={this.showEditPicture}>
        {this.state.pictureVisible ? <span id="over" onClick={this.uploadPicture}>Change picture</span> : <span id="over">Change picture</span>}
        <input type="hidden" style={{display:"none"}} name="MAX_FILE_SIZE" value="30000" />
        <input ref={(input) => { this.inputPicture = input; }} id="upload-profile-picture-one" name="userfile" style={{display:"none"}} type="file" accept=".jpg, .jpeg, .png"/>
        <div className="picture-background-profile" style={{ backgroundImage: "url(" + url + ")" }}></div>
      </div>
    )
  }
}

class MyProfile extends Component {
  constructor(props) {
    super(props);
    this.state = {
    }
  }
  render() {
    var profilePictures = [
      require('../../../design/pictures/Register-robson-hatsukami-morgan-250757-unsplash.jpg'),
      require('../../../design/pictures/Profile-molly-belle-73279-unsplash.jpg'),
      require('../../../design/pictures/Login-sorasak-217807-unsplash.jpg'),
    ];
    return (
      <div>
        <h1>My profile</h1>
        <div className="pictures">
          <EditPicture className="profile-picture-one" urlPicture={profilePictures[2]} />
          <div className="picture-sub-area">
            <EditPicture className="profile-picture-two" urlPicture={profilePictures[2]} />
            <EditPicture className="profile-picture-three" urlPicture={profilePictures[1]} />
            <EditPicture className="profile-picture-four" urlPicture={profilePictures[2]} />
            <EditPicture className="profile-picture-five" urlPicture={profilePictures[2]} />
          </div>
        </div>
      </div>
    )
  }
}

export default MyProfile
