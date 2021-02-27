package main

import (
	"fmt"

	"github.com/Paypi/paypi-go/gql"
)

var (
	Key    string
	apiUrl string = "localhost:8080/graphql"
)

type AuthenticatedOutput struct {
	ClientToken string
	PublicKey   string
}

// Authenticate checks if a given api token is authorized.
// This checks if the token is not blocked in any way and has
// card details setup on their account to make payments.
func Authenticate(clientToken string) (AuthenticatedOutput, error) {
	// Send request to API to authenticate clientToken

	type Resp struct {
		AuthenticateService struct {
			PublicKey string
		}
	}
	resp := Resp{}
	err := gql.MakeRequest(gql.GqlQuery{
		Query: `
			mutation AuthenticateService($secret: String!) { 
					authenticateService(input: {
						apiSecret: $secret
					}) {
						publicKey
					}
			}
		`,
		Variables: map[string]interface{}{
			"secret": clientToken,
		},
	}, &resp)
	if err != nil {
		fmt.Print(err.Error())
		return AuthenticatedOutput{}, err
	}

	return AuthenticatedOutput{
		ClientToken: clientToken,
		PublicKey:   resp.AuthenticateService.PublicKey,
	}, nil
}

type MakeChargeInput struct {
	ChargeIdentifier string
}

type MakeChargeOutput struct {
	Success bool
}

// MakeCharge requests to take payments from the user, totaling the cost of the
// given ChargeIdentifier.
func (a AuthenticatedOutput) MakeCharge(input MakeChargeInput) (MakeChargeOutput, error) {
	// Send request to make a charge
	return MakeChargeOutput{}, nil
}

func main() {
	Authenticate("my-key")
}
