import React, { Component } from 'react';
import './Register.css'

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
  }

  handleUserInput (e) {
    this.setState({
      [e.target.name]: e.target.value,
    });
  }
  handleSubmit(e) {
    console.log(this.state)
    e.preventDefault();
  }
  render() {
    return (
      <div className="background">
        <a href="/login"><p id="member">Already a member ?</p></a>
        <div id="register" className="card card-2">
          <h2 className="title">Register</h2>
          <h3 className="sub-title">A simple dating website</h3>
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
            <p className="error">{this.state.newError}</p>
          </form>
        </div>
      </div>
    )
  }
}

// https://learnetto.com/blog/how-to-do-simple-form-validation-in-reactjs

export default Register;
