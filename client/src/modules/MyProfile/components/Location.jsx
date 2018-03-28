import React, { Component } from 'react';
import GoogleMapReact from 'google-map-react';
import './Location.css';
import Pin from '../../../design/icons/map-pin-64.png';
import api from '../../../library/api'

const GetGeocode = (address, updateData, updateState) => {
    fetch(`https://maps.googleapis.com/maps/api/geocode/json?address=${address}&key=AIzaSyCPhgHvPYOdkj1t5RLcvlRP_sTt6hgK71o`)
    .then(response => {
      return response.json();
    })
    .then(data => {
      if (!data.results[0] || !data.results[0].geometry.location) {
        updateData('error', 'Error: Invalid address');
        return;
      }
      const location = data.results[0].geometry.location;
      data.results[0].address_components.forEach((elem) => {
        if (elem.types[0] === "locality") {
          updateData("city", elem.long_name);
        } else if (elem.types[0] === "postal_code") {
          updateData("zip", elem.long_name);
        } else if (elem.types[0] === "country") {
          updateData("country", elem.long_name);
        }
      });
      updateData('lat', location.lat);
      updateData('lng', location.lng);
      updateData('geolocalisation_allowed', true);
      updateData('address', '');
      updateData('newAddress', true);
      updateData('error', '');
    })
}

const UpdateLocation = async (args, updateGlobalState, updateState) => {
  let res = await api.location(args);
  if (res) {
    if (res.status >= 400) {
      const response = await res.json();
      if (res.status >= 400) {
        throw new Error("Bad response from server - UpdateLocation has failed");
      } else if (res.status >= 400) {
        updateGlobalState('newError', 'Location: ' + response.error);
        return
      }
    }
    updateState('newAddress', false);
    return;
  }
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
      lat: this.props.lat,
      lng: this.props.lng,
      geolocalisation_allowed: this.props.geolocalisation_allowed,
      error: '',
      newAddress: false,
    }
    this.handleUserInput = this.handleUserInput.bind(this);
    this.handleSubmit = this.handleSubmit.bind(this);
    this.updateState = this.updateState.bind(this);
  }
  handleUserInput = (e) => {
    const fieldName = e.target.name;
    var data = e.target.value;
    this.setState({
      [fieldName]: data,
    });
  }
  updateState = (fieldName, data) => {
    this.setState({
      [fieldName]: data,
    });
  }
  handleSubmit(e) {
    this.setState({
      newError: '',
    });
    UpdateLocation({
      lat: this.state.lat,
      lng: this.state.lng,
      city: this.state.city,
      zip: this.state.zip,
      country: this.state.country,
    }, this.props.updateState, this.updateState);
    e.preventDefault();
  }
  componentWillReceiveProps() {
    this.updateState('lat', this.props.lat);
    this.updateState('lng', this.props.lng);
    this.updateState('geolocalisation_allowed', this.props.geolocalisation_allowed);
  }
  render () {
    return (
      <div>
        {this.state.geolocalisation_allowed ? (<Map lat={this.state.lat} lng={this.state.lng}/>) : null}
        <form className="profile-personal-data" onSubmit={this.handleSubmit}>
          <input className="field-input" placeholder="Enter your location address" type="text" name="address"
            value={this.state.address} onChange={this.handleUserInput}/>
          <span id="search-location" onClick={() => GetGeocode(this.state.address, this.updateState)} title="Search for address location">{this.state.address ? '↪' : null}</span>
          <span id="map-error">{this.state.error}</span>
          {this.state.newAddress ? (<input className="submit-profile" type="submit" value="Set" title="Save as location"/>) : null}
        </form>
      </div>
    );
  }
}

export default Location;
