package dto

type VueDictionary struct {
	DicNo  string       `json:"dicNo"`
	Config string       `json:"config"`
	Data   []DictDetail `json:"data"`
}
type DictDetail struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}
