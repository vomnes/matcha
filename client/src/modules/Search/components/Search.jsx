import React, { Component } from 'react';
import Browsing from './Browsing.jsx'
import DataMap from './DataMap.jsx'
import List from './List.jsx'
import './Search.css';
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
      console.log(response);
      updateState("me", response);
      return;
    }
  }
}

const GetMatch = async (updateState) => {
  let res = await api.match();
  if (res) {
    const response = await res.json();
    if (res.status >= 500) {
      throw new Error("Bad response from server - GetMatch has failed");
    } else if (res.status >= 400) {
      console.log(response.error);
    } else {
      updateState("profiles", response);
      return;
    }
  }
}

class Search extends Component {
  constructor(props) {
    super(props);
    this.state = {
      me: {},
    }
    this.updateState = this.updateState.bind(this);
  }
  updateState(key, content) {
    this.setState({
      [key]: content,
    });
  }
  componentDidMount() {
    GetMe(this.updateState);
    GetMatch(this.updateState);
  }
  render() {
    return (
      <div>
        <Browsing />
        <div id="result-area">
          <DataMap lat={this.state.me.lat} lng={this.state.me.lng} profiles={this.state.profiles}/>
          <List profiles={this.state.profiles}/>
        </div>
      </div>
    )
  }
}

export default Search
