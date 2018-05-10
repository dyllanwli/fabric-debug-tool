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
// 金票垫付申请结构体
type BizGldPmoaccAply struct{
	ObjectType string     `json:"docType"`             // 类型
	Id         int64      `json:"id"`                  // 主键
	CreateTime int64     `json:"createTime"`          // 创建时间
	UpdateTime int64     `json:"updateTime"`          // 更新时间
	CreateUser string     `json:"createUser"`          // 创建人
	UpdateUser string     `json:"updateUser"`          // 更新人
	ExpdId     string     `json:"expdId"`              // 扩展ID
	DelInd     string     `json:"delInd"`              // 删除标志
	Version    int32      `json:"version"`             // 版本号
	TenantId   string     `json:"tenantId"`            // 租户ID

	GldId           string   `json:"gldId"`            // 金票编号
	PmoaccAplyPsn   string   `json:"pmoaccAplyPsn"`    // 垫付申请人
	Pmoacc          string   `json:"pmoacc"`           // 垫付方
	PmoaccAmt       float64  `json:"pmoaccAmt"`        // 垫付金额
	PmoaccDt        int64    `json:"pmoaccDt"`         // 垫付日期
	PmoaccAplyDt    int64    `json:"pmoaccAplyDt"`     // 垫付申请日期
	PcsgStCd        string   `json:"pcsgStCd"`         // 处理状态/处理流程

	// 2017.11.27 新增
	Opin               string    `json:"opin"`          // 意见
}

// 用于解析queryString
type QryStrBizGldPmoaccAply struct{
	Selector          BizGldPmoaccAply         `json:"selector"`
}

// 判断BizGldAccvalRltv变量是否为空
func isEmptyBizGldPmoaccAply(arg BizGldPmoaccAply) bool {
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
func tranfBizGldPmoaccAplyToQryStr(arg BizGldPmoaccAply)( string){
	queryString := "{\"selector\":{\"docType\":\"BizGldPmoaccAply\""                  // 封装，头部

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
// 保存或更新gldPmoaccAply
func saveOrUpdateGldPmoaccAply(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	var err error
	fmt.Println("- start saveOrUpdate gldPmoaccAply")

	if len(args) < 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}

	gldPmoaccAply := args[0]
//	fmt.Println("- the received gldPmoaccAply args is :",gldPmoaccAply)
	bizGldPmoaccAply1 := BizGldPmoaccAply{}
	err = json.Unmarshal([]byte(gldPmoaccAply), &bizGldPmoaccAply1)
	if err != nil {
		return shim.Error(err.Error())
	}

	if bizGldPmoaccAply1.GldId == "" {
		return shim.Error("gldId can't be null")
	} else {
		bizGldPmoaccAplyFromState, err := stub.GetState("BizGldPmoaccAply" + bizGldPmoaccAply1.GldId)
		if err != nil {
			return shim.Error("Failed to get bizGldPmoaccAply:" + err.Error())
		} else if bizGldPmoaccAplyFromState == nil {
			bizGldPmoaccAply1.ObjectType = "BizGldPmoaccAply"
			bizGldPmoaccAply1.DelInd = "0"                                                      // 删除标志，初始值为0
			bizGldPmoaccAplyToState, err := json.Marshal(bizGldPmoaccAply1)
			err = stub.PutState("BizGldPmoaccAply" +bizGldPmoaccAply1.GldId, bizGldPmoaccAplyToState)
			if err != nil {
				return shim.Error(err.Error())
			}

			timestamp := time.Now().Format("2006-01-02 15:04:05")
			fmt.Printf("- save successfully ! %v \n",timestamp)
	//		fmt.Printf("- the key of record is : BizGldPmoaccAply%v \n",bizGldPmoaccAply1.GldId)
	//		fmt.Printf("- the value of record is ： %v \n",string(bizGldPmoaccAplyToState))

			return shim.Success(nil)
		} else {
			bizGldPmoaccAply2 := BizGldPmoaccAply{}
			err = json.Unmarshal([]byte(bizGldPmoaccAplyFromState), &bizGldPmoaccAply2)
			value1 := reflect.ValueOf(&bizGldPmoaccAply1).Elem()
			value2 := reflect.ValueOf(&bizGldPmoaccAply2).Elem()
			for i:=0; i<value1.NumField(); i++{
				if !isEmpty( value1.Field(i).Interface() ){
					value2.Field(i).Set( value1.Field(i) )
				}
			}
			bizGldPmoaccAplyToState,err := json.Marshal(bizGldPmoaccAply2)
			err = stub.PutState("BizGldPmoaccAply" + bizGldPmoaccAply1.GldId, bizGldPmoaccAplyToState)
			if err != nil {
				return shim.Error(err.Error())
			}

			timestamp := time.Now().Format("2006-01-02 15:04:05")
			fmt.Printf("- update successfully ! %v \n",timestamp)
	//		fmt.Printf("- the key of record is : BizGldPmoaccAply%v \n",bizGldPmoaccAply1.GldId)
	//		fmt.Printf("- the value of record is ： %v \n",string(bizGldPmoaccAplyToState))

		}
	}
	return shim.Success(nil)
}

// 删除gldPmoaccAply
func deleteGldPmoaccAply(stub shim.ChaincodeStubInterface, args []string) pb.Response {

	var err error
	var gldId, jsonResp string

	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting gldId of the gldPmoaccAply to delete")
	}

	gldId = strings.ToLower(args[0])

	bizGldPmoaccAplyFromState, err := stub.GetState("BizGldPmoaccAply" + gldId)
	if err != nil {
		jsonResp = "{\"Error\":\"Failed to get state for " + gldId + "\"}"
		return shim.Error(jsonResp)
	} else if bizGldPmoaccAplyFromState == nil {
		jsonResp = "{\"Error\":\"bizGldPmoaccAply does not exist: " + gldId + "\"}"
		return shim.Error(jsonResp)
	}
	bizGldPmoaccAply := BizGldPmoaccAply{}

	err = json.Unmarshal(bizGldPmoaccAplyFromState,&bizGldPmoaccAply)

	bizGldPmoaccAply.DelInd = "1"
	bizGldPmoaccAplyToState, err := json.Marshal(bizGldPmoaccAply)
	err = stub.PutState("BizGldPmoaccAply" + bizGldPmoaccAply.GldId, bizGldPmoaccAplyToState)         // 删除并非真的从数据中删除记录，只是将删除标志位置位而已
	if err != nil {
		return shim.Error(err.Error())
	}

	timestamp := time.Now().Format("2006-01-02 15:04:05")
	fmt.Printf("- delete successfully ! %v \n",timestamp)
//	fmt.Printf("- the key of record is : BizGldPmoaccAply%v \n",bizGldPmoaccAply.GldId)
//	fmt.Printf("- the value of record is ： %v \n",string(bizGldPmoaccAplyToState))

	return shim.Success(nil)
}


// 通过gldId查询gldPmoaccAply
func queryGldPmoaccAplyByGldId(stub shim.ChaincodeStubInterface, args []string) pb.Response {

	var err error
	var gldId, jsonResp string

	fmt.Println("- start query gldPmoaccAply by gldId")

	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting gldId of the gldPmoaccAply to query")
	}

	gldId = strings.ToLower(args[0])

	bizGldPmoaccAplyFromState, err := stub.GetState("BizGldPmoaccAply" + gldId)

	if err != nil {
		jsonResp = "{\"Error\":\"Failed to get state for " + gldId + "\"}"
		return shim.Error(jsonResp)
	} else if bizGldPmoaccAplyFromState == nil {
		jsonResp = "{\"Error\":\"bizGldPmoaccAply does not exist: " + gldId + "\"}"
		return shim.Error(jsonResp)
	}


	var buffer bytes.Buffer                                                                     // 缓存，用来保存查询结果
	buffer.WriteString("{\"Key\":")
	buffer.WriteString("\"BizGldPmoaccAply")
	buffer.WriteString(gldId)
	buffer.WriteString("\"")
	buffer.WriteString(", \"Value\":")
	buffer.WriteString(string(bizGldPmoaccAplyFromState))
	buffer.WriteString("}")

	timestamp := time.Now().Format("2006-01-02 15:04:05")
	fmt.Printf("- query by gldId successfully ! %v \n",timestamp)
//	fmt.Printf("- the query result is : \n")
//	fmt.Printf("  %v \n",buffer.String())

	return shim.Success(buffer.Bytes())
}



// 通过gldId范围查询gldPmoaccAply(复杂查询)
func queryGldPmoaccAplyByGldIdRange(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	var err error
	var startGldId,endGldId string

	fmt.Println("- start query gldPmoaccAply by gldId range")

	if len(args)<2{
		return shim.Error("Incorrect number of arguments. Range query expect startGldId and endGldId of the gldPmoaccAply to query")
	}
	startGldId = strings.ToLower(args[0])
	endGldId = strings.ToLower(args[1])

	startKey := "BizGldPmoaccAply" + startGldId
	endKey := "BizGldPmoaccAply" + endGldId

	resultsIterator, err := stub.GetStateByRange(startKey, endKey)
	if err != nil {
		return shim.Error(err.Error())
	}
	defer resultsIterator.Close()                                            // 迭代器用完需关闭！！

	if !resultsIterator.HasNext(){
		timestamp := time.Now().Format("2006-01-02 15:04:05")
		fmt.Printf("- can not find gldPmoaccAply by gldId range ! %v \n",timestamp)
		return shim.Success([]byte("can not find gldPmoaccAply by gldId range"))
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

// 通过querystring查询gldPmoaccAply(复杂查询)
func queryGldPmoaccAplyByQryStr(stub shim.ChaincodeStubInterface, args []string) pb.Response {

	var err error
	var queryString string

	//	fmt.Println("- start query gldAccvalRltv by queryString, the queryString is : \n %v \n",queryString )
	fmt.Println("- start query gldPmoaccAply by queryString" )

	if len(args)!= 1{
		return shim.Error("Incorrect number of arguments. Expecting querystring of the gldPmoaccAply to query")
	}
	queryString = args[0]

	if queryString == ""{                                                           // 判定querystring是否为空
		return shim.Error("The queryString can not be null ")
	}
	qryStrBizGldPmoaccAply := QryStrBizGldPmoaccAply{}
	err = json.Unmarshal([]byte(queryString), &qryStrBizGldPmoaccAply)
	if err != nil {
		return shim.Error(err.Error())
	}
	if isEmptyBizGldPmoaccAply(qryStrBizGldPmoaccAply.Selector){                    // 判断queryString中有没有查询条件
		return shim.Error("The selector of queryString can not be null ")
	}

	querystring := tranfBizGldPmoaccAplyToQryStr(qryStrBizGldPmoaccAply.Selector)

	resultsIterator, err := stub.GetQueryResult(querystring)              // 返回迭代器，querystring格式：{"selector":{"docType":"BizGldAccvalRltv","CreateUser":"jzk",...}}
	if err != nil {
		return shim.Error(err.Error())
	}
	defer resultsIterator.Close()

	if !resultsIterator.HasNext(){
		timestamp := time.Now().Format("2006-01-02 15:04:05")
		fmt.Printf("- can not find gldPmoaccAply by querystring ! %v \n",timestamp)
		return shim.Success([]byte("can not find gldPmoaccAply by querystring"))
	}
	var buffer bytes.Buffer
	buffer.WriteString("[")                                                        // 封装成对象数组json串格式
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

// 通过condition查询gldPmoaccAply(复杂查询)
func queryGldPmoaccAplyByCond(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	var err error

	fmt.Println("- start query gldPmoaccAply by conditions" )

	if len(args) < 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}
	gldPmoaccAply := args[0]                                                        // 查询条件以json串传递，格式为：{"docType":"BizGldAccvalRltv","id":"id0001","createTime":"20171110",...}
	bizGldPmoaccAply := BizGldPmoaccAply{}
	err = json.Unmarshal([]byte(gldPmoaccAply), &bizGldPmoaccAply)
	if err != nil {
		return shim.Error(err.Error())
	}

	if isEmptyBizGldPmoaccAply(bizGldPmoaccAply){
		return shim.Error("The query conditions can not be null")
	}

	queryString := tranfBizGldPmoaccAplyToQryStr(bizGldPmoaccAply)
	resultsIterator, err := stub.GetQueryResult(queryString)                        // 返回迭代器
	if err != nil {
		return shim.Error(err.Error())
	}
	defer resultsIterator.Close()

	if !resultsIterator.HasNext(){
		timestamp := time.Now().Format("2006-01-02 15:04:05")
		fmt.Printf("- can not find gldPmoaccAply by querystring ! %v \n",timestamp)
		return shim.Success([]byte("can not find gldPmoaccAply by conditions"))
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
	fmt.Printf("- query by conditions successfully ! %v \n",timestamp)
//	fmt.Printf("- the query result is : \n")
//	fmt.Printf("  %v \n",buffer.String())

	return shim.Success(buffer.Bytes())


}

// 通过gldId查询gldPmoaccAply的历史（复杂查询）
func queryGldPmoaccAplyHsyByGldId(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	var err error
	var gldId string

	fmt.Println("- start query gldPmoaccAply history by gldId")

	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting gldId of the gldPmoaccAply to query")
	}

	gldId = strings.ToLower(args[0])

	resultsIterator, err := stub.GetHistoryForKey("BizGldPmoaccAply"+ gldId)           // 返回迭代器
	if err != nil {
		return shim.Error(err.Error())
	}
	defer resultsIterator.Close()

	if !resultsIterator.HasNext(){
		timestamp := time.Now().Format("2006-01-02 15:04:05")
		fmt.Printf("- can not find gldPmoaccAply history by gldId ! %v \n",timestamp)
		return shim.Success([]byte("can not find gldPmoaccAply history by gldId"))
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
		if response.IsDelete {                                                           // 若已删除则value为null
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



