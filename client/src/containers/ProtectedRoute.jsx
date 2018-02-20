import React from 'react';
import {
  Route,
  Redirect,
} from 'react-router-dom';
import PageLayout from '../layouts/PageLayout';

// Check if the user is logged
const ProtectedRoute = ({
    isLoggedIn,
    component: Component,
    ...rest,
  }) => {
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

  export default ProtectedRoute;
