package config

import (
	"os"
	"testing"
)

func TestConfigLoadDefaults(t *testing.T) {
	// Clear relevant environment variables
	os.Unsetenv("PORT")
	os.Unsetenv("WORKSPACES_DIR")
	os.Unsetenv("HOST_WORKSPACES_DIR")
	os.Unsetenv("SANDBOX_RUNTIME")
	os.Unsetenv("SANDBOX_CLEANUP_WORKSPACE")
	os.Unsetenv("SANDBOX_READ_ONLY_ROOTFS")
	os.Unsetenv("SANDBOX_TMPFS_SIZE")
	os.Unsetenv("SANDBOX_DROP_CAPABILITIES")
	os.Unsetenv("SANDBOX_NETWORK_DISABLED")
	os.Unsetenv("SANDBOX_CPU_LIMIT")
	os.Unsetenv("SANDBOX_MEMORY_LIMIT")
	os.Unsetenv("SANDBOX_PIDS_LIMIT")
	os.Unsetenv("SANDBOX_NO_NEW_PRIVILEGES")

	cfg := Load()

	if cfg.Port != "8080" {
		t.Errorf("Expected default Port to be 8080, got %s", cfg.Port)
	}
	if cfg.WorkspacesDir != "workspaces" {
		t.Errorf("Expected default WorkspacesDir to be workspaces, got %s", cfg.WorkspacesDir)
	}
	if cfg.HostWorkspacesDir != "workspaces" {
		t.Errorf("Expected default HostWorkspacesDir to be workspaces, got %s", cfg.HostWorkspacesDir)
	}
	if cfg.Runtime != "runc" {
		t.Errorf("Expected default Runtime to be runc, got %s", cfg.Runtime)
	}
	if !cfg.CleanupWorkspace {
		t.Error("Expected default CleanupWorkspace to be true")
	}
	if !cfg.ReadOnlyRootfs {
		t.Error("Expected default ReadOnlyRootfs to be true")
	}
	if cfg.TmpfsSize != "64m" {
		t.Errorf("Expected default TmpfsSize to be 64m, got %s", cfg.TmpfsSize)
	}
	if !cfg.DropCapabilities {
		t.Error("Expected default DropCapabilities to be true")
	}
	if !cfg.NetworkDisabled {
		t.Error("Expected default NetworkDisabled to be true")
	}
	if cfg.CPULimit != 1.0 {
		t.Errorf("Expected default CPULimit to be 1.0, got %f", cfg.CPULimit)
	}
	if cfg.MemoryLimit != 512*1024*1024 {
		t.Errorf("Expected default MemoryLimit to be 536870912, got %d", cfg.MemoryLimit)
	}
	if cfg.PidsLimit != 100 {
		t.Errorf("Expected default PidsLimit to be 100, got %d", cfg.PidsLimit)
	}
	if !cfg.NoNewPrivileges {
		t.Error("Expected default NoNewPrivileges to be true")
	}
}

func TestConfigLoadEnvOverrides(t *testing.T) {
	os.Setenv("PORT", "9090")
	os.Setenv("WORKSPACES_DIR", "/local/workspaces")
	os.Setenv("HOST_WORKSPACES_DIR", "/host/workspaces")
	os.Setenv("SANDBOX_RUNTIME", "runsc")
	os.Setenv("SANDBOX_CLEANUP_WORKSPACE", "false")
	os.Setenv("SANDBOX_READ_ONLY_ROOTFS", "false")
	os.Setenv("SANDBOX_TMPFS_SIZE", "128m")
	os.Setenv("SANDBOX_DROP_CAPABILITIES", "false")
	os.Setenv("SANDBOX_NETWORK_DISABLED", "false")
	os.Setenv("SANDBOX_CPU_LIMIT", "2.5")
	os.Setenv("SANDBOX_MEMORY_LIMIT", "1073741824")
	os.Setenv("SANDBOX_PIDS_LIMIT", "200")
	os.Setenv("SANDBOX_NO_NEW_PRIVILEGES", "false")

	defer func() {
		os.Unsetenv("PORT")
		os.Unsetenv("WORKSPACES_DIR")
		os.Unsetenv("HOST_WORKSPACES_DIR")
		os.Unsetenv("SANDBOX_RUNTIME")
		os.Unsetenv("SANDBOX_CLEANUP_WORKSPACE")
		os.Unsetenv("SANDBOX_READ_ONLY_ROOTFS")
		os.Unsetenv("SANDBOX_TMPFS_SIZE")
		os.Unsetenv("SANDBOX_DROP_CAPABILITIES")
		os.Unsetenv("SANDBOX_NETWORK_DISABLED")
		os.Unsetenv("SANDBOX_CPU_LIMIT")
		os.Unsetenv("SANDBOX_MEMORY_LIMIT")
		os.Unsetenv("SANDBOX_PIDS_LIMIT")
		os.Unsetenv("SANDBOX_NO_NEW_PRIVILEGES")
	}()

	cfg := Load()

	if cfg.Port != "9090" {
		t.Errorf("Expected Port to be 9090, got %s", cfg.Port)
	}
	if cfg.WorkspacesDir != "/local/workspaces" {
		t.Errorf("Expected WorkspacesDir to be /local/workspaces, got %s", cfg.WorkspacesDir)
	}
	if cfg.HostWorkspacesDir != "/host/workspaces" {
		t.Errorf("Expected HostWorkspacesDir to be /host/workspaces, got %s", cfg.HostWorkspacesDir)
	}
	if cfg.Runtime != "runsc" {
		t.Errorf("Expected Runtime to be runsc, got %s", cfg.Runtime)
	}
	if cfg.CleanupWorkspace {
		t.Error("Expected CleanupWorkspace to be false")
	}
	if cfg.ReadOnlyRootfs {
		t.Error("Expected ReadOnlyRootfs to be false")
	}
	if cfg.TmpfsSize != "128m" {
		t.Errorf("Expected TmpfsSize to be 128m, got %s", cfg.TmpfsSize)
	}
	if cfg.DropCapabilities {
		t.Error("Expected DropCapabilities to be false")
	}
	if cfg.NetworkDisabled {
		t.Error("Expected NetworkDisabled to be false")
	}
	if cfg.CPULimit != 2.5 {
		t.Errorf("Expected CPULimit to be 2.5, got %f", cfg.CPULimit)
	}
	if cfg.MemoryLimit != 1073741824 {
		t.Errorf("Expected MemoryLimit to be 1073741824, got %d", cfg.MemoryLimit)
	}
	if cfg.PidsLimit != 200 {
		t.Errorf("Expected PidsLimit to be 200, got %d", cfg.PidsLimit)
	}
	if cfg.NoNewPrivileges {
		t.Error("Expected NoNewPrivileges to be false")
	}
}
