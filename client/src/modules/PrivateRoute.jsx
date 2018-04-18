import React, { Component } from 'react';
import {
  Route,
  Redirect,
} from 'react-router-dom';
import PageLayout from '../layouts/PageLayout';
import api from '../library/api'

const GetMe = async () => {
  let res = await api.me();
  if (res) {
    const response = await res.json();
    if (res.status >= 500) {
      throw new Error("Bad response from server - GetMe has failed");
    } else if (res.status >= 400) {
      return res.status;
    } else {
      return response;
    }
  }
}

// Check if the user is logged
class PrivateRoute extends Component {
  constructor() {
    super();
    this.state = {
      token: localStorage.getItem(`matcha_token`),
    }
    this.updateState = this.updateState.bind(this);
  }
  updateState(key, content) {
    this.setState({
      [key]: content,
    });
  }
  handleRedirect() {
    GetMe().then((data) => {
      if (data >= 400) {
        localStorage.removeItem('matcha_token');
        window.location.replace(`/login`);
      } else {
        if (data.redirect && data.redirect.length) {
          if (!window.location.href.includes("/profile")) {
            window.location.replace(`/profile?empty=${data.redirect.join('|')}`);
          }
        } else {
          this.updateState('myProfileData', data);
        }
      }
    });
  }
  componentWillMount() {
    this.handleRedirect();
  }
  render() {
    const {component: Component, ...rest} = this.props;

    if (this.state.token) {
      const wsConn = new WebSocket(`ws://localhost:8081/ws/${this.state.token}`);
      return (
        <Route
          {...rest}
          render={(props) => (
            <PageLayout wsConn={wsConn} myProfileData={GetMe()}>
              <Component myProfileData={GetMe()} {...props} wsConn={wsConn}/>
            </PageLayout>
          )}
        />
      );
    }
    return (
      <Route
        {...rest}
        render={(props) => (
          <Redirect to={{
            pathname: '/login',
            state: { from: props.location }
          }}/>
        )}
      />
    );
  }
}

export default PrivateRoute;
