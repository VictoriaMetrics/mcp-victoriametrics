---
draft: false
page: blog blog_post
authors:
  - Jose Gomez-Selles
date: 2026-01-13
title: "What's New in VictoriaMetrics Cloud Q4 2025? New tiers, more deployment options, IaC and alerting rules."
enableComments: true
featuredPost: true
summary: "In the last quarter of 2025, VictoriaMetrics Cloud brings many great features: New powerful Capacity Tiers, the expansion to the us-east-1 (N.Virginia) AWS region in the US, new Notification Groups, a Terraform provider to complete your IaC, 9 brand new Alerting Rule Integrations and much more."
categories:
  - Product News
tags:
  - victoriametrics
  - cloud
  - monitoring
  - observability
  - aws
  - roadmap
images:
  - /blog/q4-2025-whats-new-victoriametrics-cloud/preview.webp
---

2025 has been quite a year! As we enter into 2026, we want to share all the great features that
[VictoriaMetrics Cloud](https://console.victoriametrics.cloud/) has brought in the last quarter.

Remember that this **Quarterly Live Update** is available in video format as well here:

{{< youtube RzB4go1awDg >}}

Let's get to it!

## New VictoriaMetrics Single Tiers

One of the main changes of this last quarter of 2025 has been the introduction
of **new _Capacity Tiers_ for VictoriaMetrics Single** type instances in
VictoriaMetrics Cloud. In a nutshell, after a benchmarking exercise, we've adjusted our tiers so
they follow a linear progression in computing power, quotas/limits and costs for our users.
We believe this work gives more flexibility while easing the price model
progression as users scale their Observability needs.

<img src="/blog/q4-2025-whats-new-victoriametrics-cloud/cus.webp"
     style="width:80%; display:block; margin:auto;" />

Apart from adding more power and adjusting prices, the UI for instance creation
has been simplified and the benchmark results have been [published](https://docs.victoriametrics.com/victoriametrics-cloud/deployments/tiers-and-types/).

For Cluster type deployments (which provide higher SLAs and multitenancy), the main change
is that interested users can now open a request for a PoC in which our team will tailor an
instance for their requirements, adapting to their needs.

Learn more about the new tiers in the full announcement [here](https://victoriametrics.com/blog/new-tiers-in-victoriametrics-cloud-2025/).

## New US region available: us-east-1 (N.Virginia)

VictoriaMetrics Cloud keeps expanding thanks to the fact that more and more users are relying on
our platform. It was not long ago that we opened a [new region in Singapore](https://victoriametrics.com/blog/cloud-new-region-ap-southeast-1/),
and now we are happily announcing that we are also **open for business in the AWS us-east-1
(N.Virginia) region**.

This region is special because it is one of the most used worldwide, and we definitely didn't want
to miss the party (with all the fun it entails)! We believe this is great news for those who
are already in this region and want to benefit from AWS PrivateLink and lower latency for their
observability data.

<img src="/blog/q4-2025-whats-new-victoriametrics-cloud/new-region.webp"
     style="width:70%; display:block; margin:auto;" />

Feel free to give it a go whether you already are a VictoriaMetrics Cloud user, or are just
considering VictoriaMetrics Cloud. In this case, let me kindly remind you that a simple [Sign Up](https://console.victoriametrics.cloud/signUp)
gives you a free trial for a month (without requiring you to add a payment method).

As usual, [contact us](https://victoriametrics.com/contact-us/) if you'd like to see us in other
regions!

## Notification groups

We've been after improving our notifications system for quite a while now. That's why we are very
excited to rollout an improved system that helps your team stay informed about critical events,
product updates, and account activities, while giving the ability to customize them within
different groups. In this way, different users can define how (and what!) to get notified.
Notifications are divided into the following interest groups:
- **System alerts**: for critical updates about platform health and infrastructure
- **Billing**: is related to notifications about invoices, usage limits, and payment issues
- **Product News**: updates on new features, improvements, and releases
- **Activity**: events related to user actions and deployments

<img src="/blog/q4-2025-whats-new-victoriametrics-cloud/notifications_crop.webp"
     style="width:100%; display:block; margin:auto;" />

These categories can be used to define which updates and notifications are relevant to you,
and set to which email addresses and/or slack channels you want to get your system alerts,
at Organization or User level.

## Terraform integration: more ways to automate your VictoriaMetrics Cloud instances

We do love automating our own actions and that's why we've always wanted to help (and encourage)
VictoriaMetrics Cloud users to do the same. Historically, the main way to do that has been via
the API, (see our [docs](https://docs.victoriametrics.com/victoriametrics-cloud/api/))
and it's been great to see how everyone has been using it. In case you haven't tried it, check everything
you can do, the Swagger way, [here](https://console.victoriametrics.cloud/api-docs).

Not long ago we also made public the **client library for the VictoriaMetrics Cloud API in go**,
which provides an idiomatic way to _talk_ with the Product. Check how to integrate in the
[Github repo](https://github.com/VictoriaMetrics/victoriametrics-cloud-api-go).

So what was left? Exactly, being part of the Terraform ecosystem!

<img src="/blog/q4-2025-whats-new-victoriametrics-cloud/terraform.webp"
     style="width:80%; display:block; margin:auto;" />

### Registry

Since last December, our users can deploy and manage their instances via the VictoriaMetrics Cloud
provider, available on the Terraform [registry](https://registry.terraform.io/providers/VictoriaMetrics/victoriametricscloud/latest).
With the implementation made Open Source, and available in [this repository](https://github.com/VictoriaMetrics/terraform-provider-victoriametricscloud),
this provider gives a wrapper of the mentioned VictoriaMetrics Cloud API that exposes all
resources and data sources (deployments, access tokens, rules, regions...) so you can manage
VictoriaMetrics Cloud as part of your **Infrastructure as Code**.

After saving your API key, this is all it takes for you to add the VictoriaMetrics Cloud provider
to your configuration:

```terraform
terraform {
  required_providers {
    victoriametricscloud = {
      source  = "VictoriaMetrics/victoriametricscloud"
      version = "0.0.1"
    }
  }
}

provider "victoriametricscloud" {
  api_key = var.api_key
}
```

### Terraform Integration

But that's not all. As usual, we wanted to guide users to easily integrate VictoriaMetrics Cloud
with the terraform provider. Check out the [integration page](https://console.victoriametrics.cloud/integrations/terraform?integrationType=other)
inside the product, or in the [docs](https://docs.victoriametrics.com/victoriametrics-cloud/integrations/terraform/)
with the same embedded experience.

## Integrations: New alerting rules

Following with integrations, our last main addition this quarter is the inclusion of **9 new
Alerting Rules integrations**. As with previous integrations of this kind, we are making it
easier for those who want to get notified when something is wrong with their stack by importing
more of the popular [Awesome Prometheus Alerts](https://samber.github.io/awesome-prometheus-alerts/)
into the [integrations](https://console.victoriametrics.cloud/integrations?integrationType=rules) section.

<img src="/blog/q4-2025-whats-new-victoriametrics-cloud/integrations.webp"
     style="width:100%; display:block; margin:auto;" />

This time, users can benefit from getting out of the box alerting rules for:

- **CI/CD:** ArgoCD, FluxCD
- **Observability agents:** OpenTelemetry collector
- **Networking/Mesh:** CoreDNS, Istio, Consul
- **Databases/storage:** etcd, PostgreSQL, ClickHouse

As you can see, by picking from a list of alerts it's much easier to plan for those situations that
may get tricky. Of course, we recommend to only add those alerts that
are important to you and actually actionable. This will avoid unnecessary fatigue, letting you
focus on what matters, while being certain that you'll get notified for important events.

<img src="/blog/q4-2025-whats-new-victoriametrics-cloud/otel_rules.webp"
     style="width:100%; display:block; margin:auto;" />


## Other Improvements

We don't want to take more of your precious time, but this is not all. Among all these great
features and additions we just commented, we have also been working on improving our **resiliency**,
making upgrades across the platform, revising our Disaster Recovery procedures and adding more Scalability in
some of our components.

There's a new metric available for users in the Overview and Monitoring pages of every deployment:
**Data read rate**. It shows the amount of data scanned per time unit, and represents the
reading effort that the deployment is doing. While mostly influenced by the querying path, it also
accounts for VictoriaMetrics periodic merges and deduplications. It took us some effort and
discussions but we believe users will finally have a way to understand how their reading load impacts
the system performance. This information will help to decide when upgrading or downgrading between tiers.

Other minor improvements include: **Access Tokens** enhancements in limits, the inclusion of the last usage
time in the API and bug fixes. Organization Id is now available for **hyperlinks**, and new users can see
now that higher **quotas are automatically enabled** upon subscription. Last but not least, we've also improved
the user experience in Stripe card expiration **notifications**.

> [!TIP] Did you notice the ~Easter~ Christmas Egg? :) Let us know in the comments!

## What’s Next?

We're so excited to announce that VictoriaLogs is ready, as we've finalized the work on backups and tier definition.
This means that **in the following weeks** you can expect VictoriaLogs being part of the
VictoriaMetrics Cloud product. This wouldn't have been possible without the feedback and help
from our early adopters, so we want to publicly **thank you** for all your patience and meaningful
reflections. Thanks to you all!

2026 is going to be a great year. See you soon!

And if you haven't done it already, remember that there's a free trial when you
[Sign up for VictoriaMetrics Cloud](https://console.victoriametrics.cloud/signUp)!
