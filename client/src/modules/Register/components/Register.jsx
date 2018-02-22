import React, { Component } from 'react';
import './Register.scss'

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
      <div id="register">
        <h2>Register</h2>
        <form onSubmit={this.handleSubmit}>
          <fieldset>
           <label for="name">Username</label><input type="text" name="username" value={this.state.name} onChange={(event) => this.handleUserInput(event)}/>
           {/* <label for="name">Password</label><input type="password" name="password"/> */}
           {/* <label for="name">Re-enter password</label><input type="password" name="re-password"/> */}
           {/* <label for="name">Email address</label><input type="email" name="email" required/> */}
         </fieldset>
         <input type="submit" value="Registered"/>
         <p>Content Error</p>
        </form>
      </div>
    )
  }
}

// https://learnetto.com/blog/how-to-do-simple-form-validation-in-reactjs

export default Register;
