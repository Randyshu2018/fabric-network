package transaction

import (
	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
	"encoding/json"
	"fmt"
	"bytes"
)

const (
	KEYSOTRE_PREFIX = "keystore"
	COMPOSITE_KEY = "uid~id"
)

// Get returns the value of the specified asset key
func (t *SimpleChainCode) saveKey(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	if len(args) != 2 {
		return shim.Error("Incorrect arguments. Expecting one")
	}
	// 获取键
	key := args[0]
	if len(key) == 0 {
		return shim.Error("Key can not be empty")
	}
	// 获取SecretKey对象
	secretKey := &Secretkey{}
	err :=json.Unmarshal([]byte(args[1]), secretKey)
	if err != nil {
		return shim.Error("Failed to convert keystore object , source : "+ args[1])
	}
	// 将SecretKey对象转换为JSON对象，并返回字节数组
	keyBytes ,err := json.Marshal(secretKey)
	if err != nil {
		return shim.Error("Failed to marshal keystore object , source : "+ string(keyBytes))
	}
	
	// 密钥上链
	err = stub.PutState(key,keyBytes)

	if err != nil {
		return shim.Error("Failed to save secretKey : "+ string(keyBytes))
	}

	//  ==== Index the marble to enable color-based range queries, e.g. return all blue marbles ====
	//  An 'index' is a normal key/value entry in state.
	//  The key is a composite key, with the elements that you want to range query on listed first.
	//  In our case, the composite key is based on indexName~color~name.
	//  This will enable very efficient state range queries based on composite keys matching indexName~color~*
	colorNameIndexKey, err := stub.CreateCompositeKey(COMPOSITE_KEY, []string{secretKey.UserID, secretKey.ID})
	if err != nil {
		return shim.Error(err.Error())
	}
	//  Save index entry to state. Only the key name is needed, no need to store a duplicate copy of the marble.
	//  Note - passing a 'nil' value will effectively delete the key from state, therefore we pass null character as value
	value := []byte{0x00}
	stub.PutState(colorNameIndexKey, value)
	
	// 返回交易ID
	return shim.Success([]byte(stub.GetTxID()))
}

func (t *SimpleChainCode) getKeyByRange(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	if len(args) < 2 {
		return shim.Error("Incorrect number of arguments. Expecting 2")
	}

	startKey := args[0]
	endKey := args[1]

	resultsIterator, err := stub.GetStateByRange(startKey, endKey)
	if err != nil {
		return shim.Error(err.Error())
	}
	defer resultsIterator.Close()

	buffer, err := constructQueryResponseFromIterator(resultsIterator)
	if err != nil {
		return shim.Error(err.Error())
	}

	fmt.Printf("- getKeyByRange queryResult:\n%s\n", buffer.String())

	return shim.Success(buffer.Bytes())
}

// ===========格式化查询结果集================================================================================
func constructQueryResponseFromIterator2(stub shim.ChaincodeStubInterface,resultsIterator shim.StateQueryIteratorInterface) (*bytes.Buffer, error) {
	var buffer bytes.Buffer
	buffer.WriteString("{")
	buffer.WriteString("\"data\":")

	buffer.WriteString("[")

	bArrayMemberAlreadyWritten := false
	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()
		// get the color and name from color~name composite key
		objectType, compositeKeyParts, err := stub.SplitCompositeKey(queryResponse.Key)
		if err != nil {
			return nil,err
		}
		uid := compositeKeyParts[0]
		id := compositeKeyParts[1]
		fmt.Printf("- found a key from index:%s uid:%s id:%s\n", objectType, uid, id)
		if err != nil {
			return nil, err
		}
		if bArrayMemberAlreadyWritten == true {
			buffer.WriteString(",")
		}
		buffer.WriteString("{\"key\":")
		buffer.WriteString("\"")
		buffer.WriteString(id)
		buffer.WriteString("\"")
		keyBytes,err := stub.GetState(id)
		if err != nil {
			return nil, err
		}
		buffer.WriteString(", \"record\":")
		buffer.WriteString(string(keyBytes))
		buffer.WriteString("}")
		bArrayMemberAlreadyWritten = true
	}
	buffer.WriteString("]")

	return &buffer, nil
}


// ==== Example: GetStateByPartialCompositeKey/RangeQuery =========================================
// transferMarblesBasedOnColor will transfer marbles of a given color to a certain new owner.
// Uses a GetStateByPartialCompositeKey (range query) against color~name 'index'.
// Committing peers will re-execute range queries to guarantee that result sets are stable
// between endorsement time and commit time. The transaction is invalidated by the
// committing peers if the result set has changed between endorsement time and commit time.
// Therefore, range queries are a safe option for performing update transactions based on query results.
// ===========================================================================================
func (t *SimpleChainCode) getKeyByPartialCompositeKey(stub shim.ChaincodeStubInterface, args []string) pb.Response {

	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}

	uid := args[0]
	fmt.Println("- start transferMarblesBasedOnColor ", uid)

	// Query the color~name index by color
	// This will execute a key range query on all keys starting with 'color'
	coloredMarbleResultsIterator, err := stub.GetStateByPartialCompositeKey(COMPOSITE_KEY, []string{uid})
	if err != nil {
		return shim.Error(err.Error())
	}
	defer coloredMarbleResultsIterator.Close()

	buffer, err := constructQueryResponseFromIterator2(stub, coloredMarbleResultsIterator)
	if err != nil {
		return shim.Error(err.Error())
	}

	fmt.Printf("- getKeyByRange queryResult:\n%s\n", buffer.String())

	return shim.Success(buffer.Bytes())
}
