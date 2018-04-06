import fetch from 'isomorphic-fetch';

const existingTags = (conf) => {
  let token = localStorage.getItem('matcha_token');
  return fetch (
    `${conf.BACK_URL}/v1/users/data/tags`,
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

export default existingTags;
