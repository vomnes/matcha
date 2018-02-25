import React from 'react';
import {
  BrowserRouter as Router,
  Switch,
  Route,
} from 'react-router-dom'
import PrivateRoute from './PrivateRoute'
import Register from './Register'
import Login from './Login'

const Home = () => (
  <div>
    <h1>Zelcome Home</h1>
  </div>
)

const Main = () => (
  <Router>
    <Switch>
      <Route exact path='/' component={Register}/>
      <Route exact path='/login' component={Login}/>
      <Route exact path='/register' component={Register}/>
      <PrivateRoute exact path='/home' component={Home}/>
    </Switch>
  </Router>
)

export default Main;
