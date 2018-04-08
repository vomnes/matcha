const objectToBase64 = (object) => {
    return Buffer.from(JSON.stringify(object)).toString("base64")
}

export default { objectToBase64 };
