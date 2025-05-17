# syntax=docker/dockerfile:1

FROM golang:1.23

WORKDIR /source-code

COPY go.mod go.sum ./

RUN go mod download

COPY *.go ./
COPY schemas/*.go ./schemas/
COPY db/*.go ./db/
COPY river-data/*.go ./river-data/


# RUN CGO_ENABLED=0 GOOS=linux go build -o /forecaster.exe
RUN CGO_ENABLED=0 GOOS=linux go build -o /forecaster.exe

WORKDIR /

RUN rm -rf /app

CMD [ "/forecaster.exe" ]