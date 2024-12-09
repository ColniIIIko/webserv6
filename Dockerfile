FROM golang:1.22.4-bullseye
 
WORKDIR /app
 
COPY . .

COPY go.mod .
 
RUN go mod download

COPY main.go .
 
EXPOSE 8080
 
CMD [ "go", "run", "main.go" ]