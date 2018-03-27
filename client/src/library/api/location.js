import fetch from 'isomorphic-fetch';

const location = (params, conf) => {
  let token = localStorage.getItem('matcha_token');
  return fetch (
    `${conf.BACK_URL}/v1/profiles/edit/location`,
    {
      credentials: 'include',
      method: `POST`,
      body: JSON.stringify({
        lat: params.lat,
        lng: params.lng,
        city: params.city,
        zip: params.zip,
        country: params.country,
      }),
      headers: {
        'Authorization': 'Bearer ' + token,
      },
    },
  );
};

export default location;
