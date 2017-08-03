# DRAX

This is DRAX, the [DC/OS](https://dcos.io) Resilience Automated Xenodiagnosis tool. It helps to test DC/OS deployments by applying a [Chaos Monkey](http://techblog.netflix.com/2012/07/chaos-monkey-released-into-wild.html)-inspired, proactive and invasive testing approach.

Well, actually DRAX is a reverse acronym inspired by the Guardians of the Galaxy character Drax the Destroyer.

You might have heard of Netflix's [Chaos Monkey](http://techblog.netflix.com/2012/07/chaos-monkey-released-into-wild.html) or it's containerized [variant](https://medium.com/production-ready/chaos-monkey-for-fun-and-profit-87e2f343db31). Maybe you've seen a [gaming version](https://www.wehkamplabs.com/blog/2016/06/02/docker-and-zombies/) of it or stumbled upon a [lower-level species](http://probablyfine.co.uk/2016/05/30/announcing-byte-monkey/). In any case I assume you're somewhat familiar with chaos-based resilience testing.

DRAX is a DC/OS-specific resilience testing tool that works mainly on the task-level. Future work may include node-level up to cluster-level.

## Installation and usage

Note that DRAX assumes a running [DC/OS 1.9](https://dcos.io/) cluster.

### Production

Launch DRAX using the DC/OS CLI via the Marathon app spec provided:

    $ dcos marathon app add marathon-drax.json

Now you can (modulo the public node of your cluster) do the following:

    $ http http://ec2-52-38-188-110.us-west-2.compute.amazonaws.com:7777/stats
    HTTP/1.1 200 OK
    Content-Length: 10
    Content-Type: application/javascript
    Date: Mon, 13 Jun 2016 14:39:11 GMT

    {"gone":0}

If you launched DRAX via Marathon, you can also trigger a POST to the /rampage continuously by deploying a DC/OS job.  The example job is triggering the destruction every business hour from Monday till Friday: 

    $ dcos job add metronome-drax.json

### Testing and development

Get DRAX and build from source:

    $ go get github.com/dcos-labs/drax
    $ go build
    $ MARATHON_URL=http://localhost:8080 ./drax
    INFO[0000] This is DRAX in version 0.4.0                 main=init
    INFO[0000] Listening on port 7777                        main=init
    INFO[0000] On destruction level 0                        main=init
    INFO[0000] Using Marathon at  http://localhost:8080      main=init
    INFO[0000] I will destroy 2 tasks on a rampage           main=init

And in a different terminal session:

    $ http http://localhost:7777/stats
    HTTP/1.1 200 OK
    Content-Length: 10
    Content-Type: application/javascript
    Date: Mon, 13 Jun 2016 14:39:11 GMT

    {"gone":0}

For Go development, be aware of the following dependencies (not using explicit vendoring ATM):

- [github.com/gambol99/go-marathon](https://github.com/gambol99/go-marathon), an API library for working with Marathon.
- [github.com/Sirupsen/logrus](https://github.com/Sirupsen/logrus), a logging library.

### Configuration

Note that the following environment variables are pre-set in the [Marathon app spec](marathon-drax.json) and yours to overwrite.


#### Number of target tasks

To specify how many tasks DRAX is supposed to destroy in one rampage, use `NUM_TARGETS`. For example, `NUM_TARGETS=5 drax` means that (up to) 5 tasks will be destroyed, unless the overall number of tasks is less, of course.

#### Log level

To influence the log level, use the `LOG_LEVEL` env variable, for example `LOG_LEVEL=DEBUG drax` would give you fine-grained log messages (defaults to `INFO`).

## API

### /health [GET]

Will return a HTTP 200 code and `I am Groot` if DRAX is healthy.

### /stats [GET]

Will return runtime statistics, such as killed containers or apps and will report from the beginning of time (well, beginning of time for DRAX anyways).

    $ http http://localhost:7777/stats
    HTTP/1.1 200 OK
    Content-Length: 10
    Content-Type: application/javascript
    Date: Mon, 13 Jun 2016 14:39:11 GMT

    {"gone":2}

### /rampage [POST]

Will trigger a destruction. Invoke with:

    $ http POST localhost:7777/rampage
    HTTP/1.1 200 OK
    Content-Length: 121
    Content-Type: application/javascript
    Date: Mon, 13 Jun 2016 12:15:19 GMT

    {"success":true,"goners":["webserver.0fde0035-315f-11e6-aad0-1e9bbbc1653f","dummy.11a7c3bb-315f-11e6-aad0-1e9bbbc1653f"]}
