import { cacheMount, State } from "dacc";
import path from "path";
import { fileURLToPath } from "url";

async function main() {
    const tag = "ghcr.io/r2d4/llb"
    const version = "1.0.3"

    const base = "golang"
    const baseVersion = "1.23-alpine"

    const s = await new State().from(`${base}:${baseVersion}`)

    s.workdir("/app")
        .merge(
            s.parallel(
                s => s.copy(["go.mod", "go.sum"], ".").run("go mod download").with(cacheMount("/go/pkg/mod")),
                ...["cmd", "pkg"].map(dir => (s: State) => s.copy(dir, dir))
            )
        )
        .diff(
            s => s.run("CGO_ENABLED=0 go build -o /go/bin/llb ./cmd/llb")
                .with(cacheMount("/root/.cache/go-build"))
                .with(cacheMount("/go/pkg/mod"))
        ).entrypoint(["/go/bin/llb"])

    const __filename = fileURLToPath(import.meta.url);
    const __dirname = path.dirname(__filename);
    const contextPath = path.resolve(__dirname, "..", "..")

    s.image.build({ contextPath, tag: [`${tag}:${version}`] })
}

void main()