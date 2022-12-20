package servicecognito

const (
	cognitoUserAttributeFirstName     = "given_name"
	cognitoUserAttributeLastName      = "family_name"
	cognitoUserAttributeName          = "name"
	cognitoUserAttributeEmailVerified = "email_verified"
	cognitoUserAttributeEmail         = "email"
	cognitoUserAttributeStore         = "custom:store"
	cognitoUserAttributeCompany       = "custom:company"

	cognitoPool = "eu-central-1_P9Ik6952f"
)

type User struct {
	Username  string `json:"username"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Status    string `json:"status"`
	Company   string `json:"company"`
}
