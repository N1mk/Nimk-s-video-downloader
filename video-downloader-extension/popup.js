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
      console.error('Не удалось загрузить конфиг:', response.status);
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
      document.getElementById('settingsMenu').style.style.display = 'none';
    } else {
      statusDiv.innerText = `Ошибка сохранения: ${response.status} ${response.statusText}`;
    }
  } catch (error) {
    console.error(error);
  }
});

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
    statusDiv.innerText = `Вы не на странице плеера.\nТекущий хост: ${host}\nПуть: ${path}`;
    return;
  }

  statusDiv.innerText = 'Отправка запроса на загрузку...';

  try {
    const response = await fetch('http://localhost:8080/download', {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json'
      },
      body: JSON.stringify({ 
        url: cleanUrl,
        extension: selectedExtension 
      })
    });

    if (response.ok) {
      statusDiv.innerText = `Загрузка начата [${selectedExtension.toUpperCase()}]!\n${cleanUrl}`;
    } else {
      statusDiv.innerText = `Ошибка: ${response.status} ${response.statusText}`;
    }
  } catch (error) {
    statusDiv.innerText = `Ошибка подключения: убедитесь, что go приложение запущено`;
    console.error(error);
  }
});
