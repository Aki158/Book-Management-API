package structs

type Response struct {
	StudentArr []Student `json:"students"`
	TotalCount int `json:"totalCount"`
}

type Student struct {
	Id int `json:"id"`
	Name string `json:"name"`
	LoginId string `json:"loginId"`
	ClassRoom Class `json:"classroom"`
}

type Class struct {
	Id int `json:"id"`
	Name string `json:"name"`
}
