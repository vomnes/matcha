import React, { Component } from 'react';
import GoogleMapReact from 'google-map-react';
import './Location.css';
import Pin from '../../../design/icons/map-pin-64.png';

const GetGeocode = (address, updateDate) => {
    fetch(`https://maps.googleapis.com/maps/api/geocode/json?address=${address}&key=AIzaSyCPhgHvPYOdkj1t5RLcvlRP_sTt6hgK71o`)
    .then(response => {
      return response.json();
    })
    .then(data => {
      const location = data.results[0].geometry.location;
      updateDate('lat', location.lat);
      updateDate('lng', location.lng);
      updateDate('address', '');
    })
}

const PositionMark = ({ text }) => {
    return (
      <div>
        <img alt="Location pin" src={Pin} className="map-pin"/>
      </div>
    )
}

class Map extends Component {
  render() {
    if (this.props.lat && this.props.lng) {
      const lat = this.props.lat;
      const lng = this.props.lng;
      return (
        <div id="profile-map">
          <GoogleMapReact
            bootstrapURLKeys={{ key: 'AIzaSyCPhgHvPYOdkj1t5RLcvlRP_sTt6hgK71o' }}
            center={{ lat: lat, lng: lng }}
            defaultZoom={11}
          >
          <PositionMark
            lat={lat}
            lng={lng}
            text={'Your location'}
          />
          </GoogleMapReact>
        </div>
      )
    }
    return null;
  }
}

class Location extends Component {
  constructor(props) {
    super(props);
    this.state = {
      address: '',
      lat: '',
      lng: '',
      center: {lat: 48.856614, lng: 2.3522219},
    }
    this.handleUserInput = this.handleUserInput.bind(this);
  }
  handleUserInput = (e) => {
    const fieldName = e.target.name;
    var data = e.target.value;
    this.setState({
      [fieldName]: data,
    });
  }
  updateDate = (fieldName, data) => {
    this.setState({
      [fieldName]: data,
    });
  }
  render () {
    return (
      <div>
        <Map lat={this.state.lat} lng={this.state.lng}/>
        <div className="profile-personal-data">
          <input className="field-input" placeholder="Enter your location address" type="text" name="address"
            value={this.state.address} onChange={this.handleUserInput}/>
          <span id="search-location" onClick={() => GetGeocode(this.state.address, this.updateDate)} title="Search for address location">{this.state.address ? '↪' : null}</span>
        </div>
      </div>
    );
  }
}

export default Location;
