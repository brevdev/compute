# Lambda Labs SECURITY.md for Brev Cloud SDK

This document outlines how the Lambda Labs integration complies with Brev Cloud SDK’s security architecture using Lambda Cloud primitives such as firewall rules and virtual network isolation.

---

## 🔐 Network Security

### Default Rules

- **Inbound:** All inbound traffic is **denied by default**. Lambda instances are provisioned without open ports unless explicitly defined. (Lambda Labs by default allows SSH/ICMP; Brev overrides this to achieve a true deny-all inbound posture.)
- **Outbound:** All outbound traffic is **unrestricted by default**. Lambda Cloud imposes no egress restrictions, so instances can freely initiate outbound connections.

### Explicit Inbound Access

- Brev allows inbound access only through explicitly defined `FirewallRule` resources.
- Each `FirewallRule` maps to a Lambda Cloud firewall rule allowing a specific port and protocol from an authorized source.
- No ports are open unless explicitly configured by Brev.

### Implementation Mapping

- **Inbound Deny:** Achieved by removing default allow rules (e.g., SSH) and not attaching any ingress rules unless defined in `FirewallRule`.
- **Outbound Allow:** No changes needed; outbound access is unrestricted by Lambda by default.
- **FirewallRule Mapping:** Each Brev `FirewallRule` maps directly to a Lambda Labs firewall rule (e.g., TCP/22 from 203.0.113.0/24).

---

## 🔒 Data Protection

### Encryption at Rest

- All persistent storage is encrypted at rest using industry-standard encryption (e.g., AES-256) by default.
- Lambda Cloud does not currently support customer-managed encryption keys.

### Encryption in Transit

- All Brev SDK and control plane communication uses HTTPS with TLS 1.2+.
- Users are encouraged to use secure protocols (e.g., SSH, HTTPS) for application-level traffic to and from their instances.

### Key Management

- Provider-managed keys are used for at-rest encryption.
- No customer-managed key option is currently available.

---

## ✅ Implementation Checklist

- [x] **Inbound Deny by Default** using Lambda Labs firewall configuration
- [x] **Outbound Allow by Default** (no action required)
- [x] **Explicit Inbound Access** mapped to Brev `FirewallRule` ➝ Lambda firewall rule
- [x] **Persistent Disk Encryption** enabled by default (AES-256)
- [x] **TLS 1.2+** enforced for all Brev SDK and API communication
- [x] **Provider-Managed Keys** for encryption at rest

---

## 🛡️ Security Contact

For security issues or vulnerability disclosures related to this integration, please contact:

**Email:** [brev@nvidia.com](mailto:brev@nvidia.com)

Please report any vulnerabilities privately before disclosing them publicly.
