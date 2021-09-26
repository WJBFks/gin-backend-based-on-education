package main

import (
	"bytes"
	"encoding/json"
	"fmt"

	"github.com/hyperledger/fabric-chaincode-go/shim"
	"github.com/hyperledger/fabric-protos-go/peer"
)

type Chaincode struct {
}

func (t *Chaincode) Init(stub shim.ChaincodeStubInterface) peer.Response {
	fmt.Println(" ==== Init ====")

	return shim.Success(nil)
}

func (t *Chaincode) Invoke(stub shim.ChaincodeStubInterface) peer.Response {
	// 获取用户意图
	fun, args := stub.GetFunctionAndParameters()

	if fun == "TestAdd" {
		return t.TestAdd(stub, args)
	} else if fun == "TestQueryByName" {
		return t.TestQueryByName(stub, args)
	}

	return shim.Error("指定的函数名称错误")

}

type Test struct {
	Name   string `json:"Name"`   // 姓名
	Gender string `json:"Gender"` // 性别
	Age    string `json:"Age"`    // age
}

func (t *Chaincode) TestAdd(stub shim.ChaincodeStubInterface, args []string) peer.Response {
	if len(args) != 1 {
		return shim.Error("给定的参数个数不符合要求")
	}

	var tData Test

	err := json.Unmarshal([]byte(args[0]), &tData)
	if err != nil {
		return shim.Error("反序列化信息时发生错误")
	}

	_, bl := TestPut(stub, tData)
	if !bl {
		return shim.Error("保存信息时发生错误")
	}

	// err = stub.SetEvent(args[1], []byte{})
	// if err != nil {
	// 	return shim.Error(err.Error())
	// }

	return shim.Success([]byte("信息添加成功"))
}

func TestPut(stub shim.ChaincodeStubInterface, td Test) ([]byte, bool) {

	b, err := json.Marshal(td)
	if err != nil {
		return nil, false
	}

	fmt.Println(td.Name)

	// 保存edu状态
	err = stub.PutState(td.Name, b)
	if err != nil {
		return nil, false
	}

	return b, true
}

func TestFindByName(stub shim.ChaincodeStubInterface, name string) (Test, bool) {
	var td Test
	// 根据身份证号码查询信息状态
	b, err := stub.GetState(name)
	if err != nil {
		return td, false
	}

	if b == nil {
		return td, false
	}

	// 对查询到的状态进行反序列化
	err = json.Unmarshal(b, &td)
	if err != nil {
		return td, false
	}

	// 返回结果
	return td, true
}

func (t *Chaincode) TestQueryByName(stub shim.ChaincodeStubInterface, args []string) peer.Response {

	if len(args) != 1 {
		return shim.Error("给定的参数个数不符合要求")
	}
	Name := args[0]

	// 拼装CouchDB所需要的查询字符串(是标准的一个JSON串)
	// queryString := fmt.Sprintf("{\"selector\":{\"docType\":\"eduObj\", \"CertNo\":\"%s\"}}", CertNo)
	queryString := fmt.Sprintf("{\"selector\":{\"Name\":\"%s\"}}", Name)

	// 查询数据
	result, err := getEduByQueryString(stub, queryString)
	if err != nil {
		return shim.Error("Query Error")
	}
	if result == nil {
		return shim.Error("No Infor")
	}
	return shim.Success(result)
}

func getEduByQueryString(stub shim.ChaincodeStubInterface, queryString string) ([]byte, error) {

	resultsIterator, err := stub.GetQueryResult(queryString)
	if err != nil {
		return nil, err
	}
	defer resultsIterator.Close()

	// buffer is a JSON array containing QueryRecords
	var buffer bytes.Buffer

	bArrayMemberAlreadyWritten := false
	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()
		if err != nil {
			return nil, err
		}
		// Add a comma before array members, suppress it for the first array member
		if bArrayMemberAlreadyWritten == true {
			buffer.WriteString(",")
		}

		// Record is a JSON object, so we write as-is
		buffer.WriteString(string(queryResponse.Value))
		bArrayMemberAlreadyWritten = true
	}

	fmt.Printf("- getQueryResultForQueryString queryResult:\n%s\n", buffer.String())

	return buffer.Bytes(), nil

}

func main() {
	err := shim.Start(new(Chaincode))
	if err != nil {
		fmt.Printf("启动Chaincode时发生错误: %s", err)
	}
}
