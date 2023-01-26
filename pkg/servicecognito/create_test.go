package servicecognito_test

import (
	"context"
	"log"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/asia-loop-gmbh/asia-loop-utils-go/v7/pkg/servicecognito"
)

func TestCreateUser(t *testing.T) {
	if testing.Short() {
		t.Skip()
	}
	ctx := context.TODO()

	err := servicecognito.DeleteUser(ctx, &servicecognito.DeleteUserData{
		Username: "lenamtruong+unitest@gmail.com",
	})
	if err != nil {
		log.Printf("it's ok, just a cleanup")
	}

	err = servicecognito.CreateUser(ctx, &servicecognito.CreateUserData{
		Username:          "lenamtruong+unitest@gmail.com",
		TemporaryPassword: "Toor***???1234",
		CompanyKey:        "SIEMENS",
		StoreKey:          "ERLANGEN",
		FirstName:         "Nam Truong",
		LastName:          "Le",
		Role:              "COMPANY",
	})
	assert.NoError(t, err)

	user, err := servicecognito.GetUser(ctx, &servicecognito.GetUserData{
		Username: "lenamtruong+unitest@gmail.com",
	})
	assert.NoError(t, err)
	assert.Equal(t, "lenamtruong+unitest@gmail.com", user.Username)
	assert.Equal(t, "Nam Truong", user.FirstName)
	assert.Equal(t, "Le", user.LastName)
	assert.Equal(t, "SIEMENS", user.Company)

	err = servicecognito.DeleteUser(ctx, &servicecognito.DeleteUserData{
		Username: "lenamtruong+unitest@gmail.com",
	})
	assert.NoError(t, err)

	_, err = servicecognito.GetUser(ctx, &servicecognito.GetUserData{
		Username: "lenamtruong+unitest@gmail.com",
	})
	assert.NoError(t, err)
}
