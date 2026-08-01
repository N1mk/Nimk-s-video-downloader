const SERVER_URL = 'http://127.0.0.1:8080';

chrome.runtime.onMessage.addListener((message, sender, sendResponse) => {
  handleMessageAsync(message, sendResponse);
  return true; 
});

async function handleMessageAsync(message, sendResponse) {
  try {
    if (message.action === 'fetchGet') {
      const response = await fetch(`${SERVER_URL}${message.path}`);
      if (!response.ok) return sendResponse({ success: false, status: response.status });
      const data = await response.json();
      sendResponse({ success: true, data });
    }

    if (message.action === 'fetchPost') {
      const response = await fetch(`${SERVER_URL}${message.path}`, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify(message.body || {})
      });
      if (!response.ok) return sendResponse({ success: false, status: response.status });
      
      let data = null;
      const text = await response.text();
      if (text) data = JSON.parse(text);
      sendResponse({ success: true, data });
    }
  } catch (error) {
    sendResponse({ success: false, error: error.message });
  }
}
