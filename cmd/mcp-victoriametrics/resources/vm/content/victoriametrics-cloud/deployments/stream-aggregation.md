---
weight: 7
title: "Stream aggregation"
menu:
  docs:
    parent: "victoriametrics-deployments"
    weight: 7
    name: "Stream aggregation"
tags:
  - metrics
  - cloud
---

[Stream aggregation](https://docs.victoriametrics.com/victoriametrics/stream-aggregation/) aggregates
incoming samples on the fly, before they are written to storage, reducing the number of stored series
and samples. VictoriaMetrics Cloud lets you configure it directly from the console, for **single-node
VictoriaMetrics deployments**.

> [!NOTE]
> Stream aggregation is only available for single-node `VictoriaMetrics` deployments.

> [!WARNING]
> Enabling stream aggregation increases CPU and memory usage on the deployment. If the deployment
> runs into resource limits, you may need to upgrade it to a higher tier.

## How to configure stream aggregation

1. Go to the "Settings" tab on the [deployment page](https://console.victoriametrics.cloud/deployments)
   of a single-node deployment.
2. Open the "Advanced settings" section and switch to the "Stream aggregation" tab.
3. Upload a YAML configuration file, or paste/edit it directly in the built-in editor. See the
   [stream aggregation configuration reference](https://docs.victoriametrics.com/victoriametrics/stream-aggregation/#stream-aggregation-config)
   for the full config syntax and available aggregation functions.
4. Optionally, click "Validate configuration" to check the config for syntax errors before saving.
5. Configure the following options as needed:
   - **Keep input** - keep raw samples matching aggregation rules. By default, matched samples are
     aggregated and dropped.
   - **Drop input** - drop raw samples that don't match any aggregation rule. By default, unmatched
     samples are stored as is.
   - **Ignore old samples** - ignore samples with timestamps outside the current aggregation interval.
   - **Ignore first intervals** - number of aggregation intervals to skip after the deployment
     restarts. Increase it if aggregation results are incorrect after restarts due to delayed data.
6. Click "Save" to apply the configuration.

To disable stream aggregation, clear the configuration and save the form.

## Monitoring stream aggregation

Once configured, stream aggregation metrics are available on the "Monitoring" tab of the
[deployment page](https://console.victoriametrics.cloud/deployments), under the "Stream aggregation"
section: matched, produced and dropped samples, flush timeouts and duration, samples lag and memory
usage. See [stream aggregation troubleshooting](https://docs.victoriametrics.com/victoriametrics/stream-aggregation/#troubleshooting)
for help interpreting these metrics.
