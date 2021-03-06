import React from 'react';
import './List.css';
import Pin from '../../../design/icons/map-pin-64.png';
import Star from '../../../design/icons/star-128.png';
import FilterLogo from '../../../design/icons/filter.svg'
import UpBlack from '../../../design/icons/caret-arrow-up-black.svg'
import UpRed from '../../../design/icons/caret-arrow-up-red.svg'
import DownBlack from '../../../design/icons/sort-down-black.svg'
import DownRed from '../../../design/icons/sort-down-red.svg'
import utils from '../../../library/utils/pictures.js'

const Item = (props) => {
  return (
    <div>
      <a href={`/profile/${props.username}` + (props.optionsBase64 ? '/' + props.optionsBase64 : '')} title={`Click to see ${props.firstname}'s profile`}>
        <div className="profile-element" id={props.username} style={ props.isSelectedOnMap ? { backgroundColor: "#EAEAEA" } : null }>
          <div className="picture-list">
            <div className="picture-list-background" style={{ backgroundImage: "url(" + utils.pictureURLFormated(props.picture, "h=100&q=100") + ")" }}></div>
          </div>
          <span className="name-list">{props.name}</span>
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
      </a>
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
        optionsBase64={props.optionsBase64}
      />
    );
      index++;
    });
    showData = (
      <div id="view-list">
        {listProfiles}
        {!props.allDataCollected ? (<span id="more-data" title="Click to load more data" onClick={props.loadMoreData}>Load more profiles</span>) : (<span id="no-more-data" title="No more profiles">&#8226;&#8226;&#8226;</span>)}
      </div>)
  } else {
    showData = (<span id="no-data" title="No profiles with the search">&#8226;&#8226;&#8226;</span>)
  }
  var listResponsiveStyle;
  if (window.innerWidth <= 650) {
    listResponsiveStyle = {height: "100px"}
    if (window.innerHeight >  500) {
      listResponsiveStyle = {height: (window.innerHeight - 393) + "px"}
    }
  }
  return (
    <div id="list" style={listResponsiveStyle}>
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
