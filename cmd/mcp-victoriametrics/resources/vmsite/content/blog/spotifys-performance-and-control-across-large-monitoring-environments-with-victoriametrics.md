---
draft: false
page: blog blog_post
authors:
  - Adam Yates
date: 2025-12-16
featuredPost: true
enableComments: true
title: "Spotify’s performance & control across large monitoring environments with VictoriaMetrics"
summary: "Spotify needed to replace its legacy in-house time series database to overcome stability and performance limitations, which would bring about query delays and timeouts. The Spotify observability team chose VictoriaMetrics to support efficient metric ingestion, querying, and alerting at scale."
categories: 
 - Product News
 - Company News
 - Customer Stories
tags:
 - spotify
 - victoriametrics
 - case study
 - time series database
 - monitoring
 - observability
images:
 - /blog/spotifys-performance-and-control-across-large-monitoring-environments-with-victoriametrics/spotify-case-study.webp
---

When your active time series is in the billions and the total number of data points you need to  monitor runs into the tens of trillions, you need a high-performance observability solution with operational simplicity.

Streaming behemoth Spotify is one such case. Their observability team chose VictoriaMetrics as the fastest monitoring and observability solution on the market.

## Spotify’s challenges
Spotify needed to replace its legacy in-house time series database (Heroic), which had become outdated, difficult to maintain, and inefficient at scale.

The goal was to implement a modern time-series database (TSDB) that could efficiently handle large-scale metric ingestion and querying, improve dashboard and alert performance, reduce operational overhead, and align with open observability standards such as Prometheus, OTel, and Grafana.

Difficulties Spotify’s observability team faced:
- **Stability and performance limitations** in its previous in-house TSDB, leading to query delays and timeouts
- **Limited feature parity** with modern observability systems
- A bespoke, **closed-source architecture** that restricted community support and maintainability
- **Growing maintenance overhead** as team familiarity with the legacy system decreased
- **Latency issues** with the existing alert engine
- Inconsistent metric models and **difficulty handling high-cardinality data**
- **Limited compatibility** with Prometheus and related open standards

Spotify evaluated multiple vendors and technologies before selecting VictoriaMetrics.

The alternative systems they tested during the evaluation phase showed limitations in scalability, compatibility with existing tooling, and flexibility of deployment models.

<p align="center" style="color: #71797E; font-family: Arial, sans-serif;"><i>VictoriaMetrics is “a robust, efficient, and flexible platform aligned with Spotify’s<br/> operational and architectural requirements”</i></p>
<p align="right">Lauren Roshore, Engineering Manager, Observability</p>

Spotify’s observability team had several evaluation criteria:
- Performance (data ingestion and query speed)
- Scalability for large, distributed workloads
- Cost efficiency (storage, licensing)
- Flexibility between self-managed and managed deployment models
- Compatibility with open-source standards
- Alerting infrastructure compatibility
- Operational maintainability

From the many different observability solutions on the market, VictoriaMetrics came out on top to support Spotify’s scalability and performance goals.

## Outcome of VictoriaMetrics adoption

<p align="center" style="color: #71797E; font-family: Arial, sans-serif;"><i>“Spotify’s transition to VictoriaMetrics has resulted in significant performance improvements across its monitoring stack, greater efficiency in engineering operations, and enhanced scalability to support future growth.”</i></p>
<p align="right">Lauren Roshore, Engineering Manager, Observability</p>

The solution provided a robust, efficient, and flexible platform aligned with the team’s operational and architectural requirements.

Some of the key benefits VictoriaMetrics now brings to Spotify’s observability:
- Significant improvements in data ingestion and query performance
- Prometheus-compatible APIs and query language 
- Simplified architecture for easier deployment and management
- Enhanced data retention and cost efficiency through downsampling and control features
- Support for both cloud and self-hosted deployments, offering high operational visibility
- Scalable, performant alerting infrastructure
- A predictable and transparent licensing model
- Noticeable improvements in dashboard responsiveness and alert evaluation times

Spotify is not stopping there in the coming months and years that involve VictoriaMetrics and observability in general. Their plans include UX and alert-annotation enhancements for a better on-call experience, anomaly detection in time-series data for advanced analytics, adoption of OTel, and stronger integration between reliability tooling (SLOs) and VictoriaMetrics/Grafana.

If you want to learn more about Spotify’s observability journey, [join us for our quarterly meetup on December 18, 2025](https://www.youtube.com/watch?v=yuZ_JkOx1uo). At the meetup, Spotify’s Observability Engineering Manager, Lauren Roshore, will explain “How & why we use VictoriaMetrics".