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


