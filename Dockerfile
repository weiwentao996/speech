FROM golang:1.19
RUN go env -w GO111MODULE=on
RUN go env -w GOPROXY=https://goproxy.cn,direct

#speechsdk start
ENV SPEECHSDK_ROOT="$HOME/speechsdk"
RUN apt-get update && apt-get install -y build-essential libssl-dev libasound2 wget \
    && mkdir -p "$SPEECHSDK_ROOT" \
    && wget -O SpeechSDK-Linux.tar.gz https://aka.ms/csspeech/linuxbinary \
    && tar --strip 1 -xzf SpeechSDK-Linux.tar.gz -C "$SPEECHSDK_ROOT" \
    && ls -l "$SPEECHSDK_ROOT" \
    && rm SpeechSDK-Linux.tar.gz \
	&& cd ~ \
	&& wget -O - https://www.openssl.org/source/openssl-1.1.1u.tar.gz | tar zxf - \
	&& cd openssl-1.1.1u \
	&& ./config --prefix=/usr/local \
	&& make -j $(nproc) \
	&& make install_sw install_ssldirs \
	&& ldconfig -v
ENV SSL_CERT_DIR="/etc/ssl/certs"
ENV CGO_CFLAGS="-I$SPEECHSDK_ROOT/include/c_api"
ENV CGO_LDFLAGS="-L$SPEECHSDK_ROOT/lib/x64 -lMicrosoft.CognitiveServices.Speech.core"
ENV LD_LIBRARY_PATH="$SPEECHSDK_ROOT/lib/x64:$LD_LIBRARY_PATH"
#speechsdk end

# Set the working directory to /app
WORKDIR /app

# Copy the current directory contents into the container at /app
COPY . /app

RUN go build -o speech main.go

CMD ["./speech"]