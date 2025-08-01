# V1 Design Notes

This README captures the thinking behind the first version of the `compute` interface and some of the decisions that shaped it. It's meant to provide context as we move toward more general, multi-cloud support.

---

## 🌍 Background

The first version of this interface was designed around **AWS EC2**. At the time, AWS was the only provider we were supporting, so we built around its APIs and assumptions. The design has evolved, however you will notice AWS's momentum in the API.

---

## 🧱 Known Quirks

### Location vs SubLocation vs AZ
- The difference between `Location`, `SubLocation`, and `AvailableAzs` is unclear.
- Some providers don’t expose availability zones or don’t map cleanly to this model.

### Universal Tagging Assumption
- Tag updates are assumed to work everywhere.
- Many providers lack first-class tag support.

### Lifecycle Status Mapping
- LifecycleStatus enums are based on AWS terms.

### InstanceTypeID Generation
- Relies on manual string formatting (e.g., `location-subLoc-type`).

### Error Modeling
- Uses a few top-level errors (e.g., `ErrOutOfQuota`) with no structured data.
- Makes it hard to reason about retryability or provider-specific failure modes.

### Inconsistent Use of “Disk”, “Volume”, and “Storage”
The terminology around instance-attached storage is one of the more confusing parts of the v1 design. The interface uses three overlapping terms:
- Disk: Used in Instance and CreateInstanceAttrs (AdditionalDisks, DiskSize)
- Volume: Used in other methods
- Storage: Used in capabilities (SupportedStorage), and types (StorageFilters)

#### Unclear Ownership
- DiskSize appears both in CreateInstanceAttrs and UpdateInstanceRetireableArgs, but it’s unclear if it applies only to the root volume or all attached volumes.
- AdditionalDisks allows specifying multiple disks, but there’s no clear linkage to volume IDs or mount behavior post-creation.

#### Provider-Specific Mismatch
- Some clouds (e.g. AWS) treat root and attached volumes differently (with separate APIs).
- Others (e.g. Lambda) don’t expose volumes at all — only a total storage value.
- Elastic volumes, ephemeral disks, and NVMe local storage are not modeled cleanly in v1.
