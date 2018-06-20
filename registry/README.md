# registry
This package contains container registries such as ecr, docker hub etc.

### Pulling a manifest

```
curl -H "Authorization: Bearer ${token}" \
    https://registry.hub.docker.com/v2/cockroachdb/cockroach/manifests/latest
```
