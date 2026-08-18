package cdkutil

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strings"

	"github.com/aws/aws-cdk-go/awscdk/v2/awsecrassets"
	"github.com/aws/aws-cdk-go/awscdk/v2/awslambda"
	"github.com/aws/jsii-runtime-go"
)

// LambdaImageAsset creates a DockerImageCode with a custom asset hash based on
// only the Go files the lambda transitively imports. This prevents CDK from
// rebuilding containers when unrelated files change.
//
// repoRoot should be the absolute path to the repository root (where go.mod lives).
// dockerfile is relative to repoRoot.
// lambdaPkg is the Go import path of the lambda's main package.
func LambdaImageAsset(repoRoot string, dockerfile string, lambdaPkg string, buildArgs *map[string]*string) awslambda.DockerImageCode {
	hash := computeLambdaHash(repoRoot, lambdaPkg)

	return awslambda.DockerImageCode_FromImageAsset(jsii.String(repoRoot), &awslambda.AssetImageCodeProps{
		File:      jsii.String(dockerfile),
		BuildArgs: buildArgs,
		ExtraHash: jsii.String(hash),
		Invalidation: &awsecrassets.DockerImageAssetInvalidationOptions{
			BuildArgs:      jsii.Bool(false),
			ExtraHash:      jsii.Bool(true),
			File:           jsii.Bool(false),
			NetworkMode:    jsii.Bool(false),
			Outputs:        jsii.Bool(false),
			Platform:       jsii.Bool(false),
			Target:         jsii.Bool(false),
			BuildSecrets:   jsii.Bool(false),
			BuildSsh:       jsii.Bool(false),
			RepositoryName: jsii.Bool(false),
		},
	})
}

// computeLambdaHash returns a deterministic hash of all .go files that the
// given package transitively depends on, plus go.mod and go.sum.
func computeLambdaHash(repoRoot string, lambdaPkg string) string {
	cmd := exec.Command("go", "list", "-deps", "-f", "{{.Dir}}", lambdaPkg)
	cmd.Dir = repoRoot
	out, err := cmd.Output()
	if err != nil {
		panic(fmt.Sprintf("go list failed for %s: %v", lambdaPkg, err))
	}

	modulePath := repoRoot + "/"
	var dirs []string
	for _, line := range strings.Split(strings.TrimSpace(string(out)), "\n") {
		if strings.HasPrefix(line, modulePath) {
			dirs = append(dirs, line)
		}
	}

	var files []string
	for _, dir := range dirs {
		matches, _ := filepath.Glob(filepath.Join(dir, "*.go"))
		for _, m := range matches {
			rel, _ := filepath.Rel(repoRoot, m)
			files = append(files, rel)
		}
	}

	files = append(files, "go.mod", "go.sum")
	sort.Strings(files)

	h := sha256.New()
	for _, f := range files {
		content, err := os.ReadFile(filepath.Join(repoRoot, f))
		if err != nil {
			panic(fmt.Sprintf("failed to read %s: %v", f, err))
		}
		h.Write([]byte(f + "\n"))
		h.Write(content)
	}
	return hex.EncodeToString(h.Sum(nil))[:16]
}
