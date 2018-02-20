import fetch from 'isomorphic-fetch';
// npm install --save isomorphic-fetch es6-promise

const login = (params, conf) => {
  return fetch (
    `${conf.BACK_URL}/v1/accounts/login`,
    {
      credentials: 'include',
      method: `POST`,
      body: JSON.stringify({
        email: params.email,
        hashedPassword: params.password,
      }),
    },
  );
};

export default login;
