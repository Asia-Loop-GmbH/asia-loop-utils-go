package woo

type ProductCategory struct {
	TermID         int    `json:"term_id"`
	Name           string `json:"name"`
	TermTaxonomyID int    `json:"term_taxonomy_id"`
	Taxonomy       string `json:"taxonomy"`
	Description    string `json:"description"`
}
