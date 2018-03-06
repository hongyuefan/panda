package types

type GetPet struct {
	UId string `json:"number"`
}
type RspGetPets struct {
	Total int       `json:"total"`
	Pets  []*GetPet `json:"data"`
}
