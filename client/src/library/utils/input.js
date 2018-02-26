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

const formatInput = (fieldName, data) => {
  if (fieldName === "username") {
    return blockForbiddenKeys(data, /[0-9a-zA-Z.\-_]/i, 64);
  } else if (fieldName === "firstname" || fieldName === "lastname") {
    data = formatName(data);
    return blockForbiddenKeys(data, /[a-zA-Z-]/i, 64);
  } else if (fieldName === "email") {
    data = data.toLowerCase();
    return blockForbiddenKeys(data, /[a-zA-Z@.]/i, 254);
  }
  return data
}

export default { formatName, blockForbiddenKeys, formatInput };
