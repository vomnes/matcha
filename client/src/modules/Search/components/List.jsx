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
      <a href={`/profile/${props.username}`}><span className="name-list" title={`See ${props.firstname}'s profile`}>{props.name}</span></a>
      <span className="age-list">{props.age} year old</span>
      <div className="distance-list">
        <img alt="Distance icon" src={Pin} className="list-icon"/>
        <span className="list-value"> {props.distance} km</span>
      </div>
      <div className="rating-list">
        <img alt="Rating icon" src={Star} className="list-icon"/>
        <span className="list-value"> {Math.round(props.rating * 100) / 100} / 5</span>
      </div>
    </div>
  )
}

const List = (props) => {
  var listProfiles = [];
  if (props.profiles) {
    var index = 0;
    props.profiles.forEach((profile) => {
      listProfiles.push(<Item
        key={index}
        picture={profile.picture_url}
        firstname={profile.firstname}
        name={`${profile.firstname} ${profile.lastname}`}
        age={profile.age}
        distance={profile.distance}
        rating={profile.rating}
        username={profile.username}/>);
      index++;
    });
  }
  return (
    <div id="list">
      <div id="view-list">
        {listProfiles}
      </div>
    </div>
  )
}

export default List;
