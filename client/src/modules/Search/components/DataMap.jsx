import React, { Component } from 'react';
import GoogleMap from 'google-map-react';
import './DataMap.css';
import Pin from '../../../design/icons/maps-and-flags-red.svg';
import MyPin from '../../../design/icons/maps-and-flags-blue.svg';
import mapStyle from './styleMap.js';
import utils from '../../../library/utils/pictures.js'

const PositionMark = (props) => {
  const selectProfile = () => {
    props.updateState('selectProfile', props.username);
    document.getElementById(props.username).scrollIntoView({ behavior: "smooth", inline: "center" });
  }
  if (props.picture !== undefined) {
    return (
      <div title={`See ${props.firstname} ${props.lastname}`} className="map-mark" onClick={selectProfile}>
        <img alt="Location pin" src={Pin} className="map-pin"/>
        <div className="picture-pin">
          <div className="picture-pin-background" style={{ backgroundImage: "url(" + utils.pictureURLFormated(props.picture) + ")" }}></div>
        </div>
      </div>
    )
  } else {
    return (
      <div title="Your location" style={{ cursor: "default" }}>
        <img alt="Pin" src={MyPin} className="map-pin"/>
      </div>
    )
  }
}

const Map = (props) => {
  if (props.lat && props.lng) {
    const lat = props.lat;
    const lng = props.lng;
    const mapOptions = {
      styles: mapStyle,
    };
    var marks = [];
    if (props.profiles && props.profiles.length > 0) {
      var index = 0;
      props.profiles.forEach((profile) => {
        marks.push(
          <PositionMark
            lat={profile.latitude}
            lng={profile.longitude}
            picture={profile.picture_url}
            username={profile.username}
            firstname={profile.firstname}
            lastname={profile.lastname}
            key={index}
            updateState={props.updateState}
          />);
        index++;
      });
    }
    return (
      <GoogleMap
        bootstrapURLKeys={{ key: 'AIzaSyCPhgHvPYOdkj1t5RLcvlRP_sTt6hgK71o' }}
        center={{ lat: lat, lng: lng }}
        defaultZoom={12}
        options={mapOptions}
      >
      {marks}
      <PositionMark
        lat={lat}
        lng={lng}
        text={'Your location'}
      />
      </GoogleMap>
    )
  }
  return null;
}

class DataMap extends Component {
  constructor(props) {
    super(props);
    this.state = {
    }
    this.updateMapHeight = this.updateMapHeight.bind(this);
  }
  updateMapHeight = () => {
    var offsetTop = 163;
    if (this.instance) {
      offsetTop  = this.instance.getBoundingClientRect().top;
    }
    this.setState({
      mapHeight: window.innerHeight - offsetTop - 15
    });
  }
  componentWillMount() {
      this.updateMapHeight();
  }
  componentDidMount() {
      window.addEventListener("resize", this.updateMapHeight);
  }
  render() {
    return (
      <div id="data-map" style={{ height: window.innerWidth > 650 ? this.state.mapHeight : 200 }}>
        <Map lat={this.props.lat} lng={this.props.lng} profiles={this.props.profiles} updateState={this.props.updateState}/>
      </div>
    )
  }
}

export default DataMap;
