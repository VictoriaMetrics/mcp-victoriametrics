---
weight: 1
title: Grafana
menu:
  docs:
    identifier: integrations-vl-grafana
    parent: integrations-vl
    weight: 1
---


The [VictoriaLogs Grafana plugin](https://grafana.com/grafana/plugins/victoriametrics-logs-datasource/) allows Grafana 
to query, visualize, and interact with [VictoriaLogs](https://docs.victoriametrics.com/victorialogs/),
a high-performance log storage and processing system.

<img alt="Grafana Dashboard Screenshot" width="100%" src="https://github.com/VictoriaMetrics/victorialogs-datasource/blob/main/src/img/dashboard.png?raw=true">

## Capabilities

1. Use [LogsQL](https://docs.victoriametrics.com/victorialogs/logsql/) to filter, aggregate, and transform logs data to gain insights into application behavior.
1. Use Explore mode with Grafana.
1. Show live-streaming logs.
1. Build dashboards and setup alerts.
1. Use Ad Hoc filters.

<!-- vlcompat:min-version:start -->
The plugin requires VictoriaLogs `v1.54.0` or newer: see [VictoriaLogs compatibility](https://github.com/VictoriaMetrics/victorialogs-datasource/tree/main#victorialogs-compatibility) for what each feature needs and how the plugin reports an unsupported version.
<!-- vlcompat:min-version:end -->

Try it at [VictoriaMetrics playground](https://play-grafana.victoriametrics.com/d/be5zidev72m80f/k8s-logs-demo)!

## Installation

For detailed instructions on how to install the plugin on Grafana Cloud or locally, please checkout the [Plugin installation docs](https://grafana.com/docs/grafana/latest/plugins/installation/).
For installation options in Docker or Kubernetes refer to [these docs](https://github.com/VictoriaMetrics/victorialogs-datasource?tab=readme-ov-file#installation).

## Multitenancy configuration

VictoriaLogs datasource supports [multitenancy](https://docs.victoriametrics.com/victorialogs/#multitenancy) via `multitenancyHeaders` datasource configuration option. 
To configure it, set AccountID and ProjectID in the `multitenancyHeaders` section.

<img alt="Multitenancy configuration" src="https://github.com/VictoriaMetrics/victorialogs-datasource/blob/main/src/img/multitenancy_configuration.png?raw=true">

### Manual configuration via UI

Once the plugin is installed on your Grafana instance, follow [these instructions](https://grafana.com/docs/grafana/latest/datasources/add-a-data-source/)
to add a new VictoriaLogs data source, and enter configuration options.

### Configuration via file

Provisioning of Grafana plugin requires creating [datasource config file](http://docs.grafana.org/administration/provisioning/#datasources):

```yaml
apiVersion: 1
datasources:
  - name: VictoriaLogs
    type: victoriametrics-logs-datasource
    access: proxy
    url: http://victorialogs:9428
    isDefault: true
    jsonData:
      # Multitenancy settings, see https://docs.victoriametrics.com/victorialogs/#multitenancy
      # to use the multitenancy, uncomment lines below: AccountID and ProjectID
      multitenancyHeaders:
        #AccountID: 0
        #ProjectID: 0
```

## Authentication

The datasource supports the authentication methods provided by Grafana's standard HTTP settings
(Basic auth, TLS client certificates, custom HTTP headers), plus forwarding the OAuth identity
of the signed-in Grafana user.

### Forward OAuth Identity

If Grafana itself authenticates users via OAuth/OIDC (for example, [generic OAuth](https://grafana.com/docs/grafana/latest/setup-grafana/configure-security/configure-authentication/generic-oauth/)
with Keycloak, Okta, Auth0 or another provider), the datasource can reuse the access token of the logged-in user.
With **Forward OAuth Identity** enabled, Grafana adds the user's `Authorization: Bearer <access token>` header
(and `X-Id-Token`, when available) to every request the plugin sends to the configured URL:
data queries, health check ("Save & test") and autocomplete/resource requests.

To enable it in the UI, turn on the **Forward OAuth Identity** toggle in the Auth section of the datasource settings.
To enable it via provisioning, set `jsonData.oauthPassThru: true`:

```yaml
apiVersion: 1
datasources:
  - name: VictoriaLogs
    type: victoriametrics-logs-datasource
    access: proxy
    # vmauth validates the forwarded JWT and routes the request to VictoriaLogs
    url: http://vmauth:8427
    jsonData:
      oauthPassThru: true
```

This pairs naturally with [vmauth](https://docs.victoriametrics.com/victoriametrics/vmauth/) in front of VictoriaLogs:
vmauth can verify the forwarded JWT via [OIDC discovery](https://docs.victoriametrics.com/victoriametrics/vmauth/#oidc-discovery)
and authorize/route requests based on token claims via [JWT claim matching](https://docs.victoriametrics.com/victoriametrics/vmauth/#jwt-claim-matching):

```yaml
# vmauth -auth.config
users:
  - jwt:
      oidc:
        issuer: "https://sso.example.com/realms/main"
      match_claims:
        # authorize and route requests based on any JWT claims
        team: "observability"
      # Required for tokens issued by general-purpose identity providers
      # (Keycloak, Okta, Auth0, ...) which don't embed the VictoriaMetrics-specific
      # `vm_access` claim, see the note below.
      default_vm_access_claim: {}
    url_prefix: "http://victorialogs:9428/"
```

**Note about the `vm_access` claim:** since vmauth `v1.147.0`, a JWT token without the `vm_access` claim
falls through to `unauthorized_user` and is rejected unless `default_vm_access_claim` is configured
in the `jwt` section. Tokens issued by general-purpose identity providers normally don't contain `vm_access`,
so either set `default_vm_access_claim` (as in the example above) or configure your identity provider
to embed the [`vm_access` claim](https://docs.victoriametrics.com/victoriametrics/vmauth/#jwt-claim-based-request-templating).
The rejection is silent by default — run vmauth with `-logInvalidAuthTokens` to see the reason.

**Limitations:** the forwarded token exists only within the session of a signed-in user.
Features that run without a user context — alerting and recording rules, public dashboards, reporting —
send requests without the token and will be rejected by the auth proxy.
For those, use a service credential (for example, Basic auth or a bearer token on a separate datasource).

## Building queries

VictoriaLogs query language is [LogsQL](https://docs.victoriametrics.com/victorialogs/logsql/).
Queries can be built using raw LogsQL or via QueryBuilder.

See panels examples at [VictoriaMetrics playground](https://play-grafana.victoriametrics.com/d/be5zidev72m80f/k8s-logs-demo)
and LogsQL examples [here](https://docs.victoriametrics.com/victorialogs/logsql-examples/).

### Query builder

> **Note:** available since [v0.27.1](https://github.com/VictoriaMetrics/victorialogs-datasource/releases/tag/v0.27.1).

A visual alternative to the LogsQL code editor: build queries by chaining typed pipes through an interactive UI with field/value dropdowns and keyboard navigation.

<img alt="Query builder" width="100%" src="https://github.com/VictoriaMetrics/victorialogs-datasource/blob/main/src/img/query_builder.gif?raw=true">


#### Switching between Builder and Code

Use the **Builder** / **Code** toggle in the query editor header to switch modes.

* **Builder → Code** is always safe — the current builder state is serialized to LogsQL.
* **Code → Builder** shows a confirmation dialog and discards the manual query and any stream filters. Automatic conversion of raw LogsQL into builder pipes is not supported yet — it is planned for a future release.

#### Available pipe types

* **Filter**:
  * **exact** — `field:in(values)`. Stream fields are available in the field picker (under the "Stream fields" group) — selecting one converts the filter into a stream filter automatically.
  * **phrase** — substring match: `field:value`
  * **range** — numeric comparison: `>`, `>=`, `<`, `<=`
  * **regexp** — regex match: `~`, `!~`
* **Modify** — `rename`, `copy`, `delete`, `keep`, `replace`, `replace_regexp`, `extract`.
* **Aggregate** — `count`, `sum`, `avg`, `min`, `max`, `median`, `quantile`, `rate`, `histogram`. Each supports optional `by`, `if`, and `as` clauses.
* **Sort** — `sort by field asc|desc`, with optional `limit`, `offset`, `partition by`, `rank as`.
* **Limit** — `limit`, `offset`, `first`, `last`, `top`.
* **Custom** — raw LogsQL escape hatch for syntax not covered by the structured pipes (not validated).

See the [LogsQL pipes reference](https://docs.victoriametrics.com/victorialogs/logsql/#pipes) for the underlying semantics.

#### Adding and editing pipes

* Click the **+** button (Add pipe) to open the pipe picker. Pipes are grouped by category and searchable by name.
* Hover the **|** separator between pipes to insert a new pipe at that position.
* Click any placeholder chip to edit it. Field-name and field-value chips fetch live suggestions from VictoriaLogs; stream fields are surfaced first under a separate group as the faster lookup.
* Multi-value chips (e.g., `field:in(...)`) accept comma-separated values rendered as inline pills.
* Each pipe has a **✕** button to delete it and a **⋯** button to add optional clauses available for that pipe type (e.g., `by`, `if`, `as` on aggregations).
* The trash icon at the bottom-right clears all pipes via the **Clear query** action.

#### Keyboard navigation

| Key | Action |
|-----|--------|
| **Tab** | Confirm the current segment and move to the next one. After the last segment, opens the Add pipe menu. |
| **Enter** | Pick the highlighted dropdown option. |
| **↑ / ↓** | Navigate dropdown options. |
| **Escape** | Cancel the current edit and revert the value. |
| **Backspace** | In an empty multi-value chip — remove the last entered value. |
| **Shift+Enter** | Run the query. |

### Logs panel

For using [Logs panel](https://grafana.com/docs/grafana/latest/panels-visualizations/visualizations/logs/)
switch to `Raw Logs` query type:

<img alt="Logs panel" src="https://github.com/VictoriaMetrics/victorialogs-datasource/blob/main/src/img/panel_logs.png?raw=true">

### Time series panel

For using [Time series panel](https://grafana.com/docs/grafana/latest/panels-visualizations/visualizations/time-series/)
switch to `Range` query type:

<img alt="Time series panel" src="https://github.com/VictoriaMetrics/victorialogs-datasource/blob/main/src/img/panel_time_series.png?raw=true">

### Stats panel

For using [Stats panel](https://grafana.com/docs/grafana/latest/panels-visualizations/visualizations/stat/)
switch to `Instant` query type:

<img alt="Stats panel" src="https://github.com/VictoriaMetrics/victorialogs-datasource/blob/main/src/img/panel_stat.png?raw=true">

For enabling background visualization switch to `Range` query type.

### Table panel

For using [Table panel](https://grafana.com/docs/grafana/latest/panels-visualizations/visualizations/table/)
switch to `Raw Logs` query type:

<img alt="Table panel" src="https://github.com/VictoriaMetrics/victorialogs-datasource/blob/main/src/img/panel_table.png?raw=true">

And apply `Transformations` by labels:

<img alt="Transformations" src="https://github.com/VictoriaMetrics/victorialogs-datasource/blob/main/src/img/panel_table_transformation.png?raw=true">

### Heatmap panel
For using [Heatmap panel](https://grafana.com/docs/grafana/latest/panels-visualizations/visualizations/heatmap/)
switch to `Range` query type and use [`histogram` stats function](https://docs.victoriametrics.com/victorialogs/logsql/#histogram-stats). Use the step to adjust the granularity of the heatmap.
Query example:
```logsql
{collector="otel-collector", k8s.namespace.name="play-otel", service.name="checkout"} "duration:" 
  | extract "duration: <duration>" 
  | stats histogram(duration)
```
<img alt="Heatmap panel" width="100%" src="https://github.com/VictoriaMetrics/victorialogs-datasource/blob/main/src/img/panel_heatmap.png?raw=true">

### Log level rules

The **Log level rules** section in the datasource configuration allows you to assign log levels based on custom field conditions. This helps classify logs dynamically (e.g., as `error`, `info`, `debug`, etc.) using rules you define.

#### How to use

1. Open the datasource settings.

2. Scroll to the **Log level rules** section.

3. Click **"Add rule"** to define a new rule.

4. For each rule, configure the following:

    * **Enable switch** – enable or disable the rule.
    * **Field name** – the log field the condition will evaluate.
    * **Operator** – choose from: `Equals`, `Not equal`, `Matches regex`, `Less than`, `Greater than`, `Word filter`
    * **Value** – the value to compare the field against.
    * **Log level** – level to assign if the condition matches: `critical`, `warning`, `error`, `info`, `debug`, `trace`, `unknown`
    * **Delete button** – remove the rule.

5.  After adding or editing rules, click **"Save & test"** to apply the changes.

**Rule priority**: If multiple rules match a log entry, the **first matching rule** (top to bottom) takes precedence.

6. To define rules via the provision file, use the following format of the provision file:

```yaml
apiVersion: 1
datasources:
  # <string, required> Name of the VictoriaLogs datasource
  # displayed in Grafana panels and queries.
  - name: VictoriaLogs
    # <string, required> Sets the data source type.
    type: victoriametrics-logs-datasource
    # <string, required> Sets the access mode, either
    # proxy or direct (Server or Browser in the UI).
    access: proxy
    # <string> Sets URL for sending queries to VictoriaLogs server.
    # see https://docs.victoriametrics.com/victorialogs/querying/
    url: https://play-vmlogs.victoriametrics.com
    # <string> Sets the pre-selected datasource for new panels.
    # You can set only one default data source per organization.
    isDefault: true
    jsonData:
      logLevelRules:
        - field: "_stream_id"
          value: "123123"
          level: "error"
          operator: "regex"
          enabled: true
```
Where:
* `field` is the name of the log field to evaluate (e.g. `_stream_id`, `status_code`, `message`).
* `value` is the value to compare against. Can be a string or a number, depending on the field.
* `level` is the log level to assign if the condition matches. Valid values: `critical`, `error`, `warning`, `info`, `debug`, `trace`, `unknown`.
* `operator` is the comparison operator to use, such as `equals`, `notEquals`, `regex`, `lessThan`, `greaterThan`, `wordFilter`. The `wordFilter` operator follows [LogsQL word filter](https://docs.victoriametrics.com/victorialogs/logsql/#word-filter) semantics: the value matches whole tokens, not substrings (e.g. `error` matches `an error happened` but not `terror`).
* `enabled` is a boolean flag to enable or disable the rule. Defaults to `true` if omitted.

### Variables
VictoriaLogs datasource supports [variables](https://grafana.com/docs/grafana/latest/variables/) in queries.

#### Ad Hoc filter
You can use [Ad Hoc filters](https://grafana.com/docs/grafana/latest/visualizations/dashboards/variables/add-template-variables/#add-ad-hoc-filters) to filter logs in Dashboards.
By default, ad-hoc filters are applied to all panels in the Dashboard as the [extra_filters](https://docs.victoriametrics.com/victorialogs/querying/#extra-filters) query parameter. If you navigate to the Explore page from the Dashboard, ad-hoc filters are also applied there.
<img alt="Ad hoc filters" src="https://github.com/VictoriaMetrics/victorialogs-datasource/blob/main/src/img/extra_filters_explore_page.png?raw=true">

##### Ad-hoc filters mode

> **Note:** available since [v0.28.0](https://github.com/VictoriaMetrics/victorialogs-datasource/releases/tag/v0.28.0).

Each panel exposes an **Ad-hoc filters** selector in the Query Editor *Options* section that controls how ad-hoc values reach the query:

* **Extra Filters** — *default*. Filters are sent as the `extra_filters` query parameter. They are applied at the top of the pipeline and do **not** propagate into `join` / `union` subqueries.
* **Root Query** — filters are prepended to the query expression (e.g. `level:="error" | <your-query>`). Use this when you want ad-hoc values to participate in the main pipeline as ordinary filters.
* **Off** — disables automatic ad-hoc filter injection for the panel. Choose this when you intentionally interpolate ad-hoc values yourself (e.g. through template variables) to avoid double filtering, or when filtering on fields produced by pipe transformations would silently drop all rows.


#### Word filter
If you want to filter logs by a specific word or phrase, you can use [word filter](https://docs.victoriametrics.com/victorialogs/logsql/#word-filter).
To do this, you need to create a [text box variable](https://grafana.com/docs/grafana/latest/visualizations/dashboards/variables/add-template-variables/#add-a-text-box-variable).
After it is created, you need to add this to the query: `<q> | $word_filter`. Type any word or phrase you want to look for. Set to "*" to not filter. By default, the given word is searched in the _msg field.
<img alt="Word filter" src="https://github.com/VictoriaMetrics/victorialogs-datasource/blob/main/src/img/word_filter.png?raw=true">

#### Variable interpolation

VictoriaLogs datasource supports automatic variable interpolation with the following rules:

**1. Field Value Variables:**
- `field:$var` → `field:in("v1", ..., "vN")`
- `field:=$var` → `field:in("v1", ..., "vN")`
- Values are quoted and escaped
- Empty values or "All" expand to `in(*)`

**2. Function Contexts:**
- `in($var)` and `contains_any($var)` expand to quoted lists
- Values maintain proper quoting within function calls

**3. Inequality Operators in filters:**
- `field:!$var` → `!field:in("v1", ..., "vN")`
- `field:!=$var` → `!field:in("v1", ..., "vN")`

**4. Stream Filters:**
- `{tag=$var}` → `{tag in(...)}`
- `{field!=$var}` → `{field not_in(...)}`

**5. Restrictions:**
- Avoid using variables as values in regexp filters (e.g., `field:~$var`), in this case, you will get a warning message
- Invalid usage will show an error message

## Incremental hits loading

In Explore, `Raw Logs` queries and their logs volume histogram are loaded incrementally by default:

- the logs are requested first; the histogram request starts once the logs have answered;
- when VictoriaLogs has not started answering the histogram request within 3 seconds, the request is cancelled and
  the histogram is loaded bar by bar from the newest to the oldest, so a specific moment can be zoomed into before
  the whole range is loaded;
- when VictoriaLogs has not started answering the logs request itself within 3 seconds, it is cancelled and the
  histogram is loaded bar by bar right away; the logs are requested again for the whole range once the histogram
  is complete (or stopped).
  With the **Logs volume** panel collapsed the logs request simply keeps running.

While the histogram is loaded bar by bar, a status row in the query editor shows the progress
(`Logs volume is loading incrementally: 12 / 28 bars`) with **Pause**/**Resume** and **Stop** controls. **Pause** lets the bar in flight finish
and holds the next one; **Stop** keeps the bars loaded so far and lets the logs load, unlike the Explore **Cancel**
button which discards everything. Zooming into the histogram while it loads cancels the current requests and starts
over for the new time range.

The behaviour can be switched off per query with the **Incremental hits loading** option in the query editor.
With a **Group hits by** field, the top-20 groups are picked within every bar separately, so the tail
group `other` may differ from the one of a single whole-range request.

## Drilldown

> **Note:** available since [v0.33.0](https://github.com/VictoriaMetrics/victorialogs-datasource/releases/tag/v0.33.0).
> Requires VictoriaLogs v1.21.0 or newer; filtering by a message pattern requires VictoriaLogs v1.33.0 or newer.

The **Drilldown** button in the Explore query editor opens a drawer that explores the logs without writing LogsQL:
it breaks the logs down by field values, by stream fields and by message patterns, narrows the selection with
filter chips and hands the finished selection back to the editor.

<img alt="Drilldown button" width="100%" src="https://github.com/VictoriaMetrics/victorialogs-datasource/blob/main/src/img/drilldown_button.png?raw=true">

> **Note:** for now, the button is available only in the `Code` mode of the Explore query editor. The drawer does not
> take the current query and always starts from the `*` query over the editor's time range.

### Main view

The main view shows one tab per breakdown field:

- the first tab is the service field, when the logs have one (`service.name`, `service_name`, `service` or `app`),
  otherwise the first stream field, or the first non-internal field if there are no stream fields;
- **Stream fields** lists the [stream fields](https://docs.victoriametrics.com/victorialogs/keyconcepts/#stream-fields)
  with their hit counts, and **Break down by** on a stream field shows its values;
- **Patterns** groups the log messages by their shape, with numbers, IPs, UUIDs and the like collapsed into
  placeholders such as `<N>`. The list is estimated from a 1-in-10 sample of the logs, so the counts are marked with `~`;
  filter for a pattern to see its exact count;
- the **+** button adds a tab for any other field.

Every tab shows the volume of the top values over time, and a table with the count, the share and a sparkline
of every value. Expanding a row shows a sample of its logs. The **Filter for** / **Filter out** icons at the end
of a row add a filter chip for the value. The level buttons above the tabs narrow the selection by the log level
derived from the [log level rules](#log-level-rules) without drilling in.

<img alt="Drilldown main view" width="100%" src="https://github.com/VictoriaMetrics/victorialogs-datasource/blob/main/src/img/drilldown_main.png?raw=true">

### Details view

Adding a filter chip or a pattern filter opens the details view of the selection: a volume chart by log level and the tabs

- **Logs** - a sample of the selected logs;
- **Stream fields** and **Fields** - the breakdown of the selection by every stream field and by every other field,
  with a chart of the top values per field. Fields that have the same value in every selected log are listed as
  constant fields, and fields with too many or too long values are listed without a summary;
- **Patterns** - the message patterns of the selection.

<img alt="Drilldown details view" width="100%" src="https://github.com/VictoriaMetrics/victorialogs-datasource/blob/main/src/img/drilldown_details.png?raw=true">

<img alt="Drilldown patterns" width="100%" src="https://github.com/VictoriaMetrics/victorialogs-datasource/blob/main/src/img/drilldown_patterns.png?raw=true">

A chip can be edited in place (field, operator `=`, `!=`, `=~`, `!~`, and value) and removed. Removing the last chip
and pattern filter returns to the main view.

### Time range

The drawer starts with the editor's time range. Selecting a range on any chart zooms the drawer into it, and the
time range picker at the top moves, zooms out or resets the range. The drawer's time range does not change the
Explore time range, and **Go to editor** does not transfer it.

### Go to editor

**Go to editor** replaces the editor query with the drawer's selection: the filter chips become the filter
chips of the query editor, the pattern filters are appended to the query expression, and the query is run.
The original expression and stream filters of the editor are replaced, not merged.

<img alt="Query editor after Go to editor" width="100%" src="https://github.com/VictoriaMetrics/victorialogs-datasource/blob/main/src/img/drilldown_go_to_editor.png?raw=true">

The drawer runs at most 6 queries at a time, so a wide time range with many fields takes a while to load completely.

## Line limits

Every VictoriaLogs query is sent with a limit on the maximum number of log lines returned. Keeping this limit low
reduces the rendering load on the browser and prevents the page from freezing when displaying heavy log results.

There are two levels of control:

- **Maximum lines** — the datasource-wide default, configured in the datasource settings. It applies to every query
  that does not set its own limit. The default value is `1000` and the maximum allowed value is `10000`.
- **Line limit** — a per-query override available in the query editor options. When set, it takes precedence over the
  datasource-wide **Maximum lines** for that query. It is also capped at `10000`.

Increase these limits only if needed; large values may make the browser sluggish.

<img alt="Line limits configuration" src="https://github.com/VictoriaMetrics/victorialogs-datasource/blob/main/src/img/line_limits.png?raw=true">

### Configuring the default limit via file

When [provisioning the datasource via file](#configuration-via-file), set the default limit through the `maxLines`
field in `jsonData`:

```yaml
apiVersion: 1
datasources:
  - name: VictoriaLogs
    type: victoriametrics-logs-datasource
    access: proxy
    url: http://victorialogs:9428
    jsonData:
      # Default maximum number of log lines returned per query.
      # Defaults to 1000, maximum allowed value is 10000.
      maxLines: 1000
```

## Correlations

Signals can be correlated together if they share the same list of attributes, so they can uniquely identify the
same system or event. Grafana provides various interfaces for the [correlations](https://grafana.com/docs/grafana/latest/administration/correlations/)
feature for interactive links between visualizations.

### Trace to logs

Tempo, Jaeger, and Zipkin data sources support [Trace to logs](https://grafana.com/docs/grafana/latest/explore/trace-integration/#trace-to-logs)
feature for navigating from a span in a trace directly to logs relevant for that span. _(Supported since Grafana v12.2.0)_.

<img alt="Derived fields configuration" width="100%" src="https://github.com/VictoriaMetrics/victorialogs-datasource/blob/main/src/img/trace_to_logs.png?raw=true">

An example of the correlation query in traces datasource is the following:
```
trace_id:="${__trace.traceId}" AND span_id:="${__span.spanId}"
```

### Log to metrics

Log to traces correlation is possible via Derived Fields functionality. But for it to work log entries and time series
in metrics datasource should share common labels that could be used as filters. See example of building logs to metrics
correlation in [this issue](https://github.com/VictoriaMetrics/VictoriaMetrics/issues/9592#issuecomment-3202104607).

### Log to traces

Log to traces correlation is possible via Derived Fields functionality. See its description in the sections below.

## OpenTelemetry preset

> **Note:** available since [v0.27.1](https://github.com/VictoriaMetrics/victorialogs-datasource/releases/tag/v0.27.1).

The **OpenTelemetry preset** in the datasource configuration automatically sets up [Derived Fields](#derived-fields) for trace IDs and [Log level rules](#log-level-rules) for OTel severity — without manual configuration.

### How to enable

1. Open the datasource settings.
2. Scroll to the **OpenTelemetry preset** section.
3. Toggle **Enable OpenTelemetry preset**.
4. Click **Save & test** to save your datasource URL first if you have not done so — the preset needs a reachable VictoriaLogs backend to detect the field format.

### What gets auto-detected

After enabling the preset, the datasource queries recent logs to detect:

* **Field format** — `snake_case` (`trace_id`) or `camelCase` (`traceId`).
* **Severity field** — the first matching candidate among `severity_text`, `SeverityText`, `severity`, `Severity`.
* **Severity value case** — whether the severity field contains `UPPERCASE`, `lowercase`, or mixed values, to generate the correct number of log level rules.

The detected configuration is shown in the UI. Your manually created derived fields and log level rules are preserved.

### Traces datasource

Select an optional **Traces datasource** to make the generated `trace_id` derived field link directly to the matching trace. Currently only Jaeger is supported.

### Changing the severity field

If the auto-detected severity field is wrong, click **Change** next to the severity row to pick a different field from the dropdown. The preset will re-run detection using the selected field.

### Provisioning

The preset can be enabled via [Grafana provisioning](https://grafana.com/docs/grafana/latest/administration/provisioning/#data-sources).
Auto-detection runs only in the datasource settings UI, so a provisioned datasource must specify the `detection` block
explicitly — `enabled: true` alone has no effect:

```yaml
apiVersion: 1
datasources:
  - name: VictoriaLogs
    type: victoriametrics-logs-datasource
    access: proxy
    url: http://victorialogs:9428
    jsonData:
      otelPreset:
        enabled: true
        # `uid` of a Jaeger datasource the generated trace ID derived field
        # links to (currently only Jaeger is supported).
        # Omit to skip trace linking.
        tracesDatasourceUid: jaeger
        detection:
          # Log field holding the trace ID, e.g. trace_id, traceId or TraceId.
          traceIdField: trace_id
          # Omit `severity` to skip generating log level rules.
          severity:
            # Log field holding the severity, e.g. severity_text or severity_number.
            field: severity_text
            # `string` for text severities (info, WARN, ...) or
            # `number` for numeric OTel severities (1-24).
            valueCase: string
            # `manual` because the field is set explicitly here
            # (`auto` is reserved for UI auto-detection).
            source: manual
```

Manually configured derived fields and log level rules are merged with the generated ones, user-defined entries win on conflict.

## Derived Fields

In VictoriaLogs datasource settings, you can configure rules of extracting values from a log message to create a link with that value.

For example, if log entries have field `trace_id` then we can configure a Derived Field to make a link to Jaeger datasource
for viewing an associated trace:

<img alt="Derived fields configuration" width="100%" src="https://github.com/VictoriaMetrics/victorialogs-datasource/blob/main/src/img/derived_fields_cfg.png?raw=true">

Once configured, in Explore mode or in Logs panel log entries with field `trace_id` will also get a link that would
open a Jaeger datasource and search for the `trace_id` value:

<img alt="Derived fields explore" src="https://github.com/VictoriaMetrics/victorialogs-datasource/blob/main/src/img/derived_fields_explore.png?raw=true">

If the trace ID is not stored in a separate field but in a log message itself, then use `Regex in log line` option
to specify a regex expression for extracting the trace value from log message.

The target of a Derived Field can be a tracing datasource (Tempo, Jaeger, Zipkin) or another VictoriaLogs datasource
(logs → logs). When the target is VictoriaLogs, set `url` to a LogsQL query such as `trace_id:${__value.raw}`, and
clicking the field opens Explore with that query — letting you pivot from a value like `trace_id` or `request_id`
to all related logs.

Derived Fields can also be configured via provisioning:

```yaml
apiVersion: 1
datasources:
  - name: VictoriaLogs
    type: victoriametrics-logs-datasource
    url: http://victorialogs:9428
    jsonData:
      derivedFields:
        - name: trace_id
          # Take the value from the `trace_id` log field (or use matcherType: regex
          # with a capture group to extract it from the raw log line instead).
          matcherType: label
          matcherRegex: trace_id
          # Value handed to the target datasource: the trace ID for a tracing
          # datasource, or a LogsQL query like 'trace_id:${__value.raw}' for logs → logs.
          url: '${__value.raw}'
          urlDisplayLabel: View trace
          # `uid` of the target datasource (a tracing or VictoriaLogs datasource).
          datasourceUid: tempo
```

Learn more about [Derived Fields in Grafana](https://grafana.com/docs/grafana/next/datasources/loki/configure-loki-data-source/#derived-fields).

## License

This project is licensed under
the [Apache 2.0 license](https://github.com/VictoriaMetrics/victorialogs-datasource/blob/main/LICENSE).
