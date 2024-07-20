package main

import (
	"context"
)

func main() {
	ctx := context.Background()
	initDB()
	loadCasbinPolicies()
	loadProblemPackages(ctx)
	println("init success")
}
