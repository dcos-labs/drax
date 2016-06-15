# DRAX

This is DRAX, the [DC/OS](https://dcos.io) Resilience Automated Xenodiagnosis tool. It helps to test DC/OS deployments by applying a [Chaos Monkey](http://techblog.netflix.com/2012/07/chaos-monkey-released-into-wild.html)-inspired, proactive and invasive testing approach.

![DRAX logo](img/drax-logo.png)

Well, actually DRAX is a reverse acronym inspired by the Guardians of the Galaxy character Drax the Destroyer.


You might have heard of Netflix's [Chaos Monkey](http://techblog.netflix.com/2012/07/chaos-monkey-released-into-wild.html) or it's containerized [variant](https://medium.com/production-ready/chaos-monkey-for-fun-and-profit-87e2f343db31). Maybe you've seen a [gaming version](https://www.wehkamplabs.com/blog/2016/06/02/docker-and-zombies/) of it or stumbled upon a [lower-level species](http://probablyfine.co.uk/2016/05/30/announcing-byte-monkey/). In any case I assume you're somewhat familiar with chaos-based resilience testing.

DRAX is a DC/OS-specific resilience testing tool that works mainly on the task-level. Future work may include node-level up to cluster-level.

## Installation and usage

Note that DRAX assumes a running [DC/OS 1.7](https://dcos.io/releases/1.7.0/) cluster.

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

### Testing and development

Get DRAX and build from source:

    $ go get github.com/dcos-labs/drax
    $ go build
    $ MARATHON_URL=http://localhost:8080 ./drax
    INFO[0000] Using Marathon at  http://localhost:8080      main=init
    INFO[0000] On destruction level 0                        main=init
    INFO[0000] I will destroy 2 tasks on a rampage           main=init
    INFO[0000] This is DRAX in version 0.3.0 listening on port 7777

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

#### Destruction level

You can influence the default destruction setting for DRAX via the env variable `DESTRUCTION_LEVEL`: 

    0 == destroy random tasks of any app
    1 == destroy random tasks of specific app

#### Number of target tasks

To specify how many tasks DRAX is supposed to destroy in one rampage, use `NUM_TARGETS`. For example, `NUM_TARGETS=5 drax` means that (up to) 5 tasks will be destroyed, unless the overall number of tasks is less, of course.

#### Log level

To influence the log level, use the `LOG_LEVEL` env variable, for example `LOG_LEVEL=DEBUG drax` would give you fine-grained log messages (defaults to `INFO`).

### Roadmap

- add seeds (hello world dummy, NGINX, Marvin): shell script + DC/OS CLI and walkthrough examples
- Weave [Scope](https://www.weave.works/products/weave-scope/) demo
- tests, tutorial, blog post
- node/cluster level rampages

## API

### /health [GET]

Will return a HTTP 200 code and `I am Groot` if DRAX is healthy.

### /stats [GET]

Will return runtime statistics, such as killed containers or apps over a report period specified with the `runs` parameter. For example, `/stats?runs=2` will report over the past two runs and if the `runs` parameter is not or wrongly specified it will report from the beginning of time (well, beginning of time for DRAX anyways).

    $ http http://localhost:7777/stats
    HTTP/1.1 200 OK
    Content-Length: 10
    Content-Type: application/javascript
    Date: Mon, 13 Jun 2016 14:39:11 GMT
    
    {"gone":2}

### /rampage [POST]

Will trigger a destruction run on a certain destruction level (see also configuration section above for the default value). 

#### Target any (non-framework) app

To target any non-framework app, set the level of destruction (using the `level` parameter) to `0`, for example, `/rampage?level=0` will destroy random tasks of any apps.

Invoke with default level (any tasks in any app):

    $ http POST localhost:7777/rampage
    HTTP/1.1 200 OK
    Content-Length: 121
    Content-Type: application/javascript
    Date: Mon, 13 Jun 2016 12:15:19 GMT
    
    {"success":true,"goners":["webserver.0fde0035-315f-11e6-aad0-1e9bbbc1653f","dummy.11a7c3bb-315f-11e6-aad0-1e9bbbc1653f"]}

#### Target a specific (non-framework) app

To target a specific (non-framework) app, set the level of destruction to `1` and specify the Marathon app id using the the `app` parameter. For example, `/rampage?level=1&app=dummy` will destroy random tasks of the app with the Marathon ID `/dummy`.

Invoke like so (to destroy tasks of app `/dummy`):

    $ cat rp.json
    {
     "level" : "1",
     "app" : "dummy"
    }
    $ http POST localhost:7777/rampage < rp.json
    HTTP/1.1 200 OK
    Content-Length: 117
    Content-Type: application/javascript
    Date: Mon, 13 Jun 2016 13:05:31 GMT
    
    {"success":true,"goners":["dummy.59dca877-3165-11e6-aad0-1e9bbbc1653f","dummy.e96ffce3-3164-11e6-aad0-1e9bbbc1653f"]}

