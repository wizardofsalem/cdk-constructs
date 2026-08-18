// Package cdkapp provides a wrapper for creating a CDK app with environment configuration.
package cdkapp

import (
	"fmt"

	"github.com/aws/aws-cdk-go/awscdk/v2"
	"github.com/aws/jsii-runtime-go"
)

// Env represents the deployment environment.
type Env int

const (
	Dev Env = iota
	Prod
)

func (env Env) String() string {
	switch env {
	case Dev:
		return "Dev"
	case Prod:
		return "Prod"
	default:
		panic("unknown env")
	}
}

// ParseEnv converts a string to an Env value.
func ParseEnv(s string) (Env, error) {
	switch s {
	case "dev":
		return Dev, nil
	case "prod":
		return Prod, nil
	default:
		return 0, fmt.Errorf("invalid env: %q (expected 'dev' or 'prod')", s)
	}
}

// CDKApp wraps an AWS CDK App with its resolved environment (dev/prod).
type CDKApp struct {
	App         awscdk.App
	Environment Env
}

// GetRemovalPolicyFromEnv returns RETAIN for prod, DESTROY for dev.
func GetRemovalPolicyFromEnv(env Env) awscdk.RemovalPolicy {
	switch env {
	case Prod:
		return awscdk.RemovalPolicy_RETAIN
	default:
		return awscdk.RemovalPolicy_DESTROY
	}
}

// NewCDKApp creates a new CDK app and reads the "env" context variable
// (passed via -c env=dev or -c env=prod) to determine the environment.
func NewCDKApp() CDKApp {
	app := awscdk.NewApp(nil)
	env := configureEnvironment(app)
	return CDKApp{App: app, Environment: env}
}

func configureEnvironment(app awscdk.App) Env {
	envInterface := app.Node().TryGetContext(jsii.String("env"))
	if envInterface == nil {
		panic("Could not find CDK environment — pass -c env=dev or -c env=prod")
	}

	envString, ok := envInterface.(string)
	if !ok {
		panic(fmt.Sprintf("env context value is not a string: %v", envInterface))
	}

	environment, err := ParseEnv(envString)
	if err != nil {
		panic(err)
	}

	return environment
}
