package pkg_test

import (
	"testing"

	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/tormath1/sectl/pkg"
)

func TestGetStatus(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		// mock the actual filesystem
		fs := afero.NewMemMapFs()
		err := fs.MkdirAll("/sys/fs/selinux/", 0777)
		require.Nil(t, err)
		err = fs.MkdirAll("/etc/selinux/", 0777)
		require.Nil(t, err)

		enf, err := fs.Create("/sys/fs/selinux/enforce")
		require.Nil(t, err)

		_, err = enf.Write([]byte{49})
		require.Nil(t, err)
		enf.Close()

		p, err := fs.Create("/sys/fs/selinux/policyvers")
		require.Nil(t, err)

		_, err = p.Write([]byte{51, 51})
		require.Nil(t, err)
		p.Close()

		c, err := fs.Create("/sys/fs/selinux/checkreqprot")
		require.Nil(t, err)

		_, err = c.Write([]byte{48})
		require.Nil(t, err)
		c.Close()

		d, err := fs.Create("/sys/fs/selinux/deny_unknown")
		require.Nil(t, err)

		_, err = d.Write([]byte{48})
		require.Nil(t, err)
		d.Close()

		conf, err := fs.Create("/etc/selinux/config")
		require.Nil(t, err)

		_, err = conf.Write([]byte(`# This file controls the state of SELinux on the system on boot.

# SELINUX can take one of these three values:
#	enforcing - SELinux security policy is enforced.
#	permissive - SELinux prints warnings instead of enforcing.
#	disabled - No SELinux policy is loaded.
SELINUX=permissive

# SELINUXTYPE can take one of these four values:
#	targeted - Only targeted network daemons are protected.
#	strict   - Full SELinux protection.
#	mls      - Full SELinux protection with Multi-Level Security
#	mcs      - Full SELinux protection with Multi-Category Security
#	           (mls, but only one sensitivity level)
SELINUXTYPE=mcs`))
		require.Nil(t, err)
		conf.Close()

		status, err := pkg.GetStatus(fs, "/etc/selinux/config")
		require.Nil(t, err)

		assert.Equal(t, "enforcing", status.CurrentMode)
		assert.Equal(t, "permissive", status.Mode)
	})
}
