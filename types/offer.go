package types

type RspTransPet struct {
	TxHash string `json:"txhash"`
}

type RspAddOffer struct {
	OfferId string `json:"offerId"`
}

type RspUpdateOffer struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}
type GetOffers struct {
	Id         int64  `json:"OfferId"`
	Years      int    `json:"years"`
	Pid        string `json:"petId"`
	Uid        string `json:"memberId"`
	Price      string `json:"price"`
	Status     int    `json:"status"`
	CreateTime int64  `json:"create_time"`
	UpdateTime int64  `json:"update_time"`
}
type RspGetOffers struct {
	Total  int64       `json:"total"`
	Offers []GetOffers `json:"offers"`
}
