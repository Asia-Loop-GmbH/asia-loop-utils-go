package db

type CustomizablePrice struct {
	Value        string            `bson:"value" json:"value"`
	CustomValues map[string]string `bson:"customValues" json:"customValues"`
}

func (p *CustomizablePrice) GetPrice(key string) string {
	custom, ok := p.CustomValues[key]
	if ok && custom != "0.00" {
		return custom
	}
	return p.Value
}
