const getURLParameter = (name) => {
  var url = new URL(window.location.href);
  return url.searchParams.get(name);
}
