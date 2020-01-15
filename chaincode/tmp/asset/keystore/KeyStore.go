package keystore

/* 密钥仓库 */
type KeyStore struct {
	/* 标识 */
	ID string `json:"id"`
	/* 用户标识 */
	UserID string `json:"uid"`
	/*
	 * 加密后的密钥：
	 * 1.生成密钥
	 * 2.用户公钥加密
	 * 3.BASE64加密
	 */
	EncryptedKey string  `json:"encryptedKey"`
}