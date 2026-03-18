---
build:
  list: never
  publishResources: false
  render: never
sitemap:
  disable: true
---

VictoriaMetrics Cloud is a DBaaS (Database as a Service) product for VictoriaMetrics.
This means that you need to [create a deployment](https://docs.victoriametrics.com/victoriametrics-cloud/get-started/quickstart/#creating-deployments) before sending data to it.
It is a fully managed service that allows you to deploy and manage your own instances of VictoriaMetrics and VictoriaLogs
in the cloud.

> [!TIP]
> Regular use cases of VictoriaMetrics Cloud are monitoring applications, infrastructure, or services,
running in different setups like on-prem, Private, Public or Hybrid Cloud, Edge devices or IoT.

## Deployment types
At the moment, the following `Deployment types` are available in VictoriaMetrics Cloud:

- `VictoriaMetrics` single-node and cluster instances
- `VictoriaLogs` single-node instances
- **Coming soon**: `VictoriaTraces` instances

## Capacity Tiers
When deploying a new VictoriaMetrics Cloud instance, users need to choose from a set of configurations,
called `Capacity Tiers`, that are tested to serve a certain load and conditions, with a fixed CPU and Memory
allocation.

> [!TIP] DBaaS means Flexibility
> This approach ensures that pay-as-you-go users can rely on transparent and predictable pricing,
> without fearing unexpected cost increases overnight.

- While [VictoriaMetrics Capacity Tiers](https://docs.victoriametrics.com/victoriametrics-cloud/deployments/victoriametrics-tiers/)
for `single-node` instances, provide benchmarking capacities across various magnitudes, such as
Active Time Series, Churn Rate or Ingestion Rate, `cluster` instances are normally tailored for
each use case, to ensure cost savings and performance.
- [VictoriaLogs Capacity Tiers](https://docs.victoriametrics.com/victoriametrics-cloud/deployments/victorialogs-tiers/)
for `single-node` instances are presented in a simpler way, based on CPU and RAM capacities. The
`Cluster` version is not available at the moment.

You can also read more about the [differences between single and cluster instances](https://docs.victoriametrics.com/victoriametrics-cloud/deployments/single-or-cluster/).

### Upgrading between Single-node Capacity tiers
Should your deployment need more CPU/RAM to handle the load, or perform certain activities (such
as ad-hoc heavy investigations) you can always upgrade or downgrade between tiers in the `Settings`
menu of your deployment.

## Limits
Apart from the capacities per tier described in the corresponding pages for
[VictoriaMetrics](https://docs.victoriametrics.com/victoriametrics-cloud/deployments/victoriametrics-tiers/)
and [VictoriaLogs](https://docs.victoriametrics.com/victoriametrics-cloud/deployments/victorialogs-tiers/),
VictoriaMetrics Cloud also enforces hard limits in order to ensure operations. Users can expect
receiving alerts when these thresholds are surpassed, and new data eventually being rejected.
The majority of these limits are available under the `Monitor` tab of the
[deployments](https://console.victoriametrics.cloud/deployments) section for each of your instances.

- Read more about [Limits that apply to deployments](https://docs.victoriametrics.com/victoriametrics-cloud/deployments/limits/)

## More information about deployments
- [How to create a deployment](https://docs.victoriametrics.com/victoriametrics-cloud/get-started/quickstart/#creating-deployments)
- [How to use Access tokens to write and read data](https://docs.victoriametrics.com/victoriametrics-cloud/deployments/access-tokens/)
