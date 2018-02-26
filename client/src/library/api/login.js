import fetch from 'isomorphic-fetch';

const login = (params, conf) => {
  return fetch (
    `${conf.BACK_URL}/v1/accounts/login`,
    {
      credentials: 'include',
      method: `POST`,
      body: JSON.stringify({
        username: params.username,
        password: params.password,
        uuid: params.uuid,
      }),
    },
  );
};

export default login;
