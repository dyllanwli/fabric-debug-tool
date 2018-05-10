package main

import (
	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
	"fmt"
	"encoding/json"
	"reflect"
	"time"
	"strings"
	"strconv"
	"bytes"
)

/**************************************************************************************/
// 金票签收记录结构体
type BizGldSigntoacptRcrd struct {
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

	GldId          string  `json:"gldId"`          // 金票编号
	GldAmt         float64 `json:"gldAmt"`         // 金票金额
	SignToacpt     string  `json:"signToacpt"`     // 签收方
	SignToacptAmt  float64 `json:"signToacptAmt"`  // 签收金额
	SignToacptTm   int64   `json:"signToacptTm"`   // 签收时间
	SignToacptStCd string  `json:"signToacptStCd"` // 签收状态：0.成功签收  1.拒绝签收
}


// 用于解析queryString
type QryStrBizGldSigntoacptRcrd struct{
	Selector          BizGldSigntoacptRcrd         `json:"selector"`
}

// 判断BizGldSigntoacptRcrd变量是否为空
func isEmptyBizGldSigntoacptRcrd(arg BizGldSigntoacptRcrd) bool {
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

// BizGldSigntoacptRcrd变量变为queryString
func tranfBizGldSigntoacptRcrdToQryStr(arg BizGldSigntoacptRcrd)( string){
	queryString := "{\"selector\":{\"docType\":\"BizGldSigntoacptRcrd\""                  // 封装，头部

	value := reflect.ValueOf(arg)
	typ := reflect.TypeOf(arg)
	for i:=1;i<value.NumField();i++{                                                     // Field(0)跳过
		if !isEmpty( value.Field(i).Interface() ){
			keyname := typ.Field(i).Name
			keystring := strFirstToLower(keyname)
			valuestring := interfaceTostring(value.Field(i).Interface())
			partstring := fmt.Sprintf(",\"%v\":\"%v\"",keystring,valuestring)
			queryString = queryString + partstring
		}
	}
	queryString = queryString + "}}"                                                     // 封装，尾部

	return queryString
}
/**************************************************************************************/
// 保存或更新gldSigntoacptRcrd
func saveOrUpdateGldSigntoacptRcrd(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	var err error
	fmt.Println("- start saveOrUpdate gldSigntoacptRcrd")

	if len(args) < 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}

	gldSigntoacptRcrd := args[0]
//	fmt.Println("- the received gldSigntoacptRcrd args is :",gldSigntoacptRcrd)
	bizGldSigntoacptRcrd1 := BizGldSigntoacptRcrd{}
	err = json.Unmarshal([]byte(gldSigntoacptRcrd), &bizGldSigntoacptRcrd1)
	if err != nil {
		return shim.Error(err.Error())
	}

	if bizGldSigntoacptRcrd1.GldId == "" {
		return shim.Error("gldId can't be null")
	} else {
		bizGldSigntoacptRcrdFromState, err := stub.GetState("BizGldSigntoacptRcrd" + bizGldSigntoacptRcrd1.GldId)
		if err != nil {
			return shim.Error("Failed to get gldSigntoacptRcrd:" + err.Error())
		} else if bizGldSigntoacptRcrdFromState == nil {
			bizGldSigntoacptRcrd1.ObjectType = "BizGldSigntoacptRcrd"
			bizGldSigntoacptRcrd1.DelInd = "0"
			bizGldSigntoacptRcrdToState, err := json.Marshal(bizGldSigntoacptRcrd1)
			err = stub.PutState("BizGldSigntoacptRcrd" + bizGldSigntoacptRcrd1.GldId, bizGldSigntoacptRcrdToState)
			if err != nil {
				return shim.Error(err.Error())
			}

			timestamp := time.Now().Format("2006-01-02 15:04:05")
			fmt.Printf("- save successfully ! %v \n",timestamp)
	//		fmt.Printf("- the key of record is : BizGldSigntoacptRcrd%v \n",bizGldSigntoacptRcrd1.GldId)
	//		fmt.Printf("- the value of record is ： %v \n",string(bizGldSigntoacptRcrdToState))

			return shim.Success(nil)
		} else {
			bizGldSigntoacptRcrd2 := BizGldSigntoacptRcrd{}
			err = json.Unmarshal([]byte(bizGldSigntoacptRcrdFromState), &bizGldSigntoacptRcrd2)
			value1 := reflect.ValueOf(&bizGldSigntoacptRcrd1).Elem()
			value2 := reflect.ValueOf(&bizGldSigntoacptRcrd2).Elem()
			for i:=0; i<value1.NumField(); i++{
				if !isEmpty( value1.Field(i).Interface() ){
					value2.Field(i).Set( value1.Field(i) )
				}
			}
			bizGldSigntoacptRcrdToState,err := json.Marshal(bizGldSigntoacptRcrd2)
			err = stub.PutState("BizGldSigntoacptRcrd" + bizGldSigntoacptRcrd1.GldId, bizGldSigntoacptRcrdToState)

			if err != nil {
				return shim.Error(err.Error())
			}

			timestamp := time.Now().Format("2006-01-02 15:04:05")
			fmt.Printf("- update successfully ! %v \n",timestamp)
	//		fmt.Printf("- the key of record is : BizGldSigntoacptRcrd%v \n",bizGldSigntoacptRcrd1.GldId)
	//		fmt.Printf("- the value of record is ： %v \n",string(bizGldSigntoacptRcrdToState))

		}
	}
	return shim.Success(nil)
}


// 通过gldId查询gldSigntoacptRcrd
func queryGldSigntoacptRcrdByGldId(stub shim.ChaincodeStubInterface, args []string) pb.Response {

	var err error
	var gldId, jsonResp string

	fmt.Println("- start query gldSigntoacptRcrd by gldId")

	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting gldId of the gldSigntoacptRcrd to query")
	}

	gldId = strings.ToLower(args[0])

	bizGldSigntoacptRcrdFromState, err := stub.GetState("BizGldSigntoacptRcrd" + gldId)

	if err != nil {
		jsonResp = "{\"Error\":\"Failed to get state for " + gldId + "\"}"
		return shim.Error(jsonResp)
	} else if bizGldSigntoacptRcrdFromState == nil {
		jsonResp = "{\"Error\":\"bizGldSigntoacptRcrd does not exist: " + gldId + "\"}"
		return shim.Error(jsonResp)
	}


	var buffer bytes.Buffer                                                            // 缓存，用来保存查询结果
	buffer.WriteString("{\"Key\":")
	buffer.WriteString("\"BizGldSigntoacptRcrd")
	buffer.WriteString(gldId)
	buffer.WriteString("\"")
	buffer.WriteString(", \"Value\":")
	buffer.WriteString(string(bizGldSigntoacptRcrdFromState))
	buffer.WriteString("}")

	timestamp := time.Now().Format("2006-01-02 15:04:05")
	fmt.Printf("- query by gldId successfully ! %v \n",timestamp)
//	fmt.Printf("- the query result is : \n")
//	fmt.Printf("  %v \n",buffer.String())

	return shim.Success(buffer.Bytes())
}




// 通过gldId范围查询gldSigntoacptRcrd(复杂查询)
func queryGldSigntoacptRcrdByGldIdRange(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	var err error
	var startGldId,endGldId string

	fmt.Println("- start query gldSigntoacptRcrd by gldId range")

	if len(args)<2{
		return shim.Error("Incorrect number of arguments. Range query expect startGldId and endGldId of the gldSigntoacptRcrd to query")
	}
	startGldId = strings.ToLower(args[0])
	endGldId = strings.ToLower(args[1])

	startKey := "BizGldSigntoacptRcrd" + startGldId
	endKey := "BizGldSigntoacptRcrd" + endGldId

	resultsIterator, err := stub.GetStateByRange(startKey, endKey)
	if err != nil {
		return shim.Error(err.Error())
	}
	defer resultsIterator.Close()                                            // 迭代器用完需关闭！！

	if !resultsIterator.HasNext(){
		timestamp := time.Now().Format("2006-01-02 15:04:05")
		fmt.Printf("- can not find gldSigntoacptRcrd by gldId range ! %v \n",timestamp)
		return shim.Success([]byte("can not find gldSigntoacptRcrd by gldId range"))
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
	fmt.Printf("- query range by startGldId and endGldId successfully ! %v \n",timestamp)
//	fmt.Printf("- the query result is : \n")
//	fmt.Printf("  %v \n",buffer.String())

	return shim.Success(buffer.Bytes())                                      // 最终，结果以[]byte形式返回

}

// 通过querystring查询gldSigntoacptRcrd(复杂查询)
func queryGldSigntoacptRcrdByQryStr(stub shim.ChaincodeStubInterface, args []string) pb.Response {

	var err error
	var queryString string

	//	fmt.Println("- start query gldAccvalRltv by queryString, the queryString is : \n %v \n",queryString )
	fmt.Println("- start query gldSigntoacptRcrd by queryString" )

	if len(args)!= 1{
		return shim.Error("Incorrect number of arguments. Expecting querystring of the gldSigntoacptRcrd to query")
	}
	queryString = args[0]

	if queryString == ""{                                                                   // 判定querystring是否为空
		return shim.Error("The queryString can not be null ")
	}
	qryStrBizGldSigntoacptRcrd := QryStrBizGldSigntoacptRcrd{}
	err = json.Unmarshal([]byte(queryString), &qryStrBizGldSigntoacptRcrd)
	if err != nil {
		return shim.Error(err.Error())
	}
	if isEmptyBizGldSigntoacptRcrd(qryStrBizGldSigntoacptRcrd.Selector){                    // 判断queryString中有没有查询条件
		return shim.Error("The selector of queryString can not be null ")
	}

	querystring := tranfBizGldSigntoacptRcrdToQryStr(qryStrBizGldSigntoacptRcrd.Selector)

	resultsIterator, err := stub.GetQueryResult(querystring)                     // 返回迭代器，querystring格式：{"selector":{"docType":"BizGldAccvalRltv","CreateUser":"jzk",...}}
	if err != nil {
		return shim.Error(err.Error())
	}
	defer resultsIterator.Close()

	if !resultsIterator.HasNext(){
		timestamp := time.Now().Format("2006-01-02 15:04:05")
		fmt.Printf("- can not find gldSigntoacptRcrd by querystring ! %v \n",timestamp)
		return shim.Success([]byte("can not find gldSigntoacptRcrd by querystring"))
	}
	var buffer bytes.Buffer
	buffer.WriteString("[")                                                              // 封装成对象数组json串格式
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

// 通过condition查询gldSigntoacptRcrd(复杂查询)
func queryGldSigntoacptRcrdByCond(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	var err error

	fmt.Println("- start query gldSigntoacptRcrd by conditions" )

	if len(args) < 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}
	gldSigntoacptRcrd := args[0]                                                             // 查询条件以json串传递，格式为：{"docType":"BizGldAccvalRltv","id":"id0001","createTime":"20171110",...}
	bizGldSigntoacptRcrd := BizGldSigntoacptRcrd{}
	err = json.Unmarshal([]byte(gldSigntoacptRcrd), &bizGldSigntoacptRcrd)
	if err != nil {
		return shim.Error(err.Error())
	}

	if isEmptyBizGldSigntoacptRcrd(bizGldSigntoacptRcrd){
		return shim.Error("The query conditions can not be null")
	}

	queryString := tranfBizGldSigntoacptRcrdToQryStr(bizGldSigntoacptRcrd)
	resultsIterator, err := stub.GetQueryResult(queryString)                           // 返回迭代器
	if err != nil {
		return shim.Error(err.Error())
	}
	defer resultsIterator.Close()

	if !resultsIterator.HasNext(){
		timestamp := time.Now().Format("2006-01-02 15:04:05")
		fmt.Printf("- can not find gldSigntoacptRcrd by querystring ! %v \n",timestamp)
		return shim.Success([]byte("can not find gldSigntoacptRcrd by conditions"))
	}
	var buffer bytes.Buffer
	buffer.WriteString("[")                                                         // 封装成对象数组json串格式
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

// 通过gldId查询gldSigntoacptRcrd的历史（复杂查询）
func queryGldSigntoacptRcrdHsyByGldId(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	var err error
	var gldId string

	fmt.Println("- start query gldSigntoacptRcrd history by gldId")

	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting gldId of the gldSigntoacptRcrd to query")
	}

	gldId = strings.ToLower(args[0])

	resultsIterator, err := stub.GetHistoryForKey("BizGldSigntoacptRcrd"+ gldId)      // 返回迭代器
	if err != nil {
		return shim.Error(err.Error())
	}
	defer resultsIterator.Close()

	if !resultsIterator.HasNext(){
		timestamp := time.Now().Format("2006-01-02 15:04:05")
		fmt.Printf("- can not find gldSigntoacptRcrd history by gldId ! %v \n",timestamp)
		return shim.Success([]byte("can not find gldSigntoacptRcrd history by gldId"))
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
	fmt.Printf("- query history by gldId successfully ! %v \n",timestamp)
//	fmt.Printf("- the query result is : \n")
//	fmt.Printf("  %v \n",buffer.String())

	return shim.Success(buffer.Bytes())

}




