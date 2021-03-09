# Make Env File

Make Env File is a tool to generate an environment file (ex: .env) from another environment file (ie: .env.example).

Sometime applications have a default environment file that must be copied to launch them. It can be useful to have a dedicated tool to do that instead of creating a quick script every time it is required.

## Features

### Check the two files

Make Env File displays differences between the two files:

- variables that have different values.
- variables in source file but not in dest file (if already created).
- variables in dest file (if already created) but not in source file.

For security reasons no values are displayed in the log.

### Copy the file

Make Env File creates a `dest` env file from a `source` env file if not existing. Once created it will not be overwritten by other instance of Make Env File.

`dest` file can overwritten it if the option `overwrite` is activacted.

In every case, if files are the same it does nothing.

## How to use

### Get help

To get help, execute:

```
make-env-file -h
```

It should display:

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

Create a file if it is not existing:

```
make-env-file -source=.env.example -dest=.env
```

Create a file or overwrite if it is already existing:

```
make-env-file -source=.env.example -dest=.env -overwrite
```
