import fetch from 'isomorphic-fetch';

const getMessages = (username, conf) => {
  let token = localStorage.getItem('matcha_token');
  return fetch (
    `${conf.BACK_URL}/v1/chat/messages/${username}`,
    {
      credentials: 'include',
      method: `GET`,
      headers: {
        'Authorization': 'Bearer ' + token,
      },
    },
  );
};

export default getMessages;
