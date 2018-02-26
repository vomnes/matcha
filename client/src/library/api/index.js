import login from './login';
import register from './register';
import forgotpassword from './forgotpassword';
import resetpassword from './resetpassword';
import logout from './logout';

const apiConf = {
  BACK_URL: `http://localhost:8080`,
};

export default {
  login: (params) => login(params, apiConf),
  register: (params) => register(params, apiConf),
  forgotpassword: (params) => forgotpassword(params, apiConf),
  resetpassword: (params) => resetpassword(params, apiConf),
  logout: () => logout(apiConf),
}
