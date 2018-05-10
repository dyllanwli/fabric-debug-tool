package main

import (
	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
	"encoding/json"
	"fmt"
	"strings"
	"reflect"
	"time"
	"strconv"
	"bytes"
)
/**************************************************************************************/
// 金票确认结构体
type BizGldCfm struct{
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

	CfmAplyId          string    `json:"cfmAplyId"`     // 确认申请编号
	PcsStCd            string    `json:"pcsStCd"`       // 流程状态
	Opin               string    `json:"opin"`          // 意见

	// 2017.11.27 新增
	RspbPsnId          string    `json:"rspbPsnId"`     // 经办人编号
	HdlInstId          string    `json:"hdlInstId"`     // 经办机构编号
	HdlDt              int64     `json:"hdlDt"`         // 经办日期（营业日）
}

// 用于解析queryString
type QryStrBizGldCfm struct{
	Selector          BizGldCfm         `json:"selector"`
}

// 判断BizGldCfm变量是否为空
func isEmptyBizGldCfm(arg BizGldCfm) bool {
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

// BizGldCfm变量变为queryString
func tranfBizGldCfmToQryStr(arg BizGldCfm)( string){
	queryString := "{\"selector\":{\"docType\":\"BizGldCfm\""

	value := reflect.ValueOf(arg)
	typ := reflect.TypeOf(arg)
	for i:=1;i<value.NumField();i++{                                                         // Field(0)跳过
		if !isEmpty( value.Field(i).Interface() ){
			keyname := typ.Field(i).Name
			keystring := strFirstToLower(keyname)
			valuestring := interfaceTostring(value.Field(i).Interface())
			partstring := fmt.Sprintf(",\"%v\":\"%v\"",keystring,valuestring)
			queryString = queryString + partstring
		}
	}
	queryString = queryString + "}}"

	return queryString
}

/**************************************************************************************/
// 保存或更新gldCfm
func saveOrUpdateGldCfm(stub shim.ChaincodeStubInterface, args []string) pb.Response {

	var err error
	fmt.Println("- start saveOrUpdate gldCfm")

	if len(args) < 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}

	gldCfm := args[0]
//	fmt.Println("- the received gldCfm args is :",gldCfm)
	bizGldCfm1 := BizGldCfm{}
	err = json.Unmarshal([]byte(gldCfm), &bizGldCfm1)
	if err != nil {
		return shim.Error(err.Error())
	}

	if bizGldCfm1.CfmAplyId == "" {
		return shim.Error("cfmAplyId can't be null")
	} else {
		bizGldCfmFromState, err := stub.GetState("BizGldCfm" +bizGldCfm1.CfmAplyId)
		if err != nil {
			return shim.Error("Failed to get bizGldCfm:" + err.Error())
		} else if bizGldCfmFromState == nil {
			bizGldCfm1.ObjectType = "BizGldCfm"
			bizGldCfm1.DelInd = "0"                                                           // 删除标志，初始值为0
			bizGldCfmToState, err := json.Marshal(bizGldCfm1)
			err = stub.PutState("BizGldCfm" +bizGldCfm1.CfmAplyId, bizGldCfmToState)
			if err != nil {
				return shim.Error(err.Error())
			}

			timestamp := time.Now().Format("2006-01-02 15:04:05")
			fmt.Printf("- save successfully ! %v \n",timestamp)
	//		fmt.Printf("- the key of record is : BizGldCfm%v \n",bizGldCfm1.CfmAplyId)
	//		fmt.Printf("- the value of rcord is ： %v \n",string(bizGldCfmToState))

			return shim.Success(nil)
		} else {
			bizGldCfm2 := BizGldCfm{}
			err = json.Unmarshal([]byte(bizGldCfmFromState), &bizGldCfm2)
			value1 := reflect.ValueOf(&bizGldCfm1).Elem()
			value2 := reflect.ValueOf(&bizGldCfm2).Elem()
			for i:=0; i<value1.NumField(); i++{
				if !isEmpty( value1.Field(i).Interface() ){
					value2.Field(i).Set( value1.Field(i) )
				}
			}
			bizGldCfmToState,err := json.Marshal(bizGldCfm2)
			err = stub.PutState("BizGldCfm" + bizGldCfm1.CfmAplyId, bizGldCfmToState)
			if err != nil {
				return shim.Error(err.Error())
			}

			timestamp := time.Now().Format("2006-01-02 15:04:05")
			fmt.Printf("- update successfully ! %v \n",timestamp)
	//		fmt.Printf("- the key of record is : BizGldCfm%v \n",bizGldCfm1.CfmAplyId)
	//		fmt.Printf("- the value of record is ： %v \n",string(bizGldCfmToState))

		}
	}
	return shim.Success(nil)
}

// 删除gldCfm
func deleteGldCfm(stub shim.ChaincodeStubInterface, args []string) pb.Response {

	var err error
	var cfmAplyId, jsonResp string

	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting cfmAplyId of the gldCfm to delete")
	}

	cfmAplyId = strings.ToLower(args[0])                                                      // 传入的不再是json串，而是cfmAplyId，string格式

	bizGldCfmFromState, err := stub.GetState("BizGldCfm" + cfmAplyId)
	if err != nil {
		jsonResp = "{\"Error\":\"Failed to get state for " + cfmAplyId + "\"}"
		return shim.Error(jsonResp)
	} else if bizGldCfmFromState == nil {
		jsonResp = "{\"Error\":\"bizGldCfm does not exist: " + cfmAplyId + "\"}"
		return shim.Error(jsonResp)
	}
	bizGldCfm := BizGldCfm{}

	err = json.Unmarshal(bizGldCfmFromState,&bizGldCfm)

	bizGldCfm.DelInd = "1"
	bizGldCfmToState, err := json.Marshal(bizGldCfm)
	err = stub.PutState("BizGldCfm" + bizGldCfm.CfmAplyId, bizGldCfmToState)                  // 删除并非真的从数据中删除记录，只是将删除标志位置位而已
	if err != nil {
		return shim.Error(err.Error())
	}

	timestamp := time.Now().Format("2006-01-02 15:04:05")
	fmt.Printf("- delete successfully ! %v \n",timestamp)
//	fmt.Printf("- the key of record is : BizGldCfm%v \n",bizGldCfm.CfmAplyId)
//	fmt.Printf("- the value of record is ： %v \n",string(bizGldCfmToState))

	return shim.Success(nil)
}


// 通过gldId查询gldCfm
func queryGldCfmByCfmAplyId(stub shim.ChaincodeStubInterface, args []string) pb.Response {

	var err error
	var cfmAplyId, jsonResp string

	fmt.Println("- start query gldCfm by cfmAplyId")

	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting cfmAplyId of the gldCfm to query")
	}

	cfmAplyId = strings.ToLower(args[0])

	bizBizGldCfmFromState, err := stub.GetState("BizGldCfm" + cfmAplyId)

	if err != nil {
		jsonResp = "{\"Error\":\"Failed to get state for " + cfmAplyId + "\"}"
		return shim.Error(jsonResp)
	} else if bizBizGldCfmFromState == nil {
		jsonResp = "{\"Error\":\"BizGldCfm does not exist: " + cfmAplyId + "\"}"
		return shim.Error(jsonResp)
	}

	var buffer bytes.Buffer
	buffer.WriteString("{\"Key\":")
	buffer.WriteString("\"BizGldCfm")
	buffer.WriteString(cfmAplyId)
	buffer.WriteString("\"")
	buffer.WriteString(", \"Value\":")
	buffer.WriteString(string(bizBizGldCfmFromState))
	buffer.WriteString("}")

	timestamp := time.Now().Format("2006-01-02 15:04:05")
	fmt.Printf("- query by cfmAplyId successfully ! %v \n",timestamp)
//	fmt.Printf("- the query result is : \n")
//	fmt.Printf("  %v \n",buffer.String())

	return shim.Success(buffer.Bytes())

}

// 通过gldId范围查询gldCfm(复杂查询)
func queryGldCfmByCfmAplyIdRange(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	var err error
	var startCfmAplyId,endCfmAplyId string

	fmt.Println("- start query gldCfm by cfmAplyId range")

	if len(args)<2{
		return shim.Error("Incorrect number of arguments. Range query expect startCfmAplyId and endCfmAplyId of the gldCfm to query")
	}
	startCfmAplyId = strings.ToLower(args[0])
	endCfmAplyId = strings.ToLower(args[1])

	startKey := "BizGldCfm" + startCfmAplyId
	endKey := "BizGldCfm" + endCfmAplyId

	resultsIterator, err := stub.GetStateByRange(startKey, endKey)
	if err != nil {
		return shim.Error(err.Error())
	}
	defer resultsIterator.Close()                                                            // 迭代器用完需关闭！！

	if !resultsIterator.HasNext(){
		timestamp := time.Now().Format("2006-01-02 15:04:05")
		fmt.Printf("- can not find gldCfm by cfmAplyId range ! %v \n",timestamp)
		return shim.Success([]byte("can not find gldCfm by cfmAplyId range"))
	}
	var buffer bytes.Buffer
	buffer.WriteString("[")                                                               // 对查询结果进行封装，封装为对象数组json格式
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
	fmt.Printf("- query range by startCfmAplyId and endCfmAplyId successfully ! %v \n",timestamp)
//	fmt.Printf("- the query result is : \n")
//	fmt.Printf("  %v \n",buffer.String())

	return shim.Success(buffer.Bytes())                                                         // 最终，结果以[]byte形式返回

}

// 通过querystring查询gldCfm(复杂查询)
func queryGldCfmByQryStr(stub shim.ChaincodeStubInterface, args []string) pb.Response {

	var err error
	var queryString string

	//	fmt.Println("- start query gldAccvalRltv by queryString, the queryString is : \n %v \n",queryString )
	fmt.Println("- start query gldCfm by queryString" )

	if len(args)!= 1{
		return shim.Error("Incorrect number of arguments. Expecting querystring of the gldCfm to query")
	}
	queryString = args[0]

	if queryString == ""{                                                                        // 判定querystring是否为空
		return shim.Error("The queryString can not be null ")
	}
	qryStrBizGldCfm := QryStrBizGldCfm{}
	err = json.Unmarshal([]byte(queryString), &qryStrBizGldCfm)
	if err != nil {
		return shim.Error(err.Error())
	}
	if isEmptyBizGldCfm(qryStrBizGldCfm.Selector){                                               // 判断queryString中有没有查询条件
		return shim.Error("The selector of queryString can not be null ")
	}

	querystring := tranfBizGldCfmToQryStr(qryStrBizGldCfm.Selector)

	resultsIterator, err := stub.GetQueryResult(querystring)                                     // 返回迭代器，querystring格式：{"selector":{"docType":"BizGldAccvalRltv","CreateUser":"jzk",...}}
	if err != nil {
		return shim.Error(err.Error())
	}
	defer resultsIterator.Close()

	if !resultsIterator.HasNext(){
		timestamp := time.Now().Format("2006-01-02 15:04:05")
		fmt.Printf("- can not find gldCfm by querystring ! %v \n",timestamp)
		return shim.Success([]byte("can not find gldCfm by querystring"))
	}
	var buffer bytes.Buffer
	buffer.WriteString("[")                                                                   // 封装成对象数组json串格式
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

// 通过condition查询gldCfm(复杂查询)
func queryGldCfmByCond(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	var err error

	fmt.Println("- start query gldCfm by conditions" )

	if len(args) < 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}
	gldCfm := args[0]                                                                             // 查询条件以json串传递，格式为：{"docType":"BizGldAccvalRltv","id":"id0001","createTime":"20171110",...}
	bizGldCfm := BizGldCfm{}
	err = json.Unmarshal([]byte(gldCfm), &bizGldCfm)
	if err != nil {
		return shim.Error(err.Error())
	}

	if isEmptyBizGldCfm(bizGldCfm){
		return shim.Error("The query conditions can not be null")
	}

	queryString := tranfBizGldCfmToQryStr(bizGldCfm)
	resultsIterator, err := stub.GetQueryResult(queryString)                                      // 返回迭代器
	if err != nil {
		return shim.Error(err.Error())
	}
	defer resultsIterator.Close()

	if !resultsIterator.HasNext(){
		timestamp := time.Now().Format("2006-01-02 15:04:05")
		fmt.Printf("- can not find gldCfm by querystring ! %v \n",timestamp)
		return shim.Success([]byte("can not find gldCfm by conditions"))
	}
	var buffer bytes.Buffer
	buffer.WriteString("[")                                                                   // 封装成对象数组json串格式
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

// 通过gldId查询gldCfm的历史（复杂查询）
func queryGldCfmHsyByGldId(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	var err error
	var cfmAplyId string

	fmt.Println("- start query gldCfm history by cfmAplyId")

	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting cfmAplyId of the gldCfm to query")
	}

	cfmAplyId = strings.ToLower(args[0])

	resultsIterator, err := stub.GetHistoryForKey("BizGldCfm"+ cfmAplyId)                             // 返回迭代器
	if err != nil {
		return shim.Error(err.Error())
	}
	defer resultsIterator.Close()

	if !resultsIterator.HasNext(){
		timestamp := time.Now().Format("2006-01-02 15:04:05")
		fmt.Printf("- can not find gldCfm history by cfmAplyId ! %v \n",timestamp)
		return shim.Success([]byte("can not find gldCfm history by cfmAplyId"))
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
		if response.IsDelete {                                                                      // 若已删除则value为null
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
	fmt.Printf("- query history by cfmAplyId successfully ! %v \n",timestamp)
//	fmt.Printf("- the query result is : \n")
//	fmt.Printf("  %v \n",buffer.String())

	return shim.Success(buffer.Bytes())


}

