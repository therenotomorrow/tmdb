# tmdb

<p>
<a href="https://goreportcard.com/report/github.com/therenotomorrow/tmdb" target="_blank">
    <img src="https://goreportcard.com/badge/github.com/therenotomorrow/tmdb" alt="Go Report Card">
</a>
<a href="https://codecov.io/gh/therenotomorrow/tmdb" target="_blank">
    <img src="https://codecov.io/gh/therenotomorrow/tmdb/graph/badge.svg?token=QF1GAMYEM1" alt="Codecov">
</a>
<a href="https://github.com/therenotomorrow/tmdb/releases" target="_blank">
    <img src="https://img.shields.io/github/v/release/therenotomorrow/tmdb" alt="GitHub Releases">
</a>
<a href="https://roadmap.sh/projects/tmdb-cli" target="_blank">
    <img src="https://img.shields.io/badge/project-tmdb_cli-blue" alt="Project Link">
</a>
</p>

> One of [roadmap.sh](https://roadmap.sh/projects) project. This is my small hobby.

## Goal

`tmdb` is a command-line tool that fetches and displays movie data
from [The Movie Database (TMDB)](https://www.themoviedb.org/) right in your terminal. Available commands:

- `playing`: Now Playing Movies
- `popular`: Popular Movies
- `top`: Top Rated Movies
- `upcoming`: Upcoming Movies

## System Requirements

```shell
go version
# go version go1.24.3
just --version
# just 1.40.0
```

## Development

### Download sources

```shell
PROJECT_ROOT=tmdb
git clone https://github.com/therenotomorrow/tmdb.git "$PROJECT_ROOT"
cd "$PROJECT_ROOT"
```

### Setup dependencies

```shell
# copy environment
cp .env.example .env
vim .env

# check code integrity
just code test build
```

### Taste it :heart:

```shell
# run application
./bin/tmdb -help

# you could inject `TMDB_TOKEN` inside binary file (could be unsafe)
just build 'your-token-value'
```

### Receipts

```shell
# more just recipes available with
just
```

### Setup safe development (optional)

```shell
git config --local core.hooksPath .githooks
```
