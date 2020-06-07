# Piadina

Piadina is a command line tool for querying and sorting results from [pkg.go.dev](https://pkg.go.dev). It may evolve
and include more features in the future.

## Usage

```shell script
piadina search-term
```

For instance, if you are looking for a logging package you might want to try the following command.

```shell script
piadina log
```

Results are automatically sorted by popularity, which is considered to be the number of imports.

## Unit tests

```shell script
go test ./...
```
