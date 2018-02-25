// import React, { Component } from 'react';

const register = (params, conf) => {
  return fetch (
    `${conf.BACK_URL}/v1/accounts/register`,
    {
      credentials: 'include',
      method: `POST`,
      body: JSON.stringify({
        username: params.username,
        email: params.email,
        lastname: params.lastname,
        firstname: params.firstname,
        password: params.password,
        rePassword: params.rePassword
      }),
    },
  );
};

export default register;
