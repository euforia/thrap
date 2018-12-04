# thrap

Thrap is a tool accelerate the software development process by automating common tasks and
simpifying various integration points

## Getting Started

### Pre-requisites

The following are required for thrap:

- Docker (api version >= 1.37)

### Installation

Download the appropriate binary based on your platform from [here](https://github.com/euforia/thrap/releases)
and copy it into your path:

```shell
$ mv thrap /usr/local/bin/
```

### Configuration

The commnand below will perform the initial configuration asking questions as needed.

```shell
$ thrap configure
```

You are now ready to use thrap.

### Initialize a new project

Start by initializing a new project:

```shell
$ mkdir my-project
$ cd my-project
$ thrap stack init
```

This will create the initial set of base files and configurations.

### Build your project (locally)

Once the project is initialized, you can make code changes as needed.  When ready, the project stack can
be built using the following command:

```shell
$ thrap stack -p <profile> build
```

This starts all necessary services, builds all containers, exiting after all head containers have 
completed building.

### Deploy your project (locally)

Once built, deploy your project:

```shell
$ thrap stack -p <profile> deploy
...
```

### Check project status

Check the status of your stack:

```shell
$ thrap stack status
...
```

## Development

#### Install dependencies

```shell
$ make deps
...
```

#### Run tests

```shell
$ make test
...
```

#### Make binary

```shell
$ make thrap
...
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
