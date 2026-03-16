---
title: Exploring Data
weight: 10
disableToc: true
menu:
  docs:
    weight: 4
    parent: cloud
    identifier: exploring-data
    pageRef: /victoriametrics-cloud/exploring-data/
tags:
  - metrics
  - logs
  - cloud
  - enterprise
  - guide
---

VictoriaMetrics Cloud helps users to analyze time series or logging data and troubleshoot
queries through the built-in `Explore` utility. This functionality is powered by
[VMUI](https://docs.victoriametrics.com/victoriametrics/single-server-victoriametrics/#vmui)
components and available for both metrics and logs.

You can `Explore` your data accessing to this section in the two following ways:
1. Explore page at [console.victoriametrics.cloud/explore](https://console.victoriametrics.cloud/explore)
1. Per deployment, via a dedicated URL pattern: `console.victoriametrics.cloud/deployment/<DEPLOYMENT_ID>/explore`

## What is VMUI?

[VMUI](https://docs.victoriametrics.com/victoriametrics/single-server-victoriametrics/#vmui) is the
native user interface for VictoriaMetrics, designed to help users explore, troubleshoot and optimize
their queries and metrics. In VictoriaMetrics Cloud, this UI is integrated into the `Explore` view,
offering an accessible toolset to [get instant value from data](https://docs.victoriametrics.com/victoriametrics-cloud/get-started/features/#get-instant-value-from-your-data).

The VMUI documentation may be found [here](https://docs.victoriametrics.com/victoriametrics/single-server-victoriametrics/#vmui),
which is maintained and updated alongside product releases.

The same approach is used for the VictoriaLogs User Interface, sometimes known as `VLUI`, which provides
visualization and analysis capabilities for logs data stored in VictoriaLogs deployments.

![Logs Overview](https://docs.victoriametrics.com/victoriametrics-cloud/exploring-data/explore-logs-overview.webp)
<figcaption style="text-align: center; font-style: italic;">Logs Overview in VictoriaMetrics Cloud</figcaption>

### Playground
The quickest way to discover VMUI is by directly interacting with it. If you are curious, the available
playgrounds for [VictoriaMetrics](https://play.victoriametrics.com/) and [VictoriaLogs](https://play-vmlogs.victoriametrics.com/)
allow you to check real examples of different installations of the VictoriaStack.
There, you can easily testing and learn the query engine or the relabeling debugger among other tools
and pages provided by VMUI.

