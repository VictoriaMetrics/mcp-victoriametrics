---
draft: false
page: blog blog_post
authors:
  - Denys Holius
date: 2025-12-01
enableComments: true
title: "Our latest updates across the VictoriaMetrics Observability ecosystem"
summary: "The VictoriaMetrics ecosystem continues to evolve rapidly, and the latest updates bring meaningful improvements across metrics, logs, and traces. Read the announcement for details."
categories: 
 - Product News
 - Company News
tags:
 - victoriametrics
 - victorialogs
 - victoriatraces
 - observability
 - open source
 - kubernetes
 - new features
images:
 - /blog/our-latest-updates-across-the-victoriametrics-observability-ecosystem/preview.webp
---

We’re excited to announce a set of updates across the entire VictoriaMetrics open source products suite — including VictoriaMetrics, VictoriaLogs, VictoriaTraces, the VictoriaMetrics Kubernetes Operator. These improvements bring better performance, stronger security, enhanced metadata visibility, and a smoother experience when running observability at scale.

## VictoriaMetrics v1.130.0 & v1.129.0 — Metadata, Performance, and Compatibility Improvements

The latest VictoriaMetrics releases introduce major enhancements across some of its core components: vmagent, vmalert, vmctl, as well as their toolset.

<strong>New in </strong><a href="https://docs.victoriametrics.com/victoriametrics/changelog/#v11300" target="_blank" rel="noopener noreferrer"><strong>v1.130.0</strong></a>:
- <strong>Metric metadata support (Prometheus-compatible API)</strong> </br >
    VictoriaMetrics can now scrape and store metadata using vmsingle or vmagent and expose it via `/api/v1/metadata`. Enable it with `-enableMetadata=true` across vmsingle, vmagent, and vminsert.
- <strong>Namespace metadata for Kubernetes targets</strong> </br >
    vmagent now can attach namespace metadata to discovered pods, services, endpoints, and more allowing richer filtering and discovery rules.
- <a href="https://docs.victoriametrics.com/victoriametrics/#vmui" target="_blank" rel="noopener noreferrer"><strong>VMUI</strong></a><strong> improvements</strong></br >
    The VictoriaMetrics UI is now faster and more intuitive.
- <strong>Better alert templating behavior in </strong><a href="https://docs.victoriametrics.com/victoriametrics/vmalert/" target="_blank" rel="noopener noreferrer"><strong>vmalert</strong></a></br >
    Templating errors now surface as annotation values while alerts continue to be generated, aligning with Prometheus behavior.

<strong>Highlights from </strong><a href="https://docs.victoriametrics.com/victoriametrics/changelog/#v11290" target="_blank" rel="noopener noreferrer"><strong>v1.129.0</strong></a>:
- <strong>vmalert</strong>
  - Faster startup for groups with long intervals.
  - New `now` template function for easier time calculations.
  - Added support for `alert_relabel_configs` per notifier.
  - Fixed UI search and HTML formatting in annotations.
- <strong>vmctl</strong> improvements including multiple `--remote-read-filter-label` flags.</br >
    This is useful in order to narrow down the data being migrated by using more precise filters.
- Improved performance for [stream aggregation](https://docs.victoriametrics.com/victoriametrics/stream-aggregation/) in vmagent with multiple rules. Push operation could utilize all CPU cores now.
- Load balancing fixes in vmauth (3% -> 0.04% deviation).
- s390x builds added for Linux.
- New `/remotewrite-*` endpoints for relabel config inspection.

## VictoriaLogs v1.38.0 — Security & Operational Intelligence Enhancements

<strong>New Features</strong>
- <a href="https://docs.victoriametrics.com/victorialogs/#how-to-delete-logs" target="_blank" rel="noopener noreferrer"><strong>Log deletion via HTTP API</strong></a></br >
    Enables GDPR-compliant deletion of stored logs. Intended for rare but critical actions such as data-breach cleanup or regulatory enforcement.
- <strong>Per-query redaction of sensitive fields</strong></br >
    Platform owners can hide specific fields (emails, credentials, SSNs, etc.) for some users while full data remains visible to authorized ones — essential for managed platforms and multi-tenant deployments.
- <strong>Slow query detection</strong></br >
    <a href="https://victoriametrics.com/products/victorialogs/" target="_blank" rel="noopener noreferrer"><strong>VictoriaLogs</strong></a> now supports slow query detection and automatically logs queries that exceed the configured execution time thresholds. This helps investigate, optimize slow queries and understand which queries cause long durations or spikes in CPU or memory usage.

Check out the full [Changelog](https://docs.victoriametrics.com/victorialogs/changelog/#v1380) for this release!

## VictoriaTraces v0.5.0 — OTLP/gRPC Support for Tracing Pipelines

VictoriaTraces now includes better compatibility with modern tracing infrastructures.

<strong>Key Feature</strong>
- **OTLP/gRPC ingestion for both single-node and cluster deployments**
    Requires -otlpGRPCListenAddr on VictoriaTraces or vtinsert. This simplifies ingestion from OpenTelemetry SDKs and collectors.

[This update](https://docs.victoriametrics.com/victoriatraces/changelog/#v050) makes adopting VictoriaTraces easier for teams standardizing on OpenTelemetry.

Speaking of gRPC, Zhu Jiekun [published an article](https://victoriametrics.com/blog/opentelemetry-without-grpc-go/) in which he revealed in more detail how support for OTLP/gRPC in the HTTP/2 + easyproto way was added to VictoriaTraces. 

## VictoriaMetrics Kubernetes Operator — VictoriaLogs and VictoriaTraces support

<a href="https://docs.victoriametrics.com/operator/changelog/#v0630" target="_blank" rel="noopener noreferrer"><strong>v0.63.0</strong></a><strong> major update</strong>
- This is a very significant release in that the operator now supports [VictoriaLogs](https://docs.victoriametrics.com/victorialogs/) and [VictoriaTraces](https://docs.victoriametrics.com/victoriatraces/) in both Single Node and Cluster-based versions.

<a href="https://docs.victoriametrics.com/operator/changelog/#v0650" target="_blank" rel="noopener noreferrer"><strong>v0.65.0</strong></a><strong> — Better Prometheus-Operator Compatibility</strong>
- **scrapeClass & scrapeClassName support across VM*Scrape objects** </br >
    Added to VMServiceScrape, VMPodScrape, VMProbe, VMScrapeConfig, VMStaticScrape, and VMNodeScrape. This enhancement smooths migration paths for users switching from Prometheus Operator to VictoriaMetrics.

## What’s Next

Across the VictoriaMetrics ecosystem, we continue working on:

- Deeper cross-linking between metrics, logs, and traces
- Further metadata improvements
- Enhanced performance for large ingest pipelines
- Improved security and compliance tooling

Thank you for choosing VictoriaMetrics — and stay tuned for more updates!
