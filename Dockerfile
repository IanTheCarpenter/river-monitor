FROM golang:1.23

WORKDIR /source-code

COPY go.mod go.sum ./
COPY db/go.mod go.sum ./

RUN go mod download

COPY *.go ./
COPY schemas/*.go ./
COPY db/*.go ./
COPY schemas/*.go ./


# RUN CGO_ENABLED=0 GOOS=linux go build -o /forecaster.exe
RUN go build -o /forecaster.exe

WORKDIR /

RUN rm -rf /app

CMD [ "/forecaster.exe" ]