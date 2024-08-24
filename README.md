# <img src="./dev/assets/logo.svg" width="60"> Kantoku

A platform for distributed task execution

<img src="./dev/assets/pipeline_example.gif" width="600">

# Todo

## Back-end

### Core
- [x] Basic
    - [x] Task Dependencies
    - [x] Services
        - Common service structure for different types of processors
            - We need a single package that would take care of running and gracefully shutting down  
              all types of processors
        - Cross-session service identification
            - Each service instance must have its persistent ID that can be used as a consumer group label for queues.
        - Service discovery
            - We need to collect information about services (their presence, id, other properties)

- [x] Task Restarts

  A restart seems to be an easy concept but in the world of kantoku a few problems arise:
    - Should task be recreated (or sending a new `task_ready` event is enough)? Should output resources be recreated (if
      yes, how to manage dependents of the failed task?)?
    - Restarting a task parallelly can cause double output resource resolution. So, it's crucial to guarantee some kind
      of synchronization of the action. It can be done by delegating the restarting logic to scheduler (actually it's
      not a complete solution yet, synchronization among multiple scheduler instances must be also provided) or by
      setting some flag in the restarted task's info, e.g. `restarted: true` (this approach also has a tricky part: the
      task storage should guarantee atomic writes and such stuff).
    - Restart context id and restart parent id: restarts represent a simple tree.

- [x] Investigate alternatives for `records`
    - [x] Consider using raw mql requests (can be used with Postgres via FerretDB)
    - [ ] ~~Consider using postgres + jsonb~~

- [ ] Implement auto tasks parsing from source code
    - [ ] Use "source" column in the specifications table to identify which tasks are expected to be deleted
    - [ ] Either parse task from source code structures with embedded `kantoku.Task`
      or register tasks manually in some package-wise router (in this case, code generation for task files is
      preferrable)

- [ ] Pipelines

  Pipelines are DAGs of tasks. It's pretty simple to implement but it requires good visualization along with support for
  partial restarts (a task restart without output reallocation).
    - Preprocessors: pure functions written in some scripting language that are used for a slight data rectification in
      order to match a dependent task input's type.

- [ ] Processor Nodes
  Nodes that would read the events and process them with multiple handlers (status update/garbage collector/etc)

- [ ] WebSocket API
    - Use Centrifugo for notifications about updates in a context/task/namespace

- [ ] CLI (ktk)
    - Design a `ktk` cli utility similar to kubectl

- [ ] API
    - [ ] Make it follow REST principles (stop using `POST` for all requests)
    - [ ] Consider using jsonrpc

- [ ] Spawners

  A spawner is a new entity. It is a program that spawns new tasks when some event happens. Might be a kafka message, a cron event or a telegram
  message.

- [ ] Batching

  APIs and services should support batched execution. It would reduce amount of network requests to databases and
  overhead of small requests

- [ ] Garbage collector

  Garbage collection is necessary to limit the amount of used resources

- [ ] Scheduler
    - [ ] Dependencies
        - [ ] Make resource dependency resolver instant using resource initialization notifier.
        - [ ] What other resolvers can be useful?
        - [ ] Should dependencies fail or have any other statuses besides `PENDING` and `READY`?
        - [ ] Investigate approaches for a better dependencies engine
    - [ ] Make it scalable

- [ ] Executor

  Executor interface is pretty simple and the current implementation is sufficient for basic usage. In the future (
  meaning that this task shouldn't be taken until more fundamental problems are solved) it'd be nice to have a more
  complex and environment-aware version. Some points on what'd be cool to see:
    - Select the running machine based on some set of labels in task's info
    - Shard tasks by their execution context / resource usage. This would improve usage of a local resource cache
    - Implement Leader/Follower decomposition.

### Infra

- [ ] Save logs to some persistent storage (OpenSearch/Elastic/Loki)
- [ ] Use Prometheus and Grafana for metrics
- [ ] When environment is not set up (Docker isn't running), services start to ddos the components (Consul, Postgres,
  etc). We need timeouts.

### Tests

- [ ] Add tests for all services
- [ ] Add tests for kantoku (end2end)
- [ ] Write unit tests for common stuff

### Docs

- [ ] Create a JS emulator of the whole system and show it using reactflow

## Web UI

### Task List

- [x] Fix `failed` and `ok` statuses
- [x] Add `running` status (a.k.a received)
- [x] Make ID a hyperlink to task (and remove the eye action)
- [x] Make one-sided filter on updated_at
- [ ] Add updated_at column and make it sortable

### Task Show

- [x] Show the task's status (ok, failed, running...)
- [x] Show type (aka spec) and info about it
- [x] List dependencies
- [x] Show `updated_at`
- [x] List parameters (aka inputs) and outputs as a form
- [x] Show all info as json editor, make its style fit the antd's theme
- [ ] Show specifications in "Specifications" tab
- [ ] Show graphs on frontend
    - [ ] Use context_id and parent_id to connect tasks

### Links

[Jira](https://r-ischenko.atlassian.net/jira/software/projects/KAN/boards/1)

[Miro](https://miro.com/app/board/uXjVNS1PReQ=/)
