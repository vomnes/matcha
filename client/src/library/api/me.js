import fetch from 'isomorphic-fetch';

const me = (conf) => {
  let token = localStorage.getItem('matcha_token');
  return fetch (
    `${conf.BACK_URL}/v1/users/data/me`,
    {
      credentials: 'include',
      method: `GET`,
      headers: {
        'Authorization': 'Bearer ' + token,
      },
    },
  );
};

export default me;
