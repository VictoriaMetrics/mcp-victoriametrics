

---
title: VictoriaMetrics Common
menu:
  docs:
    parent: helm
    weight: 15
    identifier: helm-victoria-metrics-common
url: /helm/victoria-metrics-common/
aliases:
  - /helm/victoriametrics-common/
tags:
  - metrics
  - kubernetes
  - logs
---

![Version](https://img.shields.io/badge/0.0.46-gray?logo=Helm&labelColor=gray&link=https%3A%2F%2Fdocs.victoriametrics.com%2Fhelm%2Fvictoria-metrics-common%2Fchangelog%2F%230046)
![ArtifactHub](https://img.shields.io/badge/ArtifactHub-informational?logoColor=white&color=417598&logo=artifacthub&link=https%3A%2F%2Fartifacthub.io%2Fpackages%2Fhelm%2Fvictoriametrics%2Fvictoria-metrics-common)
![License](https://img.shields.io/github/license/VictoriaMetrics/helm-charts?labelColor=green&label=&link=https%3A%2F%2Fgithub.com%2FVictoriaMetrics%2Fhelm-charts%2Fblob%2Fmaster%2FLICENSE)
![Slack](https://img.shields.io/badge/Join-4A154B?logo=slack&link=https%3A%2F%2Fslack.victoriametrics.com)
![X](https://img.shields.io/twitter/follow/VictoriaMetrics?style=flat&label=Follow&color=black&logo=x&labelColor=black&link=https%3A%2F%2Fx.com%2FVictoriaMetrics)
![Reddit](https://img.shields.io/reddit/subreddit-subscribers/VictoriaMetrics?style=flat&label=Join&labelColor=red&logoColor=white&logo=reddit&link=https%3A%2F%2Fwww.reddit.com%2Fr%2FVictoriaMetrics)

VictoriaMetrics Common - contains shared templates for all Victoria Metrics helm charts

## Supported templates

| Name                   | Description                                                                                                        |
|------------------------|--------------------------------------------------------------------------------------------------------------------|
| `vm.license.volume`    | renders volume for license secret if either `.Values.global.license` or `.Values.license` is set                   |
| `vm.license.mount`     | renders volume mount for license secret if either `.Values.global.license` or `.Values.license` is set             |
| `vm.license.flag`      | renders vm command line flags if either `.Values.global.license` or `.Values.license` is set                       |
| `vm.image`             | renders container image depending on `.Values.global.image`, `.Values.image`, `.Values.license` and chart params   |
