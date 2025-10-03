server:
	test_server/rtsp-simple-server

stream:
	ffmpeg -stream_loop -1 \
	-re -i ./test_server/stream.mp4 \
	-c:v libx264 \
 	-f rtsp rtsp://localhost:8554/stream

run:
	go run cmd/main.go

dump:
	rm project.md || echo "cleared"
	code2prompt -O project.md . -e *.sum -F markdown


