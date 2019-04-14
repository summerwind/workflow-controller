package run

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/summerwind/workflow-controller/pkg/terraform"
)

func Reconcile() error {
	var (
		configPath string
		commit     string
		err        error
	)

	state := State{}
	err = json.NewDecoder(os.Stdin).Decode(&state)
	if err != nil {
		return err
	}

	r := state.Resource

	if r.Spec.Source.Git != nil {
		dir, err := ioutil.TempDir("", "terraform-controller")
		if err != nil {
			return err
		}
		defer os.RemoveAll(dir)

		err = os.Chdir(dir)
		if err != nil {
			return err
		}

		log.Print("Fetching source...")
		commit, err = checkout(r.Spec.Source.Git)
		if err != nil {
			return err
		}

		configPath = filepath.Join(dir, r.Spec.Source.Git.Path)
	} else {
		configPath = r.Spec.Source.File.Path
	}

	_, err = os.Stat(configPath)
	if err != nil {
		return err
	}

	tf := &terraform.Terraform{
		WorkDir: filepath.Dir(configPath),
		Logger:  os.Stderr,
	}

	log.Print("Initializing terraform...")
	err = tf.Init()
	if err != nil {
		return err
	}

	if r.Spec.Workspace != "" {
		log.Print("Changing workspace...")
		err = tf.SelectWorkspace(r.Spec.Workspace, true)
		if err != nil {
			return err
		}
	}

	log.Print("Applying terraform configuration...")
	err = tf.Apply(r.Spec.Vars)
	if err != nil {
		return err
	}

	state.Resource.Status.LastApplyTime = time.Now().Unix()
	state.Resource.Status.LastApplyCommit = commit

	err = json.NewEncoder(os.Stdout).Encode(&state)
	if err != nil {
		return err
	}

	return nil
}

func checkout(git *RunSpecSourceGit) (string, error) {
	err := exec.Command("git", "init").Run()
	if err != nil {
		return "", fmt.Errorf("failed to initialize source repository: %v", err)
	}
	err = exec.Command("git", "remote", "add", "origin", git.URL).Run()
	if err != nil {
		return "", fmt.Errorf("failed to add remote repository: %v", err)
	}

	if strings.HasPrefix(git.Ref, "refs/tags/") {
		err = exec.Command("git", "fetch", "--tags", "origin", fmt.Sprintf("+%s:", git.Ref)).Run()
		if err != nil {
			return "", fmt.Errorf("failed to fetch tags: %v", err)
		}
		err = exec.Command("git", "checkout", "-qf", "FETCH_HEAD").Run()
		if err != nil {
			return "", fmt.Errorf("failed to checkout: %v", err)
		}
	} else {
		err = exec.Command("git", "fetch", "--no-tags", "origin", fmt.Sprintf("+%s:", git.Ref)).Run()
		if err != nil {
			return "", fmt.Errorf("failed to fetch: %v", err)
		}
		err = exec.Command("git", "reset", "--hard", "-q", "FETCH_HEAD").Run()
		if err != nil {
			return "", fmt.Errorf("failed to reset repository: %v", err)
		}
	}

	err = exec.Command("git", "submodule", "update", "--init", "--recursive").Run()
	if err != nil {
		return "", fmt.Errorf("failed to update submodule: %v", err)
	}

	commit, err := exec.Command("git", "rev-parse", "HEAD").Output()
	if err != nil {
		return "", fmt.Errorf("failed to get commit hash: %v", err)
	}

	return string(commit), nil
}
