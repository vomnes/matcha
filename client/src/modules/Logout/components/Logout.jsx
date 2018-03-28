import React, { Component } from 'react';
import api from '../../../library/api'

class Logout extends Component {
  constructor(props) {
    super(props);
    this.state = {
      hasError: false
    };
    const createError = () => {
      this.setState({ hasError: true });
    }
    api.logout()
    .then(function(response) {
      if (response.status >= 400) {
        localStorage.removeItem('matcha_token');
        createError();
        throw new Error("Bad response from server - Logout");
      } else {
        document.location = "/login";
        return;
      }
    })
  }

  render() {
    console.log(this.state.hasError);
    if (!this.state.hasError) {
      return (
        <span>
          Logging out...
        </span>
      );
    } else {
      return (
        <div>
          <span>An error has occured during logout - Please contact us - </span>
          <a href='/home'><u>Click here to go back on home page</u></a>
        </div>
      );
    }
  }
}

export default Logout;
