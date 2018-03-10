import React, { Component } from 'react';
import Browsing from './Browsing.jsx'
import DataMap from './DataMap.jsx'
import List from './List.jsx'
import './Search.css';

class Search extends Component {
  constructor(props) {
    super(props);
    this.state = {
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
  render() {
    return (
      <div>
        <Browsing />
        <div id="result-area">
          <DataMap />
          <List />
        </div>
      </div>
    )
  }
}

export default Search
