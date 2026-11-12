---
title: Deployment
menu:
  docs:
    weight: 2
    identifier: victoriastack-deployment
    parent: victoriastack
tags:
  - metrics
  - logs
  - traces
---

You can deploy the components of VictoriaStack together or separately, in many different ways and environments:
- [Local or bare metal](#local--bare-metal) deployments take more work to compose, but give you the most flexibility.
- [Kubernetes](#kubernetes) deployments are easier to set up and run thanks to Helm charts and the operator.

Choose what suits your needs best.

## Single-node or cluster

Every database has a single-node and a clustered version, and both are Apache 2 licensed. You can mix the modes. For example, run clustered
VictoriaMetrics next to single-node VictoriaLogs.

|                   | Single-node                                                | Cluster                                                                                                                                                                                                                  |
|-------------------|------------------------------------------------------------|--------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|
| **What runs**     | one binary per database                                    | separate insert, storage, and select components                                                                                                                                                                          |
| **How to scale**  | vertically: add CPU, RAM, and disk                         | horizontally: add insert/select replicas and storage shards                                                                                                                                                              |
| **Operations**    | one process to run, upgrade, and back up                   | more components, each scaled and upgraded independently                                                                                                                                                                  |
| **Multitenancy**  | VictoriaLogs and VictoriaTraces only                       | all three databases                                                                                                                                                                                                      |
| **Docs**          | [metrics](https://docs.victoriametrics.com/victoriametrics/single-server-victoriametrics/), [logs](https://docs.victoriametrics.com/victorialogs/), [traces](https://docs.victoriametrics.com/victoriatraces/) | [metrics](https://docs.victoriametrics.com/victoriametrics/cluster-victoriametrics/), [logs](https://docs.victoriametrics.com/victorialogs/cluster/), [traces](https://docs.victoriametrics.com/victoriatraces/cluster/) |

<br>

**Start with the single-node version.** It handles workloads far larger than most teams expect, and it keeps
the setup simple. A single node scales vertically, so growing means simply using a bigger machine.
[High availability](#high-availability) can be achieved by running multiple nodes,
[replicating traffic](https://docs.victoriametrics.com/victoriastack/ingestion/#native-ingestion) to each node and load balancing
read requests between them.

Switch to the clustered version when:

- one machine can no longer keep up with ingestion or queries, or the data no longer fits its disk;
- you need [multitenancy](https://docs.victoriametrics.com/victoriastack/multitenancy/) in VictoriaMetrics, which the
  single-node version doesn't support.

The clustered versions use the same three component types in every database, named after their signal:
`vminsert`, `vmselect`, and `vmstorage` for metrics, `vlinsert`, `vlselect`, and `vlstorage` for logs,
`vtinsert`, `vtselect`, and `vtstorage` for traces.
See [Architecture](https://docs.victoriametrics.com/victoriastack/#architecture) for what each component does.

## High availability

All three databases achieve high availability the same way: run two or more independent installations and replicate
every record to each of them. An agent in front of the databases does the replication and buffers the data if
destination is temporarily unavailable, while
[vmauth](https://docs.victoriametrics.com/victoriametrics/vmauth/#load-balancing) balances read requests among the
healthy ones. There is no consensus protocol and no coordination between the installations, which keeps the scheme
simple to run and to troubleshoot.

| Signal      | HA docs                                                                                                              | What replicates the data                                                                               |
|-------------|----------------------------------------------------------------------------------------------------------------------|--------------------------------------------------------------------------------------------------------|
| **Metrics** | [VictoriaMetrics](https://docs.victoriametrics.com/victoriametrics/single-server-victoriametrics/#high-availability) | [vmagent](https://docs.victoriametrics.com/victoriametrics/vmagent/#replication-and-high-availability) |
| **Logs**    | [VictoriaLogs](https://docs.victoriametrics.com/victorialogs/#high-availability)                                     | [vlagent](https://docs.victoriametrics.com/victorialogs/vlagent/#replication-and-high-availability)    |
| **Traces**  | [VictoriaTraces](https://docs.victoriametrics.com/victoriatraces/#high-availability)                                 | [vtagent](https://docs.victoriametrics.com/victoriatraces/vtagent/#replication-and-high-availability)  |

<br>

The scheme works the same way for single-node and for clustered installations, within one availability zone or across
several.

Replication protects against losing a node or a zone. It does not protect against a bad config or an accidental deletion,
since those reach every replica. Keep [backups](#backups) as well.

On Kubernetes, the operator can manage this whole topology as a single resource.
See [high availability in Kubernetes](#high-availability-in-kubernetes).

## Local / Bare Metal

Every VictoriaStack component has extensive Quick Start instructions for running in various environments:
- [VictoriaMetrics Quick Start](https://docs.victoriametrics.com/victoriametrics/quick-start/) (metrics)
- [VictoriaLogs Quick Start](https://docs.victoriametrics.com/victorialogs/quick-start/) (logs)
- [VictoriaTraces Quick Start](https://docs.victoriametrics.com/victoriatraces/quick-start/) (traces)
- [vmauth Quick Start](https://docs.victoriametrics.com/victoriametrics/vmauth/) (authentication, routing, load balancing)
- [vmalert Quick Start](https://docs.victoriametrics.com/victoriametrics/vmalert/#quickstart) (alerting)

### Binary

Download the binaries from the GitHub releases pages:
- [VictoriaMetrics releases](https://github.com/VictoriaMetrics/VictoriaMetrics/releases)
- [VictoriaLogs releases](https://github.com/VictoriaMetrics/VictoriaLogs/releases)
- [VictoriaTraces releases](https://github.com/VictoriaMetrics/VictoriaTraces/releases)

The releases pages also contain the `vmutils`, `vlutils` and `vtutils` (for metrics, logs and traces) packages with components like
[vmauth](https://docs.victoriametrics.com/victoriametrics/vmauth/), [vmalert](https://docs.victoriametrics.com/victoriametrics/vmalert/),
[vmagent](https://docs.victoriametrics.com/victoriametrics/vmagent/), [vlagent](https://docs.victoriametrics.com/victorialogs/vlagent/), [vtagent](https://docs.victoriametrics.com/victoriatraces/vtagent/), etc.

Choose the archive that matches the operating system and CPU architecture of your machine.
See detailed instructions on downloading and running binaries:
- [VictoriaMetrics: try it locally](https://docs.victoriametrics.com/victoriametrics/quick-start/#try-it-locally)
- [VictoriaLogs: try it locally](https://docs.victoriametrics.com/victorialogs/quickstart/#try-it-locally)
- [VictoriaTraces: try it locally](https://docs.victoriametrics.com/victoriatraces/quick-start/#starting-victoriatraces-single-node-from-a-binary)

### Docker

VictoriaStack components Docker images are published on [Docker Hub](https://hub.docker.com/u/victoriametrics) and 
[Quay](https://quay.io/organization/victoriametrics).

Every database keeps its data in a directory inside the container, so mount a volume for it. Without a volume the data
stays in the container's writable layer and disappears together with the container:

```sh
docker run -it --rm -p 8428:8428 -v ./victoria-metrics-data:/victoria-metrics-data \
  victoriametrics/victoria-metrics:latest

docker run -it --rm -p 9428:9428 -v ./victoria-logs-data:/victoria-logs-data \
  victoriametrics/victoria-logs:latest

docker run -it --rm -p 10428:10428 -v ./victoria-traces-data:/victoria-traces-data \
  victoriametrics/victoria-traces:latest
```

Each command runs a single-node database and publishes its HTTP port: 8428 for metrics, 9428 for logs, 10428 for traces.
Both the API and the built-in vmui are served on that port. Pin an explicit image tag instead of `latest`, so that
pulling an image never turns into an unplanned [upgrade](#upgrades).

### Docker Compose

The [QuickStart page](https://docs.victoriametrics.com/victoriastack/quick-start/) has Docker Compose examples for running the
whole stack, in single-node and clustered versions.

## Kubernetes

[VictoriaMetrics operator](#victoriametrics-operator) runs every component of the stack through Kubernetes custom resources,
and the [k8s-stack chart](#helm-charts) installs the operator together with those resources. 

Pick the entry point that matches how you manage Kubernetes:

- **[k8s-stack chart](#helm-charts)**: the fastest way to install monitoring stack for Kubernetes. It brings the operator, the
  databases, agents, alerting rules, Grafana, and the scrape configs for the cluster itself. Every component is still
  configured through its custom resource, under the `spec` key in the chart values.
- **[Operator](#victoriametrics-operator)**: is a more generic solution for building a platform with most flexibility using custom resources. 
  Keep the manifests in Git and apply them with ArgoCD, Flux, or similar tooling. This option fits installations that need a
  layout the chart doesn't cover.
- **Single-component charts**: use [one of the other charts](https://docs.victoriametrics.com/helm/), such as
  `victoria-metrics-single` or `victoria-logs-single`, when you need one component rather than the whole stack.

### Helm Charts

With [VictoriaMetrics K8s Stack](https://docs.victoriametrics.com/helm/victoria-metrics-k8s-stack/) chart you can install
and upgrade a complete VictoriaStack in Kubernetes via [Helm](https://helm.sh/). Alongside the VictoriaStack components
the chart installs the operator, Grafana, kube-state-metrics, and node-exporter as dependency charts. Each of them can be
disabled in the values, for example to reuse a Grafana or an operator you already run.

By default, the chart installs everything needed to collect metrics, evaluate alerting rules, and visualize the data in
Grafana. Logs and traces storage is disabled by default. Enable it in the chart values:

```yaml
# single-node deployments; use vlcluster and vtcluster for the clustered versions
vlsingle:
  enabled: true
vtsingle:
  enabled: true
```

The chart settings configure every component separately, so you can install only what you need.

See the full list of [available charts](https://docs.victoriametrics.com/helm/) for the rest of the components.


### VictoriaMetrics Operator

[VictoriaMetrics Kubernetes Operator](https://docs.victoriametrics.com/operator/) manages VictoriaStack components in
Kubernetes and OpenShift clusters via custom resources. Describe the desired state via
[resources](https://docs.victoriametrics.com/operator/resources/) such as `VMCluster`, `VLSingle`, or `VMAgent`, and the
operator will create and update the installation behind it.

For metrics, the operator picks up scrape configuration on its own: it discovers `VMServiceScrape` and similar resources, and
[converts existing prometheus-operator objects](https://docs.victoriametrics.com/operator/integrations/prometheus/#objects-conversion)
such as `ServiceMonitor` into their VictoriaMetrics equivalents. Existing Prometheus setups keep working after the
migration.

See [Quick Start instructions](https://docs.victoriametrics.com/operator/quick-start/) for the operator.

### High availability in Kubernetes

The operator describes a multi-zone installation as a single resource:
[VMDistributed](https://docs.victoriametrics.com/operator/resources/vmdistributed/) for metrics
and [VLDistributed](https://docs.victoriametrics.com/operator/resources/vldistributed/) for logs. Each of them takes a
list of zones and creates, for every zone, the database plus an agent that replicates writes into it. All zones use the
same backend type: either the clustered or the single-node version of the database. A shared vmauth in front of the
zones acts as the global read and write entrypoint, so producers and Grafana keep using one address no matter how many
zones exist.

The operator applies changes to the zones sequentially, and each zone can be switched between serving reads, writes,
both, or nothing at all. An upgrade or a maintenance window therefore takes one zone out of rotation at a time, while
the remaining zones keep ingesting and answering queries.

The [k8s-stack chart](https://docs.victoriametrics.com/helm/victoria-metrics-k8s-stack/) creates these resources from
its values:

```yaml
# metrics: mutually exclusive with vmsingle and vmcluster
vmdistributed:
  enabled: true

# logs: mutually exclusive with vlsingle and vlcluster; needs chart v0.90.0 or newer
vldistributed:
  enabled: true
```

The operator also documents high availability
for [every individual resource](https://docs.victoriametrics.com/operator/resources/#high-availability) it manages,
including vmagent, vlagent, vmauth, vmalert, and Alertmanager.

## Security

All three databases expect to run in a trusted network, unreachable from the internet. They perform no authentication
of their own, so put [vmauth](https://docs.victoriametrics.com/victoriametrics/vmauth/) or another authorization proxy
in front of them. It authenticates requests, terminates TLS, and routes them to the right component.
See [Access control](https://docs.victoriametrics.com/victoriastack/access-control/) for the configuration.

Two more things to set up for every component:

- **Protect the internal HTTP endpoints.** Stack components expose sensitive endpoints for config hot reloading,
  metrics scraping, profiling, etc. See more about [Sensitive endpoints](https://docs.victoriametrics.com/victoriastack/access-control/#sensitive-endpoints).
- **Enable TLS on components that serve traffic without a proxy in front** with the `-tls`, `-tlsCertFile`, and
  `-tlsKeyFile` flags. vmauth additionally supports mTLS and automatic issuing of TLS certificates.

Security fixes are published in the latest release of every database, so staying current is part of keeping an
installation secure. See [Upgrades](#upgrades).

## Upgrades

Upgrades and downgrades are safe unless the release notes say otherwise. Skipping versions is safe too. Check the
changelog before upgrading: [metrics](https://docs.victoriametrics.com/victoriametrics/changelog/),
[logs](https://docs.victoriametrics.com/victorialogs/changelog/),
[traces](https://docs.victoriametrics.com/victoriatraces/changelog/). The databases are released separately, so
upgrade each of them on its own schedule.

Upgrading a single-node database takes three steps:

1. Send `SIGINT` to the process to stop it gracefully.
2. Wait until the process exits. This can take a few seconds.
3. Start the new version with the same data directory and flags.

Most agents keep buffering data while a database restarts, so nothing is lost during the gap.
See the detailed docs on upgrading
[VictoriaMetrics](https://docs.victoriametrics.com/victoriametrics/single-server-victoriametrics/#how-to-upgrade-victoriametrics),
[VictoriaLogs](https://docs.victoriametrics.com/victorialogs/#upgrading), [VictoriaTraces](https://docs.victoriametrics.com/victoriatraces/#upgrading).

In Kubernetes, the operator manages the restarts. Pin the version of every component
via [`spec.image.tag`](https://docs.victoriametrics.com/operator/resources/#managing-versions-of-vm), and the operator
rolls out the change. Multi-zone installations managed
by [VMDistributed or VLDistributed](#high-availability-in-kubernetes) upgrade one zone at a time, so the remaining
zones keep serving traffic.

Teams that prefer rare, low-risk updates can follow the long-term support release lines instead of the latest release.
LTS releases are published for VictoriaMetrics and VictoriaLogs and are part of the Enterprise package,
see [Long-term support releases](https://docs.victoriametrics.com/victoriastack/enterprise/#long-term-support-releases).

## Backups

Backups in VictoriaStack are based on data snapshots. A snapshot is created instantly within a second and puts no load on the
database when doing so. Upload the snapshots to durable storage such as S3 to keep them as backups:
- [VictoriaMetrics backups management: vmbackup, vmrestore, vmbackupmanager](https://docs.victoriametrics.com/victoriametrics/vmbackup/)
- [VictoriaLogs backup and restore](https://docs.victoriametrics.com/victorialogs/#backup-and-restore)
- [VictoriaTraces backup and restore](https://docs.victoriametrics.com/victoriatraces/#backup-and-restore)

Besides the snapshots API, VictoriaMetrics has more advanced tools for backup management:
[vmbackup](https://docs.victoriametrics.com/victoriametrics/vmbackup/) creates backups and
[vmrestore](https://docs.victoriametrics.com/victoriametrics/vmrestore/) restores them.
[vmbackupmanager](https://docs.victoriametrics.com/victoriametrics/vmbackupmanager/) adds scheduling and backup
retention, so backups need no additional tooling. It is a part of
the [Enterprise package](https://docs.victoriametrics.com/victoriastack/enterprise/), while vmbackup and vmrestore
are Apache 2 licensed. All three tools are incorporated into Helm charts and the operator.

The same tooling is not available for VictoriaLogs and VictoriaTraces yet. Follow the progress
in [issue #123](https://github.com/VictoriaMetrics/VictoriaLogs/issues/123).


## Monitoring

Every component of VictoriaStack was built with monitoring in mind. It exposes its own metrics in Prometheus exposition
format and logs every significant event, so its internal processes stay transparent.

VictoriaStack components ship with production-ready Grafana dashboards, available at
[grafana.com/victoriametrics/dashboards](https://grafana.com/orgs/victoriametrics/dashboards). The VictoriaMetrics team
maintains and updates them.

Alerting rules come with the components too. Here are the rules recommended for running with the stack:
- [VictoriaMetrics rules](https://github.com/VictoriaMetrics/VictoriaMetrics/tree/master/deployment/docker/rules)
- [VictoriaLogs rules](https://github.com/VictoriaMetrics/VictoriaLogs/tree/master/deployment/docker/rules)
- [VictoriaTraces rules](https://github.com/VictoriaMetrics/VictoriaTraces/tree/master/deployment/docker/rules)

The rules and the dashboards are closely interconnected: an alert points at the panel that explains it.

The VictoriaMetrics team runs the same stack internally, in VictoriaMetrics Cloud, and for Enterprise customers, so the
rules and the dashboards keep up with the components they monitor.
