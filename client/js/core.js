'use strict';

let view_share;

// Page Setup
onReady(()=>{
  // Click Listeners
  onClick('#btn_generate', generate_random_id);
  onClick('.btn_copy', (e)=>{
    const target = document.getElementById(e.srcElement.dataset.target);
    target.style.animation = '';
    if (navigator.clipboard) {
      navigator.clipboard.writeText(target.value).then(()=>{
        target.style.animation = 'copy 0.5s linear';
      });
    } else {
      target.select();
      document.execCommand("copy");
      setTimeout(()=>{target.style.animation = 'copy 0.5s linear';}, 1);
    }
  });
  onClick('.tab_group > .tab', (e, el)=>{
    if (!el.classList.contains('active') && !el.classList.contains('disabled')) {
      const c = el.parentElement.children;
      for (let i = c.length - 1; i >= 0; i--) {
        c[i].classList.remove('active');
      }
      el.classList.add('active');
      if (el.dataset.target && el.parentElement.dataset.tgroup) {
        if (el.dataset.target === 'tab_image' && !checkCanvas()) log('Canvas data access denied: randomization detected');
        selectAll(`div[data-group="${el.parentElement.dataset.tgroup}"]`, (t)=>{ t.classList.add('hide'); });
        document.getElementById(el.dataset.target).classList.remove('hide');
        ui_big_query_watch();
      }
    }
  });
  let is_dev = false;
  let to_dev = 7;
  onClick('#header', ()=>{
    if (is_dev) {
      error(`Last 3 entries:\n\n${glogs.map((e, i)=>{return `[${i}] ${e}`;}).slice(glogs.length-3).join('\n')}`, true);
    } else if (to_dev <= 0) {
      is_dev = true;
      error('Hello Dev!', true);
    } else {
      to_dev--;
      setTimeout(()=>{ to_dev++; }, 5E3)
    }
  });

  onClick('#btn_create', ()=>{
    // Verify entries
    if (!checkInput('id_input')) { error('Please enter a valid ID'); return; }
    switch (document.querySelector('#container_shared > .tab_group > .tab.active').dataset.target.substring(4)) {
      case 'link':
        if (!checkInput('link_data_input')) { error('Please enter a valid link'); return; }
        CreateLinkUI({
          t: 'l',
          d: getInput('link_data_input')
        });
        break;
      case 'message':
        if (!checkInput('message_data_input')) { error('Please enter a valid message'); return; }
        CreateLinkUI({
          t: 't',
          d: getInput('message_data_input')
        });
        break;
      case 'image':
        if (!checkInput('image_data_input')) { error('Please select a valid image'); return; }
        const img_data = document.getElementById('image_preview').src;
        if (img_data.length === 0 || img_data.length > 1500000) { error('Please pick a smaller image (~1.5MB)'); return; }
        CreateLinkUI({
          t: 'i',
          i: [img_data]
        });
        break;
      case 'file':
        if (!checkInput('file_data_input')) { error('Please select a valid file'); return; }
        const f = document.getElementById('file_data_input').files[0];
        if (f.size > 1500000) { error('Please pick a smaller file (~1.5MB)'); return; }
        fileReader('file_data_input').then((d)=>{
          CreateLinkUI({
            t: 'f',
            f: [{
              n: f.name,
              d: d
            }]
          });
        }).catch((e)=>{
          error('Cannot read the file data');
          console.error(e);
        });
        break;
      default:
        ui_reset();
        error(`Cannot find the selected type.`);
        return;
    }
  });

  onClick('.btn_reset', ()=>{
    ui_reset();
    view_share=false;
  });

  onClick('#btn_accept_share', ()=>{
    const errorFn = () => {
      document.getElementById('view_error_message').innerText = 'We detected something suspicious, for security reasons we prevented this share from displaying';
      hide('#data_view');
      show('#container_error');
    }
    hide('#accept_view');
    show('#data_view');
    document.getElementById('view_share_id').innerText = ''; // Clear share ID
    let parent;
    switch (view_share.t) {
      case 'l':
        setInput('link_view', view_share.d);
        if (checkInput('link_view')) {
          document.getElementById('link_view_open').href = view_share.d;
          show('#view_link');
        } else errorFn();
        break;
      case 't':
        setInput('message_view', view_share.d);
        if (checkInput('message_view')) {
          show('#view_message');
        } else errorFn();
        break;
      case 'i':
        parent = document.getElementById('view_image');
        view_share.i.forEach((i)=>{
          if (i && i.startsWith('data:image')) {
            const newimg = document.getElementById('sample_image_preview').cloneNode(true);
            newimg.id = '';
            newimg.classList.remove('hide');
            newimg.children[0].src = i;
            newimg.children[1].href = i;
            let it = i.substring(11,i.indexOf(';'));
            if (it === 'jpeg') it = 'jpg';
            newimg.children[1].download = `image.${it}`;
            parent.appendChild(newimg);
          } else errorFn();
        });
        show('#view_image');
        break;
      case 'f':
        parent = document.getElementById('view_file');
        view_share.f.forEach((f)=>{
          if (f && f.d && f.n && f.d.startsWith('data:')) {
            const newfile = document.getElementById('sample_btn_dl_file').cloneNode(true);
            newfile.id = '';
            newfile.classList.remove('hide');
            newfile.href = f.d;
            newfile.download = f.n;
            newfile.children[0].innerText = f.n.substring(0,20) + (f.n.length > 20 ? '...' : '');
            parent.appendChild(newfile);
          } else errorFn();
        });
        show('#view_file');
        break;
    }
    view_share = '';
  });

  // Dynamic Listeners
  document.getElementById('id_input').addEventListener('input', ()=>{
    const id_input = document.getElementById('id_input');
    const link = document.getElementById('id_link');
    link.value = `https://${window.location.hostname}/#${id_input.value}`
  });
  document.getElementById('message_data_input').addEventListener('input', ()=>{
    ui_big_query_watch();
  });
  const ldi = 'link_data_input';
  document.getElementById(ldi).addEventListener('change', ()=>{
    if (!checkInput(ldi)) {
      const orig = getInput(ldi);
      setInput(ldi, `https://${orig}`);
      if (!checkInput(ldi)) {
        setInput(ldi, orig);
      }
    }
  });
  const idi = 'image_data_input';
  document.getElementById(idi).addEventListener('input', ()=>{
    hide('#image_preview');
    document.getElementById(idi).parentElement.querySelector('.size').innerHTML = 0;
    if (document.getElementById(idi).files.length === 1) fileReader(idi).then((d)=>{
      const size = (x) => { return Math.round(x.length / (102.4 * 1024)) / 10; }
      document.getElementById(idi).parentElement.querySelector('.size').innerHTML = size(d);
      if (d.length > 1200000) {
        hide('#app > div');
        show('#container_loader');
        percentage(0, 'Rendering the image, please wait...');
        compressImage(d).then((data)=>{
          document.getElementById('image_preview').src = data;
          document.getElementById(idi).parentElement.querySelector('.tune').innerHTML = ` compressed to ${size(data)}MB`;
          show('#image_preview');
        }).catch((e)=>{
          error(typeof e === 'string' ? e : 'Error compressing the image');
          document.getElementById(idi).value = '';
          document.getElementById(idi).parentElement.querySelector('.size').innerHTML = 0;
        }).finally(()=>{
          hide('#container_loader');
          show('#app > .default');
        });
      } else {
        document.getElementById('image_preview').src = d;
        show('#image_preview');
      }
    });
  });
  const fdi = 'file_data_input';
  document.getElementById(fdi).addEventListener('input', ()=>{
    document.getElementById(fdi).parentElement.querySelector('.size').innerHTML = 0;
    if (document.getElementById(fdi).files.length === 1) {
      if (document.getElementById(fdi).files[0].size > 1500000) {
        document.getElementById(fdi).value = '';
        error('The file can be at most ~1.5MB');
      } else document.getElementById(fdi).parentElement.querySelector('.size').innerHTML = Math.round(document.getElementById(fdi).files[0].size / (102.4 * 1024)) / 10;
    }
  });

  window.onoffline = () => {
    error(`You're offline`);
    document.getElementById('btn_create').classList.add('disabled');
  }
  window.ononline = () => {
    if (crypto.subtle !== undefined) document.getElementById('btn_create').classList.remove('disabled');
  }

  if (crypto.subtle === undefined) {
    error('Sorry, this browser does not support web cryptography. Please download the latest version of Firefox or Chrome.');
    document.getElementById('btn_create').classList.add('disabled');
  }

  if (window.location.hash && window.location.hash.length > 1) {
    const hash = window.location.hash.substring(1);
    log('Found a share hash');
    document.getElementById('view_share_id').innerText = hash.substring(0,40) + (hash.length > 40 ? '...' : '');
    GetLinkUI(hash).then((e)=>{
      const errorFn = () => {
        document.getElementById('view_error_message').innerText = 'Unsupported Share';
        hide('#container_loader');
        show('#container_error');
      }
      if (e) {
        // Decrypted
        view_share = e;
        let stype = '';
        switch (e.t) {
          case 'l':
            stype = 'a link';
            break;
          case 't':
            stype = 'a message';
            break;
          case 'i':
            if (!e.i || !e.i.length) { errorFn(); return; }
            if (e.i.length > 1) {
              stype = 'some images';
            } else stype = 'an image';
            break;
          case 'f':
            if (!e.f || !e.f.length) { errorFn(); return; }
            if (e.f.length > 1) {
              stype = 'some files';
            } else stype = 'a file';
            break;
          default:
            errorFn();
            return;
        }
        document.getElementById('view_share_type').innerText = stype;
        hide('#container_loader');
        show('#container_view');
        window.history.replaceState( {} , 'R3', '/' ); // Remove URL from history
      } else {
        // Error
        hide('#container_loader');
        show('#container_error');
      }
    });
  } else {
    ui_reset();
  }

  log('Page ready');
});
