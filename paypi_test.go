package paypi_test

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/paypi/paypi-go/gql/mocks"

	"github.com/paypi/paypi-go"
)

func TestAuthenticate(t *testing.T) {
	t.Run("check authed requests are allowed", func(t *testing.T) {
		gqlMock := mocks.GqlClient{}

		gqlMock.On("MakeRequest", mock.Anything, mock.Anything).Return(nil).Run(func(args mock.Arguments) {
			arg := args.Get(1).(*paypi.CheckSubscriberSecretResponse)
			arg.CheckSubscriberSecret.IsAuthed = true
		})

		paypi.Key = "TestKey"
		paypi.GqlClient = &gqlMock

		out, err := paypi.Authenticate("1234")

		assert.Nil(t, err)
		assert.Equal(t, out.ClientToken, "1234")

	})

	t.Run("check unauthed requests error", func(t *testing.T) {
		gqlMock := mocks.GqlClient{}

		var outErr error = errors.New("An error from backend")
		gqlMock.On("MakeRequest", mock.Anything, mock.Anything).Return(outErr).Run(func(args mock.Arguments) {
			arg := args.Get(1).(*paypi.CheckSubscriberSecretResponse)
			arg.CheckSubscriberSecret.IsAuthed = false
		})

		paypi.Key = "TestKey"
		paypi.GqlClient = &gqlMock

		out, err := paypi.Authenticate("12345")

		assert.Equal(t, err, outErr)
		assert.Equal(t, out.ClientToken, "")
	})
}
