---
title: Quick Start
menu:
  docs:
    weight: 1
    identifier: victoriastack-quick-start
    parent: victoriastack
tags:
  - metrics
  - logs
  - traces
---

> If you're looking for a visual impression of the VictoriaStack,
> visit [Playgrounds page](https://docs.victoriametrics.com/playgrounds/)
> and try publicly running installations of VictoriaMetrics, VictoriaLogs and VictoriaTraces fed
> with [OpenTelemetry Astronomy Shop demo](https://github.com/open-telemetry/opentelemetry-demo) data.

You can deploy VictoriaStack as binaries on bare metal, as docker containers, or in Kubernetes via Helm charts and the
operator. See the [deployment docs](https://docs.victoriametrics.com/victoriastack/deployment/) for details.

The sections below deploy VictoriaStack with Docker Compose, the easiest way to run and experiment with it locally.
You will need [Docker](https://docs.docker.com/get-started/get-docker/)
with [Compose v2](https://docs.docker.com/compose/) and `git`.

## Single-node version

The simplest setup runs VictoriaMetrics, VictoriaLogs, and VictoriaTraces as single-node deployments.
Start it locally with the [otel-single-node.yml](https://github.com/VictoriaMetrics/VictoriaStack/blob/main/deployment/otel-single-node.yml)
Docker Compose file:

```sh
git clone https://github.com/VictoriaMetrics/VictoriaStack.git && cd VictoriaStack
make otel-single-up
```

> Note: this Docker Compose environment is simplified for demonstration purposes.
> It is **not intended for production use**.

The command starts:

- a minimized version of [OpenTelemetry Demo](https://github.com/VictoriaMetrics/VictoriaStack/blob/main/deployment/otel-demo-mini.yml) with an OpenTelemetry collector;
- single-node VictoriaMetrics, VictoriaLogs, and VictoriaTraces for storing and querying signals;
- vmalert for evaluating alerting and recording rules;
- Grafana for visualizing signals;
- vmauth as a single entrypoint that authorizes and routes requests to the appropriate component.

![Diagram showing how the demo works. The OTel Demo sends metrics, logs, and traces to the OTel Collector, which connects to vmauth. vmauth then routes the telemetry to the three databases. Grafana and vmalert also connect to the databases via vmauth.](quick-start-compose-single.webp)


The OpenTelemetry demo will keep pushing telemetry signals to the stack through the OpenTelemetry collector, and vmalert
will keep evaluating rules. To stop the stack, run:

```sh
make otel-single-down
```

> The single-node version is the simplest way to start experimenting and to learn how the components fit together.
To scale the stack horizontally, see the [clustered version](#clustered-version).

After the demo has been running for a few minutes, the links in the chapters below become available. 
If demo is failing, check how to [get help](https://docs.victoriametrics.com/victoriastack/#get-help).

### VictoriaMetrics UI

Visit [VictoriaMetrics UI - http://localhost:8427/metrics/select/vmui](http://localhost:8427/metrics/select/vmui/?#/cardinality)
to explore the stored metrics or [run a query to plot charts](http://localhost:8427/metrics/select/vmui/#/?g0.range_input=15m&g0.relative_time=last_15_minutes&g0.tab=0&g0.expr=sum%28rate%28%7B__name__%3D%22http.server.request.body.size_bucket%22%7D%5B5m%5D%29%29&g0.step_input=3s):
![Screenshot of VMUI for VictoriaMetrics.](quick-start-metrics-vmui-demo.webp)
{width="500"}

### VictoriaLogs UI

Visit [VictoriaLogs UI - http://localhost:8427/logs/select/vmui](http://localhost:8427/logs/select/vmui) to see
recently ingested logs or [run a query to plot charts](http://localhost:8427/logs/select/vmui/?#/?g0.range_input=15m&g0.relative_time=last_15_minutes&step=5s&query=*+%7C+by%28service.name%29+count%28%29&group_hits=service.name):
![Screenshot of VMUI for VictoriaLogs.](quick-start-logs-vmui-demo.webp)
{width="500"}

### VictoriaTraces UI

Visit [VictoriaTraces UI - http://localhost:8427/traces/select/vmui](http://localhost:8427/traces/select/vmui) to see recent traces:
![Screenshot of VMUI for VictoriaTraces.](quick-start-traces-vmui-demo.webp) 

Or explore traces in Grafana's Explore [via Jaeger](http://localhost:3000/explore?schemaVersion=1&panes=%7B%220ut%22:%7B%22datasource%22:%22victoriatraces-jaeger-ds%22,%22queries%22:%5B%7B%22refId%22:%22A%22,%22datasource%22:%7B%22type%22:%22jaeger%22,%22uid%22:%22victoriatraces-jaeger-ds%22%7D,%22queryType%22:%22search%22,%22service%22:%22shipping%22%7D%5D,%22range%22:%7B%22from%22:%22now-1h%22,%22to%22:%22now%22%7D,%22compact%22:false%7D%7D&orgId=1)
or [Tempo](http://localhost:3000/explore?schemaVersion=1&panes=%7B%220ut%22:%7B%22datasource%22:%22victoriatraces-tempo-ds%22,%22queries%22:%5B%7B%22refId%22:%22A%22,%22datasource%22:%7B%22type%22:%22tempo%22,%22uid%22:%22victoriatraces-tempo-ds%22%7D,%22queryType%22:%22traceqlSearch%22,%22limit%22:20,%22tableType%22:%22traces%22,%22metricsQueryType%22:%22range%22,%22serviceMapUseNativeHistograms%22:false,%22filters%22:%5B%7B%22id%22:%2219638448%22,%22operator%22:%22%3D%22,%22scope%22:%22span%22%7D%5D%7D%5D,%22range%22:%7B%22from%22:%22now-1h%22,%22to%22:%22now%22%7D,%22compact%22:false%7D%7D&orgId=1)
(_login with admin:admin_):
![Screenshot of Grafana's Explore page showing VictoriaTraces trace data.](quick-start-traces-grafana-explore.webp)
{width="500"}

### Grafana

Visit [Grafana - http://localhost:3000](http://localhost:3000/explore) with preconfigured datasources
  and a simple demo dashboard (_login with admin:admin_):
  ![Screenshot of Grafana showing the demo dashboard.](quick-start-grafana-demo.webp)
  {width="500"}


### Alerting UI

The demo will spawn two vmalert instances:
1. To run alerting rules [alerts-demo.yml](https://github.com/VictoriaMetrics/VictoriaStack/blob/main/deployment/rules/alerts-demo.yml)
   on metrics generated by OpenTelemetry demo and stored in VictoriaMetrics.
   Check rules state at [VictoriaMetrics alerting UI](http://localhost:8427/metrics/select/vmui/#/rules) or in [Grafana alerting](http://localhost:3000/alerting/list).
   ![Screenshot of VMUI for VictoriaMetrics showing the Alerting page.](quick-start-vmui-alerts.webp)
   {width="500"}

2. To run recording rules [logs-to-metrics.yml](https://github.com/VictoriaMetrics/VictoriaStack/blob/main/deployment/rules/logs-to-metrics.yml)
   converting logs from VictoriaLogs to metrics and writing them to VictoriaMetrics. Check results by querying
   [service_name:logs_count:1m](http://localhost:8427/metrics/select/vmui/#/?g0.range_input=30m&g0.relative_time=last_30_minutes&g0.tab=0&g0.expr=service_name%3Alogs_count%3A1m)
   or [severity_text:quotes_count:1m](http://localhost:8427/metrics/select/vmui/#/?g0.range_input=30m&g0.relative_time=last_30_minutes&g0.tab=0&g0.expr=severity_text%3Aquotes_count%3A1m)
   metric names in VictoriaMetrics UI.

## Clustered version

Run the stack locally with [otel-cluster.yml](https://github.com/VictoriaMetrics/VictoriaStack/blob/main/deployment/otel-cluster.yml)
Docker Compose file:

```sh
git clone https://github.com/VictoriaMetrics/VictoriaStack.git && cd VictoriaStack
make otel-cluster-up
```

It runs the same setup as the [single-node version](#single-node-version), but each database runs as a cluster of three
independently scalable components: insert, select, and storage.
In this topology, vmauth also balances requests among insert and select components.

![Diagram showing how the demo works in cluster mode. The OTel Demo sends metrics, logs, and traces to the OTel Collector, which connects to vmauth. vmauth routes the telemetry to the three clusters. Grafana and vmalert also connect to the three clusters via vmauth.](quick-start-compose-cluster.webp)

vmauth routes the requests, so the UI links stay
the [same as in the single-node version](#single-node-version).

The generator keeps pushing telemetry signals to the stack through the OpenTelemetry collector, and vmalert keeps
evaluating alerting rules. To stop the stack, run:

```sh
make otel-cluster-down
```

## Scalability

VictoriaMetrics, VictoriaLogs, and VictoriaTraces are highly scalable and simple to operate:

- one VictoriaMetrics cluster can handle thousands of tenants, billions of time series, millions of samples per second,
  thousands of read requests per second;
- one VictoriaLogs cluster can handle thousands of tenants, hundreds of terabytes of logs per day,
  thousands of read requests per second;
- one VictoriaTraces cluster can handle thousands of tenants, billions of spans per day.

[**inserts**, **selects** and **storages**](https://docs.victoriametrics.com/victoriastack/#architecture) can be scaled independently:

* add more **insert** replicas to increase ingestion throughput;
* add more **select** replicas to increase query throughput;
* add more **storage** shards to increase overall performance and overcome hardware limits like CPU, memory, IO,
  or filesystem capacity.

Components **insert** and **select** are stateless and can be scaled in response to demand. The **storage** component
is stateful and represents a storage shard. It is recommended carefully choosing the number of storage shards, as
scaling down is a complex operation.

## Retention and storage

All three databases store data on the local filesystem and delete it automatically once it falls outside the
configured retention:

| Database                                                                                                     | Default retention | How to configure                                                                                                            |
|--------------------------------------------------------------------------------------------------------------|-------------------|-----------------------------------------------------------------------------------------------------------------------------|
| [VictoriaMetrics](https://docs.victoriametrics.com/victoriametrics/single-server-victoriametrics/#retention) | 1 month           | `-retentionPeriod` (`1d`, `10y`)                                                                                            |
| [VictoriaLogs](https://docs.victoriametrics.com/victorialogs/#retention)                                     | 7 days            | `-retentionPeriod` (`1d`, `10y`), or a disk limit via `-retention.maxDiskSpaceUsageBytes`, `-retention.maxDiskUsagePercent` |
| [VictoriaTraces](https://docs.victoriametrics.com/victoriatraces/#retention)                                 | 7 days            | `-retentionPeriod` (`1d`, `10y`), or a disk limit via `-retention.maxDiskSpaceUsageBytes`, `-retention.maxDiskUsagePercent` |

<br>

Disk usage depends on your data and on how well it compresses. Estimate it with the capacity planning docs
for [VictoriaMetrics](https://docs.victoriametrics.com/victoriametrics/single-server-victoriametrics/#capacity-planning)
and [VictoriaLogs](https://docs.victoriametrics.com/victorialogs/#capacity-planning).

VictoriaStack has no object storage tiering: databases keep all data on local disks, and S3 or GCS is used
for [backups](https://docs.victoriametrics.com/victoriastack/deployment/#backups) only. Support for historical data in object
storage is in progress, see [this issue](https://github.com/VictoriaMetrics/VictoriaLogs/issues/48).

To store metrics longer Enterprise adds [downsampling](https://docs.victoriametrics.com/victoriametrics/single-server-victoriametrics/#downsampling)
and per-tenant [retention filters](https://docs.victoriametrics.com/victoriametrics/single-server-victoriametrics/#retention-filters).

## Next steps

The examples above run on generated telemetry. To build an installation of your own, go through the following steps:

1. [Choose deployment type](https://docs.victoriametrics.com/victoriastack/deployment/): Kubernetes operator, Helm charts, or
   binaries.
1. [Ingest your own data](https://docs.victoriametrics.com/victoriastack/ingestion/): point the telemetry agents you already run at
   VictoriaStack, or use vmagent and vlagent for higher throughput and reliable delivery.
1. [Query it back](https://docs.victoriametrics.com/victoriastack/querying/): MetricsQL for metrics, LogsQL for logs and
   traces.
1. [Build dashboards](https://docs.victoriametrics.com/victoriastack/ui/): Grafana or Perses for dashboarding, vmui for ad-hoc
   exploration.
1. [Set up alerting](https://docs.victoriametrics.com/victoriastack/alerting/): vmalert for rules evaluation, Alertmanager for
   notifications.
1. [Protect the installation](https://docs.victoriametrics.com/victoriastack/access-control/): vmauth for authentication, request
   routing, and TLS termination.
1. [Separate data between tenants](https://docs.victoriametrics.com/victoriastack/multitenancy/): one installation can serve
   thousands of isolated tenants.