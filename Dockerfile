FROM golang:1.25-bookworm AS go-builder

RUN apt-get update && apt-get install -y \
    gcc libc6-dev git libsqlite3-dev ffmpeg \
    && rm -rf /var/lib/apt/lists/*

WORKDIR /app
COPY . .

RUN rm -f go.mod go.sum || true
RUN go mod init impossible-bot && \
RUN go get go.mau.fi/whatsmeow@latest && \
RUN go get go.mongodb.org/mongo-driver/mongo@latest && \
RUN go get go.mongodb.org/mongo-driver/bson@latest && \
RUN go get github.com/redis/go-redis/v9@latest && \
RUN go get github.com/gin-gonic/gin@latest && \
RUN go get github.com/mattn/go-sqlite3@latest && \
RUN go get github.com/lib/pq@latest && \
RUN go get github.com/gorilla/websocket@latest && \
RUN go get google.golang.org/protobuf/proto@latest && \
RUN go get github.com/showwin/speedtest-go && \
RUN go get google.golang.org/genai && \
RUN go mod tidy

RUN apt-get update && apt-get install -y git && rm -rf /var/lib/apt/lists/*


RUN apt-get update && apt-get install -y \
    ffmpeg imagemagick curl sqlite3 libsqlite3-0 \
    nodejs npm \
    atomicparsley \
    ca-certificates libgomp1 megatools libwebp-dev webp \
    libwebpmux3 libwebpdemux2 libsndfile1 \
    && rm -rf /var/lib/apt/lists/*

RUN ln -sf /usr/bin/nodejs /usr/local/bin/node

RUN curl -L https://github.com/yt-dlp/yt-dlp/releases/latest/download/yt-dlp -o /usr/local/bin/yt-dlp \
    && chmod a+rx /usr/local/bin/yt-dlp

RUN pip3 install --no-cache-dir \
    torch torchaudio --index-url https://download.pytorch.org/whl/cpu \
    && pip3 install --no-cache-dir \
    fastapi uvicorn python-multipart requests \
    faster-whisper scipy gTTS playwright librosa

RUN playwright install --with-deps chromium

WORKDIR /app

CMD ["/app/bot"]