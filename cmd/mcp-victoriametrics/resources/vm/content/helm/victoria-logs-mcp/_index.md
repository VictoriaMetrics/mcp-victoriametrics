

---
weight: 17
title: VictoriaLogs MCP
menu:
  docs:
    parent: helm
    weight: 17
    identifier: helm-victoria-logs-mcp
url: /helm/victoria-logs-mcp/
tags:
  - logs
  - kubernetes
  - mcp
  - observability
  - victorialogs
---

![Version](https://img.shields.io/badge/0.1.0-gray?logo=Helm&labelColor=gray&link=https%3A%2F%2Fdocs.victoriametrics.com%2Fhelm%2Fvictoria-logs-mcp%2Fchangelog%2F%23010)
![ArtifactHub](https://img.shields.io/badge/ArtifactHub-informational?logoColor=white&color=417598&logo=artifacthub&link=https%3A%2F%2Fartifacthub.io%2Fpackages%2Fhelm%2Fvictoriametrics%2Fvictoria-logs-mcp)
![License](https://img.shields.io/github/license/VictoriaMetrics/helm-charts?labelColor=green&label=&link=https%3A%2F%2Fgithub.com%2FVictoriaMetrics%2Fhelm-charts%2Fblob%2Fmaster%2FLICENSE)
![Slack](https://img.shields.io/badge/Join-4A154B?logo=slack&link=https%3A%2F%2Fslack.victoriametrics.com)
![X](https://img.shields.io/twitter/follow/VictoriaMetrics?style=flat&label=Follow&color=black&logo=x&labelColor=black&link=https%3A%2F%2Fx.com%2FVictoriaMetrics)
![Reddit](https://img.shields.io/reddit/subreddit-subscribers/VictoriaMetrics?style=flat&label=Join&labelColor=red&logoColor=white&logo=reddit&link=https%3A%2F%2Fwww.reddit.com%2Fr%2FVictoriaMetrics)

A Helm chart for VictoriaLogs MCP server

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

* Rollout VictoriaLogs MCP server.

## How to install

Access a Kubernetes cluster.

### Setup chart repository (can be omitted for OCI repositories)

Add a chart helm repository with follow commands:

```console
helm repo add vm https://victoriametrics.github.io/helm-charts/

helm repo update
```
List versions of `vm/victoria-logs-mcp` chart available to installation:

```console
helm search repo vm/victoria-logs-mcp -l
```

### Install `victoria-logs-mcp` chart

Export default values of `victoria-logs-mcp` chart to file `values.yaml`:

  - For HTTPS repository

    ```console
    helm show values vm/victoria-logs-mcp > values.yaml
    ```
  - For OCI repository

    ```console
    helm show values oci://ghcr.io/victoriametrics/helm-charts/victoria-logs-mcp > values.yaml
    ```

Change the values according to the need of the environment in ``values.yaml`` file.

> Consider setting `.Values.nameOverride` to a small value like `vlm` to avoid hitting resource name limits of 63 characters

Test the installation with command:

  - For HTTPS repository

    ```console
    helm install vlm vm/victoria-logs-mcp -f values.yaml -n NAMESPACE --debug
    ```

  - For OCI repository

    ```console
    helm install vlm oci://ghcr.io/victoriametrics/helm-charts/victoria-logs-mcp -f values.yaml -n NAMESPACE --debug
    ```

Install chart with command:

  - For HTTPS repository

    ```console
    helm install vlm vm/victoria-logs-mcp -f values.yaml -n NAMESPACE
    ```

  - For OCI repository

    ```console
    helm install vlm oci://ghcr.io/victoriametrics/helm-charts/victoria-logs-mcp -f values.yaml -n NAMESPACE
    ```

Get the pods lists by running this commands:

```console
kubectl get pods -A | grep 'vlm'
```

Get the application by running this command:

```console
helm list -f vlm -n NAMESPACE
```

See the history of versions of `vlm` application with command.

```console
helm history vlm -n NAMESPACE
```

## How to uninstall

Remove application with command.

```console
helm uninstall vlm -n NAMESPACE
```

## Parameters

The following tables lists the configurable parameters of the chart and their default values.

Change the values according to the need of the environment in ``victoria-logs-mcp/values.yaml`` file.

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
    <tr id="env">
      <td><a href="#env"><pre class="chroma"><code><span class="line"><span class="cl"><span class="nt">env</span><span class="p">:</span><span class="w"> </span><span class="p">[]</span></span></span></code></pre>
</a></td>
      <td><em><code>(list)</code></em></td>
    </tr>
    <tr id="fullnameoverride">
      <td><a href="#fullnameoverride"><pre class="chroma"><code><span class="line"><span class="cl"><span class="nt">fullnameOverride</span><span class="p">:</span><span class="w"> </span><span class="s2">&#34;&#34;</span></span></span></code></pre>
</a></td>
      <td><em><code>(string)</code></em></td>
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
      <td><a href="#image-repository"><pre class="chroma"><code><span class="line"><span class="cl"><span class="nt">image.repository</span><span class="p">:</span><span class="w"> </span><span class="l">victoriametrics/mcp-victorialogs</span></span></span></code></pre>
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
      <td><em><code>(int)</code></em></td>
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
    <tr id="scrape-enabled">
      <td><a href="#scrape-enabled"><pre class="chroma"><code><span class="line"><span class="cl"><span class="nt">scrape.enabled</span><span class="p">:</span><span class="w"> </span><span class="kc">false</span></span></span></code></pre>
</a></td>
      <td><em><code>(bool)</code></em></td>
    </tr>
    <tr id="securitycontext">
      <td><a href="#securitycontext"><pre class="chroma"><code><span class="line"><span class="cl"><span class="nt">securityContext</span><span class="p">:</span><span class="w"> </span>{}</span></span></code></pre>
</a></td>
      <td><em><code>(object)</code></em></td>
    </tr>
    <tr id="service-port">
      <td><a href="#service-port"><pre class="chroma"><code><span class="line"><span class="cl"><span class="nt">service.port</span><span class="p">:</span><span class="w"> </span><span class="m">8080</span></span></span></code></pre>
</a></td>
      <td><em><code>(int)</code></em></td>
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
    <tr id="serviceaccount-name">
      <td><a href="#serviceaccount-name"><pre class="chroma"><code><span class="line"><span class="cl"><span class="nt">serviceAccount.name</span><span class="p">:</span><span class="w"> </span><span class="s2">&#34;&#34;</span></span></span></code></pre>
</a></td>
      <td><em><code>(string)</code></em></td>
    </tr>
    <tr id="tolerations">
      <td><a href="#tolerations"><pre class="chroma"><code><span class="line"><span class="cl"><span class="nt">tolerations</span><span class="p">:</span><span class="w"> </span><span class="p">[]</span></span></span></code></pre>
</a></td>
      <td><em><code>(list)</code></em></td>
    </tr>
    <tr id="vl-bearertoken">
      <td><a href="#vl-bearertoken"><pre class="chroma"><code><span class="line"><span class="cl"><span class="nt">vl.bearerToken</span><span class="p">:</span><span class="w"> </span><span class="s2">&#34;&#34;</span></span></span></code></pre>
</a></td>
      <td><em><code>(string)</code></em></td>
    </tr>
    <tr id="vl-entrypoint">
      <td><a href="#vl-entrypoint"><pre class="chroma"><code><span class="line"><span class="cl"><span class="nt">vl.entrypoint</span><span class="p">:</span><span class="w"> </span><span class="s2">&#34;&#34;</span></span></span></code></pre>
</a></td>
      <td><em><code>(string)</code></em></td>
    </tr>
    <tr id="vl-headers">
      <td><a href="#vl-headers"><pre class="chroma"><code><span class="line"><span class="cl"><span class="nt">vl.headers</span><span class="p">:</span><span class="w"> </span><span class="p">[]</span></span></span></code></pre>
</a></td>
      <td><em><code>(list)</code></em></td>
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

