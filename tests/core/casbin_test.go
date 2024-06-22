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

	policies, err := enforcer.GetFilteredPolicy(1, `testData`)
	if err != nil {
		t.Error(err)
	}
	t.Logf("Policies: %v", policies)

	subject := auth.CasbinSubject{Age: 30}
	allow, err := enforcer.Enforce(subject, `testData`, http.MethodGet)
	if err != nil {
		t.Error(err)
	}
	if !allow {
		t.Error("Expected to allow")
	}

	subject = auth.CasbinSubject{Age: 30, Role: "admin"}
	allow, err = enforcer.Enforce(subject, `adminRequired`, http.MethodDelete)
	if err != nil {
		t.Error(err)
	}
	if !allow {
		t.Error("Expected to allow")
	}
}
