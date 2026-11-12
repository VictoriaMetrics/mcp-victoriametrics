---
build:
  list: never
  publishResources: false
  render: never
sitemap:
  disable: true
---

The following `tip` changes can be tested by building vmestimator from the latest commits according to the following docs:

* [How to build vmestimator](https://docs.victoriametrics.com/victoriametrics/vmestimator/#how-to-build-from-sources)

Metrics of the latest version of vmestimator cluster are available for viewing at our
[sandbox](https://play-grafana.victoriametrics.com/d/mkv22l4).

## tip

## [v0.1.18](https://github.com/VictoriaMetrics/vmestimator/releases/tag/v0.1.18)

Released at 2026-10-26

* SECURITY: upgrade Go builder from Go1.27.2 to Go1.27.3. See [the list of issues addressed in Go1.27.3](https://github.com/golang/go/issues?q=milestone%3AGo1.27.3%20label%3ACherryPickApproved).
* SECURITY: upgrade base docker image (Alpine) from 3.24.2 to 3.24.3. See [Alpine 3.24.3 release notes](https://alpinelinux.org/posts/Alpine-3.21.9-3.22.7-3.23.7-3.24.3-released.html).

## [v0.1.17](https://github.com/VictoriaMetrics/vmestimator/releases/tag/v0.1.17)

Released at 2026-10-15

**Update Note 1:** [vmestimator](https://docs.victoriametrics.com/victoriametrics/vmestimator/): vmestimator now refuses to start if the `-config` file contains streams with identical `labels`, `group_by`, `interval` and `filter`. Such streams produce identical time series, so differences in other fields (`group_limit`, `buckets`, `hll_precision`, `hll_sparse`, `calculate_churn`, `min_cardinality`) were ambiguous. The `group_by` order is ignored, and an omitted `interval` is treated as the default `5m`. Before upgrading, merge or remove duplicate streams in the config.

**Update Note 2:** [vmestimator](https://docs.victoriametrics.com/victoriametrics/vmestimator/): the `churn_interval` stream config option is replaced with `calculate_churn: true`. Churn is now always calculated over the stream `interval`, and the `churn_interval` label is removed from `cardinality_churn_ratio` metrics. Configs that still contain `churn_interval` fail to load. Before upgrading, replace `churn_interval: <duration>` with `calculate_churn: true` and update queries that rely on the `churn_interval` label. In cluster mode, upgrade all nodes together: snapshots from older storage nodes do not carry the churn flag, so churn is not reported for them until they are upgraded.

* SECURITY: upgrade Go builder from Go1.26.6 to Go1.27.2. See [the list of issues addressed in Go1.27.2](https://github.com/golang/go/issues?q=milestone%3AGo1.27.2%20label%3ACherryPickApproved).

* FEATURE: [dashboards](https://github.com/VictoriaMetrics/vmestimator/tree/main/dashboards): add "Trend" section with a long-term cardinality overview (90 days). Add dedicated "Global Churn %" and "Global Churn" (absolute) panels. Add "Top 10 Churn %" panel per group_by_keys. Remove the `instance` variable.
* FEATURE: [vmestimator](https://docs.victoriametrics.com/victoriametrics/vmestimator/): add per-stream `min_cardinality` config field to override the global `-cardinalityMetrics.minCardinality` flag for individual streams. When unset or `0`, the global flag value is used. A negative value disables the threshold for the stream. See [#46](https://github.com/VictoriaMetrics/vmestimator/issues/46).
* FEATURE: [vmestimator](https://docs.victoriametrics.com/victoriametrics/vmestimator/): support reloading the `-config` file without restart on `SIGHUP` signal, on request to `/-/reload` endpoint (protected with `-reloadAuthKey`) and every `-configCheckInterval`. Streams with unchanged settings keep their accumulated estimates and churn state, while added or changed streams start from scratch. See [#57](https://github.com/VictoriaMetrics/vmestimator/issues/57).

## [v0.1.16](https://github.com/VictoriaMetrics/vmestimator/releases/tag/v0.1.16)

Released at 2026-09-08

* FEATURE: [vmestimator](https://docs.victoriametrics.com/victoriametrics/vmestimator/): add `churn_interval` stream config option. When set, vmestimator emits `cardinality_churn_ratio` metrics measuring how quickly the series set changes within the measurement window. The ratio is computed as `(2*|A∪B| - |A| - |B|) / |A∪B|` where `A` is the retained HLL snapshot and `B` is the current sketch, resulting in a value in `[0, 1]` where `0` means a stable series set and `1` means complete churn. See [churn-calculation](https://docs.victoriametrics.com/victoriametrics/vmestimator/#churn-calculation) and [#24](https://github.com/VictoriaMetrics/vmestimator/issues/24).
* FEATURE: [vmestimator](https://docs.victoriametrics.com/victoriametrics/vmestimator/): simplify overflow handling in the deduplicator. Expose `vmestimator_deduplication_passed_total`, `vmestimator_deduplication_dropped_total`, `vmestimator_deduplication_bloom_filter_size`, and `vmestimator_deduplication_bloom_filter_max_size` metrics for monitoring deduplication effectiveness. See [#48](https://github.com/VictoriaMetrics/vmestimator/pull/48).

## [v0.1.15](https://github.com/VictoriaMetrics/vmestimator/releases/tag/v0.1.15)

Released at 2026-08-17

* FEATURE: [vmestimator](https://docs.victoriametrics.com/victoriametrics/vmestimator/): add time series deduplicator to drop duplicate series within a configurable time window before forwarding to estimators. Use `-deduplication.interval` to enable it, `-deduplication.maxSizeBytes` or `-deduplication.maxSize` to cap memory usage. When cardinality approaches 80% of the limit, the deduplicator gradually falls back to pass-through mode to keep bloom filters below saturation. See [#44](https://github.com/VictoriaMetrics/vmestimator/pull/44).

## [v0.1.14](https://github.com/VictoriaMetrics/vmestimator/releases/tag/v0.1.14)

* FEATURE: [vmestimator](https://docs.victoriametrics.com/victoriametrics/vmestimator/): add `-workers` flag for processing time series insertions concurrently. This improves throughput under high ingestion load. See [#42](https://github.com/VictoriaMetrics/vmestimator/pull/42).

## [v0.1.13](https://github.com/VictoriaMetrics/vmestimator/releases/tag/v0.1.13)

* FEATURE: [vmestimator](https://docs.victoriametrics.com/victoriametrics/vmestimator/): add `-cardinalityMetrics.minCardinality` flag to suppress group-by cardinality estimates below the given threshold. When estimates are filtered, the number of suppressed groups is exposed via the `vmestimator_cardinality_estimates_dropped` metric. See [#41](https://github.com/VictoriaMetrics/vmestimator/pull/41).

## [v0.1.12](https://github.com/VictoriaMetrics/vmestimator/releases/tag/v0.1.12)

* BUGFIX: [vmestimator](https://docs.victoriametrics.com/victoriametrics/vmestimator/): fix panic when stream config contains __label__ pseudo-label.

## [v0.1.11](https://github.com/VictoriaMetrics/vmestimator/releases/tag/v0.1.11)

**Update Note 1:** This version introduces a breaking change in the cluster protocol. Please update everything to the latest version in one go and ignore errors during deployment.

* FEATURE: [vmestimator](https://docs.victoriametrics.com/victoriametrics/vmestimator/): add `filter` field to stream configuration for filtering time series by label matchers before counting. Supports equality (`=`), negative equality (`!=`), regexp (`=~`), and negative regexp (`!~`) matchers in MetricsQL selector syntax, e.g. `{job="api",env!~"dev|staging"}`. See [#29](https://github.com/VictoriaMetrics/vmestimator/pull/29).
* FEATURE: [vmestimator](https://docs.victoriametrics.com/victoriametrics/vmestimator/): decouple metrics generation logic from estimator. Do it via snapshot. This refactoring opens ways to implement upcoming profiler related features.

* BUGFIX: [vmestimator](https://docs.victoriametrics.com/victoriametrics/vmestimator/): revert prev sketches reuse of HLL sketches as it hurts estimate precision.

## [v0.1.10](https://github.com/VictoriaMetrics/vmestimator/releases/tag/v0.1.10)

* FEATURE: [vmestimator](https://docs.victoriametrics.com/victoriametrics/vmestimator/): Reduce memory pressure on rotation by reusing HLL sketches from previous group.

## [v0.1.9](https://github.com/VictoriaMetrics/vmestimator/releases/tag/v0.1.9)

Released at 2026-08-03

* FEATURE: [vmestimator](https://docs.victoriametrics.com/victoriametrics/vmestimator/): add `__label__` pseudo-label for `group_by`. It estimates unique values per label name, making it easy to spot high-cardinality labels like `trace_id` or `user_id`. Can be combined with explicit keys, e.g. `["job", "__label__"]`. See [#26](https://github.com/VictoriaMetrics/vmestimator/issues/26).

* BUGFIX: [vmestimator](https://docs.victoriametrics.com/victoriametrics/vmestimator/): fix snapshot merging collision when multiple estimators share the same `group_by` keys but use different `interval` values. Previously, snapshots from such estimators were incorrectly merged into a single entry, causing the interval of one to overwrite the other.

## [v0.1.8](https://github.com/VictoriaMetrics/vmestimator/releases/tag/v0.1.8)

Released at 2026-07-17

* FEATURE: [alerts](https://github.com/VictoriaMetrics/vmestimator/blob/main/deployment/docker/rules/alerts-cardinality.yml): add `GlobalChurnTooHigh` alert to detect when global series churn or cardinality growth exceeds 10%. See [#13](https://github.com/VictoriaMetrics/vmestimator/pull/13).
* FEATURE: [alerts](https://github.com/VictoriaMetrics/vmestimator/blob/main/deployment/docker/rules/alerts-cardinality.yml): update `JobChurnTooHigh` alert (previously `JobTooHighChurnRate`) to also detect rapid cardinality growth. The threshold was raised from 10% to 20% and the `for` window extended from 15m to 30m to reduce false positives. See [#13](https://github.com/VictoriaMetrics/vmestimator/pull/13).
* FEATURE: [alerts](https://github.com/VictoriaMetrics/vmestimator/blob/main/deployment/docker/rules/alerts-cardinality.yml): replace the static-threshold `JobTooHighCardinality` alert with adaptive 3-sigma anomaly detection alerts `GlobalCardinalityTooHigh` and `JobCardinalityTooHigh`. See [#17](https://github.com/VictoriaMetrics/vmestimator/pull/17).

* BUGFIX: [vmestimator](https://docs.victoriametrics.com/victoriametrics/vmestimator/): add `interval` label to `vmestimator_estimator_insert_total` metric in order to avoid exposing duplicate series when multiple streams share the same `group_by` but use different intervals. The insert-rate panel in the bundled Grafana dashboard now deduplicates per-interval series with `max without(interval)` before summing the remaining series, so the rate is not multiplied by the number of per-interval streams. See [#20](https://github.com/VictoriaMetrics/vmestimator/issues/20).
* BUGFIX: [vmestimator](https://docs.victoriametrics.com/victoriametrics/vmestimator/): fix inaccurate cardinality estimates when the number of unique series was in the upper ~10% of the sparse mode. Estimates in this region could be off due to premature transition to a less precise (dense) counting mode. See [hyperloglog#d44d606f](https://github.com/makasim/hyperloglog/commit/d44d606f7e8bdd78d2b56c27fc6fe82f3981d4a6).

## [v0.1.7](https://github.com/VictoriaMetrics/vmestimator/releases/tag/v0.1.7)

Released at 2026-07-09

* FEATURE: [vmestimator](https://docs.victoriametrics.com/victoriametrics/vmestimator/): randomize estimator iteration order to reduce lock contention. See [#15](https://github.com/VictoriaMetrics/vmestimator/pull/15).

## [v0.1.6](https://github.com/VictoriaMetrics/vmestimator/releases/tag/v0.1.6)

Released at 2026-07-07

* FEATURE: [vmestimator](https://docs.victoriametrics.com/victoriametrics/vmestimator/): initial release
