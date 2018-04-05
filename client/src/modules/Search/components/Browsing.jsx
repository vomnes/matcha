import React, { Component } from 'react';
import './Browsing.css';
import Slider from 'rc-slider';
import 'rc-slider/assets/index.css';
import AgeLogo from '../../../design/icons/avatar.svg'
import StarLogo from '../../../design/icons/star-128.png'
import DistanceLogo from '../../../design/icons/distance.svg'
import Downshift from 'downshift'
import FilterLogo from '../../../design/icons/filter.svg'

const GetGeocode = (location, updateState, globalUpdateState) => {
    fetch(`https://maps.googleapis.com/maps/api/geocode/json?address=${location}&key=AIzaSyCPhgHvPYOdkj1t5RLcvlRP_sTt6hgK71o`)
    .then(response => {
      return response.json();
    })
    .then(data => {
      if (!data.results[0] || !data.results[0].formatted_address) {
        return;
      }
      const location = data.results[0].geometry.location;
      globalUpdateState('lat', location.lat);
      globalUpdateState('lng', location.lng);
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
    var div = 1
    if (this.props.div !== undefined && this.props.div > 0) {
      div = this.props.div
    }
    this.props.globalUpdateState(this.props.title, {min: value[0] / div, max: value[1] / div});
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
    this.props.globalUpdateState(this.props.title, {max: value});
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

const MyTag = (props) => {
  return (<span key={props.index} tagid={props.tagid} className="picture-tag matched-tag" onClick={() => props.deleteTag(props.tagName, props.tagid)} title={`Click to delete ${props.tagName}`} style={{ cursor: "pointer" }}>#{props.tagName}</span>)
}

const Tags = (props) => {
  var index = 0;
  var tags = [];
  props.tags.forEach((tag) => {
    tags.push(<MyTag key={index} index={index} tagid={tag.id} tagName={tag.name} deleteTag={props.deleteTag}/>);
    index += 1;
  });
  return (
    <div id="hide-scroll">
      <div id="list-tags">
        {tags}
      </div>
    </div>
  )
}

class BasicAutocomplete extends Component {
  constructor(props) {
    super(props);
    this.state = {
      menuTagsIsOpen: false,
    };
    this.outerClick = this.outerClick.bind(this);
  }
  outerClick = () => {
    this.setState({
        menuTagsIsOpen: false,
        selectedItem: ''
    });
  }
  onChange = (selectedItem, func) => {
    if (selectedItem) {
      this.props.appendTag(selectedItem[0], selectedItem[1]);
    }
    func.clearSelection();
  }
  render() {
    return (
      <Downshift
        onChange={this.onChange}
        onOuterClick={this.outerClick}
        render={({
          getInputProps,
          getItemProps,
          isOpen = this.state.menuTagsIsOpen,
          inputValue,
          selectedItem,
          highlightedIndex,
        }) => (
          <div id="new-tag-area">
            <input {...getInputProps({placeholder: 'Which tag ?', className: 'picture-tag matched-tag', id: 'browsing-new-tag'})} />
            {isOpen ? (
              <div id="list-tag-area">
                {this.props.items
                  .filter(i => !inputValue || i.name.toLowerCase().includes(inputValue.toLowerCase()))
                  .map((item, index) => {
                    let nbTags = this.props.tags.length;
                    for (var i = 0; i < nbTags; i++) {
                      if (this.props.tags[i].name === item.name) {
                        return null;
                      }
                    }
                    if (item.name === selectedItem || index >= 10) {
                      return null;
                    }
                    return (
                      <div
                        {...getItemProps({item: [item.name, item.id], className: 'picture-tag matched-tag', id: 'browsing-new-tag-list'})}
                        key={item.name}
                        style={{
                          backgroundColor:
                            highlightedIndex === index ? '#cb2d3e' : '#e63946',
                        }}
                      >
                        {item.name}
                      </div>
                    )
                  }
                )}
              </div>
            ) : null}
          </div>
        )}
      />
    )
  }
}

class Browsing extends Component {
  constructor(props) {
    super(props);
    this.state = {
      location: '',
      newLocation: 'Enter a new location',
      menuTagsIsOpen: false,
      tags: []
    };
    this.handleUserInput = this.handleUserInput.bind(this);
    this.updateState = this.updateState.bind(this);
    this.appendTag = this.appendTag.bind(this);
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
  deleteTag = (name, id) => {
    const tags = this.state.tags.filter(item => item.name !== name)
    this.updateState('tags', tags);
    this.props.updateState('tagsIds', tags.map((item) => item.id));
  }
  appendTag = (name, id) => {
    const nbTags = this.state.tags.length;
    for (var i = 0; i < nbTags; i++) {
      if (this.state.tags[i].name === name) {
        return;
      }
    }
    var newTags;
    const newElem = { name, id }
    if (this.state.tags === null) {
      newTags = [newElem];
    } else {
      newTags = this.state.tags.concat(newElem);
    }
    this.updateState('tags', newTags);
    this.props.updateState('tagsIds', this.state.tags.map((item) => item.id).concat(id));
  }

  render() {
    const marksRating = {
      10: '',
      20: '',
      30: '',
      40: '',
      50: '',
    };
    if (!this.props.age) {
      return (
        <div id="browsing">
          <div id="parameters">
          </div>
        </div>
      )
    }
    return (
      <div id="browsing">
        <div id="parameters">
          <SliderDynamicBounds id="range-age" min={16} max={100} globalUpdateState={this.props.updateState} defaultValue={[this.props.age - 3, this.props.age + 3]} logo={AgeLogo} title="age" unit="year old"/>
          <SliderDynamicBounds id="range-rating" min={1} max={50} globalUpdateState={this.props.updateState} defaultValue={[25, 50]} logo={StarLogo} title="rating" unit="star(s)" div={10} marks={marksRating}/>
          <SimpleSlider id="range-distance" min={0} max={150} globalUpdateState={this.props.updateState} defaultValue={50} logo={DistanceLogo} title="distance" unit="km"/>
          <input id="location" type="text" name="location"
            placeholder={this.state.newLocation} value={this.state.location} onChange={this.handleUserInput}/>
          <span id="send-location" onClick={() => GetGeocode(this.state.location, this.updateState, this.props.updateState)} title="Search for location">{this.state.location ? '↪' : null}</span>
          <div className="limit" style={{ width: "70%", margin: "2.5px" }}></div>
          <div id="browsing-tags">
            <Tags tags={this.state.tags} deleteTag={this.deleteTag}/>
            <BasicAutocomplete
              items={
                [
                  {
                    id: '1',
                    name: 'apple',
                  },
                  {
                    id: '2',
                    name: 'orange',
                  },
                  {
                    id: '3',
                    name: 'carrot',
                  },
                  {
                    id: '1',
                    name: 'apple1',
                  },
                  {
                    id: '2',
                    name: 'orange2',
                  },
                  {
                    id: '3',
                    name: 'carrot3',
                  },
                  {
                    id: '1',
                    name: 'apple4',
                  },
                  {
                    id: '2',
                    name: 'orange5',
                  },
                  {
                    id: '6',
                    name: 'carrotaaaaaaaaaaaaaaaaaaaaaaaa6',
                  },
                ]}
              appendTag={this.appendTag}
              tags={this.state.tags}
            />
            <div className="update-btn" id="filter-btn-big" title="Update data with filter">
              <img alt="Filter logo" title="Update data with filter" src={FilterLogo} className="filter-logo"/>
              <span className="filter-text">Update</span>
            </div>
          </div>
        </div>
      </div>
    )
  }
}

export default Browsing;
