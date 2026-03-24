# Lab 1

## Requirements

- [Go](https://go.dev/dl/) 1.26+
- [golangci-lint](https://golangci-lint.run/) 2.x

## Run the application

```bash
make build
./bin/app
```

Or without building:

```bash
go run ./cmd/app/
```

## Run tests

```bash
make test
```

## Run linter

```bash
make lint
```

## Run everything

```bash
make all
```

Runs `fmt`, `lint`, `test`, and `build` sequentially.
