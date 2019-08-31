# Destributed Graph

A highly-destributed graph using (arrango db graph)[https://www.arangodb.com/docs/stable/graphs.html] for directed data.

[![CircleCI](https://circleci.com/gh/dgoldstein1/destrib-graph.svg?style=svg)](https://circleci.com/gh/dgoldstein1/destrib-graph)
[![Maintainability](https://api.codeclimate.com/v1/badges/3ef17277612516e345de/maintainability)](https://codeclimate.com/github/dgoldstein1/destrib-graph/maintainability)
[![Test Coverage](https://api.codeclimate.com/v1/badges/3ef17277612516e345de/test_coverage)](https://codeclimate.com/github/dgoldstein1/destrib-graph/test_coverage)

## Install

```sh
go install github.com/dgoldstein1/destrib-graph
```

or

```sh
docker pull dgoldstein1/destrib-graph:latest
```


## Run it

```sh
export GRAPH_DB_STORE_PORT="5001" # port served on
export GRAPH_DOCS_DIR="./api/*" # location of docs (warning: this entire dir is served up to the browser)
./destrib-graph server
```


## Development

#### Local Development

- Install [inotifywait](https://linux.die.net/man/1/inotifywait)
```sh
./watch_dev_changes.sh
```

#### Testing

```sh
go test $(go list ./... | grep -v /vendor/)
```

## Generating New Documentation

```sh
pip install PyYAML
python api/swagger-yaml-to-html.py < api/swagger.yml > api/index.html
```


## Authors

* **David Goldstein** - [DavidCharlesGoldstein.com](http://www.davidcharlesgoldstein.com/?github-destrib-graph) - [Decipher Technology Studios](http://deciphernow.com/)

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details
