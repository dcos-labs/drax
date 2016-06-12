/*
Package DRAX (DC/OS Resilience Automated Xenodiagnosis) implements a
chaosmonkey-like testing functionality for DC/OS clusters.

It provides the following functionality:

- Via the environment variable DESTRUCTION_LEVEL the destruction level
  is set, with 0 == destroy random tasks, 1 == destroy random apps, and
  2 == destroy random apps and services.
- It will expose metrics via the `/stats` endpoint.
- It will expose health status via the `/health` endpoint.
*/
package main
