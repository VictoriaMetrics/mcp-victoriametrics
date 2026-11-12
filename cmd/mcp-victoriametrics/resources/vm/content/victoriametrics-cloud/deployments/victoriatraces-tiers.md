---
weight: 3
title: "VictoriaTraces Deployments"
menu:
  docs:
    parent: "deployments"
    weight: 3
    name: "VictoriaTraces"
tags:
  - traces
  - cloud
  - enterprise
---

VictoriaMetrics Cloud users [launch](https://docs.victoriametrics.com/victoriametrics-cloud/get-started/quickstart/#creating-deployments)
and interact with dedicated instances managed by our team. This section provides relevant information
to help choosing which **single-node** `Capacity Tier` is the best for your needs.

> [!TIP]
> Keep in mind that you can always move to a `bigger` or `smaller` single-node instance!
> Sometimes you may need more CPU/RAM to handle load permanently, or perform certain activities
> (such as ad-hoc heavy investigations). In these cases you can always upgrade or downgrade
> between single-node tiers in the `Settings` menu of your deployment.

## Selecting a VictoriaTraces Capacity Tier

VictoriaTraces Capacity Tiers are created based on computing resources (CPU and RAM). Other
parameters, such as span ingestion or reading rate are shown for users as informational magnitudes,
and most of the limits shown in this document and in the application should be treated as **soft
limits**. This means that if your experience with a given instance is satisfactory, we generally
don't force users to upgrade.

{{% collapse name="Why are tiers only based on CPU and Memory?" %}}
Trace data varies a lot between applications. The number of spans per trace, the number of resource
and span attributes, and how much of that data is high cardinality all differ per service and
instrumentation library. This produces a complexity hard to model, which affects both compression
(and subsequently storage) together with how computing resources are used.

However, our premise is that by understanding the CPU and Memory used/available, users can know
how their specific load and usage affects to the resource consumption and mindfully adapt,
minimizing costs.

{{% /collapse %}}

>[!TIP]
> Not sure about which tier to use? [Contact Us](https://console.victoriametrics.cloud/contact_support/)
for a dedicated PoC. The process takes up to a week.

### Capacities

To help users decide which instance they should start with, our team has performed a benchmarking
exercise that may represent some use cases. The tested parameters and results are shown in this
document as reference, but **changes are expected as we learn more about usage**. However, the
**main parameters that define VictoriaTraces Capacity Tiers are: number of CPUs and GBs of Memory**
allocated per tier.

| **Capacity**                              | **Description**                  |
|-------------------------------------------|-----------------------------------|
| **Span Ingestion Rate**                   | Number of trace spans ingested per second.|
| **Data Ingestion Rate**                   | Amount of data ingested per second.       |
| **New Streams over 24h**                  | Number of new [stream fields in the last 24 hours](https://docs.victoriametrics.com/victoriatraces/keyconcepts/#stream-fields). In VictoriaTraces, `service.name` and the span `name` are used as stream fields.|
| **Data Read Rate**                        | Amount of data scanned per time unit. It represents the reading effort that the deployment is doing. While mostly influenced by the querying path, it also accounts for VictoriaTraces periodic merges and deduplications.|
| **Data Read per Query**                   | Amount of data scanned from the database per each query.       |

<br></br>

> [!TIP]
> These metrics are available, in real time, under the `Monitoring Page` for each deployment.

### VictoriaTraces Capacity Tiers

VictoriaTraces tiers in VictoriaMetrics Cloud were suited to, at least, handle the load in tests done
under the following conditions:

> [!TIP]
> The numbers in both tables below were measured **while ingesting and querying at the same time**.
> If your workload is ingestion-heavy but queries infrequently (or against small time ranges), your
> deployment may comfortably handle a higher ingestion rate than the table suggests, since it isn't
> competing with querying for CPU and memory. See
> [What happens if a parameter in a Capacity Tier is surpassed?](https://docs.victoriametrics.com/victoriametrics-cloud/deployments/victoriatraces-tiers/#what-happens-if-a-parameter-in-a-capacity-tier-is-surpassed)
> below.

#### Ingestion

| CPUs | RAM (GB) | Span Ingestion Rate (spans/s) | Data Ingestion Rate (MB/s) | New Streams over 24h |
|------|----------|-------------------------------|----------------------------|----------------------|
| 1    | 4        | 800                           | 3.5                        | 150000               |
| 2    | 8        | 1600                          | 7.0                        | 200000               |
| 3    | 12       | 2400                          | 10.5                       | 250000               |
| 4    | 16       | 3200                          | 14.0                       | 300000               |
| 5    | 20       | 4000                          | 17.5                       | 350000               |

#### Querying

| CPUs | RAM (GB) | Data Read Rate (GB/s) | Data Read per query (GB) |
|------|----------|-----------------------|--------------------------|
| 1    | 4        | 4.68                  | 4.0                      |
| 2    | 8        | 11.49                 | 6.0                      |
| 3    | 12       | 16.37                 | 9.2                      |
| 4    | 16       | 19.64                 | 11.9                     |
| 5    | 20       | 21.9                  | 14.9                     |

### What happens if a parameter in a Capacity Tier is surpassed?

It depends. This model gives users the control to maximize the value from the Cloud service based
on their needs without needing to pay more just because one parameter is exceeded.

Depending on every specific usage, the numbers shown may strongly vary. Trace workloads are bursty
by nature: a deployment may stay far below its ingestion capacity most of the time and peak during
incidents, when querying also increases. Since ingestion and querying compete for the same CPU and
memory, you may get very high ingestion peaks without needing to upgrade (and paying more) when
keeping querying low.

**If you exceed a capacity, performance may degrade and SLAs are not guaranteed - consider moving to a higher tier.**
> [!TIP]
> Upgrades between Capacity tiers are easily done with a click in `Deployment Settings`.

## Selecting Retention and Storage

The other parameter needed to be set in deployment is the storage.

VictoriaTraces stores ingested data in per-day partitions, and old partitions are removed according
to the configured retention. Read more about how VictoriaTraces handles
[retention](https://docs.victoriametrics.com/victoriatraces/#retention) and
[partitions](https://docs.victoriametrics.com/victoriatraces/#partitions-lifecycle).

### Flexible storage helps to reduce costs and adapt it to your needs.
Keeping in mind that storage can always be increased (but not downsized) **users are recommended to start
small and scale as needed**.

For example, the full amount of storage needed for 6 months retention will only be reached after those 6 months of operation.
You may save costs by avoiding reserving all of that storage from the beginning.

Feel free to contact [support](mailto:support-cloud@victoriametrics.com) should you need more information or guidance.

## Limits
VictoriaMetrics Cloud also enforces hard limits in order to ensure operations. Users can expect
receiving alerts when these thresholds are surpassed, and new data eventually being rejected. Each
capacity tier has its own limits.

At the moment, VictoriaTraces is not enforcing limits apart from Maximum Access Tokens and the
number of concurrent requests per Access Token. This may change in the future.

| CPUs | Access token limit | Concurrent requests per token |
|------|--------------------|-------------------------------|
| 1    | 20                 | 32                            |
| 2    | 40                 | 48                            |
| 3    | 60                 | 64                            |
| 4    | 80                 | 72                            |
| 5    | 100                | 98                            |

To request higher limits, contact support through the [Contact Support](https://console.victoriametrics.cloud/contact_support/) button in the Cloud UI.
More information about limits may be found [here](https://docs.victoriametrics.com/victoriametrics-cloud/deployments/limits/).

## Sending and reading trace data

VictoriaTraces deployments in VictoriaMetrics Cloud accept trace spans in the
[OpenTelemetry protocol (OTLP)](https://opentelemetry.io/docs/specs/otlp/) **over HTTP**, and expose
both the Tempo and the Jaeger querying APIs:

| Operation | Endpoint |
|-----------|----------|
| Write spans (OTLP/HTTP) | `<DEPLOYMENT_ENDPOINT_URL>/insert/opentelemetry/v1/traces` |
| Read (Tempo API) | `<DEPLOYMENT_ENDPOINT_URL>/select/tempo` |
| Read (Jaeger API) | `<DEPLOYMENT_ENDPOINT_URL>/select/jaeger` |

> [!NOTE]
> The OTLP gRPC endpoint (port `4317`) described in the
> [VictoriaTraces documentation](https://docs.victoriametrics.com/victoriatraces/data-ingestion/opentelemetry/)
> is not exposed in VictoriaMetrics Cloud. Configure your collectors and SDKs to use the OTLP/HTTP
> exporter.

Both read endpoints require an Access Token with **read** access, and the write endpoint requires one
with **write** access. See [Access Tokens](https://docs.victoriametrics.com/victoriametrics-cloud/deployments/access-tokens/).

## Advanced Parameters: Flags

Additionally, VictoriaMetrics Cloud exposes certain parameters (or [command-line flags](https://docs.victoriametrics.com/victoriatraces/#list-of-command-line-flags))
that **advanced users** can tweak on their own under the `Advanced settings` section of every deployment
after creation.

> [!WARNING] Changing default command-line flags may lead to errors
> Modifying Advanced parameters can result into changes in resource consumption usage.
