import fetch from 'isomorphic-fetch';

const tag = (method, params, conf) => {
  let token = localStorage.getItem('matcha_token');
  return fetch (
    `${conf.BACK_URL}/v1/profiles/edit/tag`,
    {
      credentials: 'include',
      method: method,
      body: JSON.stringify({
        tag_name: params.tagName,
        tag_id: params.tagID,
      }),
      headers: {
        'Authorization': 'Bearer ' + token,
      },
    },
  );
};

export default tag;
