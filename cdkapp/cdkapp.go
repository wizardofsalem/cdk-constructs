// Package cdkapp provides a wrapper for creating a CDK app with environment configuration.
package cdkapp

import (
	"fmt"

	"github.com/wizardofsalem/cdk-constructs/awsutil"

	"github.com/aws/aws-cdk-go/awscdk/v2"
	"github.com/aws/jsii-runtime-go"
)

// CDKApp wraps an AWS CDK App with its resolved environment (dev/prod).
type CDKApp struct {
	App         awscdk.App
	Environment awsutil.Env
}

// GetRemovalPolicyFromEnv returns RETAIN for prod, DESTROY for dev.
func GetRemovalPolicyFromEnv(env awsutil.Env) awscdk.RemovalPolicy {
	switch env {
	case awsutil.Prod:
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

func configureEnvironment(app awscdk.App) awsutil.Env {
	envInterface := app.Node().TryGetContext(jsii.String("env"))
	if envInterface == nil {
		panic("Could not find CDK environment — pass -c env=dev or -c env=prod")
	}

	envString, ok := envInterface.(string)
	if !ok {
		panic(fmt.Sprintf("env context value is not a string: %v", envInterface))
	}

	if envString != "dev" && envString != "prod" {
		panic(fmt.Sprintf("invalid env value: %q (expected 'dev' or 'prod')", envString))
	}

	environment, err := awsutil.ParseEnv(envString)
	if err != nil {
		panic(err)
	}

	return environment
}
