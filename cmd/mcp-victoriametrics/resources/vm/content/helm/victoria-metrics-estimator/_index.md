

---
weight: 20
title: victoria-metrics-estimator
menu:
  docs:
    parent: helm
    weight: 20
    identifier: helm-victoria-metrics-estimator
url: /helm/victoria-metrics-estimator/
tags:
  - metrics
  - cardinality
  - kubernetes
---

![Version](https://img.shields.io/badge/0.2.2-gray?logo=Helm&labelColor=gray&link=https%3A%2F%2Fdocs.victoriametrics.com%2Fhelm%2Fvictoria-metrics-estimator%2Fchangelog%2F%23v022)
![ArtifactHub](https://img.shields.io/badge/ArtifactHub-informational?logoColor=white&color=417598&logo=artifacthub&link=https%3A%2F%2Fartifacthub.io%2Fpackages%2Fhelm%2Fvictoriametrics%2Fvictoria-metrics-estimator)
![License](https://img.shields.io/github/license/VictoriaMetrics/helm-charts?labelColor=green&label=&link=https%3A%2F%2Fgithub.com%2FVictoriaMetrics%2Fhelm-charts%2Fblob%2Fmaster%2FLICENSE)
![Slack](https://img.shields.io/badge/Join-4A154B?logo=slack&link=https%3A%2F%2Fslack.victoriametrics.com)
![X](https://img.shields.io/twitter/follow/VictoriaMetrics?style=flat&label=Follow&color=black&logo=x&labelColor=black&link=https%3A%2F%2Fx.com%2FVictoriaMetrics)
![Reddit](https://img.shields.io/reddit/subreddit-subscribers/VictoriaMetrics?style=flat&label=Join&labelColor=red&logoColor=white&logo=reddit&link=https%3A%2F%2Fwww.reddit.com%2Fr%2FVictoriaMetrics)

A Helm chart for VictoriaMetrics vmestimator

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

Install the chart in single-node mode:

```sh
helm install vmestimator vm/victoria-metrics-estimator
```

Configure vmagent to replicate remote write traffic to:

```text
http://vmestimator-single-vmestimator:8490/cardinality/api/v1/write
```

For cluster mode:

```sh
helm install vmestimator vm/victoria-metrics-estimator --set mode=cluster
```

Send remote write traffic to the `vmestimator-storage-vmestimator-insert` service and scrape cardinality estimates from the `vmestimator-select-vmestimator` service. See the [vmestimator documentation](https://docs.victoriametrics.com/victoriametrics/vmestimator/) for architecture and vmagent configuration details.

## How to install

Access a Kubernetes cluster.

### Setup chart repository (can be omitted for OCI repositories)

Add a chart helm repository with follow commands:

```console
helm repo add vm https://victoriametrics.github.io/helm-charts/

helm repo update
```
List versions of `vm/victoria-metrics-estimator` chart available to installation:

```console
helm search repo vm/victoria-metrics-estimator -l
```

### Install `victoria-metrics-estimator` chart

Export default values of `victoria-metrics-estimator` chart to file `values.yaml`:

  - For HTTPS repository

    ```console
    helm show values vm/victoria-metrics-estimator > values.yaml
    ```
  - For OCI repository

    ```console
    helm show values oci://ghcr.io/victoriametrics/helm-charts/victoria-metrics-estimator > values.yaml
    ```

Change the values according to the need of the environment in ``values.yaml`` file.

> Consider setting `.Values.nameOverride` to a small value like `vme` to avoid hitting resource name limits of 63 characters

Test the installation with command:

  - For HTTPS repository

    ```console
    helm install vme vm/victoria-metrics-estimator -f values.yaml -n NAMESPACE --debug --dry-run
    ```

  - For OCI repository

    ```console
    helm install vme oci://ghcr.io/victoriametrics/helm-charts/victoria-metrics-estimator -f values.yaml -n NAMESPACE --debug --dry-run
    ```

Install chart with command:

  - For HTTPS repository

    ```console
    helm install vme vm/victoria-metrics-estimator -f values.yaml -n NAMESPACE
    ```

  - For OCI repository

    ```console
    helm install vme oci://ghcr.io/victoriametrics/helm-charts/victoria-metrics-estimator -f values.yaml -n NAMESPACE
    ```

Get the pods lists by running this commands:

```console
kubectl get pods -A | grep 'vme'
```

Get the application by running this command:

```console
helm list -f vme -n NAMESPACE
```

See the history of versions of `vme` application with command.

```console
helm history vme -n NAMESPACE
```

## How to uninstall

Remove application with command.

```console
helm uninstall vme -n NAMESPACE
```

## Parameters

<table class="helm-vars">
  <thead>
    <th class="helm-vars-key">Key</th>
    <th class="helm-vars-description">Description</th>
  </thead>
  <tbody>
    <tr id="config-data">
      <td><a href="#config-data"><pre class="chroma"><code><span class="line"><span class="cl"><span class="nt">config.data</span><span class="p">:</span><span class="w">
</span></span></span><span class="line"><span class="cl"><span class="w">    </span><span class="nt">streams</span><span class="p">:</span><span class="w">
</span></span></span><span class="line"><span class="cl"><span class="w">        </span>- <span class="nt">interval</span><span class="p">:</span><span class="w"> </span><span class="l">5m</span><span class="w">
</span></span></span><span class="line"><span class="cl"><span class="w">        </span>- <span class="nt">group_by</span><span class="p">:</span><span class="w">
</span></span></span><span class="line"><span class="cl"><span class="w">            </span>- <span class="l">job</span><span class="w">
</span></span></span><span class="line"><span class="cl"><span class="w">          </span><span class="nt">interval</span><span class="p">:</span><span class="w"> </span><span class="l">5m</span><span class="w">
</span></span></span><span class="line"><span class="cl"><span class="w">        </span>- <span class="nt">group_by</span><span class="p">:</span><span class="w">
</span></span></span><span class="line"><span class="cl"><span class="w">            </span>- <span class="l">__name__</span><span class="w">
</span></span></span><span class="line"><span class="cl"><span class="w">          </span><span class="nt">interval</span><span class="p">:</span><span class="w"> </span><span class="l">5m</span></span></span></code></pre>
</a></td>
      <td><em><code>(object)</code></em><p>vmestimator stream configuration</p>
</td>
    </tr>
    <tr id="config-existingconfigmap">
      <td><a href="#config-existingconfigmap"><pre class="chroma"><code><span class="line"><span class="cl"><span class="nt">config.existingConfigMap</span><span class="p">:</span><span class="w"> </span><span class="s2">&#34;&#34;</span></span></span></code></pre>
</a></td>
      <td><em><code>(string)</code></em><p>Existing ConfigMap containing the vmestimator configuration</p>
</td>
    </tr>
    <tr id="config-key">
      <td><a href="#config-key"><pre class="chroma"><code><span class="line"><span class="cl"><span class="nt">config.key</span><span class="p">:</span><span class="w"> </span><span class="l">streams.yaml</span></span></span></code></pre>
</a></td>
      <td><em><code>(string)</code></em><p>Key in the ConfigMap containing the vmestimator configuration</p>
</td>
    </tr>
    <tr id="extraobjects">
      <td><a href="#extraobjects"><pre class="chroma"><code><span class="line"><span class="cl"><span class="nt">extraObjects</span><span class="p">:</span><span class="w"> </span><span class="p">[]</span></span></span></code></pre>
</a></td>
      <td><em><code>(list)</code></em><p>Additional Kubernetes resources rendered with tpl</p>
</td>
    </tr>
    <tr id="fullnameoverride">
      <td><a href="#fullnameoverride"><pre class="chroma"><code><span class="line"><span class="cl"><span class="nt">fullnameOverride</span><span class="p">:</span><span class="w"> </span><span class="s2">&#34;&#34;</span></span></span></code></pre>
</a></td>
      <td><em><code>(string)</code></em><p>Override the base name used for chart-wide shared resources (ConfigMap, ServiceAccount, NetworkPolicy) and, unless a component sets its own fullnameOverride, as the base for that component&rsquo;s resource names too</p>
</td>
    </tr>
    <tr id="global-cluster-dnsdomain">
      <td><a href="#global-cluster-dnsdomain"><pre class="chroma"><code><span class="line"><span class="cl"><span class="nt">global.cluster.dnsDomain</span><span class="p">:</span><span class="w"> </span><span class="l">cluster.local</span></span></span></code></pre>
</a></td>
      <td><em><code>(string)</code></em><p>Kubernetes cluster DNS suffix</p>
</td>
    </tr>
    <tr id="global-compatibility-openshift-adaptsecuritycontext">
      <td><a href="#global-compatibility-openshift-adaptsecuritycontext"><pre class="chroma"><code><span class="line"><span class="cl"><span class="nt">global.compatibility.openshift.adaptSecurityContext</span><span class="p">:</span><span class="w"> </span><span class="l">auto</span></span></span></code></pre>
</a></td>
      <td><em><code>(string)</code></em><p>OpenShift security context compatibility mode</p>
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
      <td><em><code>(string)</code></em><p>Image registry shared by all components</p>
</td>
    </tr>
    <tr id="global-imagepullsecrets">
      <td><a href="#global-imagepullsecrets"><pre class="chroma"><code><span class="line"><span class="cl"><span class="nt">global.imagePullSecrets</span><span class="p">:</span><span class="w"> </span><span class="p">[]</span></span></span></code></pre>
</a></td>
      <td><em><code>(list)</code></em><p>Image pull secrets shared by all components</p>
</td>
    </tr>
    <tr id="image-pullpolicy">
      <td><a href="#image-pullpolicy"><pre class="chroma"><code><span class="line"><span class="cl"><span class="nt">image.pullPolicy</span><span class="p">:</span><span class="w"> </span><span class="l">IfNotPresent</span></span></span></code></pre>
</a></td>
      <td><em><code>(string)</code></em><p>Image pull policy</p>
</td>
    </tr>
    <tr id="image-registry">
      <td><a href="#image-registry"><pre class="chroma"><code><span class="line"><span class="cl"><span class="nt">image.registry</span><span class="p">:</span><span class="w"> </span><span class="l">docker.io</span></span></span></code></pre>
</a></td>
      <td><em><code>(string)</code></em><p>Image registry</p>
</td>
    </tr>
    <tr id="image-repository">
      <td><a href="#image-repository"><pre class="chroma"><code><span class="line"><span class="cl"><span class="nt">image.repository</span><span class="p">:</span><span class="w"> </span><span class="l">victoriametrics/vmestimator</span></span></span></code></pre>
</a></td>
      <td><em><code>(string)</code></em><p>Image repository</p>
</td>
    </tr>
    <tr id="image-tag">
      <td><a href="#image-tag"><pre class="chroma"><code><span class="line"><span class="cl"><span class="nt">image.tag</span><span class="p">:</span><span class="w"> </span><span class="s2">&#34;&#34;</span></span></span></code></pre>
</a></td>
      <td><em><code>(string)</code></em><p>Image tag; defaults to Chart.appVersion</p>
</td>
    </tr>
    <tr id="imagepullsecrets">
      <td><a href="#imagepullsecrets"><pre class="chroma"><code><span class="line"><span class="cl"><span class="nt">imagePullSecrets</span><span class="p">:</span><span class="w"> </span><span class="p">[]</span></span></span></code></pre>
</a></td>
      <td><em><code>(list)</code></em><p>Image pull secrets</p>
</td>
    </tr>
    <tr id="ingress-annotations">
      <td><a href="#ingress-annotations"><pre class="chroma"><code><span class="line"><span class="cl"><span class="nt">ingress.annotations</span><span class="p">:</span><span class="w"> </span>{}</span></span></code></pre>
</a></td>
      <td><em><code>(object)</code></em><p>Ingress annotations</p>
</td>
    </tr>
    <tr id="ingress-classname">
      <td><a href="#ingress-classname"><pre class="chroma"><code><span class="line"><span class="cl"><span class="nt">ingress.className</span><span class="p">:</span><span class="w"> </span><span class="s2">&#34;&#34;</span></span></span></code></pre>
</a></td>
      <td><em><code>(string)</code></em><p>Ingress class name</p>
</td>
    </tr>
    <tr id="ingress-enabled">
      <td><a href="#ingress-enabled"><pre class="chroma"><code><span class="line"><span class="cl"><span class="nt">ingress.enabled</span><span class="p">:</span><span class="w"> </span><span class="kc">false</span></span></span></code></pre>
</a></td>
      <td><em><code>(bool)</code></em><p>Enable ingress for the public single/select service</p>
</td>
    </tr>
    <tr id="ingress-hosts[0]-host">
      <td><a href="#ingress-hosts[0]-host"><pre class="chroma"><code><span class="line"><span class="cl"><span class="nt">ingress.hosts[0].host</span><span class="p">:</span><span class="w"> </span><span class="l">vmestimator.local</span></span></span></code></pre>
</a></td>
      <td><em><code>(string)</code></em></td>
    </tr>
    <tr id="ingress-hosts[0]-paths[0]-path">
      <td><a href="#ingress-hosts[0]-paths[0]-path"><pre class="chroma"><code><span class="line"><span class="cl"><span class="nt">ingress.hosts[0].paths[0].path</span><span class="p">:</span><span class="w"> </span><span class="l">/</span></span></span></code></pre>
</a></td>
      <td><em><code>(string)</code></em></td>
    </tr>
    <tr id="ingress-hosts[0]-paths[0]-pathtype">
      <td><a href="#ingress-hosts[0]-paths[0]-pathtype"><pre class="chroma"><code><span class="line"><span class="cl"><span class="nt">ingress.hosts[0].paths[0].pathType</span><span class="p">:</span><span class="w"> </span><span class="l">Prefix</span></span></span></code></pre>
</a></td>
      <td><em><code>(string)</code></em></td>
    </tr>
    <tr id="ingress-tls">
      <td><a href="#ingress-tls"><pre class="chroma"><code><span class="line"><span class="cl"><span class="nt">ingress.tls</span><span class="p">:</span><span class="w"> </span><span class="p">[]</span></span></span></code></pre>
</a></td>
      <td><em><code>(list)</code></em></td>
    </tr>
    <tr id="mode">
      <td><a href="#mode"><pre class="chroma"><code><span class="line"><span class="cl"><span class="nt">mode</span><span class="p">:</span><span class="w"> </span><span class="l">single</span></span></span></code></pre>
</a></td>
      <td><em><code>(string)</code></em><p>Deployment mode: single or cluster</p>
</td>
    </tr>
    <tr id="nameoverride">
      <td><a href="#nameoverride"><pre class="chroma"><code><span class="line"><span class="cl"><span class="nt">nameOverride</span><span class="p">:</span><span class="w"> </span><span class="s2">&#34;&#34;</span></span></span></code></pre>
</a></td>
      <td><em><code>(string)</code></em><p>Override chart name</p>
</td>
    </tr>
    <tr id="networkpolicy-annotations">
      <td><a href="#networkpolicy-annotations"><pre class="chroma"><code><span class="line"><span class="cl"><span class="nt">networkPolicy.annotations</span><span class="p">:</span><span class="w"> </span>{}</span></span></code></pre>
</a></td>
      <td><em><code>(object)</code></em><p>Additional annotations</p>
</td>
    </tr>
    <tr id="networkpolicy-egress">
      <td><a href="#networkpolicy-egress"><pre class="chroma"><code><span class="line"><span class="cl"><span class="nt">networkPolicy.egress</span><span class="p">:</span><span class="w"> </span><span class="p">[]</span></span></span></code></pre>
</a></td>
      <td><em><code>(list)</code></em><p>Egress rules</p>
</td>
    </tr>
    <tr id="networkpolicy-enabled">
      <td><a href="#networkpolicy-enabled"><pre class="chroma"><code><span class="line"><span class="cl"><span class="nt">networkPolicy.enabled</span><span class="p">:</span><span class="w"> </span><span class="kc">false</span></span></span></code></pre>
</a></td>
      <td><em><code>(bool)</code></em><p>Enable NetworkPolicy</p>
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
      <td><em><code>(object)</code></em><p>Additional labels</p>
</td>
    </tr>
    <tr id="podsecuritycontext-enabled">
      <td><a href="#podsecuritycontext-enabled"><pre class="chroma"><code><span class="line"><span class="cl"><span class="nt">podSecurityContext.enabled</span><span class="p">:</span><span class="w"> </span><span class="kc">true</span></span></span></code></pre>
</a></td>
      <td><em><code>(bool)</code></em></td>
    </tr>
    <tr id="podsecuritycontext-fsgroup">
      <td><a href="#podsecuritycontext-fsgroup"><pre class="chroma"><code><span class="line"><span class="cl"><span class="nt">podSecurityContext.fsGroup</span><span class="p">:</span><span class="w"> </span><span class="m">2000</span></span></span></code></pre>
</a></td>
      <td><em><code>(int)</code></em></td>
    </tr>
    <tr id="podsecuritycontext-runasgroup">
      <td><a href="#podsecuritycontext-runasgroup"><pre class="chroma"><code><span class="line"><span class="cl"><span class="nt">podSecurityContext.runAsGroup</span><span class="p">:</span><span class="w"> </span><span class="m">1000</span></span></span></code></pre>
</a></td>
      <td><em><code>(int)</code></em></td>
    </tr>
    <tr id="podsecuritycontext-runasnonroot">
      <td><a href="#podsecuritycontext-runasnonroot"><pre class="chroma"><code><span class="line"><span class="cl"><span class="nt">podSecurityContext.runAsNonRoot</span><span class="p">:</span><span class="w"> </span><span class="kc">true</span></span></span></code></pre>
</a></td>
      <td><em><code>(bool)</code></em></td>
    </tr>
    <tr id="podsecuritycontext-runasuser">
      <td><a href="#podsecuritycontext-runasuser"><pre class="chroma"><code><span class="line"><span class="cl"><span class="nt">podSecurityContext.runAsUser</span><span class="p">:</span><span class="w"> </span><span class="m">1000</span></span></span></code></pre>
</a></td>
      <td><em><code>(int)</code></em></td>
    </tr>
    <tr id="probe-liveness-httpget">
      <td><a href="#probe-liveness-httpget"><pre class="chroma"><code><span class="line"><span class="cl"><span class="nt">probe.liveness.httpGet</span><span class="p">:</span><span class="w"> </span>{}</span></span></code></pre>
</a></td>
      <td><em><code>(object)</code></em></td>
    </tr>
    <tr id="probe-liveness-initialdelayseconds">
      <td><a href="#probe-liveness-initialdelayseconds"><pre class="chroma"><code><span class="line"><span class="cl"><span class="nt">probe.liveness.initialDelaySeconds</span><span class="p">:</span><span class="w"> </span><span class="m">30</span></span></span></code></pre>
</a></td>
      <td><em><code>(int)</code></em></td>
    </tr>
    <tr id="probe-liveness-periodseconds">
      <td><a href="#probe-liveness-periodseconds"><pre class="chroma"><code><span class="line"><span class="cl"><span class="nt">probe.liveness.periodSeconds</span><span class="p">:</span><span class="w"> </span><span class="m">30</span></span></span></code></pre>
</a></td>
      <td><em><code>(int)</code></em></td>
    </tr>
    <tr id="probe-readiness-httpget">
      <td><a href="#probe-readiness-httpget"><pre class="chroma"><code><span class="line"><span class="cl"><span class="nt">probe.readiness.httpGet</span><span class="p">:</span><span class="w"> </span>{}</span></span></code></pre>
</a></td>
      <td><em><code>(object)</code></em></td>
    </tr>
    <tr id="probe-readiness-initialdelayseconds">
      <td><a href="#probe-readiness-initialdelayseconds"><pre class="chroma"><code><span class="line"><span class="cl"><span class="nt">probe.readiness.initialDelaySeconds</span><span class="p">:</span><span class="w"> </span><span class="m">5</span></span></span></code></pre>
</a></td>
      <td><em><code>(int)</code></em></td>
    </tr>
    <tr id="probe-readiness-periodseconds">
      <td><a href="#probe-readiness-periodseconds"><pre class="chroma"><code><span class="line"><span class="cl"><span class="nt">probe.readiness.periodSeconds</span><span class="p">:</span><span class="w"> </span><span class="m">5</span></span></span></code></pre>
</a></td>
      <td><em><code>(int)</code></em></td>
    </tr>
    <tr id="securitycontext-allowprivilegeescalation">
      <td><a href="#securitycontext-allowprivilegeescalation"><pre class="chroma"><code><span class="line"><span class="cl"><span class="nt">securityContext.allowPrivilegeEscalation</span><span class="p">:</span><span class="w"> </span><span class="kc">false</span></span></span></code></pre>
</a></td>
      <td><em><code>(bool)</code></em></td>
    </tr>
    <tr id="securitycontext-capabilities-drop[0]">
      <td><a href="#securitycontext-capabilities-drop[0]"><pre class="chroma"><code><span class="line"><span class="cl"><span class="nt">securityContext.capabilities.drop[0]</span><span class="p">:</span><span class="w"> </span><span class="l">ALL</span></span></span></code></pre>
</a></td>
      <td><em><code>(string)</code></em></td>
    </tr>
    <tr id="securitycontext-enabled">
      <td><a href="#securitycontext-enabled"><pre class="chroma"><code><span class="line"><span class="cl"><span class="nt">securityContext.enabled</span><span class="p">:</span><span class="w"> </span><span class="kc">true</span></span></span></code></pre>
</a></td>
      <td><em><code>(bool)</code></em></td>
    </tr>
    <tr id="securitycontext-readonlyrootfilesystem">
      <td><a href="#securitycontext-readonlyrootfilesystem"><pre class="chroma"><code><span class="line"><span class="cl"><span class="nt">securityContext.readOnlyRootFilesystem</span><span class="p">:</span><span class="w"> </span><span class="kc">true</span></span></span></code></pre>
</a></td>
      <td><em><code>(bool)</code></em></td>
    </tr>
    <tr id="select-affinity">
      <td><a href="#select-affinity"><pre class="chroma"><code><span class="line"><span class="cl"><span class="nt">select.affinity</span><span class="p">:</span><span class="w"> </span>{}</span></span></code></pre>
</a></td>
      <td><em><code>(object)</code></em></td>
    </tr>
    <tr id="select-command">
      <td><a href="#select-command"><pre class="chroma"><code><span class="line"><span class="cl"><span class="nt">select.command</span><span class="p">:</span><span class="w"> </span><span class="p">[]</span></span></span></code></pre>
</a></td>
      <td><em><code>(list)</code></em></td>
    </tr>
    <tr id="select-dnsconfig">
      <td><a href="#select-dnsconfig"><pre class="chroma"><code><span class="line"><span class="cl"><span class="nt">select.dnsConfig</span><span class="p">:</span><span class="w"> </span>{}</span></span></code></pre>
</a></td>
      <td><em><code>(object)</code></em></td>
    </tr>
    <tr id="select-env">
      <td><a href="#select-env"><pre class="chroma"><code><span class="line"><span class="cl"><span class="nt">select.env</span><span class="p">:</span><span class="w"> </span><span class="p">[]</span></span></span></code></pre>
</a></td>
      <td><em><code>(list)</code></em></td>
    </tr>
    <tr id="select-envfrom">
      <td><a href="#select-envfrom"><pre class="chroma"><code><span class="line"><span class="cl"><span class="nt">select.envFrom</span><span class="p">:</span><span class="w"> </span><span class="p">[]</span></span></span></code></pre>
</a></td>
      <td><em><code>(list)</code></em></td>
    </tr>
    <tr id="select-extraargs">
      <td><a href="#select-extraargs"><pre class="chroma"><code><span class="line"><span class="cl"><span class="nt">select.extraArgs</span><span class="p">:</span><span class="w">
</span></span></span><span class="line"><span class="cl"><span class="w">    </span><span class="nt">envflag.enable</span><span class="p">:</span><span class="w"> </span><span class="kc">true</span><span class="w">
</span></span></span><span class="line"><span class="cl"><span class="w">    </span><span class="nt">envflag.prefix</span><span class="p">:</span><span class="w"> </span><span class="l">VM_</span><span class="w">
</span></span></span><span class="line"><span class="cl"><span class="w">    </span><span class="nt">loggerFormat</span><span class="p">:</span><span class="w"> </span><span class="l">json</span></span></span></code></pre>
</a></td>
      <td><em><code>(object)</code></em><p>Additional command-line arguments</p>
</td>
    </tr>
    <tr id="select-extravolumemounts">
      <td><a href="#select-extravolumemounts"><pre class="chroma"><code><span class="line"><span class="cl"><span class="nt">select.extraVolumeMounts</span><span class="p">:</span><span class="w"> </span><span class="p">[]</span></span></span></code></pre>
</a></td>
      <td><em><code>(list)</code></em></td>
    </tr>
    <tr id="select-extravolumes">
      <td><a href="#select-extravolumes"><pre class="chroma"><code><span class="line"><span class="cl"><span class="nt">select.extraVolumes</span><span class="p">:</span><span class="w"> </span><span class="p">[]</span></span></span></code></pre>
</a></td>
      <td><em><code>(list)</code></em></td>
    </tr>
    <tr id="select-fullnameoverride">
      <td><a href="#select-fullnameoverride"><pre class="chroma"><code><span class="line"><span class="cl"><span class="nt">select.fullnameOverride</span><span class="p">:</span><span class="w"> </span><span class="s2">&#34;&#34;</span></span></span></code></pre>
</a></td>
      <td><em><code>(string)</code></em><p>Overrides the full name of the select component</p>
</td>
    </tr>
    <tr id="select-nodeselector">
      <td><a href="#select-nodeselector"><pre class="chroma"><code><span class="line"><span class="cl"><span class="nt">select.nodeSelector</span><span class="p">:</span><span class="w"> </span>{}</span></span></code></pre>
</a></td>
      <td><em><code>(object)</code></em></td>
    </tr>
    <tr id="select-podannotations">
      <td><a href="#select-podannotations"><pre class="chroma"><code><span class="line"><span class="cl"><span class="nt">select.podAnnotations</span><span class="p">:</span><span class="w"> </span>{}</span></span></code></pre>
</a></td>
      <td><em><code>(object)</code></em></td>
    </tr>
    <tr id="select-poddisruptionbudget-enabled">
      <td><a href="#select-poddisruptionbudget-enabled"><pre class="chroma"><code><span class="line"><span class="cl"><span class="nt">select.podDisruptionBudget.enabled</span><span class="p">:</span><span class="w"> </span><span class="kc">false</span></span></span></code></pre>
</a></td>
      <td><em><code>(bool)</code></em></td>
    </tr>
    <tr id="select-poddisruptionbudget-maxunavailable">
      <td><a href="#select-poddisruptionbudget-maxunavailable"><pre class="chroma"><code><span class="line"><span class="cl"><span class="nt">select.podDisruptionBudget.maxUnavailable</span><span class="p">:</span><span class="w"> </span><span class="m">0</span></span></span></code></pre>
</a></td>
      <td><em><code>(int)</code></em></td>
    </tr>
    <tr id="select-poddisruptionbudget-minavailable">
      <td><a href="#select-poddisruptionbudget-minavailable"><pre class="chroma"><code><span class="line"><span class="cl"><span class="nt">select.podDisruptionBudget.minAvailable</span><span class="p">:</span><span class="w"> </span><span class="m">0</span></span></span></code></pre>
</a></td>
      <td><em><code>(int)</code></em></td>
    </tr>
    <tr id="select-poddisruptionbudget-unhealthypodevictionpolicy">
      <td><a href="#select-poddisruptionbudget-unhealthypodevictionpolicy"><pre class="chroma"><code><span class="line"><span class="cl"><span class="nt">select.podDisruptionBudget.unhealthyPodEvictionPolicy</span><span class="p">:</span><span class="w"> </span><span class="kc">null</span></span></span></code></pre>
</a></td>
      <td><em><code>(string)</code></em></td>
    </tr>
    <tr id="select-podlabels">
      <td><a href="#select-podlabels"><pre class="chroma"><code><span class="line"><span class="cl"><span class="nt">select.podLabels</span><span class="p">:</span><span class="w"> </span>{}</span></span></code></pre>
</a></td>
      <td><em><code>(object)</code></em></td>
    </tr>
    <tr id="select-replicacount">
      <td><a href="#select-replicacount"><pre class="chroma"><code><span class="line"><span class="cl"><span class="nt">select.replicaCount</span><span class="p">:</span><span class="w"> </span><span class="m">2</span></span></span></code></pre>
</a></td>
      <td><em><code>(int)</code></em><p>Number of selector nodes in cluster mode</p>
</td>
    </tr>
    <tr id="select-resources">
      <td><a href="#select-resources"><pre class="chroma"><code><span class="line"><span class="cl"><span class="nt">select.resources</span><span class="p">:</span><span class="w"> </span>{}</span></span></code></pre>
</a></td>
      <td><em><code>(object)</code></em></td>
    </tr>
    <tr id="select-service-annotations">
      <td><a href="#select-service-annotations"><pre class="chroma"><code><span class="line"><span class="cl"><span class="nt">select.service.annotations</span><span class="p">:</span><span class="w"> </span>{}</span></span></code></pre>
</a></td>
      <td><em><code>(object)</code></em></td>
    </tr>
    <tr id="select-service-labels">
      <td><a href="#select-service-labels"><pre class="chroma"><code><span class="line"><span class="cl"><span class="nt">select.service.labels</span><span class="p">:</span><span class="w"> </span>{}</span></span></code></pre>
</a></td>
      <td><em><code>(object)</code></em></td>
    </tr>
    <tr id="select-service-port">
      <td><a href="#select-service-port"><pre class="chroma"><code><span class="line"><span class="cl"><span class="nt">select.service.port</span><span class="p">:</span><span class="w"> </span><span class="m">8490</span></span></span></code></pre>
</a></td>
      <td><em><code>(int)</code></em></td>
    </tr>
    <tr id="select-service-type">
      <td><a href="#select-service-type"><pre class="chroma"><code><span class="line"><span class="cl"><span class="nt">select.service.type</span><span class="p">:</span><span class="w"> </span><span class="l">ClusterIP</span></span></span></code></pre>
</a></td>
      <td><em><code>(string)</code></em></td>
    </tr>
    <tr id="select-tolerations">
      <td><a href="#select-tolerations"><pre class="chroma"><code><span class="line"><span class="cl"><span class="nt">select.tolerations</span><span class="p">:</span><span class="w"> </span><span class="p">[]</span></span></span></code></pre>
</a></td>
      <td><em><code>(list)</code></em></td>
    </tr>
    <tr id="select-topologyspreadconstraints">
      <td><a href="#select-topologyspreadconstraints"><pre class="chroma"><code><span class="line"><span class="cl"><span class="nt">select.topologySpreadConstraints</span><span class="p">:</span><span class="w"> </span><span class="p">[]</span></span></span></code></pre>
</a></td>
      <td><em><code>(list)</code></em></td>
    </tr>
    <tr id="serviceaccount-annotations">
      <td><a href="#serviceaccount-annotations"><pre class="chroma"><code><span class="line"><span class="cl"><span class="nt">serviceAccount.annotations</span><span class="p">:</span><span class="w"> </span>{}</span></span></code></pre>
</a></td>
      <td><em><code>(object)</code></em><p>Service account annotations</p>
</td>
    </tr>
    <tr id="serviceaccount-automounttoken">
      <td><a href="#serviceaccount-automounttoken"><pre class="chroma"><code><span class="line"><span class="cl"><span class="nt">serviceAccount.automountToken</span><span class="p">:</span><span class="w"> </span><span class="kc">false</span></span></span></code></pre>
</a></td>
      <td><em><code>(bool)</code></em><p>Mount the service account token</p>
</td>
    </tr>
    <tr id="serviceaccount-create">
      <td><a href="#serviceaccount-create"><pre class="chroma"><code><span class="line"><span class="cl"><span class="nt">serviceAccount.create</span><span class="p">:</span><span class="w"> </span><span class="kc">false</span></span></span></code></pre>
</a></td>
      <td><em><code>(bool)</code></em><p>Create a service account</p>
</td>
    </tr>
    <tr id="serviceaccount-name">
      <td><a href="#serviceaccount-name"><pre class="chroma"><code><span class="line"><span class="cl"><span class="nt">serviceAccount.name</span><span class="p">:</span><span class="w"> </span><span class="s2">&#34;&#34;</span></span></span></code></pre>
</a></td>
      <td><em><code>(string)</code></em><p>Service account name</p>
</td>
    </tr>
    <tr id="servicemonitor-annotations">
      <td><a href="#servicemonitor-annotations"><pre class="chroma"><code><span class="line"><span class="cl"><span class="nt">serviceMonitor.annotations</span><span class="p">:</span><span class="w"> </span>{}</span></span></code></pre>
</a></td>
      <td><em><code>(object)</code></em><p>Additional annotations</p>
</td>
    </tr>
    <tr id="servicemonitor-enabled">
      <td><a href="#servicemonitor-enabled"><pre class="chroma"><code><span class="line"><span class="cl"><span class="nt">serviceMonitor.enabled</span><span class="p">:</span><span class="w"> </span><span class="kc">false</span></span></span></code></pre>
</a></td>
      <td><em><code>(bool)</code></em><p>Create a VMServiceScrape for the public single/select service</p>
</td>
    </tr>
    <tr id="servicemonitor-extralabels">
      <td><a href="#servicemonitor-extralabels"><pre class="chroma"><code><span class="line"><span class="cl"><span class="nt">serviceMonitor.extraLabels</span><span class="p">:</span><span class="w"> </span>{}</span></span></code></pre>
</a></td>
      <td><em><code>(object)</code></em><p>Additional labels</p>
</td>
    </tr>
    <tr id="servicemonitor-interval">
      <td><a href="#servicemonitor-interval"><pre class="chroma"><code><span class="line"><span class="cl"><span class="nt">serviceMonitor.interval</span><span class="p">:</span><span class="w"> </span><span class="s2">&#34;&#34;</span></span></span></code></pre>
</a></td>
      <td><em><code>(string)</code></em><p>Scrape interval</p>
</td>
    </tr>
    <tr id="servicemonitor-metricrelabelings">
      <td><a href="#servicemonitor-metricrelabelings"><pre class="chroma"><code><span class="line"><span class="cl"><span class="nt">serviceMonitor.metricRelabelings</span><span class="p">:</span><span class="w"> </span><span class="p">[]</span></span></span></code></pre>
</a></td>
      <td><em><code>(list)</code></em><p>Metric relabeling rules</p>
</td>
    </tr>
    <tr id="single-affinity">
      <td><a href="#single-affinity"><pre class="chroma"><code><span class="line"><span class="cl"><span class="nt">single.affinity</span><span class="p">:</span><span class="w"> </span>{}</span></span></code></pre>
</a></td>
      <td><em><code>(object)</code></em><p>Pod affinity</p>
</td>
    </tr>
    <tr id="single-command">
      <td><a href="#single-command"><pre class="chroma"><code><span class="line"><span class="cl"><span class="nt">single.command</span><span class="p">:</span><span class="w"> </span><span class="p">[]</span></span></span></code></pre>
</a></td>
      <td><em><code>(list)</code></em><p>Override the container command</p>
</td>
    </tr>
    <tr id="single-dnsconfig">
      <td><a href="#single-dnsconfig"><pre class="chroma"><code><span class="line"><span class="cl"><span class="nt">single.dnsConfig</span><span class="p">:</span><span class="w"> </span>{}</span></span></code></pre>
</a></td>
      <td><em><code>(object)</code></em><p>Custom DNS config for pod</p>
</td>
    </tr>
    <tr id="single-env">
      <td><a href="#single-env"><pre class="chroma"><code><span class="line"><span class="cl"><span class="nt">single.env</span><span class="p">:</span><span class="w"> </span><span class="p">[]</span></span></span></code></pre>
</a></td>
      <td><em><code>(list)</code></em><p>Additional environment variables</p>
</td>
    </tr>
    <tr id="single-envfrom">
      <td><a href="#single-envfrom"><pre class="chroma"><code><span class="line"><span class="cl"><span class="nt">single.envFrom</span><span class="p">:</span><span class="w"> </span><span class="p">[]</span></span></span></code></pre>
</a></td>
      <td><em><code>(list)</code></em><p>Additional environment sources</p>
</td>
    </tr>
    <tr id="single-extraargs">
      <td><a href="#single-extraargs"><pre class="chroma"><code><span class="line"><span class="cl"><span class="nt">single.extraArgs</span><span class="p">:</span><span class="w">
</span></span></span><span class="line"><span class="cl"><span class="w">    </span><span class="nt">envflag.enable</span><span class="p">:</span><span class="w"> </span><span class="kc">true</span><span class="w">
</span></span></span><span class="line"><span class="cl"><span class="w">    </span><span class="nt">envflag.prefix</span><span class="p">:</span><span class="w"> </span><span class="l">VM_</span><span class="w">
</span></span></span><span class="line"><span class="cl"><span class="w">    </span><span class="nt">loggerFormat</span><span class="p">:</span><span class="w"> </span><span class="l">json</span></span></span></code></pre>
</a></td>
      <td><em><code>(object)</code></em><p>Additional command-line arguments</p>
</td>
    </tr>
    <tr id="single-extravolumemounts">
      <td><a href="#single-extravolumemounts"><pre class="chroma"><code><span class="line"><span class="cl"><span class="nt">single.extraVolumeMounts</span><span class="p">:</span><span class="w"> </span><span class="p">[]</span></span></span></code></pre>
</a></td>
      <td><em><code>(list)</code></em><p>Additional volume mounts</p>
</td>
    </tr>
    <tr id="single-extravolumes">
      <td><a href="#single-extravolumes"><pre class="chroma"><code><span class="line"><span class="cl"><span class="nt">single.extraVolumes</span><span class="p">:</span><span class="w"> </span><span class="p">[]</span></span></span></code></pre>
</a></td>
      <td><em><code>(list)</code></em><p>Additional volumes</p>
</td>
    </tr>
    <tr id="single-fullnameoverride">
      <td><a href="#single-fullnameoverride"><pre class="chroma"><code><span class="line"><span class="cl"><span class="nt">single.fullnameOverride</span><span class="p">:</span><span class="w"> </span><span class="s2">&#34;&#34;</span></span></span></code></pre>
</a></td>
      <td><em><code>(string)</code></em><p>Overrides the full name of the single-node component</p>
</td>
    </tr>
    <tr id="single-nodeselector">
      <td><a href="#single-nodeselector"><pre class="chroma"><code><span class="line"><span class="cl"><span class="nt">single.nodeSelector</span><span class="p">:</span><span class="w"> </span>{}</span></span></code></pre>
</a></td>
      <td><em><code>(object)</code></em><p>Pod node selector</p>
</td>
    </tr>
    <tr id="single-podannotations">
      <td><a href="#single-podannotations"><pre class="chroma"><code><span class="line"><span class="cl"><span class="nt">single.podAnnotations</span><span class="p">:</span><span class="w"> </span>{}</span></span></code></pre>
</a></td>
      <td><em><code>(object)</code></em><p>Pod annotations</p>
</td>
    </tr>
    <tr id="single-poddisruptionbudget-enabled">
      <td><a href="#single-poddisruptionbudget-enabled"><pre class="chroma"><code><span class="line"><span class="cl"><span class="nt">single.podDisruptionBudget.enabled</span><span class="p">:</span><span class="w"> </span><span class="kc">false</span></span></span></code></pre>
</a></td>
      <td><em><code>(bool)</code></em><p>Enable a PodDisruptionBudget</p>
</td>
    </tr>
    <tr id="single-poddisruptionbudget-maxunavailable">
      <td><a href="#single-poddisruptionbudget-maxunavailable"><pre class="chroma"><code><span class="line"><span class="cl"><span class="nt">single.podDisruptionBudget.maxUnavailable</span><span class="p">:</span><span class="w"> </span><span class="m">0</span></span></span></code></pre>
</a></td>
      <td><em><code>(int)</code></em><p>Maximum unavailable pods</p>
</td>
    </tr>
    <tr id="single-poddisruptionbudget-minavailable">
      <td><a href="#single-poddisruptionbudget-minavailable"><pre class="chroma"><code><span class="line"><span class="cl"><span class="nt">single.podDisruptionBudget.minAvailable</span><span class="p">:</span><span class="w"> </span><span class="m">0</span></span></span></code></pre>
</a></td>
      <td><em><code>(int)</code></em><p>Minimum available pods</p>
</td>
    </tr>
    <tr id="single-poddisruptionbudget-unhealthypodevictionpolicy">
      <td><a href="#single-poddisruptionbudget-unhealthypodevictionpolicy"><pre class="chroma"><code><span class="line"><span class="cl"><span class="nt">single.podDisruptionBudget.unhealthyPodEvictionPolicy</span><span class="p">:</span><span class="w"> </span><span class="kc">null</span></span></span></code></pre>
</a></td>
      <td><em><code>(string)</code></em><p>Policy for evicting unhealthy pods</p>
</td>
    </tr>
    <tr id="single-podlabels">
      <td><a href="#single-podlabels"><pre class="chroma"><code><span class="line"><span class="cl"><span class="nt">single.podLabels</span><span class="p">:</span><span class="w"> </span>{}</span></span></code></pre>
</a></td>
      <td><em><code>(object)</code></em><p>Pod labels</p>
</td>
    </tr>
    <tr id="single-replicacount">
      <td><a href="#single-replicacount"><pre class="chroma"><code><span class="line"><span class="cl"><span class="nt">single.replicaCount</span><span class="p">:</span><span class="w"> </span><span class="m">1</span></span></span></code></pre>
</a></td>
      <td><em><code>(int)</code></em><p>Number of single-node replicas</p>
</td>
    </tr>
    <tr id="single-resources">
      <td><a href="#single-resources"><pre class="chroma"><code><span class="line"><span class="cl"><span class="nt">single.resources</span><span class="p">:</span><span class="w"> </span>{}</span></span></code></pre>
</a></td>
      <td><em><code>(object)</code></em><p>Container resources</p>
</td>
    </tr>
    <tr id="single-service-annotations">
      <td><a href="#single-service-annotations"><pre class="chroma"><code><span class="line"><span class="cl"><span class="nt">single.service.annotations</span><span class="p">:</span><span class="w"> </span>{}</span></span></code></pre>
</a></td>
      <td><em><code>(object)</code></em><p>Service annotations</p>
</td>
    </tr>
    <tr id="single-service-labels">
      <td><a href="#single-service-labels"><pre class="chroma"><code><span class="line"><span class="cl"><span class="nt">single.service.labels</span><span class="p">:</span><span class="w"> </span>{}</span></span></code></pre>
</a></td>
      <td><em><code>(object)</code></em><p>Service labels</p>
</td>
    </tr>
    <tr id="single-service-port">
      <td><a href="#single-service-port"><pre class="chroma"><code><span class="line"><span class="cl"><span class="nt">single.service.port</span><span class="p">:</span><span class="w"> </span><span class="m">8490</span></span></span></code></pre>
</a></td>
      <td><em><code>(int)</code></em><p>Service port</p>
</td>
    </tr>
    <tr id="single-service-type">
      <td><a href="#single-service-type"><pre class="chroma"><code><span class="line"><span class="cl"><span class="nt">single.service.type</span><span class="p">:</span><span class="w"> </span><span class="l">ClusterIP</span></span></span></code></pre>
</a></td>
      <td><em><code>(string)</code></em><p>Service type</p>
</td>
    </tr>
    <tr id="single-tolerations">
      <td><a href="#single-tolerations"><pre class="chroma"><code><span class="line"><span class="cl"><span class="nt">single.tolerations</span><span class="p">:</span><span class="w"> </span><span class="p">[]</span></span></span></code></pre>
</a></td>
      <td><em><code>(list)</code></em><p>Pod tolerations</p>
</td>
    </tr>
    <tr id="single-topologyspreadconstraints">
      <td><a href="#single-topologyspreadconstraints"><pre class="chroma"><code><span class="line"><span class="cl"><span class="nt">single.topologySpreadConstraints</span><span class="p">:</span><span class="w"> </span><span class="p">[]</span></span></span></code></pre>
</a></td>
      <td><em><code>(list)</code></em><p>Pod topology spread constraints</p>
</td>
    </tr>
    <tr id="storage-affinity">
      <td><a href="#storage-affinity"><pre class="chroma"><code><span class="line"><span class="cl"><span class="nt">storage.affinity</span><span class="p">:</span><span class="w"> </span>{}</span></span></code></pre>
</a></td>
      <td><em><code>(object)</code></em></td>
    </tr>
    <tr id="storage-command">
      <td><a href="#storage-command"><pre class="chroma"><code><span class="line"><span class="cl"><span class="nt">storage.command</span><span class="p">:</span><span class="w"> </span><span class="p">[]</span></span></span></code></pre>
</a></td>
      <td><em><code>(list)</code></em></td>
    </tr>
    <tr id="storage-dnsconfig">
      <td><a href="#storage-dnsconfig"><pre class="chroma"><code><span class="line"><span class="cl"><span class="nt">storage.dnsConfig</span><span class="p">:</span><span class="w"> </span>{}</span></span></code></pre>
</a></td>
      <td><em><code>(object)</code></em></td>
    </tr>
    <tr id="storage-env">
      <td><a href="#storage-env"><pre class="chroma"><code><span class="line"><span class="cl"><span class="nt">storage.env</span><span class="p">:</span><span class="w"> </span><span class="p">[]</span></span></span></code></pre>
</a></td>
      <td><em><code>(list)</code></em></td>
    </tr>
    <tr id="storage-envfrom">
      <td><a href="#storage-envfrom"><pre class="chroma"><code><span class="line"><span class="cl"><span class="nt">storage.envFrom</span><span class="p">:</span><span class="w"> </span><span class="p">[]</span></span></span></code></pre>
</a></td>
      <td><em><code>(list)</code></em></td>
    </tr>
    <tr id="storage-extraargs">
      <td><a href="#storage-extraargs"><pre class="chroma"><code><span class="line"><span class="cl"><span class="nt">storage.extraArgs</span><span class="p">:</span><span class="w">
</span></span></span><span class="line"><span class="cl"><span class="w">    </span><span class="nt">cardinalityMetrics.exposeAt</span><span class="p">:</span><span class="w"> </span><span class="l">/cardinality/metrics</span><span class="w">
</span></span></span><span class="line"><span class="cl"><span class="w">    </span><span class="nt">envflag.enable</span><span class="p">:</span><span class="w"> </span><span class="kc">true</span><span class="w">
</span></span></span><span class="line"><span class="cl"><span class="w">    </span><span class="nt">envflag.prefix</span><span class="p">:</span><span class="w"> </span><span class="l">VM_</span><span class="w">
</span></span></span><span class="line"><span class="cl"><span class="w">    </span><span class="nt">loggerFormat</span><span class="p">:</span><span class="w"> </span><span class="l">json</span></span></span></code></pre>
</a></td>
      <td><em><code>(object)</code></em><p>Additional command-line arguments</p>
</td>
    </tr>
    <tr id="storage-extravolumemounts">
      <td><a href="#storage-extravolumemounts"><pre class="chroma"><code><span class="line"><span class="cl"><span class="nt">storage.extraVolumeMounts</span><span class="p">:</span><span class="w"> </span><span class="p">[]</span></span></span></code></pre>
</a></td>
      <td><em><code>(list)</code></em></td>
    </tr>
    <tr id="storage-extravolumes">
      <td><a href="#storage-extravolumes"><pre class="chroma"><code><span class="line"><span class="cl"><span class="nt">storage.extraVolumes</span><span class="p">:</span><span class="w"> </span><span class="p">[]</span></span></span></code></pre>
</a></td>
      <td><em><code>(list)</code></em></td>
    </tr>
    <tr id="storage-fullnameoverride">
      <td><a href="#storage-fullnameoverride"><pre class="chroma"><code><span class="line"><span class="cl"><span class="nt">storage.fullnameOverride</span><span class="p">:</span><span class="w"> </span><span class="s2">&#34;&#34;</span></span></span></code></pre>
</a></td>
      <td><em><code>(string)</code></em><p>Overrides the full name of the storage component</p>
</td>
    </tr>
    <tr id="storage-nodeselector">
      <td><a href="#storage-nodeselector"><pre class="chroma"><code><span class="line"><span class="cl"><span class="nt">storage.nodeSelector</span><span class="p">:</span><span class="w"> </span>{}</span></span></code></pre>
</a></td>
      <td><em><code>(object)</code></em></td>
    </tr>
    <tr id="storage-podannotations">
      <td><a href="#storage-podannotations"><pre class="chroma"><code><span class="line"><span class="cl"><span class="nt">storage.podAnnotations</span><span class="p">:</span><span class="w"> </span>{}</span></span></code></pre>
</a></td>
      <td><em><code>(object)</code></em></td>
    </tr>
    <tr id="storage-poddisruptionbudget-enabled">
      <td><a href="#storage-poddisruptionbudget-enabled"><pre class="chroma"><code><span class="line"><span class="cl"><span class="nt">storage.podDisruptionBudget.enabled</span><span class="p">:</span><span class="w"> </span><span class="kc">false</span></span></span></code></pre>
</a></td>
      <td><em><code>(bool)</code></em></td>
    </tr>
    <tr id="storage-poddisruptionbudget-maxunavailable">
      <td><a href="#storage-poddisruptionbudget-maxunavailable"><pre class="chroma"><code><span class="line"><span class="cl"><span class="nt">storage.podDisruptionBudget.maxUnavailable</span><span class="p">:</span><span class="w"> </span><span class="m">0</span></span></span></code></pre>
</a></td>
      <td><em><code>(int)</code></em></td>
    </tr>
    <tr id="storage-poddisruptionbudget-minavailable">
      <td><a href="#storage-poddisruptionbudget-minavailable"><pre class="chroma"><code><span class="line"><span class="cl"><span class="nt">storage.podDisruptionBudget.minAvailable</span><span class="p">:</span><span class="w"> </span><span class="m">0</span></span></span></code></pre>
</a></td>
      <td><em><code>(int)</code></em></td>
    </tr>
    <tr id="storage-poddisruptionbudget-unhealthypodevictionpolicy">
      <td><a href="#storage-poddisruptionbudget-unhealthypodevictionpolicy"><pre class="chroma"><code><span class="line"><span class="cl"><span class="nt">storage.podDisruptionBudget.unhealthyPodEvictionPolicy</span><span class="p">:</span><span class="w"> </span><span class="kc">null</span></span></span></code></pre>
</a></td>
      <td><em><code>(string)</code></em></td>
    </tr>
    <tr id="storage-podlabels">
      <td><a href="#storage-podlabels"><pre class="chroma"><code><span class="line"><span class="cl"><span class="nt">storage.podLabels</span><span class="p">:</span><span class="w"> </span>{}</span></span></code></pre>
</a></td>
      <td><em><code>(object)</code></em></td>
    </tr>
    <tr id="storage-replicacount">
      <td><a href="#storage-replicacount"><pre class="chroma"><code><span class="line"><span class="cl"><span class="nt">storage.replicaCount</span><span class="p">:</span><span class="w"> </span><span class="m">3</span></span></span></code></pre>
</a></td>
      <td><em><code>(int)</code></em><p>Number of storage nodes in cluster mode</p>
</td>
    </tr>
    <tr id="storage-resources">
      <td><a href="#storage-resources"><pre class="chroma"><code><span class="line"><span class="cl"><span class="nt">storage.resources</span><span class="p">:</span><span class="w"> </span>{}</span></span></code></pre>
</a></td>
      <td><em><code>(object)</code></em></td>
    </tr>
    <tr id="storage-service-annotations">
      <td><a href="#storage-service-annotations"><pre class="chroma"><code><span class="line"><span class="cl"><span class="nt">storage.service.annotations</span><span class="p">:</span><span class="w"> </span>{}</span></span></code></pre>
</a></td>
      <td><em><code>(object)</code></em></td>
    </tr>
    <tr id="storage-service-labels">
      <td><a href="#storage-service-labels"><pre class="chroma"><code><span class="line"><span class="cl"><span class="nt">storage.service.labels</span><span class="p">:</span><span class="w"> </span>{}</span></span></code></pre>
</a></td>
      <td><em><code>(object)</code></em></td>
    </tr>
    <tr id="storage-service-port">
      <td><a href="#storage-service-port"><pre class="chroma"><code><span class="line"><span class="cl"><span class="nt">storage.service.port</span><span class="p">:</span><span class="w"> </span><span class="m">8490</span></span></span></code></pre>
</a></td>
      <td><em><code>(int)</code></em><p>Port for remote write and selector queries</p>
</td>
    </tr>
    <tr id="storage-service-type">
      <td><a href="#storage-service-type"><pre class="chroma"><code><span class="line"><span class="cl"><span class="nt">storage.service.type</span><span class="p">:</span><span class="w"> </span><span class="l">ClusterIP</span></span></span></code></pre>
</a></td>
      <td><em><code>(string)</code></em><p>Service type</p>
</td>
    </tr>
    <tr id="storage-tolerations">
      <td><a href="#storage-tolerations"><pre class="chroma"><code><span class="line"><span class="cl"><span class="nt">storage.tolerations</span><span class="p">:</span><span class="w"> </span><span class="p">[]</span></span></span></code></pre>
</a></td>
      <td><em><code>(list)</code></em></td>
    </tr>
    <tr id="storage-topologyspreadconstraints">
      <td><a href="#storage-topologyspreadconstraints"><pre class="chroma"><code><span class="line"><span class="cl"><span class="nt">storage.topologySpreadConstraints</span><span class="p">:</span><span class="w"> </span><span class="p">[]</span></span></span></code></pre>
</a></td>
      <td><em><code>(list)</code></em></td>
    </tr>
    <tr id="uselegacynaming">
      <td><a href="#uselegacynaming"><pre class="chroma"><code><span class="line"><span class="cl"><span class="nt">useLegacyNaming</span><span class="p">:</span><span class="w"> </span><span class="kc">false</span></span></span></code></pre>
</a></td>
      <td><em><code>(bool)</code></em><p>Use legacy release-name based resource naming instead of operator-style names</p>
</td>
    </tr>
  </tbody>
</table>

