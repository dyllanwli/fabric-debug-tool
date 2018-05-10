
package main
import (
	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
	"encoding/json"
	"fmt"
	"time"
	"strings"
	"reflect"
	"bytes"
	"strconv"
)


/**************************************************************************************/
// 金票单据关联结构体
type BizGldDocRltv struct{
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

	GldId           string       `json:"gldId"`        // 金票编号
	RcvbId          string       `json:"rcvbId"`       // 应收账款编号
	EstbAmt         float64      `json:"estbAmt"`      // 开立金额
}

// 用于解析queryString
type QryStrBizGldDocRltv struct{
	Selector          BizGldDocRltv         `json:"selector"`
}

// 判断BizGldAccvalRltv变量是否为空
func isEmptyBizGldDocRltv(arg BizGldDocRltv) bool {
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

// BizGldAccvalRltv变量变为queryString
func tranfBizGldDocRltvToQryStr(arg BizGldDocRltv)( string){
	queryString := "{\"selector\":{\"docType\":\"BizGldDocRltv\""                     // 封装，头部

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
// 保存gldDocRltv
func saveGldDocRltv(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	var err error
	fmt.Println("- start save gldDocRltv")

	if len(args) < 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}
	gldDocRltv := args[0]
//	fmt.Println("- the received gldDocRltv args is :",gldDocRltv)
	bizGldDocRltv1 := BizGldDocRltv{}
	err = json.Unmarshal([]byte(gldDocRltv), &bizGldDocRltv1)
	if err != nil {
		return shim.Error(err.Error())
	}

	if bizGldDocRltv1.GldId == "" {
		return shim.Error("gldId can't be null")
	} else {
		bizGldDocRltv1.ObjectType = "BizGldDocRltv"
		bizGldDocRltv1.DelInd = "0"
		bizGldDocRltvToState, err := json.Marshal(bizGldDocRltv1)
		err = stub.PutState("BizGldDocRltv" + bizGldDocRltv1.GldId, bizGldDocRltvToState)
		if err != nil {
			return shim.Error(err.Error())
		}

		timestamp := time.Now().Format("2006-01-02 15:04:05")
		fmt.Printf("- save successfully ! %v \n",timestamp)
//		fmt.Printf("- the key of record is : BizGldDocRltv%v \n",bizGldDocRltv1.GldId)
//		fmt.Printf("- the value of record is ： %v \n",string(bizGldDocRltvToState))

		return shim.Success(nil)
	}
}


// 通过gldId查询gldDocRltv
func queryGldDocRltvByGldId(stub shim.ChaincodeStubInterface, args []string) pb.Response {

	var err error
	var gldId, jsonResp string

	fmt.Println("- start query gldDocRltv by gldId")

	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting gldId of the gldDocRltv to query")
	}

	gldId = strings.ToLower(args[0])

	bizGldDocRltvFromState, err := stub.GetState("BizGldDocRltv" + gldId)

	if err != nil {
		jsonResp = "{\"Error\":\"Failed to get state for " + gldId + "\"}"
		return shim.Error(jsonResp)
	} else if bizGldDocRltvFromState == nil {
		jsonResp = "{\"Error\":\"bizGldDocRltv does not exist: " + gldId + "\"}"
		return shim.Error(jsonResp)
	}

	var buffer bytes.Buffer                                                            // 缓存，用来保存查询结果
	buffer.WriteString("{\"Key\":")
	buffer.WriteString("\"BizGldDocRltv")
	buffer.WriteString(gldId)
	buffer.WriteString("\"")
	buffer.WriteString(", \"Value\":")
	buffer.WriteString(string(bizGldDocRltvFromState))
	buffer.WriteString("}")

	timestamp := time.Now().Format("2006-01-02 15:04:05")
	fmt.Printf("- query by gldId successfully ! %v \n",timestamp)
//	fmt.Printf("- the query result is : \n")
//	fmt.Printf("  %v \n",buffer.String())

	return shim.Success(buffer.Bytes())
}

// 通过gldId范围查询gldDocRltv(复杂查询)
func queryGldDocRltvByGldIdRange(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	var err error
	var startGldId,endGldId string

	fmt.Println("- start query gldDocRltv by gldId range")

	if len(args)<2{
		return shim.Error("Incorrect number of arguments. Range query expect startGldId and endGldId of the gldDocRltv to query")
	}
	startGldId = strings.ToLower(args[0])
	endGldId = strings.ToLower(args[1])

	startKey := "BizGldDocRltv" + startGldId
	endKey := "BizGldDocRltv" + endGldId

	resultsIterator, err := stub.GetStateByRange(startKey, endKey)
	if err != nil {
		return shim.Error(err.Error())
	}
	defer resultsIterator.Close()                                            // 迭代器用完需关闭！！

	if !resultsIterator.HasNext(){
		timestamp := time.Now().Format("2006-01-02 15:04:05")
		fmt.Printf("- can not find gldDocRltv by gldId range ! %v \n",timestamp)
		return shim.Success([]byte("can not find gldDocRltv by gldId range"))
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

// 通过querystring查询gldDocRltv(复杂查询)
func queryGldDocRltvByQryStr(stub shim.ChaincodeStubInterface, args []string) pb.Response {

	var err error
	var queryString string

	//	fmt.Println("- start query gldAccvalRltv by queryString, the queryString is : \n %v \n",queryString )
	fmt.Println("- start query gldDocRltv by queryString" )

	if len(args)!= 1{
		return shim.Error("Incorrect number of arguments. Expecting querystring of the gldDocRltv to query")
	}
	queryString = args[0]

	if queryString == ""{                                                     // 判定querystring是否为空
		return shim.Error("The queryString can not be null ")
	}
	qryStrBizGldDocRltv := QryStrBizGldDocRltv{}
	err = json.Unmarshal([]byte(queryString), &qryStrBizGldDocRltv)
	if err != nil {
		return shim.Error(err.Error())
	}
	if isEmptyBizGldDocRltv(qryStrBizGldDocRltv.Selector){                    // 判断queryString中有没有查询条件
		return shim.Error("The selector of queryString can not be null ")
	}

	querystring := tranfBizGldDocRltvToQryStr(qryStrBizGldDocRltv.Selector)

	resultsIterator, err := stub.GetQueryResult(querystring)                  // 返回迭代器，querystring格式：{"selector":{"docType":"BizGldAccvalRltv","CreateUser":"jzk",...}}
	if err != nil {
		return shim.Error(err.Error())
	}
	defer resultsIterator.Close()

	if !resultsIterator.HasNext(){
		timestamp := time.Now().Format("2006-01-02 15:04:05")
		fmt.Printf("- can not find gldDocRltv by querystring ! %v \n",timestamp)
		return shim.Success([]byte("can not find gldDocRltv by querystring"))
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

// 通过condition查询gldDocRltv(复杂查询)
func queryGldDocRltvByCond(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	var err error

	fmt.Println("- start query gldDocRltv by conditions" )

	if len(args) < 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}
	gldDocRltv := args[0]                                                         // 查询条件以json串传递，格式为：{"docType":"BizGldAccvalRltv","id":"id0001","createTime":"20171110",...}
	bizGldDocRltv := BizGldDocRltv{}
	err = json.Unmarshal([]byte(gldDocRltv), &bizGldDocRltv)
	if err != nil {
		return shim.Error(err.Error())
	}

	if isEmptyBizGldDocRltv(bizGldDocRltv){
		return shim.Error("The query conditions can not be null")
	}

	queryString := tranfBizGldDocRltvToQryStr(bizGldDocRltv)
	resultsIterator, err := stub.GetQueryResult(queryString)                     // 返回迭代器
	if err != nil {
		return shim.Error(err.Error())
	}
	defer resultsIterator.Close()

	if !resultsIterator.HasNext(){
		timestamp := time.Now().Format("2006-01-02 15:04:05")
		fmt.Printf("- can not find gldDocRltv by querystring ! %v \n",timestamp)
		return shim.Success([]byte("can not find gldDocRltv by conditions"))
	}
	var buffer bytes.Buffer
	buffer.WriteString("[")                                                  // 封装成对象数组json串格式
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

// 通过gldId查询gldDocRltv的历史（复杂查询）
func queryGldDocRltvHsyByGldId(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	var err error
	var gldId string

	fmt.Println("- start query gldDocRltv history by gldId")

	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting gldId of the gldDocRltv to query")
	}

	gldId = strings.ToLower(args[0])

	resultsIterator, err := stub.GetHistoryForKey("BizGldDocRltv"+ gldId)                  // 返回迭代器
	if err != nil {
		return shim.Error(err.Error())
	}
	defer resultsIterator.Close()

	if !resultsIterator.HasNext(){
		timestamp := time.Now().Format("2006-01-02 15:04:05")
		fmt.Printf("- can not find gldDocRltv history by gldId ! %v \n",timestamp)
		return shim.Success([]byte("can not find gldDocRltv history by gldId"))
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