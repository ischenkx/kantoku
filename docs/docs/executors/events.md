---
title: Task Events Protocol
---

**_Task Events Protocol_** - is an executor protocol that guarantees certain events to be emitted before or after some execution steps.

**_Inherits:_** Default

**_Rules:_**
- All events are sent to a known constant set of channels. Each event is bound to one of those channels.
- The following events must be sent during the task execution:
  - `received_task` - right after the moment when the task was popped from the Inputs
  - `executed_task` - after execution but before writing the result to the Outputs
  - `sent_outputs` - right after the moment when the result was written to the outputs
  - `error` - at any moment if something goes unrecoverably wrong

```go
for task := range tasks {
  emit("received_task", task.ID())
  result, err := execute(task)
  emit("executed_task", task.ID())
  send(task.ID(), result, err)
  emit("sent_outputs", task.ID())
}
```