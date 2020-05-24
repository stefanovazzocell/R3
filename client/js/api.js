'use strict';

/*
* apiGet(link) - Retrives a link
*
* @requires link {String} to be the PBKDF2 key
* @returns  Promise
*/
async function apiGet(link) {
  return await query('POST', 'link', {
    id: link
  });
}

/*
* apiCreate(link, data, hits, ttl, edit) - Adds a new link
*
* @requires link {String} to be the PBKDF2 key
* @requires data {String} to be the link data
* @requires hits {Number} to be the max amount of views to allow
* @requires ttl  {Number} to be the ttl in seconds for the link
* @requires edit {String} to be the PBKDF2 key for editing (can be blank)
* @returns  Promise
*/
async function apiCreate(link, data, hits, ttl, edit) {
  return await query('POST', 'links', {
    id: link,
    payload: {
      data: data,
      hits: hits,
      ttl:  ttl,
      edit: edit
    }
  });
}

/*
* apiEdit(link, data, hits, ttl, edit, pass) - Edits a link
*
* @requires link {String} to be the PBKDF2 key
* @requires data {String} to be the link data
* @requires hits {Number} to be the max amount of views to allow
* @requires ttl  {Number} to be the ttl in seconds for the link
* @requires edit {String} to be the PBKDF2 key for editing (can be blank)
* @requires pass {String} to be the link editing password
* @returns  Promise
*/
async function apiEdit(link, data, hits, ttl, edit, pass) {
  return await query('PUT', 'link', {
    id: link,
    password: pass,
    payload: {
      data: data,
      hits: hits,
      ttl:  ttl,
      edit: edit
    }
  });
}

/*
* apiDelete(link, pass) - Deletes a link
*
* @requires link {String} to be the PBKDF2 key
* @requires pass {String} to be the link editing password
* @returns  Promise
*/
async function apiDelete(link, pass) {
  return await query('Delete', 'link', {
    id: link,
    password: pass
  });
}
