---
draft: false
page: blog blog_post
authors:
  - Jose Gomez-Selles
date: 2026-11-03
title: "What's New in VictoriaMetrics Cloud Q3 2026? VictoriaTraces, a Revamped Overview and Monitoring, Alerting for Logs and Traces, and a new Beta coming soon"
enableComments: true
featuredPost: true
summary: "Q3 2026 brought VictoriaTraces to VictoriaMetrics Cloud, completing the stack with metrics, logs and traces. We also revamped the Overview and Monitoring pages for full transparency, added alerting for logs and traces, enabled stream aggregation for VictoriaMetrics Single, and announced code-named project Visum: our upcoming visualization tool."
description: "Discover VictoriaTraces in VictoriaMetrics Cloud, the new Overview and Monitoring pages, alerting rules for logs and traces, stream aggregation and the Visum beta."
categories:
  - Product News
tags:
  - victoriametrics
  - victoriatraces
  - cloud
  - monitoring
  - observability
  - tracing
  - roadmap
images:
  - /blog/q3-2026-whats-new-victoriametrics-cloud/preview.webp
---

Q3 2026 has been one of the greatest quarters for VictoriaMetrics Cloud, and that is thanks to
the whole team, including open source, of course. The big milestone is that
**VictoriaTraces is now part of [VictoriaMetrics Cloud](https://victoriametrics.com/products/cloud/)**,
which means that the full `VictoriaStack` lives in one managed place, with a single bill
and the same experience users already know.

Alongside that, we kept pushing on transparency. The **Overview and Monitoring pages have been
revamped** to show users everything we look at when we investigate a deployment ourselves. We also added
**Alerting Rules for Logs and Traces**, enabled **stream aggregation for VictoriaMetrics Single**
instances, and much more. Also, stay tuned because we're launching a closed beta for the _not-that-much_ secret
project we mentioned last quarter, code-named `Visum`: **our upcoming visualization tool**.

As usual, if you missed our **Quarterly Live Update**, you can watch the full recording here:

{{< youtube id="6mn3LdddObw" start="3204" >}}

Let's go!

## VictoriaTraces in VictoriaMetrics Cloud

Users of VictoriaMetrics Cloud expect everything to be managed. So we've worked this quarter to finally provide a
**complete and centralized observability suite** with all signals: metrics, logs and traces.

With VictoriaTraces, the full `VictoriaStack` is ready to be installed and managed by our team within
minutes. Users only need to pick the type, the region and a tier, and the deployment will be up and
running in no time.

<figure style="width:70%; margin:1.5em auto;">
  <img src="/blog/q3-2026-whats-new-victoriametrics-cloud/traces-create-deployment.webp"
       alt="Deploy VictoriaMetrics, VictoriaLogs or VictoriaTraces in a few clicks"
       style="display:block; width:100%; margin:auto;" />
  <figcaption style="margin:.6em 0 0; text-align:center; font-style:italic; font-size:.9em; opacity:.8;">
    Deploy VictoriaMetrics, VictoriaLogs or VictoriaTraces in a few clicks
  </figcaption>
</figure>

### Exploring Spans

When we were about to release VictoriaTraces in VictoriaMetrics Cloud, the first question was: what do we give
users to explore their data? Then the team [started working](https://github.com/VictoriaMetrics/VictoriaTraces/pull/248)
on it. Those of us who have been in the tracing world for years can tell that this UI is fantastic.

It is called `VTUI` (first available in open source, of course), and it's now included in
VictoriaMetrics Cloud. This means that there's no extra tooling needed anymore.

> [!TIP]
> If you like what you're seeing, wait for what's coming at the end of this post!

From the **Explore** section, now users can:

* Select the **area of time** to investigate within a heatmap that highlights errors
* **Filter** by service, operation, duration, errors, span types and more
* **Query** with [LogsQL](https://docs.victoriametrics.com/victorialogs/logsql/), helped by `Examples`,
`Autocomplete` and `Tenant selection`
* Check the `Query history`, which is super useful in the middle of an investigation

{{< carousel width="90%" height="580" >}}
/blog/q3-2026-whats-new-victoriametrics-cloud/traces-explore.webp | The Explore view: filters, heatmap and traces list
/blog/q3-2026-whats-new-victoriametrics-cloud/traces-explore-preview.webp | Trace preview, opened at the right of the Explore view
/blog/q3-2026-whats-new-victoriametrics-cloud/traces-view.webp | Trace view: span timeline, events and links
{{< /carousel >}}

One of our favorite features is the `Trace preview`. Selecting a trace shows how its duration
compares to the rest of the traces for the same operation on a color-coded heatmap, so it is immediately clear where it sits.
It also shows the `Critical Span`, the `Errors` and all the `Attributes`, without leaving the list.

For deeper investigations, the `Trace View` shows the span timeline, the
critical span, `Events` (even if they are deprecated in OpenTelemetry, and we won't open that Pandora's box
today) and `Links`.

### Sending Traces

As in every VictoriaMetrics Cloud deployment, `Integrations` are one click away. VictoriaTraces is
compatible with the **OpenTelemetry Collector**, and even faster with [vtagent](https://docs.victoriametrics.com/victoriatraces/vtagent/),
and data can be visualized in **Grafana** or **Perses** too (and soon without leaving VictoriaMetrics Cloud!).

{{< carousel >}}
/blog/q3-2026-whats-new-victoriametrics-cloud/traces-integrations.webp | Integrations available for VictoriaTraces
/blog/q3-2026-whats-new-victoriametrics-cloud/traces-otel-guide.webp | Interactive, paste-ready guide for the OpenTelemetry Collector
{{< /carousel >}}

Each guide is interactive and paste-ready. Users wanting to integrate OpenTelemetry without reading
documentation (or, let's be honest, spending tokens on it) can come here, copy the snippets and
get started easily.

### Choosing a Tier

VictoriaTraces in VictoriaMetrics Cloud is based on the open source project, and its single-node instances are already
**multitenant**, so there's no need for clusters to get started. As usual, tiers are based on
**compute units**: users pick as many resources as needed and are billed for them.

> [!TIP]
> One of the benefits of this system is that spikes in ingestion or reading are not charged.

Of course, deployments can be **scaled** up and down at any time, **backups** and **support**
are included, and storage can grow as needed, so there is no need to commit to a lot of it from day one.

{{< carousel >}}
/blog/q3-2026-whats-new-victoriametrics-cloud/traces-tiers.webp | Choose the Capacity Tier that fits the workload
/blog/q3-2026-whats-new-victoriametrics-cloud/traces-recommendations.webp | Each tier shows the load it was benchmarked for, as a recommendation.
{{< /carousel >}}

To help choosing a tier, we ran some **benchmarks**. And as every benchmark,
take them with a grain of salt. Many benchmarks out there only measure writes. In ours,
we did both **reading and writing** at the same time, and we'd actually say **we leaned towards over-reading on purpose**. This way we are not
overselling anything: no promises about ingesting petabytes of spans per millisecond that crash
as soon as somebody runs a query. Anyone may find the details and recommendations in the
[VictoriaTraces tiers documentation](https://docs.victoriametrics.com/victoriametrics-cloud/deployments/victoriatraces-tiers/).

We will keep running these exercises but, for help deciding, just write to us and we will
help pick one, or arrange a PoC for big volumes.

> [!TIP]
> Ready to see traces in action? [Sign up for VictoriaMetrics Cloud](https://console.victoriametrics.cloud/signUp?utm_source=blog&utm_medium=blog_post&utm_campaign=q3-2026-whats-new&utm_content=traces-free-trial)
> and try it with a one-month free trial, no credit card required.

## Better Overview: Everything at a Glance

We have been talking about making our tiers more transparent for a while. We are showing components,
metrics and resource usage, so it was time to improve the first thing users see when opening a deployment page.

This is how the `Overview` page looks now:

<figure style="width:90%; margin:1.5em auto;">
  <img src="/blog/q3-2026-whats-new-victoriametrics-cloud/overview-new.webp"
       alt="The new Overview page: endpoint, topology, KPIs, capacity and health at a glance"
       style="display:block; width:100%; margin:auto;" />
  <figcaption style="margin:.6em 0 0; text-align:center; font-style:italic; font-size:.9em; opacity:.8;">
    The new Overview page: endpoint, topology, KPIs, capacity and health at a glance
  </figcaption>
</figure>

Now this page shows everything relevant for users **at a glance**:

* **Basic info and topology**: copy the endpoint to integrate faster, and check the topology that has
been deployed. This is especially relevant for `Cluster` deployments, where every component and its replicas can be seen.
* **Real time KPIs**: for VictoriaMetrics instances, for example, that means ingestion rate, active time series, read rate and new series over 24 hours.
* **Resource consumption**: as our tiers are based on compute units, CPU, RAM and storage can be checked
per component. If `vmstorage` is almost out of RAM, users can open a support ticket and we will help, or scale the
tier on their own.
* **Health and alerts**, with links to navigate to relevant pages.
* **Last changes**: showing latest events from `Audit logs`, such as new deployments or access tokens created.

We also renewed the **Deployments page**, so the whole fleet may be managed in a few clicks. Many users
have different VictoriaMetrics instances, for staging and production, together with VictoriaLogs and VictoriaTraces deployments.
Now all of them can be seen in one page, with the relevant information and links of each one available at a glance.

<figure style="width:90%; margin:1.5em auto;">
  <img src="/blog/q3-2026-whats-new-victoriametrics-cloud/deployments-fleet.webp"
       alt="The new Deployments page: the whole fleet of metrics, logs and traces deployments in one place"
       style="display:block; width:100%; margin:auto;" />
  <figcaption style="margin:.6em 0 0; text-align:center; font-style:italic; font-size:.9em; opacity:.8;">
    The new Deployments page: the whole fleet of metrics, logs and traces deployments in one place
  </figcaption>
</figure>

## Better Monitoring: Do It Like the Pros

The other feature we're very proud of is the new `Monitoring` page, that we advanced in the
[Q2 update](https://victoriametrics.com/blog/q2-2026-whats-new-victoriametrics-cloud/).
It natively shows the same information we use ourselves when we investigate issues or want to know
how an instance is performing. Those who have followed our [monitoring guidelines](https://victoriametrics.com/blog/victoriametrics-monitoring/) will
recognize these visualizations from the public [Grafana dashboards](https://grafana.com/orgs/victoriametrics/dashboards). Now they are included in VictoriaMetrics Cloud... **Natively**!

{{< carousel >}}
/blog/q3-2026-whats-new-victoriametrics-cloud/monitoring-overview.webp | The Monitoring overview with the key metrics of a deployment
/blog/q3-2026-whats-new-victoriametrics-cloud/monitoring-components.webp | All the stats of each component
/blog/q3-2026-whats-new-victoriametrics-cloud/monitoring-split-view.webp | Split view to compare two components side by side
{{< /carousel >}}

We learned that sharing this with everybody is better:

* **Transparent tiers**: it is easier to understand why a tier fits (or doesn't) a given workload, and scaling and usage become more evident.
* **Per-component breakdown**: available to check each component of a deployment on its own tab.
* **All data available**: all the relevant metrics are there to explore.

And one more cool feature: the screen can be **split**, so one side can scroll while the other stays in place
for comparing. This is perfect for **correlating events** that produced errors, for example, looking at two
components side by side.

## Alerting for Logs and Traces

People were asking a lot for this one: **Recording and Alerting rules** for **VictoriaLogs and VictoriaTraces**.
As we were saying earlier, VictoriaMetrics Cloud users want everything to be managed in the Cloud, so this was _a must_.

Those already using alerting with metrics know how it works. A `vmalert` component is installed when creating
rules, either alerting or recording ones. It reads from the deployment, generates the results as new metrics,
and an optional **notifier**, based on `Alertmanager`, may be deployed to send alerts to Slack, email or any webhook.
Readers can learn more in the [alerting setup docs](https://docs.victoriametrics.com/victoriametrics-cloud/alertmanager-setup-for-deployment/).

<figure style="width:90%; margin:1.5em auto;">
  <img src="/blog/q3-2026-whats-new-victoriametrics-cloud/alerting-metrics-setup.webp"
       alt="Alerting for metrics: vmalert reads from the deployment, writes rule results back as new metrics and sends firing alerts to the optional notifier"
       style="display:block; width:100%; margin:auto;" />
  <figcaption style="margin:.6em 0 0; text-align:center; font-style:italic; font-size:.9em; opacity:.8;">
    Alerting for metrics: <code>vmalert</code> reads from the deployment, writes rule results back as new metrics and sends firing alerts to the optional notifier
  </figcaption>
</figure>

For logs and traces it works in a very similar way, with one difference:

* **Alerting rules** work directly on the VictoriaLogs or VictoriaTraces deployments.
* To have **recording rules**, or alert states **surviving `vmalert` restarts**,
a VictoriaMetrics instance is needed to store the results.

<figure style="width:90%; margin:1.5em auto;">
  <img src="/blog/q3-2026-whats-new-victoriametrics-cloud/alerting-logs-traces-setup.webp"
       alt="Alerting for logs and traces: the same flow, with an optional linked VictoriaMetrics deployment that stores recording rule results and alert state"
       style="display:block; width:100%; margin:auto;" />
  <figcaption style="margin:.6em 0 0; text-align:center; font-style:italic; font-size:.9em; opacity:.8;">
    Alerting for logs and traces: the same flow, with an optional linked VictoriaMetrics deployment that stores recording rule results and alert state
  </figcaption>
</figure>

That is why the ability to **link a metrics deployment** to logs or traces deployments was added. It
is just another optional component when configuring alerting. Both deployments need to sit in the
same region.

This is a first step, and we have more planned for our whole alerting experience, so stay tuned!

## Stream Aggregation for VictoriaMetrics Single

The last big feature of the quarter is another one that our customers were asking for a lot.
We talked about recording rules before: a `vmalert` that reads a lot of data and runs queries that
can be quite complex.

But those who know the key features of VictoriaMetrics may want to use another great
feature instead: [stream aggregation](https://docs.victoriametrics.com/victoriametrics/stream-aggregation/). The [difference between both mechanisms](https://docs.victoriametrics.com/victoriametrics/stream-aggregation/#recording-rules-alternative)
relies on **when** the computation happens.

With recording rules, all the raw series are stored first and computed afterwards.
With stream aggregation, samples are **aggregated in memory on the fly**, in a single pass, and then only the aggregates are written
(and the raw input only when users choose to keep it).

<figure style="width:90%; margin:1.5em auto;">
  <img src="/blog/q3-2026-whats-new-victoriametrics-cloud/recording-rules.webp"
       alt="Recording rules: compute after storage"
       style="display:block; width:100%; margin:auto;" />
  <figcaption style="margin:.6em 0 0; text-align:center; font-style:italic; font-size:.9em; opacity:.8;">
    Recording rules: compute after storage
  </figcaption>
</figure>

<figure style="width:90%; margin:1.5em auto;">
  <img src="/blog/q3-2026-whats-new-victoriametrics-cloud/stream-aggregation.webp"
       alt="Stream aggregation: compute before storage"
       style="display:block; width:100%; margin:auto;" />
  <figcaption style="margin:.6em 0 0; text-align:center; font-style:italic; font-size:.9em; opacity:.8;">
    Stream aggregation: compute before storage
  </figcaption>
</figure>

It is now available for **VictoriaMetrics Single instances** in the Cloud. It can be configured under
`Settings → Advanced settings → Stream aggregation (YAML)`. It is best for when:

* A recording rule is slow because it scans many raw series
* Users don't need the raw series, only the aggregate
* High-cardinality labels or scrape frequency are inflating storage

> [!TIP]
> Stream aggregation adds CPU and memory load, so a higher tier may be needed. When unsure, talk to us and we will help size it.

## Other Updates, Improvements and Bug Fixes

As always, there are many smaller changes that came from user feedback:

* **API and Terraform support for VictoriaLogs and VictoriaTraces.** The whole stack can be managed as code,
learn more in the [API](https://docs.victoriametrics.com/victoriametrics-cloud/api/) and
[Terraform](https://docs.victoriametrics.com/victoriametrics-cloud/integrations/terraform/) docs.
* **Access tokens.** Better handling of tokens, with system-managed ones instead of random ones.
* **Logs live tailing** is now available in the console UI. We know our users (and founders!) love this feature.
It comes with VLUI, but we needed to do some internal work to make it happen in VictoriaMetrics Cloud.
* **Explore page** added several UX improvements.
* **Docs** now include rule presets integrations with several other pages being improved, as usual.

## What's Next?

### Clusters, Tiers and More

We are working on several features, some of which we have been mentioning in previous updates:

* **New compute-based model**: the clusters are ready and we are finalizing the implementation. The ladder of
scaling will be much smoother, and the pricing and everything around it will be much better.
* **Clusters for VictoriaLogs and VictoriaTraces**: we will start working on them after releasing the new Metrics
clusters. We had a very interesting conversation about this just this week, so expect news soon.
* **SOC 2**: the evaluation is still ongoing.
* **S3 storage**: as soon as it is ready in open source, we will work to make it happen in VictoriaMetrics Cloud.
* **Generic OAuth**: more flexibility to authenticate. Today users can log in with email and password or
with Google, and we know that SSO is a big ask. This is definitely on the roadmap.

### The Big Thing: project code-named `Visum`

We know: clusters, tiers, all of that is exciting. But for us, the biggest thing coming is this one.

We have metrics, we have logs, and now we have traces. So we have everything centralized in VictoriaMetrics Cloud.
What else do we need? **A super cool visualization tool**.

Its project name is **`Visum`** and we will be launching a **closed beta** very soon (we are talking weeks, so it is coming soon).
It will be a single pane of glass for all observability data. Users will be able to:

* See **services end to end** and their dependencies
* Read **RED metrics** per service
* **Drill down into all signals** and correlate them, thanks to a new concept: **connection groups**
* Run **incident investigations** across metrics, logs and traces in one place
* Keep findings in **dashboards** and **notebooks**
* Integrate it with favorite AI tooling and agents

<figure style="width:80%; margin:1.5em auto;">
  <img src="/blog/q3-2026-whats-new-victoriametrics-cloud/visum-beta.webp"
       alt="Visum: service dependencies, RED metrics and drill-down into every signal, with beta access one click away"
       style="display:block; width:100%; margin:auto;" />
  <figcaption style="margin:.6em 0 0; text-align:center; font-style:italic; font-size:.9em; opacity:.8;">
    Visum: service dependencies, RED metrics and drill-down into every signal, with beta access one click away
  </figcaption>
</figure>

We would need hours to talk about everything the team is doing in this area, so stay tuned!

> [!WARNING]
> Those wanting to be among the first to try it will find a **Request beta access** button in the console. Click it and we will enable it as soon as it's ready!

In a few words: VictoriaMetrics Cloud `console` is the control plane where users deploy and manage the
`VictoriaStack`, and `Visum` will be the new portal that accesses that data, so everything
will be found there.

## See You in the Next One!

As always, thanks for being part of the VictoriaMetrics community. User feedback, bug reports, and
feature requests are what keep us moving in the right direction.

Ready to try everything we shipped this quarter?
[Sign up for VictoriaMetrics Cloud](https://console.victoriametrics.cloud/signUp?utm_source=blog&utm_medium=blog_post&utm_campaign=q3-2026-whats-new&utm_content=closing-free-trial) and get the
first month free, no credit card required!
