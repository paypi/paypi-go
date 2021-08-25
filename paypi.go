package paypi

import (
	"errors"
	"fmt"

	"github.com/paypi/paypi-go/gql"
)

var (
	//Key is the API secret Key
	Key       string
	GqlClient gql.GqlClient
)

func init() {
	GqlClient = gql.New("https://api.paypi.dev/graphql")
}

func SetConnection(url string) {
	GqlClient = gql.New(url)
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
	if Key == "" {
		return AuthenticatedOutput{}, errors.New("paypi.Key is not set")
	}

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
	UnitsUsed        int32
}

type MakeChargeOutput struct {
	Success bool
}

// MakeCharge requests to take payments from the user, totaling the cost of the
// given ChargeIdentifier.
func (a AuthenticatedOutput) MakeCharge(input MakeChargeInput) (MakeChargeOutput, error) {
	resp := MakeChargeResponse{}

	var unitsUsed int32
	if input.UnitsUsed != 0 {
		unitsUsed = input.UnitsUsed
	}
	err := GqlClient.MakeRequest(gql.GqlQuery{
		Query: `
			mutation MakeCharge($chargeIdent: String!, $subSecret: String!, $unitsUsed: Int) {
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
