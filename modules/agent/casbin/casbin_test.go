package casbin_agent

import (
	"net/http"
	"strings"
	"testing"
)

func TestKeyMatchGin(t *testing.T) {
	key2 := "/api/v1/user/:id"
	key1 := "/api/v1/user/1"
	if !KeyMatchGin(key1, key2) {
		t.Error("Expected to match")
	}
	key1 = "/api/v1/user/"
	if KeyMatchGin(key1, key2) {
		t.Error("Expected not to match")
	}

	key2 = "/api/v1/:resource/*any"
	key1 = "/api/v1/user"
	if !KeyMatchGin(key1, key2) {
		t.Error("Expected to match")
	}
	key1 = "/api/v1/user/1"
	if !KeyMatchGin(key1, key2) {
		t.Error("Expected to match")
	}
	key1 = "/api/v1/user/"
	if !KeyMatchGin(key1, key2) {
		t.Error("Expected to match")
	}
	key1 = "/api/v1/user/1/send"
	if !KeyMatchGin(key1, key2) {
		t.Error("Expected to match")
	}
	key1 = "/api/v1//"
	if KeyMatchGin(key1, key2) {
		t.Error("Expected not to match")
	}
}

func TestCasbin(t *testing.T) {
	enforcer := GetDefaultCasbinEnforcer()
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

	allow, err := enforcer.Enforce("user_test", ExtraInfo{
		IsVIP: true,
	}, `system`, `testData`, http.MethodGet)
	if err != nil {
		t.Error(err)
	}
	if !allow {
		t.Error("Expected to allow")
	}
	allow, err = enforcer.Enforce("user_test", ExtraInfo{
		IsVIP: false,
	}, `system`, `testData`, http.MethodGet)
	if err != nil {
		t.Error(err)
	}
	if allow {
		t.Error("Expected to not allow")
	}

	_, err = enforcer.AddPolicies([][]string{
		{
			RoleSubjectPrefix + `admin`, `true`, `system`,
			`/api/v1/user/*any`, strings.Join([]string{http.MethodGet, http.MethodPost, http.MethodDelete}, "|"),
			"allow",
		},
	})
	if err != nil {
		t.Error(err)
	}
	err = enforcer.SavePolicy()
	if err != nil {
		t.Error(err)
	}
	_, err = enforcer.AddGroupingPolicies([][]string{
		{`role:super`, `role:admin`, `system`},
		{`user:root`, `role:super`, `system`},
	})
	if err != nil {
		t.Error(err)
	}
	err = enforcer.SavePolicy()
	if err != nil {
		t.Error(err)
	}
	roles, err = enforcer.GetImplicitRolesForUser("user:root", "system")
	if err != nil {
		t.Error(err)
	}
	t.Logf("Implicit Roles: %v", roles)

	allow, err = enforcer.Enforce("user:root", "_", "system", "/api/v1/user", http.MethodGet)
	if err != nil {
		t.Error(err)
	}
	if !allow {
		t.Error("Expected to allow")
	}
	enforcer.ClearPolicy()
}
