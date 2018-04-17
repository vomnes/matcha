const pictureURLFormated = (pictureURL) => {
  let url = pictureURL ? `http://localhost:8080${pictureURL}` : null
  if (pictureURL && pictureURL.includes('images.unsplash.com/photo-')) {
    url = pictureURL;
  }
  return url;
}

export default { pictureURLFormated };
