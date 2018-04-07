import React, { Component } from 'react';
import Browsing from './Browsing.jsx'
import DataMap from './DataMap.jsx'
import List from './List.jsx'
import api from '../../../library/api'

const GetMe = async (updateState) => {
  let res = await api.me();
  if (res) {
    const response = await res.json();
    if (res.status >= 500) {
      throw new Error("Bad response from server - GetMe has failed");
    } else if (res.status >= 400) {
      console.log(response.error);
    } else {
      updateState("me", response);
      return;
    }
  }
}

const GetMatch = async (updateState, optionsBase64, profiles) => {
  let res = await api.match(optionsBase64);
  if (res) {
    const response = await res.json();
    if (res.status >= 500) {
      throw new Error("Bad response from server - GetMatch has failed");
    } else if (res.status >= 400) {
      console.log(response.error);
    } else {
      if (response.data === "No (more) users") {
        if (!profiles) {
          updateState("profiles", []);
        }
        updateState("allDataCollected", true);
      } else {
        if (profiles) {
          updateState("profiles", response.concat(profiles));
        } else {
          updateState("profiles", response);
        }
        if (response.length < 20) {
          updateState("allDataCollected", true);
        } else {
          updateState("allDataCollected", false);
        }
      }
      return;
    }
  }
}

const GetTags = async (updateState) => {
  let res = await api.existingTags();
  if (res) {
    const response = await res.json();
    if (res.status >= 500) {
      throw new Error("Bad response from server - GetMatch has failed");
    } else if (res.status >= 400) {
      console.log(response.error);
    } else {
      updateState("existingTags", response);
      return;
    }
  }
}

class Search extends Component {
  constructor(props) {
    super(props);
    this.state = {
      me: {},
      profiles: {},
      age: {
       min: 0,
       max: 0
      },
      rating: {
       min: 0.0,
       max: 0.0
      },
      distance: {
       max: 0
      },
      existingTags: [],
      tagsIds: [],
      lat: 0.0,
      lng: 0.0,
      sort_type: "rating",   // age, rating, distance, common_tags
      sort_direction: "normal", // normal or reverse
      start_position: 0,
      finish_position: 20,
      optionsBase64: '',
      allDataCollected: false,
      selectProfile: '',
    }
    this.updateState = this.updateState.bind(this);
    this.searchProfiles = this.searchProfiles.bind(this);
  }
  updateState(key, content) {
    this.setState({
      [key]: content,
    });
  }
  optionsToBase64 = (start_position, finish_position) => {
    const options = {
      age: this.state.age,
      rating: this.state.rating,
      distance: this.state.distance,
      tags: this.state.tagsIds,
      lat: this.state.lat,
      lng: this.state.lng,
      sort_type: this.state.sort_type,
      sort_direction: this.state.sort_direction,
      start_position,
      finish_position,
    }
    return Buffer.from(JSON.stringify(options)).toString("base64");
  }
  searchProfiles = (type) => {
    this.updateState('start_position', 0);
    this.updateState('finish_position', 20);
    const objJSONBase64 = this.optionsToBase64(0, 20);
    this.updateState('optionsBase64', objJSONBase64);
    GetMatch(this.updateState, objJSONBase64);
    if (this.state.lat && this.state.lng) {
      var me = this.state.me;
      me['lat'] = this.state.lat;
      me['lng'] = this.state.lng;
      this.updateState('me', me);
    }
  }
  loadMoreData = () => {
    this.updateState('start_position', this.state.start_position + 20);
    this.updateState('finish_position', this.state.finish_position + 20);
    const objJSONBase64 = this.optionsToBase64(this.state.start_position + 20, this.state.finish_position + 20);
    this.updateState('optionsBase64', objJSONBase64);
    GetMatch(this.updateState, objJSONBase64, this.state.profiles);
  }

  componentDidMount() {
    GetMe(this.updateState);
    GetMatch(this.updateState);
    GetTags(this.updateState);
  }
  render() {
    return (
      <div>
        <Browsing existingTags={this.state.existingTags}age={this.state.me.age} updateState={this.updateState} searchProfiles={this.searchProfiles}/>
        <div id="result-area" style={{ position: "relative", marginTop: "10px" }}>
          <DataMap lat={this.state.me.lat} lng={this.state.me.lng} profiles={this.state.profiles} updateState={this.updateState}/>
          <List profiles={this.state.profiles}
            sortType={this.state.sort_type} sortDirection={this.state.sort_direction}
            updateState={this.updateState} searchProfiles={this.searchProfiles}
            loadMoreData={this.loadMoreData}
            allDataCollected={this.state.allDataCollected}
            selectProfile={this.state.selectProfile}
            optionsBase64={this.state.optionsBase64}
          />
        </div>
      </div>
    )
  }
}

export default Search
