import React, { Component } from 'react';
import GoogleMapReact from 'google-map-react';
import './DataMap.css';
import Pin from '../../../design/icons/maps-and-flags-red.svg';
import MyPin from '../../../design/icons/maps-and-flags-blue.svg';

const PositionMark = (props) => {
    if (props.picture !== undefined) {
      return (
        <div>
          <img alt="Location pin" src={Pin} className="map-pin"/>
          <div className="picture-pin">
            <div className="picture-pin-background" style={{ backgroundImage: "url(" + props.picture.replace("h=1000&q=10", "h=40&q=100") + ")" }}></div>
          </div>
        </div>
      )
    } else {
      return (
        <div>
          <img alt="Location pin" src={MyPin} className="map-pin"/>
        </div>
      )
    }
}

const Map = (props) => {
  if (props.lat && props.lng) {
    const lat = props.lat;
    const lng = props.lng;
    var marks = [];
    if (props.profiles) {
      var index = 0;
      props.profiles.forEach((profile) => {
        marks.push(
          <PositionMark
            lat={profile.latitude}
            lng={profile.longitude}
            picture={profile.picture_url}
            text={profile.username}
            key={index}
          />);
        index++;
      });
    }
    return (
      <GoogleMapReact
        bootstrapURLKeys={{ key: 'AIzaSyCPhgHvPYOdkj1t5RLcvlRP_sTt6hgK71o' }}
        center={{ lat: lat, lng: lng }}
        defaultZoom={12}
      >
      {marks}
      <PositionMark
        lat={lat}
        lng={lng}
        text={'Your location'}
      />
      </GoogleMapReact>
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
        <Map lat={this.props.lat} lng={this.props.lng} profiles={this.props.profiles}/>
      </div>
    )
  }
}

export default DataMap;
