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

- [Make](https://www.gnu.org/software/make/) (Optional)

## Build guide

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

## TODO

- [ ] Add live mouse pointer movements
- [ ] Add user custom settings
