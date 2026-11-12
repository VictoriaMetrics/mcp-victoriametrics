---
title: UI integrations
menu:
  docs:
    weight: 4
    identifier: victoriastack-ui-integrations
    parent: victoriastack
tags:
  - metrics
  - logs
  - traces
---

VictoriaStack has no unified interface of its own that covers all three signals, but it is in-progress. For now, 
**we recommend using [Grafana](#grafana)**: it connects to all three databases as datasources, so metrics, logs, and traces end up on the
same dashboards and can be [correlated](#correlations) with each other. [Perses](#perses) is an alternative visualization
UI that can be used with VictoriaStack as well.

Every database also goes with a built-in UI called [vmui](#vmui). It is available out of the box and is great for ad-hoc
data exploration, but it doesn't allow building dashboards, correlating different signal types, or managing access.

|                                 | [Grafana](#grafana) or [Perses](#perses)   | [vmui](#vmui)                                             |
|---------------------------------|--------------------------------------------|-----------------------------------------------------------|
| **What it is for**              | dashboards and cross-signal investigations | ad-hoc exploration and troubleshooting of a single signal |
| **Signals in one place**        | all three                                  | each database has its own vmui for a specific signal      |
| **Setup**                       | install and configure the datasources      | none, it ships with the database                          |
| **Dashboards**                  | yes                                        | no                                                        |
| **Correlation between signals** | [yes](#correlations)                       | no                                                        |

<br>

## Grafana

Grafana connects to all three databases through datasource plugins:

* VictoriaMetrics via [Prometheus](https://docs.victoriametrics.com/victoriametrics/integrations/grafana/#prometheus-datasource) or [VictoriaMetrics](https://docs.victoriametrics.com/victoriametrics/integrations/grafana/#victoriametrics-datasource) datasource plugins.
* VictoriaLogs via [VictoriaLogs datasource plugin](https://docs.victoriametrics.com/victorialogs/integrations/grafana/).
* VictoriaTraces via [Jaeger](https://docs.victoriametrics.com/victoriatraces/querying/grafana/#jaeger-datasource) or [Tempo](https://docs.victoriametrics.com/victoriatraces/querying/grafana/#grafana-tempo-datasource) datasource plugins.

See an example of [provisioning Grafana](https://grafana.com/docs/grafana/latest/administration/provisioning/) with corresponding datasources:
```yaml
apiVersion: 1
datasources:
  - name: VictoriaMetrics
    type: prometheus
    url: http://vmstack:8427/metrics/select/prometheus/
  
  - name: VictoriaLogs
    type: victoriametrics-logs-datasource
    url: http://vmstack:8427/logs/
  
  - name: VictoriaTraces
    type: jaeger
    url: http://vmstack:8427/traces/select/jaeger/
```

> See practical [Grafana provisioning examples](https://github.com/VictoriaMetrics/VictoriaStack/tree/main/deployment/provisioning/datasources) for VictoriaStack.
> The [QuickStart](https://docs.victoriametrics.com/victoriastack/quick-start/) runs Grafana with these datasources and a demo
> dashboard already provisioned.

With the datasources in place, Grafana covers everything a team does with the stack day to day: the query builders in the
plugins provide a queryless experience for exploring every signal, dashboards and
[alerting views](https://docs.victoriametrics.com/victoriastack/alerting/) live in one place, and a single investigation can move
between metrics, logs, and traces. For the query languages behind these interfaces,
see [Querying](https://docs.victoriametrics.com/victoriastack/querying/).

## Perses

[Perses](https://perses.dev/) is a lightweight, Kubernetes-native dashboarding tool that keeps dashboards as code.
Like Grafana, it integrates with VictoriaStack via datasource plugins:

* VictoriaMetrics integrates via [Prometheus datasource plugin](https://docs.victoriametrics.com/victoriametrics/integrations/perses/).
* VictoriaLogs integrates via [VictoriaLogs datasource plugin](https://docs.victoriametrics.com/victorialogs/integrations/perses/).
* VictoriaTraces integrates via [Jaeger](https://github.com/perses/plugins/tree/main/jaeger) or [Tempo](https://perses.dev/docs/plugins/tempo/) datasource plugins.

Perses covers dashboards, while cross-signal [correlations](#correlations) functionality is more mature in Grafana so far.

## Correlations

Signals can be correlated when they share common attributes that identify the same system or event: a trace ID present
in a log entry to correlate traces and logs, or the same service name present in all three signals to correlate all of them.
Correlation is the reason to put all three databases into one UI, and Grafana is the recommended place for it thanks to its
[correlations](https://grafana.com/docs/grafana/latest/administration/correlations/) and derived fields features.

![ui-correlations.webp](ui-correlations.webp)

Depending on the Grafana datasource plugin, there could be multiple correlations available:

1. Trace to logs, log to trace, log to metrics - see [correlations via VictoriaLogs plugin](https://docs.victoriametrics.com/victorialogs/integrations/grafana/#correlations).
1. Trace to metrics, metric to logs, metric to traces - see [correlations via VictoriaMetrics plugin](https://docs.victoriametrics.com/victoriametrics/integrations/grafana/datasource/#correlations).
   The Prometheus datasource also supports metric to logs and metric to traces correlations.
1. Plugins Tempo, Jaeger, and Zipkin can correlate with logs or metrics using [Trace to logs](https://grafana.com/docs/grafana/latest/explore/trace-integration/#trace-to-logs)
   and [Trace to metrics](https://grafana.com/docs/grafana/latest/visualizations/explore/trace-integration/#trace-to-metrics) feature.

> See examples of correlation settings configured for each datasource plugin in the [VictoriaStack provisioning](https://github.com/VictoriaMetrics/VictoriaStack/tree/main/deployment/provisioning/datasources) demo. 

Logs and traces can also be turned into metrics
with [recording rules](https://docs.victoriametrics.com/victoriastack/alerting/#rules), which makes them available on the same
dashboards and in the same queries as the rest of the metrics.

## vmui

vmui is the web interface built into the databases. It requires no setup and gives every signal a native UI for ad-hoc
exploration, focused on the signal and use case specific to the signal:

- [VictoriaMetrics UI](https://docs.victoriametrics.com/victoriametrics/#vmui): MetricsQL queries, 
  metrics and [cardinality explorers](https://docs.victoriametrics.com/victoriametrics/#cardinality-explorer),
  top and active queries, query tracing, relabeling debugging.
- [VictoriaLogs UI](https://docs.victoriametrics.com/victorialogs/querying/#web-ui): LogsQL queries,
  streams and fields, live tailing, log context.
- [VictoriaTraces UI](https://docs.victoriametrics.com/victoriatraces/querying/#web-ui): LogsQL queries,
   filtering and exploring traces.

> Have a quick look at the vmui demo for [VictoriaMetrics](https://docs.victoriametrics.com/playgrounds/#victoriametrics-playground),
> [VictoriaLogs](https://docs.victoriametrics.com/playgrounds/#victorialogs-playground) and [VictoriaTraces](https://docs.victoriametrics.com/playgrounds/#victoriatraces-playground) at our playgrounds.

vmui complements Grafana with tools specific to each signal. VictoriaMetrics vmui supports time series visualization,
cardinality and metric usage analysis, expensive query analysis, and query performance tracing. VictoriaLogs vmui supports
log exploration, analytics, and log context. VictoriaTraces vmui supports LogsQL queries and trace exploration.

## vmanomaly UI

The [vmanomaly UI](/anomaly-detection/ui/) provides a vmui-like interface for exploring queries, comparing anomaly detection models, and backtesting configurations. Although `vmanomaly` is an Enterprise component, you can explore its UI in the public playground with [Metrics](https://play-vmanomaly.victoriametrics.com/metrics/), [Logs](https://play-vmanomaly.victoriametrics.com/logs/), or [Traces](https://play-vmanomaly.victoriametrics.com/traces/) data sources.
