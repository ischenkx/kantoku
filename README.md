# Kantoku
A platform for distributed task execution

# How it works
Base structure of the platform can be represented as three separate levels

## L0
### Cell Storage
Cell - some piece of data that can be addressed by a unique identifier

### Event bus
Event bus is responsible for communication between layers 
(commonly upper levels track events from lower levels and somehow process them)

## L1
### Executor
Executor takes a task as an input and returns the result of its execution.
