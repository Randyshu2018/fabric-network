package main

import "errors"

/**
 * 参数合法性校验
 */
func CheckBaseModel(model *BaseModel) error {
	if model.AssetType == "" {
		return  errors.New("argument assetType is missing")
	}
	if model.Proposer == "" {
		return errors.New("argument proposer is missing")
	}
	if model.Key == "" {
		return errors.New("argument key is missing")
	}
	return nil
}

/**
 * 参数合法性校验
 */
func CheckShareData(model *DataShare) error {
	if model.AssetType == "" {
		return errors.New("argument dataShare.assetType is missing")
	}
	if model.Proposer == "" {
		return errors.New("argument proposer is missing")
	}
	if model.ProposerPublicKey == "" {
		return errors.New("argument proposerPublicKey is missing")
	}
	if model.TargetKey == "" {
		return errors.New("argument targetKey is missing")
	}
	if model.Status == "" {
		return errors.New("argument status is missing")
	}
	if model.Status != "0" && model.Status != "1" && model.Status != "2" {
		return errors.New("Invalid status. Expecting \"0\" \"1\" \"2\" ")
	}
	if model.Approver == "" {
		return errors.New("argument approver is missing")
	}
	return nil
}
