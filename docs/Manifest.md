# manifest
This doc describes the thrap manifest file.

## Components
Components represent different components of the application stack. There are 4
available component types:

- ui
- api
- datastore
- serverless

### UI

### API

### Datastore
This component includes databases, caches, storage and similar constructs.  
Examples of this component type would be mysql, postgres, minio, redis, and
elasticsearch.

### Serverless
This component includes batch, periodic or event based constructs i.e. processes
that have a finite lifespan.  Examples of this component type would be a nomad
batch job, AWS lambda or even cron jobs.

### Dependencies
Dependencies are 'external' dependencies required by your stack.  These can be
third-party services such as Github or a service provided by AWS and even any  
other services within your environment in cases where one service depends on
a another one.

## Languages
Certain operations require knowing the programming language of the project.
Currently supported languages can be found [here](thrapb/thrap.go)
