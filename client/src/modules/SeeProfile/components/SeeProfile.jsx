import React, { Component } from 'react';
import './SeeProfile.css'

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

const ProfilePicture = (props) => {
  const index = props.index
  return (
    <div className="picture-area">
      <div className="picture-user-background" style={{ backgroundImage: "url(" + props.picture + ")" }}></div>
      <IndexPictures pictureArrayLength={props.pictureArrayLength} index={index}/>
      <div id="picture-previous" style={{ visibility: (index === 0) ? "hidden" :  "visible" }}>
        <span className="arrow" onClick={() => props.changePicture(0, props.pictureArrayLength)}>&#x21A9;</span>
      </div>
      <div id="picture-next" style={{ visibility: (index === (props.pictureArrayLength - 1)) ? "hidden" :  "visible" } }>
        <span className="arrow" onClick={() => props.changePicture(1, props.pictureArrayLength)}>&#x21AA;</span>
      </div>
    </div>
  )
}

class SeeProfile extends Component {
  constructor(props) {
    super(props);
    this.state = {
      indexProfilePictures: 0,
      lengthProfilePictures: 0,
    }
    this.changePicture = this.changePicture.bind(this);
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
  render() {
    var profilePictures = [
      require('../../../design/pictures/Profile-molly-belle-73279-unsplash.jpg'),
      require('../../../design/pictures/Login-sorasak-217807-unsplash.jpg'),
      require('../../../design/pictures/Register-robson-hatsukami-morgan-250757-unsplash.jpg'),
    ];
    return (
      <div>
        <ProfilePicture
          picture={profilePictures[this.state.indexProfilePictures]}
          changePicture={this.changePicture}
          pictureArrayLength={profilePictures.length}
          index={this.state.indexProfilePictures}
        />
        <div className="data-area">
          <span>Valentin Omnes - vomnes</span>
        </div>
      </div>
    )
  }
}

export default SeeProfile
