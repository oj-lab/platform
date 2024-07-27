package main

import (
	"context"

	user_service "github.com/oj-lab/oj-lab-platform/services/user"
)

func main() {
	ctx := context.Background()
	initDB()
	loadCasbinPolicies()
	err := user_service.GrantUserRole(ctx, "root", "super", "system")
	if err != nil {
		panic(err)
	}

	loadProblemPackages(ctx)
	println("init success")
}
