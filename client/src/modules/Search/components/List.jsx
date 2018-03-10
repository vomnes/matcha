import React from 'react';
import './List.css';
import Pin from '../../../design/icons/map-pin-64.png';
import Star from '../../../design/icons/star-128.png';

const Item = (props) => {
  return (
    <div className="profile-element">
      <div className="picture-list">
        <div className="picture-list-background" style={{ backgroundImage: "url(" + props.picture + ")" }}></div>
      </div>
      <a href="/profile/vomnes"><span className="name-list" title="See profile">{props.name}</span></a>
      <span className="age-list">{props.age} year old</span>
      <div className="distance-list">
        <img alt="Distance icon" src={Pin} className="list-icon"/>
        <span className="list-value">{props.distance} km</span>
      </div>
      <div className="rating-list">
        <img alt="Rating icon" src={Star} className="list-icon"/>
        <span className="list-value"> {props.rating}/5</span>
      </div>
    </div>
  )
}

const List = () => {
  var profilePictures = [
    require('../../../design/pictures/Register-robson-hatsukami-morgan-250757-unsplash.jpg'),
    require('../../../design/pictures/Profile-molly-belle-73279-unsplash.jpg'),
    require('../../../design/pictures/Login-sorasak-217807-unsplash.jpg'),
  ];
  return (
    <div id="list">
      <div id="view-list">
        <Item picture={profilePictures[0]} name="Valentin Omnes" age={22} distance={10} rating={5}/>
        <Item picture={profilePictures[1]} name="Valentin Omnes" age={22} distance={10} rating={5}/>
        <Item picture={profilePictures[2]} name="Valentin Omnes" age={22} distance={10} rating={5}/>
        <Item picture={profilePictures[0]} name="Valentin Omnes" age={22} distance={10} rating={5}/>
        <Item picture={profilePictures[1]} name="Valentin Omnes" age={22} distance={10} rating={5}/>
        <Item picture={profilePictures[2]} name="Valentin Omnes" age={22} distance={10} rating={5}/>
        <Item picture={profilePictures[0]} name="Valentin Omnes" age={22} distance={10} rating={5}/>
        <Item picture={profilePictures[1]} name="Valentin Omnes" age={22} distance={10} rating={5}/>
        <Item picture={profilePictures[2]} name="Valentin Omnes" age={22} distance={10} rating={5}/>
        <Item picture={profilePictures[0]} name="Valentin Omnes" age={22} distance={10} rating={5}/>
        <Item picture={profilePictures[1]} name="Valentin Omnes" age={22} distance={10} rating={5}/>
        <Item picture={profilePictures[2]} name="Valentin Omnes" age={22} distance={10} rating={5}/>
      </div>
    </div>
  )
}

export default List;
