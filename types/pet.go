package types

type GetPet struct {
	Pid        string `json:"petid"`
	Uid        string `json:"number"`
	PetName    string `json:"name"`
	Imag       string `json:"imgurl"`
	Years      int    `json:"years"`
	Cid        string `json:"catchid"`
	Fid        string `json:"father"`
	Status     int    `json:"status"`
	CreateTime int64  `json:"birth"`
}
type RspGetPets struct {
	Total int      `json:"total"`
	Pets  []GetPet `json:"data"`
}
