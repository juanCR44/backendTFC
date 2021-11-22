FROM golang:1.10

WORKDIR /
COPY . .
RUN go get -d github.com/gorilla/mux


EXPOSE 8000
CMD ["go","run","main.go"]
