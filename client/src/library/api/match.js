import fetch from 'isomorphic-fetch';

const match = (options, conf) => {
  let token = localStorage.getItem('matcha_token');
  return fetch (
    `${conf.BACK_URL}/v1/users/data/match`,
    {
      credentials: 'include',
      method: `GET`,
      headers: {
        'Authorization': 'Bearer ' + token,
      },
    },
  );
};

export default match;
