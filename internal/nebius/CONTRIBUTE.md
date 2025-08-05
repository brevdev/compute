# Contributing to Nebius Brev Compute SDK

Nebius has a [golang SDK](https://github.com/nebius/gosdk) that is used to interact with the Nebius API.

Get started by reading the [Nebius API documentation](https://github.com/nebius/api).

## Local Development

Place a credential file in your home directory and run the provider tests.


## Prompts

```
Please analyze the Nebius @nebius  implementation in and research their actual API capabilities to determine which features are supported vs unsupported.

**Required Research Steps:**
1. **API Documentation Analysis**: Find and review Nebius' official API documentation at @https://api.nebius.com/docs  to identify supported endpoints and features. Use the golang sdk https://github.com/nebius/gosdk as a reference.
2. **Feature-by-Feature Verification**: For each CloudClient interface method, verify if Nebius actually supports it by checking their API docs
3. **Evidence-Based Decisions**: Only mark features as "supported" if you find concrete evidence in their documentation

**Implementation Approach:**
- Use the existing NotImplCloudClient pattern for unsupported features
- Return ErrNotImplemented for features that Lambda Labs doesn't support
- Maintain full CloudClient interface compliance

**Key Questions to Answer:**
- What instance management operations does Lambda Labs actually support? (create, terminate, list, stop/start, reboot?)
- Do they support volume resizing, instance type changing, or machine images?
- What networking features do they provide? (firewall rules, security groups?)
- Do they have quota management APIs?
- What authentication and location management do they support?

**Deliverables:**
1. Evidence-based list of supported vs unsupported features
2. Refactored implementation using NotImplCloudClient
3. Updated documentation reflecting actual capabilities
4. Remove any files/methods that are completely unnecessary

Please provide specific citations from Lambda Labs' API documentation for any features you mark as supported. Make sure we list features that the client implements NOT the API.
```
