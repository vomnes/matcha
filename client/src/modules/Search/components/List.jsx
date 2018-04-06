import React from 'react';
import './List.css';
import Pin from '../../../design/icons/map-pin-64.png';
import Star from '../../../design/icons/star-128.png';
import FilterLogo from '../../../design/icons/filter.svg'
import UpBlack from '../../../design/icons/caret-arrow-up-black.svg'
import UpRed from '../../../design/icons/caret-arrow-up-red.svg'
import DownBlack from '../../../design/icons/sort-down-black.svg'
import DownRed from '../../../design/icons/sort-down-red.svg'

const Item = (props) => {
  return (
    <div>
      <div className="profile-element" id={props.username} style={ props.isSelectedOnMap ? { backgroundColor: "#EAEAEA" } : null }>
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
      <div className="limit" style={{ width: "94.25%", margin: "2.5px 0px 2.5px" }}></div>
    </div>
  )
}

const List = (props) => {
  var listProfiles = [];
  var showData;
  if (props.profiles && props.profiles.length > 0) {
    var index = 0;
    props.profiles.forEach((profile) => {
      listProfiles.unshift(<Item
        key={index}
        picture={profile.picture_url}
        firstname={profile.firstname}
        name={`${profile.firstname} ${profile.lastname}`}
        age={profile.age}
        distance={profile.distance}
        rating={profile.rating}
        username={profile.username}
        isSelectedOnMap={props.selectProfile === profile.username}
      />
    );
      index++;
    });
    showData = (
      <div id="view-list">
        {listProfiles}
        {!props.allDataCollected ? (<span id="more-data" onClick={props.loadMoreData}>Load more data</span>) : (<span id="no-more-data">&#8226;&#8226;&#8226;</span>)}
      </div>)
  } else {
    showData = (<span id="no-data">&#8226;&#8226;&#8226;</span>)
  }
  return (
    <div id="list">
      <div id="sort-list">
        <a className="sort-title" onClick={() => props.updateState("sort_type", "rating")} style={props.sortType === "rating" ? { borderRadius: "2px", border: "0.5px solid #e63946"} : null }>
          <span title="Sort by rating">Rating</span>
        </a>
        <a className="sort-title" onClick={() => props.updateState("sort_type", "age")} style={props.sortType === "age" ? { borderRadius: "2px", border: "0.5px solid #e63946"} : null }>
          <span title="Sort by age">Age</span>
        </a>
        <a className="sort-title" onClick={() => props.updateState("sort_type", "distance")} style={props.sortType === "distance" ? { borderRadius: "2px", border: "0.5px solid #e63946"} : null }>
          <span title="Sort by distance">Distance</span>
        </a>
        <a className="sort-title" onClick={() => props.updateState("sort_type", "common_tags")} style={props.sortType === "common_tags" ? { borderRadius: "2px", border: "0.5px solid #e63946"} : null }>
          <span title="Sort by common tags">Tags</span>
        </a>
        <a id="up-down-area">
          <div>
            <img alt="Filter logo" title="Ascendant sort"
              src={props.sortDirection === "normal" ? UpRed : UpBlack }
              onClick={() => props.updateState("sort_direction", 'normal')}
              className="up-down-sort"/>
            <img alt="Filter logo" title="Descendant sort"
              src={props.sortDirection === "reverse" ? DownRed : DownBlack }
              onClick={() => props.updateState("sort_direction", 'reverse')}
              className="up-down-sort" style={{ top: "8px" }}/>
          </div>
        </a>
        <a style={{ position: "absolute", paddingLeft: "20px" }}>
          <div className="update-btn" id="filter-btn-little" title="Update data with filter" onClick={props.searchProfiles}>
            <img alt="Filter logo" title="Update data with filter" src={FilterLogo} className="filter-logo-little"/>
          </div>
        </a>
      </div>
      <div className="limit" style={{ width: "94.25%", margin: "0px", marginBottom: "2.5px" }}></div>
      {showData}
    </div>
  )
}

export default List;
