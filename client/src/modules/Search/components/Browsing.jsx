import React, { Component } from 'react';
import './Browsing.css';
import Slider from 'rc-slider';
import 'rc-slider/assets/index.css';
import AgeLogo from '../../../design/icons/avatar.svg'
import StarLogo from '../../../design/icons/star-128.png'
import DistanceLogo from '../../../design/icons/distance.svg'

const GetGeocode = (location, updateState) => {
    fetch(`https://maps.googleapis.com/maps/api/geocode/json?address=${location}&key=AIzaSyCPhgHvPYOdkj1t5RLcvlRP_sTt6hgK71o`)
    .then(response => {
      return response.json();
    })
    .then(data => {
      if (!data.results[0] || !data.results[0].formatted_address) {
        return;
      }
      const location = data.results[0].geometry.location;
      updateState('lat', location.lat);
      updateState('lng', location.lng);
      updateState('location', '');
      updateState('newLocation', data.results[0].formatted_address)
    })
}

class SliderDynamicBounds extends Component {
  constructor(props) {
    super(props);
    this.state = {
      min: this.props.min,
      max: this.props.max,
      defaultValue: this.props.defaultValue,
      values: this.props.defaultValue,
    };
  }
  onSliderChange = (value) => {
    this.setState({
      values: value,
    })
  }
  onAfterChange = (value) => {
    console.log(this.props.title, value);
  }
  render() {
    const trackStyle = { backgroundColor: "#e63946", borderRadius: 0 };
    const railStyle = { borderRadius: 0 };
    const handleStyle = {
      color: "#2c3e50",
      opacity: 0,
      borderRadius: "0px",
      border: "0.5px solid",
      width: "10px",
      height: "14px",
      marginLeft: "-5px",
    };
    const dotStyle = {
      opacity: 0,
      borderRadius: "0px",
      border: "0.5px solid",
      width: "2px",
      height: "4px",
      marginLeft: "-1px",
      bottom: "0px",
    };
    const activeDotStyle = {
      color: "#ffffff",
      backgroundColor: "#ffffff",
      opacity: 1,
    }
    var div = 1
    if (this.props.div !== undefined && this.props.div > 0) {
      div = this.props.div
    }
    return (
      <div className="range-area" id={this.props.id}>
        <img alt="Range logo" title={this.props.title} src={this.props.logo} className="range-logo"/>
        <label className="range-title">{this.props.title}</label>
        <label className="range-values">{this.state.values[0] / div} to {this.state.values[1] / div} {this.props.unit}</label>
        <div className="range-slider">
          <Slider.Range
            defaultValue={this.state.defaultValue} min={this.state.min} max={this.state.max}
            marks={this.props.marks}
            onChange={this.onSliderChange} onAfterChange={this.onAfterChange}
            trackStyle={[trackStyle]}
            railStyle={railStyle}
            handleStyle={[handleStyle, handleStyle]}
            dotStyle={dotStyle}
            activeDotStyle={activeDotStyle}
          />
        </div>
      </div>
    );
  }
}

class SimpleSlider extends Component {
  constructor(props) {
    super(props);
    this.state = {
      min: this.props.min,
      max: this.props.max,
      defaultValue: this.props.defaultValue,
      value: this.props.defaultValue,
    };
  }
  onSliderChange = (value) => {
    this.setState({
      value,
    });
  }
  onAfterChange = (value) => {
    console.log(this.props.title, value);
  }
  render() {
    const trackStyle = { backgroundColor: "#e63946", borderRadius: 0 };
    const railStyle = { borderRadius: 0 };
    const handleStyle = {
      color: "#2c3e50",
      opacity: 0,
      borderRadius: "0px",
      border: "0.5px solid",
      width: "10px",
      height: "14px",
      marginLeft: "-5px",
    };
    var divider = 1
    if (this.props.div !== undefined && this.props.div > 0) {
      divider = this.props.div
    }
    return (
      <div className="range-area" id={this.props.id}>
        <img alt="Range logo" title={this.props.title} src={this.props.logo} className="range-logo"/>
        <label className="range-title">{this.props.title}</label>
        <label className="range-values">Maximum {this.state.value / divider} {this.props.unit}</label>
        <div className="range-slider">
          <Slider value={this.state.value} min={this.state.min} max={this.state.max}
            onChange={this.onSliderChange} onAfterChange={this.onAfterChange}
            trackStyle={[trackStyle]}
            railStyle={railStyle}
            handleStyle={[handleStyle, handleStyle]}
          />
        </div>
      </div>
    );
  }
}

class Browsing extends Component {
  constructor(props) {
    super(props);
    this.state = {
      newLocation: 'Enter a new location',
      location: '',
      lat: 0,
      lng: 0,
    };
    this.handleUserInput = this.handleUserInput.bind(this);
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

  render() {
    const marksRating = {
      10: '',
      20: '',
      30: '',
      40: '',
      50: '',
    };
    return (
      <div id="browsing">
        <div id="parameters">
          <SliderDynamicBounds id="range-age" min={16} max={100} defaultValue={[22, 28]} logo={AgeLogo} title="Age" unit="year old"/>
          <SliderDynamicBounds id="range-rating" min={1} max={50} defaultValue={[25, 50]} logo={StarLogo} title="Rating" unit="star(s)" div={10} marks={marksRating}/>
          <SimpleSlider id="range-distance" min={0} max={150} defaultValue={50} logo={DistanceLogo} title="Distance" unit="km"/>
          <input id="location" type="text" name="location"
            placeholder={this.state.newLocation} value={this.state.location} onChange={this.handleUserInput}/>
          <span id="send-location" onClick={() => GetGeocode(this.state.location, this.updateState)} title="Search for location">{this.state.location ? '↪' : null}</span>
          <div className="limit" style={{ width: "70%", margin: "2.5px" }}></div>
          <div id="browsing-tags">

          </div>
        </div>
      </div>
    )
  }
}

export default Browsing;
