# Tools

A library of tools to setup and run projects quickly.

## Features

- env: load environment variables from `.env` files without overriding existing ones
- command: run external commands with captured output, custom working directory and extra environment variables

## Project Structure

tools/
- command/
  - command.go		- Run and capture the output of external commands
  - command_test.go
- env/
  - env.go		- Load KEY=VALUE pairs from env files
  - env_test.go
- go.mod		- Go module file
- LICENSE		- Apache 2.0
- README.md		- This file

## Usage

```go
package main

import (
	"fmt"
	"log"

	"tools/command"
	"tools/env"
)

func main() {
	if err := env.Load(".env"); err != nil {
		log.Fatal(err)
	}

	version, err := command.Output("go", "version")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(version)

	if err := command.Run("make", "build"); err != nil {
		log.Fatal(err)
	}
}
```

## Development

```sh
go test ./...
```
