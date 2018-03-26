import React, { Component } from 'react';
import './Picture.css';
import api from '../../../library/api'

const getPicture = async (args, updatePicture, updateState) => {
  let res = await api.editpicture(args);
  if (res) {
    const json = await res.json();
    if (res.status >= 500) {
      throw new Error("Bad response from server - getPicture has failed");
    } else if (res.status >= 400) {
      updateState('newError', json.error);
      return;
    } else {
      updatePicture('picture_url_'+args.number, json.picture_url);
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
      getPicture({number: props.number, method: `POST`, base64}, props.updatePicture, props.updateState);
    };
    if (file) {
      reader.readAsDataURL(file);
    }
    e.preventDefault();
  }
  render() {
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
     <EditPicture className="one" number="1" urlPicture={"http://localhost:8080" + props.profilePictures[0]} updatePicture={props.updatePicture} updateState={props.updateState}/>
     <div className="picture-sub-area">
       <EditPicture className="two" number="2" urlPicture={"http://localhost:8080" + props.profilePictures[1]} deleteAvailable="true" deletePicture={props.clickDeletePicture} updatePicture={props.updatePicture} updateState={props.updateState}/>
       <EditPicture className="three" number="3" urlPicture={"http://localhost:8080" + props.profilePictures[2]} deleteAvailable="true" deletePicture={props.clickDeletePicture} updatePicture={props.updatePicture} updateState={props.updateState}/>
       <EditPicture className="four" number="4" urlPicture={"http://localhost:8080" + props.profilePictures[3]} deleteAvailable="true" deletePicture={props.clickDeletePicture} updatePicture={props.updatePicture} updateState={props.updateState}/>
       <EditPicture className="five" number="5" urlPicture={"http://localhost:8080" + props.profilePictures[4]} deleteAvailable="true" deletePicture={props.clickDeletePicture} updatePicture={props.updatePicture} updateState={props.updateState}/>
     </div>
   </div>
 )
}

export default Pictures;
