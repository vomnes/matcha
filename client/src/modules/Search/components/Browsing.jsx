import React, { Component } from 'react';
import './Browsing.css';
import Slider from 'rc-slider';
import 'rc-slider/assets/index.css';
import AgeLogo from '../../../design/icons/avatar.svg'
import StarLogo from '../../../design/icons/star-128.png'
import DistanceLogo from '../../../design/icons/distance.svg'

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

const Browsing = () => {
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
      </div>
    </div>
  )
}

export default Browsing;
