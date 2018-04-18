const pictureURLFormated = (pictureURL, quality) => {
  let url = pictureURL ? `http://localhost:8080${pictureURL}` : null
  if (pictureURL && pictureURL.includes('images.unsplash.com/photo-')) {
    url = pictureURL;
    if (quality) {
      url = url.replace("h=1000&q=10", quality);
    }
  }
  return url;
}

export default { pictureURLFormated };
