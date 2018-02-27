import React, { Component } from 'react';
import './SeeProfile.css'

class SeeProfile extends Component {
  render() {
    const profilePicture = require('../../../design/pictures/Profile-molly-belle-73279-unsplash.jpg');
    return (
      <div>
        <div className="picture-area">
          <div className="picture-user-background" style={{ backgroundImage: "url(" + profilePicture + ")" }}></div>
          <div id="picture-previous">
            <span className="arrow">&#x21A9;</span>
          </div>
          <div id="picture-next">
            <span className="arrow">&#x21AA;</span>
          </div>
        </div>
        <div className="data-area">
          <span>Valentin Omnes - vomnes</span>
        </div>
      </div>
    )
  }
}

export default SeeProfile
