package servicemailjet

type SendOrderVariables struct {
	FirstName     string `json:"firstName"`
	Title         string `json:"title"`
	Content       string `json:"content"`
	ActionText    string `json:"actionText"`
	ActionLink    string `json:"actionLink"`
	ActionEnabled string `json:"actionEnabled"`
}
