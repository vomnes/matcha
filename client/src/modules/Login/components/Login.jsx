import React, { Component } from 'react';
import Error from '../../../components/Error'
import './Login.css'
// import api from '../../../api'

class Login extends Component {
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
    this.props.history.push("/home");
  }
  handleSubmit(e) {
    this.setState({
      newError: '',
    });
    // signUp(
    //   this.state.username,
    //   this.state.password,
    // )
    e.preventDefault();
  }
  render() {
    return (
      <div id="background-login">
        <div id="login">
          <h1 className="up-title-form">Matcha</h1>
          <h2 className="title-form">Login</h2>
          <h3 className="sub-title-form">Welcome back !</h3>
          <form onSubmit={this.handleSubmit}>
            <input className="input-form" id="placeholder-icon-username" placeholder="Username" type="text" name="username"
              value={this.state.username} onChange={this.handleUserInput}/><br />
            <input className="input-form" id="placeholder-icon-password" placeholder="Password" type="password" name="password"
              value={this.state.password} onChange={this.handleUserInput}/><br />
            <input className="submit-form" type="submit" value="Enter"/>
            <div className="form-link">
              <a href='/forgotpassword' className="form-link-click"><span>Forgot password ?</span></a> - <a href='/register' className="form-link-click"><span>Not registered yet ?</span></a>
            </div>
          </form>
        </div>
        <Error content={this.state.newError} onClose={this.closeError}/>
      </div>
    )
  }
}

export default Login
