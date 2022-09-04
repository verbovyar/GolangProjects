package structs

type UpdateRequest struct {
	Name        string `json:"name"`
	Club        string `json:"club"`
	Nationality string `json:"nationality"`
	Id          int32  `json:"id"`
}
