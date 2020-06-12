[![Travis CI Build Status](https://travis-ci.org/xetus-oss/alertmanager-logging-receiver.svg?branch=master)](https://travis-ci.org/github/xetus-oss/alertmanager-logging-receiver)

# Alertmanager Logging Receiver

A tiny [webhook receiver](https://prometheus.io/docs/alerting/configuration/#webhook_config) for [Prometheus AlertManager](https://prometheus.io/docs/alerting/overview/) alerts that logs the received alert data in JSON format to stdout.

## Use Case

This app makes it possible -- when coupled with a log aggregator like fluentd, fluent-bit, logspout, etc... -- to use a log management tool to monitor a Prometheus environment and alert if it goes down (e.g. if a WatchDog alert doesn't fire in x minutes). While there might be other use cases, that's the only use case intended to be fulfilled by this app.

# Using This Image

Build the app, the docker image, and run a docker container locally from the app:

```sh
docker run --name am-logging-receiver -p 8080:8080 xetusoss/alertmanager-logging-receiver:latest
```

# Contributing

See the available make commands:

```sh
make help
```

## Releasing

The release process is:

1. Ensure you have a clean git state (the following should have no output):

    ```sh
    git status --porcelain
    ```

2. Ensure tests pass:

    ```sh
    make test
    ```

3. Create a tag with the desired version, push the docker images, and push the git tag

    ```sh
    git tag v${version}
    VERSION=${version} make push
    VERSION="latest" make push
    git push origin v${version}
    ```