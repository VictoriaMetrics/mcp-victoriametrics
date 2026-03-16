---
weight: 1
title: "VictoriaMetrics Deployments"
menu:
  docs:
    parent: "deployments"
    weight: 1
    name: "VictoriaMetrics"
aliases:
  - /victoriametrics-cloud/tiers-parameters/index.html
  - /victoriametrics-cloud/tiers-parameters/
  - /victoriametrics-cloud/deployments/tiers-and-types.html
  - /victoriametrics-cloud/deployments/tiers-and-types/
tags:
  - metrics
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

## Selecting a VictoriaMetrics Capacity Tier

The VictoriaMetrics Cloud team constantly benchmarks different configurations that are
able to serve relevant use cases in the industry, based on our own and users' experience. This process
results in a list of VictoriaMetrics configurations with maximum `Capacities` that are periodically
reviewed and updated.


### Capacities

The parameters that define these Capacity Tiers include:
[active time series](https://docs.victoriametrics.com/victoriametrics/keyconcepts/#time-series),
data read rate, [churn rate](https://docs.victoriametrics.com/victoriametrics/faq/#what-is-high-churn-rate),
ingestion rate or labelset size, and are described in the following table:


| **Capacity**                              | **Description**                                                                                                                                                              |
|-------------------------------------------|------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|
| **Active Time Series**                    | Number of [time series](https://docs.victoriametrics.com/victoriametrics/faq/#what-is-an-active-time-series) that received at least one data point in the last hour.         |
| **Churn Rate**                            | Number of [new active time series in the last 24 hours](https://docs.victoriametrics.com/victoriametrics/faq/#what-is-high-churn-rate).                                      |
| **Ingestion Rate**                        | Number of [time series](https://docs.victoriametrics.com/victoriametrics/keyconcepts/#time-series) samples ingested per second.                                                      |
| **Data Read Rate**                        | Amount of data scanned per time unit. It represents the reading effort that the deployment is doing. While mostly influenced by the querying path, it also accounts for VictoriaMetrics periodic merges and deduplications. |
| **Series Read per Query**                 | Number of series scanned from the database per each query.                                                                                                                        |
| **Equivalent Hosts**                      | Number of hosts that typically generate this amount of load (CPU, Memory, Disk IO...) with a 30s scrape interval. For example, [node exporter](https://prometheus.io/docs/guides/node-exporter/) exposes ~1000 time series per instance. Therefore, if you collect metrics from 50 node exporters, the approximate amount of Active Time Series is 1000 * 50 = 50,000 series. This parameter may be helpful for  planning if you don't know your needed Capacity. |

<br></br>

> [!TIP]
> These metrics are available, in real time, under the `Monitoring Page` for each deployment.

### VictoriaMetrics Cloud Capacity Tiers

VictoriaMetrics Cloud offers the following list of Capacity Tiers:

| **Active Time Series**     | **Churn Rate**  | **Ingestion rate** (datapoints/s) |  **Series Read per query** | **Data read rate** (GiB/min) | **Equivalent Hosts**  |
|----------------------------|-----------------|-----------------------------------|----------------------------|------------------------------|-----------------------|
| 500k                       | 1.25M           | 16.7k                             | 11.4k                       | 0.12                        | 500                   |
| 1M                         | 2.5M            | 33.3k                             | 23.9k                       | 0.2                         | 1000                     |
| 2M                         | 5M              | 66.6k                             | 46.6k                       | 0.3                         | 2000                     |
| 3M                         | 7M              | 100k                              | 49.6k                       | 0.6                         | 3000                     |
| 4M                         | 10M             | 133k                              | 78.4k                       | 0.8                         | 4000                     |
| 5M                         | 12.5M           | 166k                              | 81k                         | 1                           | 5000                     |
| 7.5M                       | 18.8M           | 250k                              | 122k                        | 1.5                         | 7500                     |
| 10M                        | 25M             | 333k                              | 165k                        | 3                           | 10000                     |
| 12.5M                      | 31.3M           | 417k                              | 208k                        | 4                           | 12500                     |
| 15M                        | 37.5M           | 500k                              | 244k                        | 6                           | 15000                     |

### Other Capacities common to all Tiers

| **Capacity**                              | **Value**                                                                                                                                                                      |
|-------------------------------------------|--------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|
| **Scrape Interval**                       | An average scrape interval of 30s was used in these tests to generate the target loads.            |
| **N labels**                              | A typical [node exporter](https://github.com/prometheus/node_exporter) labels set was used with the addition of between 5-10 labels.          |
| **Label value length**                    | The size of the whole time series was, in average, 500B. Extremely increasing this value in average can significantly impact performance and a higher capacity tier may be needed.                                      |
| **Read Requests per second**              | Between 1 and 5, including recording rules and dashboard queries. Extremely increasing these values can significantly impact performance. Cluster versions or intermediate agents may be needed. |
| **Write Requests per second**             | Ingestion Rate Value / 1000. Extremely increasing these values can significantly impact performance. Cluster versions or intermediate agents may be needed. |

<br></br>

> [!TIP]
> All these metrics are available under the `Monitoring Page` for each deployment.

For a detailed explanation of these parameters, visit the guide on [Understanding Your Setup Size](https://docs.victoriametrics.com/guides/understand-your-setup-size/).

### What happens if a parameter in a Capacity Tier is surpassed?

It depends. This model gives users the control to maximize the value from the Cloud service based
on their needs without needing to pay more just because one parameter is exceeded:
* For example, a 1M Active Time Series deployment configuration with low Churn and Read rates may
be perfectly suitable for 1.2M Active Time Series, without needing to upgrade to a higher Tier,
or receiving an unexpected high bill after a sudden increase in data flow.
* On the other hand, excessive churn rates or label sizes may incur on a higher resource consumption,
needing an upgrade to guarantee quality of service, even if the Active Time Series Capacity is not
surpassed.

**If you exceed a capacity, performance may degrade and SLAs are not guaranteed — consider moving to a higher tier.**
> [!TIP]
> Upgrades between Capacity tiers are easily done with a click in `Deployment Settings`.
## Selecting Retention and Storage

The other parameter needed to be set in deployment is the storage. Recommendations are given in
VictoriaMetrics Cloud for the desired retention period.

The [formula used to calculate recommended storage](https://docs.victoriametrics.com/guides/understand-your-setup-size/#retention-perioddisk-space)
depends upon **ingestion rate**, desired **retention** and datapoint size. It is assumed that each data point is 0.8 bytes based
on our own experience with VictoriaMetrics Cloud. However, several factors such as **high cardinality data increases the data point size**,
which may incur in needing more storage than recommended.

### Flexible storage helps to reduce costs and adapt it to your needs.
Keeping in mind that storage can always be increased (but not downsized) **users are recommended to start
small and scale as needed**.

For example, the full amount of storage needed for 6 months retention will only be reached after those 6 months of operation.
You may save costs by avoiding reserving all of that storage from the beginning.

Using [Downsampling](https://docs.victoriametrics.com/victoriametrics/#downsampling),
[retention filters](https://docs.victoriametrics.com/victoriametrics/#retention-filters),
[Data Deduplication](https://docs.victoriametrics.com/victoriametrics/single-server-victoriametrics/#deduplication),
or [Cardinality Explorer](https://docs.victoriametrics.com/victoriametrics/single-server-victoriametrics/#cardinality-explorer)
are encouraged to further reduce your costs.

Feel free to contact [support](mailto:support-cloud@victoriametrics.com) should you need more information or guidance.

## Limits

Apart from these recommendations, VictoriaMetrics Cloud also enforces hard limits in order to
ensure operations. Users can expect receiving alerts when these thresholds are surpassed, and new
metrics eventually being rejected.

Each capacity tier has its own limits, including limits on ingestion rate, concurrent requests,
and the number of access tokens that can be created per deployment.

To request higher limits, contact support through the [Contact Support](https://console.victoriametrics.cloud/contact_support/) button in the Cloud UI.
More information about limits may be found [here](https://docs.victoriametrics.com/victoriametrics-cloud/deployments/limits/).

## Advanced Parameters: Flags

Additionally, VictoriaMetrics Cloud exposes certain parameters (or [command-line flags](https://docs.victoriametrics.com/victoriametrics/single-server-victoriametrics/#list-of-command-line-flags))
that **advanced users** can tweak on their own under the `Advanced settings` section of every deployment
after creation.

> [!WARNING] Changing default command-line flags may lead to errors
> Modifying Advanced parameters can result into changes in resource consumption usage, causing a
> deployment not being able to compute the load they were designed to support. In these cases,
> a higher tier is most probably needed.

Some of these advanced parameters are outlined below:

| **Flag**                               | **Description**                                                             |
|----------------------------------------|-----------------------------------------------------------------------------|
| <nobr>`-maxLabelsPerTimeseries`</nobr> | Maximum number of labels per time series. Time series with excess labels are dropped. Higher values can increase [cardinality](https://docs.victoriametrics.com/victoriametrics/keyconcepts/#cardinality) and resource usage.  |
| `-maxLabelValueLen`                    | Maximum length of label values. Time series with longer values are dropped. Large label values can lead to high RAM consumption. This parameter is not exposed and can only be adjusted via [support](mailto:support-cloud@victoriametrics.com). **In general, label values with high values `~>1kb` are not supported**. |

## Terms and definitions

  - [Time series](https://docs.victoriametrics.com/victoriametrics/keyconcepts/#time-series)
  - [Labels](https://docs.victoriametrics.com/victoriametrics/keyconcepts/#labels)
  - [Active time series](https://docs.victoriametrics.com/victoriametrics/faq/#what-is-an-active-time-series)
  - [Churn rate](https://docs.victoriametrics.com/victoriametrics/faq/#what-is-high-churn-rate)
  - [Cardinality](https://docs.victoriametrics.com/victoriametrics/keyconcepts/#cardinality)
