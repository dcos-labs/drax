# DRAX

This is DRAX, the [DC/OS](https://dcos.io) Resilience Automated Xenodiagnosis tool. Well, actually DRAX is a reverse acronym inspired by the Guardians of the Galaxy character Drax the Destroyer.

## Installation

From source, which will get you always the latest version:

    $ go get github.com/dcos-labs/drax

Via Marathon app spec:

    $ TBD

## Dependencies

- [DC/OS 1.7](https://dcos.io/releases/1.7.0/)
- [github.com/gambol99/go-marathon](https://github.com/gambol99/go-marathon), an API library for working with Marathon.
- [github.com/Sirupsen/logrus](https://github.com/Sirupsen/logrus), a logging library.


## API

### /health

Will return a HTTP 200 code and `I am Groot` if DRAX is healthy.