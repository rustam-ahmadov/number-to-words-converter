FROM golang:1.22
LABEL authors="rustam"
WORKDIR /app
COPY . .
RUN go build .
EXPOSE 8080
CMD ["./m"]