---
sidebar_position: 1
title: Default
---

**Subject:** Kantoku platform (Inputs, Outputs, Events), Executor

**Rules:**

- The only way to receive a task is to pop it from the Inputs
- A task is equivalent to a function _F: Bytes → (Bytes, Error)_, where _Error_ is a more verbose _Status_
- If task is received it must be executed (in other words the function _F_ must be evaluated with task.Data as an
  argument) and the result must be put in the Outputs with the task’s ID as the key