const statusDiv = document.getElementById('status');
const pathInput = document.getElementById('downloadPathInput');
let statusIntervalId = null;

function loadConfig() {
  chrome.runtime.sendMessage({ action: 'fetchGet', path: '/config' }, response => {
    if (response && response.success) {
      if (response.data && response.data.download_path !== undefined) {
        pathInput.value = response.data.download_path;
      }
    } else {
      statusDiv.innerText = 'Предупреждение: Не удалось подключиться к серверу для загрузки настроек';
    }
  });
}
loadConfig();

document.getElementById('toggleSettingsBtn').addEventListener('click', () => {
  const settingsMenu = document.getElementById('settingsMenu');
  settingsMenu.style.display = settingsMenu.style.display === 'block' ? 'none' : 'block';
});

document.getElementById('saveConfigBtn').addEventListener('click', () => {
  statusDiv.innerText = 'Сохранение конфигурации...';
  chrome.runtime.sendMessage({ 
    action: 'fetchPost', 
    path: '/config', 
    body: { download_path: pathInput.value } 
  }, response => {
    if (response && response.success) {
      statusDiv.innerText = 'Конфигурация успешно сохранена!';
      document.getElementById('settingsMenu').style.display = 'none';
    } else {
      statusDiv.innerText = `Ошибка сохранения: ${response?.status || 'нет связи'}`;
    }
  });
});

document.getElementById('openLogsBtn').addEventListener('click', () => {
  statusDiv.innerText = 'Открытие файла логов...';
  chrome.runtime.sendMessage({ action: 'fetchPost', path: '/logs' }, response => {
    if (response && response.success) {
      statusDiv.innerText = 'Файл логов открыт!';
    } else {
      statusDiv.innerText = `Ошибка открытия файла логов: ${response?.status || 'нет связи'}`;
    }
  });
});

document.getElementById('updateLoaderBtn').addEventListener('click', () => {
  statusDiv.innerText = 'Обновление загрузчика...';
  chrome.runtime.sendMessage({ action: 'fetchPost', path: '/update' }, response => {
    if (response && response.success) {
      statusDiv.innerText = 'Загрузчик успешно обновлен!';
    } else {
      statusDiv.innerText = `Ошибка обновления: ${response?.status || 'нет связи'}`;
    }
  });
});

document.getElementById('sendBtn').addEventListener('click', async () => {
  const formatSelect = document.getElementById('formatSelect');
  const selectedExtension = formatSelect.value;

  statusDiv.innerText = 'Проверяем вкладку...';
  
  const [tab] = await chrome.tabs.query({ active: true, currentWindow: true });
  
  if (!tab || !tab.url) {
    statusDiv.innerText = 'Не удалось определить URL страницы';
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
    statusDiv.innerText = `Вы не на странице видео/трека`;
    return;
  }

  statusDiv.innerText = 'Начало загрузки...';
  const randomId = Math.floor(Math.random() * (9999 - 1000 + 1)) + 1000;

  chrome.runtime.sendMessage({
    action: 'fetchPost',
    path: '/download',
    body: { url: cleanUrl, extension: selectedExtension, id: randomId }
  }, response => {
    if (response && response.success) {
      statusDiv.innerText = `Загрузка начата!`;
      startStatusPolling(randomId);
    } else {
      statusDiv.innerText = `Ошибка подключения: убедитесь, что go приложение запущено`;
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
      statusDiv.innerText = 'Загрузка заняла слишком много времени\nПроверьте папку загрузок';
      return;
    }

    chrome.runtime.sendMessage({
      action: 'fetchPost',
      path: '/status',
      body: { url: "", extension: "", id: id }
    }, response => {
      if (response && response.success && response.data) {
        const data = response.data;
        
        if (data.status === 'in process') {
          statusDiv.innerText = `Видео скачивается... (ID: ${id})`;
        } else if (data.status === 'error') {
          statusDiv.innerText = `Ошибка: ${data.error || 'Неизвестная ошибка при загрузке'}`;
          clearInterval(statusIntervalId);
        } else if (data.status === 'complete') {
          statusDiv.innerText = 'Загрузка завершена!';
          clearInterval(statusIntervalId);
        }
      }
    });
  }, 10000);
}
