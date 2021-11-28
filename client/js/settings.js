'use strict';

/*
* Encryption
*/

// The salt to be used when deriving the hash
const HASH_SALT = new Uint8Array([82,242,11,190,119,15,58,152,115,230,184,149,107,12,5,37,184,242,159,111,72,180,65,53,104,78,252,123,188,17,71,187,216,128,141,148,126,110,15,113,175,70,216,37,211,247,93,216,210,197,189,100,37,81,113,113,173,8,184,97,225,223,24,69]);
// The iteration for the key generation the hash
const HASH_ROUNDS = 1000000;

// The iteration for the key generation the encryption key
const KEY_ROUNDS = 1000000;

/*
* API
*/

// The API endpoint
const API_ENDPOINT = `https://api.rkt.one/v1/`;
