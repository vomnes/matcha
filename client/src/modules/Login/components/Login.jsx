import React, { Component } from 'react';
import Modal from '../../../components/Modal'
import './Login.css'
import api from '../../../library/api'
import utils from '../../../library/utils/input.js'

const getSuccessContent = () => {
  const query = new URLSearchParams(window.location.search);
  const value = query.get('code');
  if (value === '1') {
    return "Account created";
  }
  return "";
}

const signIn = (username, password, createError, redirectHome) => {
  api.login({
      username,
      password,
      'uuid': window.navigator.userAgent.replace(/\D+/g, ''),
  }).then(function(response) {
    if (response.status >= 500) {
      throw new Error("Bad response from server - SignIn has failed");
    } else if (response.status >= 400) {
      response.json().then(function(data) {
        createError(data.error);
        return;
      });
    } else {
      response.json().then(function(data) {
        localStorage.setItem('matcha_token', data.token);
        redirectHome();
      });
      return;
    }
  })
}

class Login extends Component {
  constructor (props) {
    super(props);
    if (localStorage.getItem(`matcha_token`)) {
      this.props.history.push("/home");
    }
    this.state = {
      username: '',
      password: '',
      newError: '',
      newSuccess: getSuccessContent(),
    }
    this.handleUserInput = this.handleUserInput.bind(this);
    this.handleSubmit = this.handleSubmit.bind(this);
    this.createError = this.createError.bind(this);
    this.closeModal = this.closeModal.bind(this);
    this.redirectHome = this.redirectHome.bind(this);
  }
  handleUserInput (e) {
    const fieldName = e.target.name
    var data = e.target.value
    data = utils.formatInput(fieldName, data)
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
  closeModal(event) {
    this.setState({
      newError: '',
      newSuccess: ''
    });
    event.preventDefault();
  }
  redirectHome() {
    this.props.history.push("/profile");
  }
  handleSubmit(e) {
    this.setState({
      newError: '',
    });
    signIn(
      this.state.username,
      this.state.password,
      this.createError,
      this.redirectHome,
    )
    e.preventDefault();
  }
  render() {
    return (
      <div id="background-login">
        <div id="box-center-top-form">
          <h1 className="up-title-form">Matcha</h1>
          <h2 className="title-form">Login</h2>
          <h3 className="sub-title-form">Welcome back !</h3>
          <form onSubmit={this.handleSubmit}>
            <input className="input-form" id="placeholder-icon-username" placeholder="Username" type="text" name="username" autoComplete="username"
              value={this.state.username} onChange={this.handleUserInput} required/><br />
            <input className="input-form" id="placeholder-icon-password" placeholder="Password" type="password" name="password" autoComplete="current-password"
              value={this.state.password} onChange={this.handleUserInput} required/><br />
            <input className="submit-form" type="submit" value="Enter"/>
            <div className="form-link">
              <a href='/forgotpassword' className="form-link-click"><span>Forgot password ?</span></a> - <a href='/register' className="form-link-click"><span>Not registered yet ?</span></a>
            </div>
          </form>
        </div>
        <Modal type="error" content={this.state.newError} onClose={this.closeModal}/>
        <Modal type="success" content={this.state.newSuccess} onClose={this.closeModal}/>
      </div>
    )
  }
}

export default Login
