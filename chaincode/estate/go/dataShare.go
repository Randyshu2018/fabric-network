package main

type DataShare struct {
	/**
     * 主键
     */
	ID string `json:"id"`
	
	/**
	 * 提案人
	 */
	Proposer string `json:"proposer"`
	
	/**
	 * 提案人公钥
	 */
	ProposerPublicKey string `json:"proposerPublicKey"`

	/**
	 * 请求授权数据对应的键
	 */
	TargetKey string `json:"targetKey"`

	/**
	 * 状态
	 * 0 申请 
	 * 1:同意 
	 * 2:拒绝
	 */
	Status string `json:"status"`

	/**
	 * 提案人公钥加密的对称密钥
	 */
	EncryptedKey string `json:"encryptedKey"`

	/**
	 * 受理人
	 */
	Approver string `json:"approver"`

	/**
	 * targetKey所在的资产类型
	 */
	AssetType string `json:"assetType"`
} 