package servicecognito

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/cognitoidentityprovider"
)

type GetUserData struct {
	Username string `json:"username"`
}

func GetUser(ctx context.Context, data *GetUserData) (*User, error) {
	c, err := getClient(ctx)
	if err != nil {
		return nil, err
	}
	out, err := c.AdminGetUser(ctx, &cognitoidentityprovider.AdminGetUserInput{
		Username:   aws.String(data.Username),
		UserPoolId: aws.String(cognitoPool),
	})
	if err != nil {
		return nil, err
	}
	return &User{
		Username:  *out.Username,
		FirstName: getCognitoUserAttribute(out, cognitoUserAttributeFirstName),
		LastName:  getCognitoUserAttribute(out, cognitoUserAttributeLastName),
		Status:    string(out.UserStatus),
		Company:   getCognitoUserAttribute(out, cognitoUserAttributeCompany),
	}, nil
}

func getCognitoUserAttribute(user *cognitoidentityprovider.AdminGetUserOutput, attr string) string {
	for _, item := range user.UserAttributes {
		if *item.Name == attr {
			return *item.Value
		}
	}
	return ""
}
