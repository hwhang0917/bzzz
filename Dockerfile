FROM node:22-alpine AS client-build
WORKDIR /app/client
COPY client/package*.json ./
RUN npm ci --only=production
COPY client/ ./
RUN npm run build

FROM golang:1.24 AS server-build
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download && go mod verify
COPY *.go ./
COPY --from=client-build /app/client/dist ./client/dist
RUN CGO_ENABLED=0 GOOS=linux go build \
    -ldflags='-w -s -extldflags "-static"' \
    -a -installsuffix cgo \
    -o main .

FROM scratch AS runtime
COPY --from=server-build /app/main /main

EXPOSE 8080

CMD ["./main"]
