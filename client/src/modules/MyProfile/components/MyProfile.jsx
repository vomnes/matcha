import React, { Component } from 'react';
import './MyProfile.css';
import InfoLogo from '../../../design/icons/information-128.png'

const ConfirmModal = (props) => {
  if (props.number) {
    var content;
    if (props.type === "deletePicture") {
        content = "Delete the picture number " + props.number;
    }
    return (
      <div className="modal-confirm-area">
        <div className="modal-confirm">
          <img alt="Information logo" className="modal-confirm-icon" width="32px" src={InfoLogo}/>
          <span className="modal-confirm-content">{content}</span>
          <span className="modal-confirm-cancel" onClick={props.cancelAction}>Cancel</span>
          <span className="modal-confirm-confirm" onClick={props.confirmAction}>Confirm</span>
        </div>
      </div>
    )
  }
  return null;
}

const ButtonEdit = (props) => {
  var deleteBtn = null;
  if (props.deleteAvailable) {
    deleteBtn = <span className="over-delete" title="Delete picture" onClick={() => props.deletePicture(props.id)}>x</span>;
  }
  if (props.pictureVisible) {
    return (
      <div>
        <span className="over-edit" onClick={props.uploadPicture} title="Upload a new picture">Change picture</span>
        {deleteBtn}
      </div>
    );
  }
  return null;
}

const NoPicture = (props) => {
  if (!props.url) {
    return (
      <div className="no-picture"><span onClick={props.uploadPicture} title="Upload a picture">+</span></div>
    )
  }
  return null
}

class EditPicture extends Component {
  constructor(props) {
    super(props);
    this.state = {
      pictureVisible: false,
    }
    this.showEditPicture = this.showEditPicture.bind(this);
    this.hideEditPicture = this.hideEditPicture.bind(this);
    this.uploadPicture = this.uploadPicture.bind(this);
  }
  showEditPicture() {
    this.setState({
      pictureVisible: true,
    })
  }
  hideEditPicture() {
    this.setState({
      pictureVisible: false,
    })
  }
  uploadPicture() {
    this.inputPicture.click();
  }
  createPicture(e) {
    const file = e.target.files[0];
    var base64;
    var reader = new FileReader();
    reader.onloadend = function () {
      base64 = reader.result;
      console.log(base64);
    };
    if (file) {
      reader.readAsDataURL(file);
    }
    e.preventDefault();
  }
  render() {
    var url = this.props.urlPicture
    return (
      <div className={"profile-picture-" + this.props.className} onMouseEnter={this.showEditPicture} onMouseLeave={this.hideEditPicture}>
        <ButtonEdit
          pictureVisible={this.state.pictureVisible}
          uploadPicture={this.uploadPicture}
          deleteAvailable={this.props.deleteAvailable}
          deletePicture={this.props.deletePicture}
          id={this.props.className}
        />
        <input type="hidden" style={{display:"none"}} name="MAX_FILE_SIZE" value="30000" />
        <input
          ref={(input) => { this.inputPicture = input; }}
          onChange={this.createPicture}
          id={"upload-profile-picture-" + this.props.className}
          name="userfile"
          style={{display:"none"}} type="file" accept=".jpg, .jpeg, .png"/>
        <NoPicture
          url={url}
          uploadPicture={this.uploadPicture}
        />
        <div className="picture-background-profile" style={{ backgroundImage: "url(" + url + ")" }}></div>
      </div>
    )
  }
}

class MyProfile extends Component {
  constructor(props) {
    super(props);
    this.state = {
      variableModal: '1',
    }
    this.clickDeletePicture = this.clickDeletePicture.bind(this);
    this.cancelAction = this.cancelAction.bind(this);
    this.confirmDeletePicture = this.confirmDeletePicture.bind(this);
  }
  cancelAction() {
    this.setState({
      variableModal: '',
    })
  }
  clickDeletePicture(id) {
    this.setState({
      variableModal: id,
    })
  }
  confirmDeletePicture() {
    console.log("Picture " + this.state.variableModal + " deleted");
    this.setState({
      variableModal: '',
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
        <h1>My profile</h1>
        <div className="pictures">
          <EditPicture className="one" urlPicture={profilePictures[2]} />
          <div className="picture-sub-area">
            <EditPicture className="two" urlPicture={profilePictures[2]} deleteAvailable="true" deletePicture={this.clickDeletePicture}/>
            <EditPicture className="three" urlPicture={profilePictures[0]} deleteAvailable="true" deletePicture={this.clickDeletePicture}/>
            <EditPicture className="four" deleteAvailable="true" deletePicture={this.clickDeletePicture}/>
            <EditPicture className="five" urlPicture={null} deleteAvailable="true" deletePicture={this.clickDeletePicture}/>
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
