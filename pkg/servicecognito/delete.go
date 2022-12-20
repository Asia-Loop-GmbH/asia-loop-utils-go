package servicecognito

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/cognitoidentityprovider"
)

type DeleteUserData struct {
	Username string `json:"username"`
}

func DeleteUser(ctx context.Context, data *DeleteUserData) error {
	c, err := getClient(ctx)
	if err != nil {
		return err
	}
	_, err = c.AdminDeleteUser(ctx, &cognitoidentityprovider.AdminDeleteUserInput{
		UserPoolId: aws.String(cognitoPool),
		Username:   aws.String(data.Username),
	})
	if err != nil {
		return err
	}
	return nil
}
