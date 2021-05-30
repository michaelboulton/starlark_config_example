package main

import (
	"context"
	"encoding/json"
	"fmt"

	_ "github.com/argoproj/argo-workflows/v3/pkg/apiclient/workflowtemplate"
	"github.com/argoproj/argo-workflows/v3/pkg/apis/workflow/v1alpha1"
	_ "github.com/argoproj/argo-workflows/v3/pkg/apis/workflow/v1alpha1"
	"github.com/bazelbuild/rules_go/go/tools/bazel"
	proto2 "github.com/golang/protobuf/proto"
	"github.com/pkg/errors"
	"github.com/stripe/skycfg"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoregistry"
	"gopkg.in/yaml.v2"

	_ "github.com/michaelboulton/starlark_config_example/v1/prometheus"
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
	config, err := skycfg.Load(
		ctx,
		file,
		skycfg.WithGlobals(skycfg.UnstablePredeclaredModules(nil)),
	)
	if err != nil {
		return errors.Wrap(err, "loading files")
	}

	err = RegisterGogoProtoType(&v1alpha1.Workflow{})
	if err != nil {
		return errors.Wrap(err, "")
	}

	messages, err := config.Main(ctx)
	if err != nil {
		return errors.Wrap(err, "running file")
	}

	for _, m := range messages {
		messageV2 := proto2.MessageV2(m)

		out, err := toYaml(messageV2)
		if err != nil {
			return err
		}

		fmt.Println(messageV2.ProtoReflect().Descriptor().FullName())
		fmt.Println(string(out))
	}

	return nil
}

func RegisterGogoProtoType(oldMessage proto2.Message) error {
	message := proto2.MessageV2(oldMessage).ProtoReflect()
	err := protoregistry.GlobalTypes.RegisterMessage(message.Type())
	if err != nil {
		return errors.Wrapf(err, "registering message %s", message.Descriptor().FullName())
	}

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
