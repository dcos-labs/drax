/*
Package DRAX (DC/OS Resilience Automated Xenodiagnosis) implements a
chaosmonkey-like testing functionality for DC/OS clusters.

To launch it locally use the following:

	$ MARATHON_URL=http://localhost:8080 drax

To launch it into the cluster (via the DC/OS CLI) use:

	$ dcos marathon app add marathon-drax.json

You can then use the API to destroy tasks and check stats,
see https://github.com/dcos-labs/drax#api for details.

*/
package main
