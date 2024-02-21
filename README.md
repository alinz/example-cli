# STEPS

# Step 1

get self-update library using

```
go get selfupdate.blockthrough.com@latest
```

# Step 2

install the self-update cli

```
go install selfupdate.blockthrough.com/cmd/selfupdate@latest
```

# Step 3

create public/private key pair

```
selfupdate crypto keys
```

# Step 4

create two `repository secret` variables using the provide the link below. Make sure the change the `example-cli` with your project name

- EXAMPLE_CLI_SELFUPDATE_PRIVATE_KEY, and put the content of `selfupdate.key`
- EXAMPLE_CLI_SELFUPDATE_PUBLIC_KEY, and put the content of `selfupdate.pub`

```
https://github.com/alinz/example-cli/settings/secrets/actions
```

> NOTE: make sure to keep `selfupdate.key` safe as it will be used to sign the binary of your program

# Step 5

Create a fine-grain github token by going to this [URL](https://github.com/settings/tokens?type=beta)

- Click on `Generate new token`
- Give your token a name
- Set your token an expiration based on your need
- Select your Organization from `Resource owner` that your project is
- From `Repository access` select `Only select repositories` and chose the repository and project that your cli is
- From `Repository permission`, make sure `Contents` has `Read-only` access
- Click on `Generate token`
- Save that token as you won't be able to access it

# Step 6

add the following Go code to your main execution, it is recommended to add it at the beginning of the main function. Make sure to set these variables `ownerName`, `repoName`, `execName`, according to your project needs. `Version` and `PublicKey` will be set during the build times so leave them empty

```golang
package main

import (
	"context"
	"fmt"
	"os"

	"selfupdate.blockthrough.com"
)

const (
	Version   = ""
	PublicKey = ""

	ownerName = "alinz"
	repoName  = "example-cli"
	execName  = "example-cli"
)

func main() {
	runUpdate()

	// this is the rest of your main function
    // ...
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
		fmt.Fprintf(os.Stderr, "Warning: EXAMPLE_CLI_GH_TOKEN env is not set, selfupdating is disabled")
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
```

# Step 7

Create a new Github actions workflow as follows

```yaml
name: Build and Release
on:
  push:
    tags:
      - "v*"
jobs:
  build:
    runs-on: ubuntu-latest

    env:
      EXAMPLE_CLI_OWNER: alinz
      EXAMPLE_CLI_REPO: example-cli
      EXAMPLE_CLI_GH_TOKEN: ${{ secrets.GITHUB_TOKEN }}
      EXAMPLE_CLI_PRIVATE_KEY: ${{ secrets.EXAMPLE_CLI_SELFUPDATE_PRIVATE_KEY }}
      EXAMPLE_CLI_PUBLIC_KEY: ${{ secrets.EXAMPLE_CLI_SELFUPDATE_PUBLIC_KEY }}
      EXAMPLE_CLI_VERSION: ${{ github.ref_name }}

    steps:
      - name: Setup Go 1.22
        uses: actions/setup-go@v3
        with:
          go-version: ^1.22

      - name: Setup Repo
        uses: actions/checkout@v3

      - name: Install latest selfupdate cli
        run: go install selfupdate.blockthrough.com/cmd/selfupdate@latest

      - name: Create a Release
        run: |
          selfupdate github release \
          --owner ${{ env.EXAMPLE_CLI_OWNER }} \
          --repo ${{ env.EXAMPLE_CLI_REPO }} \
          --version ${{ env.EXAMPLE_CLI_VERSION }} \
          --title ${{ env.EXAMPLE_CLI_VERSION }} \
          --token ${{ env.EXAMPLE_CLI_GH_TOKEN }}

      - name: Build Darwin arm64
        run: |
          GOOS=darwin GOARCH=arm64 go build \
          -ldflags "-X main.Version=${{ env.EXAMPLE_CLI_VERSION }} -X main.PublicKey=${{ env.EXAMPLE_CLI_PUBLIC_KEY }}" \
          -o ./example-cli-darwin-arm64 \
          ./main.go

      - name: Upload Darwin arm64
        run: |
          selfupdate github upload \
          --owner blockthrough \
          --repo selfupdate.go \
          --filename example-cli-darwin-arm64.sign \
          --version ${{ env.EXAMPLE_CLI_VERSION }} \
          --token ${{ env.EXAMPLE_CLI_GH_TOKEN }} \
          --key ${{ env.EXAMPLE_CLI_PRIVATE_KEY }} < ./example-cli-darwin-arm64

          selfupdate github upload \
          --owner blockthrough \
          --repo selfupdate.go \
          --filename example-cli-darwin-arm64 \
          --version ${{ env.EXAMPLE_CLI_VERSION }} \
          --token ${{ env.EXAMPLE_CLI_GH_TOKEN }} < ./example-cli-darwin-arm64
```

# Step 8

Make sure everyone sets `EXAMPLE_CLI_GH_TOKEN` in their system as an env and sets it to the fine-grain token which was described in `Step 5`
