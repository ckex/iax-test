package models

type  Result struct {
	Ads     []Ads `json:"ads"`
	Success bool `json:"success"`
}

type Ads struct {
	BidId    string `json:"bidId"`
	D        string `json:"d"`
	Eid      string `json:"eid"`
	Height   int `json:"height"`
	Id       string `json:"id"`
	ImpId    string `json:"impId"`
	Link     string `json:"link"`
	Noclick  bool `json:"noclick"`
	P        int `json:"p"`
	Position string `json:"position"`
	PubId    int `json:"pubid"`
	Rt       int `json:"rt"`
	Sts      string `json:"sts"`
	Width    int `json:"width"`
}
