package main

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/bazelbuild/rules_go/go/tools/bazel"
	"github.com/golang/protobuf/proto"
	_ "github.com/golang/protobuf/ptypes/wrappers"
	_ "github.com/michaelboulton/starlark_config_example/v1"
	"github.com/pkg/errors"
	"github.com/stripe/skycfg"
	"gopkg.in/yaml.v2"
)

func main() {
	err := run()
	if err != nil {
		panic(err)
	}
}

func run() error {
	file, err := bazel.Runfile("cfg/prometheus.bzl")
	if err != nil {
		return errors.Wrap(err, "finding file")
	}

	ctx := context.Background()
	config, err := skycfg.Load(ctx, file,
		skycfg.WithGlobals(skycfg.UnstablePredeclaredModules(nil)),
	)
	if err != nil {
		return errors.Wrap(err, "loading files")
	}

	messages, err := config.Main(ctx)
	if err != nil {
		return errors.Wrap(err, "running file")
	}

	out, err := toYaml(messages[0])
	if err != nil {
		return err
	}

	fmt.Println(string(out))

	return nil
}

// toYaml converts it to json, then back into Go, then into yaml
// This is because the go-yaml package does not really do what you would expect https://github.com/go-yaml/yaml/issues/424
func toYaml(message proto.Message) ([]byte, error) {
	asJson, err := json.MarshalIndent(message, "", " ")
	if err != nil {
		return nil, errors.Wrap(err, "marshaling to json")
	}

	var buf map[string]interface{}
	err = json.Unmarshal(asJson, &buf)
	if err != nil {
		return nil, errors.Wrap(err, "unmarshaling")
	}

	out, err := yaml.Marshal(buf)
	if err != nil {
		return nil, errors.Wrap(err, "marshaling to yaml")
	}

	return out, nil
}
