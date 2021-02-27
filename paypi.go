package paypi

import (
	"fmt"

	"github.com/Paypi/paypi-go/gql"
)

var (
	//Key is the API secret Key
	Key string
	// GqlClient is the graphql client to use
	GqlClient gql.GqlClient
)

func init() {
	GqlClient = gql.New("http://localhost:8080/graphql")
}

type CheckSubscriberSecretResponse struct {
	CheckSubscriberSecret struct {
		IsAuthed bool
	}
}

type AuthenticatedOutput struct {
	ClientToken string
}

// Authenticate checks if a given api token is authorized.
// This checks if the token is not blocked in any way and has
// card details setup on their account to make payments.
func Authenticate(clientToken string) (AuthenticatedOutput, error) {
	// Send request to API to authenticate clientToken

	resp := CheckSubscriberSecretResponse{}

	err := GqlClient.MakeRequest(gql.GqlQuery{
		Query: `
			mutation AuthenticateClient($serviceSecret: String!, $subSecret: String!) { 
					checkSubscriberSecret(input: {
						serviceSecret: $serviceSecret
						subscriptionSecret: $subSecret
					}) {
						isAuthed
					}
			}
		`,
		Variables: map[string]interface{}{
			"serviceSecret": Key,
			"subSecret":     clientToken,
		},
	}, &resp)
	if err != nil {
		return AuthenticatedOutput{}, err
	}

	if resp.CheckSubscriberSecret.IsAuthed {
		return AuthenticatedOutput{
			ClientToken: clientToken,
		}, nil
	}

	return AuthenticatedOutput{}, ErrInvalidToken
}

type MakeChargeResponse struct {
	MakeCharge struct {
		Success bool
	}
}

type MakeChargeInput struct {
	ChargeIdentifier string
	UnitsUsed        float64
}

type MakeChargeOutput struct {
	Success bool
}

// MakeCharge requests to take payments from the user, totaling the cost of the
// given ChargeIdentifier.
func (a AuthenticatedOutput) MakeCharge(input MakeChargeInput) (MakeChargeOutput, error) {
	resp := MakeChargeResponse{}

	var unitsUsed float64
	if input.UnitsUsed != 0 {
		unitsUsed = input.UnitsUsed
	}
	err := GqlClient.MakeRequest(gql.GqlQuery{
		Query: `
			mutation MakeCharge($chargeIdent: String!, $subSecret: String!, $unitsUsed: Float) {
					makeCharge(input: {
						chargeIdentifier: $chargeIdent
						subscriptionSecret: $subSecret
						unitsUsed: $unitsUsed
					}) {
						success
					}
			}
		`,
		Variables: map[string]interface{}{
			"chargeIdent": input.ChargeIdentifier,
			"subSecret":   a.ClientToken,
			"unitsUsed":   unitsUsed,
		},
	}, &resp)
	if err != nil {
		fmt.Print(err.Error())
		return MakeChargeOutput{Success: false}, err
	}

	if resp.MakeCharge.Success {
		return MakeChargeOutput{Success: true}, nil
	}

	return MakeChargeOutput{Success: false}, ErrUnableToMakeCharge
}
