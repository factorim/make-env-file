# Make Env File

Make Env File is a tool to generate an environment file (ex: .env) from another environment file (ex: .env.example).

Sometime applications have a default environment file that must be copied to launch them. It can be useful to have a dedicated tested tool to do that instead of creating a quick script every time it is required.

## Features

### Check the two files

Make Env File displays differences between the two files:

- variables that have different values.
- variables in source file but not in dest file (if already created).
- variables in dest file (if already created) but not in source file.

For security reasons no values are displayed in the log.

If format errors are found in the files, the tool will stop and return an error.

### Copy the file

Make Env File creates a `dest` env file from a `source` env file if not existing. Once created other instance of Make Env File will do nothing.

`dest` file can overwritten it if the option `overwrite` is activacted.

In every case, if files are the same it does nothing.

## Install

Clone the repository and execute:

```
go install
```

It will install the Make Env File in `GOPATH`. Then it allows to run:

```
make-env-file -h
```

## Build binary

### Build with Docker

Execute the `make` command to generate a `make-env-file` file:

```
make build
```

### Build manually

Execute the `go buid` command to generate a `make-env-file` file:

```
go build -o make-env-file main.go
```

## How to use

### Get help

To get help, execute:

```
make-env-file -h
```

It displays:

```
Usage of make-env:
  -dest string
        config file destination. (default ".env")
  -overwrite
        overwrite if config are not equal.
  -source string
        config file source. (default ".env.example")
```

### Examples

**Check**

Create a `dest` file if it is not existing:

```
 make-env-file -source=./example/.env.example -dest=./example/.env2
```

Results:

- `.env2` is the copy of `env.example`.

**Create**

Check an already existing `dest` file:

```
 make-env-file -source=./example/.env.example -dest=./example/.env
```

Results:

- differences between `.env` and `env.example` are displayed.
- `.env` and `env.example` stay different.

**Overwrite**

Create a `dest` file or overwrite if it is already existing:

```
 make-env-file -source=./example/.env.example -dest=./example/.env -overwrite
```

Results:

- differences between `.env` and `env.example` are displayed.
- `.env` is overwritten by `env.example`.
