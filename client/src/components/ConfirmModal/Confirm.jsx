import React, { Component } from 'react';
import PropTypes from 'prop-types';
import './Confirm.css'
import InfoLogo from '../../design/icons/information-128.png'

class ConfirmModal extends Component {
  render() {
    if (this.props.number) {
      var content;
      if (this.props.type === "deletePicture") {
          content = "Delete the picture number " + this.props.number;
      }
      return (
        <div className="modal-confirm-area">
          <div className="modal-confirm">
            <img alt="Information logo" className="modal-confirm-icon" width="32px" src={InfoLogo}/>
            <span className="modal-confirm-content">{content}</span>
            <span className="modal-confirm-cancel" onClick={this.props.cancelAction}>Cancel</span>
            <span className="modal-confirm-confirm" onClick={this.props.confirmAction}>Confirm</span>
          </div>
        </div>
      )
    }
    return null;
  }
}

ConfirmModal.propTypes = {
  number: PropTypes.string,
  cancelAction: PropTypes.func.isRequired,
  confirmAction: PropTypes.func.isRequired,
};

export default ConfirmModal;
