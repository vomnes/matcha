import fetch from 'isomorphic-fetch';

const match = (options, conf) => {
  let token = localStorage.getItem('matcha_token');
  var headers = {
    'Authorization': 'Bearer ' + token,
  }
  if (options) {
    headers['Search-Parameters'] = options;
  }
  console.log(headers['Search-Parameters']);
  return fetch (
    `${conf.BACK_URL}/v1/users/data/match`,
    {
      credentials: 'include',
      method: `GET`,
      headers,
    },
  );
};

export default match;
