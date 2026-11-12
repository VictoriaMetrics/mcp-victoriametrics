---
draft: false
page: "blog blog_post"
authors:
  - "Michelle Sebek"
date: 2026-11-05
title: "VictoriaMetrics Is Expanding in North America"
summary: "VictoriaMetrics Expands in North America as Spotify and MongoDB Use Our Open-Source Stack to Run Observability at Scale"
description: "Announcing building momentum for our next phase of growth with a fast-growing open-source community, continued product innovation and expanding enterprise adoption. We're increasing investment in North America to bring cost-efficient observability technology to more engineering teams managing increasingly complex IT and AI infrastructure."
enableComments: true
categories:
  - "Company News"
tags:
  - "victoriametrics"
  - "open source"
  - "observability"
keywords:
  - "observability solutions"
  - "logs"
  - "metrics"
  - "traces"
  - "usa"
  - "north america"
images:
  - "/blog/victoriametrics-is-expanding-in-north-america/preview.webp"
---

We’re an open-source and Enterprise observability solutions provider for infrastructure, application, and AI workloads, built to give engineering teams more visibility with less infrastructure and less operational overhead.

Spotify, MongoDB and Grammarly run observability at scale on our open-source stack, and our open-source components have over 1B downloads.

We're headquartered in the US, bootstrapped and profitable, and today we're announcing a larger investment in our North American team.

VictoriaMetrics started in 2018 with four engineers who had run observability systems that couldn't keep up at scale. They wrote a time series database that used less memory, less disk and fewer moving parts, and they made it available under Apache 2.0. That is still how we build. Every product we ship has to be simple to run, reliable when traffic spikes and efficient with the hardware you pay for.

## Who runs it in production: When enterprise complements open source

The clearest proof is the companies that moved their own observability data onto it. Spotify R&D replaced its in-house time series database with VictoriaMetrics Enterprise and documented ingestion of about [78 million data points per second](https://victoriametrics.com/case-studies/spotify), with dashboards and queries up to 10x faster. Grammarly, now Superhuman, [reported a 10x cost reduction after migrating](https://www.grammarly.com/blog/engineering/monitoring-with-victoriametrics/), and its engineering team called the move one of the smoothest transitions in the company's history. DreamHost [cut memory use by 80% while running 76 million active time series](https://victoriametrics.com/blog/dreamhost-case-study/).

MongoDB runs VictoriaMetrics Enterprise, self-managed on its own infrastructure, for global [metrics observability across AWS, Google Cloud and Microsoft Azure](https://victoriametrics.com/case-studies/mongodb/). Data stays in its home region while engineers query across all three clouds, and the team kept its Grafana dashboards, alerting rules and recording rules unchanged during the move.

<p align="center" style="color: #71797E; font-family: Arial, sans-serif;"><i>"VictoriaMetrics Enterprise has given us the visibility and flexibility 
to operate a globally distributed platform with confidence. We were able to build on the tools and workflows our engineers already knew, while gaining the scalability and query performance we needed to support our next stage of growth and continue delivering the experience our customers expect," said persephone thorn, Manager, Observability at MongoDB.</i></p>

Practitioners on [PeerSpot report the same pattern](https://www.peerspot.com/products/victoriametrics-reviews). One software engineer who replaced Prometheus with VictoriaMetrics measured about 70% less RAM, 75% less disk, 3x faster writes and 7x faster p95 query latency. An SRE lead at a financial services firm with more than 10,000 employees has run it for four years as a high-availability cluster split across data centers, and rated it 10 out of 10. Another reviewer runs it as long-term metric storage for an AI platform with many GPU nodes.

## Open source first

Everything above starts in the open. Our repositories have passed 21,000 GitHub stars, and in the first half of 2026 the team shipped more than 70 releases and merged 1,715 pull requests, not counting bots. You can read every line, run it in production without a contract, and file an issue that a core engineer will answer.

The same engine runs three ways. Open source is where most teams start and where many stay, self-hosted anywhere. Enterprise adds long-term support releases, advanced features and direct access to the engineers who wrote the code, self-hosted on-prem or in your own AWS, Google Cloud or Azure account. [VictoriaMetrics Cloud](https://victoriametrics.com/products/cloud/) is for teams that would rather we operate the infrastructure. You pick the operating model, and the software doesn't get worse in any of them.

## More than metrics

Metrics are where we earned trust, and we've applied the same approach to logs and traces. [VictoriaLogs](https://victoriametrics.com/products/victorialogs/) is generally available as open source, in Enterprise and in VictoriaMetrics Cloud, and runs as a single binary or a cluster. [VictoriaTraces](https://victoriametrics.com/products/victoriatraces/) stores distributed traces on the VictoriaLogs engine, so teams already running VictoriaLogs reuse the same engine and the same query language, LogsQL. VictoriaTraces is open source, and Enterprise and now supported in VictoriaMetrics Cloud alongside metrics and logs, ingesting OpenTelemetry over OTLP/HTTP. You don't need to adopt all three. Start with the signal that hurts most, and add the next one without touching the rest of your setup.

## Expanding in North America

Our global headcount grew 50% this year without venture funding, and we aren’t raising any. Many of our Enterprise customers started with open source, and their contracts fund the engineering behind every release, open source included.

In 2026 we added to our US-based staff across sales, customer success and developer advocacy. We’re still [hiring in North America](https://careers.victoriametrics.com/), including support and engineering staff who work North American business hours.

<p align="center" style="color: #71797E; font-family: Arial, sans-serif;"><i>“As open-source adoption grows and our solution expands beyond metrics into logs and traces, we’re investing in the people, technology and community needed to support engineering teams at every stage of scale,” said Artem Navoiev, CEO and co-founder of VictoriaMetrics.</i></p>

## If you already run VictoriaMetrics

For teams already running VictoriaMetrics, the expansion means more of us working North American hours, so more of your questions get answered inside your own working day. Enterprise customers also get long-term support release lines for VictoriaMetrics, with bug and security fixes for 12 months, so upgrades happen on your schedule.

[We'd also like to hear how you run it](https://www.peerspot.com/interviews/new?product_id=41853). If you have a migration story, a scale number you're proud of, or a complaint, tell us in the [community Slack](https://slack.victoriametrics.com/) or at the booth.

## See you at KubeCon next week

Next week we'll be at KubeCon + CloudNativeCon North America in Salt Lake City, booth #218. Bring the query that's been slowing your dashboards down and ask an engineer to run it, or meet the people building our North American team. If you want to try it before then, the [quick start](https://docs.victoriametrics.com/victoriametrics/quick-start/) takes a few minutes, and [VictoriaMetrics Cloud](https://console.victoriametrics.cloud/signUp) has a free trial.

Simple, reliable & efficient observability for everyone. By engineers, for engineers.

## Read more

- [VictoriaLogs vs Loki, Elasticsearch, OpenSearch, SigNoz and ClickStack on structured log ingestion](https://victoriametrics.com/blog/structured-log-ingestion-benchmark/)
- [Benchmarking Kubernetes log collectors: vlagent, Vector, Fluent Bit, OpenTelemetry Collector, and more](https://victoriametrics.com/blog/log-collectors-benchmark-2026/)
- [Announcing 1B+ downloads and product development with logs, traces, metrics](https://victoriametrics.com/blog/announcing-1b-downloads-and-product-development-with-logs-traces-metrics)
