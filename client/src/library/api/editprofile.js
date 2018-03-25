import fetch from 'isomorphic-fetch';

const editprofile = (ip, conf) => {
  let token = localStorage.getItem('matcha_token');
  return fetch (
    `${conf.BACK_URL}/v1/profiles/edit`,
    {
      credentials: 'include',
      method: `GET`,
      headers: {
        'Content-Type': 'application/json',
        'Authorization': 'Bearer ' + token,
        'ip': ip,
      },
    },
  );
};

export default editprofile;
