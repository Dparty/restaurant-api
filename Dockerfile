FROM golang:1.21 as build-stage
WORKDIR /app
COPY go.mod go.sum ./
RUN go env -w GOPRIVATE=github.com/Dparty/*
RUN git config --global url."https://chenyunda218:${TOKEN}@github.com".insteadOf "https://github.com"
RUN go mod download
COPY . .
RUN go build -o /main
EXPOSE 8080
CMD [ "/main" ]