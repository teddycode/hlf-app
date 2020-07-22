package schema

type SensorSwag struct {
	StarTime string `json:"star_time" example:"1594382265"`
	EndTime  string `json:"end_time" example:"1595382265"`
	Point    string `json:"point" example:"19372180"`
	PageSize string `json:"page_size" example:"10"` // 页大小
	BookMark string `json:"book_mark" example:"s~19372180~1594382265"`  // 书签bookmark
}

type PicSwag struct {
	StarTime string `json:"star_time" example:"1594382265"`
	EndTime  string `json:"end_time" example:"1595382265"`
	Point    string `json:"point" example:"0018DE743E31"`
	PageSize string `json:"page_size" example:"10"` // 页大小
	BookMark string `json:"book_mark" example:"p~0018DE743E31~1594382265"`  // 书签bookmark
}

type FarmSwag struct {
	StarTime string `json:"star_time" example:"1594382265"`
	EndTime  string `json:"end_time" example:"9999999999"`
	User     string `json:"user" example:"admin"`
	PageSize string `json:"page_size" example:"10"` // 页大小
	BookMark string `json:"book_mark" example:"f~admin~1594382265"`  // 书签bookmark
}
