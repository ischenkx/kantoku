---
sidebar_position: 1
---

# What's Kantoku?

Kantoku is a platform for distributed task execution.

Our main goals:
1. The system is fully customizable and acts as an interface - you can replace any of the components with your implementations
2. Tasks are easily defined in code
3. There is a default implementation that lets you run tasks out of the box
4. There are client libraries (which let you create tasks) in multiple languages

## How is it going to work?
As we simplified the components of the system to most primary ones, we reduced the number to only 4:
1. Inputs - we put tasks and their data (basically arguments) here
2. Executors - these are user-defined programs that process tasks taken from inputs
3. Outputs - executors put the results of the processed tasks here
4. Events - on all of these steps kantoku logs some information to event storage

:::note
Events are not needed for the system itself. At the moment we think that they are too useful to get rid of, but this point is open to discussion.
:::

Note that all of these components are logical - they are user defined by kantoku design. It means that you can use our default implementation of Inputs, or Kafka Queue, or make requests to your own service doing some complex logic. 
Our only requirement is that you provide the necessary functions for components communication.

:::caution
TODO: add links to component pages
:::

## Customization
Here is what you need to implement to build all components from scratch:
- Inputs, Outputs, and Events services
- Executor services (depending on your Inputs, there can be different types, multiple instances of executors of the same type, and so on)
- Custom Plugins

## Plugins
Plugins are a way to add functionality to the process of task creation. Similarly to the client library they are defined in a specific programming language. It means that plugin has to be implemented in all languages you want to use it.

Plugins are called before/after important parts of task creation.

## How it should work for a user
Note that we call ‘user’ the person who builds functionality on top of the kantoku. For example, they send tasks for recommendation prediction and then process their results.

We plan to make libraries for some programming languages (probably go, python) with similar interfaces.

To create a task:

```go
spec := kantoku.Task("test-task", []byte('Hello world'))
               .With(loggerPlugin.LogEverything())

// kan is an instance of Kantoku
result, err := kan.Spawn(spec)
```

Here we register a plugin, then create a task, add a plugin option to it, and finally spawn it.

The result contains the id of spawned task and some info about task creation.

 