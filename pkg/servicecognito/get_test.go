package servicecognito_test

import (
	"context"
	"log"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/asia-loop-gmbh/asia-loop-utils-go/pkg/servicecognito"
	commoncontext "github.com/nam-truong-le/lambda-utils-go/v2/pkg/context"
)

func TestGetUser(t *testing.T) {
	if testing.Short() {
		t.Skip()
	}
	ctx := context.WithValue(context.TODO(), commoncontext.FieldStage, "dev")
	user, err := servicecognito.GetUser(ctx, &servicecognito.GetUserData{
		Username: "lenamtruong@gmail.com",
	})
	assert.NoError(t, err)
	assert.Equal(t, "lenamtruong@gmail.com", user.Username)
	log.Printf("%v", user)
}
