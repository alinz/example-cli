package main

import (
	"context"
	"fmt"
	"os"

	"selfupdate.blockthrough.com"
)

var (
	Version   = ""
	PublicKey = ""

	ownerName = "alinz"
	repoName  = "example-cli"
	execName  = "example-cli"
)

func main() {
	runUpdate()

	// this is the rest of your main function
	fmt.Println("Version", Version)
	fmt.Println("args:", os.Args)
}

func runUpdate() {
	// In order for selfupdating to work, the following conditions must be met:
	// 1. Version must be set
	// 2. EXAMPLE_CLI_GH_TOKEN must be set
	// 3. PublicKey must be set
	// for setting up the token please refer to
	// "Create a Fine-Grained Personal Access Tokens" in README.md
	ghToken, ok := os.LookupEnv("EXAMPLE_CLI_GH_TOKEN")
	if !ok {
		fmt.Fprintf(os.Stderr, "Warning: EXAMPLE_CLI_GH_TOKEN env is not set, selfupdating is disabled\n")
		return
	}

	selfupdate.Auto(
		context.Background(), // Context
		ownerName,            // Owner Name
		repoName,             // Repo Name
		Version,              // Current Version
		execName,             // Executable Name,
		ghToken,              // Github Token
		PublicKey,            // Public Key
	)
}
