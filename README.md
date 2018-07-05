# thrap


## Usage 

### Configure
```shell
$ thrap configure
```

### Initialize a new project
```shell
$ mkdir my-project
$ cd my-project
$ thrap stack init
```

### Build your project (locally)

Build your stack by running the following your project directory:

```shell
$ thrap stack build
```

### Deploy you project (locally)

Once built, deploy your project:

```shell
$ thrap stack deploy
```

### Check project status

Check the status of your stack:

```shell
$ thrap stack status
```



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
