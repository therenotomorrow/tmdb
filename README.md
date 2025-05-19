tmdb
======

<p>
<a href="https://goreportcard.com/report/github.com/therenotomorrow/tmdb" target="_blank">
    <img src="https://goreportcard.com/badge/github.com/therenotomorrow/tmdb" alt="Report">
</a>
<a href="https://roadmap.sh/projects/task-tracker" target="_blank">
    <img src="https://img.shields.io/badge/project-tmdb_cli-blue" alt="Project">
</a>
</p>

> One of [roadmap.sh](https://roadmap.sh/projects) project. This is my small hobby.

Goal
----

`tmdb` is a project to fetch data from The Movie Database (TMDB) and display it in the terminal:

- Now Playing Movies `playing`
- Popular Movies `popular`
- Top Rated Movies `top`
- Upcoming Movies `upcoming`

System Requirements
-------------------

```shell
go version
# go version go1.24.x ...

just --version
# just 1.40.0
```

Development
-----------

Download sources

```shell
PROJECT_ROOT=tmdb
git clone https://github.com/therenotomorrow/tmdb.git "$PROJECT_ROOT"
cd "$PROJECT_ROOT"
```

Taste it :heart:

```shell
# copy environment
cp .env.example .env
vim .env

# check code integrity
just lint build

# run application
./bin/tmdb -help

# you could inject `TMDB_TOKEN` value in binary (could be unsafe)
just build 'token-value'
```

Setup safe development

```shell
git config --local core.hooksPath .githooks
```
