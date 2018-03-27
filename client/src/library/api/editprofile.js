import fetch from 'isomorphic-fetch';

const editprofile = (params, conf) => {
  let token = localStorage.getItem('matcha_token');
  return fetch (
    `${conf.BACK_URL}/v1/profiles/edit/data`,
    {
      credentials: 'include',
      method: `POST`,
      body: JSON.stringify({
        lastname:       params.lastname,
        firstname:      params.firstname,
        email:          params.email,
        biography:      params.biography,
        birthday:       params.birthday,
        genre:          params.genre,
        interesting_in: params.interesting_in,
      }),
      headers: {
        'Authorization': 'Bearer ' + token,
      },
    },
  );
};

export default editprofile;
