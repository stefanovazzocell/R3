'use strict';

/*
* selectAll(query, fn) - Query multiple elements and loop through them
*
* @requires query {String} the query string
* @requires fn {Function} the function to call
*/
function selectAll(query, fn) {
  Array.from(document.querySelectorAll(query)).forEach(fn);
}

/*
* show(query) - Shows all the elements that match the query
*
* @requires query {String} the query string
*/
function show(query) {
  selectAll(query, (el)=>{el.classList.remove('hide')});
}

/*
* hide(query) - Hides all the elements that match the query
*
* @requires query {String} the query string
*/
function hide(query) {
  selectAll(query, (el)=>{el.classList.add('hide')});
}

/*
* onClick(btn, fn) - Handles clicks for a button
*
* @requires btn {String} the selector for the button
* @requires fn {Function} the function to call
*/
function onClick(btn, fn) {
  selectAll(btn, (el)=>{ el.addEventListener('click', (e)=>{ fn(e, el); }); });;
}

/*
* checkInput(input) - Checks input for validity
*
* @requires input {String} the selector for the input
*
* @returns {bool} true if valid, false otherwise
*/
function checkInput(input) {
  return document.getElementById(input).checkValidity();
}

/*
* getInput(input) - Gets an input value
*
* @requires input {String} the selector for the input
*
* @returns {string|int} the value of the input
*/
function getInput(input) {
  return document.getElementById(input).value;
}

/*
* getBtnInput(input) - Gets an input value from a .tab element
*
* @requires input {String} the selector for the input
*
* @returns {string|int} the value of the input
*/
function getBtnInput(input) {
  return document.getElementById(input).querySelector('div.active').dataset.value;
}

/*
* setInput(input, val) - Sets an input value
*
* @requires input {String} the selector for the input
* @requires val {String|int} the input value
*/
function setInput(input, val) {
  document.getElementById(input).value = val;
}

/*
* fileReader(query) - Reads a file given an element id for a picker
*
* @requires query {String} the selector for the input
*
* @returns {Promise}
*/
async function fileReader(input) {
  return new Promise((resolve)=>{
    const reader = new FileReader();
    reader.onload = (e) => {
      resolve(e.target.result);
    }
    reader.readAsDataURL(document.getElementById(input).files[0]);
  });
}

/*
* percentage(val, msg) - Sets the loading percentage to a given value
*
* @requires val {Int} percentage completed
* @requires msg {String} message to show
*/
function percentage(val, msg) {
  selectAll('.loader > .bar', (el)=>{ el.style.width = `${val}%`; });
  selectAll('.loader > .message', (el)=>{ el.innerText = msg; });
}

/*
* error(msg) - Shows an error message
*
* @requires msg {String} message to show
* @requires nolog {Bool|undefined} true to disable logging
*/
function error(msg, nolog) {
  let alert = document.getElementById('sample').cloneNode(true);
  alert.id = '';
  alert.children[0].innerText = msg;
  alert.onclick = ()=>{ alert.remove(); };
  setTimeout(()=>{ alert.remove(); }, 5000 * Math.ceil(msg.length / 50));
  document.getElementById('alerts').appendChild(alert);
  if (nolog !== true) log(`[error] ${msg}`);
}

/*
* generate_random_id() - Generates a random ID
*/
function generate_random_id() {
  const id_input = document.getElementById('id_input');
  const link = document.getElementById('id_link');
  id_input.value = randomKey();
  link.value = `https://${window.location.hostname}/#${id_input.value}`
}

/*
* ui_reset() - Resets the UI
*/
function ui_reset() {
  window.location.hash = '';
  generate_random_id();
  selectAll('.tab.default', (el)=>{el.click();});
  selectAll('#container_shared > div > input, #container_shared > div > textarea', (el)=>{ el.value = ''; });
  hide('#container_settings > .loader, #container_loader, #container_final, #container_view, #container_error, #image_preview');
  selectAll('#app > .default, #container_settings > .tab_group', (el)=>{ el.classList.remove('disabled'); });
  show('#app > .default, #btn_create');
}

/*
* ui_big_query_watch() - Watches for Bigqueries
*/
function ui_big_query_watch() {
  const isActive = (eid)=>{return !document.getElementById(eid).classList.contains('hide')};
  if (isActive('tab_image') || isActive('tab_file') || document.getElementById('message_data_input').value.length > 400000) {
    selectAll('.tab.high', (el)=>{ el.classList.add('disabled'); el.classList.remove('active'); });
    selectAll('#settings_time > .tab.default', (el)=>{ el.classList.add('active'); });
  } else {
    selectAll('.tab.high', (el)=>{ el.classList.remove('disabled'); });
  }
}

/*
* CreateLinkUI(data) - Encrypts data and creates a link for it
*
* @requires data {Object} the data to encrypt
* @returns  Promise
*/
async function CreateLinkUI(data) {
  percentage(0, 'Generating keys...');
  hide('#btn_create');
  show('#loader');
  selectAll('#container_id, #container_shared, #container_settings > .tab_group', (el)=>{ el.classList.add('disabled'); });
  data = JSON.stringify(data);
  const k = getInput('id_input');
  const e0 = performance.now();
  const hashP = SimpleHash(k);
  const encryptedP = SimpleEncrypt(data, await SimpleKey(k, false));
  percentage(40, 'Encrypting data...');
  const output = await Promise.all([
    hashP,
    encryptedP
  ]);
  const e1 = performance.now();
  percentage(60, 'Creating Link...');
  log(`[perf] Encryption ${e1-e0}ms`);
  const r0 = performance.now();
  const response = await apiCreate(output[0], output[1], parseInt(getBtnInput('settings_views')), parseInt(getBtnInput('settings_time')) * 60, '');
  const r1 = performance.now();
  log(`[perf] Network ${r1-r0}ms`);
  percentage(100, 'Done.');
  if (response.status !== 200) {
    percentage(100, 'Error.');
    response.json().then((j)=>{
      error(j.error ? j.error : "Couldn't reach the server; try again later.");
    }).catch((e)=>{
      error("Couldn't reach the server; try again later.");
    }).finally(()=>{
      ui_reset();
    });
  } else {
    hide('#app > div');
    show('#container_final');
  }
}

/*
* GetLinkUI(share_id) - Retrives the encrypted data and decrypts it
*
* @requires share_id {String} the share id
* @returns  Promise
*/
async function GetLinkUI(share_id) {
  const errorFn = (e) => { document.getElementById('view_error_message').innerText = e; }
  percentage(5, 'Generating keys...');
  const h0 = performance.now();
  const hash = await SimpleHash(share_id);
  const h1 = performance.now();
  log(`[perf] Hash Derivation ${h1-h0}ms`);
  percentage(40, 'Retriving data...');
  const n0 = performance.now();
  const response = await apiGet(hash);
  const n1 = performance.now();
  log(`[perf] Network ${n1-n0}ms`);
  if (response.status !== 200) {
    percentage(100, 'Error.');
    try {
      const json = await response.json();
      errorFn(json.error ? json.error : "Couldn't reach the server; try again later.");
      return false;
    } catch (e) {
      console.error(e);
      errorFn("Couldn't reach the server; try again later.");
      return false;
    }
  } else {
    try {
      const json_res = JSON.parse((await response.json()).data);
      percentage(80, 'Decrypting Data...');
      const e0 = performance.now();
      const data = await SimpleDecrypt(json_res, await SimpleKey(share_id, json_res.salt));
      const e1 = performance.now();
      log(`[perf] Encryption ${e1-e0}ms`);
      percentage(100, 'Done.');
      return JSON.parse(data);
    } catch (e) {
      console.error(e);
      percentage(100, 'Error.');
      errorFn("Error reading the share.");
      return false;
    }
  }
}

/*
* compressImage(data) - Compresses an image
*
* @returns {Promise}
*/
async function compressImage(data) {
  const start_size = data.length;
  return new Promise((resolve, reject)=>{
    const img = new Image();
    let count = 0;
    img.onload=()=>{
      const t0 = performance.now();
      const canvas = document.createElement('canvas');
      const ctx = canvas.getContext('2d');
      let width = img.width; canvas.width = width;
      let height = img.height; canvas.height = height;
      let cur_data = '';
      let m = 1;
      const tuner = (min, max, q) => {
        count++;
        if (count > 20) {
          reject('Unable to compress the image.');
          return;
        }
        m = min + (max - min) / 2;
        ctx.clearRect(0, 0, width, height);
        canvas.width = width * m;
        canvas.height = height * m;
        ctx.drawImage(img, 0, 0, width * m, height * m);
        cur_data = canvas.toDataURL('image/jpeg', q);
        if (1000000 <= cur_data.length && cur_data.length <= 1200000) {
          const t1 = performance.now();
          resolve(cur_data);
          log(`[perf] Compression ${t1-t0}ms, ${Math.round(q*100)}% quality, ${Math.round(m*100)}% size, ${Math.round(start_size / (1024 * 1024))}->${Math.round(cur_data.length / (1024 * 1024))}Mb (${Math.round(100*cur_data.length/start_size)}%).`);
          return;
        }
        if (cur_data.length < 1000000) { min = m; }
        if (1200000 < cur_data.length) { q = Math.max(0.3, q - 0.05); max = m; }
        percentage(10 + Math.min(Math.abs((cur_data.length - 1200000) / (start_size - 1200000)), 1) * 90, 'Tuning the image, please wait...');
        setTimeout(()=>{ tuner(min, max, q); }, 5);
      }
      percentage(10, 'Tuning the image, please wait...');
      setTimeout(()=>{ tuner(0.05, 1, 0.8); }, 5);
    }
    img.src = data;
  });
}

/*
* promptCanvas() - Prompts for Canvas access on Firefox
*/
function promptCanvas() {
  log('Prompting for canvas access');
  const canvas = document.createElement('canvas');
  const ctx = canvas.getContext('2d');
  canvas.height = 10; canvas.width = 10;
  canvas.toDataURL('image/jpeg', 0.5);
}
