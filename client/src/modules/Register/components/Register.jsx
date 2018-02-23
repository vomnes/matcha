import React, { Component } from 'react';
import './Register.css'

class Register extends Component {
  constructor (props) {
    super(props);
    this.state = {
      name: '',
    }
  }

  handleUserInput (e) {
    const name = e.target.name;
    const value = e.target.value;
    this.setState({[name]: value});
  }
  handleSubmit(event) {
    alert(event.username.value);
  }
  render() {
    return (
      <div className="background">
        <div id="register" class="card card-2">
          <h2 className="title">Register</h2>
          <h3 className="sub-title">A simple dating website</h3>
          <form onSubmit={this.handleSubmit}>
            <input class="input-form" placeholder="Username" type="text" name="username" value={this.state.name} onChange={(event) => this.handleUserInput(event)}/><br />
            <input class="input-form" placeholder="Email address" type="email" name="email" required/><br />
            <input class="input-form" placeholder="Password" type="password" name="password"/><br />
            <input class="input-form" placeholder="Re-enter Password" type="password" name="re-password"/><br />
            <input type="submit" value="Registered"/>
            <p>Content Error</p>
          </form>
        </div>
      </div>
    )
  }
}

// https://learnetto.com/blog/how-to-do-simple-form-validation-in-reactjs

export default Register;
