package servicewoo

type MetaData struct {
	Key          string      `json:"key"`
	Value        interface{} `json:"value"`
	DisplayKey   string      `json:"display_key"`
	DisplayValue string      `json:"display_value"`
}
