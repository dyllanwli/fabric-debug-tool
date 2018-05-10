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
// 金票确认表明细结构体
type BizGldCfmDtl struct {
	ObjectType string     `json:"docType"`             // 类型
	Id         int64      `json:"id"`                  // 主键
	CreateTime int64      `json:"createTime"`          // 创建时间   后台自动生成
	UpdateTime int64      `json:"updateTime"`          // 更新时间
	CreateUser string     `json:"createUser"`          // 创建人
	UpdateUser string     `json:"updateUser"`          // 更新人
	ExpdId     string     `json:"expdId"`              // 扩展ID
	DelInd     string     `json:"delInd"`              // 删除标志
	Version    int32      `json:"version"`             // 版本号
	TenantId   string     `json:"tenantId"`            // 租户ID

	CfmAplyId     string   `json:"cfmAplyId"`      // 确认申请编号
	GldId         string   `json:"gldId"`          // 金票编号
	OriGldId      string   `json:"oriGldId"`       // 原金票编号
	Estb          string   `json:"estb"`           // 开立方
	Pypt          string   `json:"pypt"`           // 支付方
	RcPty         string   `json:"rcPty"`          // 接收方/持票人
	Fctr          string   `json:"fctr"`           // 保理商
	GldAmt        float64  `json:"gldAmt"`         // 金票金额
	GldBal        float64  `json:"gldBal"`         // 金票余额
	PyAmt         float64  `json:"pyAmt"`          // 支付金额
	FncAmt        float64  `json:"fncAmt"`         // 融资金额
	FncBal        float64  `json:"fncBal"`         // 融资余额
	PmoaccTamt    float64  `json:"pmoaccTamt"`     // 垫付总金额
	PmoaccTotBal  float64  `json:"pmoaccTotBal"`   // 垫付总余额
	PymTamt       float64  `json:"pymTamt"`        // 付款总金额
	EstbDay       int64    `json:"estbDay"`        // 开立日
	ExDay         int64    `json:"exDay"`          // 到期日
	SignToacptDay int64    `json:"signToacptDay"`  // 签收日

	// 2017.11.27 新增
	RspbPsnId          string    `json:"rspbPsnId"`     // 经办人编号
	HdlInstId          string    `json:"hdlInstId"`     // 经办机构编号
	HdlDt              int64     `json:"hdlDt"`         // 经办日期（营业日）
}

// 用于解析queryString
type QryStrBizGldCfmDtl struct{
	Selector          BizGldCfmDtl         `json:"selector"`
}

// 判断BizGldCfmDtl变量是否为空
func isEmptyBizGldCfmDtl(arg BizGldCfmDtl) bool {
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

// BizGldCfmDtl变量变为queryString
func tranfBizGldCfmDtlToQryStr(arg BizGldCfmDtl)( string){
	queryString := "{\"selector\":{\"docType\":\"BizGldCfmDtl\""                       // 封装，头部

	value := reflect.ValueOf(arg)
	typ := reflect.TypeOf(arg)
	for i:=1;i<value.NumField();i++{
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
// 保存或更新gldCfmDtl
func saveOrUpdateGldCfmDtl(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	var err error
	fmt.Println("- start saveOrUpdate gldCfmDtl")

	if len(args) < 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}

	gldCfmDtl := args[0]
//	fmt.Println("- the received gldCfmDtl args is :",gldCfmDtl)
	bizGldCfmDtl1 := BizGldCfmDtl{}
	err = json.Unmarshal([]byte(gldCfmDtl), &bizGldCfmDtl1)
	if err != nil {
		return shim.Error(err.Error())
	}

	if bizGldCfmDtl1.GldId == "" {
		return shim.Error("gldId can't be null")
	} else {
		bizGldCfmDtlFromState, err := stub.GetState("BizGldCfmDtl" + bizGldCfmDtl1.GldId)    // worldstate中，key是ObjectType+gldId形式
		if err != nil {
			return shim.Error("Failed to get bizGldCfmDtl:" + err.Error())
		} else if bizGldCfmDtlFromState == nil {
			bizGldCfmDtl1.ObjectType = "BizGldCfmDtl"
			bizGldCfmDtl1.DelInd = "0"
			bizGldCfmDtlToState, err := json.Marshal(bizGldCfmDtl1)
			err = stub.PutState("BizGldCfmDtl" + bizGldCfmDtl1.GldId, bizGldCfmDtlToState)
			if err != nil {
				return shim.Error(err.Error())
			}

			timestamp := time.Now().Format("2006-01-02 15:04:05")
			fmt.Printf("- save successfully ! %v \n",timestamp)
//			fmt.Printf("- the key of record is : BizGldCfmDtl%v \n",bizGldCfmDtl1.GldId)
//			fmt.Printf("- the value of record is ： %v \n",string(bizGldCfmDtlToState))

			return shim.Success(nil)
		} else {
			bizGldCfmDtl2 := BizGldCfmDtl{}
			err = json.Unmarshal([]byte(bizGldCfmDtlFromState), &bizGldCfmDtl2)
			value1 := reflect.ValueOf(&bizGldCfmDtl1).Elem()
			value2 := reflect.ValueOf(&bizGldCfmDtl2).Elem()
			for i:=0; i<value1.NumField(); i++{
				if !isEmpty( value1.Field(i).Interface() ){
					value2.Field(i).Set( value1.Field(i) )
				}
			}
			bizGldCfmDtlToState,err := json.Marshal(bizGldCfmDtl2)
			err = stub.PutState("BizGldCfmDtl" + bizGldCfmDtl1.GldId, bizGldCfmDtlToState)

			if err != nil {
				return shim.Error(err.Error())
			}

			timestamp := time.Now().Format("2006-01-02 15:04:05")
			fmt.Printf("- update successfully ! %v \n",timestamp)
	//		fmt.Printf("- the key of record is : BizGldCfmDtl%v \n",bizGldCfmDtl1.GldId)
	//		fmt.Printf("- the value of record is ： %v \n",string(bizGldCfmDtlToState))

		}
	}
	return shim.Success(nil)
}

// 删除gldCfmDtl
func deleteGldCfmDtl(stub shim.ChaincodeStubInterface, args []string) pb.Response {

	var err error
	var gldId, jsonResp string

	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting gldId of the bizGldCfmDtl to delete")
	}

	gldId = strings.ToLower(args[0])                                                            // 传入的直接就是string格式的gldID

	bizGldCfmDtlFromState, err := stub.GetState("BizGldCfmDtl" + gldId)
	if err != nil {
		jsonResp = "{\"Error\":\"Failed to get state for " + gldId + "\"}"
		return shim.Error(jsonResp)
	} else if bizGldCfmDtlFromState == nil {
		jsonResp = "{\"Error\":\"bizGldCfmDtl does not exist: " + gldId + "\"}"
		return shim.Error(jsonResp)
	}
	bizGldCfmDtl := BizGldCfmDtl{}

	err = json.Unmarshal(bizGldCfmDtlFromState,&bizGldCfmDtl)                                   // worldstate中，数据以k-v形式存在，其中v是json形式

	bizGldCfmDtl.DelInd = "1"
	bizGldCfmDtlToState, err := json.Marshal(bizGldCfmDtl)
	err = stub.PutState("BizGldCfmDtl" + bizGldCfmDtl.GldId, bizGldCfmDtlToState)
	if err != nil {
		return shim.Error(err.Error())
	}

	timestamp := time.Now().Format("2006-01-02 15:04:05")
	fmt.Printf("- delete successfully ! %v \n",timestamp)
//	fmt.Printf("- the key of record is : BizGldCfmDtl%v \n",bizGldCfmDtl.GldId)
//	fmt.Printf("- the value of record is ： %v \n",string(bizGldCfmDtlToState))

	return shim.Success(nil)
}


// 通过gldId查询gldCfmDtl
func queryGldCfmDtlByGldId(stub shim.ChaincodeStubInterface, args []string) pb.Response {

	var err error
	var gldId, jsonResp string

	fmt.Println("- start query gldCfmDtl by gldId")

	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting gldId of the gldCfmDtl to query")
	}

	gldId = strings.ToLower(args[0])

	bizGldCfmDtlFromState, err := stub.GetState("BizGldCfmDtl" + gldId)

	if err != nil {
		jsonResp = "{\"Error\":\"Failed to get state for " + gldId + "\"}"
		return shim.Error(jsonResp)
	} else if bizGldCfmDtlFromState == nil {
		jsonResp = "{\"Error\":\"bizGldCfmDtl does not exist: " + gldId + "\"}"
		return shim.Error(jsonResp)
	}

	var buffer bytes.Buffer                                                                         // 缓存，用来保存查询结果
	buffer.WriteString("{\"Key\":")
	buffer.WriteString("\"BizGldCfmDtl")
	buffer.WriteString(gldId)
	buffer.WriteString("\"")
	buffer.WriteString(", \"Value\":")
	buffer.WriteString(string(bizGldCfmDtlFromState))
	buffer.WriteString("}")

	timestamp := time.Now().Format("2006-01-02 15:04:05")
	fmt.Printf("- query by gldId successfully ! %v \n",timestamp)
//	fmt.Printf("- the query result is : \n")
//	fmt.Printf("  %v \n",buffer.String())

	return shim.Success(buffer.Bytes())

}

// 通过gldId范围查询gldCfmDtl(复杂查询)
func queryGldCfmDtlByGldIdRange(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	var err error
	var startGldId,endGldId string

	fmt.Println("- start query gldCfmDtl by gldId range")

	if len(args)<2{
		return shim.Error("Incorrect number of arguments. Range query expect startGldId and endGldId of the gldCfmDtl to query")
	}
	startGldId = strings.ToLower(args[0])
	endGldId = strings.ToLower(args[1])

	startKey := "BizGldCfmDtl" + startGldId
	endKey := "BizGldCfmDtl" + endGldId

	resultsIterator, err := stub.GetStateByRange(startKey, endKey)
	if err != nil {
		return shim.Error(err.Error())
	}
	defer resultsIterator.Close()                                                                // 迭代器用完需关闭！！

	if !resultsIterator.HasNext(){
		timestamp := time.Now().Format("2006-01-02 15:04:05")
		fmt.Printf("- can not find gldCfmDtl by gldId range ! %v \n",timestamp)
		return shim.Success([]byte("can not find gldCfmDtl by gldId range"))
	}
	var buffer bytes.Buffer
	buffer.WriteString("[")                                                                   // 对查询结果进行封装，封装为对象数组json格式
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

	return shim.Success(buffer.Bytes())                                                            // 最终，结果以[]byte形式返回


}

// 通过querystring查询gldCfmDtl(复杂查询)
func queryGldCfmDtlByQryStr(stub shim.ChaincodeStubInterface, args []string) pb.Response {

	var err error
	var queryString string

	//	fmt.Println("- start query gldAccvalRltv by queryString, the queryString is : \n %v \n",queryString )
	fmt.Println("- start query gldCfmDtl by queryString" )

	if len(args)!= 1{
		return shim.Error("Incorrect number of arguments. Expecting querystring of the gldCfmDtl to query")
	}
	queryString = args[0]

	if queryString == ""{                                                                         // 判定querystring是否为空
		return shim.Error("The queryString can not be null ")
	}
	qryStrBizGldCfmDtl := QryStrBizGldCfmDtl{}
	err = json.Unmarshal([]byte(queryString), &qryStrBizGldCfmDtl)
	if err != nil {
		return shim.Error(err.Error())
	}
	if isEmptyBizGldCfmDtl(qryStrBizGldCfmDtl.Selector){                                          // 判断queryString中有没有查询条件
		return shim.Error("The selector of queryString can not be null ")
	}

	querystring := tranfBizGldCfmDtlToQryStr(qryStrBizGldCfmDtl.Selector)

	resultsIterator, err := stub.GetQueryResult(querystring)                                      // 返回迭代器，querystring格式：{"selector":{"docType":"BizGldAccvalRltv","CreateUser":"jzk",...}}
	if err != nil {
		return shim.Error(err.Error())
	}
	defer resultsIterator.Close()

	if !resultsIterator.HasNext(){
		timestamp := time.Now().Format("2006-01-02 15:04:05")
		fmt.Printf("- can not find gldCfmDtl by querystring ! %v \n",timestamp)
		return shim.Success([]byte("can not find gldCfmDtl by querystring"))
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

// 通过condition查询gldCfmDtl(复杂查询)
func queryGldCfmDtlByCond(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	var err error

	fmt.Println("- start query gldCfmDtl by conditions" )

	if len(args) < 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}
	gldCfmDtl := args[0]                                                                             // 查询条件以json串传递，格式为：{"docType":"BizGldAccvalRltv","id":"id0001","createTime":"20171110",...}
	bizGldCfmDtl := BizGldCfmDtl{}
	err = json.Unmarshal([]byte(gldCfmDtl), &bizGldCfmDtl)
	if err != nil {
		return shim.Error(err.Error())
	}

	if isEmptyBizGldCfmDtl(bizGldCfmDtl){
		return shim.Error("The query conditions can not be null")
	}

	queryString := tranfBizGldCfmDtlToQryStr(bizGldCfmDtl)
	resultsIterator, err := stub.GetQueryResult(queryString)                                         // 返回迭代器
	if err != nil {
		return shim.Error(err.Error())
	}
	defer resultsIterator.Close()

	if !resultsIterator.HasNext(){
		timestamp := time.Now().Format("2006-01-02 15:04:05")
		fmt.Printf("- can not find gldCfmDtl by querystring ! %v \n",timestamp)
		return shim.Success([]byte("can not find gldCfmDtl by conditions"))
	}
	var buffer bytes.Buffer
	buffer.WriteString("[")                                                                      // 封装成对象数组json串格式
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

// 通过gldId查询gldCfmDtl的历史（复杂查询）
func queryGldCfmDtlHsyByGldId(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	var err error
	var gldId string

	fmt.Println("- start query gldCfmDtl history by gldId")

	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting gldId of the gldCfmDtl to query")
	}

	gldId = strings.ToLower(args[0])

	resultsIterator, err := stub.GetHistoryForKey("BizGldCfmDtl"+ gldId)                            // 返回迭代器
	if err != nil {
		return shim.Error(err.Error())
	}
	defer resultsIterator.Close()

	if !resultsIterator.HasNext(){
		timestamp := time.Now().Format("2006-01-02 15:04:05")
		fmt.Printf("- can not find gldCfmDtl history by gldId ! %v \n",timestamp)
		return shim.Success([]byte("can not find gldCfmDtl history by gldId"))
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
		if response.IsDelete {                                                                     // 若已删除则value为null
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