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
