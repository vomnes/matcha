import React, { Component } from 'react';
import PropTypes from 'prop-types';
import './Modal.css'

class Modal extends Component {
  render() {
    if (!this.props.content) {
      return null;
    }
    if (this.props.type === 'error') {
      return (
        <div className="modal-popup error-color" style={{ top: this.props.online ? "-60px" :  "-5px" } }>
          <span className="modal-warning">&#9888;</span>
          <span>Error : {this.props.content}</span>
          <span className="modal-close" onClick={this.props.onClose}>&times;</span>
        </div>
      )
    } else {
      return (
        <div className="modal-popup success-color" style={{ top: this.props.online ? "-60px" :  "-5px" } }>
          <span>{this.props.content}</span>
          <span className="modal-close" onClick={this.props.onClose}>&times;</span>
        </div>
      )
    }
  }
}

Modal.propTypes = {
  type: PropTypes.string,
  onClose: PropTypes.func.isRequired,
  content: PropTypes.string
};

export default Modal;
