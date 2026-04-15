---
weight: 10
title: VictoriaMetrics Cloud Billing
menu:
  docs:
    parent: "cloud"
    weight: 10
    name: Billing
tags:
  - metrics
  - cloud
  - enterprise
---

## Pricing model

VictoriaMetrics Cloud pricing is based on a fixed tier model, where majority of the costs are known
at deployment time. The cost per deployment consists of:
- **Compute**: The cost of deployment installation. Users select and deploy a [Capacity Tier](https://docs.victoriametrics.com/victoriametrics-cloud/deployments/#capacity-tiers), with a fixed cost.
- **Storage**: The storage used by the deployment. Users select their desired storage, with a fixed cost.
- **Network**: External (egress) network usage. Users are charged $0.09 per GB for egress traffic,
which matches AWS’ rate. Even for very large deployments, from our experience this amount is
typically below 0.5% of the total cost. The total cost will be included in your monthly invoice.

This breakdown is made to help understanding and managing your costs. Usage data is sent hourly to
the payment provider (AWS or Stripe). Detailed billing information is available via the [Billing Page](https://console.victoriametrics.cloud/billing) of your VictoriaMetrics Cloud account.

### Why is this important?
Each deployment operates with predefined configurations and limits, **protecting you from unexpected
overages** caused by factors such as:

* Data ingestion spikes.
* Cardinality explosions.
* Accidental heavy queries.

> [!TIP]
> This ensures predictable costs and proactive alerts for workload anomalies.

### Detailed pricing structure

Pricing begins at ~**$190/month** for the smallest [tiers](https://docs.victoriametrics.com/victoriametrics-cloud/deployments/#capacity-tiers) of VictoriaMetrics and VictoriaLogs.
To view other tiers and their costs, navigate to the [Create New Deployment](https://console.victoriametrics.cloud/deployments/create) section in the VictoriaMetrics Cloud application.

Our aim is to make pricing information easy to access and understand. If you have any questions or feedback on our pricing, please contact us.

## Usage Reports

The [Usage Reports](https://console.victoriametrics.cloud/billing/usage) section in the billing area provides a breakdown of:

* Storage Costs
* Compute Costs
* Networking Costs
* Applied Credits

Your Final Monthly Cost is calculated as `usage - credits` and reflects the amount billed by your payment provider.

A graph is also available to display the daily cost breakdown for the selected month.


## Payment Methods

VictoriaMetrics Cloud supports the following payment options:

- Credit Card
- AWS Marketplace
- ACH Transfers

You can add multiple payment methods and set one as the primary. Backup payment methods are used if the primary fails. More details are available via the [Payment Methods](https://console.victoriametrics.cloud/billing) tab of the Billing Page.

__Note__: VictoriaMetrics Cloud does not store or process your payment information. We rely on trusted API providers (Stripe, AWS) for secure payment processing.

### Credit Card

Credit cards can be added through [Stripe](https://stripe.com/) integration.

### AWS Marketplace

Payments made via [AWS Marketplace](https://aws.amazon.com/marketplace/pp/prodview-atfvt3b73m2z4?sr=0-1&ref_=beagle&applicationId=AWSMPContessa) include billing details in the AWS portal. AWS finalizes monthly bills at the start of the next month, typically charging between the 3rd and 5th business day. Visit the [AWS Knowledge Center](https://aws.amazon.com/premiumsupport/knowledge-center/) for more information.

### ACH Transfers

ACH payments are supported. Contact [VictoriaMetrics Cloud Support](https://docs.victoriametrics.com/victoriametrics-cloud/support/) for setup assistance.



## Invoices

[Invoices](https://console.victoriametrics.cloud/billing/invoices) are emailed monthly to users who pay via Credit Card or ACH Transfers. Notification email addresses can be updated in the [VictoriaMetrics Cloud Notifications](https://docs.victoriametrics.com/victoriametrics-cloud/setup-notifications/) section.

Invoices are also accessible on the Invoices Page, which provides:

* Invoice Period
* Invoice Status
* Downloadable PDF Links

For AWS Marketplace billing, check the AWS Portal for invoice information.

---

## FAQ

Check for pricing or billing Frequently Asked Questions at the [Pricing and Billing part of the VictoriaMetrics Cloud FAQ](https://docs.victoriametrics.com/victoriametrics-cloud/cloud-faq/#pricing--billing)


