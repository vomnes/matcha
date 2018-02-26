import fetch from 'isomorphic-fetch';

const logout = (conf) => {
  return fetch (
    `${conf.BACK_URL}/v1/accounts/logout`,
    {
      credentials: 'include',
      method: `POST`,
      headers: {
        'Authorization': 'Bearer ' + localStorage.getItem('matcha_token')
      }
    },
  );
};

export default logout;
