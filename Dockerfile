FROM golang:1.18-alpine
RUN mkdir /app
COPY . /app
WORKDIR /app
RUN go build -o wppserver cmd/main.go
VOLUME [ "/app/dbdata" ]
ENTRYPOINT [ "/app/wppserver" ]
CMD [ "-mode=run" ]