import React, { Component } from 'react';
import Modal from '../../../components/Modal'
import './Register.css'
import api from '../../../library/api'
import utils from '../../../library/utils/input.js'

const signUp = (username, firstname, lastname, email, password, rePassword, createError, redirectLogin) => {
  api.register({
      username,
      firstname,
      lastname,
      email,
      password,
      rePassword
  }).then(function(response) {
    if (response.status >= 500) {
      throw new Error("Bad response from server");
    } else if (response.status >= 400) {
      response.json().then(function(data) {
        createError(data.error);
        return;
      });
    } else {
      redirectLogin();
      return;
    }
  })
}

const formatInput = (fieldName, data) => {
  if (fieldName === "username") {
    return utils.blockForbiddenKeys(data, /[0-9a-zA-Z.\-_]/i, 64);
  } else if (fieldName === "firstname" || fieldName === "lastname") {
    data = utils.formatName(data);
    return utils.blockForbiddenKeys(data, /[a-zA-Z-]/i, 64);
  } else if (fieldName === "email") {
    data = data.toLowerCase();
    return utils.blockForbiddenKeys(data, /[a-zA-Z@.]/i, 254);
  }
}

class Register extends Component {
  constructor (props) {
    super(props);
    if (localStorage.getItem(`matcha_token`)) {
      this.props.history.push("/home");
    }
    this.state = {
      username: '',
      firstname: '',
      lastname: '',
      email: '',
      password: '',
      rePassword: '',
      newError: '',
    }
    this.handleUserInput = this.handleUserInput.bind(this);
    this.handleSubmit = this.handleSubmit.bind(this);
    this.createError = this.createError.bind(this);
    this.closeError = this.closeError.bind(this);
    this.redirectLogin = this.redirectLogin.bind(this);
  }

  handleUserInput (e) {
    const fieldName = e.target.name
    var data = e.target.value
    data = formatInput(fieldName, data)
    if (data === -1) {
      return;
    }
    this.setState({
      [fieldName]: data,
    });
  }
  createError(content) {
    this.setState({
      newError: content,
    });
  }
  closeError(event) {
    this.setState({
      newError: '',
    });
    event.preventDefault();
  }
  redirectLogin() {
    this.props.history.push("/login?code=1");
  }
  handleSubmit(e) {
    this.setState({
      newError: '',
    });
    signUp(
      this.state.username,
      this.state.firstname,
      this.state.lastname,
      this.state.email,
      this.state.password,
      this.state.rePassword,
      this.createError,
      this.redirectLogin
    )
    e.preventDefault();
  }
  render() {
    return (
      <div>
        <div id="background-register"></div>
        <a href="/login"><p id="member">Already a member ?</p></a>
        <div id="box-left-side-form">
          <h1 className="up-title-form">Matcha</h1>
          <h2 className="title-form">Register</h2>
          <h3 className="sub-title-form">A simple dating website</h3>
          <form onSubmit={this.handleSubmit}>
            <input className="input-form" id="placeholder-icon-username" placeholder="Username" type="text" name="username"
              pattern="[a-zA-Z0-9\.\-_]{6,64}" title="Username must be between 6 and 64 characters and contain only lowercase and uppercase characters, digit, dot, dash and underscore."
              value={this.state.username} onChange={this.handleUserInput} required/><br />
            <input className="input-form placeholder-icon-name" placeholder="First name" type="text" name="firstname"
              pattern="[a-zA-Z\-]{6,64}" title="Firstname must be between 6 and 64 characters and contain only lowercase and uppercase characters and dash."
              value={this.state.firstname} onChange={this.handleUserInput} required/><br />
            <input className="input-form placeholder-icon-name" placeholder="Last name" type="text" name="lastname"
              pattern="[a-zA-Z\-]{6,64}" title="Lastname must be between 6 and 64 characters and contain only lowercase and uppercase characters and dash."
              value={this.state.lastname} onChange={this.handleUserInput} required/><br />
            <input className="input-form" id="placeholder-icon-email" placeholder="Email address" minLength="6" maxLength="254" type="email" name="email"
              value={this.state.email} onChange={this.handleUserInput} required/><br />
            <input className="input-form" id="placeholder-icon-password" placeholder="Password" type="text" name="password"
              pattern="^(?=.*\d)(?=.*[a-z])(?=.*[A-Z])[0-9a-zA-Z]{8,100}$" title="Must contain only and at least one number and one uppercase and lowercase letter, and at least 8 or more characters"
              value={this.state.password} onChange={this.handleUserInput} required/><br />
            <input className="input-form" id="placeholder-icon-re-password" placeholder="Re-enter Password" type="text" name="rePassword"
              pattern="^(?=.*\d)(?=.*[a-z])(?=.*[A-Z])[0-9a-zA-Z]{8,100}$" title="Must contain only and at least one number and one uppercase and lowercase letter, and at least 8 or more characters"
              value={this.state.rePassword} onChange={this.handleUserInput} required/><br />
            <input className="submit-form" type="submit" value="Register"/>
          </form>
          <hr className="form-line"/>
          <div className="form-link">
            <a href='/login' className="form-link-btn"><span>Already a member ?</span></a>
          </div>
        </div>
        <Modal type="error" content={this.state.newError} onClose={this.closeError}/>
      </div>
    )
  }
}

export default Register;
