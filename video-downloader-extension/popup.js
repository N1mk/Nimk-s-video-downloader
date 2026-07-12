let statusIntervalId = null;
let statusTimeoutId = null;
const HOST_NAME = "com.video.downloader";

document.getElementById('sendBtn').addEventListener('click', async () => {
    const statusDiv = document.getElementById('status');
    const formatSelect = document.getElementById('formatSelect');
    const selectedExtension = formatSelect.value;

    clearStatusPolling();

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
            cleanUrl = `https://youtube.com{videoId}`;
        }
    } else if (host.includes('youtu.be')) {
        const videoId = path.substring(1);
        if (videoId) {
            cleanUrl = `https://youtube.com{videoId}`;
        }
    }

    if (!cleanUrl) {
        statusDiv.innerText = `Вы не на странице плеера.\nТекущий хост: ${host}\nПуть: ${path}`;
        return;
    }

    statusDiv.innerText = 'Отправка запроса на загрузку...';

    const downloadMessage = {
        "message type": "download",
        "link": cleanUrl,
        "extension": selectedExtension
    };

    chrome.runtime.sendNativeMessage(HOST_NAME, downloadMessage, (response) => {
        if (chrome.runtime.lastError) {
            statusDiv.innerText = `Ошибка: ${chrome.runtime.lastError.message}`;
            console.error(chrome.runtime.lastError);
            return;
        }
        
        statusDiv.innerText = `Запрос отправлен [${selectedExtension.toUpperCase()}]. Начинаем проверку статуса...`;
        
        startStatusPolling(cleanUrl, statusDiv);
    });
});

function startStatusPolling(videoUrl, statusElement) {
    const pollMessage = {
        "message type": "check status",
        "link": videoUrl,
        "extension": ""
    };

    const fetchStatus = () => {
        chrome.runtime.sendNativeMessage(HOST_NAME, pollMessage, (response) => {
            if (chrome.runtime.lastError) {
                statusElement.innerText = `Ошибка опроса: ${chrome.runtime.lastError.message}`;
                clearStatusPolling();
                return;
            }

            if (response) {
                const hostStatus = response.status || 'не указан';
                const hostError = response.error || 'нет';
                
                statusElement.innerText = `Статус: ${hostStatus}\nОшибка: ${hostError}`;
            }
        });
    };

    statusIntervalId = setInterval(fetchStatus, 10000);

    statusTimeoutId = setTimeout(() => {
        clearStatusPolling();
        statusElement.innerText += '\n[Время отслеживания статуса истекло (5 мин)]';
    }, 300000);
}

function clearStatusPolling() {
    if (statusIntervalId) {
        clearInterval(statusIntervalId);
        statusIntervalId = null;
    }
    if (statusTimeoutId) {
        clearTimeout(statusTimeoutId);
        statusTimeoutId = null;
    }
}
