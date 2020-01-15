package main

type BaseModel struct {
	// 公钥加密后的对称密钥
	EncryptedKey string  `json:"encryptedKey"`
	// JSON数据中需加密的属性
	CryptoFields  []string  `json:"cryptoFields"`
	// JSON数据摘要
	Digest string `json:"digest"`
	// 摘要签名
	Signature string `json:"signature"`
	// 资产类型
	AssetType string `json:"assetType"`
	// 键
	Key string `json:"key"`
	// 提案人
	Proposer string `json:"proposer"`
}