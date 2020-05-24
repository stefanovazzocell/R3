'use strict';

/*
* ToBase64(u8) - Converts a Uint8Array to a base64 string
*
* @requires u8 {Object} a bits array to convert
* @returns  String
*/
function ToBase64(u8) {
  u8 = new Uint8Array(u8);
  const CHUNK_SIZE = 0x8000;
  let index = 0;
  const length = u8.length;
  let result = '';
  let slice;
  while (index < length) {
    slice = u8.subarray(index, Math.min(index + CHUNK_SIZE, length));
    result += String.fromCharCode.apply(null, slice);
    index += CHUNK_SIZE;
  }
  return btoa(result);
}

/*
* FromBase64(b64) - Converts a base64 string to a Uint8Array
*
* @requires b64 {string} the string to convert
* @returns  Uint8Array
*/
function FromBase64(b64) {
  return new Uint8Array(atob(b64).split('').map((c) => { return c.charCodeAt(0); }));
}

/*
* query(type, url, data) - Converts a base64 string to a Uint8Array
*
* @requires type {string} the request type (GET, POST, PUT, DELETE, ...)
* @requires url  {string} the url path
* @requires data {Object} the data to send (if any)
* @returns  Promise
*/
async function query(type, url, data) {
  // TODO: Change link
  const response = await fetch(`${API_ENDPOINT}${url}`, {
    method: type,
    headers: {
      'Content-Type': 'application/json'
    },
    body: (data ? JSON.stringify(data) : undefined)
  });
  return await response;
}

/*
* randomInt(min, max) - Gets a random int in a given range [min, max)
*
* @requires start {int} the start
* @returns  Int
*/
function randomInt(min, max) {
  return Math.floor((window.crypto.getRandomValues(new Uint32Array(1))[0] / 2**32) * (max - min) + min);
}

/*
* randomString(len) - Generates a random string (exclude I/l)
*
* @requires len {int} the length
* @returns  String
*/
function randomString(len) {
  const chars = 'abcdefghijkmnopqrstuvwxyzABCDEFGHJKLMNOPQRSTUVWXYZ0123456789';
  let out = '';
  for (let i = len; i > 0; i--) {
    out += chars.charAt(randomInt(0, 60));
  }
  return out;
}

/*
* randomKey() - Generates a random key
*
* @returns  String
*/
function randomKey() {
  return [randomString(3), randomString(3), randomString(3), randomString(3)].join('.');
}

/*
* onReady(callback) - executes callback when page ready
*
* @requires callback {Function} to be the function to callback
*/
function onReady(callback) {
  if(document.readyState === 'interactive' || document.readyState === 'complete') {
    callback();
  } else document.addEventListener('DOMContentLoaded', callback);
}

/*
* log(msg) - Logs a message
*
* @requires msg {String} to be the message to log
*/
let glogs = [];
function log(msg) {
  console.info(`${msg}`);
  glogs.push(msg);
}
