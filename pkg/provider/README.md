# provider

This package contains 1 directory per provider type. The following is the current
support matrix:

- [x] registry
- [ ] orchestrator
- [ ] secrets
- [ ] vcs

## Type

A type defines a resource type which then contain providers.

### Requirements

- Should contain a list of supported providers
- Must be self contained without depending on internal libs