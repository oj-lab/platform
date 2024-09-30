package main

import (
	"context"
	"os"

	user_service "github.com/oj-lab/platform/services/user"
)

func main() {
	args := os.Args
	ctx := context.Background()
	if len(args) > 1 && args[1] == "problem_only" {
		loadProblemPackages(ctx)
	} else {
		initDB()
		loadCasbinPolicies()
		err := user_service.GrantUserRole(ctx, "root", "super", "system")
		if err != nil {
			panic(err)
		}
		loadProblemPackages(ctx)
	}

	println("init success")
}
