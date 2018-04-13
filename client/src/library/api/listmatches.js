import fetch from 'isomorphic-fetch';

const listMatches = (conf) => {
  let token = localStorage.getItem('matcha_token');
  return fetch (
    `${conf.BACK_URL}/v1/chat/matches`,
    {
      credentials: 'include',
      method: `GET`,
      headers: {
        'Authorization': 'Bearer ' + token,
      },
    },
  );
};

export default listMatches;
