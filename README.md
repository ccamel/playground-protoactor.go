# playground-protoactor.go

![Build](https://github.com/ccamel/playground-protoactor.go/workflows/Build/badge.svg)
[![gitmoji](https://img.shields.io/badge/gitmoji-%20😜%20😍-FFDD67.svg?style=flat-square)](https://gitmoji.carloscuesta.me)

> My playground I use for playing with fancy and exciting technologies. This one's for experimenting a platform actor in
> go named [protoactor](https://github.com/AsynkronIT/protoactor-go), in which actors follow DDD - Event Sourcing / CQRS principles.

## Introduction

I have a deep-rooted passion for the Actor Model. Having been a big fan and an avid admirer of the [Erlang](https://www.erlang.org) Actor system, its elegance, and efficiency in managing concurrent processes, I've always been intrigued by the power and simplicity it offers.

[Erlang](https://www.erlang.org)'s approach to concurrency, fault tolerance, and system distribution has set a high standard in the realm of Actor-based systems. Drawing inspiration from this, I've created this repository as an experimentation to explore these compelling concepts in a different landscape - the world of Go. It's a personal quest to understand how the robustness and versatility of the Actor Model can be harnessed in Go.

If you're new to the Actor Model, take a moment to read about its theory. This project isn't just about understanding the theory; it's about diving into how to use it effectively and why it's a valuable tool in building concurrent systems.

## Why the Actor Model (in Go)?

Go, with its native support for concurrency, seems like the perfect ground for experimenting with the Actor Model. The Actor Model offers an excellent abstraction, allowing us to build complex, concurrent systems in a more manageable and less error-prone way. It's about encapsulating state and behavior in actors, promoting message-driven interactions, and ensuring that our systems are scalable and maintainable.

In this repository, I'll be using the [protoactor](https://proto.actor/) library, a Go implementation of the Actor Model which provides a platform for building applications using the Actor Model. It's key features include:

- Minimalistic API - small and easy to understand and use.
- Build on existing technologies.
- Protobuf all the way for maximum performance and interoperability.
- Pluggable serialization.

## Objectives

This project aims to:

- Explore [protoactor](https://github.com/AsynkronIT/protoactor-go), a framework that brings the Actor Model to Go.
- Bring some of the concepts from the [Erlang](https://www.erlang.org) Actor Model to Go, replicating some of its elegance.
- Integrate Domain-Driven Design (DDD) principles, particularly focusing on Event Sourcing in the Actor-based systems.
- Explore persistence options for actors, including in-memory and disk-based storage, possibly using existing database systems (e.g., [PostgreSQL](https://www.postgresql.org/), [CockroachDB](https://www.cockroachlabs.com/), [MongoDB](https://www.mongodb.com/), etc.).
- Build and demonstrate a sample application to showcase the practical application of these concepts in Go.
- etc.

## Architecture

### Actor Model

An actor is a highly intelligent abstraction that provides a simplified approach to constructing complex concurrent systems. It is a protected state with only a reference available for communication through message exchanges.

Messages sent to an actor are placed in the actor's mailbox. When the actor is ready, it will consume the message. The consumption of a message can lead to:

- A change in state,
- Responding to the sender (who is also an actor) with another message,
- Altering the actor's behavior.
- Actors exist within an actor system - a hierarchical structure, akin to a tree. There will only be one actor with a given name/identifier throughout the entire actor system.

### Event Sourcing

Event Sourcing involves integrating its principles into the actor system, treating each message sent to an actor as a command. While an actor can handle all states of the domain, it's more relevant to consider a separate actor, known as an aggregate, responsible for managing the state of a domain identifier. This aggregate actor is an event-sourced actor, also referred to as a persistent actor.

In this context, events are generated from the command, representing the effect of the command. These events are then persisted. After successful persistence, they are used to change the actor’s state. The aggregate, therefore, manages the state of a domain identifier, responds to commands, and generates events. These events are sent to the actor who issued the command, and the state's persistence is achieved by consuming the Event Stream of the aggregate for the specified identifier.

``` mermaid
flowchart LR
    O(["actor X"])
    A(["init/usr/book_lend"])
    B(["init/usr/book_lend/_id"])

    es[(Event\nStore)]

    O --> | RegisterBook _id | A
    A -.-> | spawn _id | B
    B --> | restore state | B
    A --> | RegisterBook _id | B
    es -.-> | event stream | B
```
