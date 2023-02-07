# artblocks-stats

Small API to gather specific arblocks collection holders and their distribution by scores based on [Artacle API](https://artacle.github.io/api-docs/).

[![Go Report Card](https://goreportcard.com/badge/github.com/zd4rova/artblocks-stats)](https://goreportcard.com/report/github.com/zd4rova/artblocks-stats)
[![codecov](https://codecov.io/gh/zd4r/artblocks-stats/branch/master/graph/badge.svg?token=5KTBZW0IH6)](https://app.codecov.io/gh/zd4r/artblocks-stats)
## Starting project
Just run:
```bash
$ make compose-build-up
```
And apply migrations with `docker usage` of [migrate](https://github.com/golang-migrate/migrate):
```bash
$ docker run -v migrations:/migrations --network host migrate/migrate -path=/migrations/ -database 'postgres://user:pass@localhost:5432/holders?sslmode=disable' up
```
Swagger API specification can be found at [http://localhost:8080/api-docs/](http://localhost:8080/api-docs/) (with default service port configuration).
