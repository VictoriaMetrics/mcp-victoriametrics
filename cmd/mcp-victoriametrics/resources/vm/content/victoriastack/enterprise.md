---
title: Enterprise
menu:
  docs:
    weight: 10
    identifier: victoriastack-enterprise
    parent: victoriastack
tags:
  - metrics
  - logs
  - traces
---

VictoriaStack is Apache-2.0 licensed and complete on its own: the open source version runs in production at any scale
described in these docs. The [Enterprise](https://victoriametrics.com/products/enterprise/) offering adds two things:
**technical support from the engineers who build the stack**, and a set of features for Enterprise-level installations.

Enterprise components include all the features of the community edition - start with the community edition today and 
add Enterprise later without re-architecting, re-ingesting, or rewriting anything.

## Technical support

The main value of Enterprise is support from the people who develop VictoriaMetrics, VictoriaLogs, and VictoriaTraces.
An Enterprise contract includes:

- **A strict SLA on response times**, with online support and dedicated chat rooms, so questions go straight to the
  engineering team instead of a ticket queue.
- **Priority bug fixing** and prioritized feature requests for issues that affect your installation.
- **CVE fixes in the Enterprise version within 3 business days.**
- **Architectural guidance**: regular audits of the setup and its capacity, plus verification of the security of the
  installation and of data access.

This is also the reason many teams move to Enterprise without needing any of the extra features: the software already
works for them, and what they buy is the guarantee that somebody who knows it is reachable when it doesn't.

## Long-term support releases

VictoriaStack releases often, and following the latest release is the recommended way to stay current.
[LTS lines for VictoriaMetrics](https://docs.victoriametrics.com/victoriametrics/lts-releases/) and [VictoriaLogs](https://docs.victoriametrics.com/victorialogs/lts-releases/)
are for teams that prefer rare, low-risk updates instead: every LTS line receives bugfixes and security fixes for 12 months, a new LTS line appears every
6 months, and two lines are supported at any moment. That leaves up to 6 months to migrate between lines.

Every fix that lands in an LTS release also lands in the latest release, so users of the community edition get the same
fixes by [upgrading](https://docs.victoriametrics.com/victoriastack/deployment/#upgrades) to the latest version.

## Additional features

Enterprise adds the following capabilities, most of them to VictoriaMetrics:

| Feature                                                                                                                                                                                                                                                         | What it gives you                                                                                                                                                                                                                                                                |
|-----------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|----------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|
| [Downsampling](https://docs.victoriametrics.com/victoriametrics/single-server-victoriametrics/#downsampling)                                                                                                                                                    | keeps long retention affordable by storing older data at lower resolution                                                                                                                                                                                                        |
| [Retention filters](https://docs.victoriametrics.com/victoriametrics/single-server-victoriametrics/#retention-filters)                                                                                                                                          | different retention per team, environment, or [tenant](https://docs.victoriametrics.com/victoriastack/multitenancy/#per-tenant-retention)                                                                                                                                        |
| [Anomaly detection](https://docs.victoriametrics.com/victoriastack/anomaly-detection/)                                                                                                                                                                          | [vmanomaly](https://docs.victoriametrics.com/anomaly-detection/) learns expected data behavior and writes a [unified anomaly score](https://docs.victoriametrics.com/anomaly-detection/faq/#what-is-anomaly-score) to VictoriaMetrics for contextual alerting and investigation. |
| [vmbackupmanager](https://docs.victoriametrics.com/victoriametrics/vmbackupmanager/)                                                                                                                                                                            | scheduled [backups](https://docs.victoriametrics.com/victoriastack/deployment/#backups) with their own retention, to S3, GCS, or Azure Blob Storage                                                                                                                              |
| [mTLS](https://docs.victoriametrics.com/victoriametrics/single-server-victoriametrics/#mtls-protection) and [automatic TLS certificates](https://docs.victoriametrics.com/victoriametrics/single-server-victoriametrics/#automatic-issuing-of-tls-certificates) | mutual TLS between components and clients, including [mTLS-based routing](https://docs.victoriametrics.com/victoriametrics/vmauth/#mtls-based-request-routing) and IP filters in vmauth                                                                                          |
| [Kafka](https://docs.victoriametrics.com/victoriametrics/integrations/kafka/) and [Google PubSub](https://docs.victoriametrics.com/victoriametrics/integrations/pubsub/)                                                                                        | ingestion from and to message queues                                                                                                                                                                                                                                             |

<br>

See the [full list of Enterprise features](https://docs.victoriametrics.com/victoriametrics/enterprise/#victoriametrics-enterprise-features)
in the product docs, and the [Enterprise page](https://victoriametrics.com/products/enterprise/) for the offering itself.

## Community and Enterprise are compatible

Nothing has to be planned in advance to keep the Enterprise path open:

- **The features overlap, they don't diverge.** Enterprise includes everything the community edition has, so no
  configuration, query, dashboard, or alerting rule stops working.
- **Switching is a change of image and a license key.** Enterprise executables and Docker images carry the `enterprise`
  suffix in their names and tags, and the license is passed via the `-license` or `-licenseFile` flag. The data stays
  where it is: no migration, no re-ingestion.
- **Editions mix per component.** Only the components that use Enterprise features need an Enterprise license. For example, [`vmanomaly`](https://docs.victoriametrics.com/victoriastack/anomaly-detection/) requires a license itself, while it can read from Community or Enterprise single-node or clustered VictoriaMetrics, VictoriaLogs, or VictoriaTraces and write results to VictoriaMetrics.
- **Enterprise components are free to evaluate.** They can be used in non-production setups with
  a [trial license](https://victoriametrics.com/products/enterprise/trial/); production use requires an Enterprise
  contract.

[VictoriaMetrics Cloud](https://docs.victoriametrics.com/victoriametrics-cloud/) is built on top of VictoriaMetrics
Enterprise, so it is a third option for teams that prefer not to operate the stack themselves.

## How to get it

- Request a [trial license](https://victoriametrics.com/products/enterprise/trial/) to evaluate Enterprise components in
  a non-production setup.
- Contact [sales@victoriametrics.com](mailto:sales@victoriametrics.com) to discuss support and an Enterprise contract.
