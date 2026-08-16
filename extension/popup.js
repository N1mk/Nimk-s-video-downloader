const statusDiv = document.getElementById('status');
const downloadPathInput = document.getElementById('downloadPathInput');
const maxRetryCountInput = document.getElementById('maxRetryCountInput')
let statusIntervalId = null;

const jobStatusInProcess = 0
const jobStatusError     = 1
const jobStatusRetrying  = 2
const jobStatusComplete  = 3
const selectEl = document.getElementById('formatSelect');
const wrapperEl = document.getElementById('selectWrapper');


if (selectEl && wrapperEl) {
  selectEl.addEventListener('focus', () => {
    wrapperEl.classList.add('is-open');
  });

  selectEl.addEventListener('change', () => {
    selectEl.blur();
  });

  selectEl.addEventListener('blur', () => {
    wrapperEl.classList.remove('is-open');
  });
}

function loadConfig() {
  chrome.runtime.sendMessage({ action: 'fetchGet', path: '/config' }, response => {
    if (response && response.success) {
      if (response.data && response.data.download_path !== undefined) {
        downloadPathInput.value = response.data.download_path;
        maxRetryCountInput.value = response.data.max_retry_count;
      }
    } else {
      statusDiv.innerText = 'Warning: Cannot connect to the application';
    }
  });
}
loadConfig();

document.getElementById('toggleSettingsBtn').addEventListener('click', () => {
  const settingsMenu = document.getElementById('settingsMenu');
  settingsMenu.style.display = settingsMenu.style.display === 'block' ? 'none' : 'block';
});

document.getElementById('saveConfigBtn').addEventListener('click', () => {
  statusDiv.innerText = 'Saving the configuration...';
  chrome.runtime.sendMessage({ 
    action: 'fetchPost', 
    path: '/config', 
    body: { download_path: downloadPathInput.value, max_retry_count: maxRetryCountInput.value } 
  }, response => {
    if (response && response.success) {
      statusDiv.innerText = 'Configuration successfully saved!';
    } else {
      statusDiv.innerText = `Save error: ${response?.status || 'no conntection'}`;
    }
  });
});

document.getElementById('openLogsBtn').addEventListener('click', () => {
  statusDiv.innerText = 'Opening the log file...';
  chrome.runtime.sendMessage({ action: 'fetchPost', path: '/logs' }, response => {
    if (response && response.success) {
      statusDiv.innerText = 'Log file opened!';
    } else {
      statusDiv.innerText = `Log file opening error: ${response?.status || 'no connection'}`;
    }
  });
});

document.getElementById('updateLoaderBtn').addEventListener('click', () => {
  statusDiv.innerText = 'Updating the loader...';
  chrome.runtime.sendMessage({ action: 'fetchPost', path: '/update' }, response => {
    if (response && response.success) {
      statusDiv.innerText = 'Loader successfully updated!';
    } else {
      statusDiv.innerText = `Update error: ${response?.status || 'no connection'}`;
    }
  });
});

document.getElementById('sendBtn').addEventListener('click', async () => {
  const formatSelect = document.getElementById('formatSelect');
  const selectedExtension = formatSelect.value;

  statusDiv.innerText = 'Checking the tab...';
  
  const [tab] = await chrome.tabs.query({ active: true, currentWindow: true });
  
  if (!tab || !tab.url) {
    statusDiv.innerText = 'Cannot determine the page URL';
    return;
  }

  const urlObj = new URL(tab.url);
  let cleanUrl = '';
  const host = urlObj.hostname;
  const path = urlObj.pathname;

  if ((host.includes('youtube.com') || host.includes('://youtube.com')) && path === '/watch') {
    const videoId = urlObj.searchParams.get('v');
    if (videoId) cleanUrl = `https://youtube.com/watch?v=${videoId}`;
  } else if (host.includes('youtu.be')) {
    const videoId = path.substring(1);
    if (videoId) cleanUrl = `https://youtube.com/watch?v=${videoId}`;
  } else if (host.includes('soundcloud.com')) {
    const isPlaylist = path.includes('/sets') || path.includes('/albums') || path.includes('/playlists');
    const systemPaths = ['/discover', '/stream', '/charts', '/search', '/upload', '/messages', '/notifications', '/you'];
    const isSystem = systemPaths.some(p => path.startsWith(p));
    const pathSegments = path.split('/').filter(Boolean);
    if (!isPlaylist && !isSystem && pathSegments.length === 2) {
      cleanUrl = `https://soundcloud.com/${pathSegments[0]}/${pathSegments[1]}`;
    }
  } else if (host.includes('rutube.ru')) {
    const rutubeRegex = /^\/video\/([a-zA-Z0-9]+)/;
    const match = path.match(rutubeRegex);
    if (match) cleanUrl = `https://rutube.ru/video/${match[1]}/?r=wd`;
  }

  if (!cleanUrl) {
    statusDiv.innerText = `You are not on the player page`;
    return;
  }

  statusDiv.innerText = 'Starting download...';
  const randomId = Math.floor(Math.random() * (9999 - 1000 + 1)) + 1000;

  chrome.runtime.sendMessage({
    action: 'fetchPost',
    path: '/download',
    body: { url: cleanUrl, extension: selectedExtension, id: randomId }
  }, response => {
    if (response && response.success) {
      statusDiv.innerText = `Download started!`;
      startStatusPolling(randomId);
    } else {
      statusDiv.innerText = `Conection error: make sure the application is running`;
    }
  });
});

function startStatusPolling(id) {
  if (statusIntervalId) clearInterval(statusIntervalId);

  const startTime = Date.now();
  const maxDuration = 5 * 60 * 1000;

  statusIntervalId = setInterval(() => {
    if (Date.now() - startTime > maxDuration) {
      clearInterval(statusIntervalId);
      statusDiv.innerText = 'Download took too much time\nCheck the destination';
      return;
    }

    chrome.runtime.sendMessage({
      action: 'fetchPost',
      path: '/status',
      body: { url: "", extension: "", id: id }
    }, response => {
      if (response && response.success && response.data) {
        const data = response.data;
        
        if (data.status === jobStatusInProcess) {
          statusDiv.innerText = `Video downloading... (ID: ${id})`;
        } else if (data.status === jobStatusError) {
          statusDiv.innerText = `Error: ${data.error || 'unknown error while downloading'}`;
          clearInterval(statusIntervalId);
        } else if (data.status === jobStatusComplete) {
          statusDiv.innerText = 'Download complete!';
          clearInterval(statusIntervalId);
        }
      }
    });
  }, 10000);
}