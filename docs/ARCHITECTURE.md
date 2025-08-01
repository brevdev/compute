## 🏗️ Architecture: How the Cloud SDK Fits into Brev

### Stateless by Design

The `cloud` SDK is **stateless**. It acts as a thin layer over cloud provider APIs, abstracting away provider-specific quirks while exposing a consistent interface.

        ┌───────────────────────────────┐
        │         Brev Database         │
        │ - Credentials                 │
        │ - Instances                   │
        │ - Instance Types              │
        └────────────┬──────────────────┘
                     │
           (reads/writes state)
                     │
        ┌────────────▼────────────┐
        │      Brev Backend       │
        │ - Sync loops            │
        └────────────┬────────────┘
                     │
       (stateless)   │  uses credential + region
                     ▼
             ┌──────────────┐
             │  compute SDK │
             │ (stateless)  │
             └──────┬───────┘
                    │
        ┌───────────▼────────────┐
        │   Cloud Provider API   │
        │ (AWS, GCP, etc.)       │
        └────────────────────────┘
