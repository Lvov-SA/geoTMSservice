FROM golang:1.24-alpine
ENV CGO_ENABLED=1
RUN echo "https://dl-cdn.alpinelinux.org/alpine/edge/main" >> /etc/apk/repositories && \
    echo "https://dl-cdn.alpinelinux.org/alpine/edge/community" >> /etc/apk/repositories && \
    echo "https://dl-cdn.alpinelinux.org/alpine/edge/testing" >> /etc/apk/repositories && \
    apk update && \
    apk add --no-cache gdal gdal-dev g++ pkgconfig
RUN apk add --no-cache \
    gdal-tools \
    proj \
    geos \
    hdf5 \
    gdal-driver-png \
    gdal-driver-gif \
    gdal-driver-jpeg \
    netcdf
RUN apk add sqlite
WORKDIR /app
COPY go.mod ./ 
RUN go mod tidy
RUN go mod download
COPY . .

RUN cd cmd && go build -o app main.go
