package schema

type Blockchain struct {
	Height       string `json:"height"`       // 区块高度
	Messages     string `json:"messages"`     // 上链信息数量
	transactions string `json:"transactions"` // 交易数量
	Nodes        string `json:"nodes"`        // 节点数量
}

type VerifySwag struct {
	Hash string `json:"hash" example:"b9c52e66c1ebfc826e324a394a106f9dc9550fed4390808b2d8932ff91c92b5a"`
}

type UploadSwag struct {
	Type string `json:"type" example:"p"`
	Raw  string `json:"raw" example:"{\"point\":\"point001\",\"type\":\"temperature\",\"value\":\"26.2\",\"unit\":\"C\"}"`
}

type BCSensor struct {
	Point string `json:"point" example:"point001"`
	Type  string `json:"type"  example:"temperature"`
	Value string `json:"value" example:"26.3"`
	Unit  string `json:"unit"  example:"℃"`
}

type BCPic struct {
	Point string `json:"point"  example:"point001"`
	Name  string `json:"name"  example:"b9c52e66c1ebfc826e324a394a106f9dc9550fed4390808b2d8932ff91c92b5a"`
	Size  string `json:"size" example:"1024"`
	Type  string `json:"type" example:"sensor"` // 0: sensor, 1:pic,
}

//type BCFarmSwag struct {
//	Name string `json:"name"`
//
//}
