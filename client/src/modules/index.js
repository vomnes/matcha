import React from 'react';
import {
  BrowserRouter as Router,
  Switch,
  Route,
} from 'react-router-dom'
import PrivateRoute from './PrivateRoute'
import Register from './Register'
import Login from './Login'
import ForgotPassword from './ForgotPassword'
import ResetPassword from './ResetPassword'
import Logout from './Logout'
import SeeProfile from './SeeProfile'
import MyProfile from './MyProfile'

const Home = () => (
  <div>
    <h1>Zelcome Home</h1>
  </div>
)

const Main = () => (
  <Router>
    <Switch>
      <Route exact path='/login' component={Login}/>
      <Route exact path='/register' component={Register}/>
      <Route exact path='/forgotpassword' component={ForgotPassword}/>
      <Route exact path='/resetpassword/:token' component={ResetPassword}/>
      <Route exact path='/logout' component={Logout}/>
      <PrivateRoute exact path='/home' component={Home}/>
      <PrivateRoute exact path='/profile/:userId' component={SeeProfile}/>
      <PrivateRoute path='/profile' component={MyProfile}/>
      <Route path='/' component={Register}/>
    </Switch>
  </Router>
)

export default Main;
