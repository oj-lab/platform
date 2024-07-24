package core_test

import (
	"net/http"
	"testing"

	casbin_agent "github.com/oj-lab/oj-lab-platform/modules/agent/casbin"
)

func TestKeyMatchGin(t *testing.T) {
	key2 := "/api/v1/user/:id"
	key1 := "/api/v1/user/1"
	if !casbin_agent.KeyMatchGin(key1, key2) {
		t.Error("Expected to match")
	}
	key1 = "/api/v1/user/"
	if casbin_agent.KeyMatchGin(key1, key2) {
		t.Error("Expected not to match")
	}

	key2 = "/api/v1/:resource/*any"
	key1 = "/api/v1/user"
	if !casbin_agent.KeyMatchGin(key1, key2) {
		t.Error("Expected to match")
	}
	key1 = "/api/v1/user/1"
	if !casbin_agent.KeyMatchGin(key1, key2) {
		t.Error("Expected to match")
	}
	key1 = "/api/v1/user/"
	if !casbin_agent.KeyMatchGin(key1, key2) {
		t.Error("Expected to match")
	}
	key1 = "/api/v1/user/1/send"
	if !casbin_agent.KeyMatchGin(key1, key2) {
		t.Error("Expected to match")
	}
	key1 = "/api/v1//"
	if casbin_agent.KeyMatchGin(key1, key2) {
		t.Error("Expected not to match")
	}
}

func TestCasbin(t *testing.T) {
	enforcer := casbin_agent.GetDefaultCasbinEnforcer()
	_, err := enforcer.AddPolicy(
		`user_test`, `r.ext.IsVIP == true`, `system`, `testData`, http.MethodGet, "allow")
	if err != nil {
		panic(err)
	}

	// Callback like SavePolicy should trigger the watcher to update the policy
	err = enforcer.SavePolicy()
	if err != nil {
		t.Error(err)
	}
	roles := enforcer.GetRolesForUserInDomain("user_test", `system`)
	t.Logf("Roles: %v", roles)

	policies, err := enforcer.GetFilteredPolicy(3, `testData`)
	if err != nil {
		t.Error(err)
	}
	t.Logf("Policies: %v", policies)

	allow, err := enforcer.Enforce("user_test", casbin_agent.ExtraInfo{
		IsVIP: true,
	}, `system`, `testData`, http.MethodGet)
	if err != nil {
		t.Error(err)
	}
	if !allow {
		t.Error("Expected to allow")
	}
}