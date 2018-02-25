import React, { Component } from 'react';
import PropTypes from 'prop-types';
import './Error.css'

class Error extends Component {
  render() {
    if (!this.props.content) {
      return null;
    }
    return (
      <div className="error-popup">
        <span className="error-warning">&#9888;</span>
        <span>Error : {this.props.content}</span>
        <span className="error-close" onClick={this.props.onClose}>&times;</span>
      </div>
    )
  }
}

Error.propTypes = {
  onClose: PropTypes.func.isRequired,
  content: PropTypes.string
};

export default Error;
