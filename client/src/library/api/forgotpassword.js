import fetch from 'isomorphic-fetch';

const forgotpassword = (params, conf) => {
  return fetch (
    `${conf.BACK_URL}/v1/mails/forgotpassword`,
    {
      credentials: 'include',
      method: `POST`,
      body: JSON.stringify({
        email: params.email,
      }),
    },
  );
};

export default forgotpassword;
