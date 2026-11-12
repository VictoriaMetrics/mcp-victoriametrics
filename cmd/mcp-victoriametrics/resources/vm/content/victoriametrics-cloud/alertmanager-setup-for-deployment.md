---
weight: 15
title: Setup Alertmanager & VMAlert for VictoriaMetrics Cloud
menu:
  docs:
    parent: "cloud"
    weight: 15
tags:
   - metrics
   - logs
   - traces
   - cloud
   - enterprise
   - guide
aliases:
  - /victoriametrics-cloud/alertmanager-setup-for-deployment/index.html
  - /managed-victoriametrics/alertmanager-setup-for-deployment/index.html
---
VictoriaMetrics Cloud supports configuring alerting rules, powered by vmalert, and sending notifications with hosted Alertmanager.
This is supported for `VictoriaMetrics`, `VictoriaLogs` and `VictoriaTraces` deployments.

## Configure Alertmanager

You have two options to configure Cloud Alertmanager:

1. From integrations section: Menu **"Integrations" `->` "Cloud Alertmanager" `->` "New configuration"**:
   ![Setup for deployment integrations](alertmanager-setup-for-deployment_integrations.webp?card)
2. From deployment page: **"Deployment page" `->` "Rules" tab `->` "Settings" `->` "Connect notifier" `/` "New notifier"**:
   ![Setup for deployment connect notifier](alertmanager-setup-for-deployment_connect_notifier.webp)

For creating a new configuration, you need to provide the following parameters:

- **Name of the configuration** (it only affects the display in the user interface)
- **Configuration file** in [specified format](https://docs.victoriametrics.com/victoriametrics-cloud/alertmanager-setup-for-deployment/#alertmanager-config-specification)

Before saving the configuration, you can validate it by clicking the "Test configuration" button.

After creating the configuration, you can connect it to one or multiple deployments.
In order to do this you need to go to the "Deployment page" `->` "Rules" tab `->` "Settings" `,
select the created notifier and confirm the action:

![Select notifier](alertmanager-setup-for-deployment_select_notifier.webp)

Alertmanager is now set up for your deployment, and you be able to get notifications from it.

### Alertmanager config specification

VictoriaMetrics Cloud supports Alertmanager with standard [configuration specification](https://prometheus.io/docs/alerting/latest/configuration/).

But with respect to the specification, there are the following limitations:

1. Only allowed receivers:
   * `discord_configs`
   * `pagerduty_configs`
   * `slack_configs`
   * `webhook_configs`
   * `opsgenie_configs`
   * `wechat_configs`
   * `pushover_configs`
   * `victorops_configs`
   * `telegram_configs`
   * `webex_configs`
   * `msteams_configs`
   * `msteamsv2_configs`
   * `jira_configs`
   * `rocketchat_configs`
   * `mattermost_configs`
   * `incidentio_configs`
2. All configuration params with `_file` suffix are not allowed for security reasons.
3. The maximum file size is 20mb.

### Configuration example

Here is an example of Alertmanager configuration:

```yaml
route:
 receiver: slack-infra
 repeat_interval: 1m
 group_interval: 30s
 routes:
 - matchers:
   - team = team-1
   receiver: dev-team-1
   continue: true
 - matchers:
   - team = team-2
   receiver: dev-team-2
   continue: true
receivers:
- name: slack-infra
  slack_configs:
  - api_url: https://hooks.slack.com/services/valid-url
    channel: infra
    title: |-
        [{{ .Status | toUpper -}}
        {{ if eq .Status "firing" }}:{{ .Alerts.Firing | len }}{{- end -}}
        ]
        {{ if ne .Status "firing" -}}
          :lgtm:
          {{- else if eq .CommonLabels.severity "critical" -}}
          :fire:
          {{- else if eq .CommonLabels.severity "warning" -}}
          :warning:
          {{- else if eq .CommonLabels.severity "info" -}}
          :information_source:
          {{- else -}}
          :question:
        {{- end }}
    text: |
        {{ range .Alerts }}
        {{- if .Annotations.summary }}
           Summary:  {{ .Annotations.summary }}
        {{- end }}
        {{- if .Annotations.description }}
            Description: {{ .Annotations.description }}
        {{- end }}
        {{- end }}
    actions:
    - type: button
      text: 'Query :mag:'
      url: '{{ (index .Alerts 0).GeneratorURL }}'
    - type: button
      text: 'Silence :no_bell:'
      url: '{{ template "__silenceURL" . }}'
- name: dev-team-1
  slack_configs:
  - api_url: https://hooks.slack.com/services/valid-url
    channel: dev-alerts
- name: dev-team-2
  slack_configs:
  - api_url: https://hooks.slack.com/services/valid-url
    channel: dev-alerts
```

### Custom Alertmanager

If for some reason Cloud Alertmanager is not suitable for you, you can use VictoriaMetrics Cloud with any external Alertmanager hosted in your infrastructure.

For that select Custom Alertmanager instead of Cloud Alertmanager when [creating the Alertmanager](https://docs.victoriametrics.com/victoriametrics-cloud/alertmanager-setup-for-deployment/#configure-alertmanager):

![Custom AlertManager](alertmanager-setup-for-deployment_custom_am.webp)

Limitations for the Custom Alertmanager:

- Your custom Alertmanager should be available from the Internet via **HTTPS** with **Basic authentication** or **Bearer token authentication**.
- You will not be able to use "Alerts" tab on the deployment page.

You can test the connection to your custom Alertmanager by clicking the "Test connection" button.
`/api/v2/status` endpoint is used to verify that connection to Alertmanager is working.

## Configure alerting and recording rules

Alerting and recording rules could be uploaded on **"Deployment page" `->` "Rules" tab `->` "Settings"**:

![Upload rules](alertmanager-setup-for-deployment_upload_rules.webp)

You can click on the upload area or drag and drop the files with rules there.

Files should be in the [Prometheus alerting rules definition format](https://prometheus.io/docs/prometheus/latest/configuration/alerting_rules/)
or [Prometheus recording rules definition format](https://prometheus.io/docs/prometheus/latest/configuration/recording_rules/).

For `VictoriaLogs` and `VictoriaTraces` deployments, each rule group must set `type: vlogs` and use
[LogsQL](https://docs.victoriametrics.com/victorialogs/logsql/) expressions instead of MetricsQL/PromQL
ones (`VictoriaTraces` rule groups also use `type: vlogs`, since both query languages share the same
underlying evaluation engine). See the [example below](#example-of-alerting-rules-for-victorialogs-and-victoriatraces).

There are limitations for the rules files:

1. All files may contain no more than 100 rules in total. If you need to upload more rules contact us via [support-cloud@victoriametrics.com](mailto:support-cloud@victoriametrics.com).
2. The maximum file size is 20mb.
3. The names of the groups in the files should be unique.

You can also use API for uploading rules. Switch to **"Upload with API"** on the page and follow the instructions:

- Choose the API key for uploading rules
- After that you can copy the curl command for uploading rules and execute it in your terminal

![Upload with API](alertmanager-setup-for-deployment_upload_with_api.webp)

You can use the following API endpoints for the automation with rules:

* POST: `/api/v1/deployments/{deploymentId}/rule-sets/files/{fileName}` - create/update rules file
* DELETE `/api/v1/deployments/{deploymentId}/rule-sets/files/{fileName}` - delete rules file

For more details, please check [OpenAPI Reference](https://console.victoriametrics.cloud/api-docs).

### Example of alerting rules

#### VictoriaMetrics

Here is an example of alerting rules in the Prometheus alerting rules format:

```yaml
groups:
  - name: examples
    concurrency: 2
    interval: 10s
    rules:
      - alert: never-firing
        expr: foobar > 0
        for: 30s
        labels:
          severity: warning
        annotations:
          summary: empty result rule
      - alert: always-firing
        expr: vector(1) > 0
        for: 30s
        labels:
          severity: critical
        annotations:
          summary: "rule must be always at firing state"
```

#### VictoriaLogs and VictoriaTraces {#example-of-alerting-rules-for-victorialogs-and-victoriatraces}

Here is an example of an alerting rule for a `VictoriaLogs` or `VictoriaTraces` deployment, using
`type: vlogs` and a LogsQL expression instead of a PromQL one:

```yaml
groups:
  - name: vlogs-group
    type: vlogs
    rules:
      - alert: too-many-errors
        expr: |
          _time:5m error
          | stats count() as errors
          | filter errors > 100
        for: 1m
        labels:
          severity: warning
        annotations:
          summary: "more than 100 errors logged in the last 5 minutes"
```

## Linking a metrics deployment for VictoriaLogs and VictoriaTraces alerting

`VictoriaLogs` and `VictoriaTraces` deployments don't store Prometheus-style metrics themselves, so
vmalert needs a `VictoriaMetrics` deployment to write recording rule output and alert state to. You can
link one on **"Deployment page" `->` "Rules" tab `->` "Settings"**, under "Configure remote-write
target". Only `VictoriaMetrics` deployments in the same region can be linked.

Without a linked metrics deployment:

- Recording rules are not supported: files containing recording rules are rejected on upload.
- Alerting rules still run normally and fire alerts to Alertmanager, but alert state won't survive
  vmalert restarts, since there's nowhere to persist it.

Unlinking an already-linked deployment deletes any recording rule files, since their output would
otherwise have nowhere to be written.

## Troubleshooting

### Rules execution state

The state of created rules is located in the `Rules` section of your deployment:

![Rules state](alertmanager-setup-for-deployment_rules_state.webp)

### Debug

It's possible to debug the alerting stack with logs for vmalert and Alertmanager, which are accessible in the `Logs` section of the deployment.

![Troubleshoot logs](alertmanager-setup-for-deployment_troubleshoot_logs.webp)

### Monitoring

Alertmanager and vmalert errors are tracked by a built-in monitoring system.
Deployment's `Alerts` section has information about active incidents and incident history log.
