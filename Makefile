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


