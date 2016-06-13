# DRAX

This is DRAX, the [DC/OS](https://dcos.io) Resilience Automated Xenodiagnosis tool. It helps to test DC/OS deployments by applying a [Chaos Monkey](http://techblog.netflix.com/2012/07/chaos-monkey-released-into-wild.html)-inspired, proactive and invasive testing approach.

![DRAX logo](img/drax-logo.png)

Well, actually DRAX is a reverse acronym inspired by the Guardians of the Galaxy character Drax the Destroyer.

## Installation and usage

From source, which will get you always the latest version:

    $ go get github.com/dcos-labs/drax
    $ go build
    $ ./drax
    INFO[0000] On destruction level 0                        main=init
    This is DRAX in version 0.1.0 listening on port 7777 with default level 0

Via Marathon app spec:

    $ dcos marathon app add marathon-drax.json

### Dependencies

- [DC/OS 1.7](https://dcos.io/releases/1.7.0/)
- [github.com/gambol99/go-marathon](https://github.com/gambol99/go-marathon), an API library for working with Marathon.
- [github.com/Sirupsen/logrus](https://github.com/Sirupsen/logrus), a logging library.

### Configuration

You can influence what DRAX is supposed to destroy via the env variable `DESTRUCTION_LEVEL`: 

    0 == destroy random tasks
    1 == destroy random task of specific app
    2 == destroy random apps and services

So for example you want DRAX to totally go berserk, use this to launch it from the command line: `DESTRUCTION_LEVEL=2 drax`.

Next, you can influence how many tasks DRAX is supposed to destroy in one rampage via the env variable `NUM_TARGETS`, for example `NUM_TARGETS=5 drax` means that (up to) 5 tasks will be destroyed, unless the overall number of tasks is less, of course.

Further, in order to influence the log level, use the `LOG_LEVEL` env variable, for example `LOG_LEVEL=DEBUG drax` would give you fine-grained log messages.

## API

### /health [GET]

Will return a HTTP 200 code and `I am Groot` if DRAX is healthy.

### /stats [GET]

Will return runtime statistics, such as killed containers or apps over a report period specified with the `runs` parameter. For example, `/stats?runs=2` will report over the past two runs and if the `runs` parameter is not or wrongly specified it will report from the beginning of time (well, beginning of time for DRAX anyways).

### /rampage [POST]

Will trigger a destruction run on a certain destruction level (see also configuration section above for the default value). 

#### Target any (non-framework) app

To target any non-framework app, set the level of destruction (using the `level` parameter) to `0`, for example, `/rampage?level=0` will destroy random tasks of any apps.

To test it locally, run:

    $ MARATHON_URL=http://localhost:8080 drax

And invoke with default level (any taks on any app):

    $ http POST localhost:7777/rampage
    HTTP/1.1 200 OK
    Content-Length: 121
    Content-Type: application/javascript
    Date: Mon, 13 Jun 2016 12:15:19 GMT
    
    {"success":true,"goners":["webserver.0fde0035-315f-11e6-aad0-1e9bbbc1653f","dummy.11a7c3bb-315f-11e6-aad0-1e9bbbc1653f"]}

#### Target a specific (non-framework) app

To target a specific (non-framework) app, set the level of destruction to `1` and specify the Marathon app id using the the `app` parameter. For example, `/rampage?level=1&app=dummy` will destroy random tasks of the app with the Marathon ID `/dummy`.

To test it locally, run:

    $ MARATHON_URL=http://localhost:8080 drax

And invoke like so (to destroy tasks of app `/dummy`)

    $ http -f POST localhost:7777/rampage -- level=1 app=dummy
