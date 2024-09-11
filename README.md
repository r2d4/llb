# Low-level Build (LLB) API for Docker

This repository provides the code for a LLB API frontend. It accepts a base64 encoded, binary encoded [protobuf Buildkit Definition](https://github.com/moby/buildkit/blob/9e14164a1099d3e41b58fc879cbdd6f2b2edb04e/solver/pb/ops.proto#L285-L293).

It aims to be a generic [Buildkit frontend](https://docs.docker.com/build/buildkit/frontend/) for developers looking to manipulate the build DAG with the finest granularity.

This may be useful for developers looking at access the low-level build API in Buildkit in languages other than Go. 

## Usage

This frontend accepts a JSON file with the following schema:
- **imageConfig**: a partial [OCI image config](https://github.com/opencontainers/image-spec/blob/main/config.md) that will be written to the exported image upon build
- **definition**: A base64-encoded binary-encoded protobuf [buildkit.Definition](https://github.com/moby/buildkit/blob/3a7055008a5e58a2abbe0e0c21c919d9e014e062/solver/pb/ops.proto#L284-L293) message. The message itself contains marshaled buildkit operations, a metadata map, and a sources map (currently unused by this frontend).
```
#syntax=ghcr.io/r2d4/llb:1.0.3
{
    "imageConfig": {
        "created": "1970-01-01T00:00:00.000Z",
        "os": "linux",
        "architecture": "arm64",
        "config": {
            "Cmd": [
                "sh"
            ],
            "Env": [],
            "WorkingDir": "/",
            "Labels": {
                "r2d4.dacc.version": "1.0.3",
                "r2d4.dacc.builder": "ghcr.io/r2d4/llb",
                "build.date": "2024-09-11T22:07:24.099Z"
            }
        }
    },
    "definition": "CkUaMQovZG9ja2VyLWltYWdlOi8vZG9ja2VyLmlvL2xpYnJhcnkvYnVzeWJveDp1Y2xpYmNSDgoFYXJtNjQSBWxpbnV4WgAKSwpJCkdzaGEyNTY6NTE5NzdlMjUwYjY4MjYzOWNmMjJiNTBhNzllMGMzNzM1ZGJmYjU0YjA1MzUzODcyM2FjYmZmOTQxZmQ0Yjc4MhKYAQpHc2hhMjU2OjUxOTc3ZTI1MGI2ODI2MzljZjIyYjUwYTc5ZTBjMzczNWRiZmI1NGIwNTM1Mzg3MjNhY2JmZjk0MWZkNGI3ODISTRI5Cg5sbGIuY3VzdG9tbmFtZRInW2Zyb21dIGRvY2tlci5pby9saWJyYXJ5L2J1c3lib3g6dWNsaWJjKhAKDHNvdXJjZS5pbWFnZRABEoABCkdzaGEyNTY6N2FkNThmNjQzOTU5OTI0OTlhMjc3MWZiZmExMDdmNjVkZTM1OThmMmU4NWFjNjhmYmM2YWNjZjVhMjVmZjI1NBI1Kg8KC2NvbnN0cmFpbnRzEAEqDAoIcGxhdGZvcm0QASoUChBtZXRhLmRlc2NyaXB0aW9uEAE="
}
```
To generate this configuration file, you can run the [build/src/hello-world](./build/src/hello-world.ts).

For an implementation, see the [dacc](https://github.com/r2d4/dacc) project on GitHub.

## Self-hosted compiler
The frontend is now compiled with [dacc](https://github.com/r2d4/dacc). You can see the build definition in [build/src/main.ts](build/src/main.ts).

### Motivation
The Buildkit Gateway service can be difficult to access for languages other than Go. It requires compiling a series of protobuf and gRPC infrastructure that's currently only written in Go. 

This provides an alternative to writing a frontend in Go. Instead, you only have to generate your favorite language bindings for `ops.proto`. Generating the DAG still requires some work. Higher-level libraries may be helpful.
