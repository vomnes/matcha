import React from 'react';
import {
  BrowserRouter as Router,
  Switch,
  Route,
} from 'react-router-dom'
import PrivateRoute from './PrivateRoute'

const Home = () => (
  <div>
    <h1>Welcome to the Matcha Website!</h1>
  </div>
)

const Login = () => (
  <div>
    <h1>Login</h1>
  </div>
  // localStorage.setItem(`matcha_token`, `<token>`);
)

const Main = () => (
  <Router>
    <Switch>
      <Route exact path='/' component={Home}/>
      <Route exact path='/login' component={Login}/>
      <PrivateRoute exact path='/home' component={Home}/>
    </Switch>
  </Router>
)

export default Main;
