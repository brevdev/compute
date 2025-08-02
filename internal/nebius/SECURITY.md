# Nebius SECURITY.md for Brev Cloud SDK

This document explains how Nebius VMs meet Brev Cloud SDK’s security requirements using Nebius primitives like Security Groups, VPCs, and projects.

---

## Network Security

### Default Rules

* **Inbound:** All inbound traffic is **denied by default** using a custom Nebius Security Group with no inbound rules.
* **Outbound:** We explicitly **allow all outbound traffic** by adding a wide egress rule (all ports/protocols to `0.0.0.0/0`).

### Explicit Access

* All inbound access must be added manually via Brev’s `FirewallRule` interface.
* These are mapped to Nebius Security Group rules that allow specific ports and sources.

### Isolation

* Each cluster uses its own Security Group.

---

## Cluster Security

* Instances in the same cluster:

  * Share a Security Group.
  * Can talk to each other using a "self" rule (Nebius allows rules that permit traffic from the same group).
* No traffic is allowed from outside the cluster unless explicitly opened.
* Different clusters use different Security Groups to ensure isolation.

---

## Data Protection

### At Rest

* Nebius encrypts all persistent disks by default using AES-256 or equivalent.

### In Transit

* All Brev SDK API calls use HTTPS (TLS 1.2+).
* Internal instance traffic should use secure protocols (e.g., SSH, HTTPS).

---

## Implementation Checklist

* [x] Default deny-all inbound using custom Nebius Security Group
* [x] Allow-all outbound via security group egress rule
* [x] `FirewallRule` maps to explicit Nebius SG ingress rule
* [x] Instances in the same cluster can talk via shared SG "self" rule
* [x] Different clusters are isolated using separate SGs or VPCs
* [x] Disk encryption enabled by default
* [x] TLS used for all API and external communication

---

## Security Contact

* Email: [brev@nvidia.com](mailto:brev@nvidia.com)
* Please report vulnerabilities privately before disclosing publicly.
