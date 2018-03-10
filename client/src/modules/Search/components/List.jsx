import React, { Component } from 'react';
import './List.css';

const List = () => {
  var profilePictures = [
    require('../../../design/pictures/Register-robson-hatsukami-morgan-250757-unsplash.jpg'),
    require('../../../design/pictures/Profile-molly-belle-73279-unsplash.jpg'),
    require('../../../design/pictures/Login-sorasak-217807-unsplash.jpg'),
  ];
  return (
    <div id="list">
      <span>Profiles availables</span>
      <div id="view-list">
        <div className="profile-element">
          <div className="picture-list">
            <div className="picture-list-background" style={{ backgroundImage: "url(" + profilePictures[0] + ")" }}></div>
          </div>
          <span>Valentin Omnes 23 year old 1km 5/5</span>
        </div>
      </div>
    </div>
  )
}

export default List;
