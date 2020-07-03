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
