import login from './login';

const apiConf = {
  BACK_URL: `localhost:8080`,
};

export default {
  login: (params) => login(params, apiConf),
}
