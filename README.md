# Low-level Build (LLB) API for Docker

This repository provides the code for a LLB API frontend. It accepts a base64 encoded, binary encoded [protobuf Buildkit Definition](https://github.com/moby/buildkit/blob/9e14164a1099d3e41b58fc879cbdd6f2b2edb04e/solver/pb/ops.proto#L285-L293).

It aims to be a generic [Buildkit frontend](https://docs.docker.com/build/buildkit/frontend/) for developers looking to manipulate the build DAG with the finest granularity.

This may be useful for developers looking at access the low-level build API in Buildkit in languages other than Go. 

The buildkit frontend image can be found at `ghcr.io/r2d4/llb:1.0.3`.

For an implementation, see the [dacc](https://github.com/r2d4/dacc) project on GitHub.

## Self-hosted compiler
The frontend is now compiled with [dacc](https://github.com/r2d4/dacc). You can see the build definition in [build/src/main.ts](build/src/main.ts).

### Motivation
The Buildkit Gateway service can be difficult to access for languages other than Go. It requires compiling a series of protobuf and gRPC infrastructure that's currently only written in Go. 

This provides an alternative to writing a frontend in Go. Instead, you only have to generate your favorite language bindings for `ops.proto`. Generating the DAG still requires some work. Higher-level libraries may be helpful.
