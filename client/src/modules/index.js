import React from 'react';
import {
  BrowserRouter as Router,
  Switch,
  Route,
} from 'react-router-dom'
import PrivateRoute from './PrivateRoute'
import Register from './Register'

const Login = () => (
  <div>
    <h1>Login</h1>
  </div>
  // localStorage.setItem(`matcha_token`, `<token>`);
)

const Main = () => (
  <Router>
    <Switch>
      <Route exact path='/' component={Register}/>
      <Route exact path='/login' component={Login}/>
      <PrivateRoute exact path='/home' component={Register}/>
    </Switch>
  </Router>
)

export default Main;
