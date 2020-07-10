package schema

type Blockchain struct {
	Height       string `json:"height"`  // 区块高度
	Messages     string `json:"messages"`  // 上链信息数量
	transactions string `json:"transactions"` // 交易数量
	Nodes        string `json:"nodes"` // 节点数量
}

type VerifySwag struct {
	Hash string `json:"hash" example:"45d44ca55d"`
}

type UploadSwag struct {
	Type  string `json:"type" example:"p"`
	Point string `json:"point" example:"point001"`
	Raw   string `json:"raw" example:"{\"point\":\"point001\",\"type\":\"temperature\",\"point\":\"26.2\",\"unit\":\"C\"}"`
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
