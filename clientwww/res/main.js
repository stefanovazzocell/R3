"use strict";

// ==========
// Settings

const API_ENDPOINT = "https://api.rkt.one/v2/";

// ==========
// API

/**
 * Decode a share data
 * @param {Uint8Array} data decrypted data
 * @returns {(Object|false)} object containing the user data, false if corrupted
 */
function DecodeShare(data) {
  if (data.length == 0) {
    console.warn("Share has no data")
    return false
  }
  if (data[0] > 2) {
    console.warn(`Unsupported share type "${data[0]}"`)
    return false
  }
  switch(data[0]) {
    case 1:
      return {
        t: 1,
        d: (new TextDecoder()).decode(data.slice(1))
      }
    case 2:
      // A file
      return {
        t: 2,
        d: decodeFiles(data.slice(1))
      }
    case 0:
      const u = validateURL((new TextDecoder()).decode(data.slice(1)));
      if (u !== false) {
        return {
          t: 0,
          d: u
        }
      }
  }
  return false;
}

/**
 * Encodes some data for sharing
 * @param {(String|FileList)} d data to encode
 * @return {Object} Object containing Uint8Array encoded data and type name
 */
async function EncodeShare(d) {
  if (typeof d == "string") {
    const u = validateURL(d);
    return {
      t: (u ? "Link" : "Message"),
      d: concatU8(new Uint8Array([(u ? 0 : 1)]), (new TextEncoder()).encode(u ? u : d))
    }
  } else {
    return {
      t: "Files",
      d: concatU8(new Uint8Array([2]), await loadFiles(d))
    }
  }
}

/**
 * Decode files
 * @param {Uint8Array} u8 the data to decode
 * @returns {(Array|false)} an Array representing a list of files or false if failure
 */
 function decodeFiles(u8) {
  let i = 0;
  let files = [];
  while (i < u8.byteLength) {
    let file = {};
    let offset = 0;
    let bytes = new Uint8Array(0);
    do {
      if (i+1 > u8.byteLength) return false;
      offset = u8[i];
      if (offset > 0) {
        if (i+1+offset > u8.byteLength) return false;
        bytes = concatU8(bytes, u8.slice(i+1, i+1+offset));
      }
      i = i+offset+1;
    } while (offset === 255);
    file.name = (new TextDecoder()).decode(bytes);
    bytes = new Uint8Array(0);
    do {
      if (i+1 > u8.byteLength) return false;
      offset = u8[i];
      if (offset > 0) {
        if (i+1+offset > u8.byteLength) return false;
        bytes = concatU8(bytes, u8.slice(i+1, i+1+offset));
      }
      i = i+offset+1;
    } while (offset === 255);
    file.mimetype = (new TextDecoder()).decode(bytes);
    bytes = new Uint8Array(0);
    do {
      if (i+3 > u8.byteLength) return false;
      offset = (u8[i])+(u8[i+1]<<8&0xff00)+(u8[i+2]<<16&0xff0000);
      if (offset > 0) {
        if (i+3+offset > u8.byteLength) return false;
        bytes = concatU8(bytes, u8.slice(i+3, i+3+offset));
      }
      i = i+offset+3;
    } while (offset === 16777215);
    file.blob = new Blob([bytes.buffer], {type: file.mimetype});
    file.blob.name = file.name;
    file.url = URL.createObjectURL(file.blob);
    files.push(file);
  }
  return files;
}

/**
 * Encode files to Uint8Array
 * @param {Array} files array representing a list of files to encode
 * @returns {Promise} Promise Uint8Array representing encoded files
 */
 async function encodeFiles(files) {
  var out = new Uint8Array(0);
  for (let i = 0; i < files.length; i++) {
    let offset = 0;
    let name = files[i].name.slice();
    if (name.length == 0) out = concatU8(out, new Uint8Array([0]));
    while (name.length > 0) {
      offset = Math.min(name.length, 255);
      out = concatU8(out, new Uint8Array([offset]));
      out = concatU8(out, new Uint8Array((new TextEncoder()).encode(name.slice(0,offset))));
      name = name.slice(offset);
    }
    let mime = files[i].type;
    if (mime.length == 0) out = concatU8(out, new Uint8Array([0]));
    while (mime.length > 0) {
      offset = Math.min(mime.length, 255);
      out = concatU8(out, new Uint8Array([offset]));
      out = concatU8(out, new Uint8Array((new TextEncoder()).encode(mime.slice(0,offset))));
      mime = mime.slice(offset);
    }
    let data = new Uint8Array(await files[i].arrayBuffer);
    if (data.byteLength == 0) out = concatU8(out, new Uint8Array([0]));
    while (data.byteLength > 0) {
      offset = Math.min(data.byteLength, 16777215);
      out = concatU8(out, new Uint8Array([offset&0xff, offset>>8&0xff, offset>>16&0xff]));
      out = concatU8(out, data.slice(0,offset));
      data = data.slice(offset);
    }
  }
  return out;
}

// ==========
// Crypto

/**
 * Computes a secure hash using PBKDF2
 * @param {String} key the user starting password
 * @param {boolean} limitOutput true if output should be limited to first 6 bytes, false otherwise
 * @returns {Promise} Promise a Uint8Array containing the hash
 */
async function Hash(key, limitOutput) {
  const passwordBuffer = (new TextEncoder()).encode(`hashgen::${key}::salty`);
  const importedKey = await crypto.subtle.importKey("raw", passwordBuffer, "PBKDF2", false, ["deriveBits"]);

  const bytes = await crypto.subtle.deriveBits({
    name:     "PBKDF2",
    hash:     "SHA-512",
    salt:     new Uint8Array([82, 242, 11, 190, 119, 15, 58, 152, 115, 230, 184, 149, 107, 12, 5, 37, 184, 242, 159, 111, 72, 180, 65, 53, 104, 78, 252, 123, 188, 17, 71, 187, 216, 128, 141, 148, 126, 110, 15, 113, 175, 70, 216, 37, 211, 247, 93, 216, 210, 197, 189, 100, 37, 81, 113, 113, 173, 8, 184, 97, 225, 223, 24, 69]),
    iterations: 1000000
  }, importedKey, 512);

  return ToBase64(bytes.slice(0, (limitOutput ? 6 : undefined)))
}

/**
 * Computes a secure hash using PBKDF2
 * @param {String} key the user starting password
 * @param {string=""} salt the optional base64 salt to use, if undefined generates random salt
 * @returns {Promise} Promises an object with a key and the base64 encoded salt
 */
async function KeyDerivation(key, salt) {
  const keyBuffer = (new TextEncoder()).encode(`keyder::${key}::salted`);
  const importedKey = await crypto.subtle.importKey("raw", keyBuffer, "PBKDF2", false, ["deriveKey"]);
  const saltBuffer = (salt ? salt : window.crypto.getRandomValues(new Uint8Array(32)));

  return {
    key: await crypto.subtle.deriveKey({
      name:       "PBKDF2",
      hash:       "SHA-256",
      salt:       saltBuffer,
      iterations: 1000000
    }, importedKey, { name: "AES-GCM", length: 256 }, false, [ "encrypt", "decrypt" ]),
    salt: saltBuffer
  };
}

/**
 * Encrypts some data with a given key
 * @param {Uint8Array} data 
 * @param {Object} keyObj 
 * @returns {Promise} Promise containing base64 encoded salt, iv, and encrypted data
 */
async function Encrypt(data, keyObj) {
  const iv = window.crypto.getRandomValues(new Uint8Array(12));
  const encryptedContent = await window.crypto.subtle.encrypt({
    name: "AES-GCM",
    iv: iv,
  }, keyObj.key, data);

  let out = new Uint8Array(44 + encryptedContent.byteLength);
  out.set(new Uint8Array(keyObj.salt), 0);
  out.set(new Uint8Array(iv), 32);
  out.set(new Uint8Array(encryptedContent), 44);
  return ToBase64(out);
}

/**
 * Decrypts some data with a given key
 * @param {String} data 
 * @param {String} key 
 * @returns {Promise} Promise containing Uint8Array data
 */
async function Decrypt(data, key) {
  const dataBuffer = FromBase64(data);
  const keyObj = await KeyDerivation(key, dataBuffer.slice(0, 32));

  return new Uint8Array(await window.crypto.subtle.decrypt({
    name: "AES-GCM",
    iv: dataBuffer.slice(32, 44),
  }, keyObj.key, dataBuffer.slice(44)));
}

// ==========
// Utilities

/**
 * Converts an Uint8Array to a base64 string
 * @param {Uint8Array} u8 the data to covert to base64
 * @returns {Promise} string of bytes encoded in Base64 
 */
async function ToBase64(u8) {
  const b64url = await new Promise((r) => {
    const reader = new FileReader()
    reader.onload = () => r(reader.result)
    reader.readAsDataURL(new Blob([u8]))
  })
  return b64url.split(",", 2)[1]
}

/**
 * Converts a base64 string to an Uint8Array
 * @param {String} str the base64 string
 * @returns {Uint8Array} the array of bytes from the decoded base64 
 */
function FromBase64(str) {
  return new Uint8Array(atob(str).split("").map((c) => { return c.charCodeAt(0); }));
}

/**
 * 
 * @param {String} path path to make the request to
 * @param {Object} data to send to server
 * @returns {Promise} Promise Object containing results, if false an error has occurred 
 */
async function Query(path, data) {
  try {
    const res = await fetch(`${API_ENDPOINT}${path}`, {
      method: "POST",
      headers: {
        "Content-Type": "application/json"
      },
      body: (data ? JSON.stringify(data) : undefined)
    })
    if (res.status != 200) {
      console.error(`Got ${res.status} status`);
      return false;
    }
    return res.json()
  } catch (e) {
    console.error(e);
    return false;
  }
}

/**
 * Validates a url or returns false
 * @param {String} url the url to validate
 * @returns {(string|false)} a valid url from url or false if not a url
 */
function validateURL(url) {
  if (url.length < 4 || url.length > 2048 || url.indexOf(" ") != -1 || url.indexOf("\n") != -1) {
    return false;
  }
  const i = document.createElement("input", {"type": "url", "value": url});
  i.type = "url";
  i.value = (url.startsWith("http") ? url : `https://${url}`);
  return (i.checkValidity() ? i.value : false);
}

/**
 * Concatenates two Uint8Array
 * @param {Uint8Array} a the first Uint8Array to concatenate 
 * @param {Uint8Array} b the second Uint8Array to concatenate 
 * @returns {Uint8Array} Uint8Array concatenated from a, b
 */
function concatU8(a, b) {
  var out = new Uint8Array(a.byteLength + b.byteLength);
  out.set(new Uint8Array(a), 0);
  out.set(new Uint8Array(b), a.byteLength);
  return out;
}

/**
 * Gets a secure random int in a given range [min, max)
 * @param {Number} min minimum (included)
 * @param {Number} max maximum (excluded)
 * @returns {Number} a number in the given range
 */
function randomInt(min, max) {
  return Math.floor((window.crypto.getRandomValues(new Uint32Array(1))[0] / 2**32) * (max - min) + min);
}

/**
 * Generate a random secure string
 * @param {Number} len length of the string to generate
 * @returns {String} the generated string
 */
function randomString(len) {
  const chars = "abcdefghijkmnopqrstuvwxyzABCDEFGHJKLMNOPQRSTUVWXYZ0123456789~";
  let out = "";
  for (let i = len; i > 0; i--) {
    out += chars.charAt(randomInt(0, 61)).replace("~", "");
  }
  return out;
}

/**
 * Generates a random key string
 * @returns {String}
 */
function randomKey() {
  return [randomString(3), randomString(3), randomString(3), randomString(3)].join(".");
}

/**
 * Reads all file and returns an encoded Uint8Array
 * @param {FileList} files files to read
 * @returns {Promise} A promise that resolves in an encoded Uint8Array
 */
async function loadFiles(files) {
  let promises = [];
  for (let i = 0; i < files.length; i++) {
    promises.push(loadFile(files[i]));
  }
  return encodeFiles(await Promise.all(promises));
}

/**
 * Reads a file and returns its data
 * @param {File} file file to read
 * @returns {Promise} A promise that resolves in the file data in Uint8Array
 */
async function loadFile(file) {
  return new Promise((resolve)=>{
    const reader = new FileReader();
    reader.onload = (e)=>{
      resolve({
        name: file.name,
        type: file.type,
        arrayBuffer: e.target.result
      });
    }
    reader.readAsArrayBuffer(file);
  });
}

// ==========
// UI

const shareText = select("#stext");
const shareFile = select("#sfiles");
const shareID = select("#sid");
const shareTTL = select("#sttl");
const shareViews = select("#sviews");
const shareEditPass = select("#sepass");
const shareURL = select("#surl");
const omni = select("#omni");
const filesBox = select("#filesbox");
const shareBtn = select("#s");
const sharedPanel = select("#sd");
const sharingPanel = select("#sharing");
const createPanel = select("#create");
const viewPanel = select("#view");
const rurl = select('#rurl');
const rurld = select('#rurld');
const rtext = select('#rtext');
const rfiles = select('#rfiles');

/**
 * Executes callback when page ready
 * @param {Function} callback to call when page is ready
 */
function onReady(callback) {
  if(document.readyState === "interactive" || document.readyState === "complete") {
    callback();
  } else document.addEventListener("DOMContentLoaded", callback);
}

/**
 * Query multiple elements and loop through them
 * @param {String} query the query string
 * @param {Function} fn the function to call
 */
function selectAll(query, fn) {
  Array.from(document.querySelectorAll(query)).forEach(fn);
}

/**
 * Query an element
 * @param {String} query the query string
 * @returns {Element} the requested element
 */
function select(query) {
  return document.querySelector(query);
}

/**
 * Handles clicks for an element
 * @param {String} sel the selector for the element
 * @param {Function} fn the function to call
 */
function onClick(sel, fn) {
  selectAll(sel, (el)=>{ el.addEventListener("click", fn); })
}

/**
 * Shows an alert
 * @param {String} msg the message to show in the alert 
 */
function alert(msg) {
  const el = document.getElementById("sample").cloneNode(true);
  el.id = "";
  el.children[0].innerText = msg;
  el.onclick = ()=>{ el.remove(); };
  setTimeout(()=>{ el.remove(); }, 2000 + 250 * Math.ceil(msg.split(" ").length));
  document.getElementById("alerts").appendChild(el);
  console.warn(msg);
}

/**
 * Calculates share stats, updates the UI, returns the share data
 * @returns {Promise} a share Object or nil if no data available
 */
async function shareStats() {
  var d = shareText.value;
  if (d.length == 0) {
    d = shareFile.files;
    if (d.length > 0) progress("Loading files...");
  }
  selectAll(".bq",(el)=>{ el.disabled = false; });
  if (d.length == 0) {
    shareBtn.disabled = true;
    select("#sp").innerHTML = `Share`;
    select("#storagep").innerHTML = `0% Storage Used`;
    filesBox.classList.add("hide");
    selectAll("#filesbox > div:not(#samplef)", (e)=>{e.remove();});
    shareText.classList.remove("hide");
    return;
  } else {
    const share = await EncodeShare(d);
    if (share.d.byteLength > 10485760) {
      shareBtn.disabled = true;
      select("#sp").innerHTML = "Oversize Share";
    } else {
      select("#sp").innerHTML = share.t;
      shareBtn.disabled = false;
    }
    if (share.d.byteLength > 10240) {
      selectAll(".bq",(el)=>{ el.disabled = true; });
      shareTTL.value = Math.min(3600, shareTTL.value);
    }
    select("#storagep").innerHTML = `${Math.ceil(share.d.byteLength / 104857.6)}% Storage Used`;
    if (share.t == "Files") {
      shareText.classList.add("hide");
      selectAll("#filesbox > div:not(#samplef)", (e)=>{e.remove();});
      for (let i = 0; i < shareFile.files.length; i++) {
        const el = document.getElementById("samplef").cloneNode(true);
        el.id = "";
        el.innerText = shareFile.files[i].name;
        filesBox.appendChild(el);
      }
      progress("");
      filesBox.classList.remove("hide");
    } else {
      filesBox.classList.add("hide");
      shareFile.value = "";
      selectAll("#filesbox > div:not(#samplef)", (e)=>{e.remove();});
      shareText.classList.remove("hide");
    }
    return share;
  }
}

/**
 * Resets the UI
 */
function resetUI() {
  shareText.value = "";
  shareFile.value = "";
  shareURL.value = "";
  rurld.value = "";
  rurl.innerText = "";
  rtext.innerText = "";
  rfiles.childNodes.forEach((el)=>{if (el.href.length > 1) URL.revokeObjectURL(el.href)});
  rfiles.innerText = "";
  shareTTL.value = 3600;
  shareViews.value = 10;
  shareEditPass.value = "";
  shareID.value = randomKey();
  document.body.classList.remove("shareurl");
  sharingPanel.classList.add("hide");
  createPanel.classList.remove("hide");
  viewPanel.classList.add("hide");
  sharingPanel.innerText = "";
  sharedPanel.classList.add("hide");
  shareStats();
  select("#s").disabled = true;
}

/**
 * Shows the progress UI with a message
 * @param {String} msg the message to display
 */
function progress(msg) {
  sharingPanel.innerText = msg;
  if (msg.length > 0) {
    sharingPanel.classList.remove("hide");
    document.body.classList.add("shareurl");
  } else {
    document.body.classList.remove("shareurl");
    sharingPanel.classList.add("hide");
  }
}

// ==========
// Main

onReady(()=>{
  onClick("#header > a", resetUI);
  onClick("#rID", ()=>{ shareID.value = randomKey(); });
  onClick("#rPass", ()=>{ shareEditPass.value = randomString(32); });
  onClick("#clearp", ()=>{ shareText.value = ""; shareFile.value = ""; shareStats(); });
  onClick(".copy", (e)=>{ navigator.clipboard.writeText(select("#"+e.target.dataset["t"]).value).then(()=>{ resetUI(); }); })
  shareText.addEventListener("input", ()=>{ if (shareText.value.length < 10000) shareStats(); });
  shareText.addEventListener("change", ()=>{ shareStats(); });
  shareFile.addEventListener("change", ()=>{ if (shareFile.files.length > 0) { shareText.value = ""; shareStats(); } });
  omni.addEventListener("dblclick", ()=>{ shareFile.click(); });
  omni.ondragover = omni.ondragenter = (e)=>{e.preventDefault();};
  omni.ondrop = (e)=>{
    if (e.dataTransfer.files.length == 0) {
      shareText.value = e.dataTransfer.getData("text");
    } else {
      shareFile.files = e.dataTransfer.files;
      shareText.value = "";
    }
    e.preventDefault();
    shareStats();
  };
  filesBox.addEventListener("click", ()=>{ shareFile.click(); });
  onClick("#s", async ()=>{
    progress("Loading Share...");
    if (!shareID.checkValidity())  { progress(""); alert("Check the Share ID"); return; }
    const share = await shareStats();
    if (!share) { progress(""); alert("Can't create an empty share"); return; }
    if (share.d.byteLength > 10485760) { progress(""); alert("Your share is too large"); return; }
    progress("Encrypting Share...");
    const shareHash = Hash(shareID.value, true);
    const shareData = await Encrypt(share.d, await KeyDerivation(shareID.value));
    progress("Uploading Share...");
    const res = await Query("edit", {
      id: await shareHash,
      delete:  false,
      pass:    "",
      payload: {
        data: shareData,
        ttl:  Number(shareTTL.value),
        hits: Number(shareViews.value),
        edit: (shareEditPass.value.length > 0 ? await Hash(shareEditPass.value, false) : "")
      }
    });
    if (!res) { progress(""); alert("Network error, try again later."); return; }
    if (!res.success) { progress(""); alert(res.error); return; }
    shareURL.value = `${window.location.protocol}//${window.location.hostname}/#${shareID.value}`;
    sharingPanel.classList.add("hide");
    sharedPanel.classList.remove("hide");
  });

  if (window.location.hash.length > 1) {
    createPanel.classList.add("hide");
    viewPanel.classList.remove("hide");
    progress("Downloading Share...");
    setTimeout(async ()=>{
      const hash = window.location.hash.slice(1);
      history.replaceState('', '', '/');
      const res = await Query("get", {
        id:   await Hash(hash, true),
        pass: ""
      });
      if (!res) { progress(""); alert("Network error, try again later."); resetUI(); return; }
      if (!res.success) { progress(""); alert(res.error); resetUI(); return; }
      progress("Decrypting Share...");
      const data = await Decrypt(res.data, hash);
      progress("Loading Share...");
      const share = DecodeShare(data);
      if (!share) { progress(""); alert("This share is corrupted."); resetUI(); return; }
      if (share.t == 0) {
        rurld.value = share.d;
        rurl.classList.remove('hide');
      } else if (share.t == 1) {
        rtext.innerText = share.d;
        rtext.classList.remove('hide');
      } else {
        if (share.d.length > 1) {
          const a = document.createElement('a');
          a.href = "#";
          a.innerHTML = "<b>Save All</b>";
          a.onclick = ()=>{
            rfiles.childNodes.forEach((el)=>{if (el != a) el.click();})
          };
          rfiles.appendChild(a);
        }
        for(let i = 0; i < share.d.length; i++) {
          const a = document.createElement('a');
          a.download = share.d[i].name;
          a.href = share.d[i].url;
          a.innerText = `Save "${share.d[i].name}"`;
          rfiles.appendChild(a);
        }
        rfiles.classList.remove('hide');
      }
      progress("");
    }, 1);
  } else {
    resetUI();
  }
});