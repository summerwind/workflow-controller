package terraform

import (
	"errors"
	"fmt"
	"io"
	"os/exec"
)

type Terraform struct {
	WorkDir string
	Logger  io.Writer
}

func (t *Terraform) Init() error {
	err := t.run("terraform", "init", "-no-color")
	if err != nil {
		return fmt.Errorf("terraform init: %v", err)
	}

	return nil
}

func (t *Terraform) SelectWorkspace(workspace string, create bool) error {
	err := t.run("terraform", "workspace", "select", "-no-color", workspace)
	if err != nil {
		_, ok := err.(*exec.ExitError)
		if !ok {
			return fmt.Errorf("terraform workspace select: %v", err)
		}

		if !create {
			return errors.New("workspace not found")
		}

		err = t.run("terraform", "workspace", "new", "-no-color", workspace)
		if err != nil {
			return fmt.Errorf("terraform workspace new: %v", err)
		}

		err = t.run("terraform", "workspace", "select", "-no-color", workspace)
		if err != nil {
			return fmt.Errorf("terraform workspace select: %v", err)
		}
	}

	return nil
}

func (t *Terraform) Apply(vars map[string]string) error {
	args := []string{"apply", "-auto-approve", "-no-color"}
	for k, v := range vars {
		args = append(args, "-var", fmt.Sprintf("'%s=%s'", k, v))
	}

	err := t.run("terraform", args...)
	if err != nil {
		return fmt.Errorf("terraform apply: %v", err)
	}

	return nil
}

func (t *Terraform) run(name string, arg ...string) error {
	cmd := exec.Command(name, arg...)
	cmd.Dir = t.WorkDir

	if t.Logger != nil {
		cmd.Stdout = t.Logger
		cmd.Stderr = t.Logger
	}

	return cmd.Run()
}
