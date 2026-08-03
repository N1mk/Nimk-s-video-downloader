## **Nimk's Video Downloader**

Удобный загрузчик видео и аудио, состоящий из браузерного расширения и go-приложения.
С его помощью вы сможете скачивать видео с YouTube или YouTube Music в разных форматах.

### 🛠 Стек технологий

* **go 1.26.4**
* **yt-dlp**
* **ffmpeg**
* **html5**
* **javascript**

### 🚀 Гайд по установке

#### ДЛЯ РАБОТЫ ЗАГРУЗЧИКА ТРЕБУЕТСЯ РАСШИРЕНИЕ (Chrome/Firefox) И ПРИЛОЖЕНИЕ (Windows/Linux/MacOS)

* #### Установка расширения (chromium):
1. Откройте браузер на основе chromium
2. Перейдите в настройки расширений через UI или по ссылке [название браузера]://extensions
3. Включите режим разработчика (переключатель в правом верхнем углу)
4. Откройте в проводнике папку с загруженным файлом "nvd-extension-chrome.crx"
5. Перетащите его в открытое ранее окно браузера

* #### Установка расширения (firefox):
1. Откройте браузер на основе firefox
2. Перейдите в настройки расширений через UI или по ссылке about:addons
3. Нажмите на шестерёнку и выберите пункт "Установить дополнение из файла..."
4. Выберите скаченный ранее файл "nvd-extension-firefox.xpi"

* #### Установка приложения:
1. Запустите файл "video-downloader-build.exe"
2. Дождитесь окончания загрузки (первый запуск идёт долго из-за установки зависимостей)
3. С этого момента приложение запустится и добавится в автозагрузку

* *В случае возникновения ошибок обновите загрузчик через настройки расширения и проверьте файл логов "app.log" или откорйте его с помощью расширения*

### 💻 Пример использования

1. Зайдите на страницу с любым видео или треком на YouTube или YouTube Music
2. Откройте меню расширения и выберите формат загружаемого файла
3. Выберите директорию загрузки в настройках (все недостающие папки будут созданы)
4. Нажмите кнопку "Скачать" и ожидайте загрузки (чем больше видео, тем дольше она будет идти)
5. Проверьте папку, которую указывали в настройках

* *Во время загрузки могут создаваться временные файлы (.part) или видео в неправильном формате, не трогайте их, они сами удалятся после окончания загрузки*

#### Third-party licenses / Лицензии третьих сторон
(ENG)
* This project uses yt-dlp (yt-dlp.exe) for video downloading. Yt-dlp is licensed under The Unlicense. Official yt-dlp website: github.com/yt-dlp/yt-dlp.
* This project uses ffmpeg (ffmpeg.exe) to convert video formats. FFmpeg is licensed under GNU Lesser General Public License (LGPL) version 2.1, but this project
uses build from https://github.com/BtbN/FFmpeg-Builds/releases/download/latest/ffmpeg-master-latest-win64-gpl.zip which has GNU General Public License (GPL).
Official ffmpeg website: www.ffmpeg.org. FFmpeg source code: https://github.com/ffmpeg/ffmpeg.

(RU)
* Этот проект использует yt-dlp (yt-dlp.exe) для загрузки видео. Yt-dlp имеет лицензию "The Unlicense". Официальный сайт yt-dlp: github.com/yt-dlp/yt-dlp.
* Этот проект использует ffmpeg (ffmpeg.exe) для конвертации форматов видео. FFmpeg имеет лицензию "GNU Lesser General Public License (LGPL) version 2.1", однако
проект использует билд из https://github.com/BtbN/FFmpeg-Builds/releases/download/latest/ffmpeg-master-latest-win64-gpl.zip, который использует GNU General Public License (GPL).
Официальный сайт ffmpeg: www.ffmpeg.org. Исходный код ffmpeg: https://github.com/ffmpeg/ffmpeg.

#### License / Лицензия
(ENG)
This project is licensed under the GPL-3.0 License - see the LICENSE file for details

(RU)
Этот проект распространяется под лицензией GPL-3.0. Подробнее смотерть в файле LICENSE.

