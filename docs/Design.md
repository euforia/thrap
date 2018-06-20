# Design

Thrap consists of the following resource categories:

- vcs
- orchestrator
- registry
- secrets

### VCS
Version Control System

### Orchestrator
Application orchestrator for deployments

### Registry
Container image registry

### Secrets
Secrets provider.  

Secrets are provided to the application by making a file containing the secrets
at runtime.  The user supplies a destination file path relative to there
current working directory where the secrets are made available. 
