---
draft: false
page: blog blog_post
authors:
  - Pablo Fernandez
date: 2026-11-10
enableComments: true
title: "VictoriaMetrics October 2026 Ecosystem Updates"
summary: "October's roundup brings alerting for logs and traces in VictoriaMetrics Cloud, production-ready SSO in vmauth, peer-group outlier detection in vmanomaly, and security fixes across the ecosystem."
toc: true
categories:
  - Product News
  - Company News
tags:
  - cloud
  - victoriametrics
  - vmestimator
  - vmanomaly
  - new features
images:
  - /blog/victoriametrics-october-2026-ecosystem-updates/preview.webp
---

October is the month for things that go bump in the night, so we spent it making sure you hear about them. VictoriaMetrics Cloud now supports alerting rules for VictoriaLogs and VictoriaTraces deployments, not just metrics.

Elsewhere, vmauth SSO is out of beta and ready for production. vmanomaly learned to spot the one replica acting differently from the rest. We also did some security housekeeping across VictoriaMetrics and vmestimator. No tricks, mostly treats.

This release roundup includes updates for:

- [VictoriaMetrics Cloud](#cloud)
- VictoriaMetrics [v1.155.0](#vm11550) and [v1.154.0](#vm11540)
- Anomaly Detection [v1.31.1](#vmanomaly1311) and [v1.31.0](#vmanomaly1310)
- vmestimator [v0.1.18](#vmestimator0118) and [v0.1.17](#vmestimator0117)

## VictoriaMetrics Cloud Update {#cloud}

**Highlights**

- Alerting for logs and traces: you can now set up alerting rules on VictoriaLogs and VictoriaTraces deployments, not only on metrics. Linking a metrics deployment is optional. If you link one, it stores your recording-rule results. If you don't, alerting rules still work on their own.
- [Stream aggregation on single-node deployments](https://docs.victoriametrics.com/victoriametrics-cloud/deployments/stream-aggregation/): single-node deployments now have stream aggregation. You can find it in a new tab under Advanced settings. You can aggregate data as it is ingested, which cuts storage and query costs.
- Alerting status at a glance: Alertmanager and Alerting rules pages now show one combined health status for alerting. It uses the same badges as deployment status, so you can see quickly whether alerting is healthy.
- Refreshed rule presets: ready-made alerting rule sets have a new integrations page. You can also browse them directly on the Alerting rules settings page, so adding a preset takes fewer clicks.
- VictoriaTraces agent integration (vtagent): a new [Helm chart integration for vtagent](https://docs.victoriametrics.com/victoriametrics-cloud/integrations/vtagent/) joins the existing vmagent and vlagent guides. Sending traces to Cloud now works the same way as sending metrics and logs.

**Improvements**

- Faster onboarding in integrations: if you have no deployment, API key, or access token yet, you can create one right from the integration dropdown.
- Faster cluster updates: changes that don't touch storage nodes now go through a quicker update path.
- More tuning for single-node: [the `maxConcurrentInserts` flag](https://docs.victoriametrics.com/victoriametrics/single-server-victoriametrics/#list-of-command-line-flags) is now configurable.
- More reliable live tailing: live log tailing keeps connections alive during long sessions.

**Fixes**

- Storage recommendations are now correct when you switch between GB and TB.
- Adding rule sets from integrations to single-node deployments no longer fails.
- Public integration pages are fixed.
- Recording rule creation for logs and traces deployments is fixed.

## VictoriaMetrics 1.155.0 is out: Security fixes and safer defaults {#vm11550}

> This release includes security fixes for recently discovered vulnerabilities, including the [eMudhra private key compromise](https://bugzilla.mozilla.org/show_bug.cgi?id=2083449). It also fixes three of the four serious bugs introduced in v1.154.0 (see the [v1.154.0 update notes](https://docs.victoriametrics.com/victoriametrics/changelog/#v11540)). The remaining bug only affects users who set [the `-search.skipSlowReplicas` flag](https://docs.victoriametrics.com/victoriametrics/cluster-victoriametrics/#list-of-command-line-flags-for-vmselect) (see [#11771](https://github.com/VictoriaMetrics/VictoriaMetrics/issues/11771) for details).
>

- Potential breaking change: in [VictoriaMetrics cluster](https://docs.victoriametrics.com/victoriametrics/cluster-victoriametrics/), `/admin/tenants` requests without `start` and `end` query args now return only tenants for the current day, instead of searching the entire retention period. To retrieve all tenants, set `start` and `end` so the time range fully covers `-retentionPeriod`. [vmctl](https://docs.victoriametrics.com/victoriametrics/vmctl/victoriametrics/) calls `/admin/tenants` under the hood when migrating data, so make sure you run it with both `*-filter-time-start` and `*-filter-time-end` params. See the updated [URL format](https://docs.victoriametrics.com/victoriametrics/cluster-victoriametrics/#url-format) for details.
- v1.152.0 remains the most stable version.

[Read the changelog](https://docs.victoriametrics.com/victoriametrics/changelog/#v11550)

## VictoriaMetrics 1.154.0 is out: Production-ready SSO in vmauth {#vm11540}

- vmauth Single Sign-On (SSO) ready for production: [Single Sign-On](https://docs.victoriametrics.com/victoriametrics/vmauth/#single-sign-on-sso) is now fully functional and can be recommended for production use.
- Cancels queries when clients disconnect: vmselect and single-node VictoriaMetrics now cancel query requests when the client disconnects and during graceful shutdown. CPU and memory are no longer spent computing results nobody will read.
- Fixes the pure-Go zstd panic introduced in v1.153.0: vmsingle and vmstorage no longer crash with a `cannot decompress zstd data...` panic on `darwin`, `windows`, `freebsd`, `openbsd`, `linux/arm`, `linux/ppc64le`, `linux/386`, `linux/s390x` and `-pure` builds. The `linux/amd64` and `linux/arm64` builds weren't affected. See the [v1.153.0 update note](https://docs.victoriametrics.com/victoriametrics/changelog/#v11530) for details.
- This release also ships 18 bug and 3 security fixes

## vmanomaly 1.31.1 is out: Peer-group models no longer stall on first fit {#vmanomaly1311}

Use v1.31.1 or later for peer-group detection with a [periodic scheduler](https://docs.victoriametrics.com/anomaly-detection/components/scheduler/#periodic-scheduler). It fixes a v1.31.0 issue that could block initial training and prevent anomaly scores from being produced.

[Read the changelog](https://docs.victoriametrics.com/anomaly-detection/changelog/#v1311)

## vmanomaly 1.31.0 is out: Find unusual replicas and tune detection with your own examples {#vmanomaly1310}

October added peer-group detection and a dedicated investigation UI. The new [model-selection guide](https://docs.victoriametrics.com/anomaly-detection/guides/model-selection/) helps you choose between this approach and existing univariate and multivariate models, with six ready-to-run recipes.

<section
  class="model-comparison"
  aria-label="Three ways to detect unusual behavior"
  aria-roledescription="carousel"
>
  <p><strong>Three ways to detect unusual behavior in your metrics</strong></p>
  <div class="model-comparison__controls" hidden>
    <button type="button" aria-pressed="true">
      Peer-group — new in October
    </button>
    <button type="button" aria-pressed="false">Univariate</button>
    <button type="button" aria-pressed="false">Multivariate</button>
  </div>
  <figure
    class="model-comparison__slide"
    role="group"
    aria-roledescription="slide"
    aria-label="1 of 3: Peer-group — new in October"
  >
    <a
      href="https://docs.victoriametrics.com/anomaly-detection/guides/model-selection/peer-group-comparison.webp"
      ><img
        src="https://docs.victoriametrics.com/anomaly-detection/guides/model-selection/peer-group-comparison.webp"
        alt="One replica stays near 84% CPU while seven comparable peers stay near 50%. Its own signal is steady, but it stands out from the group."
        width="2024"
        height="773"
    /></a>
    <figcaption>
      <strong>Is one replica behaving differently from comparable peers?</strong
      ><br />One replica stays near 84% CPU while seven comparable peers stay
      near 50%. Its own signal is steady, but it stands out from the group.
    </figcaption>
  </figure>
  <figure
    class="model-comparison__slide"
    role="group"
    aria-roledescription="slide"
    aria-label="2 of 3: Univariate"
  >
    <a
      href="https://docs.victoriametrics.com/anomaly-detection/guides/model-selection/univariate-comparison.webp"
      ><img
        src="https://docs.victoriametrics.com/anomaly-detection/guides/model-selection/univariate-comparison.webp"
        alt="CPU utilization rises above the range learned from this instance&#x27;s earlier observations."
        width="2024"
        height="773"
    /></a>
    <figcaption>
      <strong
        >Has this metric changed unexpectedly compared with its own
        history?</strong
      ><br />CPU utilization rises above the range learned from this instance's
      earlier observations.
    </figcaption>
  </figure>
  <figure
    class="model-comparison__slide"
    role="group"
    aria-roledescription="slide"
    aria-label="3 of 3: Multivariate"
  >
    <a
      href="https://docs.victoriametrics.com/anomaly-detection/guides/model-selection/multivariate-comparison.webp"
      ><img
        src="https://docs.victoriametrics.com/anomaly-detection/guides/model-selection/multivariate-comparison.webp"
        alt="CPU stays low while RAM stays high. Each value is familiar on its own, but their combination is unusual."
        width="2024"
        height="773"
    /></a>
    <figcaption>
      <strong>Are related signals behaving unusually together?</strong><br />CPU
      stays low while RAM stays high. Each value is familiar on its own, but
      their combination is unusual.
    </figcaption>
  </figure>
  <p class="model-comparison__status" aria-live="polite" hidden></p>
  <p><small>Select an image to view it full size.</small></p>
</section>
<style>
  .model-comparison {
    margin: 1.5rem 0;
  }
  .model-comparison__controls {
    display: flex;
    flex-wrap: wrap;
    gap: 0.5rem;
    margin-bottom: 1rem;
  }
  .model-comparison [hidden] {
    display: none !important;
  }
  .model-comparison__controls button {
    font: inherit;
    color: inherit;
    background: transparent;
    border: 1px solid currentColor;
    border-radius: 0.5rem;
    padding: 0.6rem 0.85rem;
    cursor: pointer;
  }
  .model-comparison__controls button[aria-pressed="true"] {
    color: #fff;
    background: #174fc4;
    border-color: #174fc4;
  }
  .model-comparison__controls button:focus-visible {
    outline: 3px solid #174fc4;
    outline-offset: 3px;
  }
  .model-comparison__slide {
    margin: 0 0 1rem;
  }
  .model-comparison__slide img {
    display: block;
    width: 100%;
    height: auto;
  }
  .model-comparison__slide figcaption {
    margin-top: 0.75rem;
    color: inherit;
  }
  .model-comparison__status {
    font-size: 0.875rem;
  }
</style>
<script>
  document.querySelectorAll(".model-comparison").forEach((carousel) => {
    if (carousel.dataset.ready) return;
    carousel.dataset.ready = "true";
    const controls = carousel.querySelector(".model-comparison__controls");
    const buttons = [...controls.querySelectorAll("button")];
    const slides = [...carousel.querySelectorAll(".model-comparison__slide")];
    const status = carousel.querySelector(".model-comparison__status");
    const select = (index) => {
      slides.forEach((slide, i) => {
        slide.hidden = i !== index;
      });
      buttons.forEach((button, i) =>
        button.setAttribute("aria-pressed", String(i === index)),
      );
      status.textContent = `${index + 1} of ${slides.length}: ${buttons[index].textContent}`;
    };
    controls.hidden = false;
    status.hidden = false;
    buttons.forEach((button, index) => {
      button.addEventListener("click", () => select(index));
      button.addEventListener("keydown", (event) => {
        if (!["ArrowLeft", "ArrowRight", "Home", "End"].includes(event.key))
          return;
        event.preventDefault();
        let next =
          event.key === "ArrowLeft"
            ? (index + buttons.length - 1) % buttons.length
            : (index + 1) % buttons.length;
        if (event.key === "Home") next = 0;
        if (event.key === "End") next = buttons.length - 1;
        select(next);
        buttons[next].focus();
      });
    });
    select(0);
  });
</script>

- **Find the unusual replica.** The new [peer-group model](https://docs.victoriametrics.com/anomaly-detection/components/models/#peer-outlier) scores each member separately as peers join or leave. The [investigation view](https://docs.victoriametrics.com/anomaly-detection/ui/#peer-group-investigation) highlights affected peers and their behavior.
- **Tune detection with your own examples.** The experimental [Tune workflow](https://docs.victoriametrics.com/anomaly-detection/ui/#tune-a-model) lets you mark known anomalies or specify an expected anomaly rate, then review and apply suggested settings across univariate, multivariate, and peer-group models.
- **Estimate resources before deployment.** Experimental [sizing tools](https://docs.victoriametrics.com/anomaly-detection/components/server/#deployment-sizing) estimate CPU, memory, disk, and output requirements—or how much inference work fits within a resource budget. Also available through the [vmanomaly MCP server v0.5.0](https://github.com/VictoriaMetrics/mcp-vmanomaly/releases/tag/v0.5.0). Validate estimates against your workload and allow extra headroom, especially for peer-group models.
- **Process more data and investigate results.** Inference throughput increased by up to 25% in tested workloads, depending on the model and configuration. The [anomaly-score dashboard](https://docs.victoriametrics.com/anomaly-detection/presets/#grafana-dashboard) adds improved statistics, consistent filters, and forecast-disagreement trends.

> Before a stateful upgrade or rollback, check the [compatibility matrix](https://docs.victoriametrics.com/anomaly-detection/migration/#compatibility-matrix) and run [`--dryRun`](https://docs.victoriametrics.com/anomaly-detection/migration/#dry-run) to preview persisted-state migration.

[Read the vmanomaly changelog](https://docs.victoriametrics.com/anomaly-detection/changelog/#v1310) and the [UI v1.10.0 changelog](https://docs.victoriametrics.com/anomaly-detection/ui/#v1100)

## vmestimator 0.1.18 is out: Go and Alpine security updates {#vmestimator0118}

> If you're upgrading from v0.1.16 or earlier, review the [v0.1.17 update notes](https://github.com/VictoriaMetrics/vmestimator/releases/tag/v0.1.17) first. vmestimator now refuses to start if the `-config` file contains duplicate streams. The `churn_interval` stream option is replaced with `calculate_churn: true`, and configs that still use `churn_interval` fail to load.

- Security fixes: this release contains only security updates, with no new features or bug fixes. We recommend upgrading to pick up the latest upstream fixes.
- Upgrades the Go builder to Go 1.27.3: [vmestimator](https://docs.victoriametrics.com/victoriametrics/vmestimator/) is now built with Go 1.27.3 instead of Go 1.27.2, picking up the [issues fixed in Go 1.27.3](https://github.com/golang/go/issues?q=milestone%3AGo1.27.3%20label%3ACherryPickApproved).
- Upgrades the base Docker image to Alpine 3.24.3: The base Docker image moves from Alpine 3.24.2 to 3.24.3 and includes its latest security patches. See the [Alpine 3.24.3 release notes](https://alpinelinux.org/posts/Alpine-3.21.9-3.22.7-3.23.7-3.24.3-released.html) for details.

[Read the changelog](https://github.com/VictoriaMetrics/vmestimator/blob/main/docs/vmestimator/CHANGELOG.md#v0118)

## vmestimator 0.1.17 is out: Hot config reload without restarts {#vmestimator0117}


> **Update note:** The `churn_interval` stream config option is replaced by `calculate_churn: true`. Churn is now always calculated over the stream `interval`, and the `churn_interval` label is removed from `cardinality_churn_ratio` metrics. Update your stream configs, dashboards, and alerts accordingly.

- Supports hot config reload: vmestimator can now reload its `-config` file without restarting, so you can change streams without interrupting cardinality estimation. Reloads are triggered by a `SIGHUP` signal, by a request to the `/-/reload` endpoint (protected with `-reloadAuthKey`), or automatically every `-configCheckInterval`. See the [config reload docs](https://docs.victoriametrics.com/victoriametrics/vmestimator/#config-reload) for details.
- Adds per-stream minimum cardinality: The new `min_cardinality` stream config field overrides the global `-cardinalityMetrics.minCardinality` flag for individual streams. This gives you finer control over which series each stream reports. When the field is unset or `0`, the global flag value is used. See the [vmestimator docs](https://docs.victoriametrics.com/victoriametrics/vmestimator/) for configuration details.
- Rejects duplicate stream definitions: vmestimator now validates the config at startup and refuses to start when two streams are effectively identical. This prevents duplicate work and conflicting metrics.
- Simplifies churn calculation: The `calculate_churn: true` option replaces `churn_interval`. Churn is always measured over the stream's own `interval`, which makes the configuration simpler and the results more consistent.

[Read the changelog](https://github.com/VictoriaMetrics/vmestimator/blob/main/docs/vmestimator/CHANGELOG.md#v0117)

## Until next month

We can't promise a quiet on-call shift. We can promise that if something goes wrong, you'll find out sooner and figure out why faster. Grab some candy and enjoy the spooky season. We'll be back next month with more.

