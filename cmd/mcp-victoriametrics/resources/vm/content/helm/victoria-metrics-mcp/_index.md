

---
weight: 16
title: VictoriaMetrics MCP
description: "Deploy the MCP server for AI assistant integration."
menu:
  docs:
    parent: helm
    weight: 16
    identifier: helm-victoria-metrics-mcp
url: /helm/victoria-metrics-mcp/
tags:
  - metrics
  - kubernetes
  - mcp
  - observability
  - monitoring
  - victoriametrics
---

![Version](https://img.shields.io/badge/0.4.1-gray?logo=Helm&labelColor=gray&link=https%3A%2F%2Fdocs.victoriametrics.com%2Fhelm%2Fvictoria-metrics-mcp%2Fchangelog%2F%23v041)
![ArtifactHub](https://img.shields.io/badge/ArtifactHub-informational?logoColor=white&color=417598&logo=artifacthub&link=https%3A%2F%2Fartifacthub.io%2Fpackages%2Fhelm%2Fvictoriametrics%2Fvictoria-metrics-mcp)
![License](https://img.shields.io/github/license/VictoriaMetrics/helm-charts?labelColor=green&label=&link=https%3A%2F%2Fgithub.com%2FVictoriaMetrics%2Fhelm-charts%2Fblob%2Fmaster%2FLICENSE)
![Slack](https://img.shields.io/badge/Join-4A154B?logo=slack&link=https%3A%2F%2Fslack.victoriametrics.com)
![X](https://img.shields.io/twitter/follow/VictoriaMetrics?style=flat&label=Follow&color=black&logo=x&labelColor=black&link=https%3A%2F%2Fx.com%2FVictoriaMetrics)
![Reddit](https://img.shields.io/reddit/subreddit-subscribers/VictoriaMetrics?style=flat&label=Join&labelColor=red&logoColor=white&logo=reddit&link=https%3A%2F%2Fwww.reddit.com%2Fr%2FVictoriaMetrics)

A Helm chart for VictoriaMetrics MCP server

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

## Chart Details

This chart will do the following:

* Rollout VictoriaMetrics MCP server.

## How to install

Access a Kubernetes cluster.

### Setup chart repository (can be omitted for OCI repositories)

Add a chart helm repository with follow commands:

```console
helm repo add vm https://victoriametrics.github.io/helm-charts/

helm repo update
```
List versions of `vm/victoria-metrics-mcp` chart available to installation:

```console
helm search repo vm/victoria-metrics-mcp -l
```

### Install `victoria-metrics-mcp` chart

Export default values of `victoria-metrics-mcp` chart to file `values.yaml`:

  - For HTTPS repository

    ```console
    helm show values vm/victoria-metrics-mcp > values.yaml
    ```
  - For OCI repository

    ```console
    helm show values oci://ghcr.io/victoriametrics/helm-charts/victoria-metrics-mcp > values.yaml
    ```

Change the values according to the need of the environment in ``values.yaml`` file.

> Consider setting `.Values.nameOverride` to a small value like `vmm` to avoid hitting resource name limits of 63 characters

Test the installation with command:

  - For HTTPS repository

    ```console
    helm install vmm vm/victoria-metrics-mcp -f values.yaml -n NAMESPACE --debug --dry-run
    ```

  - For OCI repository

    ```console
    helm install vmm oci://ghcr.io/victoriametrics/helm-charts/victoria-metrics-mcp -f values.yaml -n NAMESPACE --debug --dry-run
    ```

Install chart with command:

  - For HTTPS repository

    ```console
    helm install vmm vm/victoria-metrics-mcp -f values.yaml -n NAMESPACE
    ```

  - For OCI repository

    ```console
    helm install vmm oci://ghcr.io/victoriametrics/helm-charts/victoria-metrics-mcp -f values.yaml -n NAMESPACE
    ```

Get the pods lists by running this commands:

```console
kubectl get pods -A | grep 'vmm'
```

Get the application by running this command:

```console
helm list -f vmm -n NAMESPACE
```

See the history of versions of `vmm` application with command.

```console
helm history vmm -n NAMESPACE
```

## How to uninstall

Remove application with command.

```console
helm uninstall vmm -n NAMESPACE
```

## Parameters

The following tables lists the configurable parameters of the chart and their default values.

Change the values according to the need of the environment in ``victoria-metrics-mcp/values.yaml`` file.

<table class="helm-vars">
  <thead>
    <th class="helm-vars-key">Key</th>
    <th class="helm-vars-description">Description</th>
  </thead>
  <tbody>
    <tr id="affinity">
      <td><a href="#affinity"><pre class="chroma"><code><span class="line"><span class="cl"><span class="nt">affinity</span><span class="p">:</span><span class="w"> </span>{}</span></span></code></pre>
</a></td>
      <td><em><code>(object)</code></em></td>
    </tr>
    <tr id="command">
      <td><a href="#command"><pre class="chroma"><code><span class="line"><span class="cl"><span class="nt">command</span><span class="p">:</span><span class="w"> </span><span class="p">[]</span></span></span></code></pre>
</a></td>
      <td><em><code>(list)</code></em><p>Override default container command. Use when the VictoriaMetrics binary is available at a custom path</p>
</td>
    </tr>
    <tr id="dnsconfig">
      <td><a href="#dnsconfig"><pre class="chroma"><code><span class="line"><span class="cl"><span class="nt">dnsConfig</span><span class="p">:</span><span class="w"> </span>{}</span></span></code></pre>
</a></td>
      <td><em><code>(object)</code></em></td>
    </tr>
    <tr id="env">
      <td><a href="#env"><pre class="chroma"><code><span class="line"><span class="cl"><span class="nt">env</span><span class="p">:</span><span class="w"> </span><span class="p">[]</span></span></span></code></pre>
</a></td>
      <td><em><code>(list)</code></em></td>
    </tr>
    <tr id="envfrom">
      <td><a href="#envfrom"><pre class="chroma"><code><span class="line"><span class="cl"><span class="nt">envFrom</span><span class="p">:</span><span class="w"> </span><span class="p">[]</span></span></span></code></pre>
</a></td>
      <td><em><code>(list)</code></em><p>Specify alternative source for env variables</p>
</td>
    </tr>
    <tr id="extracontainers">
      <td><a href="#extracontainers"><pre class="chroma"><code><span class="line"><span class="cl"><span class="nt">extraContainers</span><span class="p">:</span><span class="w"> </span><span class="p">[]</span></span></span></code></pre>
</a></td>
      <td><em><code>(list)</code></em><p>Extra containers to run in the MCP server pod</p>
</td>
    </tr>
    <tr id="extraobjects">
      <td><a href="#extraobjects"><pre class="chroma"><code><span class="line"><span class="cl"><span class="nt">extraObjects</span><span class="p">:</span><span class="w"> </span><span class="p">[]</span></span></span></code></pre>
</a></td>
      <td><em><code>(list)</code></em><p>Add extra specs dynamically to this chart</p>
</td>
    </tr>
    <tr id="fullnameoverride">
      <td><a href="#fullnameoverride"><pre class="chroma"><code><span class="line"><span class="cl"><span class="nt">fullnameOverride</span><span class="p">:</span><span class="w"> </span><span class="s2">&#34;&#34;</span></span></span></code></pre>
</a></td>
      <td><em><code>(string)</code></em></td>
    </tr>
    <tr id="horizontalpodautoscaler">
      <td><a href="#horizontalpodautoscaler"><pre class="chroma"><code><span class="line"><span class="cl"><span class="nt">horizontalPodAutoscaler</span><span class="p">:</span><span class="w">
</span></span></span><span class="line"><span class="cl"><span class="w">    </span><span class="nt">behavior</span><span class="p">:</span><span class="w"> </span>{}<span class="w">
</span></span></span><span class="line"><span class="cl"><span class="w">    </span><span class="nt">enabled</span><span class="p">:</span><span class="w"> </span><span class="kc">false</span><span class="w">
</span></span></span><span class="line"><span class="cl"><span class="w">    </span><span class="nt">labels</span><span class="p">:</span><span class="w"> </span>{}<span class="w">
</span></span></span><span class="line"><span class="cl"><span class="w">    </span><span class="nt">maxReplicas</span><span class="p">:</span><span class="w"> </span><span class="m">10</span><span class="w">
</span></span></span><span class="line"><span class="cl"><span class="w">    </span><span class="nt">metrics</span><span class="p">:</span><span class="w"> </span><span class="p">[]</span><span class="w">
</span></span></span><span class="line"><span class="cl"><span class="w">    </span><span class="nt">minReplicas</span><span class="p">:</span><span class="w"> </span><span class="m">2</span></span></span></code></pre>
</a></td>
      <td><em><code>(object)</code></em><p>Horizontal Pod Autoscaler. Details are <a href="https://kubernetes.io/docs/tasks/run-application/horizontal-pod-autoscale/" target="_blank">here</a></p>
</td>
    </tr>
    <tr id="horizontalpodautoscaler-behavior">
      <td><a href="#horizontalpodautoscaler-behavior"><pre class="chroma"><code><span class="line"><span class="cl"><span class="nt">horizontalPodAutoscaler.behavior</span><span class="p">:</span><span class="w"> </span>{}</span></span></code></pre>
</a></td>
      <td><em><code>(object)</code></em><p>Behavior settings for scaling by the HPA</p>
</td>
    </tr>
    <tr id="horizontalpodautoscaler-enabled">
      <td><a href="#horizontalpodautoscaler-enabled"><pre class="chroma"><code><span class="line"><span class="cl"><span class="nt">horizontalPodAutoscaler.enabled</span><span class="p">:</span><span class="w"> </span><span class="kc">false</span></span></span></code></pre>
</a></td>
      <td><em><code>(bool)</code></em><p>Use HPA for the MCP server deployment. Disables <code>replicaCount</code></p>
</td>
    </tr>
    <tr id="horizontalpodautoscaler-labels">
      <td><a href="#horizontalpodautoscaler-labels"><pre class="chroma"><code><span class="line"><span class="cl"><span class="nt">horizontalPodAutoscaler.labels</span><span class="p">:</span><span class="w"> </span>{}</span></span></code></pre>
</a></td>
      <td><em><code>(object)</code></em><p>Extra labels for HorizontalPodAutoscaler</p>
</td>
    </tr>
    <tr id="horizontalpodautoscaler-maxreplicas">
      <td><a href="#horizontalpodautoscaler-maxreplicas"><pre class="chroma"><code><span class="line"><span class="cl"><span class="nt">horizontalPodAutoscaler.maxReplicas</span><span class="p">:</span><span class="w"> </span><span class="m">10</span></span></span></code></pre>
</a></td>
      <td><em><code>(int)</code></em><p>Maximum replicas for HPA to use to scale the deployment</p>
</td>
    </tr>
    <tr id="horizontalpodautoscaler-metrics">
      <td><a href="#horizontalpodautoscaler-metrics"><pre class="chroma"><code><span class="line"><span class="cl"><span class="nt">horizontalPodAutoscaler.metrics</span><span class="p">:</span><span class="w"> </span><span class="p">[]</span></span></span></code></pre>
</a></td>
      <td><em><code>(list)</code></em><p>Metrics for HPA to use to scale the deployment. When empty the field is omitted and Kubernetes applies its default of 80% average CPU utilization (requires CPU requests to be set)</p>
</td>
    </tr>
    <tr id="horizontalpodautoscaler-minreplicas">
      <td><a href="#horizontalpodautoscaler-minreplicas"><pre class="chroma"><code><span class="line"><span class="cl"><span class="nt">horizontalPodAutoscaler.minReplicas</span><span class="p">:</span><span class="w"> </span><span class="m">2</span></span></span></code></pre>
</a></td>
      <td><em><code>(int)</code></em><p>Minimum replicas for HPA to use to scale the deployment</p>
</td>
    </tr>
    <tr id="image-pullpolicy">
      <td><a href="#image-pullpolicy"><pre class="chroma"><code><span class="line"><span class="cl"><span class="nt">image.pullPolicy</span><span class="p">:</span><span class="w"> </span><span class="l">IfNotPresent</span></span></span></code></pre>
</a></td>
      <td><em><code>(string)</code></em></td>
    </tr>
    <tr id="image-registry">
      <td><a href="#image-registry"><pre class="chroma"><code><span class="line"><span class="cl"><span class="nt">image.registry</span><span class="p">:</span><span class="w"> </span><span class="l">ghcr.io</span></span></span></code></pre>
</a></td>
      <td><em><code>(string)</code></em></td>
    </tr>
    <tr id="image-repository">
      <td><a href="#image-repository"><pre class="chroma"><code><span class="line"><span class="cl"><span class="nt">image.repository</span><span class="p">:</span><span class="w"> </span><span class="l">victoriametrics/mcp-victoriametrics</span></span></span></code></pre>
</a></td>
      <td><em><code>(string)</code></em></td>
    </tr>
    <tr id="image-tag">
      <td><a href="#image-tag"><pre class="chroma"><code><span class="line"><span class="cl"><span class="nt">image.tag</span><span class="p">:</span><span class="w"> </span><span class="s2">&#34;&#34;</span></span></span></code></pre>
</a></td>
      <td><em><code>(string)</code></em></td>
    </tr>
    <tr id="imagepullsecrets">
      <td><a href="#imagepullsecrets"><pre class="chroma"><code><span class="line"><span class="cl"><span class="nt">imagePullSecrets</span><span class="p">:</span><span class="w"> </span><span class="p">[]</span></span></span></code></pre>
</a></td>
      <td><em><code>(list)</code></em></td>
    </tr>
    <tr id="ingress-annotations">
      <td><a href="#ingress-annotations"><pre class="chroma"><code><span class="line"><span class="cl"><span class="nt">ingress.annotations</span><span class="p">:</span><span class="w"> </span>{}</span></span></code></pre>
</a></td>
      <td><em><code>(object)</code></em></td>
    </tr>
    <tr id="ingress-classname">
      <td><a href="#ingress-classname"><pre class="chroma"><code><span class="line"><span class="cl"><span class="nt">ingress.className</span><span class="p">:</span><span class="w"> </span><span class="s2">&#34;&#34;</span></span></span></code></pre>
</a></td>
      <td><em><code>(string)</code></em></td>
    </tr>
    <tr id="ingress-enabled">
      <td><a href="#ingress-enabled"><pre class="chroma"><code><span class="line"><span class="cl"><span class="nt">ingress.enabled</span><span class="p">:</span><span class="w"> </span><span class="kc">false</span></span></span></code></pre>
</a></td>
      <td><em><code>(bool)</code></em></td>
    </tr>
    <tr id="ingress-hosts[0]-host">
      <td><a href="#ingress-hosts[0]-host"><pre class="chroma"><code><span class="line"><span class="cl"><span class="nt">ingress.hosts[0].host</span><span class="p">:</span><span class="w"> </span><span class="l">chart-example.local</span></span></span></code></pre>
</a></td>
      <td><em><code>(string)</code></em></td>
    </tr>
    <tr id="ingress-hosts[0]-paths[0]-path">
      <td><a href="#ingress-hosts[0]-paths[0]-path"><pre class="chroma"><code><span class="line"><span class="cl"><span class="nt">ingress.hosts[0].paths[0].path</span><span class="p">:</span><span class="w"> </span><span class="l">/</span></span></span></code></pre>
</a></td>
      <td><em><code>(string)</code></em></td>
    </tr>
    <tr id="ingress-hosts[0]-paths[0]-pathtype">
      <td><a href="#ingress-hosts[0]-paths[0]-pathtype"><pre class="chroma"><code><span class="line"><span class="cl"><span class="nt">ingress.hosts[0].paths[0].pathType</span><span class="p">:</span><span class="w"> </span><span class="l">ImplementationSpecific</span></span></span></code></pre>
</a></td>
      <td><em><code>(string)</code></em></td>
    </tr>
    <tr id="ingress-tls">
      <td><a href="#ingress-tls"><pre class="chroma"><code><span class="line"><span class="cl"><span class="nt">ingress.tls</span><span class="p">:</span><span class="w"> </span><span class="p">[]</span></span></span></code></pre>
</a></td>
      <td><em><code>(list)</code></em></td>
    </tr>
    <tr id="initcontainers">
      <td><a href="#initcontainers"><pre class="chroma"><code><span class="line"><span class="cl"><span class="nt">initContainers</span><span class="p">:</span><span class="w"> </span><span class="p">[]</span></span></span></code></pre>
</a></td>
      <td><em><code>(list)</code></em><p>Init containers for the MCP server pod</p>
</td>
    </tr>
    <tr id="lifecycle">
      <td><a href="#lifecycle"><pre class="chroma"><code><span class="line"><span class="cl"><span class="nt">lifecycle</span><span class="p">:</span><span class="w"> </span>{}</span></span></code></pre>
</a></td>
      <td><em><code>(object)</code></em><p>Container lifecycle hooks</p>
</td>
    </tr>
    <tr id="livenessprobe-httpget-path">
      <td><a href="#livenessprobe-httpget-path"><pre class="chroma"><code><span class="line"><span class="cl"><span class="nt">livenessProbe.httpGet.path</span><span class="p">:</span><span class="w"> </span><span class="l">/health/liveness</span></span></span></code></pre>
</a></td>
      <td><em><code>(string)</code></em></td>
    </tr>
    <tr id="livenessprobe-httpget-port">
      <td><a href="#livenessprobe-httpget-port"><pre class="chroma"><code><span class="line"><span class="cl"><span class="nt">livenessProbe.httpGet.port</span><span class="p">:</span><span class="w"> </span><span class="l">http</span></span></span></code></pre>
</a></td>
      <td><em><code>(string)</code></em></td>
    </tr>
    <tr id="livenessprobe-httpget-scheme">
      <td><a href="#livenessprobe-httpget-scheme"><pre class="chroma"><code><span class="line"><span class="cl"><span class="nt">livenessProbe.httpGet.scheme</span><span class="p">:</span><span class="w"> </span><span class="l">HTTP</span></span></span></code></pre>
</a></td>
      <td><em><code>(string)</code></em></td>
    </tr>
    <tr id="livenessprobe-initialdelayseconds">
      <td><a href="#livenessprobe-initialdelayseconds"><pre class="chroma"><code><span class="line"><span class="cl"><span class="nt">livenessProbe.initialDelaySeconds</span><span class="p">:</span><span class="w"> </span><span class="m">10</span></span></span></code></pre>
</a></td>
      <td><em><code>(int)</code></em></td>
    </tr>
    <tr id="livenessprobe-periodseconds">
      <td><a href="#livenessprobe-periodseconds"><pre class="chroma"><code><span class="line"><span class="cl"><span class="nt">livenessProbe.periodSeconds</span><span class="p">:</span><span class="w"> </span><span class="m">10</span></span></span></code></pre>
</a></td>
      <td><em><code>(int)</code></em></td>
    </tr>
    <tr id="livenessprobe-timeoutseconds">
      <td><a href="#livenessprobe-timeoutseconds"><pre class="chroma"><code><span class="line"><span class="cl"><span class="nt">livenessProbe.timeoutSeconds</span><span class="p">:</span><span class="w"> </span><span class="m">1</span></span></span></code></pre>
</a></td>
      <td><em><code>(int)</code></em></td>
    </tr>
    <tr id="mcp-disable-resources">
      <td><a href="#mcp-disable-resources"><pre class="chroma"><code><span class="line"><span class="cl"><span class="nt">mcp.disable.resources</span><span class="p">:</span><span class="w"> </span><span class="kc">false</span></span></span></code></pre>
</a></td>
      <td><em><code>(bool)</code></em></td>
    </tr>
    <tr id="mcp-disable-tools">
      <td><a href="#mcp-disable-tools"><pre class="chroma"><code><span class="line"><span class="cl"><span class="nt">mcp.disable.tools</span><span class="p">:</span><span class="w"> </span><span class="p">[]</span></span></span></code></pre>
</a></td>
      <td><em><code>(list)</code></em></td>
    </tr>
    <tr id="mcp-heartbeatinterval">
      <td><a href="#mcp-heartbeatinterval"><pre class="chroma"><code><span class="line"><span class="cl"><span class="nt">mcp.heartbeatInterval</span><span class="p">:</span><span class="w"> </span><span class="l">30s</span></span></span></code></pre>
</a></td>
      <td><em><code>(string)</code></em></td>
    </tr>
    <tr id="mcp-mode">
      <td><a href="#mcp-mode"><pre class="chroma"><code><span class="line"><span class="cl"><span class="nt">mcp.mode</span><span class="p">:</span><span class="w"> </span><span class="l">http</span></span></span></code></pre>
</a></td>
      <td><em><code>(string)</code></em></td>
    </tr>
    <tr id="mcp-passthroughheaders">
      <td><a href="#mcp-passthroughheaders"><pre class="chroma"><code><span class="line"><span class="cl"><span class="nt">mcp.passthroughHeaders</span><span class="p">:</span><span class="w"> </span><span class="p">[]</span></span></span></code></pre>
</a></td>
      <td><em><code>(list)</code></em></td>
    </tr>
    <tr id="nameoverride">
      <td><a href="#nameoverride"><pre class="chroma"><code><span class="line"><span class="cl"><span class="nt">nameOverride</span><span class="p">:</span><span class="w"> </span><span class="s2">&#34;&#34;</span></span></span></code></pre>
</a></td>
      <td><em><code>(string)</code></em></td>
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
      <td><em><code>(object)</code></em></td>
    </tr>
    <tr id="podannotations">
      <td><a href="#podannotations"><pre class="chroma"><code><span class="line"><span class="cl"><span class="nt">podAnnotations</span><span class="p">:</span><span class="w"> </span>{}</span></span></code></pre>
</a></td>
      <td><em><code>(object)</code></em></td>
    </tr>
    <tr id="poddisruptionbudget">
      <td><a href="#poddisruptionbudget"><pre class="chroma"><code><span class="line"><span class="cl"><span class="nt">podDisruptionBudget</span><span class="p">:</span><span class="w">
</span></span></span><span class="line"><span class="cl"><span class="w">    </span><span class="nt">enabled</span><span class="p">:</span><span class="w"> </span><span class="kc">false</span><span class="w">
</span></span></span><span class="line"><span class="cl"><span class="w">    </span><span class="nt">labels</span><span class="p">:</span><span class="w"> </span>{}<span class="w">
</span></span></span><span class="line"><span class="cl"><span class="w">    </span><span class="nt">maxUnavailable</span><span class="p">:</span><span class="w"> </span><span class="kc">null</span><span class="w">
</span></span></span><span class="line"><span class="cl"><span class="w">    </span><span class="nt">minAvailable</span><span class="p">:</span><span class="w"> </span><span class="kc">null</span><span class="w">
</span></span></span><span class="line"><span class="cl"><span class="w">    </span><span class="nt">unhealthyPodEvictionPolicy</span><span class="p">:</span><span class="w"> </span><span class="kc">null</span></span></span></code></pre>
</a></td>
      <td><em><code>(object)</code></em><p>Pod disruption budget. Details are <a href="https://kubernetes.io/docs/tasks/run-application/configure-pdb/" target="_blank">here</a>. Exactly one of <code>minAvailable</code> or <code>maxUnavailable</code> must be set when enabled; <code>0</code> is a valid value</p>
</td>
    </tr>
    <tr id="poddisruptionbudget-labels">
      <td><a href="#poddisruptionbudget-labels"><pre class="chroma"><code><span class="line"><span class="cl"><span class="nt">podDisruptionBudget.labels</span><span class="p">:</span><span class="w"> </span>{}</span></span></code></pre>
</a></td>
      <td><em><code>(object)</code></em><p>Extra labels for PodDisruptionBudget</p>
</td>
    </tr>
    <tr id="poddisruptionbudget-maxunavailable">
      <td><a href="#poddisruptionbudget-maxunavailable"><pre class="chroma"><code><span class="line"><span class="cl"><span class="nt">podDisruptionBudget.maxUnavailable</span><span class="p">:</span><span class="w"> </span><span class="kc">null</span></span></span></code></pre>
</a></td>
      <td><em><code>(string)</code></em><p>max number or percentage of pods that can be unavailable</p>
</td>
    </tr>
    <tr id="poddisruptionbudget-minavailable">
      <td><a href="#poddisruptionbudget-minavailable"><pre class="chroma"><code><span class="line"><span class="cl"><span class="nt">podDisruptionBudget.minAvailable</span><span class="p">:</span><span class="w"> </span><span class="kc">null</span></span></span></code></pre>
</a></td>
      <td><em><code>(string)</code></em><p>min number or percentage of pods that must be available</p>
</td>
    </tr>
    <tr id="poddisruptionbudget-unhealthypodevictionpolicy">
      <td><a href="#poddisruptionbudget-unhealthypodevictionpolicy"><pre class="chroma"><code><span class="line"><span class="cl"><span class="nt">podDisruptionBudget.unhealthyPodEvictionPolicy</span><span class="p">:</span><span class="w"> </span><span class="kc">null</span></span></span></code></pre>
</a></td>
      <td><em><code>(string)</code></em><p>Defines criteria when unhealthy pods should be considered for eviction</p>
</td>
    </tr>
    <tr id="podlabels">
      <td><a href="#podlabels"><pre class="chroma"><code><span class="line"><span class="cl"><span class="nt">podLabels</span><span class="p">:</span><span class="w"> </span>{}</span></span></code></pre>
</a></td>
      <td><em><code>(object)</code></em></td>
    </tr>
    <tr id="podsecuritycontext">
      <td><a href="#podsecuritycontext"><pre class="chroma"><code><span class="line"><span class="cl"><span class="nt">podSecurityContext</span><span class="p">:</span><span class="w"> </span>{}</span></span></code></pre>
</a></td>
      <td><em><code>(object)</code></em></td>
    </tr>
    <tr id="priorityclassname">
      <td><a href="#priorityclassname"><pre class="chroma"><code><span class="line"><span class="cl"><span class="nt">priorityClassName</span><span class="p">:</span><span class="w"> </span><span class="s2">&#34;&#34;</span></span></span></code></pre>
</a></td>
      <td><em><code>(string)</code></em><p>Name of Priority Class</p>
</td>
    </tr>
    <tr id="readinessprobe-httpget-path">
      <td><a href="#readinessprobe-httpget-path"><pre class="chroma"><code><span class="line"><span class="cl"><span class="nt">readinessProbe.httpGet.path</span><span class="p">:</span><span class="w"> </span><span class="l">/health/readiness</span></span></span></code></pre>
</a></td>
      <td><em><code>(string)</code></em></td>
    </tr>
    <tr id="readinessprobe-httpget-port">
      <td><a href="#readinessprobe-httpget-port"><pre class="chroma"><code><span class="line"><span class="cl"><span class="nt">readinessProbe.httpGet.port</span><span class="p">:</span><span class="w"> </span><span class="l">http</span></span></span></code></pre>
</a></td>
      <td><em><code>(string)</code></em></td>
    </tr>
    <tr id="readinessprobe-httpget-scheme">
      <td><a href="#readinessprobe-httpget-scheme"><pre class="chroma"><code><span class="line"><span class="cl"><span class="nt">readinessProbe.httpGet.scheme</span><span class="p">:</span><span class="w"> </span><span class="l">HTTP</span></span></span></code></pre>
</a></td>
      <td><em><code>(string)</code></em></td>
    </tr>
    <tr id="readinessprobe-initialdelayseconds">
      <td><a href="#readinessprobe-initialdelayseconds"><pre class="chroma"><code><span class="line"><span class="cl"><span class="nt">readinessProbe.initialDelaySeconds</span><span class="p">:</span><span class="w"> </span><span class="m">10</span></span></span></code></pre>
</a></td>
      <td><em><code>(int)</code></em></td>
    </tr>
    <tr id="readinessprobe-periodseconds">
      <td><a href="#readinessprobe-periodseconds"><pre class="chroma"><code><span class="line"><span class="cl"><span class="nt">readinessProbe.periodSeconds</span><span class="p">:</span><span class="w"> </span><span class="m">10</span></span></span></code></pre>
</a></td>
      <td><em><code>(int)</code></em></td>
    </tr>
    <tr id="readinessprobe-timeoutseconds">
      <td><a href="#readinessprobe-timeoutseconds"><pre class="chroma"><code><span class="line"><span class="cl"><span class="nt">readinessProbe.timeoutSeconds</span><span class="p">:</span><span class="w"> </span><span class="m">1</span></span></span></code></pre>
</a></td>
      <td><em><code>(int)</code></em></td>
    </tr>
    <tr id="replicacount">
      <td><a href="#replicacount"><pre class="chroma"><code><span class="line"><span class="cl"><span class="nt">replicaCount</span><span class="p">:</span><span class="w"> </span><span class="m">1</span></span></span></code></pre>
</a></td>
      <td><em><code>(int)</code></em><p>Number of replicas. Ignored when <code>horizontalPodAutoscaler.enabled</code> is true</p>
</td>
    </tr>
    <tr id="resources">
      <td><a href="#resources"><pre class="chroma"><code><span class="line"><span class="cl"><span class="nt">resources</span><span class="p">:</span><span class="w"> </span>{}</span></span></code></pre>
</a></td>
      <td><em><code>(object)</code></em></td>
    </tr>
    <tr id="route">
      <td><a href="#route"><pre class="chroma"><code><span class="line"><span class="cl"><span class="nt">route</span><span class="p">:</span><span class="w">
</span></span></span><span class="line"><span class="cl"><span class="w">    </span><span class="nt">annotations</span><span class="p">:</span><span class="w"> </span>{}<span class="w">
</span></span></span><span class="line"><span class="cl"><span class="w">    </span><span class="nt">enabled</span><span class="p">:</span><span class="w"> </span><span class="kc">false</span><span class="w">
</span></span></span><span class="line"><span class="cl"><span class="w">    </span><span class="nt">hostnames</span><span class="p">:</span><span class="w"> </span><span class="p">[]</span><span class="w">
</span></span></span><span class="line"><span class="cl"><span class="w">    </span><span class="nt">parentRefs</span><span class="p">:</span><span class="w">
</span></span></span><span class="line"><span class="cl"><span class="w">        </span>- <span class="nt">name</span><span class="p">:</span><span class="w"> </span><span class="l">gateway</span><span class="w">
</span></span></span><span class="line"><span class="cl"><span class="w">          </span><span class="nt">sectionName</span><span class="p">:</span><span class="w"> </span><span class="l">http</span><span class="w">
</span></span></span><span class="line"><span class="cl"><span class="w">    </span><span class="nt">rules</span><span class="p">:</span><span class="w">
</span></span></span><span class="line"><span class="cl"><span class="w">        </span>- <span class="nt">matches</span><span class="p">:</span><span class="w">
</span></span></span><span class="line"><span class="cl"><span class="w">            </span>- <span class="nt">path</span><span class="p">:</span><span class="w">
</span></span></span><span class="line"><span class="cl"><span class="w">                </span><span class="nt">type</span><span class="p">:</span><span class="w"> </span><span class="l">PathPrefix</span><span class="w">
</span></span></span><span class="line"><span class="cl"><span class="w">                </span><span class="nt">value</span><span class="p">:</span><span class="w"> </span><span class="l">/</span></span></span></code></pre>
</a></td>
      <td><em><code>(object)</code></em><p>Expose the service via gateway-api HTTPRoute Requires Gateway API resources and suitable controller installed within the cluster (see: <a href="https://gateway-api.sigs.k8s.io/guides/" target="_blank">https://gateway-api.sigs.k8s.io/guides/</a>)</p>
</td>
    </tr>
    <tr id="scrape">
      <td><a href="#scrape"><pre class="chroma"><code><span class="line"><span class="cl"><span class="nt">scrape</span><span class="p">:</span><span class="w">
</span></span></span><span class="line"><span class="cl"><span class="w">    </span><span class="nt">annotations</span><span class="p">:</span><span class="w"> </span>{}<span class="w">
</span></span></span><span class="line"><span class="cl"><span class="w">    </span><span class="nt">basicAuth</span><span class="p">:</span><span class="w"> </span>{}<span class="w">
</span></span></span><span class="line"><span class="cl"><span class="w">    </span><span class="nt">enabled</span><span class="p">:</span><span class="w"> </span><span class="kc">false</span><span class="w">
</span></span></span><span class="line"><span class="cl"><span class="w">    </span><span class="nt">extraLabels</span><span class="p">:</span><span class="w"> </span>{}<span class="w">
</span></span></span><span class="line"><span class="cl"><span class="w">    </span><span class="nt">interval</span><span class="p">:</span><span class="w"> </span><span class="s2">&#34;&#34;</span><span class="w">
</span></span></span><span class="line"><span class="cl"><span class="w">    </span><span class="nt">metricRelabelConfigs</span><span class="p">:</span><span class="w"> </span><span class="p">[]</span><span class="w">
</span></span></span><span class="line"><span class="cl"><span class="w">    </span><span class="nt">namespace</span><span class="p">:</span><span class="w"> </span><span class="s2">&#34;&#34;</span><span class="w">
</span></span></span><span class="line"><span class="cl"><span class="w">    </span><span class="nt">path</span><span class="p">:</span><span class="w"> </span><span class="s2">&#34;&#34;</span><span class="w">
</span></span></span><span class="line"><span class="cl"><span class="w">    </span><span class="nt">relabelConfigs</span><span class="p">:</span><span class="w"> </span><span class="p">[]</span><span class="w">
</span></span></span><span class="line"><span class="cl"><span class="w">    </span><span class="nt">scheme</span><span class="p">:</span><span class="w"> </span><span class="s2">&#34;&#34;</span><span class="w">
</span></span></span><span class="line"><span class="cl"><span class="w">    </span><span class="nt">scrapeTimeout</span><span class="p">:</span><span class="w"> </span><span class="s2">&#34;&#34;</span><span class="w">
</span></span></span><span class="line"><span class="cl"><span class="w">    </span><span class="nt">tlsConfig</span><span class="p">:</span><span class="w"> </span>{}</span></span></code></pre>
</a></td>
      <td><em><code>(object)</code></em><p><code>VMServiceScrape</code> for the MCP server. Requires the VictoriaMetrics operator</p>
</td>
    </tr>
    <tr id="scrape-annotations">
      <td><a href="#scrape-annotations"><pre class="chroma"><code><span class="line"><span class="cl"><span class="nt">scrape.annotations</span><span class="p">:</span><span class="w"> </span>{}</span></span></code></pre>
</a></td>
      <td><em><code>(object)</code></em><p>Extra annotations for the <code>VMServiceScrape</code></p>
</td>
    </tr>
    <tr id="scrape-basicauth">
      <td><a href="#scrape-basicauth"><pre class="chroma"><code><span class="line"><span class="cl"><span class="nt">scrape.basicAuth</span><span class="p">:</span><span class="w"> </span>{}</span></span></code></pre>
</a></td>
      <td><em><code>(object)</code></em><p>Basic auth params for scraping (supports helm templates)</p>
</td>
    </tr>
    <tr id="scrape-extralabels">
      <td><a href="#scrape-extralabels"><pre class="chroma"><code><span class="line"><span class="cl"><span class="nt">scrape.extraLabels</span><span class="p">:</span><span class="w"> </span>{}</span></span></code></pre>
</a></td>
      <td><em><code>(object)</code></em><p>Extra labels for the <code>VMServiceScrape</code></p>
</td>
    </tr>
    <tr id="scrape-interval">
      <td><a href="#scrape-interval"><pre class="chroma"><code><span class="line"><span class="cl"><span class="nt">scrape.interval</span><span class="p">:</span><span class="w"> </span><span class="s2">&#34;&#34;</span></span></span></code></pre>
</a></td>
      <td><em><code>(string)</code></em><p>Scrape interval</p>
</td>
    </tr>
    <tr id="scrape-metricrelabelconfigs">
      <td><a href="#scrape-metricrelabelconfigs"><pre class="chroma"><code><span class="line"><span class="cl"><span class="nt">scrape.metricRelabelConfigs</span><span class="p">:</span><span class="w"> </span><span class="p">[]</span></span></span></code></pre>
</a></td>
      <td><em><code>(list)</code></em><p>Relabel configs applied to scraped metrics</p>
</td>
    </tr>
    <tr id="scrape-namespace">
      <td><a href="#scrape-namespace"><pre class="chroma"><code><span class="line"><span class="cl"><span class="nt">scrape.namespace</span><span class="p">:</span><span class="w"> </span><span class="s2">&#34;&#34;</span></span></span></code></pre>
</a></td>
      <td><em><code>(string)</code></em><p>Namespace to create the <code>VMServiceScrape</code> in. Defaults to the release namespace</p>
</td>
    </tr>
    <tr id="scrape-path">
      <td><a href="#scrape-path"><pre class="chroma"><code><span class="line"><span class="cl"><span class="nt">scrape.path</span><span class="p">:</span><span class="w"> </span><span class="s2">&#34;&#34;</span></span></span></code></pre>
</a></td>
      <td><em><code>(string)</code></em><p>Metrics path. Defaults to <code>/metrics</code></p>
</td>
    </tr>
    <tr id="scrape-relabelconfigs">
      <td><a href="#scrape-relabelconfigs"><pre class="chroma"><code><span class="line"><span class="cl"><span class="nt">scrape.relabelConfigs</span><span class="p">:</span><span class="w"> </span><span class="p">[]</span></span></span></code></pre>
</a></td>
      <td><em><code>(list)</code></em><p>Relabel configs applied to scrape targets</p>
</td>
    </tr>
    <tr id="scrape-scheme">
      <td><a href="#scrape-scheme"><pre class="chroma"><code><span class="line"><span class="cl"><span class="nt">scrape.scheme</span><span class="p">:</span><span class="w"> </span><span class="s2">&#34;&#34;</span></span></span></code></pre>
</a></td>
      <td><em><code>(string)</code></em><p>HTTP scheme to use for scraping</p>
</td>
    </tr>
    <tr id="scrape-scrapetimeout">
      <td><a href="#scrape-scrapetimeout"><pre class="chroma"><code><span class="line"><span class="cl"><span class="nt">scrape.scrapeTimeout</span><span class="p">:</span><span class="w"> </span><span class="s2">&#34;&#34;</span></span></span></code></pre>
</a></td>
      <td><em><code>(string)</code></em><p>Scrape timeout</p>
</td>
    </tr>
    <tr id="scrape-tlsconfig">
      <td><a href="#scrape-tlsconfig"><pre class="chroma"><code><span class="line"><span class="cl"><span class="nt">scrape.tlsConfig</span><span class="p">:</span><span class="w"> </span>{}</span></span></code></pre>
</a></td>
      <td><em><code>(object)</code></em><p>TLS configuration to use when scraping the endpoint</p>
</td>
    </tr>
    <tr id="securitycontext">
      <td><a href="#securitycontext"><pre class="chroma"><code><span class="line"><span class="cl"><span class="nt">securityContext</span><span class="p">:</span><span class="w"> </span>{}</span></span></code></pre>
</a></td>
      <td><em><code>(object)</code></em></td>
    </tr>
    <tr id="service-annotations">
      <td><a href="#service-annotations"><pre class="chroma"><code><span class="line"><span class="cl"><span class="nt">service.annotations</span><span class="p">:</span><span class="w"> </span>{}</span></span></code></pre>
</a></td>
      <td><em><code>(object)</code></em><p>Service annotations</p>
</td>
    </tr>
    <tr id="service-clusterip">
      <td><a href="#service-clusterip"><pre class="chroma"><code><span class="line"><span class="cl"><span class="nt">service.clusterIP</span><span class="p">:</span><span class="w"> </span><span class="s2">&#34;&#34;</span></span></span></code></pre>
</a></td>
      <td><em><code>(string)</code></em><p>Service ClusterIP</p>
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
    <tr id="service-healthchecknodeport">
      <td><a href="#service-healthchecknodeport"><pre class="chroma"><code><span class="line"><span class="cl"><span class="nt">service.healthCheckNodePort</span><span class="p">:</span><span class="w"> </span><span class="s2">&#34;&#34;</span></span></span></code></pre>
</a></td>
      <td><em><code>(string)</code></em><p>Health check node port for a service. Check <a href="https://kubernetes.io/docs/tasks/access-application-cluster/create-external-load-balancer/#preserving-the-client-source-ip" target="_blank">here</a> for details</p>
</td>
    </tr>
    <tr id="service-ipfamilies">
      <td><a href="#service-ipfamilies"><pre class="chroma"><code><span class="line"><span class="cl"><span class="nt">service.ipFamilies</span><span class="p">:</span><span class="w"> </span><span class="p">[]</span></span></span></code></pre>
</a></td>
      <td><em><code>(list)</code></em><p>List of service IP families. Check <a href="https://kubernetes.io/docs/concepts/services-networking/dual-stack/#services" target="_blank">here</a> for details</p>
</td>
    </tr>
    <tr id="service-ipfamilypolicy">
      <td><a href="#service-ipfamilypolicy"><pre class="chroma"><code><span class="line"><span class="cl"><span class="nt">service.ipFamilyPolicy</span><span class="p">:</span><span class="w"> </span><span class="s2">&#34;&#34;</span></span></span></code></pre>
</a></td>
      <td><em><code>(string)</code></em><p>Service IP family policy. Check <a href="https://kubernetes.io/docs/concepts/services-networking/dual-stack/#services" target="_blank">here</a> for details</p>
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
    <tr id="service-nodeport">
      <td><a href="#service-nodeport"><pre class="chroma"><code><span class="line"><span class="cl"><span class="nt">service.nodePort</span><span class="p">:</span><span class="w"> </span><span class="s2">&#34;&#34;</span></span></span></code></pre>
</a></td>
      <td><em><code>(string)</code></em><p>Service node port</p>
</td>
    </tr>
    <tr id="service-port">
      <td><a href="#service-port"><pre class="chroma"><code><span class="line"><span class="cl"><span class="nt">service.port</span><span class="p">:</span><span class="w"> </span><span class="m">8080</span></span></span></code></pre>
</a></td>
      <td><em><code>(int)</code></em></td>
    </tr>
    <tr id="service-trafficdistribution">
      <td><a href="#service-trafficdistribution"><pre class="chroma"><code><span class="line"><span class="cl"><span class="nt">service.trafficDistribution</span><span class="p">:</span><span class="w"> </span><span class="s2">&#34;&#34;</span></span></span></code></pre>
</a></td>
      <td><em><code>(string)</code></em><p>Service traffic distribution. Requires Kubernetes 1.31 or newer; rendering fails on older clusters. Details are <a href="https://kubernetes.io/docs/concepts/services-networking/service/#traffic-distribution" target="_blank">here</a></p>
</td>
    </tr>
    <tr id="service-type">
      <td><a href="#service-type"><pre class="chroma"><code><span class="line"><span class="cl"><span class="nt">service.type</span><span class="p">:</span><span class="w"> </span><span class="l">ClusterIP</span></span></span></code></pre>
</a></td>
      <td><em><code>(string)</code></em></td>
    </tr>
    <tr id="serviceaccount-annotations">
      <td><a href="#serviceaccount-annotations"><pre class="chroma"><code><span class="line"><span class="cl"><span class="nt">serviceAccount.annotations</span><span class="p">:</span><span class="w"> </span>{}</span></span></code></pre>
</a></td>
      <td><em><code>(object)</code></em></td>
    </tr>
    <tr id="serviceaccount-automount">
      <td><a href="#serviceaccount-automount"><pre class="chroma"><code><span class="line"><span class="cl"><span class="nt">serviceAccount.automount</span><span class="p">:</span><span class="w"> </span><span class="kc">true</span></span></span></code></pre>
</a></td>
      <td><em><code>(bool)</code></em></td>
    </tr>
    <tr id="serviceaccount-create">
      <td><a href="#serviceaccount-create"><pre class="chroma"><code><span class="line"><span class="cl"><span class="nt">serviceAccount.create</span><span class="p">:</span><span class="w"> </span><span class="kc">true</span></span></span></code></pre>
</a></td>
      <td><em><code>(bool)</code></em></td>
    </tr>
    <tr id="serviceaccount-extralabels">
      <td><a href="#serviceaccount-extralabels"><pre class="chroma"><code><span class="line"><span class="cl"><span class="nt">serviceAccount.extraLabels</span><span class="p">:</span><span class="w"> </span>{}</span></span></code></pre>
</a></td>
      <td><em><code>(object)</code></em><p>Extra labels for the service account</p>
</td>
    </tr>
    <tr id="serviceaccount-name">
      <td><a href="#serviceaccount-name"><pre class="chroma"><code><span class="line"><span class="cl"><span class="nt">serviceAccount.name</span><span class="p">:</span><span class="w"> </span><span class="s2">&#34;&#34;</span></span></span></code></pre>
</a></td>
      <td><em><code>(string)</code></em></td>
    </tr>
    <tr id="startupprobe">
      <td><a href="#startupprobe"><pre class="chroma"><code><span class="line"><span class="cl"><span class="nt">startupProbe</span><span class="p">:</span><span class="w"> </span>{}</span></span></code></pre>
</a></td>
      <td><em><code>(object)</code></em><p>Startup probe. Disabled by default</p>
</td>
    </tr>
    <tr id="strategy">
      <td><a href="#strategy"><pre class="chroma"><code><span class="line"><span class="cl"><span class="nt">strategy</span><span class="p">:</span><span class="w"> </span>{}</span></span></code></pre>
</a></td>
      <td><em><code>(object)</code></em><p>Deployment update strategy. Details are <a href="https://kubernetes.io/docs/concepts/workloads/controllers/deployment/#strategy" target="_blank">here</a></p>
</td>
    </tr>
    <tr id="terminationgraceperiodseconds">
      <td><a href="#terminationgraceperiodseconds"><pre class="chroma"><code><span class="line"><span class="cl"><span class="nt">terminationGracePeriodSeconds</span><span class="p">:</span><span class="w"> </span><span class="kc">null</span></span></span></code></pre>
</a></td>
      <td><em><code>(string)</code></em><p>Duration in seconds the pod needs to terminate gracefully. Raise it if you serve long-lived <code>sse</code> or streamable-http connections</p>
</td>
    </tr>
    <tr id="tolerations">
      <td><a href="#tolerations"><pre class="chroma"><code><span class="line"><span class="cl"><span class="nt">tolerations</span><span class="p">:</span><span class="w"> </span><span class="p">[]</span></span></span></code></pre>
</a></td>
      <td><em><code>(list)</code></em></td>
    </tr>
    <tr id="topologyspreadconstraints">
      <td><a href="#topologyspreadconstraints"><pre class="chroma"><code><span class="line"><span class="cl"><span class="nt">topologySpreadConstraints</span><span class="p">:</span><span class="w"> </span><span class="p">[]</span></span></span></code></pre>
</a></td>
      <td><em><code>(list)</code></em><p>Pod topologySpreadConstraints. Details are <a href="https://kubernetes.io/docs/concepts/scheduling-eviction/topology-spread-constraints/" target="_blank">here</a>. <code>labelSelector</code> defaults to the chart&rsquo;s selector labels when not set</p>
</td>
    </tr>
    <tr id="vm-bearertoken">
      <td><a href="#vm-bearertoken"><pre class="chroma"><code><span class="line"><span class="cl"><span class="nt">vm.bearerToken</span><span class="p">:</span><span class="w"> </span><span class="s2">&#34;&#34;</span></span></span></code></pre>
</a></td>
      <td><em><code>(string)</code></em></td>
    </tr>
    <tr id="vm-cloudapikey">
      <td><a href="#vm-cloudapikey"><pre class="chroma"><code><span class="line"><span class="cl"><span class="nt">vm.cloudAPIKey</span><span class="p">:</span><span class="w"> </span><span class="s2">&#34;&#34;</span></span></span></code></pre>
</a></td>
      <td><em><code>(string)</code></em></td>
    </tr>
    <tr id="vm-entrypoint">
      <td><a href="#vm-entrypoint"><pre class="chroma"><code><span class="line"><span class="cl"><span class="nt">vm.entrypoint</span><span class="p">:</span><span class="w"> </span><span class="s2">&#34;&#34;</span></span></span></code></pre>
</a></td>
      <td><em><code>(string)</code></em></td>
    </tr>
    <tr id="vm-headers">
      <td><a href="#vm-headers"><pre class="chroma"><code><span class="line"><span class="cl"><span class="nt">vm.headers</span><span class="p">:</span><span class="w"> </span><span class="p">[]</span></span></span></code></pre>
</a></td>
      <td><em><code>(list)</code></em></td>
    </tr>
    <tr id="vm-type">
      <td><a href="#vm-type"><pre class="chroma"><code><span class="line"><span class="cl"><span class="nt">vm.type</span><span class="p">:</span><span class="w"> </span><span class="l">single</span></span></span></code></pre>
</a></td>
      <td><em><code>(string)</code></em></td>
    </tr>
    <tr id="volumemounts">
      <td><a href="#volumemounts"><pre class="chroma"><code><span class="line"><span class="cl"><span class="nt">volumeMounts</span><span class="p">:</span><span class="w"> </span><span class="p">[]</span></span></span></code></pre>
</a></td>
      <td><em><code>(list)</code></em></td>
    </tr>
    <tr id="volumes">
      <td><a href="#volumes"><pre class="chroma"><code><span class="line"><span class="cl"><span class="nt">volumes</span><span class="p">:</span><span class="w"> </span><span class="p">[]</span></span></span></code></pre>
</a></td>
      <td><em><code>(list)</code></em></td>
    </tr>
  </tbody>
</table>

