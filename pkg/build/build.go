package build

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/moby/buildkit/client/llb"
	"github.com/moby/buildkit/frontend/gateway/client"
	"github.com/moby/buildkit/solver/pb"
	ocispecs "github.com/opencontainers/image-spec/specs-go/v1"
)

var _ client.BuildFunc = BuildFunc

// BuildFunc is the main entrypoint for the build.
func BuildFunc(ctx context.Context, c client.Client) (*client.Result, error) {
	cfg, err := GetConfig(ctx, c)
	if err != nil {
		return nil, fmt.Errorf("failed to get definition: %w", err)
	}
	result, err := c.Solve(ctx, client.SolveRequest{
		Definition: cfg.Definition,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to solve: %w", err)
	}
	b, err := json.Marshal(cfg.ImageConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal image config: %w", err)
	}
	result.AddMeta("containerimage.config", b)
	return result, nil
}

// GetDefinition reads the json file from the build context and returns the definition
// of the build.
func GetConfig(ctx context.Context, c client.Client) (*Config, error) {
	opts := c.BuildOpts().Opts
	filename := opts["filename"]
	if filename == "" {
		return nil, fmt.Errorf("missing filename")
	}

	src := llb.Local("dockerfile", llb.IncludePatterns([]string{filename}), llb.WithCustomName("llb-json-api"))
	def, err := src.Marshal(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal source: %w", err)
	}
	var dtDockerfile []byte
	res, err := c.Solve(ctx, client.SolveRequest{
		Definition: def.ToPB(),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to solve source: %w", err)
	}

	ref, err := res.SingleRef()
	if err != nil {
		return nil, err
	}
	dtDockerfile, err = ref.ReadFile(ctx, client.ReadRequest{
		Filename: filename,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %w", err)
	}
	return DecodeConfig(dtDockerfile)
}

type Config struct {
	ImageConfig *ocispecs.Image `json:"imageConfig"`
	Definition  *pb.Definition  `json:"definition"`
}

func (c *Config) UnmarshalJSON(data []byte) error {

	type Alias Config
	aux := &struct {
		Definition string `json:"definition"`
		*Alias
	}{
		Alias: (*Alias)(c),
	}

	if err := json.Unmarshal(data, &aux); err != nil {
		return fmt.Errorf("failed to unmarshal Config: %w", err)
	}

	// Decode base64 Definition
	decodedData, err := base64.StdEncoding.DecodeString(aux.Definition)
	if err != nil {
		return fmt.Errorf("failed to decode base64 Definition: %w", err)
	}

	// Unmarshal Definition
	c.Definition = &pb.Definition{}
	if err := c.Definition.Unmarshal(decodedData); err != nil {
		return fmt.Errorf("failed to unmarshal Definition: %w", err)
	}

	return nil
}

func DecodeConfig(data []byte) (*Config, error) {
	data = []byte(strings.Join(strings.Split(string(data), "\n")[1:], "\n"))
	var cfg Config
	if err := json.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("failed to unmarshal Config: %w", err)
	}
	return &cfg, nil
}
