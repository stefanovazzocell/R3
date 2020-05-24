'use strict';

/*
* Encryption
*/

// The salt to be used when deriving the hash
const HASH_SALT = new Uint8Array([231,68,163,92,124,82,245,221,98,93,218,109,5,22,130,112,242,31,182,6,237,254,21,135,125,104,73,150,132,40,197,174,68,111,226,211,86,146,124,245,248,153,50,123,78,132,180,89,119,16,228,177,42,19,48,221,192,245,163,23,32,195,139,19]);
// The iteration for the key generation the hash
const HASH_ROUNDS = 500000;
// The iteration for the key generation pre hash
const PRE_HASH_ROUNDS = 1000;

// The iteration for the key generation the encryption key
const KEY_ROUNDS = 800000;

/*
* API
*/

// The API endpoint
const API_ENDPOINT = `${window.location.origin}/v1/`;
