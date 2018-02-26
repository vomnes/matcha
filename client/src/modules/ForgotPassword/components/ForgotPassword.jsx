import React, { Component } from 'react';
import Modal from '../../../components/Modal'
import './ForgotPassword.css'
import api from '../../../library/api'

const sendForgotPasswordEmail = (email, createError, crateSuccess) => {
  api.forgotpassword({
      email,
  }).then(function(response) {
    if (response.status >= 500) {
      throw new Error("Bad response from server");
    } else if (response.status >= 400) {
      response.json().then(function(data) {
        createError(data.error);
        return;
      });
    } else {
      crateSuccess("Reset password email sent to " + email + ".")
      return;
    }
  })
}

class ForgotPassword extends Component {
  constructor(props) {
    super(props);
    if (localStorage.getItem(`matcha_token`)) {
      this.props.history.push("/home");
    }
    this.state = {
      email: '',
      newError: '',
      newSuccess: '',
    }
    this.handleUserInput = this.handleUserInput.bind(this);
    this.createError = this.createError.bind(this);
    this.createSuccess = this.createSuccess.bind(this);
    this.handleSubmit = this.handleSubmit.bind(this);
    this.closeModal = this.closeModal.bind(this);
  }

  handleUserInput (e) {
    this.setState({
      [e.target.name]: e.target.value,
    });
  }
  createError(content) {
    this.setState({
      newError: content,
      newSuccess: ''
    });
  }
  createSuccess(content) {
    this.setState({
      newSuccess: content,
      newError: ''
    });
  }
  closeModal(event) {
    this.setState({
      newError: '',
      newSuccess: ''
    });
    event.preventDefault();
  }
  handleSubmit(e) {
    this.setState({
      newError: '',
    });
    sendForgotPasswordEmail(
      this.state.email,
      this.createError,
      this.createSuccess
    )
    e.preventDefault();
  }
  render() {
    return (
      <div id="background-forgotpassword">
        <div id="box-center-top-form">
          <h1 className="up-title-form">Matcha</h1>
          <h2 className="title-form" style={{ fontSize: "35px" }}>Forgot your password ?</h2>
          <h3 className="sub-title-form" style={{ fontSize: "20px" }}>Enter your email address and we will send you a link to reset your password.</h3>
          <form onSubmit={this.handleSubmit}>
            <input className="input-form" id="placeholder-icon-email" placeholder="Email address" type="email" name="email"
              value={this.state.email} onChange={this.handleUserInput} required/><br />
            <input className="submit-form" type="submit" value="Send"/>
            <div className="form-link">
              <a href='/login' className="form-link-click"><span>Go back on login page</span></a>
            </div>
          </form>
        </div>
        <Modal type="error" content={this.state.newError} onClose={this.closeModal}/>
        <Modal type="success" content={this.state.newSuccess} onClose={this.closeModal}/>
      </div>
    )
  }
}

export default ForgotPassword
