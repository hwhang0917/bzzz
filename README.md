# Bzzz

![GitHub License](https://img.shields.io/github/license/hwhang0917/bzzz?style=flat)

> Socket based realtime buzzing application for fun

<div align="center">
<img src="./assets/demo.gif" width="400">

</div>


## Prerequisite

- Node.js (^22.14.0)
- NPM (^11.3.0)
- Go (^1.24.1)
- Docker (^28.1.0) (Optional)

- [Make](https://www.gnu.org/software/make/) (Optional)

## Build guide (For Development)

1. Clone this repository

```sh
git clone https://github.com/hwhang0917/bzzz.git
```

2. Build using [Make](https://www.gnu.org/software/make/)

```sh
make clean build
```

3. Build manually

> [!WARNING]
> Client has to be built first because server requires [embed](https://pkg.go.dev/embed)

```sh
# Build client
cd client && npm ci && npm run build

# Build server
go build -o build/server main.go
```

## Docker Guide

1. Clone this repository

```sh
git clone https://github.com/hwhang0917/bzzz.git
```

2. Build Docker Image

```sh
docker build -t bzzz .
```

3. Run Docker container

```sh
docker run -d -p 8080:8080 bzzz
```

## TODO

- [ ] Add live mouse pointer movements
- [ ] Add user custom settings
