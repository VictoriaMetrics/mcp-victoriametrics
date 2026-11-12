

---
weight: 21
title: VictoriaTraces Agent
description: "Deploy vtagent on Kubernetes."
menu:
  docs:
    parent: helm
    weight: 21
    identifier: helm-victoria-traces-agent
url: /helm/victoria-traces-agent/
tags:
  - traces
  - kubernetes
---

![Version](https://img.shields.io/badge/0.2.2-gray?logo=Helm&labelColor=gray&link=https%3A%2F%2Fdocs.victoriametrics.com%2Fhelm%2Fvictoria-traces-agent%2Fchangelog%2F%23v022)
![ArtifactHub](https://img.shields.io/badge/ArtifactHub-informational?logoColor=white&color=417598&logo=artifacthub&link=https%3A%2F%2Fartifacthub.io%2Fpackages%2Fhelm%2Fvictoriametrics%2Fvictoria-traces-agent)
![License](https://img.shields.io/github/license/VictoriaMetrics/helm-charts?labelColor=green&label=&link=https%3A%2F%2Fgithub.com%2FVictoriaMetrics%2Fhelm-charts%2Fblob%2Fmaster%2FLICENSE)
![Slack](https://img.shields.io/badge/Join-4A154B?logo=slack&link=https%3A%2F%2Fslack.victoriametrics.com)
![X](https://img.shields.io/twitter/follow/VictoriaMetrics?style=flat&label=Follow&color=black&logo=x&labelColor=black&link=https%3A%2F%2Fx.com%2FVictoriaMetrics)
![Reddit](https://img.shields.io/reddit/subreddit-subscribers/VictoriaMetrics?style=flat&label=Join&labelColor=red&logoColor=white&logo=reddit&link=https%3A%2F%2Fwww.reddit.com%2Fr%2FVictoriaMetrics)

VictoriaTraces Agent - accepts trace spans via OpenTelemetry protocol and replicates them across multiple VictoriaTraces instances.

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

This Helm chart deploys [vtagent](https://docs.victoriametrics.com/victoriatraces/vtagent/) as a StatefulSet.
It receives trace spans via [OpenTelemetry protocol](https://docs.victoriametrics.com/victoriatraces/data-ingestion/opentelemetry/)
from cluster workloads and forwards them to the configured VictoriaTraces destinations.
If more than one destination is specified, collected trace spans are replicated to all configured destinations.
When a destination is unavailable, trace spans are buffered on disk and sent as soon as the destination is back.

- To install single-node version of VictoriaTraces, see [these docs](https://docs.victoriametrics.com/helm/victoria-traces-single/).
- To install a VictoriaTraces cluster, see [these docs](https://docs.victoriametrics.com/helm/victoria-traces-cluster/).

## Chart configuration

The simplest working configuration includes specifying the `remoteWrite` array and setting CPU and memory resources for the chart.

Example of a minimal working configuration:

```yaml
remoteWrite:
  - url: http://victoria-traces:10428

resources:
  limits:
    cpu: 100m
    memory: 128Mi
  requests:
    cpu: 100m
    memory: 128Mi
```

If the `url` has no path, trace spans are sent to the `/insert/native` endpoint.
If multiple `remoteWrite` entries are defined, trace spans are replicated to all the specified destinations.

### Accepting trace spans

`vtagent` accepts trace spans over OTLP/HTTP on the primary HTTP listener (port `10429` by default)
at the `/insert/opentelemetry/v1/traces` path. The chart always creates a Service for it, named `vtagent-<release>`.
The Service is headless by default, so its DNS name resolves to the vtagent pod IPs.
To get a regular load-balanced ClusterIP Service instead, unset `service.clusterIP`:

```yaml
service:
  clusterIP: ""
```

To additionally accept trace spans over OTLP/gRPC, set the listen address:

```yaml
otlpGRPC:
  listenAddr: :4317
```

The listener accepts plain-text connections unless `otlpGRPC.tls` is set. To enable TLS on the OTLP/gRPC listener:

```yaml
otlpGRPC:
  listenAddr: :4317
  tls: true
  tlsCertFile: /etc/tls/tls.crt
  tlsKeyFile: /etc/tls/tls.key

extraVolumes:
  - name: tls-certs
    secret:
      secretName: tls-secret

extraVolumeMounts:
  - name: tls-certs
    mountPath: /etc/tls
    readOnly: true
```

### On-disk buffering

Trace spans that can't be delivered are buffered at the data volume, which is an `emptyDir` by default.
The buffer of each `remoteWrite` destination is limited by `maxDiskUsagePerURL`.
To keep the buffer across pod restarts, enable the persistent volume and size it for all the destinations:

```yaml
maxDiskUsagePerURL: 5GiB

persistentVolume:
  enabled: true
  size: 10Gi
```

### Basic auth

If you need to use basic auth, define the needed flags via environment variables as shown below:

```yaml
remoteWrite:
  - url: http://victoria-traces:10428

env:
  - name: VM_remoteWrite_basicAuth_password
    valueFrom:
      secretKeyRef:
        name: auth-secret
        key: VT_PASSWORD
  - name: VM_remoteWrite_basicAuth_username
    valueFrom:
      secretKeyRef:
        name: auth-secret
        key: VT_USERNAME
```

or mount secrets into pod and reference them in remoteWrite section:

```yaml
remoteWrite:
  - url: http://victoria-traces:10428
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

To define [tenant](https://docs.victoriametrics.com/victoriatraces/#multitenancy), use `ProjectID` and `AccountID` headers as shown below:

```yaml
remoteWrite:
  - url: http://victoria-traces:10428
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
  - url: https://victoria-traces:10428
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
  - url: https://victoria-traces:10428
    tlsCAFile: "/etc/tls/ca.crt"
    tlsCertFile: "/etc/tls/client.crt"
    tlsKeyFile: "/etc/tls/client.key"
```

If you want to disable TLS certificate verification (not recommended in production), you can set `tlsInsecureSkipVerify` to true.

```yaml
remoteWrite:
  - url: https://victoria-traces:10428
    tlsInsecureSkipVerify: true
```

### Extra fields

You can add custom fields to all the trace spans by including the `VT-Extra-Fields` header in `headers` section of remote write configuration.
For example:

```yaml
remoteWrite:
  - url: http://victoria-traces:10428
    headers:
      VT-Extra-Fields:
        zone: us-east1-c
        source: victoria-traces-agent
```

### HTTP listen address

The chart configures the HTTP listen address via `.Values.http` list.
Each entry in the list configures a separate HTTP listener.
The primary listener (marked with `primary: true`) is the one `service.servicePort`, `service.targetPort` and the Service Monitor apply to.

To enable TLS on the HTTP listener:

```yaml
http:
  - name: https
    primary: true
    value: :10429
    tls: true
    tlsCertFile: /path/to/tls.crt
    tlsKeyFile: /path/to/tls.key
```

### Resource naming

Resources are named `vtagent-<release>`, which matches the naming used by the VictoriaMetrics operator.
Legacy release-name based naming is not supported by this chart.

## How to install

Access a Kubernetes cluster.

### Setup chart repository (can be omitted for OCI repositories)

Add a chart helm repository with follow commands:

```console
helm repo add vm https://victoriametrics.github.io/helm-charts/

helm repo update
```
List versions of `vm/victoria-traces-agent` chart available to installation:

```console
helm search repo vm/victoria-traces-agent -l
```

### Install `victoria-traces-agent` chart

Export default values of `victoria-traces-agent` chart to file `values.yaml`:

  - For HTTPS repository

    ```console
    helm show values vm/victoria-traces-agent > values.yaml
    ```
  - For OCI repository

    ```console
    helm show values oci://ghcr.io/victoriametrics/helm-charts/victoria-traces-agent > values.yaml
    ```

Change the values according to the need of the environment in ``values.yaml`` file.

> Consider setting `.Values.nameOverride` to a small value like `vta` to avoid hitting resource name limits of 63 characters

Test the installation with command:

  - For HTTPS repository

    ```console
    helm install vta vm/victoria-traces-agent -f values.yaml -n NAMESPACE --debug --dry-run
    ```

  - For OCI repository

    ```console
    helm install vta oci://ghcr.io/victoriametrics/helm-charts/victoria-traces-agent -f values.yaml -n NAMESPACE --debug --dry-run
    ```

Install chart with command:

  - For HTTPS repository

    ```console
    helm install vta vm/victoria-traces-agent -f values.yaml -n NAMESPACE
    ```

  - For OCI repository

    ```console
    helm install vta oci://ghcr.io/victoriametrics/helm-charts/victoria-traces-agent -f values.yaml -n NAMESPACE
    ```

Get the pods lists by running this commands:

```console
kubectl get pods -A | grep 'vta'
```

Get the application by running this command:

```console
helm list -f vta -n NAMESPACE
```

See the history of versions of `vta` application with command.

```console
helm history vta -n NAMESPACE
```

## How to uninstall

Remove application with command.

```console
helm uninstall vta -n NAMESPACE
```

## Parameters

The following tables lists the configurable parameters of the chart and their default values.

Change the values according to the need of the environment in ``victoria-traces-agent/values.yaml`` file.

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
</span></span></span><span class="line"><span class="cl"><span class="w">    </span><span class="nt">envflag.prefix</span><span class="p">:</span><span class="w"> </span><span class="l">VM_</span><span class="w">
</span></span></span><span class="line"><span class="cl"><span class="w">    </span><span class="nt">loggerFormat</span><span class="p">:</span><span class="w"> </span><span class="l">json</span></span></span></code></pre>
</a></td>
      <td><em><code>(object)</code></em><p>VTAgent extra command line arguments</p>
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
      <td><em><code>(bool)</code></em><p>Use HPA for vtagent</p>
</td>
    </tr>
    <tr id="horizontalpodautoscaler-maxreplicas">
      <td><a href="#horizontalpodautoscaler-maxreplicas"><pre class="chroma"><code><span class="line"><span class="cl"><span class="nt">horizontalPodAutoscaler.maxReplicas</span><span class="p">:</span><span class="w"> </span><span class="m">10</span></span></span></code></pre>
</a></td>
      <td><em><code>(int)</code></em><p>Maximum replicas for HPA to use to scale vtagent</p>
</td>
    </tr>
    <tr id="horizontalpodautoscaler-metrics">
      <td><a href="#horizontalpodautoscaler-metrics"><pre class="chroma"><code><span class="line"><span class="cl"><span class="nt">horizontalPodAutoscaler.metrics</span><span class="p">:</span><span class="w"> </span><span class="p">[]</span></span></span></code></pre>
</a></td>
      <td><em><code>(list)</code></em><p>Metric for HPA to use to scale vtagent</p>
</td>
    </tr>
    <tr id="horizontalpodautoscaler-minreplicas">
      <td><a href="#horizontalpodautoscaler-minreplicas"><pre class="chroma"><code><span class="line"><span class="cl"><span class="nt">horizontalPodAutoscaler.minReplicas</span><span class="p">:</span><span class="w"> </span><span class="m">1</span></span></span></code></pre>
</a></td>
      <td><em><code>(int)</code></em><p>Minimum replicas for HPA to use to scale vtagent</p>
</td>
    </tr>
    <tr id="http">
      <td><a href="#http"><pre class="chroma"><code><span class="line"><span class="cl"><span class="nt">http</span><span class="p">:</span><span class="w">
</span></span></span><span class="line"><span class="cl"><span class="w">    </span>- <span class="nt">name</span><span class="p">:</span><span class="w"> </span><span class="l">http</span><span class="w">
</span></span></span><span class="line"><span class="cl"><span class="w">      </span><span class="nt">primary</span><span class="p">:</span><span class="w"> </span><span class="kc">true</span><span class="w">
</span></span></span><span class="line"><span class="cl"><span class="w">      </span><span class="nt">value</span><span class="p">:</span><span class="w"> </span><span class="p">:</span><span class="m">10429</span></span></span></code></pre>
</a></td>
      <td><em><code>(list)</code></em><p>HTTP listen addresses configuration. Each item configures an HTTP listen address with optional TLS settings. Items are used for Pod ports, command line arguments and Service port generation. OTLP/HTTP trace spans are accepted on the primary listener at <code>/insert/opentelemetry/v1/traces</code>.</p>
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
      <td><a href="#image-repository"><pre class="chroma"><code><span class="line"><span class="cl"><span class="nt">image.repository</span><span class="p">:</span><span class="w"> </span><span class="l">victoriametrics/vtagent</span></span></span></code></pre>
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
      <td><em><code>(string)</code></em><p>Image tag suffix, which is appended to <code>Chart.AppVersion</code> if no <code>image.tag</code> is defined</p>
</td>
    </tr>
    <tr id="imagepullsecrets">
      <td><a href="#imagepullsecrets"><pre class="chroma"><code><span class="line"><span class="cl"><span class="nt">imagePullSecrets</span><span class="p">:</span><span class="w"> </span><span class="p">[]</span></span></span></code></pre>
</a></td>
      <td><em><code>(list)</code></em><p>Image pull secrets. Overrides <code>global.imagePullSecrets</code></p>
</td>
    </tr>
    <tr id="maxdiskusageperurl">
      <td><a href="#maxdiskusageperurl"><pre class="chroma"><code><span class="line"><span class="cl"><span class="nt">maxDiskUsagePerURL</span><span class="p">:</span><span class="w"> </span><span class="l">1GiB</span></span></span></code></pre>
</a></td>
      <td><em><code>(string)</code></em><p>Maximum disk usage per remoteWrite URL. When exceeded, old trace spans will be deleted.</p>
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
</span></span></span><span class="line"><span class="cl"><span class="w">    </span><span class="nt">labels</span><span class="p">:</span><span class="w"> </span>{}<span class="w">
</span></span></span><span class="line"><span class="cl"><span class="w">    </span><span class="nt">policyTypes</span><span class="p">:</span><span class="w"> </span><span class="p">[]</span></span></span></code></pre>
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
    <tr id="networkpolicy-policytypes">
      <td><a href="#networkpolicy-policytypes"><pre class="chroma"><code><span class="line"><span class="cl"><span class="nt">networkPolicy.policyTypes</span><span class="p">:</span><span class="w"> </span><span class="p">[]</span></span></span></code></pre>
</a></td>
      <td><em><code>(list)</code></em><p>Policy types. When empty, they are derived from the configured <code>ingress</code> and <code>egress</code> rules. With no rules and no policy types Kubernetes treats the policy as <code>Ingress</code>, which denies all incoming traffic</p>
</td>
    </tr>
    <tr id="nodeselector">
      <td><a href="#nodeselector"><pre class="chroma"><code><span class="line"><span class="cl"><span class="nt">nodeSelector</span><span class="p">:</span><span class="w"> </span>{}</span></span></code></pre>
</a></td>
      <td><em><code>(object)</code></em><p>Pod&rsquo;s node selector. Details are <a href="https://kubernetes.io/docs/concepts/scheduling-eviction/assign-pod-node/#nodeselector" target="_blank">here</a></p>
</td>
    </tr>
    <tr id="otlpgrpc-listenaddr">
      <td><a href="#otlpgrpc-listenaddr"><pre class="chroma"><code><span class="line"><span class="cl"><span class="nt">otlpGRPC.listenAddr</span><span class="p">:</span><span class="w"> </span><span class="s2">&#34;&#34;</span></span></span></code></pre>
</a></td>
      <td><em><code>(string)</code></em><p>OTLP/gRPC listen address. The listener is disabled when empty. The recommended value is <code>:4317</code></p>
</td>
    </tr>
    <tr id="otlpgrpc-tls">
      <td><a href="#otlpgrpc-tls"><pre class="chroma"><code><span class="line"><span class="cl"><span class="nt">otlpGRPC.tls</span><span class="p">:</span><span class="w"> </span><span class="kc">false</span></span></span></code></pre>
</a></td>
      <td><em><code>(bool)</code></em><p>Enable TLS for incoming OTLP/gRPC requests. Requires <code>tlsCertFile</code> and <code>tlsKeyFile</code>. vtagent enables it by default, so the chart always passes this value explicitly</p>
</td>
    </tr>
    <tr id="otlpgrpc-tlscertfile">
      <td><a href="#otlpgrpc-tlscertfile"><pre class="chroma"><code><span class="line"><span class="cl"><span class="nt">otlpGRPC.tlsCertFile</span><span class="p">:</span><span class="w"> </span><span class="s2">&#34;&#34;</span></span></span></code></pre>
</a></td>
      <td><em><code>(string)</code></em><p>Path to file with TLS certificate for the OTLP/gRPC listener</p>
</td>
    </tr>
    <tr id="otlpgrpc-tlskeyfile">
      <td><a href="#otlpgrpc-tlskeyfile"><pre class="chroma"><code><span class="line"><span class="cl"><span class="nt">otlpGRPC.tlsKeyFile</span><span class="p">:</span><span class="w"> </span><span class="s2">&#34;&#34;</span></span></span></code></pre>
</a></td>
      <td><em><code>(string)</code></em><p>Path to file with TLS key for the OTLP/gRPC listener</p>
</td>
    </tr>
    <tr id="otlpgrpc-tlsminversion">
      <td><a href="#otlpgrpc-tlsminversion"><pre class="chroma"><code><span class="line"><span class="cl"><span class="nt">otlpGRPC.tlsMinVersion</span><span class="p">:</span><span class="w"> </span><span class="s2">&#34;&#34;</span></span></span></code></pre>
</a></td>
      <td><em><code>(string)</code></em><p>Minimum TLS version for the OTLP/gRPC listener</p>
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
      <td><em><code>(bool)</code></em><p>Create/use Persistent Volume Claim for vtagent. Empty dir if false</p>
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
      <td><em><code>(string)</code></em><p>Size of the volume. Should be big enough to keep buffers for all remoteWrite destinations, see <code>maxDiskUsagePerURL</code></p>
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
      <td><em><code>(list)</code></em><p>List of trace destinations. Trace spans will be replicated to all listed destinations.  If the url path is not specified, the trace spans will be sent to the /insert/native endpoint.</p>
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
      <td><a href="#service-clusterip"><pre class="chroma"><code><span class="line"><span class="cl"><span class="nt">service.clusterIP</span><span class="p">:</span><span class="w"> </span><span class="l">None</span></span></span></code></pre>
</a></td>
      <td><em><code>(string)</code></em><p>Service ClusterIP. The Service is headless by default, which gives each vtagent pod a stable DNS name. Set it to an empty string to get a regular ClusterIP Service. It is ignored when <code>service.type</code> is not <code>ClusterIP</code></p>
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
      <td><a href="#service-targetport"><pre class="chroma"><code><span class="line"><span class="cl"><span class="nt">service.targetPort</span><span class="p">:</span><span class="w"> </span><span class="s2">&#34;&#34;</span></span></span></code></pre>
</a></td>
      <td><em><code>(string)</code></em><p>Target port of the primary HTTP listener. Defaults to the primary <code>http</code> item name</p>
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
      <td><em><code>(bool)</code></em><p>Mount API token to pod directly. vtagent doesn&rsquo;t access Kubernetes API, so it is disabled by default</p>
</td>
    </tr>
    <tr id="serviceaccount-create">
      <td><a href="#serviceaccount-create"><pre class="chroma"><code><span class="line"><span class="cl"><span class="nt">serviceAccount.create</span><span class="p">:</span><span class="w"> </span><span class="kc">true</span></span></span></code></pre>
</a></td>
      <td><em><code>(bool)</code></em><p>Specifies whether a service account should be created</p>
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
      <td><em><code>(bool)</code></em><p>Enable deployment of Service Monitor for vtagent. This is Prometheus operator object</p>
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
      <td><em><code>(bool)</code></em><p>Use VPA for vtagent</p>
</td>
    </tr>
  </tbody>
</table>

