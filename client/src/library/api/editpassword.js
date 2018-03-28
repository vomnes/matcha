import fetch from 'isomorphic-fetch';

const editpassword = (params, conf) => {
  let token = localStorage.getItem('matcha_token');
  return fetch (
    `${conf.BACK_URL}/v1/profiles/edit/password`,
    {
      credentials: 'include',
      method: `POST`,
      body: JSON.stringify({
        password:         params.password,
        new_password:     params.new_password,
        new_rePassword:   params.new_rePassword,
      }),
      headers: {
        'Authorization': 'Bearer ' + token,
      },
    },
  );
};

export default editpassword;
