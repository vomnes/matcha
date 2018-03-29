import fetch from 'isomorphic-fetch';

const like = (method, username, conf) => {
  let token = localStorage.getItem('matcha_token');
  return fetch (
    `${conf.BACK_URL}/v1/users/`+username+`/like`,
    {
      credentials: 'include',
      method: method,
      headers: {
        'Authorization': 'Bearer ' + token,
      },
    },
  );
};

export default like;
