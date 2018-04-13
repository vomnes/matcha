import React from 'react';
import {
  Route,
  Redirect,
} from 'react-router-dom';
import PageLayout from '../layouts/PageLayout';

// Check if the user is logged
const PrivateRoute = ({ component: Component, ...rest }) => {
    const isLoggedInToken = localStorage.getItem(`matcha_token`);
    if (isLoggedInToken) {
      const wsConn = new WebSocket('ws://localhost:8081/ws', isLoggedInToken);
      return (
        <Route
          {...rest}
          render={(props) => (
            <PageLayout wsConn={wsConn}>
              <Component wsConn={wsConn} {...props}/>
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
