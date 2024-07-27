package main

import (
	"github.com/oj-lab/oj-lab-platform/cmd/web_server/handler"
	casbin_agent "github.com/oj-lab/oj-lab-platform/modules/agent/casbin"
)

func loadCasbinPolicies() {
	enforcer := casbin_agent.GetDefaultCasbinEnforcer()

	_, err := enforcer.AddGroupingPolicies([][]string{
		{`user:root`, `role:super`, `system`},
		{`role:super`, `role:admin`, `system`},
	})
	if err != nil {
		panic(err)
	}

	err = handler.AddUserCasbinPolicies()
	if err != nil {
		panic(err)
	}

	err = enforcer.SavePolicy()
	if err != nil {
		panic(err)
	}
}
