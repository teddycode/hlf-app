package schema

type Blockchain struct {
	Height       string `json:"height"`
	Messages     string `json:"messages"`
	transactions string `json:"transactions"`
	Nodes        string `json:"nodes"`
}

type VerifySwag struct {
	Hash string `json:"hash"`
}

type UploadSwag struct {
	Type  string `json:"type"`
	Point string `json:"point"`
	Raw   string `json:"raw"`
}

type BCSensor struct {
	Point string `json:"point"`
	Type  string `json:"type"`
	Value string `json:"value"`
	Unit  string `json:"unit"`
}

type BCPic struct {
	Point string `json:"point"`
	Name  string `json:"name"`
	Size  string `json:"size"`
	Type  string `json:"type"` // 0: sensor, 1:pic,
}

//type BCFarmSwag struct {
//	Name string `json:"name"`
//
//}
