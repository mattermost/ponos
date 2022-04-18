// Copyright (c) 2015-present Mattermost, Inc. All Rights Reserved.
// See LICENSE.txt for license information.
//

package terraform

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/pkg/errors"
)

type terraformOutput struct {
	Sensitive bool            `json:"sensitive"`
	Type      json.RawMessage `json:"type"`
	Value     interface{}     `json:"value"`
}

// Init invokes terraform init.
func (c *Cmd) Init(remoteKey string) error {
	_, _, err := c.run(
		"init",
		Arg("backend-config", fmt.Sprintf("bucket=%s", c.remoteStateBucket)),
		Arg("backend-config", fmt.Sprintf("key=%s/%s", remoteStateDirectory, remoteKey)),
		Arg("backend-config", fmt.Sprintf("region=%s", os.Getenv("AWS_REGION"))),
	)
	if err != nil {
		return errors.Wrap(err, "failed to invoke terraform init")
	}

	return nil
}

// Plan invokes terraform Plan.
func (c *Cmd) Plan(argExtras ...string) error {
	args := []string{
		"plan",
		Arg("input", "false"),
	}
	args = append(args, argExtras...)
	_, _, err := c.run(args...)
	if err != nil {
		return errors.Wrap(err, "failed to invoke terraform plan")
	}

	return nil
}

// Apply invokes terraform apply.
func (c *Cmd) Apply(argExtras ...string) error {
	args := []string{
		"apply",
		Arg("input", "false"),
	}
	args = append(args, argExtras...)
	args = append(args, Arg("auto-approve"))
	if _, _, err := c.run(args...); err != nil {
		return errors.Wrap(err, "failed to invoke terraform apply")
	}

	return nil
}

// ApplyTarget invokes terraform apply with the given target.
func (c *Cmd) ApplyTarget(target string) error {
	_, _, err := c.run(
		"apply",
		Arg("input", "false"),
		Arg("target", target),
		Arg("auto-approve"),
	)
	if err != nil {
		return errors.Wrap(err, "failed to invoke terraform apply")
	}

	return nil
}

// Destroy invokes terraform destroy.
func (c *Cmd) Destroy(argExtras ...string) error {
	args := []string{
		"destroy",
	}
	args = append(args, argExtras...)
	args = append(args, Arg("auto-approve"))

	if _, _, err := c.run(args...); err != nil {
		return errors.Wrap(err, "failed to invoke terraform destroy")
	}

	return nil
}

// Output invokes terraform output and returns the named value, true if it exists, and an empty
// string and false if it does not.
func (c *Cmd) Output(variable string) (string, bool, error) {
	stdout, _, err := c.run(
		"output",
		"-json",
	)
	if err != nil {
		return string(stdout), false, errors.Wrap(err, "failed to invoke terraform output")
	}

	var outputs map[string]terraformOutput
	err = json.Unmarshal(stdout, &outputs)
	if err != nil {
		return string(stdout), false, errors.Wrap(err, "failed to parse terraform output")
	}

	value, ok := outputs[variable]

	return fmt.Sprintf("%s", value.Value), ok, nil
}

// Version invokes terraform version and returns the value.
func (c *Cmd) Version(removeUpgradeWarning bool) (string, error) {
	stdout, _, err := c.run("version")
	trimmed := strings.TrimSuffix(string(stdout), "\n")
	if err != nil {
		return trimmed, errors.Wrap(err, "failed to invoke terraform version")
	}

	// The terraform version command will print an upgrade warning if running
	// an older version. Optionally attempt to remove the warning before returning.
	if removeUpgradeWarning {
		trimmed = strings.Split(trimmed, "\n")[0]
	}

	return trimmed, nil
}
