import React from 'react';
import './List.css';
import Pin from '../../../design/icons/map-pin-64.png';
import Star from '../../../design/icons/star-128.png';
import FilterLogo from '../../../design/icons/filter.svg'

const Item = (props) => {
  return (
    <div className="profile-element">
      <div className="picture-list">
        <div className="picture-list-background" style={{ backgroundImage: "url(" + props.picture + ")" }}></div>
      </div>
      <a href={`/profile/${props.username}`}><span className="name-list" title={`See ${props.firstname}'s profile`}>{props.name}</span></a>
      <span className="age-list">{props.age} year old</span>
      <div className="distance-list">
        <span className="list-value"> {props.distance} km</span>
        <img alt="Distance icon" src={Pin} className="list-icon"/>
      </div>
      <div className="rating-list">
        <span className="list-value"> {Math.round(props.rating * 10) / 10} / 5</span>
        <img alt="Rating icon" src={Star} className="list-icon"/>
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
      <div id="sort-list">
        {/* <div id="sort-elem"> */}
          <a>
            <span title="Select sort by rating">Rating</span>
          </a>
          <a>
            <span title="Select sort by age">Age</span>
          </a>
          <a>
            <span title="Select sort by distance">Distance</span>
          </a>
          <a>
            <span title="Select sort by common tags">Tags</span>
          </a>
          <a>
            <div className="update-btn" id="filter-btn-little" title="Update data with filter">
              <img alt="Filter logo" title="Update data with filter" src={FilterLogo} className="filter-logo-little"/>
            </div>
          </a>
        {/* </div> */}
        {/* <div id="sort-input">

        </div> */}
      </div>
      <div className="limit" style={{ width: "94.25%", margin: "0px", marginBottom: "2.5px" }}></div>
      <div id="view-list">
        {listProfiles}
      </div>
    </div>
  )
}

export default List;
