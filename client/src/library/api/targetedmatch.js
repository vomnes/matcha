import fetch from 'isomorphic-fetch';

const targetedmatch = (options, username, conf) => {
  let token = localStorage.getItem('matcha_token');
  var headers = {
    'Authorization': 'Bearer ' + token,
  }
  if (options) {
    headers['Search-Parameters'] = options;
  }
  return fetch (
    `${conf.BACK_URL}/v1/users/data/match/${username}`,
    {
      credentials: 'include',
      method: `GET`,
      headers,
    },
  );
};

export default targetedmatch;
