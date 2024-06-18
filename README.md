# Kantoku

A platform for distributed task execution

# Todo

## Back-end

### Infra
- [ ] Save logs to some persistent storage (OpenSearch/Elastic/Loki)
- [ ] Use Prometheus and Grafana for metrics
- [ ] When environment is not set up (Docker isn't running), services start to poll ddos the components (Consul, Postgres, etc). We need timeouts.
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
- [ ] Task Restarts
	A restart seems to be an easy concept but in the world of kantoku a few problems arise:
	- Should task be recreated (or sending a new `task_ready` event is enough)? Should output resources be recreated (if yes, how to manage dependents of the failed task?)? 
	- Restarting a task parallelly can cause double output resource resolution. So, it's crucial to guarantee some kind of synchronization of the action. It can be done by delegating the restarting logic to scheduler (actually it's not a complete solution yet, synchronization among multiple scheduler instances must be also provided) or by setting some flag in the restarted task's info, e.g. `restarted: true` (this approach also has a tricky part: the task storage should guarantee atomic writes and such stuff).
	- Restart context id and restart parent id: restarts represent a simple tree. 	  
- [ ] Implement auto tasks parsing from source code
	- [ ] Use "source" column in the specifications table to identify which tasks are expected to be deleted
	- [ ] Either parse task from source code structures with embedded `kantoku.Task`  or register tasks manually in some package-wise router (in this case, code generation for task files is preferrable)
- [ ] Investigate alternatives for `records`
	- [ ] Consider using raw mql requests (can be used with Postgres via FerretDB)
	- [ ] Consider using postgres + jsonb
- [ ] Scheduler
	- [ ] Dependencies
		- [ ] Make resource dependency resolver instant using resource initialization notifier.
		- [ ] What other resolvers can be useful?
		- [ ] Should dependencies fail or have any other statuses besides `PENDING` and `READY`?
		- [ ] Investigate approaches for a better dependencies engine 
	- [ ] Make it scalable
- [ ] Executor
	Executor interface is pretty simple and the current implementation is sufficient for basic usage. In the future (meaning that this task shouldn't be taken until more fundamental problems are solved) it'd be nice to have a more complex and environment-aware version. Some points on what'd be cool to see:
	- Select the running machine based on some set of labels in task's info
	- Shard tasks by their execution context / resource usage. This would improve usage of a local resource cache  
	- Implement Leader/Follower decomposition.
- [ ] API
	- [ ] Make it follow REST principles (stop using `POST` for all requests)
	- [ ] Consider using jsonrpc
- [ ] Spawners
    A spawner is a new entity. It is some program that spawns tasks. Might be a kafka message, a cron event or a telegram message.
- [ ] Pipelines
	Pipelines are DAGs of tasks. It's pretty simple to implement but it requires good visualization along with support for partial restarts (a task restart without output reallocation).
	- Preprocessors: pure functions written in some scripting language that are used for a slight data rectification in order to match a dependent task input's type.   
### Tests
- [ ] Add tests for all services
- [ ] Add tests for kantoku (end2end)
- [ ] Write unit tests for common stuff

## Front-end
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
