# thrap


## Development

#### Install dependencies
```shell
$ make deps
```

#### Run tests
```shell
$ make test
```

#### Make binary
```shell
$ make thrap
```

Binary called `thrap` (built to be compatible with the system it was run on)
will be available in this folder

#### Make distribution
```shell
$ make dist
```

Binaries will be available in the `dist` folder.

#### Docker
A fully containerized build can be run as follows:
```shell
docker build -t thrap -f < /path/to/dockerfile > .
```
