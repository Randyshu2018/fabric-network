package main

// 数据模型
type AssetModel struct {
	BaseModel
	// JSON数据 这里默认使用string类型，可以按需调整
	Data string `json:"data"`
}