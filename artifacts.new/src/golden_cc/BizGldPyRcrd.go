package main
import (
	"github.com/hyperledger/fabric/core/chaincode/shim"
	"fmt"
//	"bytes"
	"strings"
	pb "github.com/hyperledger/fabric/protos/peer"
	"encoding/json"
	"reflect"
	"time"
	"strconv"
	"bytes"
)

/**************************************************************************************/
// 金票支付记录结构体
type BizGldPyRcrd struct{
	ObjectType string     `json:"docType"`             // 类型
	Id         int64      `json:"id"`                  // 主键
	CreateTime int64      `json:"createTime"`          // 创建时间
	UpdateTime int64      `json:"updateTime"`          // 更新时间
	CreateUser string     `json:"createUser"`          // 创建人
	UpdateUser string     `json:"updateUser"`          // 更新人
	ExpdId     string     `json:"expdId"`              // 扩展ID
	DelInd     string     `json:"delInd"`              // 删除标志
	Version    int32      `json:"version"`             // 版本号
	TenantId   string     `json:"tenantId"`            // 租户ID

	OriGldId         string    `json:"oriGldId"`   // 原金票编号
	OriGldBal        float64   `json:"oriGldBal"`  // 原金票余额
	Pypt             string    `json:"pypt"`       // 支付方
	RcPty            string    `json:"rcPty"`      // 收票人
	PyAmt            float64   `json:"pyAmt"`      // 支付金额
	PyTm             int64     `json:"pyTm"`       // 支付时间
	NewGldId         string    `json:"newGldId"`   // 新金票编号
}

// 用于解析queryString
type QryStrBizGldPyRcrd struct{
	Selector          BizGldPyRcrd         `json:"selector"`
}

// 判断BizGldPyRcrd变量是否为空
func isEmptyBizGldPyRcrd(arg BizGldPyRcrd) bool {
	value := reflect.ValueOf(arg)
	num := 0
	for i:=0;i<value.NumField();i++{
		if isEmpty( value.Field(i).Interface() ){
			num++
		}
	}

	if num == value.NumField(){
		return true
	}else{
		return false
	}
}

// BizGldPyRcrd变量变为queryString
func tranfBizGldPyRcrdToQryStr(arg BizGldPyRcrd)( string){
	queryString := "{\"selector\":{\"docType\":\"BizGldPyRcrd\""                      // 封装，头部

	value := reflect.ValueOf(arg)
	typ := reflect.TypeOf(arg)
	for i:=1;i<value.NumField();i++{                                                  // Field(0)跳过
		if !isEmpty( value.Field(i).Interface() ){
			keyname := typ.Field(i).Name
			keystring := strFirstToLower(keyname)
			valuestring := interfaceTostring(value.Field(i).Interface())
			partstring := fmt.Sprintf(",\"%v\":\"%v\"",keystring,valuestring)
			queryString = queryString + partstring
		}
	}
	queryString = queryString + "}}"                                                   // 封装，尾部

	return queryString
}

/**************************************************************************************/
// 保存或更新gldPyRcrd
func saveOrUpdateGldPyRcrd(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	var err error
	fmt.Println("- start SaveOrUpdate gldPyRcrd ")

	if len(args) < 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}


	gldPyRcrd := args[0]
//	fmt.Println("- the received gldPyRcrd args is :",gldPyRcrd)
	bizGldPyRcrd1 := BizGldPyRcrd{}
	err = json.Unmarshal([]byte(gldPyRcrd), &bizGldPyRcrd1)
	if err != nil {
		return shim.Error(err.Error())
	}//
	if bizGldPyRcrd1.OriGldId == "" {
		return shim.Error("oriGldId can't be null")
	} else if bizGldPyRcrd1.NewGldId == "" {
		return shim.Error("newGldId can't be null")
	}else{
		bizGldPyRcrdFromState, err := stub.GetState("BizGldPyRcrd" + bizGldPyRcrd1.OriGldId + bizGldPyRcrd1.NewGldId )
		if err != nil {
			return shim.Error("Failed to get bizGldPyRcrd:" + err.Error())
		} else if bizGldPyRcrdFromState == nil {
			bizGldPyRcrd1.ObjectType = "BizGldPyRcrd"
			bizGldPyRcrd1.DelInd = "0"
			bizGldPyRcrdToState, err := json.Marshal(bizGldPyRcrd1)
			err = stub.PutState("BizGldPyRcrd" + bizGldPyRcrd1.OriGldId + bizGldPyRcrd1.NewGldId, bizGldPyRcrdToState)
			if err != nil {
				return shim.Error(err.Error())
			}

			timestamp := time.Now().Format("2006-01-02 15:04:05")
			fmt.Printf("- save successfully ! %v \n",timestamp)
	//		fmt.Printf("- the key of record is : BizGldPyRcrd%v%v \n",bizGldPyRcrd1.OriGldId,bizGldPyRcrd1.NewGldId)
	//		fmt.Printf("- the value of record is ： %v \n",string(bizGldPyRcrdToState))

			return shim.Success(nil)
		} else {
			bizGldPyRcrd2 := BizGldPyRcrd{}
			err = json.Unmarshal([]byte(bizGldPyRcrdFromState), &bizGldPyRcrd2)
			value1 := reflect.ValueOf(&bizGldPyRcrd1).Elem()
			value2 := reflect.ValueOf(&bizGldPyRcrd2).Elem()
			for i:=0; i<value1.NumField(); i++{
				if !isEmpty( value1.Field(i).Interface() ){
					value2.Field(i).Set( value1.Field(i) )
				}
			}
			bizGldPyRcrdToState,err := json.Marshal(bizGldPyRcrd2)
			err = stub.PutState("BizGldPyRcrd" + bizGldPyRcrd1.OriGldId + bizGldPyRcrd1.NewGldId, bizGldPyRcrdToState)

			if err != nil {
				return shim.Error(err.Error())
			}

			timestamp := time.Now().Format("2006-01-02 15:04:05")
			fmt.Printf("- update successfully ! %v \n",timestamp)
	//		fmt.Printf("- the key of record is : BizGldPyRcrd%v%v \n",bizGldPyRcrd1.OriGldId,bizGldPyRcrd1.NewGldId)
	//		fmt.Printf("- the value of record is ： %v \n",string(bizGldPyRcrdToState))

		}
	}
	return shim.Success(nil)

}


// 删除gldPyRcrd
func deleteGldPyRcrd(stub shim.ChaincodeStubInterface, args []string) pb.Response {

	var jsonResp string
	if len(args) < 2 {
		return shim.Error("Incorrect number of arguments. Expecting 2")
	}

	oriGldId := strings.ToLower(args[0])
	newGldId := strings.ToLower(args[1])

	bizGldPyRcrdFromState,err := stub.GetState("BizGldPyRcrd" + oriGldId + newGldId )
	if err != nil {
		jsonResp = "{\"Error\":\"Failed to get state for " + oriGldId+newGldId + "\"}"
		return shim.Error(jsonResp)
	} else if bizGldPyRcrdFromState == nil {
		jsonResp = "{\"Error\":\"bizGldPyRcrd does not exist: " + oriGldId+newGldId + "\"}"
		return shim.Error(jsonResp)
	}
	bizGldPyRcrd := BizGldPyRcrd{}
	err = json.Unmarshal([]byte(bizGldPyRcrdFromState), &bizGldPyRcrd)
	if err != nil {
		return shim.Error(err.Error())
	}
	bizGldPyRcrd.DelInd="1"
	bizGldPyRcrdToState, err := json.Marshal(bizGldPyRcrd)
	err = stub.PutState("BizGldPyRcrd" + oriGldId + newGldId, bizGldPyRcrdToState)
	if err != nil {
		return shim.Error(err.Error())
	}

	timestamp := time.Now().Format("2006-01-02 15:04:05")
	fmt.Printf("- delete successfully ! %v \n",timestamp)
//	fmt.Printf("- the key of record is : BizGldPyRcrd%v%v \n",bizGldPyRcrd.OriGldId,bizGldPyRcrd.NewGldId)
//	fmt.Printf("- the value of record is ： %v \n",string(bizGldPyRcrdToState))

	return shim.Success(nil)
}


// 通过oriGldId、newGldId查询gldPyRcrd
func queryGldPyRcrdByOriNewGldId(stub shim.ChaincodeStubInterface, args []string) pb.Response {

	var err error
	var oriGldId,newGldId, jsonResp string

	fmt.Println("- start query gldPyRcrd by oriGldId and newGldId")

	if len(args) != 2 {
		return shim.Error("Incorrect number of arguments. Expecting oriGldId,newGldId of the gldPyRcrd to query")
	}

	oriGldId = strings.ToLower(args[0])
	newGldId = strings.ToLower(args[1])

	bizGldPyRcrdFromState, err := stub.GetState("BizGldPyRcrd" + oriGldId + newGldId)

	if err != nil {
		jsonResp = "{\"Error\":\"Failed to get state for " + oriGldId + newGldId + "\"}"
		return shim.Error(jsonResp)
	} else if bizGldPyRcrdFromState == nil {
		jsonResp = "{\"Error\":\"bizGldPyRcrd does not exist: " + oriGldId + newGldId + "\"}"
		return shim.Error(jsonResp)
	}


	var buffer bytes.Buffer                                                             // 缓存，用来保存查询结果
	buffer.WriteString("{\"Key\":")
	buffer.WriteString("\"BizGldPyRcrd")
	buffer.WriteString(oriGldId+newGldId)
	buffer.WriteString("\"")
	buffer.WriteString(", \"Value\":")
	buffer.WriteString(string(bizGldPyRcrdFromState))
	buffer.WriteString("}")

	timestamp := time.Now().Format("2006-01-02 15:04:05")
	fmt.Printf("- query by oriGldId and newGldId successfully ! %v \n",timestamp)
//	fmt.Printf("- the query result is : \n")
//	fmt.Printf("  %v \n",buffer.String())

	return shim.Success(buffer.Bytes())
}




// 通过oriGldId、newGldId范围查询gldPyRcrd(复杂查询)
func queryGldPyRcrdByGldIdRange(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	var err error
	var startOriGldId,startNewGldId,endOriGldId,endNewGldId string

	fmt.Println("- start query gldPyRcrd by oriGldId and newGldId range")

	if len(args)<4{
		return shim.Error("Incorrect number of arguments. Range query expect startOriGldId,startNewGldId,endOriGldId and endNewGldId of the gldPyRcrd to query")
	}
	startOriGldId = strings.ToLower(args[0])
	startNewGldId = strings.ToLower(args[1])
	endOriGldId = strings.ToLower(args[2])
	endNewGldId = strings.ToLower(args[3])


	startKey := "BizGldPyRcrd" + startOriGldId + startNewGldId
	endKey := "BizGldPyRcrd" + endOriGldId + endNewGldId

	resultsIterator, err := stub.GetStateByRange(startKey, endKey)
	if err != nil {
		return shim.Error(err.Error())
	}
	defer resultsIterator.Close()                                            // 迭代器用完需关闭！！

	if !resultsIterator.HasNext(){
		timestamp := time.Now().Format("2006-01-02 15:04:05")
		fmt.Printf("- can not find gldPyRcrd by oriGldId and newGldId range ! %v \n",timestamp)
		return shim.Success([]byte("can not find gldPyRcrd by oriGldId and newGldId range"))
	}
	var buffer bytes.Buffer
	buffer.WriteString("[")                                               // 对查询结果进行封装，封装为对象数组json格式
	bArrayMemberAlreadyWritten := false
	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()

		if err != nil {
			return shim.Error(err.Error())
		}
		if bArrayMemberAlreadyWritten == true {
			buffer.WriteString(",")
		}
		buffer.WriteString("{\"Key\":")
		buffer.WriteString("\"")
		buffer.WriteString(queryResponse.Key)
		buffer.WriteString("\"")
		buffer.WriteString(", \"Value\":")
		buffer.WriteString(string(queryResponse.Value))
		buffer.WriteString("}")
		bArrayMemberAlreadyWritten = true
	}
	buffer.WriteString("]")

	timestamp := time.Now().Format("2006-01-02 15:04:05")
	fmt.Printf("- query range by startOriGldId,startNewGldId,endOriGldId and endNewGldId successfully ! %v \n",timestamp)
//	fmt.Printf("- the query result is : \n")
//	fmt.Printf("  %v \n",buffer.String())

	return shim.Success(buffer.Bytes())                                      // 最终，结果以[]byte形式返回

}

// 通过querystring查询gldPyRcrd(复杂查询)
func queryGldPyRcrdByQryStr(stub shim.ChaincodeStubInterface, args []string) pb.Response {

	var err error
	var queryString string

	//	fmt.Println("- start query gldAccvalRltv by queryString, the queryString is : \n %v \n",queryString )
	fmt.Println("- start query gldPyRcrd by queryString" )

	if len(args)!= 1{
		return shim.Error("Incorrect number of arguments. Expecting querystring of the gldPyRcrd to query")
	}
	queryString = args[0]

	if queryString == ""{                                                     // 判定querystring是否为空
		return shim.Error("The queryString can not be null ")
	}
	qryStrBizGldPyRcrd := QryStrBizGldPyRcrd{}
	err = json.Unmarshal([]byte(queryString), &qryStrBizGldPyRcrd)
	if err != nil {
		return shim.Error(err.Error())
	}
	if isEmptyBizGldPyRcrd(qryStrBizGldPyRcrd.Selector){                     // 判断queryString中有没有查询条件
		return shim.Error("The selector of queryString can not be null ")
	}

	querystring := tranfBizGldPyRcrdToQryStr(qryStrBizGldPyRcrd.Selector)

	resultsIterator, err := stub.GetQueryResult(querystring)      // 返回迭代器，querystring格式：{"selector":{"docType":"BizGldAccvalRltv","CreateUser":"jzk",...}}
	if err != nil {
		return shim.Error(err.Error())
	}
	defer resultsIterator.Close()

	if !resultsIterator.HasNext(){
		timestamp := time.Now().Format("2006-01-02 15:04:05")
		fmt.Printf("- can not find gldPyRcrd by querystring ! %v \n",timestamp)
		return shim.Success([]byte("can not find gldPyRcrd by querystring"))
	}
	var buffer bytes.Buffer
	buffer.WriteString("[")                                                // 封装成对象数组json串格式
	bArrayMemberAlreadyWritten := false
	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()
		if err != nil {
			return shim.Error(err.Error())
		}
		if bArrayMemberAlreadyWritten == true {
			buffer.WriteString(",")
		}
		buffer.WriteString("{\"Key\":")
		buffer.WriteString("\"")
		buffer.WriteString(queryResponse.Key)
		buffer.WriteString("\"")

		buffer.WriteString(", \"Value\":")
		buffer.WriteString(string(queryResponse.Value))
		buffer.WriteString("}")
		bArrayMemberAlreadyWritten = true
	}
	buffer.WriteString("]")

	timestamp := time.Now().Format("2006-01-02 15:04:05")
	fmt.Printf("- query by queryString successfully ! %v \n",timestamp)
//	fmt.Printf("- the query result is : \n")
//	fmt.Printf("  %v \n",buffer.String())

	return shim.Success(buffer.Bytes())

}

// 通过condition查询gldPyRcrd(复杂查询)
func queryGldPyRcrdByCond(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	var err error

	fmt.Println("- start query gldPyRcrd by conditions" )

	if len(args) < 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}
	gldPyRcrd := args[0]                                                                  // 查询条件以json串传递，格式为：{"docType":"BizGldAccvalRltv","id":"id0001","createTime":"20171110",...}
	bizGldPyRcrd := BizGldPyRcrd{}
	err = json.Unmarshal([]byte(gldPyRcrd), &bizGldPyRcrd)
	if err != nil {
		return shim.Error(err.Error())
	}

	if isEmptyBizGldPyRcrd(bizGldPyRcrd){
		return shim.Error("The query conditions can not be null")
	}

	queryString := tranfBizGldPyRcrdToQryStr(bizGldPyRcrd)
	resultsIterator, err := stub.GetQueryResult(queryString)                               // 返回迭代器
	if err != nil {
		return shim.Error(err.Error())
	}
	defer resultsIterator.Close()

	if !resultsIterator.HasNext(){
		timestamp := time.Now().Format("2006-01-02 15:04:05")
		fmt.Printf("- can not find gldPyRcrd by querystring ! %v \n",timestamp)
		return shim.Success([]byte("can not find gldPyRcrd by conditions"))
	}
	var buffer bytes.Buffer
	buffer.WriteString("[")                                                            // 封装成对象数组json串格式
	bArrayMemberAlreadyWritten := false
	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()
		if err != nil {
			return shim.Error(err.Error())
		}
		if bArrayMemberAlreadyWritten == true {
			buffer.WriteString(",")
		}
		buffer.WriteString("{\"Key\":")
		buffer.WriteString("\"")
		buffer.WriteString(queryResponse.Key)
		buffer.WriteString("\"")

		buffer.WriteString(", \"Value\":")
		buffer.WriteString(string(queryResponse.Value))
		buffer.WriteString("}")
		bArrayMemberAlreadyWritten = true
	}
	buffer.WriteString("]")

	timestamp := time.Now().Format("2006-01-02 15:04:05")
	fmt.Printf("- query by conditions successfully ! %v \n",timestamp)
//	fmt.Printf("- the query result is : \n")
//	fmt.Printf("  %v \n",buffer.String())

	return shim.Success(buffer.Bytes())


}

// 通过gldId查询gldPyRcrd的历史（复杂查询）
func querGldPyRcrdHsyByGldId(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	var err error
	var oriGldId,newGldId string

	fmt.Println("- start query gldPyRcrd history by oriGldId and newGldId")

	if len(args) != 2 {
		return shim.Error("Incorrect number of arguments. Expecting oriGldId and new newGldId of the gldPyRcrd to query")
	}

	oriGldId = strings.ToLower(args[0])
	newGldId = strings.ToLower(args[1])

	resultsIterator, err := stub.GetHistoryForKey("BizGldPyRcrd" + oriGldId + newGldId)      // 返回迭代器
	if err != nil {
		return shim.Error(err.Error())
	}
	defer resultsIterator.Close()

	if !resultsIterator.HasNext(){
		timestamp := time.Now().Format("2006-01-02 15:04:05")
		fmt.Printf("- can not find gldPyRcrd history by oriGldId and newGldId ! %v \n",timestamp)
		return shim.Success([]byte("can not find gldPyRcrd history by oriGldId and newGldId"))
	}
	var buffer bytes.Buffer
	buffer.WriteString("[")
	bArrayMemberAlreadyWritten := false
	for resultsIterator.HasNext() {
		response, err := resultsIterator.Next()
		if err != nil {
			return shim.Error(err.Error())
		}
		if bArrayMemberAlreadyWritten == true {
			buffer.WriteString(",")
		}
		buffer.WriteString("{\"TxId\":")
		buffer.WriteString("\"")
		buffer.WriteString(response.TxId)
		buffer.WriteString("\"")

		buffer.WriteString(", \"Value\":")
		if response.IsDelete {
			buffer.WriteString("null")
		} else {
			buffer.WriteString(string(response.Value))
		}

		buffer.WriteString(", \"Timestamp\":")
		buffer.WriteString("\"")
		buffer.WriteString(time.Unix(response.Timestamp.Seconds, int64(response.Timestamp.Nanos)).String())
		buffer.WriteString("\"")

		buffer.WriteString(", \"IsDelete\":")
		buffer.WriteString("\"")
		buffer.WriteString(strconv.FormatBool(response.IsDelete))
		buffer.WriteString("\"")

		buffer.WriteString("}")
		bArrayMemberAlreadyWritten = true
	}
	buffer.WriteString("]")

	timestamp := time.Now().Format("2006-01-02 15:04:05")
	fmt.Printf("- query history by oriGldId and newGldId successfully ! %v \n",timestamp)
//	fmt.Printf("- the query result is : \n")
//	fmt.Printf("  %v \n",buffer.String())

	return shim.Success(buffer.Bytes())

}







