# Simple Service

This is a lightweight service that returns either the name of the Kubernetes pod it's running in or the hostname.

```json
{
  "host": "this-is-a-simple-service"
}
```

Endpoints served by this service,

- `/` - main root that returns the `hostname`
- `/healthz` - health endpoint

## Build

- Use `make build` to build the processor.
- To build and push the Docker images use `PUSH_MULTIARCH=true make docker`. By default, it only builds `linux/amd64` & `linux/arm64`.
    - If podman is available it will use `podman build` otherwise will fallback to `docker buildx`.
    - The images get pushed to `australia-southeast1-docker.pkg.dev/solo-test-236622/apac` but you can override this with the env var `REPO`.
- Run make help for all the build directives.

## Testing

`make test` will run the test suite.