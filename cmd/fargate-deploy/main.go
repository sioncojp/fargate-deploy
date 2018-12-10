package main

import (
	"os"

	"github.com/sioncojp/fargate-deploy"
)

func main() {
	os.Exit(fargatedeploy.Run())
}
