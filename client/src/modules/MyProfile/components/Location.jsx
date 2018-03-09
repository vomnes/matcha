import React, { Component } from 'react';

const GetGeocode = () => {
    fetch('https://maps.googleapis.com/maps/api/geocode/json?address=Brest&key=AIzaSyCPhgHvPYOdkj1t5RLcvlRP_sTt6hgK71o')
    .then(response => {
      return response.json();
    })
    .then(data => {
      console.log(data.results[0].geometry.location);
    })
}

class Location extends Component {
  render () {
    GetGeocode();
    return null;
  }
}

export default Location;
