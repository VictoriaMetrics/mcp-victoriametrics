---
draft: false
page: "blog blog_post"
authors:
  - "Artur Minchukov"
date: 2026-10-27
title: "Finding Why Payments Fail with VictoriaLogs Drilldown"
summary: "The VictoriaLogs data source for Grafana now ships with Drilldown, a point-and-click way to explore logs without writing queries. In this post, we use it to investigate failing payments in the OpenTelemetry demo shop, going from a single warning to the release, the log message, and the line of code behind it, all without typing a line of LogsQL."
description: "Explore VictoriaLogs in Grafana without writing queries. Learn how to use the new Drilldown feature to find problems in a few clicks and without queries."
enableComments: true
categories:
  - "VictoriaLogs"
  - "Observability"
tags:
  - "VictoriaLogs"
  - "Grafana"
  - "Drilldown"
  - "OpenTelemetry"
  - "Logs"
keywords:
  - "VictoriaLogs Drilldown"
  - "VictoriaLogs Grafana data source"
  - "log patterns"
  - "stream fields"
  - "OpenTelemetry demo logs"
images:
  - /blog/finding-why-payments-fail-with-victorialogs-drilldown/preview.webp
---

Grafana's <a href="https://grafana.com/docs/grafana/latest/visualizations/simplified-exploration/logs/" target="_blank" rel="noopener noreferrer">Logs Drilldown</a> is a point-and-click way to investigate logs without writing queries. Drilldown lets you break the logs down by service, label, pattern, or any other field, and narrow the selection one click at a time until the problem is in view.

It's a powerful tool, but it only works with Loki. So we built the same experience into the VictoriaLogs data source.

We're happy to announce that the VictoriaLogs data source for Grafana now ships with Drilldown integration. It runs entirely inside Grafana Explore, so there's nothing extra to install.

> *Drilldown is available since VictoriaLogs data source [v0.33.1](https://github.com/VictoriaMetrics/victorialogs-datasource/releases/tag/v0.33.1) and VictoriaLogs [v1.33.0](https://github.com/VictoriaMetrics/VictoriaLogs/releases/tag/v1.33.0).*

In this post, we show how it works by taking a payment incident apart with
[VictoriaLogs Drilldown](https://docs.victoriametrics.com/victorialogs/integrations/grafana/#drilldown) without
writing a single query. The application is the
[OpenTelemetry demo](https://github.com/open-telemetry/opentelemetry-demo) shop with its logs sent
to VictoriaLogs, and the incident is staged with one of the demo's built-in feature flags. The
investigation itself is real, and every screenshot comes from that setup. You can try it out yourself on the
[VictoriaMetrics Grafana playground](https://play-grafana.victoriametrics.com/explore?schemaVersion=1&panes=%7B%226hn%22%3A%7B%22datasource%22%3A%22PD775F2863313E6C7%22%2C%22queries%22%3A%5B%7B%22refId%22%3A%22A%22%2C%22datasource%22%3A%7B%22type%22%3A%22victoriametrics-logs-datasource%22%2C%22uid%22%3A%22PD775F2863313E6C7%22%7D%2C%22expr%22%3A%22*%22%2C%22editorMode%22%3A%22code%22%2C%22queryType%22%3A%22instant%22%7D%5D%2C%22range%22%3A%7B%22from%22%3A%22now-1h%22%2C%22to%22%3A%22now%22%7D%2C%22panelsState%22%3A%7B%22logs%22%3A%7B%22sortOrder%22%3A%22Descending%22%2C%22columns%22%3A%5B%22Time%22%2C%22Line%22%5D%2C%22visualisationType%22%3A%22logs%22%2C%22labelFieldName%22%3A%22labels%22%7D%7D%2C%22compact%22%3Afalse%7D%7D&orgId=1). 

## Opening Drilldown

To open Drilldown, go to **Explore**, pick the VictoriaLogs data source, switch the query editor to
**Code** mode and click the **Drilldown** button in the editor toolbar.

![Explore with the VictoriaLogs query editor: switching to Code mode, clicking the Drilldown button, and the drawer opening on the Services tab](/blog/finding-why-payments-fail-with-victorialogs-drilldown/01-entry.webp)
<figcaption style="text-align: center; font-style: italic;">Drilldown opens from the query editor in Explore. The button appears only in Code mode.</figcaption>

The main Drilldown page has tabs, one for each way to break down the logs. **Services**
groups them by service, **Stream fields** by the stream fields the logs are stored with, and
**Patterns** by the pattern of the log message. The **+** button next to the tabs adds a tab for any
other field so that the same table can be built for a namespace, a host, or a status code.

The first tab, **Services**, lists every service by its `service.name` field with a count, a share
of the total, and a small chart per row. A stacked chart above the table shows the same volume over
time. The chart in each row shows how that service's logs are spread over the range, and the color
is the log level: green for info, yellow for warning, red for error, blue for debug. For
OpenTelemetry logs, the level is set up by the
[OpenTelemetry preset](https://docs.victoriametrics.com/victorialogs/integrations/grafana/#opentelemetry-preset)
of the data source. The level buttons above the tabs filter the table by log level.

![Drilldown main view, Services tab, last hour, no filters](/blog/finding-why-payments-fail-with-victorialogs-drilldown/01b-drilldown-open.webp)

The chart looks normal, with no spikes and almost no yellow or red in any row. So we filter the
logs by warnings and errors to see which services have them at all.

## Filtering by level and by service

We click the **error** button above the tabs. The table and the chart now count only error-level
entries, and there are none. We add **warning**. One row is left, `payment`, with 100% of the
warnings in the last hour. The chart shows that they started about half an hour ago, within a
single minute, and have not stopped.

![Services tab: clicking error empties the table, adding warning leaves a single row, payment](/blog/finding-why-payments-fail-with-victorialogs-drilldown/02-levels.webp)

![Services tab with error and warning selected: payment is the only row, with 100% of the warnings, and the chart starts at one specific minute](/blog/finding-why-payments-fail-with-victorialogs-drilldown/02b-error-warning.webp)
<figcaption style="text-align: center; font-style: italic;">Logs with error and warning levels</figcaption>

To see the details of a service, click the **Filter for** icon in its row. This opens the details
view of Drilldown, which shows information only for the selected filter: a volume chart by level,
and tabs with the logs, the stream fields, the other fields, and the patterns of that selection.

![Hovering the payment row and clicking Filter for: the service.name = payment chip appears and the drawer switches to the details view](/blog/finding-why-payments-fail-with-victorialogs-drilldown/03-payment.webp)

![Payment details view with error and warning selected: the volume chart and the Logs tab with the warning lines](/blog/finding-why-payments-fail-with-victorialogs-drilldown/03b-payment-details.webp)
<figcaption style="text-align: center; font-style: italic;">The details view: everything here is about the payment service.</figcaption>

## Stream fields and breakdown by a field

Did the payment service itself change, or is it reacting to something else? We do not know whether
anything was deployed. The level filter hides everything except the bad entries, so we switch
**error** and **warning** off again. The chart fills with the normal traffic, and the warnings sit
on top of it as a thin yellow layer that begins at the same minute.

![Payment details view with no level filter: info traffic for the whole hour and warnings on top from one minute onward](/blog/finding-why-payments-fail-with-victorialogs-drilldown/04-levels-off.webp)
<figcaption style="text-align: center; font-style: italic;">Service <b>payment</b> without level filters</figcaption>

We open the **Stream fields** tab.
[Stream fields](https://docs.victoriametrics.com/victorialogs/keyconcepts/#stream-fields) are the
fields that identify the application instance a log comes from, and VictoriaLogs stores the logs
of every instance as a separate stream. For OpenTelemetry logs, the resource attributes play this
role: `service.name`, `service.namespace`, `service.version`, and the rest. Over the whole hour,
each of them had a single value. The exception is `service.version`, which has two: `2.0.2`
and `2.0.3`.

![Switching error and warning off and opening the Stream fields tab: the service.version card shows two values](/blog/finding-why-payments-fail-with-victorialogs-drilldown/04-stream-fields.webp)

![service.version card with two values: 2.0.2 ends at the moment 2.0.3 starts; the volume chart above shows the warnings beginning at the same minute](/blog/finding-why-payments-fail-with-victorialogs-drilldown/05-stream-fields-version.webp)

We click `service.version` and see that the warnings appeared only in the new version, so we add
that version to the filters.

![Clicking the service.version card title: the tab switches to the breakdown by version](/blog/finding-why-payments-fail-with-victorialogs-drilldown/06-version-field.webp)

![Breakdown by service.version: 2.0.2 with info only, 2.0.3 with info and warning, each with a chart and sample lines](/blog/finding-why-payments-fail-with-victorialogs-drilldown/06-version-breakdown.webp)
<figcaption style="text-align: center; font-style: italic;">The breakdown settles it: the old version never warned, the new one does.</figcaption>

## Patterns

2.0.3 has produced a couple of hundred warning lines. We switch **warning** back on, so the
selection is the new version's warnings and nothing else, and open the **Patterns** tab. It groups
log messages that differ only in numbers, IDs, timestamps, and the like into one row, with a count
and a share per pattern. Here all the warnings collapse into a single row:

```
Payment request failed. Invalid token. app.loyalty.level=gold
```

![In the breakdown: filtering for 2.0.3, switching warning on and opening the Patterns tab](/blog/finding-why-payments-fail-with-victorialogs-drilldown/06-version-patterns.webp)

It covers 100% of the warnings. The count reads `~150`. Patterns are estimated from a
[1-in-10 sample](https://docs.victoriametrics.com/victorialogs/logsql/#sample-pipe) of the logs.
Grouping every line by shape across a large selection would be slow, and a sample is enough to
show which patterns dominate, so the count is approximate. Clicking **Filter for pattern** on a
row adds the pattern to the filters, and the filtered view shows the exact number.

![Patterns tab: a single pattern with a ~ count and 100% share](/blog/finding-why-payments-fail-with-victorialogs-drilldown/07-patterns.webp)
<figcaption style="text-align: center; font-style: italic;">Every warning message has the same pattern. The ~ marks an estimate from a sample.</figcaption>

## Log lines

We found the most common warning pattern. Now let's look at the logs to understand what's causing the failures. We open the
**Logs** tab and click any line to expand it. The expanded entry lists all fields of the log, and
a search box above them narrows the list: typing `err` leaves only the error fields.

![From the Patterns tab to the Logs tab, then clicking a line to expand it](/blog/finding-why-payments-fail-with-victorialogs-drilldown/08-logline.webp)

`err.stack` points at the exact line in the payment service, `charge.js:37`, and the message names
the trait: `app.loyalty.level=gold`. The new version rejects the payment token for customers with
the gold loyalty level, which is why only some users saw the error.

![Logs tab with one entry expanded and the field search set to err: err.message, err.stack with charge.js:37:13, err.type](/blog/finding-why-payments-fail-with-victorialogs-drilldown/08-log-entry.webp)
<figcaption style="text-align: center; font-style: italic;">One entry is enough. The stack points at the code, and the message names the loyalty level.</figcaption>

## Back to the editor

The cause is clear. To continue in Explore, we click **Go to editor**: the drawer closes and the
query editor gets the same filters as chips.

![Explore after Go to editor: the three filters appear as chips in the query editor, the logs volume and log lines match the drawer](/blog/finding-why-payments-fail-with-victorialogs-drilldown/10-go-to-editor.webp)
<figcaption style="text-align: center; font-style: italic;">The same selection is now a query in the editor.</figcaption>

## What we found

Payment version 2.0.3 rejects charge requests from gold-level customers with "Invalid token".
Customers at other levels pay normally. The next move belongs to whoever owns the payment service:
roll back to 2.0.2 and watch the warning count on the same Drilldown screen drop to zero.

We got there without typing a query: the level buttons pointed at the service, the stream fields
showed the deploy, the patterns showed one message, and one expanded log line named the cause.

Drilldown is also enabled on the
[VictoriaMetrics Grafana playground](https://play-grafana.victoriametrics.com/explore?schemaVersion=1&panes=%7B%226hn%22%3A%7B%22datasource%22%3A%22PD775F2863313E6C7%22%2C%22queries%22%3A%5B%7B%22refId%22%3A%22A%22%2C%22datasource%22%3A%7B%22type%22%3A%22victoriametrics-logs-datasource%22%2C%22uid%22%3A%22PD775F2863313E6C7%22%7D%2C%22expr%22%3A%22*%22%2C%22editorMode%22%3A%22code%22%2C%22queryType%22%3A%22instant%22%7D%5D%2C%22range%22%3A%7B%22from%22%3A%22now-1h%22%2C%22to%22%3A%22now%22%7D%2C%22panelsState%22%3A%7B%22logs%22%3A%7B%22sortOrder%22%3A%22Descending%22%2C%22columns%22%3A%5B%22Time%22%2C%22Line%22%5D%2C%22visualisationType%22%3A%22logs%22%2C%22labelFieldName%22%3A%22labels%22%7D%7D%2C%22compact%22%3Afalse%7D%7D&orgId=1):
open Explore, switch to Code mode, and click **Drilldown** to try it on live data without installing
anything. If something gets in the way, open an issue in the
[data source repository](https://github.com/VictoriaMetrics/victorialogs-datasource).
