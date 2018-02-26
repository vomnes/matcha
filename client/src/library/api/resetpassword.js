import fetch from 'isomorphic-fetch';

const resetpassword = (params, conf) => {
  return fetch (
    `${conf.BACK_URL}/v1/accounts/resetpassword`,
    {
      credentials: 'include',
      method: `POST`,
      body: JSON.stringify({
        randomToken: params.randomToken,
        password: params.password,
        rePassword: params.rePassword
      }),
    },
  );
};

export default resetpassword;
