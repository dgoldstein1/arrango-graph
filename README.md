# Destributed Graph

A highly-destributed graph using (arrango db graph)[https://www.arangodb.com/docs/stable/graphs.html] for directed data.

- circle ci
- maintaibaility
- test coverage

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
