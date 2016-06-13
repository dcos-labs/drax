# DRAX

This is DRAX, the [DC/OS](https://dcos.io) Resilience Automated Xenodiagnosis tool. Well, actually DRAX is a reverse acronym inspired by the Guardians of the Galaxy character Drax the Destroyer.

## Installation and usage

From source, which will get you always the latest version:

    $ go get github.com/dcos-labs/drax
    $ go build
    $ ./drax
    INFO[0000] On destruction level 0                        main=init
    This is DRAX in version 0.1.0 listening on port 7777 with default level 0

Via Marathon app spec:

    $ TBD

### Dependencies

- [DC/OS 1.7](https://dcos.io/releases/1.7.0/)
- [github.com/gambol99/go-marathon](https://github.com/gambol99/go-marathon), an API library for working with Marathon.
- [github.com/Sirupsen/logrus](https://github.com/Sirupsen/logrus), a logging library.

### Configuration

You can influence what DRAX is supposed to destroy via the env variable `DESTRUCTION_LEVEL`: 

    0 == destroy random tasks
    1 == destroy random apps
    2 == destroy random apps and services

So for example you want DRAX to totally go berserk, use this to launch it from the command line: `DESTRUCTION_LEVEL=2 drax`.

## API

### /health [GET]

Will return a HTTP 200 code and `I am Groot` if DRAX is healthy.

### /stats [GET]

Will return runtime statistics, such as killed containers or apps over a report period specified with the `runs` parameter. For example, `/stats?runs=2` will report over the past two runs and if the `runs` parameter is not or wrongly specified it will report from the beginning of time (well, beginning of time for DRAX anyways).

### /rampage [POST]

Will trigger a destruction run on the current destruction level (see configuration section, above). You can explicitly set the level of destruction using the `level` parameter, for example, `/rampage?level=1` will destroy random apps (but no services/frameworks).