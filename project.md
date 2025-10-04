Project Path: rtsp-viewer

Source Tree:

```txt
rtsp-viewer
├── Makefile
├── README.md
├── components
│   ├── entry.go
│   ├── logger.go
│   └── viewer.go
├── ffmpeg
├── go.mod
├── main.go
├── model
│   ├── errors.go
│   ├── model.go
│   └── types.go
├── stream
│   ├── ffmpeg.go
│   └── stream.go
└── test_server
    ├── rtsp-simple-server
    ├── rtsp-simple-server.yaml
    └── stream.mp4

```

`Makefile`:

```
server:
	cd ./test_server && ./rtsp-simple-server rtsp-simple-server.yaml

.PHONY: stream
stream:
	ffmpeg -stream_loop -1 \
	-re -i ./test_server/stream.mp4 \
	-c:v libx264 \
 	-f rtsp rtsp://test:test@localhost:8554/stream

run:
	go run main.go

dump:
	rm project.md || echo "cleared"
	code2prompt -O project.md . -e *.sum -F markdown



```

`README.md`:

```md
# Графическое приложение для подключения к IP-камере и вывода двух RTSP-потоков

# Цель
Реализовать десктопное приложение, позволяющее подключаться к IP-камере по заданным параметрам
и одновременно получать два RTSP-потока (высокое и низкое качество)
с валидацией входных данных и недопущением одновременной загрузки идентичных потоков.

# Технологии
- Язык: Golang
- Фреймворки/библиотеки: любые (на ваше усмотрение)
- Платформа: Windows (кроссплатформенность приветствуется, но не обязательна)

# Функциональные требования
- Форма подключения:
    - Поля ввода: 'IP', 'Port', 'Login', 'Password'
    - Поля ввода RTSP: 'RTSP URI #1 (High)', 'RTSP URI #2 (Low)'
    - Кнопки: 'Проверить', 'Подключиться', 'Отключиться'
    - Переключатель 'Показать пароль'

- Валидация:
    - 'IP': корректный IPv4
    - 'Port': целое 1–65535
    - 'Login'/'Password': непустые значения
    - 'RTSP URI': схема 'rtsp://', корректный URL, допускаются плейсхолдеры вида '{login}', '{password}', '{ip}', '{port}' для автоподстановки
    - Проверка уникальности потоков: если оба 'RTSP URI' указывают на один и тот же поток - уведомить пользователя и запретить подключение к обоим одновременно
        - Минимально: строковое сравнение после нормализации (обрезка пробелов, регистр, trailing slash)
        - Опционально (бонус): детект совпадения потоков при воспроизведении (сравнение SSRC/SDP/параметров SPS/PPS, разрешения/битрейта, PID/track id)

- Подключение и воспроизведение:
    - Асинхронное подключение по обеим ссылкам
    - Отображение двух независимых превью: 'High' и 'Low'
    - Отображение статуса каждого потока (Подключение / Воспроизведение / Ошибка / Отключено)
    - Корректная обработка ошибок: неверные учетные данные, таймауты, недоступный хост, 401/403, разрыв соединения
    - Возможность остановки потоков ('Отключиться')

- Пользовательские уведомления:
    - Четкие сообщения об ошибке валидации
    - Предупреждение, если потоки совпадают, с инструкцией, что исправить
    - Не логировать 'Password'; в логах допускается маскировка


# Нефункциональные требования
- UX: отзывчивый UI, отсутствие блокировки интерфейса во время сетевых операций
- Производительность: разумное использование ресурсов при одновременном проигрывании двух потоков
- Стабильность: автоматическое восстановление при кратковременных обрывах (бонус)
- Код: читаемость, разбиение на слои (UI / сервис воспроизведения / валидация / конфигурация), обработка ошибок без немых catch
- Безопасность: не хранить пароль в открытом виде на диск; опционально - шифрование настроек


# Рекомендации по стеку
- Golang:
    - UI: 'fyne', 'wails', 'gioui.org'
    - Медиа: 'pion/rtsp', 'gstreamer' bindings, интеграция с 'ffmpeg'/'vlc' (через сабпроцесс и оверлей/виджет)


# Валидация и логика подстановки
- Допускается использовать шаблоны в 'RTSP URI':
    - Пример: 'rtsp://{login}:{password}@{ip}:{port}/Streaming/Channels/101'
    - При нажатии 'Проверить' приложение выполняет подстановку и показывает итоговые URI (без раскрытия пароля в UI - маскировка '')

- Нормализация URI перед сравнением:
    - Обрезка пробелов, приведение схемы/хоста к одному регистру, удаление лишних '/'
    - Опционально: сортировка query-параметров

- Если после нормализации 'RTSP URI #1 == RTSP URI #2':
    - Показать сообщение: «Оба RTSP-URI указывают на один и тот же поток. Измените один из URI, чтобы избежать двойной загрузки.»
    - Блокировать 'Подключиться' до исправления


# UI-макет
- Верхняя панель: поля 'IP', 'Port', 'Login', 'Password' (+ чекбокс 'Показать пароль')
- Средняя панель: 'RTSP URI #1 (High)', 'RTSP URI #2 (Low)' + кнопка 'Проверить'
- Нижняя панель: два видеоконтейнера: слева 'High', справа 'Low'; под каждым - индикатор статуса
- Кнопки справа: 'Подключиться', 'Отключиться'
- Окно слева: Подсказки, логи и сообщения валидации

```
+---------------------------------------------------------------------------------------------------+
|                                   IP Camera Viewer (пример)                                       |
+---------------------------------------------------------------------------------------------------+
| IP: [...............]   Port: [.....]   Login: [........]   Password: [********]  [ ] Показать    |
+---------------------------------------------------------------------------------------------------+
| RTSP URI 1 (High): [rtsp://{login}:{password}@{ip}:{port}/Streaming/Channels/101]  [ Проверить ]  |
| RTSP URI 2 (Low) : [rtsp://{login}:{password}@{ip}:{port}/Streaming/Channels/102]                 |
+---------------------------------------------------------------------------------------------------+
|  Preview High                                   |  Preview Low                                    |
|  +-------------------------------------------+  |  +-------------------------------------------+  |
|  |                Video Area                 |  |  |                Video Area                 |  |
|  |                                           |  |  |                                           |  |
|  |                                           |  |  |                                           |  |
|  |                                           |  |  |                                           |  |
|  |                                           |  |  |                                           |  |
|  |                                           |  |  |                                           |  |
|  +-------------------------------------------+  |  +-------------------------------------------+  |
|  Статус: ... (OK/ERR)                           |   Статус: ... (OK/ERR)                          |
+---------------------------------------------------------------------------------------------------+
| Подсказки/сообщения валидации:                                            |  +-----------------+  |
| - Некорректный IP / порт / логин / пароль / RTSP-схема                    |  | [Подключиться]  |  |
| - Внимание: оба RTSP-URI совпадают после нормализации                     |  | [Отключиться]   |  |
| - Логи                                                                    |  +-----------------+  |
+---------------------------------------------------------------------------------------------------+
```

# Ошибки и сообщения
- Валидация:
    - «Некорректный IPv4-адрес»
    - «Порт должен быть числом 1–65535»
    - «Логин/пароль не должны быть пустыми»
    - «RTSP-URI должен начинаться с rtsp://»
    - «Оба RTSP-URI совпадают после нормализации»
- Подключение:
    - «Ошибка авторизации (401/403)»
    - «Таймаут подключения»
    - «Хост недоступен / сетевой сбой»
    - «Не удалось декодировать поток» (с рекомендацией по установке кодеков/библиотек)


# Тестовые данные и проверка без реальной камеры
- Можно использовать публичные RTSP тест-потоки или локальный сервер:
    - Пример шаблонов: 'rtsp://{login}:{password}@{ip}:{port}/Streaming/Channels/101' и 'rtsp://{login}:{password}@{ip}:{port}/Streaming/Channels/102'
    - Локально: использовать 'rtsp-simple-server' или 'ffmpeg' для публикации двух разных профилей (бонус)


# Критерии приемки (Definition of Done)
- Валидация всех полей работает до попытки подключения
- Предотвращается одновременная загрузка идентичных потоков (минимум по нормализованной строке)
- Два видео-превью воспроизводятся параллельно; статусы отображаются корректно
- Ошибки сети/авторизации/декодирования отображаются человеку понятными сообщениями
- UI не блокируется во время подключения/воспроизведения
- В репозитории есть 'README' с инструкциями по сборке/запуску и перечислением зависимостей
- Пароль не попадает в логи и не сохраняется в открытом виде


# Что предоставить
- Ссылка на репозиторий с исходным кодом
- 'README' (необязательно):
    - Требования к окружению
    - Команды сборки/запуска
    - Как подставлять учетные данные и проверять работу без камеры
    - Известные ограничения
- Скриншоты работающего приложения (и/или короткое видео, 30–60 сек)
- Краткое описание архитектуры (1–2 абзаца, компоненты и их ответственность)
- При использовании нативных библиотек: инструкции по установке/распаковке DLL/SO


# Критерии оценки
- Оценивание:
    - Функциональность: соответствие требованиям (40%)
    - Качество кода и архитектуры: читаемость, разделение ответственности (30%)
    - UX и стабильность: отзывчивость, обработка ошибок (20%)
    - Документация (необязательно): полнота 'README' и простота запуска (5%)
    - Бонусы (необязательно): детект одинаковых потоков по характеристикам, автореконнект (5%)


# Бонус-задачи (необязательно)
- Автоматическое определение, что 'High' и 'Low' фактически один и тот же поток (по SDP/SSRC/профилю/разрешению/битрейту)
- Автопереподключение при обрывах с экспоненциальной паузой
- Сохранение последней рабочей конфигурации (без пароля) и «Быстрый старт»
- Отображение базовой телеметрии: FPS, разрешение, битрейт


# Подсказки по реализации
- Выполняйте сетевые операции в отдельных потоках/тасках; используйте токены отмены
- Централизуйте валидацию (отдельный модуль/сервис)
- Логи храните в памяти/файле с ротацией; пароли маскируйте
- Для 'FFmpeg/LibVLC' проверьте наличие бинарных зависимостей и корректные пути загрузки

```

`components/entry.go`:

```go
package components

import "fyne.io/fyne/v2/widget"

func NewEntry(placeholder string, onChanged func(s string)) *widget.Entry {
	e := widget.NewEntry()
	e.PlaceHolder = placeholder
	e.OnChanged = onChanged

	return e
}

func NewPwdEntry(placeholder string, onChanged func(s string)) *widget.Entry {
	e := widget.NewPasswordEntry()
	e.PlaceHolder = placeholder
	e.OnChanged = onChanged

	return e
}

```

`components/logger.go`:

```go
package components

import (
	"fmt"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/widget"
)

type Logger struct {
	label *widget.Label
}

func NewLogger(label string) *Logger {
	return &Logger{label: widget.NewLabel(label)}
}

func (l *Logger) Content() fyne.CanvasObject {
	return l.label
}

func (l *Logger) Log(msgs ...any) {
	l.label.SetText(fmt.Sprintf("%v", msgs...))
}

```

`components/viewer.go`:

```go
package components

import (
	"errors"
	"image"
	"io"
	"log"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
)

type Viewer struct {
	width  int
	height int
	image  *image.RGBA
	viewer *canvas.Image
}

func NewViewer(width, height int) *Viewer {
	img := image.NewRGBA(image.Rect(0, 0, width, height))
	viewer := canvas.NewImageFromImage(img)
	viewer.SetMinSize(fyne.NewSize(float32(width), float32(height)))

	return &Viewer{
		image:  img,
		width:  width,
		height: height,
		viewer: viewer,
	}
}

func (v *Viewer) View(r io.Reader) {
	buf := make([]byte, v.width*v.height*4)
	for {
		if _, err := io.ReadFull(r, buf); err != nil {
			switch {
			case errors.Is(err, io.EOF):
				v.refresh(make([]byte, v.width*v.height*4))
			default:
				log.Println("view: ", err)
			}

			break
		}

		v.refresh(buf)
	}
}

func (v *Viewer) refresh(buf []byte) {
	copy(v.image.Pix, buf)
	fyne.Do(func() {
		v.viewer.Refresh()
	})
}

func (v *Viewer) Content() fyne.CanvasObject {
	return v.viewer
}

```

`go.mod`:

```mod
module stream-viewer

go 1.25

require fyne.io/fyne/v2 v2.6.3

require (
	fyne.io/systray v1.11.0 // indirect
	github.com/BurntSushi/toml v1.4.0 // indirect
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/fredbi/uri v1.1.0 // indirect
	github.com/fsnotify/fsnotify v1.9.0 // indirect
	github.com/fyne-io/gl-js v0.2.0 // indirect
	github.com/fyne-io/glfw-js v0.3.0 // indirect
	github.com/fyne-io/image v0.1.1 // indirect
	github.com/fyne-io/oksvg v0.1.0 // indirect
	github.com/go-gl/gl v0.0.0-20231021071112-07e5d0ea2e71 // indirect
	github.com/go-gl/glfw/v3.3/glfw v0.0.0-20240506104042-037f3cc74f2a // indirect
	github.com/go-text/render v0.2.0 // indirect
	github.com/go-text/typesetting v0.2.1 // indirect
	github.com/godbus/dbus/v5 v5.1.0 // indirect
	github.com/hack-pad/go-indexeddb v0.3.2 // indirect
	github.com/hack-pad/safejs v0.1.0 // indirect
	github.com/jeandeaual/go-locale v0.0.0-20250612000132-0ef82f21eade // indirect
	github.com/jsummers/gobmp v0.0.0-20230614200233-a9de23ed2e25 // indirect
	github.com/kr/text v0.2.0 // indirect
	github.com/nfnt/resize v0.0.0-20180221191011-83c6a9932646 // indirect
	github.com/nicksnyder/go-i18n/v2 v2.5.1 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	github.com/rymdport/portal v0.4.1 // indirect
	github.com/srwiley/oksvg v0.0.0-20221011165216-be6e8873101c // indirect
	github.com/srwiley/rasterx v0.0.0-20220730225603-2ab79fcdd4ef // indirect
	github.com/stretchr/testify v1.10.0 // indirect
	github.com/yuin/goldmark v1.7.8 // indirect
	golang.org/x/image v0.24.0 // indirect
	golang.org/x/net v0.35.0 // indirect
	golang.org/x/sys v0.30.0 // indirect
	golang.org/x/text v0.22.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)

```

`main.go`:

```go
package main

import (
	"stream-viewer/components"
	"stream-viewer/model"
	"stream-viewer/stream"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

func main() {
	// Подсказки
	logger := components.NewLogger("Сообщения:")
	m := model.New(nil, nil, logger)

	a := app.New()
	w := a.NewWindow("RTSP Stream Viewer")
	w.Resize(fyne.NewSize(1280, 0))

	// Верхние поля
	ipEntry := components.NewEntry("IP", m.ChangeIP)
	portEntry := components.NewEntry("Port", m.ChangePort)
	loginEntry := components.NewEntry("Login", m.ChangeLogin)
	passEntry := components.NewPwdEntry("Password", m.ChangePassword)

	// RTSP URI
	rtspLowEntry := components.NewEntry("RTSP URI Low", m.ChangeUriLow)
	rtspHighEntry := components.NewEntry("RTSP URI High", m.ChangeUriHigh)

	checkBtn := widget.NewButton("Проверить", m.Validate)

	// Video area
	viewerHigh := components.NewViewer(640, 480)
	viewerLow := components.NewViewer(320, 240)

	var streamHigh *stream.Stream
	var streamLow *stream.Stream

	connectBtn := widget.NewButton("Подключится", func() {
		urlHigh := "rtsp://test:test@localhost:8554/stream"
		urlLow := "rtsp://test:test@localhost:8554/stream"

		cmdHigh := stream.NewFfmpeg(urlHigh, 640, 480)
		cmdLow := stream.NewFfmpeg(urlLow, 320, 240)

		var err error
		streamHigh, err = stream.NewStream(cmdHigh)
		if err != nil {
			logger.Log(err)
			return
		}

		streamLow, err = stream.NewStream(cmdLow)
		if err != nil {
			logger.Log(err)
			return
		}

		go viewerHigh.View(streamHigh.Out)
		go viewerLow.View(streamLow.Out)

		logger.Log("Подключено: читаем поток...")

	})

	disconnectBtn := widget.NewButton("Отключится", func() {
		err := streamHigh.Close()
		if err != nil {
			logger.Log("disconnect:", err)
			return
		}
		err = streamLow.Close()
		if err != nil {
			logger.Log("disconnect:", err)
			return
		}

		logger.Log("Отключено")
	})

	// Layout: верхние поля и RTSP
	topRow := container.NewGridWithColumns(4,
		ipEntry, portEntry,
		loginEntry, passEntry,
	)

	rtspRow := container.NewGridWithColumns(2,
		container.NewGridWithRows(2,
			rtspHighEntry,
			rtspLowEntry,
		),
		checkBtn,
	)

	previewRow := container.NewGridWithColumns(2,
		viewerHigh.Content(),
		viewerLow.Content(),
	)

	btnRow := container.NewGridWithColumns(2,
		logger.Content(),
		container.NewGridWithRows(2,
			connectBtn,
			disconnectBtn,
		),
	)

	// Основной layout
	mainContainer := container.NewVBox(
		topRow,
		rtspRow,
		previewRow,
		btnRow,
	)

	w.SetContent(mainContainer)

	w.ShowAndRun()

}

```

`model/errors.go`:

```go
package model

const (
	invalidIP       = "Некорректный IPv4-адрес"
	invalidPort     = "Порт должен быть числом 1–65535"
	invalidLogin    = "Логин не должен быть пустым"
	invalidPassword = "Пароль не должен быть пустым"
	invalidURI      = "RTSP-URI должен начинаться с rtsp://"
	invalidUriPair  = "Оба RTSP-URI совпадают после нормализации"
	ok              = "Проверка параметров пройдена успешно"
)

```

`model/model.go`:

```go
package model

import (
	"net"
	"strconv"
	"strings"
)

func New(conn Connector, view Viewer, log Logger) *Model {
	return &Model{
		connector: conn,
		viewer:    view,
		logger:    log,
	}
}

type Model struct {
	ip        string
	port      string
	login     string
	password  string
	uriHigh   string
	uriLow    string
	status    Status
	connector Connector
	viewer    Viewer
	logger    Logger
}

func (m *Model) RtspHighUri() string {
	return m.uriHigh
}
func (m *Model) RtspLowUri() string {
	return m.uriLow
}

func (m *Model) ChangeIP(ip string) {
	m.ip = strings.TrimSpace(ip)
}

func (m *Model) ChangePort(port string) {
	m.port = strings.TrimSpace(port)
}

func (m *Model) ChangeLogin(login string) {
	m.login = strings.ToLower(strings.TrimSpace(login))
}

func (m *Model) ChangePassword(password string) {
	m.password = password
}

func (m *Model) ChangeUriHigh(uri string) {
	m.uriHigh = strings.ToLower(strings.TrimSpace(uri))
}

func (m *Model) ChangeUriLow(uri string) {
	m.uriLow = strings.ToLower(strings.TrimSpace(uri))
}

func (m *Model) Validate() {
	if net.ParseIP(m.ip) == nil {
		m.logger.Log(invalidIP)
		return
	}

	portInt, err := strconv.Atoi(m.port)
	if err != nil {
		m.logger.Log(invalidPort)
		return
	}

	if portInt < 1 || portInt > 65535 {
		m.logger.Log(invalidPort)
		return
	}

	if len(m.login) == 0 {
		m.logger.Log(invalidLogin)
		return
	}

	if len(m.password) == 0 {
		m.logger.Log(invalidPassword)
		return
	}

	if !strings.HasPrefix(m.uriHigh, "rtsp://") {
		m.logger.Log(invalidURI)
		return
	}

	if !strings.HasPrefix(m.uriLow, "rtsp://") {
		m.logger.Log(invalidURI)
		return
	}

	if m.uriHigh == m.uriLow {
		m.logger.Log(invalidUriPair)
		return
	}

	m.logger.Log(ok)
}

```

`model/types.go`:

```go
package model

type Status int

const (
	Disconnected Status = iota
	Connected
)

type Connector interface {
	Connect(url string) error
	Disconnect()
}

type Viewer interface {
	View()
}

type Validator interface {
	Validate() error
}

type Logger interface {
	Log(msgs ...any)
}

```

`stream/ffmpeg.go`:

```go
package stream

import (
	"fmt"
	"os/exec"
)

func NewFfmpeg(url string, width, height int) *exec.Cmd {
	return exec.Command("ffmpeg",
		"-rtsp_transport", "tcp",
		"-i", url,
		"-loglevel", "error",
		"-fflags", "nobuffer",
		"-flags", "low_delay",
		"-vf", fmt.Sprintf("scale=%d:%d", width, height),
		"-vsync", "0",
		"-f", "rawvideo",
		"-pix_fmt", "rgba",
		"pipe:1",
	)
}

```

`stream/stream.go`:

```go
package stream

import (
	"io"
	"os/exec"
)

type Stream struct {
	cmd *exec.Cmd
	Out io.Reader
	Err io.Reader
}

func NewStream(cmd *exec.Cmd) (*Stream, error) {
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return nil, err
	}

	stderr, err := cmd.StderrPipe()
	if err != nil {
		return nil, err
	}

	s := &Stream{
		cmd: cmd,
		Out: stdout,
		Err: stderr,
	}

	if err := s.cmd.Start(); err != nil {
		return nil, err
	}

	return s, nil
}

func (s *Stream) Close() error {
	return s.cmd.Process.Kill()
}

```

`test_server/rtsp-simple-server.yaml`:

```yaml
protocols: [tcp, udp]

paths:
  stream:
    readUser: test
    readPass: test

```