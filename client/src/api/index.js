import login from './login';
import register from './register';

const apiConf = {
  BACK_URL: `http://localhost:8080`,
};

export default {
  login: (params) => login(params, apiConf),
  register: (params) => register(params, apiConf),
}
