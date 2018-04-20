import React, { Component } from 'react';
import api from '../../../library/api';
import { Redirect } from 'react-router-dom';

class Logout extends Component {
  constructor(props) {
    super(props);
    this.state = {
      hasError: false
    };
    const createError = () => {
      this.setState({ hasError: true });
    }
    try {
      api.logout()
      .then(function(response) {
        if (response.status >= 400) {
          localStorage.removeItem('matcha_token');
          createError();
          if (response.status >= 500) {
            throw new Error("Bad response from server - Logout");
          }
        } else {
          return;
        }
      })
    } catch (e) {
      console.log(e.message);
    }
  }

  render() {
    localStorage.removeItem('matcha_token');
    if (!this.state.hasError) {
      return (
        <Redirect to='/login'/>
      );
    } else {
      return (
        <Redirect to='/login'/>
      );
    }
  }
}

export default Logout;
