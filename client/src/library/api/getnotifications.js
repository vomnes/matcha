import fetch from 'isomorphic-fetch';

const getNotifications = (conf) => {
  let token = localStorage.getItem('matcha_token');
  return fetch (
    `${conf.BACK_URL}/v1/users/data/notifications`,
    {
      credentials: 'include',
      method: `GET`,
      headers: {
        'Authorization': 'Bearer ' + token,
      },
    },
  );
};

export default getNotifications;
