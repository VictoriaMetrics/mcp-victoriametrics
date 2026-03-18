---
draft: false
page: blog blog_post
authors:
  - Pablo Fernandez
date: 2026-02-25
enableComments: true
title: "VictoriaMetrics February 2026 Ecosystem Updates"
summary: "February 2026 updates deliver new LTS support, VMUI memory insights, queue alerts, jsonline output, resizable Web UI tables, and automatic snapshot expiry across the VictoriaMetrics Observability Stack."
toc: true
categories:
  - Product News
  - Company News
tags:
  - victoriametrics
  - victorialogs
  - observability
  - open source
  - kubernetes
  - new features
images:
  - /blog/victoriametrics-february-2026-ecosystem-updates/preview.webp
---

This month, we're thrilled to see [OpenAI using the VictoriaMetrics Stack](https://openai.com/index/harness-engineering/) internally — including VictoriaMetrics, VictoriaLogs, and VictoriaTraces — in their Harness engineering experiment, as shown in their architecture diagram.

It's a great way of combining observability and AI agents. The article is worth a read, but in short, they are setting up a temporary, isolated observability stack to provide Codex with metrics, logs, and traces to inform the AI and add valuable context to its coding loop. A unified stack delivers the scalability needed for ephemeral, AI-first workflows like those at OpenAI.

If you wish to replicate their setup, Alexander's complete guide on [observing AI agents using OpenTelemetry and VictoriaMetrics](https://victoriametrics.com/blog/ai-agents-observability/index.html) is the best starting point.

This roundup covers releases for:

- [VictoriaMetrics](#vm)
- [VictoriaLogs](#vl)

## VictoriaMetrics v1.136.0 and v1.135.0 {#vm}

What's new in VictoriaMetrics?

- New LTS line: v1.136 LTS is released, marking the end of support for v1.110 LTS.  
- Improved reliability and stability: addressing all known partition index issues and enhancing ingestion from slow or unreliable clients.  
- Enhanced VMUI visibility: introducing *Top Queries by memory usage*, context‑aware label autocomplete, and an improved drilldown panel that links directly to related source code.
- More flexible vmagent and vmauth configuration: including per‑backend write queues, buffered proxying, and proactive alerts for low queue space.  
- Optimized performance and usability: with smarter OpenTelemetry streaming, proxy environment variable support, and direct source code links from Grafana dashboards.  

**VictoriaMetrics [v1.136.0](https://docs.victoriametrics.com/victoriametrics/changelog/changelog_2026/#v11360) has been released**

- **New LTS line v1.136** has been started: this means v1.110 LTS has reached end of life and is no longer supported.
- Improved stability for pt-index: all known issues related to the partition index, which was introduced in [v1.133.0](https://docs.victoriametrics.com/victoriametrics/changelog/changelog_2026/#v11330), were addressed.
- Top queries by memory usage in VMUI: the [Top Queries](https://docs.victoriametrics.com/victoriametrics/single-server-victoriametrics/#top-queries) page now exposes queries with the highest average memory usage, adding a memory-based dimension to query analysis.
- Context-aware suggestions in VMUI: label autocomplete now takes already selected label filters into account, reducing noise in high-cardinality environments.
- Streaming OpenTelemetry ingestion: large OpenTelemetry requests are now streamed with proper buffer reuse and reset, lowering peak memory usage.
- Proxy environment variables support: all components now properly honor standard proxy variables `HTTP_PROXY`, `HTTPS_PROXY`, and `NO_PROXY`.
- View related source code in Logging rate panel: the `Logging rate` drilldown panel in the Grafana dashboard links directly to the corresponding source code.
- Multiple bug fixes and stability improvements across VictoriaMetrics components: including security updates to Go 1.26.0 and the base Alpine image. 

![Screenshot of Grafana](/blog/victoriametrics-february-2026-ecosystem-updates/logging-rate.webp)
<figcaption style="text-align: center; font-style: italic;">The logging rate panel shows related source code</figcaption>

Check the full [changelog](https://docs.victoriametrics.com/victoriametrics/changelog/changelog_2026/#v11360) for all changes.

**Highlights for [v1.135.0](https://docs.victoriametrics.com/victoriametrics/changelog/changelog_2026/#v11350)**

- More reliable ingestion from slow or external clients: vmauth can now buffer request bodies before forwarding them upstream. This prevents slow clients (e.g., IoT devices over poor networks) from exceeding concurrency limits and causing ingestion drops, even when backend capacity is available.
- Simpler vmagent setups for mixed backends: vmagent now supports different [`-remoteWrite.queues`](https://docs.victoriametrics.com/victoriametrics/vmagent/#common-flags) values per remote write URL. You can safely use a single vmagent instance for backends with strict ordering requirements (e.g., Mimir) and for VictoriaMetrics. Previously, users needed multiple agents to achieve the same result.
- Early warning before data loss: new vmagent alerts notify you 12 hours and 4 hours before the persistent queue runs out of space, giving time to react before metrics start getting dropped.
- Correct scrape size enforcement: [`-promscrape.maxScrapeSize`](https://docs.victoriametrics.com/victoriametrics/vmagent/#common-flags) is now applied to decompressed data. This makes scrape size limits predictable and prevents oversized payloads from slipping through due to compression.
- Correct `changes()` with variable scrape intervals: [MetricsQL `changes()`](https://docs.victoriametrics.com/victoriametrics/metricsql/#changes) function now behaves correctly when scrape intervals increase (e.g., from 30s to 60s), avoiding false positives where no real value change occurred. Important for queries that measure restarts.
- Other stability and correctness improvements.

See the [changelog](https://docs.victoriametrics.com/victoriametrics/changelog/changelog_2026/#v11350) for the full list.

## VictoriaLogs v1.46.0, v1.45.0, and v1.44.0 {#vl}

What's new in VictoriaLogs?

- Better ingestion scalability: ingestion throughput scales more efficiently on high-CPU servers.
- Kubernetes Collector stream tuning: new option (`-kubernetesCollector.streamFields`) lets you customize the default `_stream` field composition, giving operators direct control over log stream cardinality and grouping in busy clusters.
- Query & API stability: fields now sort alphabetically by default for more predictable outputs.
- Snapshot management & auto-expiry: new `/internal/partition/snapshot/delete` and `/internal/partition/snapshot/delete_stale?max_age=<d>` endpoints allow deleting individual partition snapshots and bulk-removing stale ones by age. New `-snapshotsMaxAge` flag enables automatic snapshot expiry.
- Web UI: enhanced Table View with resizable and reorderable columns, improved Group View, group-by toggle clears grouping, and hits chart defaults to `none` grouping.
- JSON nested objects: JSON ingestion can now preserve selected nested objects via `preserve_json_keys=...` (query arg) or `VL-Preserve-JSON-Keys: ...` (HTTP header), preventing dynamic keys from being flattened and inflating cardinality.

**Highlights for VictoriaLogs [v1.46.0](https://docs.victoriametrics.com/victorialogs/changelog/)**

- Query field ordering: Query responses now return fields in a stable alphabetical order by default, unless you explicitly define the output order (for example, via `fields` or `stats`), making API outputs more predictable for scripts, dashboards, and CLI use.
- vlagent JSON Lines output: vlagent now supports log delivery in jsonline format via `-remoteWrite.format=jsonline` on a per-remote-destination basis, enabling easier direct forwarding to systems such as Vector, Fluent Bit, or ClickHouse over HTTP.
- Web UI table view: table view in the Web UI now supports column resizing and drag-and-drop reordering, and it persists column layout preferences so operators can keep a practical table layout across sessions.
- Web UI group view: the group view in the Web UI now has a more polished interaction and uses zebra-row styling, making dense group listings easier to scan and requiring less visual effort.
- Web UI hits charts: the hits chart defaults to non-grouping, so the initial view shows total hits, while field-based splits appear only when explicitly selected.
- Web UI group-by behavior: clicking an already active group-by field now clears grouping instead of resetting to `_stream`, which better matches iterative investigation workflows.
- Web UI Log Context modal improvements: the Log Context modal now includes stream field chips, so nearby lines returned by `stream_context` queries display related stream metadata for easier context inspection.

Check out the [changelog](https://docs.victoriametrics.com/victorialogs/changelog/) for the full details.

**Highlights for [v1.45.0](https://docs.victoriametrics.com/victorialogs/changelog/#v1450)**

- Source code links in Grafana dashboard: the `VictoriaLogs internal logging` panel lets you jump to the code location that printed the message.
- Quality-of-life improvements for Grafana dashboards: clearer descriptions, better behavior with the Prometheus datasource, and a cluster variable to select every component in the cluster more easily.
- LogsQL adds [`ipv6_range` filter](https://docs.victoriametrics.com/victorialogs/logsql/#ipv6-range-filter): lets you [filter logs by IPv6](https://docs.victoriametrics.com/victorialogs/logsql/#ipv6-range-filter) using a range or CIDR (e.g., `client_ip:ipv6_range("2001:db8::/112")`).
- New timezone offset option: [the `/select/logsql/stats_query_range` API](https://docs.victoriametrics.com/victorialogs/querying/#querying-log-range-stats) now accepts an optional `offset` (timezone offset). This aligns "group stats by time buckets" with your timezone.
- Better ingestion on high-CPU hardware: ingestion scales even better on large servers with many CPU cores, enabling high-volume pipelines to push more logs under load.
- Offset is shown in the [Web UI](https://docs.victoriametrics.com/victorialogs/querying/#web-ui): the UI now shows offset in the bar-chart data, so bars line up with the selected timezone.
- New keyboard shortcuts: the Web UI query editor supports `Ctrl/Cmd` + `/` to comment or uncomment lines.

See the [changelog](https://docs.victoriametrics.com/victorialogs/changelog/#v1450) for more information.

**Improvements in [VictoriaLogs v1.44.0](https://docs.victoriametrics.com/victorialogs/changelog/#v1440)**

- Partition snapshot deletion: [new `/internal/partition/snapshot/delete?path=<snapshot-path>` endpoint](https://docs.victoriametrics.com/victorialogs/#partitions-lifecycle) deletes specific partition snapshots, making it easier to clean up finished backups and reclaim disk space.
- Bulk stale snapshot removal: [new `/internal/partition/snapshot/delete_stale?max_age=<d>` endpoint](https://docs.victoriametrics.com/victorialogs/#partitions-lifecycle) bulk-removes old snapshots by age (e.g., `max_age=1d`), ideal for scheduled housekeeping jobs to prevent snapshot buildup.
- Automatic snapshot expiry: `-snapshotsMaxAge` enables automatic periodic removal of snapshots, keeping storage usage predictable in long-running environments without needing external cron jobs.
- Enhanced snapshot creation: [the `/internal/partition/snapshot/create` API](https://docs.victoriametrics.com/victorialogs/#partitions-lifecycle) now supports `partition_prefix` for matching multiple daily partitions (e.g., `partition_prefix=202601`), simplifying monthly or yearly backup workflows.
- Journald remote IP support: `-journald.useRemoteIP` persists the sender’s IP address in the `remote_ip` field, aiding audits, incident response, and distinguishing logs from different hosts or forwarders.
- JSON nested objects: JSON ingestion preserves selected nested objects using `preserve_json_keys=...` (query argument) or `VL-Preserve-JSON-Keys: ...` (HTTP header), useful for dynamic keys that would otherwise increase cardinality.
- New Grafana dashboard: added [VictoriaLogs - Logging State](https://grafana.com/grafana/dashboards/24585-victorialogs-internal-state/) in [the `victorialogs-internal.json` file in the repository](https://github.com/VictoriaMetrics/VictoriaLogs/blob/master/dashboards/victorialogs-internal.json) for monitoring streams, fields, and ingestion behavior.
- [LogsQL pattern matching](https://docs.victoriametrics.com/victorialogs/logsql/#pattern-match-filter): added `pattern_match_prefix("...")` and `pattern_match_suffix("...")` to anchor matches to field value start and end, improving classification and detection precision.
- vlagent Kubernetes stream tuning: Kubernetes Collector supports `-kubernetesCollector.streamFields` to customize default `_stream` composition, controlling cardinality and grouping in busy clusters.
- Nanosecond query step precision: HTTP API accepts nanosecond `step` for `/select/logsql/hits` and `/select/logsql/stats_query_range`, enabling sub-millisecond bucketing.
- Web UI autocomplete improvements: now shows pipe titles with links to LogSQL docs, speeding up query building and onboarding.
- Web UI hits bar cumulative view: adds option for cumulative hits view to spot ramp-ups and totals across time ranges.
- Web UI configurable bar count: hits chart allows configuring bar count to balance detail and readability by time window.

Read the [changelog](https://docs.victoriametrics.com/victorialogs/changelog/#v1440) for all the details.

## Full Changelogs

For the full details on each release, read the changelogs:

- VictoriaMetrics: [v1.36.0](https://docs.victoriametrics.com/victoriametrics/changelog/changelog_2026/#v11360), [v1.35.0](https://docs.victoriametrics.com/victoriametrics/changelog/changelog_2026/#v11350)
- VictoriaLogs: [v1.46.0](https://docs.victoriametrics.com/victorialogs/changelog/), [v1.45.0](https://docs.victoriametrics.com/victorialogs/changelog/#v1450), [v1.44.0](https://docs.victoriametrics.com/victorialogs/changelog/#v1440)

Thank you for using VictoriaMetrics — [stay tuned for more updates](https://victoriametrics.com/contact-us/#communities).
