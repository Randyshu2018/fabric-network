package main

import (
	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
)

/**
 * 数据上链
 */
func (t *SimpleChainCode) put(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	if len(args) != 1 {
		return shim.Error("Incorrect arguments. Expecting one")
	}

	err := stub.PutState(args[0],[]byte(args[1]))
	if err != nil {
		return shim.Error("Failed to get key : "+ args[0])
	}
	
	return shim.Success(nil)
}