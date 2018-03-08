import React, { Component } from 'react';
import './Picture.css';

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
    reader.onloadend = () => {
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

const Pictures = (props) => {
 return (
   <div className="pictures">
     <EditPicture className="one" urlPicture={props.profilePictures[2]} />
     <div className="picture-sub-area">
       <EditPicture className="two" urlPicture={props.profilePictures[2]} deleteAvailable="true" deletePicture={props.clickDeletePicture}/>
       <EditPicture className="three" urlPicture={props.profilePictures[0]} deleteAvailable="true" deletePicture={props.clickDeletePicture}/>
       <EditPicture className="four" deleteAvailable="true" deletePicture={props.clickDeletePicture}/>
       <EditPicture className="five" urlPicture={null} deleteAvailable="true" deletePicture={props.clickDeletePicture}/>
     </div>
   </div>
 )
}

export default Pictures;
