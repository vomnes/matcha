import React, { Component } from 'react';
import Error from '../../../components/Error'
import './Register.css'
import api from '../../../api'

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

class Register extends Component {
  constructor (props) {
    super(props);
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
    this.setState({
      [e.target.name]: e.target.value,
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
      <div id="background-register">
        <a href="/login"><p id="member">Already a member ?</p></a>
        <div id="register">
          <h1 className="up-title-form">Matcha</h1>
          <h2 className="title-form">Register</h2>
          <h3 className="sub-title-form">A simple dating website</h3>
          <form onSubmit={this.handleSubmit}>
            <input className="input-form" id="placeholder-icon-username" placeholder="Username" type="text" name="username"
              value={this.state.username} onChange={this.handleUserInput}/><br />
            <input className="input-form placeholder-icon-name" placeholder="First name" type="text" name="firstname"
              value={this.state.firstname} onChange={this.handleUserInput}/><br />
            <input className="input-form placeholder-icon-name" placeholder="Last name" type="text" name="lastname"
              value={this.state.lastname} onChange={this.handleUserInput}/><br />
            <input className="input-form" id="placeholder-icon-email" placeholder="Email address" type="email" name="email"
              value={this.state.email} onChange={this.handleUserInput}/><br />
            <input className="input-form" id="placeholder-icon-password" placeholder="Password" type="password" name="password"
              value={this.state.password} onChange={this.handleUserInput}/><br />
            <input className="input-form" id="placeholder-icon-re-password" placeholder="Re-enter Password" type="password" name="rePassword"
              value={this.state.rePassword} onChange={this.handleUserInput}/><br />
            <input className="submit-form" type="submit" value="Register"/>
          </form>
          <hr className="form-line"/>
          <div className="form-link">
            <a href='/login' className="form-link-click"><span>Already a member ?</span></a>
          </div>
        </div>
        <Error content={this.state.newError} onClose={this.closeError}/>
      </div>
    )
  }
}

export default Register;
