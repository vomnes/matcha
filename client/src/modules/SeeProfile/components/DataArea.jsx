import React, { Component } from 'react';
import AgeLogo from '../../../design/icons/age-128.png'
import ConnectedLogo from '../../../design/icons/connected-128.png'
import GenderLogo from '../../../design/icons/gender-128.png'
import HeartLogo from '../../../design/icons/heart-128.png'
import LocationLogo from '../../../design/icons/location-128.png'
import Star from '../../../design/icons/star-128.png'
import './DataArea.css';

const ShowTags = (props) => {
  var index = 0;
  var elements = [];
  props.tags.forEach(function (tag) {
    elements.push(<span key={index} className="picture-tag">#{tag}</span>);
    index += 1;
  });
  return (
    <div id="picture-tags">
      {elements}
    </div>
  )
}

class DataArea extends Component {
  constructor(props) {
    super(props);
    this.state = {
    }
  }
  render() {
    return (
      <div className="data-area">
        <div className="profile-id">
          <span className="profile-title">{this.props.firstname + ' ' + this.props.lastname}</span><br />
          <span className="profile-username"> @{this.props.username}</span><br />
          <div className="profile-bio">
            <span>Greatly hearted has who believe...</span>
          </div>
          <div className="profile-data-list">
            <img alt="Age logo" title="Age" src={AgeLogo}/><span>22 years old</span><br />
            <img alt="Gender logo" title="Gender" src={GenderLogo}/><span>Male</span><br />
            <img alt="Preferred gender logo" title="Preferred gender" src={HeartLogo}/><span>Female</span><br />
            <img alt="Location logo" title="Location" src={LocationLogo}/><span>Paris, France</span><br />
            <img alt="Rating logo" title="Rating" src={Star}/><span>95/100</span><br />
            <img alt="Connected status logo" title="Connected status" src={ConnectedLogo}/>
            <span>{this.props.online ? "Online" : "Offline - Last connection 60 minutes ago"}</span><br />
          </div>
          <ShowTags tags={this.props.tags}/>
        </div>
        {this.props.usersAreConnected ? (
          <div className="profiles-linked">
            <span
              role="img" aria-labelledby="Connected with" title={'You are connected with ' + this.props.firstname + ' - Click here to take contact ;)' }>&#x1F517;</span>
          </div>
        ) : null}
      </div>
    )
  }
}

export default DataArea;
