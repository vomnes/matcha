import React, { Component } from 'react';
import AgeLogo from '../../../design/icons/age-128.png'
import ConnectedLogo from '../../../design/icons/connected-128.png'
import GenderLogo from '../../../design/icons/gender-128.png'
import HeartLogo from '../../../design/icons/heart-128.png'
import LocationLogo from '../../../design/icons/location-128.png'
import Star from '../../../design/icons/star-128.png'
import './DataArea.css';

var moment = require('moment');

const ShowTags = (props) => {
  var index = 0;
  var sharedTags = [];
  var otherTags = [];
  props.sharedTags.forEach(function (tag) {
    sharedTags.push(<span key={index} className="picture-tag matched-tag">#{tag}</span>);
    index += 1;
  });
  props.otherTags.forEach(function (tag) {
    otherTags.push(<span key={index} className="picture-tag other-tag">#{tag}</span>);
    index += 1;
  });
  return (
    <div >
      {props.sharedTags.length ? (
        <div>
          <span className="tag-title">I like as you</span>
          <div id="data-tags">
            {sharedTags}
          </div>
        </div>
      ) : null}
      {props.otherTags.length ? (
        <div>
          <span className="tag-title">I like</span>
          <div id="data-tags">
            {otherTags}
          </div>
        </div>
      ) : null}
    </div>
  )
}

class DataArea extends Component {
  constructor(props) {
    super(props);
    this.state = {};
    this.updateState = this.updateState.bind(this);
  }
  updateState(key, content) {
    this.setState({
      [key]: content,
    });
  }

  componentWillReceiveProps(nextProps) {
    this.updateState('data', nextProps.data);
  }

  render() {
    var u = this.state.data
    var online;
    if (u) {
      online = u.online ? "Online" : `Offline${u.logout_at ? ` - Last connection ${moment(u.logout_at).startOf("minute").fromNow()}` : ''}`;
    }
    return (
      <div className="data-area">
        <div className="profile-id">
          <span className="profile-title">{(u && u.firstname) || ''} {(u && u.lastname) || ''}</span><br />
          <span className="profile-username"> @{this.props.username}</span><br />
          <div className="profile-bio">
            <span>{(u && u.biography) || ''}</span>
          </div>
          <div className="limit" style={{ width: "10%" }}></div>
          <div className="profile-data-list">
            <img alt="Age logo" title="Age" src={AgeLogo}/><span>{(u && u.age) || ''} years old</span><br />
            <img alt="Gender logo" title="Gender" src={GenderLogo}/><span>{(u && u.genre) || ''}</span><br />
            <img alt="Preferred gender logo" title="Preferred gender" src={HeartLogo}/><span>{(u && u.interesting_in) || ''}</span><br />
            <img alt="Location logo" title="Location" src={LocationLogo}/><span>{(u && u.location) || ''}</span><br />
            <img alt="Rating logo" title="Rating" src={Star}/><span>{(u && u.rating) || ''}/5</span><br />
            <img alt="Connected status logo" title="Connected status" src={ConnectedLogo}/>
            <span>{online}</span><br />
          </div>
          <div className="limit" style={{ width: "10%" }}></div>
          <ShowTags sharedTags={(u && u.tags && u.tags.shared) || []} otherTags={(u && u.tags && u.tags.personal) || []}/>
        </div>
      </div>
    )
  }
}

export default DataArea;
