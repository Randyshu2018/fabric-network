package main

import (
	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
)

/**
 * 删除数据
 */
func (t *SimpleChainCode) delete(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	if len(args) != 1 {
		return shim.Error("Incorrect arguments. Expecting one")
	}

	err := stub.DelState(args[0])
	if err != nil {
		return shim.Error("Failed to delete key : "+ args[0])
	}

	return shim.Success(nil)
}