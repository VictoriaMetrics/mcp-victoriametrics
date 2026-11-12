

---
weight: 2
title: VictoriaLogs Agent
description: "Deploy vlagent on Kubernetes."
menu:
  docs:
    parent: helm
    weight: 2
    identifier: helm-victorialogs-agent
url: /helm/victoria-logs-agent
aliases:
  - /helm/victorialogs-agent/index.html
  - /helm/victorialogs-agent/
tags:
  - logs
  - kubernetes
---

![Version](https://img.shields.io/badge/0.2.10-gray?logo=Helm&labelColor=gray&link=https%3A%2F%2Fdocs.victoriametrics.com%2Fhelm%2Fvictoria-logs-agent%2Fchangelog%2F%23v0210)
![ArtifactHub](https://img.shields.io/badge/ArtifactHub-informational?logoColor=white&color=417598&logo=artifacthub&link=https%3A%2F%2Fartifacthub.io%2Fpackages%2Fhelm%2Fvictoriametrics%2Fvictoria-logs-agent)
![License](https://img.shields.io/github/license/VictoriaMetrics/helm-charts?labelColor=green&label=&link=https%3A%2F%2Fgithub.com%2FVictoriaMetrics%2Fhelm-charts%2Fblob%2Fmaster%2FLICENSE)
![Slack](https://img.shields.io/badge/Join-4A154B?logo=slack&link=https%3A%2F%2Fslack.victoriametrics.com)
![X](https://img.shields.io/twitter/follow/VictoriaMetrics?style=flat&label=Follow&color=black&logo=x&labelColor=black&link=https%3A%2F%2Fx.com%2FVictoriaMetrics)
![Reddit](https://img.shields.io/reddit/subreddit-subscribers/VictoriaMetrics?style=flat&label=Join&labelColor=red&logoColor=white&logo=reddit&link=https%3A%2F%2Fwww.reddit.com%2Fr%2FVictoriaMetrics)

VictoriaLogs Agent - accepts logs from various protocols and replicates them across multiple VictoriaLogs instances.

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

This Helm chart deploys the log receiver agent as a StatefulSet.
It receives logs using [supported protocols](https://docs.victoriametrics.com/victorialogs/data-ingestion/)
(for example: syslog, Elasticsearch, OpenTelemetry) from cluster workloads and forwards them to the configured VictoriaLogs destinations.
If more than one destination is specified, collected logs are replicated to all configured destinations.

- To quickly install single-node version of VictoriaLogs and `victoria-logs-agent`, see [these docs](https://docs.victoriametrics.com/helm/victoria-logs-single/#quick-start).
- To start with a VictoriaLogs cluster and `victoria-logs-agent`, see [these docs](https://docs.victoriametrics.com/helm/victoria-logs-cluster/#quick-start).

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

### Basic auth

If you need to use basic auth, define the needed flags via environment variables as shown below:

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

or mount secrets into pod and reference them in remoteWrite section:

```
remoteWrite:
  - url: http://victoria-logs:9428
    basicAuth.username: <name>
    basicAuth.passwordFile: /path/to/password

extraVolumes:
  - name: basic-auth-secret
    secret:
      secretName: basic-auth-secret

extraVolumeMounts:
  - name: basic-auth-secret
    mountPath: /path/to
    readOnly: true
```

### Multitenancy

To define [tenant](https://docs.victoriametrics.com/victorialogs/#multitenancy), use `ProjectID` and `AccountID` headers as shown below:

```yaml
remoteWrite:
  - url: http://localhost:9428
    headers:
      ProjectID: 12
      AccountID: 42
```

### TLS

To enable TLS verification for the remoteWrite target, you can specify the `tls`-prefixed flags inside each remoteWrite entry.

At a minimum, you should provide the `tlsCAFile` path so that the agent can verify the server's TLS certificate.
This is useful when the target endpoint uses a certificate signed by a custom or self-signed Certificate Authority (CA).

```yaml
remoteWrite:
  - url: https://victoria-logs:9428
    tlsCAFile: "/etc/tls/ca.crt"

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
    tlsCAFile: "/etc/tls/ca.crt"
    tlsCertFile: "/etc/tls/client.crt"
    tlsKeyFile: "/etc/tls/client.key"
```

If you want to disable TLS certificate verification (not recommended in production), you can set `tlsInsecureSkipVerify` to true.

This will skip verification of the server's certificate and allow connecting to targets with self-signed or invalid certificates.

```yaml
remoteWrite:
  - url: https://victoria-logs:9428
    tlsInsecureSkipVerify: true
```

### Ignore fields

VictoriaLogs efficiently compresses repeated values, such as pod and node labels.
However, if you prefer not to store certain fields, you can ignore them using the `VL-Ignore-Fields` header.
For example:

```yaml
remoteWrite:
  - url: http://victoria-logs:9428
    headers:
      VL-Ignore-Fields:
        - file
        - kubernetes.container_id
        - kubernetes.pod_annotations*
        - kubernetes.node_labels*
        - kubernetes.namespace_labels*
        - kubernetes.pod_labels*
```

This allows you to exclude unnecessary or sensitive fields from being ingested.
If sensitive data has already been ingested, see how to
[exclude logs from search result](https://docs.victoriametrics.com/victorialogs/security-and-lb/#access-control-inside-a-single-tenant).

### Extra fields

You can add custom fields to your logs by including the `VL-Extra-Fields` header in `headers` section of remote write configuration.
For example:

```yaml
remoteWrite:
  - url: http://victoria-logs:9428/
    headers:
      VL-Extra-Fields:
        zone: us-east1-c
        source: victoria-logs-agent
```

This feature lets you attach metadata to every log entry,
making it easier to filter, group, or analyze logs based on these additional attributes.

### HTTP listen address

The chart configures the HTTP listen address via `.Values.http` list.
Each entry in the list configures a separate HTTP listener.
The primary listener (marked with `primary: true`) is used for Pod health checks and Service port generation.

To enable TLS on the HTTP listener:

```yaml
http:
  - name: https
    primary: true
    value: :9429
    tls: true
    tlsCertFile: /path/to/tls.crt
    tlsKeyFile: /path/to/tls.key
```

### Syslog

The chart supports accepting logs over syslog protocol via `.Values.syslog.tcp` and `.Values.syslog.udp` lists.
Each entry configures a separate syslog listener on the given address.

To accept syslog messages over TCP and UDP:

```yaml
syslog:
  tcp:
    - name: syslog-tcp
      value: :514
  udp:
    - name: syslog-udp
      value: :514
```

To enable TLS on the syslog TCP listener:

```yaml
syslog:
  tcp:
    - name: syslog-tls
      value: :6514
      tls: true
      tlsCertFile: /path/to/tls.crt
      tlsKeyFile: /path/to/tls.key
```

## How to install

Access a Kubernetes cluster.

### Setup chart repository (can be omitted for OCI repositories)

Add a chart helm repository with follow commands:

```console
helm repo add vm https://victoriametrics.github.io/helm-charts/

helm repo update
```
List versions of `vm/victoria-logs-agent` chart available to installation:

```console
helm search repo vm/victoria-logs-agent -l
```

### Install `victoria-logs-agent` chart

Export default values of `victoria-logs-agent` chart to file `values.yaml`:

  - For HTTPS repository

    ```console
    helm show values vm/victoria-logs-agent > values.yaml
    ```
  - For OCI repository

    ```console
    helm show values oci://ghcr.io/victoriametrics/helm-charts/victoria-logs-agent > values.yaml
    ```

Change the values according to the need of the environment in ``values.yaml`` file.

> Consider setting `.Values.nameOverride` to a small value like `vla` to avoid hitting resource name limits of 63 characters

Test the installation with command:

  - For HTTPS repository

    ```console
    helm install vla vm/victoria-logs-agent -f values.yaml -n NAMESPACE --debug --dry-run
    ```

  - For OCI repository

    ```console
    helm install vla oci://ghcr.io/victoriametrics/helm-charts/victoria-logs-agent -f values.yaml -n NAMESPACE --debug --dry-run
    ```

Install chart with command:

  - For HTTPS repository

    ```console
    helm install vla vm/victoria-logs-agent -f values.yaml -n NAMESPACE
    ```

  - For OCI repository

    ```console
    helm install vla oci://ghcr.io/victoriametrics/helm-charts/victoria-logs-agent -f values.yaml -n NAMESPACE
    ```

Get the pods lists by running this commands:

```console
kubectl get pods -A | grep 'vla'
```

Get the application by running this command:

```console
helm list -f vla -n NAMESPACE
```

See the history of versions of `vla` application with command.

```console
helm history vla -n NAMESPACE
```

## How to uninstall

Remove application with command.

```console
helm uninstall vla -n NAMESPACE
```

## Parameters

The following tables lists the configurable parameters of the chart and their default values.

Change the values according to the need of the environment in ``victoria-logs-agent/values.yaml`` file.

<table class="helm-vars">
  <thead>
    <th class="helm-vars-key">Key</th>
    <th class="helm-vars-description">Description</th>
  </thead>
  <tbody>
    <tr id="affinity">
      <td><a href="#affinity"><pre class="chroma"><code><span class="line"><span class="cl"><span class="nt">affinity</span><span class="p">:</span><span class="w"> </span>{}</span></span></code></pre>
</a></td>
      <td><em><code>(object)</code></em><p>Pod affinity</p>
</td>
    </tr>
    <tr id="annotations">
      <td><a href="#annotations"><pre class="chroma"><code><span class="line"><span class="cl"><span class="nt">annotations</span><span class="p">:</span><span class="w"> </span>{}</span></span></code></pre>
</a></td>
      <td><em><code>(object)</code></em><p>Annotations to be added to the deployment</p>
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
      <td><em><code>(list)</code></em><p>Environment variables (ex.: secret tokens)</p>
</td>
    </tr>
    <tr id="extraargs">
      <td><a href="#extraargs"><pre class="chroma"><code><span class="line"><span class="cl"><span class="nt">extraArgs</span><span class="p">:</span><span class="w">
</span></span></span><span class="line"><span class="cl"><span class="w">    </span><span class="nt">envflag.enable</span><span class="p">:</span><span class="w"> </span><span class="kc">true</span><span class="w">
</span></span></span><span class="line"><span class="cl"><span class="w">    </span><span class="nt">envflag.prefix</span><span class="p">:</span><span class="w"> </span><span class="l">VL_</span><span class="w">
</span></span></span><span class="line"><span class="cl"><span class="w">    </span><span class="nt">loggerFormat</span><span class="p">:</span><span class="w"> </span><span class="l">json</span></span></span></code></pre>
</a></td>
      <td><em><code>(object)</code></em><p>VLAgent extra command line arguments</p>
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
      <td><em><code>(list)</code></em><p>Extra Volume Mounts for the container</p>
</td>
    </tr>
    <tr id="extravolumes">
      <td><a href="#extravolumes"><pre class="chroma"><code><span class="line"><span class="cl"><span class="nt">extraVolumes</span><span class="p">:</span><span class="w"> </span><span class="p">[]</span></span></span></code></pre>
</a></td>
      <td><em><code>(list)</code></em><p>Extra Volumes for the pod</p>
</td>
    </tr>
    <tr id="fullnameoverride">
      <td><a href="#fullnameoverride"><pre class="chroma"><code><span class="line"><span class="cl"><span class="nt">fullnameOverride</span><span class="p">:</span><span class="w"> </span><span class="s2">&#34;&#34;</span></span></span></code></pre>
</a></td>
      <td><em><code>(string)</code></em><p>Override resources fullname</p>
</td>
    </tr>
    <tr id="global-cluster-dnsdomain">
      <td><a href="#global-cluster-dnsdomain"><pre class="chroma"><code><span class="line"><span class="cl"><span class="nt">global.cluster.dnsDomain</span><span class="p">:</span><span class="w"> </span><span class="l">cluster.local.</span></span></span></code></pre>
</a></td>
      <td><em><code>(string)</code></em><p>K8s cluster domain suffix, used for building storage pods&rsquo; FQDN. Details are <a href="https://kubernetes.io/docs/tasks/administer-cluster/dns-custom-nameservers/" target="_blank">here</a></p>
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
    <tr id="horizontalpodautoscaler">
      <td><a href="#horizontalpodautoscaler"><pre class="chroma"><code><span class="line"><span class="cl"><span class="nt">horizontalPodAutoscaler</span><span class="p">:</span><span class="w">
</span></span></span><span class="line"><span class="cl"><span class="w">    </span><span class="nt">enabled</span><span class="p">:</span><span class="w"> </span><span class="kc">false</span><span class="w">
</span></span></span><span class="line"><span class="cl"><span class="w">    </span><span class="nt">maxReplicas</span><span class="p">:</span><span class="w"> </span><span class="m">10</span><span class="w">
</span></span></span><span class="line"><span class="cl"><span class="w">    </span><span class="nt">metrics</span><span class="p">:</span><span class="w"> </span><span class="p">[]</span><span class="w">
</span></span></span><span class="line"><span class="cl"><span class="w">    </span><span class="nt">minReplicas</span><span class="p">:</span><span class="w"> </span><span class="m">1</span></span></span></code></pre>
</a></td>
      <td><em><code>(object)</code></em><p>Horizontal Pod Autoscaler.</p>
</td>
    </tr>
    <tr id="horizontalpodautoscaler-enabled">
      <td><a href="#horizontalpodautoscaler-enabled"><pre class="chroma"><code><span class="line"><span class="cl"><span class="nt">horizontalPodAutoscaler.enabled</span><span class="p">:</span><span class="w"> </span><span class="kc">false</span></span></span></code></pre>
</a></td>
      <td><em><code>(bool)</code></em><p>Use HPA for vmagent</p>
</td>
    </tr>
    <tr id="horizontalpodautoscaler-maxreplicas">
      <td><a href="#horizontalpodautoscaler-maxreplicas"><pre class="chroma"><code><span class="line"><span class="cl"><span class="nt">horizontalPodAutoscaler.maxReplicas</span><span class="p">:</span><span class="w"> </span><span class="m">10</span></span></span></code></pre>
</a></td>
      <td><em><code>(int)</code></em><p>Maximum replicas for HPA to use to to scale vmagent</p>
</td>
    </tr>
    <tr id="horizontalpodautoscaler-metrics">
      <td><a href="#horizontalpodautoscaler-metrics"><pre class="chroma"><code><span class="line"><span class="cl"><span class="nt">horizontalPodAutoscaler.metrics</span><span class="p">:</span><span class="w"> </span><span class="p">[]</span></span></span></code></pre>
</a></td>
      <td><em><code>(list)</code></em><p>Metric for HPA to use to scale vmagent</p>
</td>
    </tr>
    <tr id="horizontalpodautoscaler-minreplicas">
      <td><a href="#horizontalpodautoscaler-minreplicas"><pre class="chroma"><code><span class="line"><span class="cl"><span class="nt">horizontalPodAutoscaler.minReplicas</span><span class="p">:</span><span class="w"> </span><span class="m">1</span></span></span></code></pre>
</a></td>
      <td><em><code>(int)</code></em><p>Minimum replicas for HPA to use to scale vmagent</p>
</td>
    </tr>
    <tr id="http">
      <td><a href="#http"><pre class="chroma"><code><span class="line"><span class="cl"><span class="nt">http</span><span class="p">:</span><span class="w">
</span></span></span><span class="line"><span class="cl"><span class="w">    </span>- <span class="nt">name</span><span class="p">:</span><span class="w"> </span><span class="l">http</span><span class="w">
</span></span></span><span class="line"><span class="cl"><span class="w">      </span><span class="nt">primary</span><span class="p">:</span><span class="w"> </span><span class="kc">true</span><span class="w">
</span></span></span><span class="line"><span class="cl"><span class="w">      </span><span class="nt">value</span><span class="p">:</span><span class="w"> </span><span class="p">:</span><span class="m">9429</span></span></span></code></pre>
</a></td>
      <td><em><code>(list)</code></em><p>HTTP listen addresses configuration. Each item configures an HTTP listen address with optional TLS settings. Items are used for Pod ports, command line arguments and Service port generation.</p>
</td>
    </tr>
    <tr id="image-pullpolicy">
      <td><a href="#image-pullpolicy"><pre class="chroma"><code><span class="line"><span class="cl"><span class="nt">image.pullPolicy</span><span class="p">:</span><span class="w"> </span><span class="l">IfNotPresent</span></span></span></code></pre>
</a></td>
      <td><em><code>(string)</code></em><p>Image pull policy</p>
</td>
    </tr>
    <tr id="image-registry">
      <td><a href="#image-registry"><pre class="chroma"><code><span class="line"><span class="cl"><span class="nt">image.registry</span><span class="p">:</span><span class="w"> </span><span class="s2">&#34;&#34;</span></span></span></code></pre>
</a></td>
      <td><em><code>(string)</code></em><p>Image registry</p>
</td>
    </tr>
    <tr id="image-repository">
      <td><a href="#image-repository"><pre class="chroma"><code><span class="line"><span class="cl"><span class="nt">image.repository</span><span class="p">:</span><span class="w"> </span><span class="l">victoriametrics/vlagent</span></span></span></code></pre>
</a></td>
      <td><em><code>(string)</code></em><p>Image repository</p>
</td>
    </tr>
    <tr id="image-tag">
      <td><a href="#image-tag"><pre class="chroma"><code><span class="line"><span class="cl"><span class="nt">image.tag</span><span class="p">:</span><span class="w"> </span><span class="s2">&#34;&#34;</span></span></span></code></pre>
</a></td>
      <td><em><code>(string)</code></em><p>Image tag</p>
</td>
    </tr>
    <tr id="image-variant">
      <td><a href="#image-variant"><pre class="chroma"><code><span class="line"><span class="cl"><span class="nt">image.variant</span><span class="p">:</span><span class="w"> </span><span class="s2">&#34;&#34;</span></span></span></code></pre>
</a></td>
      <td><em><code>(string)</code></em><p>Image tag suffix, which is appended to <code>Chart.AppVersion</code> if no <code>server.image.tag</code> is defined</p>
</td>
    </tr>
    <tr id="imagepullsecrets">
      <td><a href="#imagepullsecrets"><pre class="chroma"><code><span class="line"><span class="cl"><span class="nt">imagePullSecrets</span><span class="p">:</span><span class="w"> </span><span class="p">[]</span></span></span></code></pre>
</a></td>
      <td><em><code>(list)</code></em><p>Image pull secrets. Overrides <code>global.imagePullSecrets</code></p>
</td>
    </tr>
    <tr id="license">
      <td><a href="#license"><pre class="chroma"><code><span class="line"><span class="cl"><span class="nt">license</span><span class="p">:</span><span class="w">
</span></span></span><span class="line"><span class="cl"><span class="w">    </span><span class="nt">key</span><span class="p">:</span><span class="w"> </span><span class="s2">&#34;&#34;</span><span class="w">
</span></span></span><span class="line"><span class="cl"><span class="w">    </span><span class="nt">secret</span><span class="p">:</span><span class="w">
</span></span></span><span class="line"><span class="cl"><span class="w">        </span><span class="nt">key</span><span class="p">:</span><span class="w"> </span><span class="s2">&#34;&#34;</span><span class="w">
</span></span></span><span class="line"><span class="cl"><span class="w">        </span><span class="nt">name</span><span class="p">:</span><span class="w"> </span><span class="s2">&#34;&#34;</span></span></span></code></pre>
</a></td>
      <td><em><code>(object)</code></em><p>Enterprise license key configuration for VictoriaLogs enterprise. Required only for VictoriaLogs enterprise. Check docs <a href="https://docs.victoriametrics.com/victoriametrics/enterprise/#victorialogs-enterprise-features" target="_blank">here</a>, Request a trial license <a href="https://victoriametrics.com/products/enterprise/trial/" target="_blank">here</a></p>
</td>
    </tr>
    <tr id="license-key">
      <td><a href="#license-key"><pre class="chroma"><code><span class="line"><span class="cl"><span class="nt">license.key</span><span class="p">:</span><span class="w"> </span><span class="s2">&#34;&#34;</span></span></span></code></pre>
</a></td>
      <td><em><code>(string)</code></em><p>License key</p>
</td>
    </tr>
    <tr id="license-secret">
      <td><a href="#license-secret"><pre class="chroma"><code><span class="line"><span class="cl"><span class="nt">license.secret</span><span class="p">:</span><span class="w">
</span></span></span><span class="line"><span class="cl"><span class="w">    </span><span class="nt">key</span><span class="p">:</span><span class="w"> </span><span class="s2">&#34;&#34;</span><span class="w">
</span></span></span><span class="line"><span class="cl"><span class="w">    </span><span class="nt">name</span><span class="p">:</span><span class="w"> </span><span class="s2">&#34;&#34;</span></span></span></code></pre>
</a></td>
      <td><em><code>(object)</code></em><p>Use existing secret with license key</p>
</td>
    </tr>
    <tr id="license-secret-key">
      <td><a href="#license-secret-key"><pre class="chroma"><code><span class="line"><span class="cl"><span class="nt">license.secret.key</span><span class="p">:</span><span class="w"> </span><span class="s2">&#34;&#34;</span></span></span></code></pre>
</a></td>
      <td><em><code>(string)</code></em><p>Key in secret with license key</p>
</td>
    </tr>
    <tr id="license-secret-name">
      <td><a href="#license-secret-name"><pre class="chroma"><code><span class="line"><span class="cl"><span class="nt">license.secret.name</span><span class="p">:</span><span class="w"> </span><span class="s2">&#34;&#34;</span></span></span></code></pre>
</a></td>
      <td><em><code>(string)</code></em><p>Existing secret name</p>
</td>
    </tr>
    <tr id="maxdiskusageperurl">
      <td><a href="#maxdiskusageperurl"><pre class="chroma"><code><span class="line"><span class="cl"><span class="nt">maxDiskUsagePerURL</span><span class="p">:</span><span class="w"> </span><span class="l">1GiB</span></span></span></code></pre>
</a></td>
      <td><em><code>(string)</code></em><p>Maximum disk usage per remoteWrite URL. When exceeded, old logs will be deleted.</p>
</td>
    </tr>
    <tr id="nameoverride">
      <td><a href="#nameoverride"><pre class="chroma"><code><span class="line"><span class="cl"><span class="nt">nameOverride</span><span class="p">:</span><span class="w"> </span><span class="s2">&#34;&#34;</span></span></span></code></pre>
</a></td>
      <td><em><code>(string)</code></em><p>Override chart name</p>
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
      <td><em><code>(object)</code></em><p>Pod&rsquo;s node selector. Details are <a href="https://kubernetes.io/docs/concepts/scheduling-eviction/assign-pod-node/#nodeselector" target="_blank">here</a></p>
</td>
    </tr>
    <tr id="persistentvolume-accessmodes">
      <td><a href="#persistentvolume-accessmodes"><pre class="chroma"><code><span class="line"><span class="cl"><span class="nt">persistentVolume.accessModes</span><span class="p">:</span><span class="w">
</span></span></span><span class="line"><span class="cl"><span class="w">    </span>- <span class="l">ReadWriteOnce</span></span></span></code></pre>
</a></td>
      <td><em><code>(list)</code></em><p>Array of access modes. Must match those of existing PV or dynamic provisioner. Details are <a href="https://kubernetes.io/docs/concepts/storage/persistent-volumes/" target="_blank">here</a></p>
</td>
    </tr>
    <tr id="persistentvolume-annotations">
      <td><a href="#persistentvolume-annotations"><pre class="chroma"><code><span class="line"><span class="cl"><span class="nt">persistentVolume.annotations</span><span class="p">:</span><span class="w"> </span>{}</span></span></code></pre>
</a></td>
      <td><em><code>(object)</code></em><p>Persistent volume annotations</p>
</td>
    </tr>
    <tr id="persistentvolume-enabled">
      <td><a href="#persistentvolume-enabled"><pre class="chroma"><code><span class="line"><span class="cl"><span class="nt">persistentVolume.enabled</span><span class="p">:</span><span class="w"> </span><span class="kc">false</span></span></span></code></pre>
</a></td>
      <td><em><code>(bool)</code></em><p>Create/use Persistent Volume Claim for server component. Empty dir if false</p>
</td>
    </tr>
    <tr id="persistentvolume-existingclaim">
      <td><a href="#persistentvolume-existingclaim"><pre class="chroma"><code><span class="line"><span class="cl"><span class="nt">persistentVolume.existingClaim</span><span class="p">:</span><span class="w"> </span><span class="s2">&#34;&#34;</span></span></span></code></pre>
</a></td>
      <td><em><code>(string)</code></em><p>Existing Claim name. If defined, PVC must be created manually before volume will be bound</p>
</td>
    </tr>
    <tr id="persistentvolume-extralabels">
      <td><a href="#persistentvolume-extralabels"><pre class="chroma"><code><span class="line"><span class="cl"><span class="nt">persistentVolume.extraLabels</span><span class="p">:</span><span class="w"> </span>{}</span></span></code></pre>
</a></td>
      <td><em><code>(object)</code></em><p>Persistent volume additional labels</p>
</td>
    </tr>
    <tr id="persistentvolume-matchlabels">
      <td><a href="#persistentvolume-matchlabels"><pre class="chroma"><code><span class="line"><span class="cl"><span class="nt">persistentVolume.matchLabels</span><span class="p">:</span><span class="w"> </span>{}</span></span></code></pre>
</a></td>
      <td><em><code>(object)</code></em><p>Bind Persistent Volume by labels. Must match all labels of targeted PV</p>
</td>
    </tr>
    <tr id="persistentvolume-size">
      <td><a href="#persistentvolume-size"><pre class="chroma"><code><span class="line"><span class="cl"><span class="nt">persistentVolume.size</span><span class="p">:</span><span class="w"> </span><span class="l">10Gi</span></span></span></code></pre>
</a></td>
      <td><em><code>(string)</code></em><p>Size of the volume. Should be calculated based on the logs you send and retention policy you set</p>
</td>
    </tr>
    <tr id="persistentvolume-storageclassname">
      <td><a href="#persistentvolume-storageclassname"><pre class="chroma"><code><span class="line"><span class="cl"><span class="nt">persistentVolume.storageClassName</span><span class="p">:</span><span class="w"> </span><span class="s2">&#34;&#34;</span></span></span></code></pre>
</a></td>
      <td><em><code>(string)</code></em><p>StorageClass to use for persistent volume claim. Empty string means default StorageClass</p>
</td>
    </tr>
    <tr id="persistentvolume-volumeattributesclassname">
      <td><a href="#persistentvolume-volumeattributesclassname"><pre class="chroma"><code><span class="line"><span class="cl"><span class="nt">persistentVolume.volumeAttributesClassName</span><span class="p">:</span><span class="w"> </span><span class="kc">null</span></span></span></code></pre>
</a></td>
      <td><em><code>(string)</code></em><p>VolumeAttributesClass to use for persistent volume</p>
</td>
    </tr>
    <tr id="podannotations">
      <td><a href="#podannotations"><pre class="chroma"><code><span class="line"><span class="cl"><span class="nt">podAnnotations</span><span class="p">:</span><span class="w"> </span>{}</span></span></code></pre>
</a></td>
      <td><em><code>(object)</code></em><p>Annotations to be added to pod</p>
</td>
    </tr>
    <tr id="podlabels">
      <td><a href="#podlabels"><pre class="chroma"><code><span class="line"><span class="cl"><span class="nt">podLabels</span><span class="p">:</span><span class="w"> </span>{}</span></span></code></pre>
</a></td>
      <td><em><code>(object)</code></em><p>Extra labels for Pods only</p>
</td>
    </tr>
    <tr id="podmanagementpolicy">
      <td><a href="#podmanagementpolicy"><pre class="chroma"><code><span class="line"><span class="cl"><span class="nt">podManagementPolicy</span><span class="p">:</span><span class="w"> </span><span class="l">OrderedReady</span></span></span></code></pre>
</a></td>
      <td><em><code>(string)</code></em><p>Deploy order policy for StatefulSet pods</p>
</td>
    </tr>
    <tr id="podsecuritycontext">
      <td><a href="#podsecuritycontext"><pre class="chroma"><code><span class="line"><span class="cl"><span class="nt">podSecurityContext</span><span class="p">:</span><span class="w">
</span></span></span><span class="line"><span class="cl"><span class="w">    </span><span class="nt">enabled</span><span class="p">:</span><span class="w"> </span><span class="kc">true</span></span></span></code></pre>
</a></td>
      <td><em><code>(object)</code></em><p>Security context to be added to pod</p>
</td>
    </tr>
    <tr id="priorityclassname">
      <td><a href="#priorityclassname"><pre class="chroma"><code><span class="line"><span class="cl"><span class="nt">priorityClassName</span><span class="p">:</span><span class="w"> </span><span class="s2">&#34;&#34;</span></span></span></code></pre>
</a></td>
      <td><em><code>(string)</code></em><p>Priority class to be assigned to the pod(s)</p>
</td>
    </tr>
    <tr id="remotewrite">
      <td><a href="#remotewrite"><pre class="chroma"><code><span class="line"><span class="cl"><span class="nt">remoteWrite</span><span class="p">:</span><span class="w"> </span><span class="p">[]</span></span></span></code></pre>
</a></td>
      <td><em><code>(list)</code></em><p>List of log destinations. Logs will be replicated to all listed destinations.  If the url path is not specified, the logs will be sent to the /insert/native endpoint.</p>
</td>
    </tr>
    <tr id="replicacount">
      <td><a href="#replicacount"><pre class="chroma"><code><span class="line"><span class="cl"><span class="nt">replicaCount</span><span class="p">:</span><span class="w"> </span><span class="m">1</span></span></span></code></pre>
</a></td>
      <td><em><code>(int)</code></em><p>Replica count</p>
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
      <td><em><code>(object)</code></em><p>Security context to be added to pod&rsquo;s containers</p>
</td>
    </tr>
    <tr id="service-annotations">
      <td><a href="#service-annotations"><pre class="chroma"><code><span class="line"><span class="cl"><span class="nt">service.annotations</span><span class="p">:</span><span class="w"> </span>{}</span></span></code></pre>
</a></td>
      <td><em><code>(object)</code></em><p>Service annotations</p>
</td>
    </tr>
    <tr id="service-clusterip">
      <td><a href="#service-clusterip"><pre class="chroma"><code><span class="line"><span class="cl"><span class="nt">service.clusterIP</span><span class="p">:</span><span class="w"> </span><span class="kc">null</span></span></span></code></pre>
</a></td>
      <td><em><code>(string)</code></em><p>Service ClusterIP. When unset, the Service is a regular ClusterIP Service with legacy naming and headless (<code>None</code>) for a <code>ClusterIP</code> Service with operator naming</p>
</td>
    </tr>
    <tr id="service-enabled">
      <td><a href="#service-enabled"><pre class="chroma"><code><span class="line"><span class="cl"><span class="nt">service.enabled</span><span class="p">:</span><span class="w"> </span><span class="kc">false</span></span></span></code></pre>
</a></td>
      <td><em><code>(bool)</code></em><p>Enable agent service</p>
</td>
    </tr>
    <tr id="service-externalips">
      <td><a href="#service-externalips"><pre class="chroma"><code><span class="line"><span class="cl"><span class="nt">service.externalIPs</span><span class="p">:</span><span class="w"> </span><span class="p">[]</span></span></span></code></pre>
</a></td>
      <td><em><code>(list)</code></em><p>Service external IPs. Check <a href="https://kubernetes.io/docs/concepts/services-networking/service/#external-ips" target="_blank">here</a> for details</p>
</td>
    </tr>
    <tr id="service-externaltrafficpolicy">
      <td><a href="#service-externaltrafficpolicy"><pre class="chroma"><code><span class="line"><span class="cl"><span class="nt">service.externalTrafficPolicy</span><span class="p">:</span><span class="w"> </span><span class="s2">&#34;&#34;</span></span></span></code></pre>
</a></td>
      <td><em><code>(string)</code></em><p>Service external traffic policy. Check <a href="https://kubernetes.io/docs/tasks/access-application-cluster/create-external-load-balancer/#preserving-the-client-source-ip" target="_blank">here</a> for details</p>
</td>
    </tr>
    <tr id="service-extralabels">
      <td><a href="#service-extralabels"><pre class="chroma"><code><span class="line"><span class="cl"><span class="nt">service.extraLabels</span><span class="p">:</span><span class="w"> </span>{}</span></span></code></pre>
</a></td>
      <td><em><code>(object)</code></em><p>Service labels</p>
</td>
    </tr>
    <tr id="service-extraports">
      <td><a href="#service-extraports"><pre class="chroma"><code><span class="line"><span class="cl"><span class="nt">service.extraPorts</span><span class="p">:</span><span class="w"> </span><span class="p">[]</span></span></span></code></pre>
</a></td>
      <td><em><code>(list)</code></em><p>Extra service ports</p>
</td>
    </tr>
    <tr id="service-healthchecknodeport">
      <td><a href="#service-healthchecknodeport"><pre class="chroma"><code><span class="line"><span class="cl"><span class="nt">service.healthCheckNodePort</span><span class="p">:</span><span class="w"> </span><span class="s2">&#34;&#34;</span></span></span></code></pre>
</a></td>
      <td><em><code>(string)</code></em><p>Health check node port for a service. Check <a href="https://kubernetes.io/docs/tasks/access-application-cluster/create-external-load-balancer/#preserving-the-client-source-ip" target="_blank">here</a> for details</p>
</td>
    </tr>
    <tr id="service-internaltrafficpolicy">
      <td><a href="#service-internaltrafficpolicy"><pre class="chroma"><code><span class="line"><span class="cl"><span class="nt">service.internalTrafficPolicy</span><span class="p">:</span><span class="w"> </span><span class="s2">&#34;&#34;</span></span></span></code></pre>
</a></td>
      <td><em><code>(string)</code></em><p>Service internal traffic policy. Check <a href="https://kubernetes.io/docs/concepts/services-networking/service/#traffic-policies" target="_blank">here</a> for details</p>
</td>
    </tr>
    <tr id="service-ipfamilies">
      <td><a href="#service-ipfamilies"><pre class="chroma"><code><span class="line"><span class="cl"><span class="nt">service.ipFamilies</span><span class="p">:</span><span class="w"> </span><span class="p">[]</span></span></span></code></pre>
</a></td>
      <td><em><code>(list)</code></em><p>List of service IP families. Check <a href="https://kubernetes.io/docs/concepts/services-networking/dual-stack/#services" target="_blank">here</a> for details.</p>
</td>
    </tr>
    <tr id="service-ipfamilypolicy">
      <td><a href="#service-ipfamilypolicy"><pre class="chroma"><code><span class="line"><span class="cl"><span class="nt">service.ipFamilyPolicy</span><span class="p">:</span><span class="w"> </span><span class="s2">&#34;&#34;</span></span></span></code></pre>
</a></td>
      <td><em><code>(string)</code></em><p>Service IP family policy. Check <a href="https://kubernetes.io/docs/concepts/services-networking/dual-stack/#services" target="_blank">here</a> for details.</p>
</td>
    </tr>
    <tr id="service-loadbalancerip">
      <td><a href="#service-loadbalancerip"><pre class="chroma"><code><span class="line"><span class="cl"><span class="nt">service.loadBalancerIP</span><span class="p">:</span><span class="w"> </span><span class="s2">&#34;&#34;</span></span></span></code></pre>
</a></td>
      <td><em><code>(string)</code></em><p>Service load balancer IP</p>
</td>
    </tr>
    <tr id="service-loadbalancersourceranges">
      <td><a href="#service-loadbalancersourceranges"><pre class="chroma"><code><span class="line"><span class="cl"><span class="nt">service.loadBalancerSourceRanges</span><span class="p">:</span><span class="w"> </span><span class="p">[]</span></span></span></code></pre>
</a></td>
      <td><em><code>(list)</code></em><p>Load balancer source range</p>
</td>
    </tr>
    <tr id="service-portoverrides">
      <td><a href="#service-portoverrides"><pre class="chroma"><code><span class="line"><span class="cl"><span class="nt">service.portOverrides</span><span class="p">:</span><span class="w"> </span>{}</span></span></code></pre>
</a></td>
      <td><em><code>(object)</code></em><p>Per-port Service overrides, keyed by port name. Each entry may override port, protocol, targetPort, nodePort and/or appProtocol</p>
</td>
    </tr>
    <tr id="service-selectorlabels">
      <td><a href="#service-selectorlabels"><pre class="chroma"><code><span class="line"><span class="cl"><span class="nt">service.selectorLabels</span><span class="p">:</span><span class="w"> </span>{}</span></span></code></pre>
</a></td>
      <td><em><code>(object)</code></em><p>Extra selector labels common for service only</p>
</td>
    </tr>
    <tr id="service-serviceport">
      <td><a href="#service-serviceport"><pre class="chroma"><code><span class="line"><span class="cl"><span class="nt">service.servicePort</span><span class="p">:</span><span class="w"> </span><span class="s2">&#34;&#34;</span></span></span></code></pre>
</a></td>
      <td><em><code>(string)</code></em><p>Service port</p>
</td>
    </tr>
    <tr id="service-targetport">
      <td><a href="#service-targetport"><pre class="chroma"><code><span class="line"><span class="cl"><span class="nt">service.targetPort</span><span class="p">:</span><span class="w"> </span><span class="l">http</span></span></span></code></pre>
</a></td>
      <td><em><code>(string)</code></em><p>Target port</p>
</td>
    </tr>
    <tr id="service-trafficdistribution">
      <td><a href="#service-trafficdistribution"><pre class="chroma"><code><span class="line"><span class="cl"><span class="nt">service.trafficDistribution</span><span class="p">:</span><span class="w"> </span><span class="s2">&#34;&#34;</span></span></span></code></pre>
</a></td>
      <td><em><code>(string)</code></em><p>Traffic Distribution. Check <a href="https://kubernetes.io/docs/concepts/services-networking/service/#traffic-distribution" target="_blank">Traffic distribution</a></p>
</td>
    </tr>
    <tr id="service-type">
      <td><a href="#service-type"><pre class="chroma"><code><span class="line"><span class="cl"><span class="nt">service.type</span><span class="p">:</span><span class="w"> </span><span class="l">ClusterIP</span></span></span></code></pre>
</a></td>
      <td><em><code>(string)</code></em><p>Service type</p>
</td>
    </tr>
    <tr id="serviceaccount-annotations">
      <td><a href="#serviceaccount-annotations"><pre class="chroma"><code><span class="line"><span class="cl"><span class="nt">serviceAccount.annotations</span><span class="p">:</span><span class="w"> </span>{}</span></span></code></pre>
</a></td>
      <td><em><code>(object)</code></em><p>Annotations to add to the service account</p>
</td>
    </tr>
    <tr id="serviceaccount-automounttoken">
      <td><a href="#serviceaccount-automounttoken"><pre class="chroma"><code><span class="line"><span class="cl"><span class="nt">serviceAccount.automountToken</span><span class="p">:</span><span class="w"> </span><span class="kc">false</span></span></span></code></pre>
</a></td>
      <td><em><code>(bool)</code></em><p>Mount API token to pod directly. vlagent doesn&rsquo;t access Kubernetes API, so it is disabled by default</p>
</td>
    </tr>
    <tr id="serviceaccount-create">
      <td><a href="#serviceaccount-create"><pre class="chroma"><code><span class="line"><span class="cl"><span class="nt">serviceAccount.create</span><span class="p">:</span><span class="w"> </span><span class="kc">true</span></span></span></code></pre>
</a></td>
      <td><em><code>(bool)</code></em><p>Specifies whether a service account should be created</p>
</td>
    </tr>
    <tr id="serviceaccount-extralabels">
      <td><a href="#serviceaccount-extralabels"><pre class="chroma"><code><span class="line"><span class="cl"><span class="nt">serviceAccount.extraLabels</span><span class="p">:</span><span class="w"> </span>{}</span></span></code></pre>
</a></td>
      <td><em><code>(object)</code></em><p>Labels to add to the service account</p>
</td>
    </tr>
    <tr id="serviceaccount-name">
      <td><a href="#serviceaccount-name"><pre class="chroma"><code><span class="line"><span class="cl"><span class="nt">serviceAccount.name</span><span class="p">:</span><span class="w"> </span><span class="kc">null</span></span></span></code></pre>
</a></td>
      <td><em><code>(string)</code></em><p>The name of the service account to use. If not set and create is true, a name is generated using the fullname template</p>
</td>
    </tr>
    <tr id="servicemonitor-annotations">
      <td><a href="#servicemonitor-annotations"><pre class="chroma"><code><span class="line"><span class="cl"><span class="nt">serviceMonitor.annotations</span><span class="p">:</span><span class="w"> </span>{}</span></span></code></pre>
</a></td>
      <td><em><code>(object)</code></em><p>Service Monitor annotations</p>
</td>
    </tr>
    <tr id="servicemonitor-basicauth">
      <td><a href="#servicemonitor-basicauth"><pre class="chroma"><code><span class="line"><span class="cl"><span class="nt">serviceMonitor.basicAuth</span><span class="p">:</span><span class="w"> </span>{}</span></span></code></pre>
</a></td>
      <td><em><code>(object)</code></em><p>Basic auth params for Service Monitor</p>
</td>
    </tr>
    <tr id="servicemonitor-enabled">
      <td><a href="#servicemonitor-enabled"><pre class="chroma"><code><span class="line"><span class="cl"><span class="nt">serviceMonitor.enabled</span><span class="p">:</span><span class="w"> </span><span class="kc">false</span></span></span></code></pre>
</a></td>
      <td><em><code>(bool)</code></em><p>Enable deployment of Service Monitor for server component. This is Prometheus operator object</p>
</td>
    </tr>
    <tr id="servicemonitor-extralabels">
      <td><a href="#servicemonitor-extralabels"><pre class="chroma"><code><span class="line"><span class="cl"><span class="nt">serviceMonitor.extraLabels</span><span class="p">:</span><span class="w"> </span>{}</span></span></code></pre>
</a></td>
      <td><em><code>(object)</code></em><p>Service Monitor labels</p>
</td>
    </tr>
    <tr id="servicemonitor-metricrelabelings">
      <td><a href="#servicemonitor-metricrelabelings"><pre class="chroma"><code><span class="line"><span class="cl"><span class="nt">serviceMonitor.metricRelabelings</span><span class="p">:</span><span class="w"> </span><span class="p">[]</span></span></span></code></pre>
</a></td>
      <td><em><code>(list)</code></em><p>Service Monitor metricRelabelings</p>
</td>
    </tr>
    <tr id="servicemonitor-port">
      <td><a href="#servicemonitor-port"><pre class="chroma"><code><span class="line"><span class="cl"><span class="nt">serviceMonitor.port</span><span class="p">:</span><span class="w"> </span><span class="s2">&#34;&#34;</span></span></span></code></pre>
</a></td>
      <td><em><code>(string)</code></em><p>Service Monitor port. Uses primary http item name by default</p>
</td>
    </tr>
    <tr id="servicemonitor-relabelings">
      <td><a href="#servicemonitor-relabelings"><pre class="chroma"><code><span class="line"><span class="cl"><span class="nt">serviceMonitor.relabelings</span><span class="p">:</span><span class="w"> </span><span class="p">[]</span></span></span></code></pre>
</a></td>
      <td><em><code>(list)</code></em><p>Service Monitor relabelings</p>
</td>
    </tr>
    <tr id="servicemonitor-targetport">
      <td><a href="#servicemonitor-targetport"><pre class="chroma"><code><span class="line"><span class="cl"><span class="nt">serviceMonitor.targetPort</span><span class="p">:</span><span class="w"> </span><span class="s2">&#34;&#34;</span></span></span></code></pre>
</a></td>
      <td><em><code>(string)</code></em><p>Service Monitor target port. Overrides port when set</p>
</td>
    </tr>
    <tr id="syslog">
      <td><a href="#syslog"><pre class="chroma"><code><span class="line"><span class="cl"><span class="nt">syslog</span><span class="p">:</span><span class="w">
</span></span></span><span class="line"><span class="cl"><span class="w">    </span><span class="nt">tcp</span><span class="p">:</span><span class="w"> </span><span class="p">[]</span><span class="w">
</span></span></span><span class="line"><span class="cl"><span class="w">    </span><span class="nt">udp</span><span class="p">:</span><span class="w"> </span><span class="p">[]</span></span></span></code></pre>
</a></td>
      <td><em><code>(object)</code></em><p>Syslog TCP listen addresses configuration.</p>
</td>
    </tr>
    <tr id="tolerations">
      <td><a href="#tolerations"><pre class="chroma"><code><span class="line"><span class="cl"><span class="nt">tolerations</span><span class="p">:</span><span class="w"> </span><span class="p">[]</span></span></span></code></pre>
</a></td>
      <td><em><code>(list)</code></em><p>Node tolerations for server scheduling to nodes with taints. Details are <a href="https://kubernetes.io/docs/concepts/configuration/assign-pod-node/" target="_blank">here</a></p>
</td>
    </tr>
    <tr id="topologyspreadconstraints">
      <td><a href="#topologyspreadconstraints"><pre class="chroma"><code><span class="line"><span class="cl"><span class="nt">topologySpreadConstraints</span><span class="p">:</span><span class="w"> </span><span class="kc">null</span></span></span></code></pre>
</a></td>
      <td><em><code>(string)</code></em><p>Pod topologySpreadConstraints</p>
</td>
    </tr>
    <tr id="verticalpodautoscaler">
      <td><a href="#verticalpodautoscaler"><pre class="chroma"><code><span class="line"><span class="cl"><span class="nt">verticalPodAutoscaler</span><span class="p">:</span><span class="w">
</span></span></span><span class="line"><span class="cl"><span class="w">    </span><span class="nt">enabled</span><span class="p">:</span><span class="w"> </span><span class="kc">false</span></span></span></code></pre>
</a></td>
      <td><em><code>(object)</code></em><p>Vertical Pod Autoscaler. Requires VPA CRD (<code>autoscaling.k8s.io/v1</code>) to be installed in the cluster. Note that VPA should not be used together with HPA on the same resource metrics (CPU/memory).</p>
</td>
    </tr>
    <tr id="verticalpodautoscaler-enabled">
      <td><a href="#verticalpodautoscaler-enabled"><pre class="chroma"><code><span class="line"><span class="cl"><span class="nt">verticalPodAutoscaler.enabled</span><span class="p">:</span><span class="w"> </span><span class="kc">false</span></span></span></code></pre>
</a></td>
      <td><em><code>(bool)</code></em><p>Use VPA for vmagent</p>
</td>
    </tr>
  </tbody>
</table>

