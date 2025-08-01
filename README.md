# Brev Compute SDK (v1)

An early-stage, vendor-agnostic Go SDK for managing **clusterable, GPU-accelerated compute** across cloud providers.

> This is the internal interface we use at Brev — now open-sourced to align with NVIDIA Cloud Partners (NCPs) and collaborators building GPU infrastructure.

---

## Project Goals

- Define a clean, minimal interface for compute primitives:
  - `Instance`
  - `Storage`
  - `FirewallRule`
  - `InstanceType`
  - `Location`

- Enable **clusterable GPU workloads** across multiple providers, with shared semantics and L3 network guarantees.
- Provide a foundation for Brev's `launchables`, `instances`, and provisioning flows.

---

## Security

All cloud integrations must follow our [Security Requirements](SECURITY.md), which define:

- **Network Security**: Default "deny all inbound, allow all outbound" model
- **Cluster Security**: Internal instance communication with external isolation
- **Data Protection**: Encryption requirements for data at rest and in transit
- **Implementation Guidelines**: Security checklists for cloud provider integrations

See [SECURITY.md](docs/SECURITY.md) for complete security specifications and implementation requirements.

---

## Scope and Philosophy

- **Internal-first interface**: This repo reflects the interfaces we use within Brev — open-sourced for shared development.
- **No provider integrations (yet)**: Initial release does not include AWS, GCP, or other cloud backends. These will follow.
- **Accelerator-focused**: Designed for GPUs, but extensible to CPU or container-based workloads.

---

## Status

- Version: `v1` — internal interface, open-sourced
- Current scope: core types + interfaces + tests
- Cloud provider implementations are internal-only for now
- `v2` will be shaped by feedback and contributions from the community

---

## Who This Is For

- **NVIDIA Cloud Partners (NCPs)** looking to offer Brev-compatible GPU compute
- **Infra teams** building cluster-aware systems or abstractions on raw compute
- **Cloud providers** interested in contributing to a shared interface for accelerated compute

---

---

## Documentation

- **[V1 Design Notes](docs/V1_DESIGN_NOTES.md)**: Design decisions, known quirks, and AWS-inspired patterns in the v1 API
- **[Architecture Overview](docs/ARCHITECTURE.md)**: How the Compute SDK fits into Brev's overall architecture
- **[Security Requirements](docs/SECURITY.md)**: Security specifications and implementation requirements

---

## Get Involved

This is a foundation — we're opening it early to **learn with the community** and shape a clean, composable `v2`. If you're building GPU compute infrastructure or tooling, we'd love your input.

