import fetch from 'isomorphic-fetch';

const messages = (method, username, conf) => {
  let token = localStorage.getItem('matcha_token');
  return fetch (
    `${conf.BACK_URL}/v1/chat/messages/${username}`,
    {
      credentials: 'include',
      method: method,
      headers: {
        'Authorization': 'Bearer ' + token,
      },
    },
  );
};

export default messages;
