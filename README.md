# Admin tool
The recursive finder for files with non-valid names

## BUILD 
Requires [Go](https://golang.org/doc/install). Tested with Go 1.15.

Clone this repo locally and run test, build:
```
mkdir -p $HOME/finder-invalid-filesnames && \
cd $HOME/singlebackuper && \
git clone https://github.com/denfm/finder-invalid-filesnames  ./ && \
make test && make build && \
cd bin && ls -la
```

## Running

```
./bin/finder --parse-path=/var/www/legacy-project/uploads --config-path=/tmp/finder-legacy.csv
```

## LICENSE

See [LICENSE](./LICENSE)