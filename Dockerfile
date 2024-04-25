FROM ubuntu:latest

# Install ffmpeg
RUN apt-get update && \
    apt-get install -y ffmpeg && \
    rm -rf /var/lib/apt/lists/*

CMD ["ffmpeg", "-version"]
