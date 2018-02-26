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
        createError();
        throw new Error("Bad response from server - Logout has failed");
      } else {
        localStorage.removeItem('matcha_token');
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
          <span>An error has occured, it is not possible to logout - Please contact us - </span>
          <a href='/home'><u>Click here to go back on home page</u></a>
        </div>
      );
    }
  }
}

export default Logout;
