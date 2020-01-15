package main

import (
	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
)

/**
 * 添加私有数据
 */
func (t *SimpleChainCode) addPrivateData(stub shim.ChaincodeStubInterface, args []string) pb.Response {

	transMap, err := stub.GetTransient()
	if err != nil {
		return shim.Error("Error getting transient: " + err.Error())
	}

	if _, ok := transMap["key"]; !ok {
		return shim.Error("argument key is missing in the transient map")
	}
	if _, ok := transMap["value"]; !ok {
		return shim.Error("argument value is missing in the transient map")
	}
	
	key,_ := transMap["key"]
	value,_ := transMap["value"]

	// 数据上链
	err = stub.PutPrivateData("collectionDetails",string(key),[]byte(value))

	if err != nil {
		return shim.Error("Failed to save collectionDetails : "+ string(value))
	}
	// 返回交易ID
	return shim.Success([]byte(stub.GetTxID()))
}

/**
 * 查询私有数据
 */
func (t *SimpleChainCode) getPrivateData(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	key:=args[0]

	value, err := stub.GetPrivateData("collectionDetails",key)
	if err != nil {
		return shim.Error("Failed to get key : "+ args[0])
	}
	if value == nil {
		return shim.Error(args[0] +" not found")
	}
	return shim.Success(value)
}

func (t *SimpleChainCode) invokeWithoutPutState(stub shim.ChaincodeStubInterface,args []string) pb.Response {
	stub.SetEvent("test-event",[]byte("hello event!"))
	return shim.Success(nil)
}

