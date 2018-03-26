import React, { Component } from 'react';
import './Picture.css';
import api from '../../../library/api'

const getPicture = async (args) => {
  let res = await api.editpicture(args);
  if (res) {
    const response = await res.json();
    if (response.status >= 500) {
      throw new Error("Bad response from server - getPicture has failed");
    } else if (response.status >= 400) {
      response.json().then(function(data) {
        console.log(data.error);
        return;
      });
    } else {
      console.log(response);
      return;
    }
  }
}

// {number method, base64}

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
  createPicture(props, e) {
    const file = e.target.files[0];
    var base64;
    var reader = new FileReader();
    reader.onloadend = () => {
      base64 = reader.result;
      getPicture({number: props.number, method: `POST`, base64});
    };
    if (file) {
      reader.readAsDataURL(file);
    }
    e.preventDefault();
  }
  render() {
    console.log();
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
          onChange={(e) => this.createPicture(this.props, e)}
          id={"upload-profile-picture-" + this.props.className}
          name="userfile"
          style={{display:"none"}} type="file" accept=".jpg, .jpeg, .png"/>
        <NoPicture
          url={this.props.urlPicture}
          uploadPicture={this.uploadPicture}
        />
        <div className="picture-background-profile" style={{ backgroundImage: "url(" + this.props.urlPicture + ")" }}></div>
      </div>
    )
  }
}

const Pictures = (props) => {
 return (
   <div className="pictures">
     <EditPicture className="one" number="1" urlPicture={props.profilePictures[0]} />
     <div className="picture-sub-area">
       <EditPicture className="two" number="2" urlPicture={props.profilePictures[1]} deleteAvailable="true" deletePicture={props.clickDeletePicture}/>
       <EditPicture className="three" number="3" urlPicture={props.profilePictures[2]} deleteAvailable="true" deletePicture={props.clickDeletePicture}/>
       <EditPicture className="four" number="4" urlPicture={props.profilePictures[3]} deleteAvailable="true" deletePicture={props.clickDeletePicture}/>
       <EditPicture className="five" number="5" urlPicture={props.profilePictures[4]} deleteAvailable="true" deletePicture={props.clickDeletePicture}/>
     </div>
   </div>
 )
}

export default Pictures;
