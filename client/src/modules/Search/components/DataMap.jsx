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

const Map = (props) => {
  if (props.lat && props.lng) {
    const lat = props.lat;
    const lng = props.lng;
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

class DataMap extends Component {
  constructor(props) {
    super(props);
    this.state = {
    }
    this.updateMapHeight = this.updateMapHeight.bind(this);
  }
  updateMapHeight = () => {
    var offsetTop = 213
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
        <Map lat={this.props.lat} lng={this.props.lng}/>
      </div>
    )
  }
}

export default DataMap;
