stream1:
	ffmpeg -stream_loop -1 -re -i ./test_server/stream.mp4 \
      -c:v libx264 -preset ultrafast -tune zerolatency \
      -b:v 1000k -maxrate 1000k -bufsize 2000k \
      -vf "format=yuv420p" \
      -c:a aac -b:a 64k \
      -f rtsp -rtsp_transport tcp \
      rtsp://test:test@89.110.116.109:8554/stream1


stream2:
	ffmpeg -stream_loop -1 -re -i ./test_server/stream.mp4 \
      -c:v libx264 -preset ultrafast -tune zerolatency \
      -b:v 1000k -maxrate 1000k -bufsize 2000k \
      -vf "format=yuv420p" \
      -c:a aac -b:a 64k \
      -f rtsp -rtsp_transport tcp \
      rtsp://test:test@89.110.116.109:8554/stream2

run:
	go run cmd/main.go

dump:
	rm project.md || echo "cleared"
	code2prompt -O project.md . -e *.sum -F markdown


.PHONY: build run deps clean

# Сборка приложения
build:
	go build -o rtsp-viewer .

# Установка зависимостей
deps:
	go mod tidy
	go mod download

# Сборка с иконками (требует fyne command)
bundle:
	fyne bundle -o bundled.go icon.png
	fyne bundle -o bundled.go -append placeholder.png

# Сборка для Windows
build-windows:
	GOOS=windows GOARCH=amd64 go build -o rtsp-viewer.exe .

# Очистка
clean:
	rm -f rtsp-viewer rtsp-viewer.exe


