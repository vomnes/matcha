import fetch from 'isomorphic-fetch';

const getuser = (username, conf) => {
  let token = localStorage.getItem('matcha_token');
  return fetch (
    `${conf.BACK_URL}/v1/users/` + username,
    {
      credentials: 'include',
      method: `GET`,
      headers: {
        'Content-Type': 'application/json',
        'Authorization': 'Bearer ' + token,
      },
    },
  );
};

export default getuser;
