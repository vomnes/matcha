import React, { Component } from 'react';
import GoogleMapReact from 'google-map-react';
import './DataMap.css';
import Pin from '../../../design/icons/map-pin-64.png';

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
        <GoogleMapReact
          bootstrapURLKeys={{ key: 'AIzaSyCPhgHvPYOdkj1t5RLcvlRP_sTt6hgK71o' }}
          center={{ lat: lat, lng: lng }}
          defaultZoom={12}
        >
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
}

class DataMap extends Component {
  constructor(props) {
    super(props);
    this.state = {
    }
    this.updateMapHeight = this.updateMapHeight.bind(this);
  }
  updateMapHeight = () => {
    this.setState({
      mapHeight: window.innerHeight - 228
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
        <Map lat={48.4083868} lng={-4.5696403}/>
      </div>
    )
  }
}

export default DataMap;
