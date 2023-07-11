---
title: Executors
---

**_Executor_** - a basic processing unit that is responsible for task execution and saving of the outputs.

Executors are not restricted by any specific implementation, so anyone can write their own one using any tools they want
to. Despite the fact that Kantoku doesn’t require a concrete program to be an executor, executors can’t be anything at
all - they must follow some behavioural rules that are called _**Protocols**_.

:::tip
**_Protocols_** are like interfaces in Go (or any other programming language you like) but besides the signature of the
entity they describe its algorithmic behaviour.

Let’s take a car as an example. I could write the following protocol for an arbitrary automobile:

The subject must have a gas pedal.

The subject must move in space when the driver presses the gas pedal (or any equivalent).
:::
**_Subject_** - an object that’s described by a protocol

**_Rule_** - some piece of information that restricts the behaviour of a subject

:::info definition
Two rules are said to be incompatible if one of them does not allow to satisfy the other one (in other words none of
them contradicts the other one).
:::

**_Protocol_** - a set of rules for a given subject such that none of them contradicts any other.

:::info definition
Two protocols are said to be compatible if their union is a valid protocol.
:::

:::info definition
A protocol A inherits a protocol B if B is a subset of A
:::

---

Kantoku only works with executors that inherit the [Default](./default) protocol