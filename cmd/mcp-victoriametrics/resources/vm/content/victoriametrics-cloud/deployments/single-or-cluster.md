---
weight: 3
title: "Single or Cluster?"
menu:
  docs:
    parent: "deployments"
    weight: 3
    name: "Single or Cluster?"
tags:
  - metrics
  - cloud
  - enterprise
---

VictoriaMetrics databases offer two different deployment types: **single-node** and **cluster**.

At the moment, the following `Deployment types` are available in VictoriaMetrics Cloud:
- `VictoriaMetrics` single-node and cluster
- `VictoriaLogs` single-node
- **Coming soon**: `VictoriaTraces`


All deployments in VictoriaMetrics Cloud are based on the VictoriaMetrics [Open Source projects](https://github.com/VictoriaMetrics/VictoriaMetrics/),
and managed by the VictoriaMetrics team.

> [!TIP]
> VictoriaLogs single-node instances support multitenancy

## VictoriaMetrics single-node vs cluster versions

In a nutshell, [Single-node deployments](https://docs.victoriametrics.com/victoriametrics/single-server-victoriametrics/)
are useful for affordable and performant instances, while [Cluster deployments](https://docs.victoriametrics.com/victoriametrics/cluster-victoriametrics/)
may be the ideal choice for those use cases that require higher availability and multi-tenancy at scale.
More detailed information about the general capabilities of both tiers can be found in this [FAQ](https://docs.victoriametrics.com/victoriametrics/faq/#which-victoriametrics-type-is-recommended-for-use-in-production---single-node-or-cluster).

For simplicity, Capacity tiers based on the VictoriaMetrics Cluster version are hidden to users,
as Single-node instances cover most use cases, and the Cluster version consumes more resources,
which leads to a higher price. In most cases, the cluster version requires tailoring specific
topologies, by dedicating resources where needed (i.e. reading or storage paths).

> [!TIP]
> If you are considering using a VictoriaMetrics Cluster version, [contact-us](https://victoriametrics.com/contact-us/).

More in detail, the following topics should be considered when selecting a deployment type:

{{% collapse name="Reliability/SLA" %}}

Both instance types are highly reliable, with SLAs of 99.5% for `Single-node` deployments and 99.95%
for `Cluster` deployments.

{{% /collapse %}}

{{% collapse name="High Availability" %}}

Since `Single-node` deployments are just one instance, they cannot be highly available. In practice,
this means that during configuration changes and software upgrades, your deployment will experience
a few minutes downtime. (This period of unavailability is not included in the SLA).

On the other hand, `Cluster` deployments do not experience such downtimes.

{{% /collapse %}}

{{% collapse name="Multitenancy" %}}

While [Multitenancy](https://docs.victoriametrics.com/victoriametrics/cluster-victoriametrics/#multitenancy)
is supported in the `Cluster` version of VictoriaMetrics Cloud, it is not supported in `Single-node`
instances.

{{% /collapse %}}

{{% collapse name="Scalability" %}}

Internally, `Single-node` deployments may be scaled vertically and `Cluster` deployments horizontally.

In practice, for VictoriaMetrics Cloud tiers, this means that vertical scaling will affect by
constraining some parameters such as the maximum storage size, but horizontal scaling has no such
limitations.

{{% /collapse %}}

{{% collapse name="Data Replication" %}}

Data replication is provided for `Cluster` deployments only. `Single-node` deployments do not have
such capabilities.

{{% /collapse %}}

{{% collapse name="Enterprise features" %}}

[Enterprise features](http://docs.victoriametrics.com/victoriametrics/enterprise/#victoriametrics-enterprise-features)
are available in both `Single-node` and `Cluster` versions. Some of them may take a while to be exposed
in VictoriaMetrics Cloud. If you are missing any feature, or have any request don't hesitate to
contact us at contact us at support-cloud@victoriametrics.com.

{{% /collapse %}}

{{% collapse name="Efficiency and performance" %}}

Both `Single-node` and `Cluster` versions are highly valued for their performance in various benchmarks
and use cases in the industry. Feel free to read more about use cases and articles [here](http://docs.victoriametrics.com/victoriametrics/articles/).

{{% /collapse %}}
