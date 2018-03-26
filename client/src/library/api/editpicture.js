import fetch from 'isomorphic-fetch';

const editpicture = (params, conf) => {
  let token = localStorage.getItem('matcha_token');
  return fetch (
    `${conf.BACK_URL}/v1/profiles/picture/${params.number}`,
    {
      credentials: 'include',
      method: params.method,
      body: JSON.stringify({
        picture_base64: params.base64,
      }),
      headers: {
        'Authorization': 'Bearer ' + token,
      },
    },
  );
};

export default editpicture;
