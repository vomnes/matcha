import React from 'react';
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
      console.log(response.error);
    } else {
      return response;
    }
  }
}

// Check if the user is logged
const PrivateRoute = ({ component: Component, ...rest }) => {
    const isLoggedInToken = localStorage.getItem(`matcha_token`);
    if (isLoggedInToken) {
      const wsConn = new WebSocket(`ws://localhost:8081/ws/${isLoggedInToken}`);
      const myProfileData = GetMe();
      return (
        <Route
          {...rest}
          render={(props) => (
            <PageLayout wsConn={wsConn} myProfileData={myProfileData}>
              <Component wsConn={wsConn} myProfileData={myProfileData} {...props}/>
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
  };

  export default PrivateRoute;
