package structs

type AddRequest struct {
	Name        string `json:"name"`
	Club        string `json:"club"`
	Nationality string `json:"nationality"`
}
