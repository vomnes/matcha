import React from 'react';
import './List.css';
import Pin from '../../../design/icons/map-pin-64.png';
import Star from '../../../design/icons/star-128.png';

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
          <a href="/profile/vomnes"><span className="name-list">Valentin Omnes</span></a>
          <div className="distance-list">
            <img alt="Distance icon" src={Pin} className="list-icon"/>
            <span className="list-value">10 km</span>
          </div>
          <div className="rating-list">
            <img alt="Rating icon" src={Star} className="list-icon"/>
            <span className="list-value"> 5/5</span>
          </div>
        </div>
      </div>
    </div>
  )
}

export default List;
