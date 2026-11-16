

---
weight: 2
title: VictoriaLogs Collector
description: "Collect logs from Kubernetes containers to VictoriaLogs."
menu:
  docs:
    parent: helm
    weight: 2
    identifier: helm-victorialogs-collector
url: /helm/victoria-logs-collector
aliases:
  - /helm/victorialogs-collector/index.html
  - /helm/victorialogs-collector/
tags:
  - logs
  - kubernetes
---

![Version](https://img.shields.io/badge/0.3.9-gray?logo=Helm&labelColor=gray&link=https%3A%2F%2Fdocs.victoriametrics.com%2Fhelm%2Fvictoria-logs-collector%2Fchangelog%2F%23v039)
![ArtifactHub](https://img.shields.io/badge/ArtifactHub-informational?logoColor=white&color=417598&logo=artifacthub&link=https%3A%2F%2Fartifacthub.io%2Fpackages%2Fhelm%2Fvictoriametrics%2Fvictoria-logs-collector)
![License](https://img.shields.io/github/license/VictoriaMetrics/helm-charts?labelColor=green&label=&link=https%3A%2F%2Fgithub.com%2FVictoriaMetrics%2Fhelm-charts%2Fblob%2Fmaster%2FLICENSE)
![Slack](https://img.shields.io/badge/Join-4A154B?logo=slack&link=https%3A%2F%2Fslack.victoriametrics.com)
![X](https://img.shields.io/twitter/follow/VictoriaMetrics?style=flat&label=Follow&color=black&logo=x&labelColor=black&link=https%3A%2F%2Fx.com%2FVictoriaMetrics)
![Reddit](https://img.shields.io/reddit/subreddit-subscribers/VictoriaMetrics?style=flat&label=Join&labelColor=red&logoColor=white&logo=reddit&link=https%3A%2F%2Fwww.reddit.com%2Fr%2FVictoriaMetrics)

VictoriaLogs Collector - collects logs from Kubernetes containers and stores them to VictoriaLogs

## Prerequisites

Before installing this chart, ensure your environment meets the following requirements:

* **Kubernetes cluster** - A running Kubernetes cluster with sufficient resources
* **Helm** - Helm package manager installed and configured

Additional requirements depend on your configuration:

* **Persistent storage** - Required if you enable persistent volumes for data retention (enabled by default)
* **kubectl** - Needed for cluster management and troubleshooting

For installation instructions, refer to the official documentation:
* [Installing Helm](https://helm.sh/docs/intro/install/)
* [Installing kubectl](https://kubernetes.io/docs/tasks/tools/install-kubectl/)

## Quick start

This Helm chart deploys a log collection agent as a DaemonSet.
It collects logs from all containers in a Kubernetes cluster and forwards them to the configured VictoriaLogs destinations.
If more than one destination is specified, then the collected logs are replicated among the configured destinations.

This chart will expand its functionality as the corresponding features are added to [vlagent](https://docs.victoriametrics.com/victorialogs/vlagent/).

- To quickly install single-node version of VictoriaLogs and `victoria-logs-collector`, see [these docs](https://docs.victoriametrics.com/helm/victoria-logs-single/#quick-start).
- To start with a VictoriaLogs cluster and `victoria-logs-collector`, see [these docs](https://docs.victoriametrics.com/helm/victoria-logs-cluster/#quick-start).

## Chart configuration

The simplest working configuration includes specifying the `remoteWrite` array and setting CPU and memory resources for the chart.

Example of a minimal working configuration:

```yaml
remoteWrite:
  - url: http://victoria-logs:9428

resources:
  limits:
    cpu: 100m
    memory: 128Mi
  requests:
    cpu: 100m
    memory: 128Mi
```

If multiple `remoteWrite` entries are defined, logs are replicated to all the specified destinations.

See [victoria-logs-single](https://docs.victoriametrics.com/helm/victoria-logs-single/#how-to-install) and [victoria-logs-cluster](https://docs.victoriametrics.com/helm/victoria-logs-cluster/#how-to-install) documentation for details how to determine remote write URL for each installation type.

### On-disk buffer

If the remote VictoriaLogs instance is temporarily unavailable, vlagent buffers collected logs on disk
and sends them as soon as the connection is restored.
By default, the buffer is unlimited in size and uses the entire disk space.

To cap how much disk space the buffer can use per destination, set `maxDiskUsagePerURL`:

```yaml
remoteWrite:
  - url: http://victoria-logs:9428
    maxDiskUsagePerURL: 10GB
```

When the limit is reached, the oldest buffered logs are dropped to make room for new ones.
See [replication and high availability docs](https://docs.victoriametrics.com/victorialogs/vlagent/#replication-and-high-availability) for more details.

### Kubernetes metadata

By default, vlagent attaches Pod labels to every log entry.
You can control which metadata fields are included using the following settings:

```yaml
collector:
  includePodLabels: true              # default: true
  includePodAnnotations: false        # default: false
  includeNodeLabels: false            # default: false
  includeNodeAnnotations: false       # default: false
  includeNamespaceLabels: false       # default: false
  includeNamespaceAnnotations: false  # default: false
```

Note that even when a metadata field is excluded from log entries, it is still available for filtering via `collector.excludeFilter`.
For example, you can filter by Pod annotations without storing them in every log entry:

```yaml
collector:
  includePodAnnotations: false
  excludeFilter: "kubernetes.pod_annotations.victoriametrics.com/vlagent/exclude:=true"
```

Note that vlagent does not update Node, Pod, and Namespace metadata during runtime.
If Node, Pod, or Namespace labels change, you must restart vlagent to apply those changes.

### Ignore fields

VictoriaLogs efficiently compresses repeated values, such as Pod and Node labels.
However, if you prefer not to store certain fields, you have two options.

Use `collector.ignoreFields` to drop fields on the vlagent side before sending to any destination:

```yaml
collector:
  ignoreFields:
    - request.payload*
    - password
    - token
```

Use the `VL-Ignore-Fields` header to drop fields for a specific destination only.
In this case, the fields are still sent over the network but ignored by VictoriaLogs on ingestion:

```yaml
remoteWrite:
  - url: http://victoria-logs:9428
    headers:
      VL-Ignore-Fields:
        - kubernetes.namespace_labels*
        - kubernetes.pod_labels*
```

This allows you to exclude unnecessary or sensitive fields from being ingested.
If sensitive data has already been ingested, see how to
[exclude logs from search result](https://docs.victoriametrics.com/victorialogs/security-and-lb/#access-control-inside-a-single-tenant).

### Extra fields

You have two options to attach additional fields to log entries.

Use `collector.extraFields` to add fields on the vlagent side for all destinations:

```yaml
collector:
  extraFields: '{"env":"production","cluster":"us-east1"}'
```

Use the `VL-Extra-Fields` header to add fields for a specific destination only:

```yaml
remoteWrite:
  - url: http://victoria-logs:9428
    headers:
      VL-Extra-Fields:
        zone: us-east1-c
        source: victoria-logs-collector
```

This feature lets you attach metadata to every log entry,
making it easier to filter, group, or analyze logs based on these additional attributes.

### Filtering Kubernetes logs

To skip logs from specific containers before reading them, use `collector.excludeFilter`.
The filter is applied to container metadata fields before reading the log files,
which significantly reduces CPU and I/O usage compared to filtering after ingestion.

The filter supports any [LogsQL filter expression](https://docs.victoriametrics.com/victorialogs/logsql/#filters).
See the full list of supported metadata fields in [vlagent filtering docs](https://docs.victoriametrics.com/victorialogs/vlagent/#filtering-kubernetes-logs).

Exclude vlagent's own logs:

```yaml
collector:
  excludeFilter: "kubernetes.pod_name:=%{HOSTNAME}"
```

Exclude all logs from specific namespaces:

```yaml
collector:
  excludeFilter: "kubernetes.pod_namespace:in(kube-system, monitoring)"
```

Exclude logs from specific containers:

```yaml
collector:
  excludeFilter: "kubernetes.container_name:in(config-reloader, prometheus-config-reloader) and kubernetes.pod_namespace:=monitoring"
```

Exclude Pods by annotation - add `victoriametrics.com/vlagent/exclude: "true"` to a Pod and it will be skipped:

```yaml
collector:
  excludeFilter: "kubernetes.pod_annotations.victoriametrics.com/vlagent/exclude:=true"
```

### Sending logs to external systems

By default, vlagent uses the native VictoriaLogs protocol.
To send logs to external systems, set `format: jsonline` on the corresponding `remoteWrite` entry.

You can mix destinations with different formats in the same configuration:

```yaml
remoteWrite:
  - url: http://victoria-logs:9428
    format: native
  - url: http://vector:8080
    format: jsonline
```

See [vlagent remote write format docs](https://docs.victoriametrics.com/victorialogs/vlagent/#remote-write-format) for more details.

#### Fluent Bit

Start Fluent Bit with an HTTP input on port `8080`, then configure the collector to forward logs to it.
Fluent Bit requires the `Content-Type: application/json` header:

```yaml
remoteWrite:
  - url: http://fluent-bit:8080
    format: jsonline
    headers:
      Content-Type: application/json
```

#### Vector

Start Vector with an `http_server` source listening on port `8080` with `codec: json`, then point the collector at it:

```yaml
remoteWrite:
  - url: http://vector:8080
    format: jsonline
```

#### ClickHouse

To send logs to ClickHouse, include the `INSERT` statement as the `query` HTTP parameter in the URL.
Choose the `FORMAT` clause based on your table schema:

- `JSONAsObject` - table with a single `JSON` column.
- `JSONAsString` - table with a single `String` column.
- `JSONEachRow` - table with multiple columns, each JSON field maps to a column.

```yaml
remoteWrite:
  - url: "http://clickhouse:8123?query=INSERT%20INTO%20default.logs%20FORMAT%20JSONAsObject"
    format: jsonline
    basicAuth:
      username: clickhouse-user
      passwordFile: /etc/vlagent/ch-password.txt
```

### Multitenancy

To define [tenant](https://docs.victoriametrics.com/victorialogs/#multitenancy), set `ProjectID` and `AccountID` headers as shown below:

```yaml
remoteWrite:
  - url: http://localhost:9428
    headers:
      ProjectID: 12
      AccountID: 42
```

### Basic auth

To authenticate against VictoriaLogs using basic auth, store the password in a Kubernetes Secret
and reference it via `passwordFile` to avoid exposing credentials in values:

```yaml
remoteWrite:
  - url: http://victoria-logs:9428
    basicAuth:
      username: vlagent
      passwordFile: /etc/vlagent/password.txt

extraVolumes:
  - name: auth-secret
    secret:
      secretName: vlagent-auth

extraVolumeMounts:
  - name: auth-secret
    mountPath: /etc/vlagent
    readOnly: true
```

Alternatively, pass credentials via environment variables:

```yaml
remoteWrite:
  - url: http://victoria-logs:9428

env:
  - name: VL_remoteWrite_basicAuth_password
    valueFrom:
      secretKeyRef:
        name: auth-secret
        key: VL_PASSWORD
  - name: VL_remoteWrite_basicAuth_username
    valueFrom:
      secretKeyRef:
        name: auth-secret
        key: VL_USERNAME
```

### Bearer token auth

To authenticate using a bearer token, store the token in a Kubernetes Secret
and reference it via `bearerTokenFile`. The token file is reloaded automatically when it changes:

```yaml
remoteWrite:
  - url: http://victoria-logs:9428
    bearerTokenFile: /etc/vlagent/token.txt

extraVolumes:
  - name: auth-token
    secret:
      secretName: vlagent-token

extraVolumeMounts:
  - name: auth-token
    mountPath: /etc/vlagent
    readOnly: true
```

### TLS

To enable TLS verification for the remoteWrite target, you can specify the `tls`-prefixed params inside each remoteWrite entry.

At a minimum, you should provide the `tlsCAFile` path so that the collector can verify the server's TLS certificate.
This is useful when the target endpoint uses a certificate signed by a custom or self-signed Certificate Authority (CA).

```yaml
remoteWrite:
  - url: https://victoria-logs:9428
    tlsCAFile: /etc/tls/ca.crt

extraVolumes:
  - name: tls-certs
    secret:
      secretName: tls-secret

extraVolumeMounts:
  - name: tls-certs
    mountPath: /etc/tls
    readOnly: true
```

For mutual TLS (mTLS), additionally specify the client certificate and key:

```yaml
remoteWrite:
  - url: https://victoria-logs:9428
    tlsCAFile: /etc/tls/ca.crt
    tlsCertFile: /etc/tls/client.crt
    tlsKeyFile: /etc/tls/client.key
```

If you want to disable TLS certificate verification (not recommended in production), you can set `tlsInsecureSkipVerify` to true:

```yaml
remoteWrite:
  - url: https://victoria-logs:9428
    tlsInsecureSkipVerify: true
```

## How to install

Access a Kubernetes cluster.

### Setup chart repository (can be omitted for OCI repositories)

Add a chart helm repository with follow commands:

```console
helm repo add vm https://victoriametrics.github.io/helm-charts/

helm repo update
```
List versions of `vm/victoria-logs-collector` chart available to installation:

```console
helm search repo vm/victoria-logs-collector -l
```

### Install `victoria-logs-collector` chart

Export default values of `victoria-logs-collector` chart to file `values.yaml`:

  - For HTTPS repository

    ```console
    helm show values vm/victoria-logs-collector > values.yaml
    ```
  - For OCI repository

    ```console
    helm show values oci://ghcr.io/victoriametrics/helm-charts/victoria-logs-collector > values.yaml
    ```

Change the values according to the need of the environment in ``values.yaml`` file.

> Consider setting `.Values.nameOverride` to a small value like `vlc` to avoid hitting resource name limits of 63 characters

Test the installation with command:

  - For HTTPS repository

    ```console
    helm install vlc vm/victoria-logs-collector -f values.yaml -n NAMESPACE --debug --dry-run
    ```

  - For OCI repository

    ```console
    helm install vlc oci://ghcr.io/victoriametrics/helm-charts/victoria-logs-collector -f values.yaml -n NAMESPACE --debug --dry-run
    ```

Install chart with command:

  - For HTTPS repository

    ```console
    helm install vlc vm/victoria-logs-collector -f values.yaml -n NAMESPACE
    ```

  - For OCI repository

    ```console
    helm install vlc oci://ghcr.io/victoriametrics/helm-charts/victoria-logs-collector -f values.yaml -n NAMESPACE
    ```

Get the pods lists by running this commands:

```console
kubectl get pods -A | grep 'vlc'
```

Get the application by running this command:

```console
helm list -f vlc -n NAMESPACE
```

See the history of versions of `vlc` application with command.

```console
helm history vlc -n NAMESPACE
```

## How to uninstall

Remove application with command.

```console
helm uninstall vlc -n NAMESPACE
```

## Parameters

The following tables lists the configurable parameters of the chart and their default values.

Change the values according to the need of the environment in ``victoria-logs-collector/values.yaml`` file.

<table class="helm-vars">
  <thead>
    <th class="helm-vars-key">Key</th>
    <th class="helm-vars-description">Description</th>
  </thead>
  <tbody>
    <tr id="affinity">
      <td><a href="#affinity"><pre class="chroma"><code><span class="line"><span class="cl"><span class="nt">affinity</span><span class="p">:</span><span class="w"> </span>{}</span></span></code></pre>
</a></td>
      <td><em><code>(object)</code></em><p>Pod affinity.</p>
</td>
    </tr>
    <tr id="annotations">
      <td><a href="#annotations"><pre class="chroma"><code><span class="line"><span class="cl"><span class="nt">annotations</span><span class="p">:</span><span class="w"> </span>{}</span></span></code></pre>
</a></td>
      <td><em><code>(object)</code></em><p>Annotations to be added to the deployment.</p>
</td>
    </tr>
    <tr id="collector-excludefilter">
      <td><a href="#collector-excludefilter"><pre class="chroma"><code><span class="line"><span class="cl"><span class="nt">collector.excludeFilter</span><span class="p">:</span><span class="w"> </span><span class="s2">&#34;&#34;</span></span></span></code></pre>
</a></td>
      <td><em><code>(string)</code></em><p>LogsQL filter for excluding container logs. The filter is applied to container metadata fields (e.g., kubernetes.namespace_name, kubernetes.container_name) before reading the log files. This significantly reduces CPU and I/O usage by skipping logs from unwanted containers. See <a href="https://docs.victoriametrics.com/victorialogs/vlagent/#filtering-kubernetes-logs" target="_blank">https://docs.victoriametrics.com/victorialogs/vlagent/#filtering-kubernetes-logs</a></p>
</td>
    </tr>
    <tr id="collector-extrafields">
      <td><a href="#collector-extrafields"><pre class="chroma"><code><span class="line"><span class="cl"><span class="nt">collector.extraFields</span><span class="p">:</span><span class="w"> </span><span class="s2">&#34;&#34;</span></span></span></code></pre>
</a></td>
      <td><em><code>(string)</code></em><p>Add additional fields to the log entries. If you want to add a field per remoteWrite destination, set the VL-Extra-Fields header instead.</p>
</td>
    </tr>
    <tr id="collector-ignorefields">
      <td><a href="#collector-ignorefields"><pre class="chroma"><code><span class="line"><span class="cl"><span class="nt">collector.ignoreFields</span><span class="p">:</span><span class="w"> </span><span class="p">[]</span></span></span></code></pre>
</a></td>
      <td><em><code>(list)</code></em><p>List of fields to be ignored. A field can be ended with &lsquo;*&rsquo; to ignore all fields starting with the specified prefix. If you want to exclude Kubernetes metadata fields, use includePodLabels/includePodAnnotations/includeNodeLabels/includeNodeAnnotations instead.</p>
</td>
    </tr>
    <tr id="collector-includenodeannotations">
      <td><a href="#collector-includenodeannotations"><pre class="chroma"><code><span class="line"><span class="cl"><span class="nt">collector.includeNodeAnnotations</span><span class="p">:</span><span class="w"> </span><span class="kc">false</span></span></span></code></pre>
</a></td>
      <td><em><code>(bool)</code></em><p>Include Node annotations as additional fields in the log entries. Even this setting is disabled, Node annotations are available for filtering via .Values.excludeFilter parameter.</p>
</td>
    </tr>
    <tr id="collector-includenodelabels">
      <td><a href="#collector-includenodelabels"><pre class="chroma"><code><span class="line"><span class="cl"><span class="nt">collector.includeNodeLabels</span><span class="p">:</span><span class="w"> </span><span class="kc">false</span></span></span></code></pre>
</a></td>
      <td><em><code>(bool)</code></em><p>Include Node labels as additional fields in the log entries. Even this setting is disabled, Node labels are available for filtering via .Values.excludeFilter parameter.</p>
</td>
    </tr>
    <tr id="collector-includepodannotations">
      <td><a href="#collector-includepodannotations"><pre class="chroma"><code><span class="line"><span class="cl"><span class="nt">collector.includePodAnnotations</span><span class="p">:</span><span class="w"> </span><span class="kc">false</span></span></span></code></pre>
</a></td>
      <td><em><code>(bool)</code></em><p>Include Pod annotations as additional fields in the log entries. Even this setting is disabled, Pod annotations are available for filtering via .Values.excludeFilter parameter.</p>
</td>
    </tr>
    <tr id="collector-includepodlabels">
      <td><a href="#collector-includepodlabels"><pre class="chroma"><code><span class="line"><span class="cl"><span class="nt">collector.includePodLabels</span><span class="p">:</span><span class="w"> </span><span class="kc">true</span></span></span></code></pre>
</a></td>
      <td><em><code>(bool)</code></em><p>Include Pod labels as additional fields in the log entries. Even this setting is disabled, Pod labels are available for filtering via .Values.excludeFilter parameter.</p>
</td>
    </tr>
    <tr id="collector-msgfield">
      <td><a href="#collector-msgfield"><pre class="chroma"><code><span class="line"><span class="cl"><span class="nt">collector.msgField</span><span class="p">:</span><span class="w">
</span></span></span><span class="line"><span class="cl"><span class="w">    </span>- <span class="l">message</span><span class="w">
</span></span></span><span class="line"><span class="cl"><span class="w">    </span>- <span class="l">msg</span><span class="w">
</span></span></span><span class="line"><span class="cl"><span class="w">    </span>- <span class="l">log</span></span></span></code></pre>
</a></td>
      <td><em><code>(list)</code></em><p>List of fields to be used as <code>_msg</code> field. The first found field will be used. See <a href="https://docs.victoriametrics.com/victorialogs/keyconcepts/#message-field" target="_blank">https://docs.victoriametrics.com/victorialogs/keyconcepts/#message-field</a> for more details. Add more fields if you see the &lsquo;missing _msg field&rsquo; message in VictoriaLogs. If you set an empty list, then the default message fields will be used.</p>
</td>
    </tr>
    <tr id="collector-streamfields">
      <td><a href="#collector-streamfields"><pre class="chroma"><code><span class="line"><span class="cl"><span class="nt">collector.streamFields</span><span class="p">:</span><span class="w">
</span></span></span><span class="line"><span class="cl"><span class="w">    </span>- <span class="l">kubernetes.container_name</span><span class="w">
</span></span></span><span class="line"><span class="cl"><span class="w">    </span>- <span class="l">kubernetes.pod_name</span><span class="w">
</span></span></span><span class="line"><span class="cl"><span class="w">    </span>- <span class="l">kubernetes.pod_namespace</span></span></span></code></pre>
</a></td>
      <td><em><code>(list)</code></em><p>List of fields to be used as <code>_stream</code> field. See <a href="https://docs.victoriametrics.com/victorialogs/keyconcepts/#stream-fields" target="_blank">https://docs.victoriametrics.com/victorialogs/keyconcepts/#stream-fields</a> for more details. It is recommended to leave this with default values. If you set an empty list, then the default stream fields will be used.</p>
</td>
    </tr>
    <tr id="collector-timefield">
      <td><a href="#collector-timefield"><pre class="chroma"><code><span class="line"><span class="cl"><span class="nt">collector.timeField</span><span class="p">:</span><span class="w">
</span></span></span><span class="line"><span class="cl"><span class="w">    </span>- <span class="l">time</span><span class="w">
</span></span></span><span class="line"><span class="cl"><span class="w">    </span>- <span class="l">ts</span><span class="w">
</span></span></span><span class="line"><span class="cl"><span class="w">    </span>- <span class="l">timestamp</span></span></span></code></pre>
</a></td>
      <td><em><code>(list)</code></em><p>List of fields to be used as <code>_time</code> field. The first found field will be used. See <a href="https://docs.victoriametrics.com/victorialogs/keyconcepts/#time-field" target="_blank">https://docs.victoriametrics.com/victorialogs/keyconcepts/#time-field</a> for more details. If you set an empty list, then the default time fields will be used.</p>
</td>
    </tr>
    <tr id="defaultvolumemounts">
      <td><a href="#defaultvolumemounts"><pre class="chroma"><code><span class="line"><span class="cl"><span class="nt">defaultVolumeMounts</span><span class="p">:</span><span class="w">
</span></span></span><span class="line"><span class="cl"><span class="w">    </span>- <span class="nt">mountPath</span><span class="p">:</span><span class="w"> </span><span class="l">/var/log</span><span class="w">
</span></span></span><span class="line"><span class="cl"><span class="w">      </span><span class="nt">name</span><span class="p">:</span><span class="w"> </span><span class="l">varlog</span><span class="w">
</span></span></span><span class="line"><span class="cl"><span class="w">      </span><span class="nt">readOnly</span><span class="p">:</span><span class="w"> </span><span class="kc">true</span><span class="w">
</span></span></span><span class="line"><span class="cl"><span class="w">    </span>- <span class="nt">mountPath</span><span class="p">:</span><span class="w"> </span><span class="l">/var/lib</span><span class="w">
</span></span></span><span class="line"><span class="cl"><span class="w">      </span><span class="nt">name</span><span class="p">:</span><span class="w"> </span><span class="l">varlib</span><span class="w">
</span></span></span><span class="line"><span class="cl"><span class="w">      </span><span class="nt">readOnly</span><span class="p">:</span><span class="w"> </span><span class="kc">true</span></span></span></code></pre>
</a></td>
      <td><em><code>(list)</code></em><p>Mounts for defaultVolumes.</p>
</td>
    </tr>
    <tr id="defaultvolumes">
      <td><a href="#defaultvolumes"><pre class="chroma"><code><span class="line"><span class="cl"><span class="nt">defaultVolumes</span><span class="p">:</span><span class="w">
</span></span></span><span class="line"><span class="cl"><span class="w">    </span>- <span class="nt">hostPath</span><span class="p">:</span><span class="w">
</span></span></span><span class="line"><span class="cl"><span class="w">        </span><span class="nt">path</span><span class="p">:</span><span class="w"> </span><span class="l">/var/log</span><span class="w">
</span></span></span><span class="line"><span class="cl"><span class="w">      </span><span class="nt">name</span><span class="p">:</span><span class="w"> </span><span class="l">varlog</span><span class="w">
</span></span></span><span class="line"><span class="cl"><span class="w">    </span>- <span class="nt">hostPath</span><span class="p">:</span><span class="w">
</span></span></span><span class="line"><span class="cl"><span class="w">        </span><span class="nt">path</span><span class="p">:</span><span class="w"> </span><span class="l">/var/lib</span><span class="w">
</span></span></span><span class="line"><span class="cl"><span class="w">      </span><span class="nt">name</span><span class="p">:</span><span class="w"> </span><span class="l">varlib</span></span></span></code></pre>
</a></td>
      <td><em><code>(list)</code></em><p>Default volumes that are mounted into Pods. In most cases, these should not be changed. Use <code>extraVolumes</code>/<code>extraVolumeMounts</code> for additional custom volumes.</p>
</td>
    </tr>
    <tr id="dnsconfig">
      <td><a href="#dnsconfig"><pre class="chroma"><code><span class="line"><span class="cl"><span class="nt">dnsConfig</span><span class="p">:</span><span class="w"> </span>{}</span></span></code></pre>
</a></td>
      <td><em><code>(object)</code></em><p>Custom DNS config for pod. Details are <a href="https://kubernetes.io/docs/concepts/services-networking/dns-pod-service/#pod-s-dns-config" target="_blank">here</a></p>
</td>
    </tr>
    <tr id="env">
      <td><a href="#env"><pre class="chroma"><code><span class="line"><span class="cl"><span class="nt">env</span><span class="p">:</span><span class="w"> </span><span class="p">[]</span></span></span></code></pre>
</a></td>
      <td><em><code>(list)</code></em><p>Environment variables (ex.: secret tokens).</p>
</td>
    </tr>
    <tr id="extraargs">
      <td><a href="#extraargs"><pre class="chroma"><code><span class="line"><span class="cl"><span class="nt">extraArgs</span><span class="p">:</span><span class="w">
</span></span></span><span class="line"><span class="cl"><span class="w">    </span><span class="nt">envflag.enable</span><span class="p">:</span><span class="w"> </span><span class="kc">true</span><span class="w">
</span></span></span><span class="line"><span class="cl"><span class="w">    </span><span class="nt">envflag.prefix</span><span class="p">:</span><span class="w"> </span><span class="l">VL_</span><span class="w">
</span></span></span><span class="line"><span class="cl"><span class="w">    </span><span class="nt">loggerFormat</span><span class="p">:</span><span class="w"> </span><span class="l">json</span><span class="w">
</span></span></span><span class="line"><span class="cl"><span class="w">    </span><span class="nt">tmpDataPath</span><span class="p">:</span><span class="w"> </span><span class="l">/var/lib/vl-collector</span></span></span></code></pre>
</a></td>
      <td><em><code>(object)</code></em><p>vlagent extra command line arguments.</p>
</td>
    </tr>
    <tr id="extraobjects">
      <td><a href="#extraobjects"><pre class="chroma"><code><span class="line"><span class="cl"><span class="nt">extraObjects</span><span class="p">:</span><span class="w"> </span><span class="p">[]</span></span></span></code></pre>
</a></td>
      <td><em><code>(list)</code></em><p>Add extra specs dynamically to this chart</p>
</td>
    </tr>
    <tr id="extravolumemounts">
      <td><a href="#extravolumemounts"><pre class="chroma"><code><span class="line"><span class="cl"><span class="nt">extraVolumeMounts</span><span class="p">:</span><span class="w"> </span><span class="p">[]</span></span></span></code></pre>
</a></td>
      <td><em><code>(list)</code></em><p>Extra Volume Mounts for the container.</p>
</td>
    </tr>
    <tr id="extravolumes">
      <td><a href="#extravolumes"><pre class="chroma"><code><span class="line"><span class="cl"><span class="nt">extraVolumes</span><span class="p">:</span><span class="w"> </span><span class="p">[]</span></span></span></code></pre>
</a></td>
      <td><em><code>(list)</code></em><p>Extra Volumes for the Pod.</p>
</td>
    </tr>
    <tr id="filecollector">
      <td><a href="#filecollector"><pre class="chroma"><code><span class="line"><span class="cl"><span class="nt">fileCollector</span><span class="p">:</span><span class="w"> </span><span class="p">[]</span></span></span></code></pre>
</a></td>
      <td><em><code>(list)</code></em><p>List of file-based log collection configurations. Each entry maps directly to a set of positional <code>-fileCollector.*</code> flags, preserving the order needed for multiple glob patterns to pair correctly with their extra fields. See <a href="https://docs.victoriametrics.com/victorialogs/vlagent/#multiple-glob-patterns-and-ordering" target="_blank">https://docs.victoriametrics.com/victorialogs/vlagent/#multiple-glob-patterns-and-ordering</a></p>
</td>
    </tr>
    <tr id="fullnameoverride">
      <td><a href="#fullnameoverride"><pre class="chroma"><code><span class="line"><span class="cl"><span class="nt">fullnameOverride</span><span class="p">:</span><span class="w"> </span><span class="s2">&#34;&#34;</span></span></span></code></pre>
</a></td>
      <td><em><code>(string)</code></em><p>Override resources fullname.</p>
</td>
    </tr>
    <tr id="global-cluster-dnsdomain">
      <td><a href="#global-cluster-dnsdomain"><pre class="chroma"><code><span class="line"><span class="cl"><span class="nt">global.cluster.dnsDomain</span><span class="p">:</span><span class="w"> </span><span class="l">cluster.local.</span></span></span></code></pre>
</a></td>
      <td><em><code>(string)</code></em><p>K8s cluster domain suffix, used for building storage Pods&rsquo; FQDN. Details are <a href="https://kubernetes.io/docs/tasks/administer-cluster/dns-custom-nameservers/" target="_blank">here</a>.</p>
</td>
    </tr>
    <tr id="global-compatibility">
      <td><a href="#global-compatibility"><pre class="chroma"><code><span class="line"><span class="cl"><span class="nt">global.compatibility</span><span class="p">:</span><span class="w">
</span></span></span><span class="line"><span class="cl"><span class="w">    </span><span class="nt">openshift</span><span class="p">:</span><span class="w">
</span></span></span><span class="line"><span class="cl"><span class="w">        </span><span class="nt">adaptSecurityContext</span><span class="p">:</span><span class="w"> </span><span class="l">auto</span></span></span></code></pre>
</a></td>
      <td><em><code>(object)</code></em><p>Openshift security context compatibility configuration</p>
</td>
    </tr>
    <tr id="global-extraannotations">
      <td><a href="#global-extraannotations"><pre class="chroma"><code><span class="line"><span class="cl"><span class="nt">global.extraAnnotations</span><span class="p">:</span><span class="w"> </span>{}</span></span></code></pre>
</a></td>
      <td><em><code>(object)</code></em><p>Annotations added to all resources</p>
</td>
    </tr>
    <tr id="global-extralabels">
      <td><a href="#global-extralabels"><pre class="chroma"><code><span class="line"><span class="cl"><span class="nt">global.extraLabels</span><span class="p">:</span><span class="w"> </span>{}</span></span></code></pre>
</a></td>
      <td><em><code>(object)</code></em><p>Labels added to all resources</p>
</td>
    </tr>
    <tr id="global-image-registry">
      <td><a href="#global-image-registry"><pre class="chroma"><code><span class="line"><span class="cl"><span class="nt">global.image.registry</span><span class="p">:</span><span class="w"> </span><span class="s2">&#34;&#34;</span></span></span></code></pre>
</a></td>
      <td><em><code>(string)</code></em><p>Image registry, that can be shared across multiple helm charts</p>
</td>
    </tr>
    <tr id="global-imagepullsecrets">
      <td><a href="#global-imagepullsecrets"><pre class="chroma"><code><span class="line"><span class="cl"><span class="nt">global.imagePullSecrets</span><span class="p">:</span><span class="w"> </span><span class="p">[]</span></span></span></code></pre>
</a></td>
      <td><em><code>(list)</code></em><p>Image pull secrets, that can be shared across multiple helm charts</p>
</td>
    </tr>
    <tr id="image-pullpolicy">
      <td><a href="#image-pullpolicy"><pre class="chroma"><code><span class="line"><span class="cl"><span class="nt">image.pullPolicy</span><span class="p">:</span><span class="w"> </span><span class="l">IfNotPresent</span></span></span></code></pre>
</a></td>
      <td><em><code>(string)</code></em><p>Image pull policy.</p>
</td>
    </tr>
    <tr id="image-registry">
      <td><a href="#image-registry"><pre class="chroma"><code><span class="line"><span class="cl"><span class="nt">image.registry</span><span class="p">:</span><span class="w"> </span><span class="s2">&#34;&#34;</span></span></span></code></pre>
</a></td>
      <td><em><code>(string)</code></em><p>Image registry.</p>
</td>
    </tr>
    <tr id="image-repository">
      <td><a href="#image-repository"><pre class="chroma"><code><span class="line"><span class="cl"><span class="nt">image.repository</span><span class="p">:</span><span class="w"> </span><span class="l">victoriametrics/vlagent</span></span></span></code></pre>
</a></td>
      <td><em><code>(string)</code></em><p>Image repository.</p>
</td>
    </tr>
    <tr id="image-tag">
      <td><a href="#image-tag"><pre class="chroma"><code><span class="line"><span class="cl"><span class="nt">image.tag</span><span class="p">:</span><span class="w"> </span><span class="s2">&#34;&#34;</span></span></span></code></pre>
</a></td>
      <td><em><code>(string)</code></em><p>Image tag.</p>
</td>
    </tr>
    <tr id="image-variant">
      <td><a href="#image-variant"><pre class="chroma"><code><span class="line"><span class="cl"><span class="nt">image.variant</span><span class="p">:</span><span class="w"> </span><span class="s2">&#34;&#34;</span></span></span></code></pre>
</a></td>
      <td><em><code>(string)</code></em><p>Image tag suffix, which is appended to <code>Chart.AppVersion</code> if no <code>server.image.tag</code> is defined.</p>
</td>
    </tr>
    <tr id="license">
      <td><a href="#license"><pre class="chroma"><code><span class="line"><span class="cl"><span class="nt">license</span><span class="p">:</span><span class="w">
</span></span></span><span class="line"><span class="cl"><span class="w">    </span><span class="nt">key</span><span class="p">:</span><span class="w"> </span><span class="s2">&#34;&#34;</span><span class="w">
</span></span></span><span class="line"><span class="cl"><span class="w">    </span><span class="nt">secret</span><span class="p">:</span><span class="w">
</span></span></span><span class="line"><span class="cl"><span class="w">        </span><span class="nt">key</span><span class="p">:</span><span class="w"> </span><span class="s2">&#34;&#34;</span><span class="w">
</span></span></span><span class="line"><span class="cl"><span class="w">        </span><span class="nt">name</span><span class="p">:</span><span class="w"> </span><span class="s2">&#34;&#34;</span></span></span></code></pre>
</a></td>
      <td><em><code>(object)</code></em><p>Enterprise license key configuration for VictoriaLogs enterprise. Required only for VictoriaLogs enterprise. Check docs <a href="https://docs.victoriametrics.com/victoriametrics/enterprise/#victorialogs-enterprise-features" target="_blank">here</a>, Request a trial license <a href="https://victoriametrics.com/products/enterprise/trial/" target="_blank">here</a>.</p>
</td>
    </tr>
    <tr id="license-key">
      <td><a href="#license-key"><pre class="chroma"><code><span class="line"><span class="cl"><span class="nt">license.key</span><span class="p">:</span><span class="w"> </span><span class="s2">&#34;&#34;</span></span></span></code></pre>
</a></td>
      <td><em><code>(string)</code></em><p>License key.</p>
</td>
    </tr>
    <tr id="license-secret">
      <td><a href="#license-secret"><pre class="chroma"><code><span class="line"><span class="cl"><span class="nt">license.secret</span><span class="p">:</span><span class="w">
</span></span></span><span class="line"><span class="cl"><span class="w">    </span><span class="nt">key</span><span class="p">:</span><span class="w"> </span><span class="s2">&#34;&#34;</span><span class="w">
</span></span></span><span class="line"><span class="cl"><span class="w">    </span><span class="nt">name</span><span class="p">:</span><span class="w"> </span><span class="s2">&#34;&#34;</span></span></span></code></pre>
</a></td>
      <td><em><code>(object)</code></em><p>Use existing secret with license key.</p>
</td>
    </tr>
    <tr id="license-secret-key">
      <td><a href="#license-secret-key"><pre class="chroma"><code><span class="line"><span class="cl"><span class="nt">license.secret.key</span><span class="p">:</span><span class="w"> </span><span class="s2">&#34;&#34;</span></span></span></code></pre>
</a></td>
      <td><em><code>(string)</code></em><p>Key in secret with license key.</p>
</td>
    </tr>
    <tr id="license-secret-name">
      <td><a href="#license-secret-name"><pre class="chroma"><code><span class="line"><span class="cl"><span class="nt">license.secret.name</span><span class="p">:</span><span class="w"> </span><span class="s2">&#34;&#34;</span></span></span></code></pre>
</a></td>
      <td><em><code>(string)</code></em><p>Existing secret name.</p>
</td>
    </tr>
    <tr id="nameoverride">
      <td><a href="#nameoverride"><pre class="chroma"><code><span class="line"><span class="cl"><span class="nt">nameOverride</span><span class="p">:</span><span class="w"> </span><span class="s2">&#34;&#34;</span></span></span></code></pre>
</a></td>
      <td><em><code>(string)</code></em><p>Override chart name.</p>
</td>
    </tr>
    <tr id="networkpolicy">
      <td><a href="#networkpolicy"><pre class="chroma"><code><span class="line"><span class="cl"><span class="nt">networkPolicy</span><span class="p">:</span><span class="w">
</span></span></span><span class="line"><span class="cl"><span class="w">    </span><span class="nt">annotations</span><span class="p">:</span><span class="w"> </span>{}<span class="w">
</span></span></span><span class="line"><span class="cl"><span class="w">    </span><span class="nt">egress</span><span class="p">:</span><span class="w"> </span><span class="p">[]</span><span class="w">
</span></span></span><span class="line"><span class="cl"><span class="w">    </span><span class="nt">enabled</span><span class="p">:</span><span class="w"> </span><span class="kc">false</span><span class="w">
</span></span></span><span class="line"><span class="cl"><span class="w">    </span><span class="nt">ingress</span><span class="p">:</span><span class="w"> </span><span class="p">[]</span><span class="w">
</span></span></span><span class="line"><span class="cl"><span class="w">    </span><span class="nt">labels</span><span class="p">:</span><span class="w"> </span>{}</span></span></code></pre>
</a></td>
      <td><em><code>(object)</code></em><p>See <code>kubectl explain networkpolicy.spec</code> for more. Details are <a href="https://kubernetes.io/docs/concepts/services-networking/network-policies/" target="_blank">here</a></p>
</td>
    </tr>
    <tr id="networkpolicy-annotations">
      <td><a href="#networkpolicy-annotations"><pre class="chroma"><code><span class="line"><span class="cl"><span class="nt">networkPolicy.annotations</span><span class="p">:</span><span class="w"> </span>{}</span></span></code></pre>
</a></td>
      <td><em><code>(object)</code></em><p>Extra annotations for NetworkPolicy</p>
</td>
    </tr>
    <tr id="networkpolicy-egress">
      <td><a href="#networkpolicy-egress"><pre class="chroma"><code><span class="line"><span class="cl"><span class="nt">networkPolicy.egress</span><span class="p">:</span><span class="w"> </span><span class="p">[]</span></span></span></code></pre>
</a></td>
      <td><em><code>(list)</code></em><p>Egress rules</p>
</td>
    </tr>
    <tr id="networkpolicy-ingress">
      <td><a href="#networkpolicy-ingress"><pre class="chroma"><code><span class="line"><span class="cl"><span class="nt">networkPolicy.ingress</span><span class="p">:</span><span class="w"> </span><span class="p">[]</span></span></span></code></pre>
</a></td>
      <td><em><code>(list)</code></em><p>Ingress rules</p>
</td>
    </tr>
    <tr id="networkpolicy-labels">
      <td><a href="#networkpolicy-labels"><pre class="chroma"><code><span class="line"><span class="cl"><span class="nt">networkPolicy.labels</span><span class="p">:</span><span class="w"> </span>{}</span></span></code></pre>
</a></td>
      <td><em><code>(object)</code></em><p>Extra labels for NetworkPolicy</p>
</td>
    </tr>
    <tr id="nodeselector">
      <td><a href="#nodeselector"><pre class="chroma"><code><span class="line"><span class="cl"><span class="nt">nodeSelector</span><span class="p">:</span><span class="w"> </span>{}</span></span></code></pre>
</a></td>
      <td><em><code>(object)</code></em><p>Pod&rsquo;s node selector. Details are <a href="https://kubernetes.io/docs/concepts/scheduling-eviction/assign-pod-node/#nodeselector" target="_blank">here</a>.</p>
</td>
    </tr>
    <tr id="persistence-volume">
      <td><a href="#persistence-volume"><pre class="chroma"><code><span class="line"><span class="cl"><span class="nt">persistence.volume</span><span class="p">:</span><span class="w"> </span>{}</span></span></code></pre>
</a></td>
      <td><em><code>(object)</code></em></td>
    </tr>
    <tr id="podannotations">
      <td><a href="#podannotations"><pre class="chroma"><code><span class="line"><span class="cl"><span class="nt">podAnnotations</span><span class="p">:</span><span class="w"> </span>{}</span></span></code></pre>
</a></td>
      <td><em><code>(object)</code></em><p>Annotations to be added to Pod.</p>
</td>
    </tr>
    <tr id="podlabels">
      <td><a href="#podlabels"><pre class="chroma"><code><span class="line"><span class="cl"><span class="nt">podLabels</span><span class="p">:</span><span class="w"> </span>{}</span></span></code></pre>
</a></td>
      <td><em><code>(object)</code></em><p>Extra labels for Pods only.</p>
</td>
    </tr>
    <tr id="podmonitor-annotations">
      <td><a href="#podmonitor-annotations"><pre class="chroma"><code><span class="line"><span class="cl"><span class="nt">podMonitor.annotations</span><span class="p">:</span><span class="w"> </span>{}</span></span></code></pre>
</a></td>
      <td><em><code>(object)</code></em><p>PodMonitor annotations</p>
</td>
    </tr>
    <tr id="podmonitor-enabled">
      <td><a href="#podmonitor-enabled"><pre class="chroma"><code><span class="line"><span class="cl"><span class="nt">podMonitor.enabled</span><span class="p">:</span><span class="w"> </span><span class="kc">false</span></span></span></code></pre>
</a></td>
      <td><em><code>(bool)</code></em><p>Enable PodMonitor</p>
</td>
    </tr>
    <tr id="podmonitor-extralabels">
      <td><a href="#podmonitor-extralabels"><pre class="chroma"><code><span class="line"><span class="cl"><span class="nt">podMonitor.extraLabels</span><span class="p">:</span><span class="w"> </span>{}</span></span></code></pre>
</a></td>
      <td><em><code>(object)</code></em><p>PodMonitor labels</p>
</td>
    </tr>
    <tr id="podmonitor-interval">
      <td><a href="#podmonitor-interval"><pre class="chroma"><code><span class="line"><span class="cl"><span class="nt">podMonitor.interval</span><span class="p">:</span><span class="w"> </span><span class="s2">&#34;&#34;</span></span></span></code></pre>
</a></td>
      <td><em><code>(string)</code></em><p>Scrape interval, for example 30s</p>
</td>
    </tr>
    <tr id="podmonitor-metricrelabelings">
      <td><a href="#podmonitor-metricrelabelings"><pre class="chroma"><code><span class="line"><span class="cl"><span class="nt">podMonitor.metricRelabelings</span><span class="p">:</span><span class="w"> </span><span class="p">[]</span></span></span></code></pre>
</a></td>
      <td><em><code>(list)</code></em><p>Metric relabeling rules</p>
</td>
    </tr>
    <tr id="podmonitor-relabelings">
      <td><a href="#podmonitor-relabelings"><pre class="chroma"><code><span class="line"><span class="cl"><span class="nt">podMonitor.relabelings</span><span class="p">:</span><span class="w"> </span><span class="p">[]</span></span></span></code></pre>
</a></td>
      <td><em><code>(list)</code></em><p>Relabeling rules</p>
</td>
    </tr>
    <tr id="podmonitor-scrapetimeout">
      <td><a href="#podmonitor-scrapetimeout"><pre class="chroma"><code><span class="line"><span class="cl"><span class="nt">podMonitor.scrapeTimeout</span><span class="p">:</span><span class="w"> </span><span class="s2">&#34;&#34;</span></span></span></code></pre>
</a></td>
      <td><em><code>(string)</code></em><p>Scrape timeout, for example 10s</p>
</td>
    </tr>
    <tr id="podmonitor-vm">
      <td><a href="#podmonitor-vm"><pre class="chroma"><code><span class="line"><span class="cl"><span class="nt">podMonitor.vm</span><span class="p">:</span><span class="w"> </span><span class="kc">false</span></span></span></code></pre>
</a></td>
      <td><em><code>(bool)</code></em><p>Whether to use <a href="https://docs.victoriametrics.com/operator/resources/vmpodscrape/" target="_blank">VMPodScrape</a> from VM operator instead of PodMonitor.</p>
</td>
    </tr>
    <tr id="podsecuritycontext">
      <td><a href="#podsecuritycontext"><pre class="chroma"><code><span class="line"><span class="cl"><span class="nt">podSecurityContext</span><span class="p">:</span><span class="w">
</span></span></span><span class="line"><span class="cl"><span class="w">    </span><span class="nt">enabled</span><span class="p">:</span><span class="w"> </span><span class="kc">true</span></span></span></code></pre>
</a></td>
      <td><em><code>(object)</code></em><p>Security context to be added to Pod.</p>
</td>
    </tr>
    <tr id="priorityclassname">
      <td><a href="#priorityclassname"><pre class="chroma"><code><span class="line"><span class="cl"><span class="nt">priorityClassName</span><span class="p">:</span><span class="w"> </span><span class="s2">&#34;&#34;</span></span></span></code></pre>
</a></td>
      <td><em><code>(string)</code></em><p>Priority class to be assigned to the Pod(s).</p>
</td>
    </tr>
    <tr id="remotewrite">
      <td><a href="#remotewrite"><pre class="chroma"><code><span class="line"><span class="cl"><span class="nt">remoteWrite</span><span class="p">:</span><span class="w"> </span><span class="p">[]</span></span></span></code></pre>
</a></td>
      <td><em><code>(list)</code></em><p>List of log destinations. Logs will be replicated to all listed destinations.  If the url path is not specified, the logs will be sent to the /insert/native endpoint.</p>
</td>
    </tr>
    <tr id="resources">
      <td><a href="#resources"><pre class="chroma"><code><span class="line"><span class="cl"><span class="nt">resources</span><span class="p">:</span><span class="w"> </span><span class="kc">null</span></span></span></code></pre>
</a></td>
      <td><em><code>(string)</code></em></td>
    </tr>
    <tr id="runtimeclassname">
      <td><a href="#runtimeclassname"><pre class="chroma"><code><span class="line"><span class="cl"><span class="nt">runtimeClassName</span><span class="p">:</span><span class="w"> </span><span class="s2">&#34;&#34;</span></span></span></code></pre>
</a></td>
      <td><em><code>(string)</code></em><p>Name of the <a href="https://kubernetes.io/docs/concepts/containers/runtime-class/" target="_blank">RuntimeClass</a> used to run the pod, e.g. &ldquo;gvisor&rdquo;</p>
</td>
    </tr>
    <tr id="securitycontext">
      <td><a href="#securitycontext"><pre class="chroma"><code><span class="line"><span class="cl"><span class="nt">securityContext</span><span class="p">:</span><span class="w">
</span></span></span><span class="line"><span class="cl"><span class="w">    </span><span class="nt">enabled</span><span class="p">:</span><span class="w"> </span><span class="kc">true</span></span></span></code></pre>
</a></td>
      <td><em><code>(object)</code></em><p>Security context to be added to Pod&rsquo;s containers.</p>
</td>
    </tr>
    <tr id="serviceaccount">
      <td><a href="#serviceaccount"><pre class="chroma"><code><span class="line"><span class="cl"><span class="nt">serviceAccount</span><span class="p">:</span><span class="w">
</span></span></span><span class="line"><span class="cl"><span class="w">    </span><span class="nt">annotations</span><span class="p">:</span><span class="w"> </span>{}<span class="w">
</span></span></span><span class="line"><span class="cl"><span class="w">    </span><span class="nt">automount</span><span class="p">:</span><span class="w"> </span><span class="kc">true</span><span class="w">
</span></span></span><span class="line"><span class="cl"><span class="w">    </span><span class="nt">name</span><span class="p">:</span><span class="w"> </span><span class="s2">&#34;&#34;</span></span></span></code></pre>
</a></td>
      <td><em><code>(object)</code></em><p>Service account is needed to enrich logs with Pod metadata using Kubernetes API.</p>
</td>
    </tr>
    <tr id="tolerations">
      <td><a href="#tolerations"><pre class="chroma"><code><span class="line"><span class="cl"><span class="nt">tolerations</span><span class="p">:</span><span class="w"> </span><span class="p">[]</span></span></span></code></pre>
</a></td>
      <td><em><code>(list)</code></em><p>Node tolerations for server scheduling to nodes with taints. Details are <a href="https://kubernetes.io/docs/concepts/configuration/assign-pod-node/" target="_blank">here</a>.</p>
</td>
    </tr>
    <tr id="topologyspreadconstraints">
      <td><a href="#topologyspreadconstraints"><pre class="chroma"><code><span class="line"><span class="cl"><span class="nt">topologySpreadConstraints</span><span class="p">:</span><span class="w"> </span><span class="kc">null</span></span></span></code></pre>
</a></td>
      <td><em><code>(string)</code></em><p>Pod topologySpreadConstraints.</p>
</td>
    </tr>
    <tr id="updatestrategy">
      <td><a href="#updatestrategy"><pre class="chroma"><code><span class="line"><span class="cl"><span class="nt">updateStrategy</span><span class="p">:</span><span class="w"> </span>{}</span></span></code></pre>
</a></td>
      <td><em><code>(object)</code></em><p>DaemonSet update strategy.</p>
</td>
    </tr>
  </tbody>
</table>

