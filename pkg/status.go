package pkg

import (
	"bufio"
	"fmt"
	"os"
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
	Status                   string
	Mount                    string
	RootDirectory            string
	CurrentMode              string
	PolicyDenyUnknownStatus  string
	MemoryProtectionChecking string
	MaxKernelPolicyVersion   int
	*Config
}

func GetStatus(fs afero.Fs, configFile string) (*Status, error) {
	var s Status
	s.Status = "enabled"
	s.RootDirectory = configFile

	// TODO: can be in other location ?
	mounted := "/sys/fs/selinux"
	f, err := os.Open(mounted)
	if err != nil {
		s.Status = "disabled"
	}

	f.Close()

	s.Mount = mounted

	return &s, nil

}

// ReadConfig load the config from configFile
func ReadConfig(fs afero.Fs, configFile string) (*Config, error) {
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
