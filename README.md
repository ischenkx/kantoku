# Kantoku

A platform for distributed task execution

# Roadmap

### Refactoring
- [x] Single storage for tasks
  - [x] We need to get rid of Task.Properties and use Info as a primary source of knowledge about tasks
  - [ ] It might be better to have immutable fields in record.Record
- [x] Consistency
  - [x] Acknowledgement of successful processing in queues (events)
  - [ ] Eventual consistency of multi-service transactions (using compensating transactions aka SAGAs)
    - This requires a roll-back action for all mutations (so, we probably need to add deletion of any entities in the system)
- [ ] Logging
### Features
- [x] Task Dependencies
  - A dependency based scheduler
- [ ] Logging / Metrics (Prometheus + Grafana)
- [ ] Services
  - Common service structure for different types of processors
    - We need a single package that would take care of running and gracefully shutting down
      all types of processors
  - Cross-session service identification
    - Each service instance must have its persistent ID that can be used as a consumer group label for queues.
  - Service discovery
    - We need to collect information about services (their presence, id, other properties)
- [ ] Functional Tasks
  - It should be possible to describe tasks as Go/(Other language) functions
- [ ] Context
  - Tasks should be grouped by contexts
- [ ] Pipelines
  - A convenient way to compose tasks
- [ ] Test Coverage
- [ ] Deployment in k8s

### Links

[Jira](https://r-ischenko.atlassian.net/jira/software/projects/KAN/boards/1)

[Miro](https://miro.com/app/board/uXjVNS1PReQ=/)
