async function loadConfig() {
  const statusDiv = document.getElementById('status');
  const pathInput = document.getElementById('downloadPathInput');
  
  try {
    const response = await fetch('http://localhost:8080/config', {
      method: 'GET'
    });
    
    if (response.ok) {
      const data = await response.json();
      if (data && data.download_path !== undefined) {
        pathInput.value = data.download_path;
      }
    } else {
      console.error('Не удалось открыть конфиг:', response.status);
    }
  } catch (error) {
    statusDiv.innerText = 'Предупреждение: Не удалось подключиться к localhost:8080 для загрузки настроек';
    console.error(error);
  }
}

loadConfig();

document.getElementById('toggleSettingsBtn').addEventListener('click', () => {
  const settingsMenu = document.getElementById('settingsMenu');
  if (settingsMenu.style.display === 'block') {
    settingsMenu.style.display = 'none';
  } else {
    settingsMenu.style.display = 'block';
  }
});

document.getElementById('saveConfigBtn').addEventListener('click', async () => {
  const statusDiv = document.getElementById('status');
  const pathInput = document.getElementById('downloadPathInput');
  
  statusDiv.innerText = 'Сохранение конфигурации...';
  
  try {
    const response = await fetch('http://localhost:8080/config', {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json'
      },
      body: JSON.stringify({
        download_path: pathInput.value
      })
    });
    
    if (response.ok) {
      statusDiv.innerText = 'Конфигурация успешно сохранена!';
      document.getElementById('settingsMenu').style.display = 'none'; 
    } else {
      statusDiv.innerText = `Ошибка сохранения: ${response.status} ${response.statusText}`;
    }
  } catch (error) {
    console.error(error);
  }
});

document.getElementById('openLogsBtn').addEventListener('click', async () => {
  const statusDiv = document.getElementById('status');
  statusDiv.innerText = 'Открытие файла логов...';
  
  try {
    const response = await fetch('http://localhost:8080/logs', {
      method: 'POST'
    });
    
    if (response.ok) {
      statusDiv.innerText = 'Файл логов открыт!';
    } else {
      statusDiv.innerText = `Ошибка открытия файла логов: ${response.status} ${response.statusText}`;
    }
  } catch (error) {
    statusDiv.innerText = 'Ошибка подключения к серверу логов';
    console.error(error);
  }
});

let statusIntervalId = null;

document.getElementById('sendBtn').addEventListener('click', async () => {
  const statusDiv = document.getElementById('status');
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

  if (host.includes('youtube.com') && path === '/watch') {
    const videoId = urlObj.searchParams.get('v');
    if (videoId) {
      cleanUrl = `https://youtube.com/watch?v=${videoId}`;
    }
  } else if (host.includes('youtu.be')) {
    const videoId = path.substring(1);
    if (videoId) {
      cleanUrl = `https://youtube.com/watch?v=${videoId}`;
    }
  }

  if (!cleanUrl) {
    statusDiv.innerText = `Вы не на странице плеера`;
    return;
  }

  statusDiv.innerText = 'Начало загрузки...';

  const randomId = Math.floor(Math.random() * (9999 - 1000 + 1)) + 1000;

  try {
    const response = await fetch('http://localhost:8080/download', {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json'
      },
      body: JSON.stringify({ 
        url: cleanUrl,
        extension: selectedExtension,
        id: randomId
      })
    });

    if (response.ok) {
      statusDiv.innerText = `Загрузка начата!`;
      
      startStatusPolling(randomId);
    } else {
      statusDiv.innerText = `Ошибка: ${response.status} ${response.statusText}`;
    }
  } catch (error) {
    statusDiv.innerText = `Ошибка подключения: убедитесь, что go приложение запущено`;
    console.error(error);
  }
});

function startStatusPolling(id) {
  const statusDiv = document.getElementById('status');
  
  if (statusIntervalId) {
    clearInterval(statusIntervalId);
  }

  const startTime = Date.now();
  const maxDuration = 5 * 60 * 1000;

  statusIntervalId = setInterval(async () => {
    if (Date.now() - startTime > maxDuration) {
      clearInterval(statusIntervalId);
      statusDiv.innerText = 'Загрузка заняла слишком много времени\nПроверьте папку загрузок';
      return;
    }

    try {
      const response = await fetch('http://localhost:8080/status', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json'
        },
        body: JSON.stringify({
          url: "",
          extension: "",
          id: id
        })
      });

      if (response.ok) {
        const data = await response.json();
        
        if (data.status === 'in process') {
          statusDiv.innerText = `Видео скачивается... (ID: ${id})`;
        } else if (data.status === 'error') {
          statusDiv.innerText = `Ошибка: ${data.error || 'Неизвестная ошибка при загрузке'}`;
          clearInterval(statusIntervalId); 
        } else if (data.status === 'complete') {
          statusDiv.innerText = 'Видео скачалось!';
          clearInterval(statusIntervalId); 
        }
      } else {
        console.error('Ошибка проверки статуса загрузки:', response.status);
      }
    } catch (error) {
      console.error('Ошибка сети при проверки статуса загрузки:', error);
    }
  }, 10000);
}