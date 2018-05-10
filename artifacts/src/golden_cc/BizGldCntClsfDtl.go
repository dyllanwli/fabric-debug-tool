package main

import (
	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
	"encoding/json"
	"fmt"
	"time"
	"strings"
	"strconv"
	"bytes"
	"reflect"
)

/**************************************************************************************/
// 金票清分结构体
type BizGldCntClsfDtl struct {
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

	GldId      string  `json:"gldId"`      // 金票编号
	GldAmt     float64 `json:"gldAmt"`     // 金票金额
	GldBal     float64 `json:"gldBal"`     // 金票余额
	Pyr        string  `json:"pyr"`        // 付款方
	Rcvprt     string  `json:"rcvprt"`     // 收款方
	RcvpymtAmt float64 `json:"rcvpymtAmt"` // 收款金额
	RcvpymtTm  int64   `json:"rcvpymtTm"`  // 收款时间
	CntClsfTp  string  `json:"cntClsfTp"`  // 清分类型;0-融资清分;1-垫付清分

	// 2017.12.14  新增
	AplyId       string   `json:"aplyId"`       // 申请编号
	FncJrnlId    string   `json:"fncJrnlId"`    // 融资编号
	RcvpymtInt   float64  `json:"rcvpymtInt"`   // 收款利息
	ClientId     string   `json:"clientId"`     // 流水号
	StCd         string   `json:"stCd"`         // 清分状态（0-未转账，1-已转账）
	OriGldId     string   `json:"oriGldId"`     // 原票编号

}

// 用于解析queryString
type QryStrBizGldCntClsfDtl struct{
	Selector          BizGldCntClsfDtl         `json:"selector"`
}

// 判断BizGldCntClsfDtl变量是否为空
func isEmptyBizGldCntClsfDtl(arg BizGldCntClsfDtl) bool {
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

// BizGldCntClsfDtl变量变为queryString
func tranfBizGldCntClsfDtlToQryStr(arg BizGldCntClsfDtl)( string){
	queryString := "{\"selector\":{\"docType\":\"BizGldCntClsfDtl\""                  // 封装，头部

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
// 保存gldCntClsfDtl
func saveGldCntClsfDtl(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	var err error
	fmt.Println("- start save gldCntClsfDtl")

	if len(args) < 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}
	gldCntClsfDtl := args[0]
//	fmt.Println("- the received gldCntClsfDtl args is :",gldCntClsfDtl)
	bizGldCntClsfDtl1 := BizGldCntClsfDtl{}
	err = json.Unmarshal([]byte(gldCntClsfDtl), &bizGldCntClsfDtl1)
	if err != nil {
		return shim.Error(err.Error())
	}

	if bizGldCntClsfDtl1.GldId == "" {
		return shim.Error("gldId can't be null")
	} else {
		bizGldCntClsfDtl1.ObjectType = "BizGldCntClsfDtl"
		bizGldCntClsfDtl1.DelInd = "0"
		bizGldCntClsfDtlToState, err := json.Marshal(bizGldCntClsfDtl1)
		err = stub.PutState("BizGldCntClsfDtl" + bizGldCntClsfDtl1.GldId, bizGldCntClsfDtlToState)
		if err != nil {
			return shim.Error(err.Error())
		}

		timestamp := time.Now().Format("2006-01-02 15:04:05")
		fmt.Printf("- save successfully ! %v \n",timestamp)
	//	fmt.Printf("- the key of record is : BizGldCntClsfDtl%v \n",bizGldCntClsfDtl1.GldId)
	//	fmt.Printf("- the value of record is ： %v \n",string(bizGldCntClsfDtlToState))

		return shim.Success(nil)
	}
}


// 通过gldId查询gldCntClsfDtl
func queryGldCntClsfDtlByGldId(stub shim.ChaincodeStubInterface, args []string) pb.Response {

	var err error
	var gldId, jsonResp string

	fmt.Println("- start query gldCntClsfDtl by gldId")

	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting gldId of the gldCntClsfDtl to query")
	}

	gldId = strings.ToLower(args[0])

	bizGldCntClsfDtlFromState, err := stub.GetState("BizGldCntClsfDtl" + gldId)

	if err != nil {
		jsonResp = "{\"Error\":\"Failed to get state for " + gldId + "\"}"
		return shim.Error(jsonResp)
	} else if bizGldCntClsfDtlFromState == nil {
		jsonResp = "{\"Error\":\"bizGldCntClsfDtl does not exist: " + gldId + "\"}"
		return shim.Error(jsonResp)
	}


	var buffer bytes.Buffer                                                  // 缓存，用来保存查询结果
	buffer.WriteString("{\"Key\":")
	buffer.WriteString("\"BizGldCntClsfDtl")
	buffer.WriteString(gldId)
	buffer.WriteString("\"")
	buffer.WriteString(", \"Value\":")
	buffer.WriteString(string(bizGldCntClsfDtlFromState))
	buffer.WriteString("}")

	timestamp := time.Now().Format("2006-01-02 15:04:05")
	fmt.Printf("- query by gldId successfully ! %v \n",timestamp)
//	fmt.Printf("- the query result is : \n")
//	fmt.Printf("  %v \n",buffer.String())

	return shim.Success(buffer.Bytes())
}

// 通过gldId范围查询gldCntClsfDtl(复杂查询)
func queryGldCntClsfDtlByGldIdRange(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	var err error
	var startGldId,endGldId string

	fmt.Println("- start query gldCntClsfDtl by gldId range")

	if len(args)<2{
		return shim.Error("Incorrect number of arguments. Range query expect startGldId and endGldId of the gldCntClsfDtl to query")
	}
	startGldId = strings.ToLower(args[0])
	endGldId = strings.ToLower(args[1])

	startKey := "BizGldCntClsfDtl" + startGldId
	endKey := "BizGldCntClsfDtl" + endGldId

	resultsIterator, err := stub.GetStateByRange(startKey, endKey)
	if err != nil {
		return shim.Error(err.Error())
	}
	defer resultsIterator.Close()                                            // 迭代器用完需关闭！！

	if !resultsIterator.HasNext(){
		timestamp := time.Now().Format("2006-01-02 15:04:05")
		fmt.Printf("- can not find gldCntClsfDtl by gldId range ! %v \n",timestamp)
		return shim.Success([]byte("can not find gldCntClsfDtl by gldId range"))
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

	return shim.Success(buffer.Bytes())                                               // 最终，结果以[]byte形式返回

}

// 通过querystring查询gldCntClsfDtl(复杂查询)
func queryGldCntClsfDtlByQryStr(stub shim.ChaincodeStubInterface, args []string) pb.Response {

	var err error
	var queryString string

	//	fmt.Println("- start query gldAccvalRltv by queryString, the queryString is : \n %v \n",queryString )
	fmt.Println("- start query gldCntClsfDtl by queryString" )

	if len(args)!= 1{
		return shim.Error("Incorrect number of arguments. Expecting querystring of the gldCntClsfDtl to query")
	}
	queryString = args[0]

	if queryString == ""{                                                           // 判定querystring是否为空
		return shim.Error("The queryString can not be null ")
	}
	qryStrBizGldCntClsfDtl := QryStrBizGldCntClsfDtl{}
	err = json.Unmarshal([]byte(queryString), &qryStrBizGldCntClsfDtl)
	if err != nil {
		return shim.Error(err.Error())
	}
	if isEmptyBizGldCntClsfDtl(qryStrBizGldCntClsfDtl.Selector){                    // 判断queryString中有没有查询条件
		return shim.Error("The selector of queryString can not be null ")
	}

	querystring := tranfBizGldCntClsfDtlToQryStr(qryStrBizGldCntClsfDtl.Selector)

	resultsIterator, err := stub.GetQueryResult(querystring)                        // 返回迭代器，querystring格式：{"selector":{"docType":"BizGldAccvalRltv","CreateUser":"jzk",...}}
	if err != nil {
		return shim.Error(err.Error())
	}
	defer resultsIterator.Close()

	if !resultsIterator.HasNext(){
		timestamp := time.Now().Format("2006-01-02 15:04:05")
		fmt.Printf("- can not find gldCntClsfDtl by querystring ! %v \n",timestamp)
		return shim.Success([]byte("can not find gldCntClsfDtl by querystring"))
	}
	var buffer bytes.Buffer
	buffer.WriteString("[")                                                      // 封装成对象数组json串格式
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

// 通过condition查询gldCntClsfDtl(复杂查询)
func queryGldCntClsfDtlByCond(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	var err error

	fmt.Println("- start query gldCntClsfDtl by conditions" )

	if len(args) < 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}
	gldCntClsfDtl := args[0]                                                               // 查询条件以json串传递，格式为：{"docType":"BizGldAccvalRltv","id":"id0001","createTime":"20171110",...}
	bizGldCntClsfDtl := BizGldCntClsfDtl{}
	err = json.Unmarshal([]byte(gldCntClsfDtl), &bizGldCntClsfDtl)
	if err != nil {
		return shim.Error(err.Error())
	}

	if isEmptyBizGldCntClsfDtl(bizGldCntClsfDtl){
		return shim.Error("The query conditions can not be null")
	}

	queryString := tranfBizGldCntClsfDtlToQryStr(bizGldCntClsfDtl)
	resultsIterator, err := stub.GetQueryResult(queryString)                                // 返回迭代器
	if err != nil {
		return shim.Error(err.Error())
	}
	defer resultsIterator.Close()

	if !resultsIterator.HasNext(){
		timestamp := time.Now().Format("2006-01-02 15:04:05")
		fmt.Printf("- can not find gldCntClsfDtl by querystring ! %v \n",timestamp)
		return shim.Success([]byte("can not find gldCntClsfDtl by conditions"))
	}
	var buffer bytes.Buffer
	buffer.WriteString("[")                                                             // 封装成对象数组json串格式
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

// 通过gldId查询gldCntClsfDtl的历史（复杂查询）
func queryGldCntClsfDtlHsyByGldId(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	var err error
	var gldId string

	fmt.Println("- start query gldCntClsfDtl history by gldId")

	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting gldId of the gldCntClsfDtl to query")
	}

	gldId = strings.ToLower(args[0])

	resultsIterator, err := stub.GetHistoryForKey("BizGldCntClsfDtl"+ gldId)               // 返回迭代器
	if err != nil {
		return shim.Error(err.Error())
	}
	defer resultsIterator.Close()

	if !resultsIterator.HasNext(){
		timestamp := time.Now().Format("2006-01-02 15:04:05")
		fmt.Printf("- can not find gldCntClsfDtl history by gldId ! %v \n",timestamp)
		return shim.Success([]byte("can not find gldCntClsfDtl history by gldId"))
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
		if response.IsDelete {                                                            // 若已删除则value为null
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
