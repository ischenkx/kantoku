---
title: Dependencies and groups
---

Dependencies are objects like `(ID, resolved)` where resolved is a boolean field

Groups (think dependency groups) are objects like `(ID, pending, status)` where pending is a counter of unresolved dependencies, status is one of

- `initializing`
- `waiting`
- `scheduling`
- `scheduled`

There is a mapping `dependency.ID ⇔ group.ID`, so dependency can be a member of multiple groups, and a group can have multiple members (but they are unique).

Life-cycle of a dependency:

1. `Make() → ID` of dep; it is not saved
2. `Resolve(ID)` create resolved dep with ID, if it doesn’t exist, otherwise ensure it is resolved. Decrement counter of the linked groups

It is more complicated for a group:

1. `MakeGroup(IDs…)` makes a transaction that creates a group `(ID, -1, _initializing_)`, and maps all IDs of deps to group ID
2. There is a process that initializes groups once in the determined interval: selects all groups with `initializing` status, changes it to `waiting`, and sets the counter to the count of linked dependencies.
3. There is a similar process that searches for groups with a counter equal to zero and resolves them. When a group is taken to be resolved status is updated to `scheduling`, and when it is done to `scheduled`.
