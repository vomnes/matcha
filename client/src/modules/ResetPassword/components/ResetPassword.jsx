import React, { Component } from 'react';
import Modal from '../../../components/Modal'
import './ResetPassword.css'
// import api from '../../../library/api'

// const resetPassword = (email, createError, crateSuccess) => {
//   api.forgotpassword({
//       email,
//   }).then(function(response) {
//     if (response.status >= 500) {
//       throw new Error("Bad response from server");
//     } else if (response.status >= 400) {
//       response.json().then(function(data) {
//         createError(data.error);
//         return;
//       });
//     } else {
//       crateSuccess("Reset password email sent to " + email + ".")
//       return;
//     }
//   })
// }

class ResetPassword extends Component {
  constructor(props) {
    super(props);
    if (localStorage.getItem(`matcha_token`)) {
      this.props.history.push("/home");
    }
    this.state = {
      password: '',
      rePassword: '',
      newError: '',
      newSuccess: '',
    }
    this.handleUserInput = this.handleUserInput.bind(this);
    this.createError = this.createError.bind(this);
    this.createSuccess = this.createSuccess.bind(this);
    this.handleSubmit = this.handleSubmit.bind(this);
    this.closeModal = this.closeModal.bind(this);
    console.log(this.props.match.params.token);
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
    // resetPassword(
    //   this.state.password,
    //   this.state.rePassword,
    //   this.createError,
    //   this.createSuccess
    // )
    e.preventDefault();
  }
  render() {
    return (
      <div id="background-forgotpassword">
        <div id="box-center-top-form">
          <h1 className="up-title-form">Matcha</h1>
          <h2 className="title-form">Reset password</h2>
          <h3 className="sub-title-form" style={{ fontSize: "20px" }}>Enter your new password.</h3>
          <form onSubmit={this.handleSubmit}>
            <input className="input-form" id="placeholder-icon-password" placeholder="Password" type="password" name="password"
              value={this.state.password} onChange={this.handleUserInput}/><br />
            <input className="input-form" id="placeholder-icon-re-password" placeholder="Re-enter password" type="password" name="rePassword"
              value={this.state.rePassword} onChange={this.handleUserInput}/><br />
            <input className="submit-form" type="submit" value="Reset"/>
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

export default ResetPassword