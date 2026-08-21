## **Nimk's Video Downloader**
A simple video and audio downloader consisting of a browser extension and a Go application. It allows you to easily download video and music from YouTube, YouTube Music, SoundCloud, and Rutube.

---
(ENG)

### 🛠 Technology Stack
* **Go 1.26.4**
* **yt-dlp**
* **FFmpeg**
* **HTML5**
* **JavaScript**

### 🚀 Installation Guide

#### Important: For the downloader to work, you must install both the browser extension (Chrome/Firefox) and the desktop application (Windows/Linux/macOS)!

* #### Extension Installation (Chromium):
1. Extract the extension archive.
2. Open a Chromium-based browser.
3. Open the extension settings via the browser menu or go to `[browser name]://extensions`.
4. Enable Developer Mode (toggle in the top-right corner).
5. Click "Load unpacked" and select the folder with the extracted extension.

* #### Extension Installation (Firefox):
1. Open a Firefox-based browser.
2. Open the extension settings via the browser menu or go to `about:addons`.
3. Click the gear icon and select "Install Add-on from File...".
4. Select the `nvd-extension-firefox.xpi` file.

* #### App Installation:
1. Run the `video-downloader-build.exe` file.
2. Wait for the installation to complete. (The first launch may take a few minutes as it downloads dependencies).
3. The application will be automatically added to startup.

* *If any errors occur, update the loader via the extension settings and check the `app.log` file, or open it directly through the extension.*

### 💻 Usage Example
1. Open a video or audio player page on YouTube, YouTube Music, SoundCloud, or Rutube.
2. Open the extension menu and select your preferred format.
3. Select the download destination. (Any missing folders will be created automatically).
4. Click the "Download" button and wait.
5. Check your chosen folder.

#### Notes
* *When downloading large videos in high quality, select the default platform format (like WebM for YouTube or MP4 for Rutube) because converting takes too much time and resources.*
* *Temporary files (`.part`) or videos in transitional formats may appear during the download process. Please do not modify or delete them; they will be automatically removed once the download is complete.*

#### Third-Party Licenses
* This project uses **yt-dlp** (`yt-dlp.exe`) to download videos. yt-dlp is licensed under *The Unlicense*. Official website: [github.com/yt-dlp/yt-dlp](https://github.com/yt-dlp/yt-dlp).
* This project uses **FFmpeg** (`ffmpeg.exe`) to convert video formats. FFmpeg is licensed under the *GNU Lesser General Public License (LGPL) version 2.1*. However, this project utilizes a build from `https://github.com/BtbN/FFmpeg-Builds/releases/download/latest/ffmpeg-master-latest-win64-gpl.zip`, which is licensed under the *GNU General Public License (GPL)*. Official website: [www.ffmpeg.org](https://ffmpeg.org). Source code: [github.com/ffmpeg/ffmpeg](https://github.com/ffmpeg/ffmpeg).

#### License
This project is licensed under the **GPL-3.0 License** — see the `LICENSE` file for details.

---

(RU)

### 🛠 Стек технологий

* **Go 1.26.4**
* **yt-dlp**
* **ffmpeg**
* **html5**
* **javascript**

### 🚀 Гайд по установке

#### Важно: Для работы загрузчика требуются расширение (Chrome/Firefox) и приложение (Windows/Linux/MacOS)!

* #### Установка расширения (chromium):
1. Распакуйтее фрхив с расширением
2. Откройте браузер на основе chromium
3. Перейдите в настройки расширений через UI или по ссылке `[название браузера]://extensions`
4. Включите режим разработчика (переключатель в правом верхнем углу)
5. Нажмите кнопку "Загрузить распакованное расширение" и выберите папку распакованного расширения

* #### Установка расширения (firefox):
1. Откройте браузер на основе firefox
2. Перейдите в настройки расширений через UI или по ссылке `about:addons`
3. Нажмите на шестерёнку и выберите пункт "Установить дополнение из файла..."
4. Выберите скаченный ранее файл `nvd-extension-firefox.xpi`

* #### Установка приложения:
1. Запустите файл `video-downloader-build.exe`
2. Дождитесь окончания загрузки (первый запуск идёт долго из-за установки зависимостей)
3. С этого момента приложение запустится и добавится в автозагрузку

* *В случае возникновения ошибок обновите загрузчик через настройки расширения и проверьте файл логов `app.log` или откорйте его с помощью расширения*

### 💻 Пример использования

1. Зайдите на страницу с любым видео или треком на Youtube, Youtube Music, Soundcloud или Rutube
2. Откройте меню расширения и выберите формат загружаемого файла
3. Выберите директорию загрузки в настройках (все недостающие папки будут созданы)
4. Нажмите кнопку "Скачать" и ожидайте загрузки (чем больше видео, тем дольше она будет идти)
5. Проверьте папку, которую указывали в настройках

#### Примечания:
* *При загрузке больших видео в высоком разрешении выбирайте стандартный формат платформы (например, WebM для YouTube или MP4 для Rutube), потому что конвертация требует много времени и ресурсов*
* *Во время загрузки могут создаваться временные файлы (`.part`) или видео в неправильном формате, не трогайте их, они сами удалятся после окончания загрузки*

#### Лицензии третьих сторон

* Этот проект использует yt-dlp (`yt-dlp.exe`) для загрузки видео. Yt-dlp имеет лицензию *The Unlicense*. Официальный сайт yt-dlp: [github.com/yt-dlp/yt-dlp](https://github.com/yt-dlp/yt-dlp).
* Этот проект использует ffmpeg (`ffmpeg.exe`) для конвертации форматов видео. FFmpeg имеет лицензию *GNU Lesser General Public License (LGPL) version 2.1*, однако
проект использует билд из `https://github.com/BtbN/FFmpeg-Builds/releases/download/latest/ffmpeg-master-latest-win64-gpl.zip`, который использует *GNU General Public License (GPL)*.
Официальный сайт ffmpeg: [www.ffmpeg.org](https://ffmpeg.org). Исходный код ffmpeg: [github.com/ffmpeg/ffmpeg](https://github.com/ffmpeg/ffmpeg).

#### Лицензия

Этот проект распространяется под лицензией GPL-3.0. Подробнее смотерть в файле `LICENSE`.