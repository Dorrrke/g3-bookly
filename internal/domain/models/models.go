package models

type User struct {
	UID   string `json:"uuid,omitempty"`
	Email string `json:"email" validate:"required,email"`
	Pass  string `json:"pass" validate:"required,min=8"`
	Age   int    `json:"age" validate:"required,gte=16"`
}

type Book struct {
	BID    string `json:"bid,omitempty"`
	Lable  string `json:"lable" validate:"required,min=3"`
	Author string `json:"author" validate:"required,min=5"`
	Desc   string `json:"desc" validate:"required,min=10"`
	Age    int    `json:"age" validate:"required"`
	Count  int    `json:"count,omitempty"`
}
