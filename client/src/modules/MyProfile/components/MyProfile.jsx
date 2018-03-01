import React, { Component } from 'react';
import './MyProfile.css';

class MyProfile extends Component {
  constructor(props) {
    super(props);
    this.state = {
      pictureVisible: {
        one: false,
      },
    }
    this.showEditPicture = this.showEditPicture.bind(this);
  }
  showEditPicture() {
    let pictureVisible = Object.assign({}, this.state.pictureVisible);
    pictureVisible.one = !pictureVisible.one;
    this.setState({
      pictureVisible
    })
  }
  render() {
    var profilePictures = [
      require('../../../design/pictures/Register-robson-hatsukami-morgan-250757-unsplash.jpg'),
      require('../../../design/pictures/Profile-molly-belle-73279-unsplash.jpg'),
      require('../../../design/pictures/Login-sorasak-217807-unsplash.jpg'),
    ];
    let pictureVisible = Object.assign({}, this.state.pictureVisible);
    return (
      <div>
        <h1>My profile</h1>
        <div className="pictures">
          <div className="profile-picture-one" onMouseOver={this.showEditPicture} onMouseOut={this.showEditPicture}>
            {pictureVisible.one ? <span id="over">Change picture</span> : <span id="over">Change picture</span>}
            <div className="picture-background-profile" style={{ backgroundImage: "url(" + profilePictures[0] + ")" }}></div>
          </div>
          <div className="picture-sub-area">
            <div className="profile-picture-two">
              <div className="picture-background-profile" style={{ backgroundImage: "url(" + profilePictures[2] + ")" }}></div>
            </div>
            <div className="profile-picture-three">
              <div className="picture-background-profile" style={{ backgroundImage: "url(" + profilePictures[1] + ")" }}></div>
            </div>
            <div className="profile-picture-four">
              <div className="picture-background-profile" style={{ backgroundImage: "url(" + profilePictures[2] + ")" }}></div>
            </div>
            <div className="profile-picture-five">
              <div className="picture-background-profile" style={{ backgroundImage: "url(" + profilePictures[2] + ")" }}></div>
            </div>
          </div>
        </div>
      </div>
    )
  }
}

export default MyProfile
