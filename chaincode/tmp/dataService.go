package main

import (
	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
	"encoding/json"
	"fmt"
	"bytes"
	"strconv"
)
// 组合键
const (
	COMPOSITE_KEY_ASSET_TYPE_PROPOSER_KEY = "assettype~proposer~key"
	COMPOSITE_KEY_ASSET_TYPE_PROPOSER_STATUS_KEY = "assettype~proposer~status~key"
	COMPOSITE_KEY_ASSET_TYPE_APPROVER_STATUS_KEY = "assettype~approver~status~key"
	COMPOSITE_KEY_ASSET_TYPE_PROPOSER_STATUS_TARGET_KEY = "assettype~proposer~status~targetkey"
)

/**
 * 资产上链
 */
func (t *SimpleChainCode) saveAsset(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	// 获取DataModel对象
	model := &AssetModel{}
	err :=json.Unmarshal([]byte(args[0]), model)
	if err != nil {
		return shim.Error("Failed to convert dataModel , source : "+ args[0])
	}

	// 参数校验
	err = CheckBaseModel(&model.BaseModel)
	if err != nil {
		return shim.Error(err.Error())
	}
	
	// 将DataModel对象转换为JSON对象，并返回字节数组
	modelBytes ,err := json.Marshal(model)
	if err != nil {
		return shim.Error("Failed to marshal DataModel , source : "+ args[0])
	}
	
	// 数据上链
	err = stub.PutState(model.Key,modelBytes)

	if err != nil {
		return shim.Error("Failed to save DataModel : "+ string(modelBytes))
	}

	// 创建索引
	indexKey, err := stub.CreateCompositeKey(COMPOSITE_KEY_ASSET_TYPE_PROPOSER_KEY, []string{model.AssetType,model.Proposer,model.Key})
	if err != nil {
		return shim.Error(err.Error())
	}
	value := []byte{0x00}
	// 索引上链
	// 这里没必要再做一份数据拷贝,直接传了nil的字节数组
	stub.PutState(indexKey, value)
	
	// 返回交易ID
	return shim.Success([]byte(stub.GetTxID()))
}

/**
 * 数据申请信息上链
 */
func (t *SimpleChainCode) saveDataShare(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	
	// 获取DataModel对象
	dataModel := &DataShareModel{}
	err :=json.Unmarshal([]byte(args[0]), dataModel)
	if err != nil {
		return shim.Error("Failed to convert dataModel , source : "+ args[0])
	}

	// 参数校验
	err = CheckBaseModel(&dataModel.BaseModel)
	if err != nil {
		return shim.Error(err.Error())
	}
	
	//解析上链数据
	//var dataMap map[string]string
	dataShare := dataModel.Data
	
	// 参数校验
	err = CheckShareData(&dataShare)
	if err != nil {
		return shim.Error(err.Error())
	}

	// 初次申请需要做重复申请校验
	if dataShare.TargetKey == "0"{
		// 重复申请检测
		duplicateKey, _ := stub.CreateCompositeKey(COMPOSITE_KEY_ASSET_TYPE_PROPOSER_STATUS_TARGET_KEY, []string{dataModel.AssetType,dataModel.Proposer,"0",dataShare.TargetKey})
		duplicateKeyBytes , _ := stub.GetState(duplicateKey)
		if duplicateKeyBytes != nil {
			return shim.Error("Exist unhandled apply")
		}
	}else{
		// 申请受理后 删除原组合键
		indexKey, _ := stub.CreateCompositeKey(COMPOSITE_KEY_ASSET_TYPE_PROPOSER_STATUS_KEY, []string{dataModel.AssetType,dataModel.Proposer,"0",dataModel.Key})
		indexKey2, _ := stub.CreateCompositeKey(COMPOSITE_KEY_ASSET_TYPE_APPROVER_STATUS_KEY, []string{dataModel.AssetType,dataShare.Approver,"0",dataModel.Key})
		indexKey3, _ := stub.CreateCompositeKey(COMPOSITE_KEY_ASSET_TYPE_PROPOSER_STATUS_TARGET_KEY, []string{dataModel.AssetType,dataModel.Proposer,"0",dataShare.TargetKey})
		stub.DelState(indexKey)
		stub.DelState(indexKey2)
		stub.DelState(indexKey3)
	}
	
	// 创建索引
	// 这里根据提案人和受理人分别创建了索引
	indexKey, err := stub.CreateCompositeKey(COMPOSITE_KEY_ASSET_TYPE_PROPOSER_STATUS_KEY, []string{dataModel.AssetType,dataModel.Proposer,dataShare.Status,dataModel.Key})
	indexKey2, err := stub.CreateCompositeKey(COMPOSITE_KEY_ASSET_TYPE_APPROVER_STATUS_KEY, []string{dataModel.AssetType,dataShare.Approver,dataShare.Status,dataModel.Key})
	indexKey3, err := stub.CreateCompositeKey(COMPOSITE_KEY_ASSET_TYPE_PROPOSER_STATUS_TARGET_KEY, []string{dataModel.AssetType,dataModel.Proposer,dataShare.Status,dataShare.TargetKey})
	if err != nil {
		return shim.Error(err.Error())
	}
	// 索引上链
	// 这里没必要传DataModel的拷贝
	value := []byte{0x00}
	stub.PutState(indexKey, value)
	stub.PutState(indexKey2, value)
	stub.PutState(indexKey3, value)

	// 将DataModel对象转换为JSON对象，并返回字节数组
	modelBytes ,err := json.Marshal(dataModel)
	if err != nil {
		return shim.Error("Failed to marshal DataModel object , source : "+ string(modelBytes))
	}

	// 数据上链
	err = stub.PutState(dataModel.Key,modelBytes)

	if err != nil {
		return shim.Error("Failed to save DataModel : "+ string(modelBytes))
	}
	// 返回交易ID
	return shim.Success([]byte(stub.GetTxID()))
}

/**
 * 格式化查询结果集
 */
func constructQueryResponseFromIterator(stub shim.ChaincodeStubInterface,resultsIterator shim.StateQueryIteratorInterface, index int) (*bytes.Buffer, error) {
	var buffer bytes.Buffer

	buffer.WriteString("[")

	bArrayMemberAlreadyWritten := false
	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()
		// 组合键拆分
		objectType, compositeKeyParts, err := stub.SplitCompositeKey(queryResponse.Key)
		if err != nil {
			return nil,err
		}
		key := compositeKeyParts[index]
		fmt.Printf("- found a key from index:%s id:%s\n", objectType, index)
		if err != nil {
			return nil, err
		}
		if bArrayMemberAlreadyWritten == true {
			buffer.WriteString(",")
		}
		keyBytes,err := stub.GetState(key)
		if err != nil {
			return nil, err
		}
		buffer.WriteString(string(keyBytes))
		bArrayMemberAlreadyWritten = true
	}
	buffer.WriteString("]")

	return &buffer, nil
}

/**
 * 部分复合键查询
 */
func (t *SimpleChainCode) getKeyByPartialCompositeKey(stub shim.ChaincodeStubInterface, args []string) pb.Response {

	// 组合键
	compositeKey := args[0]
	// 数组索引
	index,err:=strconv.Atoi(args[len(args)-1])
	if err != nil {
		return shim.Error(err.Error())
	}

	fmt.Printf("- compositeKey:%s index:%s\n", compositeKey, index)

	// 基于部分组合键查询
	resultsIterator, err := stub.GetStateByPartialCompositeKey(args[0], args[1:len(args)-1])
	if err != nil {
		return shim.Error(err.Error())
	}
	defer resultsIterator.Close()

	buffer, err := constructQueryResponseFromIterator(stub, resultsIterator, index)
	if err != nil {
		return shim.Error(err.Error())
	}

	fmt.Printf("- getKeyByPartialCompositeKey queryResult:\n%s\n", buffer.String())

	return shim.Success(buffer.Bytes())
}