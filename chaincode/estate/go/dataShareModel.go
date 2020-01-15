package main

// 数据模型
type DataShareModel struct {
	BaseModel
	// 数据申请
	Data DataShare `json:"data"`
}