import login from './login';
import register from './register';
import forgotpassword from './forgotpassword';
import resetpassword from './resetpassword';
import logout from './logout';
import getprofile from './getprofile';
import editpicture from './editpicture';
import location from './location';
import editprofile from './editprofile';
import editpassword from './editpassword';
import tag from './tag';
import getuser from './getuser';
import like from './like';
import fake from './fake';
import me from './me';
import match from './match';
import existingTags from './existingtags';

const apiConf = {
  BACK_URL: `http://localhost:8080`,
};

export default {
  login: (params) => login(params, apiConf),
  register: (params) => register(params, apiConf),
  forgotpassword: (params) => forgotpassword(params, apiConf),
  resetpassword: (params) => resetpassword(params, apiConf),
  logout: () => logout(apiConf),
  getprofile: (params) => getprofile(params, apiConf),
  editpicture: (params) => editpicture(params, apiConf),
  location: (params) => location(params, apiConf),
  editprofile: (params) => editprofile(params, apiConf),
  editpassword: (params) => editpassword(params, apiConf),
  tag: (method, params) => tag(method, params, apiConf),
  getuser: (username) => getuser(username, apiConf),
  like: (method, username) => like(method, username, apiConf),
  fake: (method, username) => fake(method, username, apiConf),
  me: () => me(apiConf),
  match: (options) => match(options, apiConf),
  existingTags: () => existingTags(apiConf),
}
