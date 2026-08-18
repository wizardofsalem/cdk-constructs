// Package awsutil contains shared AWS types for environment and region.
package awsutil

import "fmt"

type Region string

const (
	USEast1 Region = "us-east-1"
	USEast2 Region = "us-east-2"
	USWest1 Region = "us-west-1"
	USWest2 Region = "us-west-2"

	EUWest1    Region = "eu-west-1"
	EUWest2    Region = "eu-west-2"
	EUWest3    Region = "eu-west-3"
	EUCentral1 Region = "eu-central-1"
	EUNorth1   Region = "eu-north-1"

	APSoutheast1 Region = "ap-southeast-1"
	APSoutheast2 Region = "ap-southeast-2"
	APNortheast1 Region = "ap-northeast-1"
	APNortheast2 Region = "ap-northeast-2"
	APSouth1     Region = "ap-south-1"

	SACentral1 Region = "sa-east-1"
	CACentral1 Region = "ca-central-1"
	AFSouth1   Region = "af-south-1"
	MESouth1   Region = "me-south-1"
)

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
		panic("Couldn't parse env as string")
	}
}

func ParseEnv(s string) (Env, error) {
	switch s {
	case "dev":
		return Dev, nil
	case "prod":
		return Prod, nil
	default:
		return 0, fmt.Errorf("invalid env: %q", s)
	}
}
