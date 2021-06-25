package pkg

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"strings"

	"github.com/spf13/afero"
)

// Config is read from configFile
type Config struct {
	// Mode is the SELinux mode (enforced / permissive)
	Mode string `json:"mode"`

	// LoadedPolicyName is the loaded policy  (mcs, targeted, etc.)
	LoadedPolicyName string `json:"loaded_policy_name"`
}

// Status is the current SELinux status
type Status struct {
	Status                   string `json:"status"`
	Mount                    string `json:"mount"`
	RootDirectory            string `json:"root_directory"`
	CurrentMode              string `json:"current_mode"`
	PolicyDenyUnknownStatus  string `json:"policy_deny_unknown_status"`
	MemoryProtectionChecking string `json:"memory_protection_checking"`
	MaxKernelPolicyVersion   string `json:"max_kernel_policy_version"`
	*Config
}

func getProperty(fs afero.Fs, name string) (string, error) {
	property, err := fs.Open(fmt.Sprintf("/sys/fs/selinux/%s", name))
	if err != nil {
		return "", fmt.Errorf("unable to open file: %v", err)
	}

	defer property.Close()

	p, err := ioutil.ReadAll(property)
	if err != nil {
		return "", fmt.Errorf("unable to read mode from file: %v", err)
	}

	return string(p), nil
}

// GetStatus reads the status from the config file and from the mounted SELinux
// fs if this one exist
func GetStatus(fs afero.Fs, configFile string) (*Status, error) {
	// TODO: can be in other location ?
	mounted := "/sys/fs/selinux"

	var s Status
	s.Status = "enabled"

	f, err := fs.Open(mounted)
	if err != nil {
		s.Status = "disabled"
		return &s, nil
	}

	f.Close()

	s.Mount = mounted
	s.RootDirectory = configFile

	c, err := readConfig(fs, configFile)
	if err != nil {
		return nil, fmt.Errorf("unable to read config: %v", err)
	}

	s.Config = c

	mode, err := getProperty(fs, "enforce")
	if err != nil {
		return nil, fmt.Errorf("unable to get enforce property: %w", err)
	}

	s.CurrentMode = "permissive"
	if mode != "0" {
		s.CurrentMode = "enforcing"
	}

	policy, err := getProperty(fs, "policyvers")
	if err != nil {
		return nil, fmt.Errorf("unable to get policy version property: %w", err)
	}

	s.MaxKernelPolicyVersion = policy

	checkreqprot, err := getProperty(fs, "checkreqprot")
	if err != nil {
		return nil, fmt.Errorf("unable to get checkreqprot property: %w", err)
	}

	s.MemoryProtectionChecking = "actual (secure)"

	if checkreqprot != "0" {
		s.MemoryProtectionChecking = "requested (insecure)"
	}

	denyUnknownStatus, err := getProperty(fs, "deny_unknown")
	if err != nil {
		return nil, fmt.Errorf("unable to get checkreqprot property: %w", err)
	}

	s.PolicyDenyUnknownStatus = "allowed"
	if denyUnknownStatus != "0" {
		// TODO: find the correct term
		s.MemoryProtectionChecking = "not allowed"
	}

	return &s, nil
}

// readConfig load the config from configFile
func readConfig(fs afero.Fs, configFile string) (*Config, error) {
	conf, err := fs.Open(configFile)
	if err != nil {
		return nil, fmt.Errorf("unable to open SELinux config: %w", err)
	}

	var c Config

	s := bufio.NewScanner(conf)
	for s.Scan() {
		line := s.Text()
		if strings.HasPrefix(line, "#") {
			continue
		}

		if len(line) == 0 {
			continue
		}

		configuration := strings.SplitN(line, "=", 2)
		if len(configuration) < 2 {
			return nil, fmt.Errorf("invalid entry, should be key=value. Got: %s", line)
		}

		switch configuration[0] {
		case "SELINUX":
			c.Mode = configuration[1]
		case "SELINUXTYPE":
			c.LoadedPolicyName = configuration[1]
		}
	}

	return &c, nil
}
