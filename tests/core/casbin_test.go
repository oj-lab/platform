package core_test

import (
	"net/http"
	"testing"

	"github.com/oj-lab/oj-lab-platform/modules/auth"
)

func TestCasbin(t *testing.T) {
	enforcer := auth.GetDefaultCasbinEnforcer()
	// Callback like SavePolicy should trigger the watcher to update the policy
	err := enforcer.SavePolicy()
	if err != nil {
		t.Error(err)
	}

	policies, err := enforcer.GetFilteredPolicy(3, `testData`)
	if err != nil {
		t.Error(err)
	}
	t.Logf("Policies: %v", policies)

	allow, err := enforcer.Enforce("admin", "_", `system`, `testData`, http.MethodGet)
	if err != nil {
		t.Error(err)
	}
	if !allow {
		t.Error("Expected to allow")
	}

	allow, err = enforcer.Enforce("test_user", "_", `system`, `adminRequired`, http.MethodDelete)
	if err != nil {
		t.Error(err)
	}
	if !allow {
		t.Error("Expected to allow")
	}
}
