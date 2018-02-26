const formatName = (string, pattern) => {
  string = string.toLowerCase();
  return string.charAt(0).toUpperCase() + string.slice(1);
}

const blockForbiddenKeys = (data, pattern, maxLength) => {
  if (data.length > maxLength) {
    return -1;
  }
  let lastChar = data[data.length - (data.length ? 1 : 0)];
  if (!lastChar) {
    return "";
  } else {
    if (!lastChar.match(pattern)) {
      return -1;
    }
  }
  return data
}

export default { formatName, blockForbiddenKeys };
