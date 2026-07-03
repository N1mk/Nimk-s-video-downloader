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
    const response = await fetch('http://localhost:8080', {
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
