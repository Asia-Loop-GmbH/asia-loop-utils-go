package servicecognito

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/cognitoidentityprovider"
	"github.com/aws/aws-sdk-go-v2/service/cognitoidentityprovider/types"
)

type CreateUserData struct {
	Username          string `json:"username"`
	TemporaryPassword string `json:"temporaryPassword"`
	CompanyKey        string `json:"companyKey"`
	StoreKey          string `json:"storeKey"`
	FirstName         string `json:"firstName"`
	LastName          string `json:"lastName"`
	Role              string `json:"role"`
}

func CreateUser(ctx context.Context, data *CreateUserData) error {
	c, err := getClient(ctx)
	if err != nil {
		return err
	}
	attrs := []types.AttributeType{
		{
			Name:  aws.String(cognitoUserAttributeCompany),
			Value: aws.String(data.CompanyKey),
		},
		{
			Name:  aws.String(cognitoUserAttributeStore),
			Value: aws.String(data.StoreKey),
		},
		{
			Name:  aws.String(cognitoUserAttributeEmail),
			Value: aws.String(data.Username),
		},
		{
			Name:  aws.String(cognitoUserAttributeEmailVerified),
			Value: aws.String("true"),
		},
		{
			Name:  aws.String(cognitoUserAttributeFirstName),
			Value: aws.String(data.FirstName),
		},
		{
			Name:  aws.String(cognitoUserAttributeLastName),
			Value: aws.String(data.LastName),
		},
		{
			Name:  aws.String(cognitoUserAttributeName),
			Value: aws.String(fmt.Sprintf("%s %s", data.FirstName, data.LastName)),
		},
	}
	_, err = c.AdminCreateUser(ctx, &cognitoidentityprovider.AdminCreateUserInput{
		UserPoolId:        aws.String(cognitoPool),
		Username:          aws.String(data.Username),
		TemporaryPassword: aws.String(data.TemporaryPassword),
		UserAttributes:    attrs,
	})
	if err != nil {
		return err
	}

	_, err = c.AdminAddUserToGroup(ctx, &cognitoidentityprovider.AdminAddUserToGroupInput{
		UserPoolId: aws.String(cognitoPool),
		Username:   aws.String(data.Username),
		GroupName:  aws.String(data.Role),
	})
	if err != nil {
		return err
	}

	return nil
}
