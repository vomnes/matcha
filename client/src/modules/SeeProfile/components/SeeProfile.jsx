import React, { Component } from 'react';
import './SeeProfile.css'
import Modal from '../../../components/Modal'

const IndexPictures = (props) => {
  var elements = [];
  for (var i = 0; i < props.pictureArrayLength; i++) {
    let color = props.index === i ? "black" : "grey";
    elements.push(<span key={i} style={{ color: color }}>&#8226;</span>);
  }
  return (
    <div id="picture-index">
      {elements}
    </div>
  )
}

class ProfilePicture extends Component {
  constructor(props) {
    super(props);
    this.state = {
      moreInformationOpen: false,
      reportedAsFakeAccount: this.props.reportedAsFakeAccount
    }
    this.openInformation = this.openInformation.bind(this);
  }
  openInformation() {
    this.setState({
      moreInformationOpen: !this.state.moreInformationOpen,
    });
  }
  render() {
    const index = this.props.index;
    return (
      <div className="picture-area">
        <div className="picture-user-background" style={{ backgroundImage: "url(" + this.props.picture + ")" }}></div>
        <div className="more-information">
          <span className="more" onClick={this.openInformation}>{this.state.moreInformationOpen ? '-' : '+'}</span>
        </div>
        <div className="information-area" style={{ visibility: this.state.moreInformationOpen ? "visible" :  "hidden" } }>
            {!this.state.reportedAsFakeAccount ?
              <span onClick={() => this.props.updateState("reportedAsFakeAccount", "You have just invalide you fake account report.")}>Invalidate fake account report</span> :
              <span onClick={() => this.props.updateState("reportedAsFakeAccount", "This profile has been declared as fake account.")}>Report as a fake account</span>
            }<br />
          <span>Unlike profile</span>
        </div>
        <IndexPictures pictureArrayLength={this.props.pictureArrayLength} index={index}/>
        <div id="picture-previous" style={{ visibility: (index === 0) ? "hidden" :  "visible" }}>
          <span className="arrow" onClick={() => this.props.changePicture(0, this.props.pictureArrayLength)}>&#x21A9;</span>
        </div>
        <div id="picture-next" style={{ visibility: (index === (this.props.pictureArrayLength - 1)) ? "hidden" :  "visible" } }>
          <span className="arrow" onClick={() => this.props.changePicture(1, this.props.pictureArrayLength)}>&#x21AA;</span>
        </div>
        <div title="Like profile" className="btn-like"
          onClick={() => this.props.likeUser()}
          style={{
            background: (this.props.liked ? "white" :  "#F80759"),
            color: (this.props.liked ? "#F80759" :  "white") }}>
          <span className="like-heart">&#9829;</span>
        </div>
      </div>
    )
  }
}

class SeeProfile extends Component {
  constructor(props) {
    super(props);
    this.state = {
      indexProfilePictures: 0,
      lengthProfilePictures: 0,
      liked: false,
      reportedAsFakeAccount: false,
      newSuccess: '',
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
    console.log(key, !this.state[key])
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
    return (
      <div>
        <ProfilePicture
          picture={profilePictures[this.state.indexProfilePictures]}
          changePicture={this.changePicture}
          pictureArrayLength={profilePictures.length}
          index={this.state.indexProfilePictures}
          liked={this.state.liked}
          likeUser={this.likeUser}
          updateState={this.updateState}
          reportedAsFakeAccount={this.reportedAsFakeAccount}
        />
        <div className="data-area">
          <span>Valentin Omnes - vomnes</span>
        </div>
        <Modal type="success" online="true" content={this.state.newSuccess} onClose={this.closeModal}/>
      </div>
    )
  }
}

export default SeeProfile
