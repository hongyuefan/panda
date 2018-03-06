package types

type GetPet struct {
	Uid           string `json:"memberid"`
	Pid           string `json:"number"`
	PetName       string `json:"name"`
	Imag          string `json:"imgurl"`
	Years         int    `json:"years"`
	Cid           string `json:"catchid"`
	Fid           string `json:"father"`
	Status        int    `json:"status"`
	TrainTotal    string `json:"train_total"`
	LastCatchTime int64  `json:"last_catch_time"`
	CreateTime    int64  `json:"birth"`
	CatchTimes    int    `json:"catch_times"`
	IsRare        int    `json:"rare"`
}
type RspGetPets struct {
	Total int      `json:"total"`
	Pets  []GetPet `json:"data"`
}

type GetPetAttr struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}
type RspGetPetAttr struct {
	Uid   string       `json:"memberid"`
	Pid   string       `json:"number"`
	Years int          `json:"years"`
	Attrs []GetPetAttr `json:"attributes"`
}
