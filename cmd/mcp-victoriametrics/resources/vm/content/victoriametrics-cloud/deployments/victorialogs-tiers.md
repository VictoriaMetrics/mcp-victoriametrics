---
weight: 2
title: "VictoriaLogs Deployments"
menu:
  docs:
    parent: "deployments"
    weight: 2
    name: "VictoriaLogs"
tags:
  - logs
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

## Selecting a VictoriaLogs Capacity Tier

VictoriaLogs Capacity Tiers are created based on computing resources (CPU and RAM). Other
parameters, such as log ingestion or reading rate are shown for users as informational magnitudes,
and most of the limits shown in this document and in the application should be treated as **soft
limits**. This means that if your experience with a given instance is satisfactory, we generally
don't force users to upgrade.

{{% collapse name="Why are tiers only based on CPU and Memory?" %}}
Logs are typically less standardized than other signals. Even in the same organization, some logs
will be generated without structure, while others may be wide-events or using an OpenTelemetry
format. This produces a complexity hard to model, which affects both compression (and subsequently
storage) together with how computing resources are used.

However, our premise is that by understanding the CPU and Memory used/available, users can known
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
**main parameters that define VictoriaLogs Capacity Tiers are: number of CPUs and GBs of Memory**
allocated per tier.

| **Capacity**                              | **Description**                  |
|-------------------------------------------|-----------------------------------|
| **Log Ingestion Rate**                    | Number of log samples ingested per second.|
| **Data Ingestion Rate**                   | Amount of data ingested per second.       |
| **New Streams over 24h**                  | Number of new [stream fields in the last 24 hours](https://docs.victoriametrics.com/victorialogs/keyconcepts/#stream-fields).                                      |
| **Data Read Rate**                        | Amount of data scanned per time unit. It represents the reading effort that the deployment is doing. While mostly influenced by the querying path, it also accounts for VictoriaLogs periodic merges and deduplications.|
| **Data Read per Query**                   | Amount of data scanned from the database per each query.       |

<br></br>

> [!TIP]
> These metrics are available, in real time, under the `Monitoring Page` for each deployment.


### VictoriaLogs Capacity Tiers

VictoriaLogs tiers in VictoriaMetrics Cloud were suited to, at least, handle the load in tests done
under the following conditions:

#### Ingestion

| CPUs | RAM (GB) | Log Ingestion Rate (logs/s) | Data Ingestion Rate (MB/s) | New Streams over 24h |
|------|----------|-----------------------------|----------------------------|----------------------|
| 1    | 3.8      | 800                         | 0.8                        | 50                   |
| 2    | 7.6      | 1200                        | 1.2                        | 60                   |
| 3    | 11.4     | 2200                        | 2.0                        | 60                   |
| 4    | 15.2     | 3200                        | 3.2                        | 80                   |
| 5    | 19.0     | 5000                        | 4.5                        | 80                   |

#### Querying

| CPUs | RAM (GB) | Data Read Rate (GB/s) | Data Read per query (GB) |
|------|----------|-----------------------|--------------------------|
| 1    | 3.8      | 4.6                   | 4                        |
| 2    | 7.6      | 11.49                 | 6                        |
| 3    | 11.4     | 16.37                 | 9.2                      |
| 4    | 15.2     | 19.64                 | 11.9                     |
| 5    | 19.0     | 21.9                  | 14.9                     |

### What happens if a parameter in a Capacity Tier is surpassed?

It depends. This model gives users the control to maximize the value from the Cloud service based
on their needs without needing to pay more just because one parameter is exceeded.

Depending on every specific usage, the numbers shown may strongly vary. For example,
In [this benchmark](https://www.truefoundry.com/blog/victorialogs-vs-loki), a 2 CPU VictoriaLogs
instance was able to cope with 66 MB/s. While the **test makes a lot of sense for comparison
purposes**, real cases are more complex, and querying at the same time would most likely have
shown a lower maximum ingestion peak for all databases. But you could get very high ingestion
peaks without needing to upgrade (and paying more) when keeping querying low.

**If you exceed a capacity, performance may degrade and SLAs are not guaranteed — consider moving to a higher tier.**
> [!TIP]
> Upgrades between Capacity tiers are easily done with a click in `Deployment Settings`.

## Selecting Retention and Storage

The other parameter needed to be set in deployment is the storage.

Read more about how to estimate the needed compute resources for a workload
[here](https://docs.victoriametrics.com/victorialogs/faq/#how-to-estimate-the-needed-compute-resources-for-the-given-workload).

### Flexible storage helps to reduce costs and adapt it to your needs.
Keeping in mind that storage can always be increased (but not downsized) **users are recommended to start
small and scale as needed**.

For example, the full amount of storage needed for 6 months retention will only be reached after those 6 months of operation.
You may save costs by avoiding reserving all of that storage from the beginning.

Feel free to contact [support](mailto:support-cloud@victoriametrics.com) should you need more information or guidance.

## Limits
VictoriaMetrics Cloud also enforces hard limits in order to
ensure operations. Users can expect receiving alerts when these thresholds are surpassed, and new
metrics eventually being rejected. Each capacity tier has its own limits.

At the moment, VictoriaLogs is not enforcing limits apart from Maximum Access Tokens. This may
change in the future.

To request higher limits, contact support through the [Contact Support](https://console.victoriametrics.cloud/contact_support/) button in the Cloud UI.
More information about limits may be found [here](https://docs.victoriametrics.com/victoriametrics-cloud/deployments/limits/).

## Advanced Parameters: Flags

Additionally, VictoriaMetrics Cloud exposes certain parameters (or [command-line flags](https://docs.victoriametrics.com/victorialogs/#list-of-command-line-flags))
that **advanced users** can tweak on their own under the `Advanced settings` section of every deployment
after creation.

> [!WARNING] Changing default command-line flags may lead to errors
> Modifying Advanced parameters can result into changes in resource consumption usage.
