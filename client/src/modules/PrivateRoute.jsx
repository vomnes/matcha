import React from 'react';
import {
  Route,
  Redirect,
} from 'react-router-dom';
import PageLayout from '../layouts/PageLayout';

// Check if the user is logged
const PrivateRoute = ({
    component: Component,
    ...rest,
  }) => {
    const isLoggedIn = localStorage.getItem(`matcha_token`);
    return (
      <Route
        {...rest}
        render={(props) => (
          isLoggedIn ? (
            <PageLayout>
              <Component {...props} />
            </PageLayout>
          ) : (
            <Redirect to={{
              pathname: '/login',
              state: { from: props.location }
            }}/>
          )
        )}
      />
    );
  };

  export default PrivateRoute;
