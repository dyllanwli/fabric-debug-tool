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
// 金票垫付申请状态结构体
type BizGldPmoaccAplyInf struct{
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

	PmoaccAplyId           string   `json:"pmoaccAplyId"`   // 垫付申请编号
	PcsStCd                string   `json:"pcsStCd"`        // 流程状态
	Opin                   string   `json:"opin"`           // 意见

	// 2017.11.27 新增
	RspbPsnId          string    `json:"rspbPsnId"`     // 经办人编号
	HdlInstId          string    `json:"hdlInstId"`     // 经办机构编号
	HdlDt              int64     `json:"hdlDt"`         // 经办日期（营业日）
}

// 用于解析queryString
type QryStrBizGldPmoaccAplyInf struct{
	Selector          BizGldPmoaccAplyInf         `json:"selector"`
}

// 判断BizGldAccvalRltv变量是否为空
func isEmptyBizGldPmoaccAplyInf(arg BizGldPmoaccAplyInf) bool {
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
func tranfBizGldPmoaccAplyInfToQryStr(arg BizGldPmoaccAplyInf)( string){
	queryString := "{\"selector\":{\"docType\":\"BizGldPmoaccAplyInf\""                  // 封装，头部

	value := reflect.ValueOf(arg)
	typ := reflect.TypeOf(arg)
	for i:=1;i<value.NumField();i++{                                                    // Field(0)跳过
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
// 保存或更新gldPmoaccAplyInf
func saveOrUpdateGldPmoaccAplyInf(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	var err error
	fmt.Println("- start saveOrUpdate gldPmoaccAplyInf")

	if len(args) < 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}

	gldPmoaccAplyInf := args[0]
//	fmt.Println("- the received gldPmoaccAplyInf args is :",gldPmoaccAplyInf)
	bizGldPmoaccAplyInf1 := BizGldPmoaccAplyInf{}
	err = json.Unmarshal([]byte(gldPmoaccAplyInf), &bizGldPmoaccAplyInf1)
	if err != nil {
		return shim.Error(err.Error())
	}

	if bizGldPmoaccAplyInf1.PmoaccAplyId == "" {
		return shim.Error("pmoaccAplyId can't be null")
	} else {
		bizGldPmoaccAplyInfFromState, err := stub.GetState("BizGldPmoaccAplyInf" + bizGldPmoaccAplyInf1.PmoaccAplyId)
		if err != nil {
			return shim.Error("Failed to get bizGldPmoaccAplyInf:" + err.Error())
		} else if bizGldPmoaccAplyInfFromState == nil {
			bizGldPmoaccAplyInf1.ObjectType = "BizGldPmoaccAplyInf"
			bizGldPmoaccAplyInf1.DelInd = "0"
			bizGldPmoaccAplyInfToState, err := json.Marshal(bizGldPmoaccAplyInf1)
			err = stub.PutState("BizGldPmoaccAplyInf" + bizGldPmoaccAplyInf1.PmoaccAplyId, bizGldPmoaccAplyInfToState)
			if err != nil {
				return shim.Error(err.Error())
			}

			timestamp := time.Now().Format("2006-01-02 15:04:05")
			fmt.Printf("- save successfully ! %v \n",timestamp)
	//		fmt.Printf("- the key of record is : BizGldPmoaccAplyInf%v \n",bizGldPmoaccAplyInf1.PmoaccAplyId)
	//		fmt.Printf("- the value of record is ： %v \n",string(bizGldPmoaccAplyInfToState))

			return shim.Success(nil)
		} else {
			bizGldPmoaccAplyInf2 := BizGldPmoaccAplyInf{}
			err = json.Unmarshal([]byte(bizGldPmoaccAplyInfFromState), &bizGldPmoaccAplyInf2)
			value1 := reflect.ValueOf(&bizGldPmoaccAplyInf1).Elem()
			value2 := reflect.ValueOf(&bizGldPmoaccAplyInf2).Elem()
			for i:=0; i<value1.NumField(); i++{
				if !isEmpty( value1.Field(i).Interface() ){
					value2.Field(i).Set( value1.Field(i) )
				}
			}
			bizGldPmoaccAplyInfToState,err := json.Marshal(bizGldPmoaccAplyInf2)
			err = stub.PutState("BizGldPmoaccAplyInf" + bizGldPmoaccAplyInf1.PmoaccAplyId, bizGldPmoaccAplyInfToState)

			if err != nil {
				return shim.Error(err.Error())
			}

			timestamp := time.Now().Format("2006-01-02 15:04:05")
			fmt.Printf("update successfully ! %v \n",timestamp)
	//		fmt.Printf("the key is : BizGldPmoaccAplyInf %v \n",bizGldPmoaccAplyInf1.PmoaccAplyId)
	//		fmt.Printf("the value is ： %v \n",string(bizGldPmoaccAplyInfToState))

		}
	}
	return shim.Success(nil)
}

// 删除gldPmoaccAplyInf
func deleteGldPmoaccAplyInf(stub shim.ChaincodeStubInterface, args []string) pb.Response {

	var err error
	var pmoaccAplyId, jsonResp string

	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting pmoaccAplyId of the bizGldPmoaccAplyInf to delete")
	}

	pmoaccAplyId = strings.ToLower(args[0])

	bizGldPmoaccAplyInfFromState, err := stub.GetState("BizGldPmoaccAplyInf" + pmoaccAplyId)
	if err != nil {
		jsonResp = "{\"Error\":\"Failed to get state for " + pmoaccAplyId + "\"}"
		return shim.Error(jsonResp)
	} else if bizGldPmoaccAplyInfFromState == nil {
		jsonResp = "{\"Error\":\"bizGldPmoaccAplyInf does not exist: " + pmoaccAplyId + "\"}"
		return shim.Error(jsonResp)
	}
	bizGldPmoaccAplyInf := BizGldPmoaccAplyInf{}

	err = json.Unmarshal(bizGldPmoaccAplyInfFromState,&bizGldPmoaccAplyInf)
	bizGldPmoaccAplyInf.DelInd = "1"
	bizGldPmoaccAplyInfToState, err := json.Marshal(bizGldPmoaccAplyInf)
	err = stub.PutState("BizGldPmoaccAplyInf" + bizGldPmoaccAplyInf.PmoaccAplyId, bizGldPmoaccAplyInfToState)
	if err != nil {
		return shim.Error(err.Error())
	}

	timestamp := time.Now().Format("2006-01-02 15:04:05")
	fmt.Printf("- delete successfully ! %v \n",timestamp)
//	fmt.Printf("- the key of record is : BizGldPmoaccAplyInf%v \n",bizGldPmoaccAplyInf.PmoaccAplyId)
//	fmt.Printf("- the value of record is ： %v \n",string(bizGldPmoaccAplyInfToState))

	return shim.Success(nil)
}


// 通过pmoaccAplyId查询gldPmoaccAplyInf
func queryGldPmoaccAplyInfByPmoaccAplyId(stub shim.ChaincodeStubInterface, args []string) pb.Response {

	var err error
	var pmoaccAplyId, jsonResp string


	fmt.Println("- start query gldPmoaccAplyInf by pmoaccAplyId")

	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting pmoaccAplyId of the gldPmoaccAplyInf to query")
	}

	pmoaccAplyId = strings.ToLower(args[0])

	bizGldPmoaccAplyInfFromState, err := stub.GetState("BizGldPmoaccAplyInf" + pmoaccAplyId)

	if err != nil {
		jsonResp = "{\"Error\":\"Failed to get state for " + pmoaccAplyId + "\"}"
		return shim.Error(jsonResp)
	} else if bizGldPmoaccAplyInfFromState == nil {
		jsonResp = "{\"Error\":\"bizGldPmoaccAplyInf does not exist: " + pmoaccAplyId + "\"}"
		return shim.Error(jsonResp)
	}


	var buffer bytes.Buffer                                                                  // 缓存，用来保存查询结果
	buffer.WriteString("{\"Key\":")
	buffer.WriteString("\"BizGldPmoaccAplyInf")
	buffer.WriteString(pmoaccAplyId)
	buffer.WriteString("\"")
	buffer.WriteString(", \"Value\":")
	buffer.WriteString(string(bizGldPmoaccAplyInfFromState))
	buffer.WriteString("}")

	timestamp := time.Now().Format("2006-01-02 15:04:05")
	fmt.Printf("- query by gldId successfully ! %v \n",timestamp)
//	fmt.Printf("- the query result is : \n")
//	fmt.Printf("  %v \n",buffer.String())

	return shim.Success(buffer.Bytes())
}



// 通过pmoaccAplyId范围查询gldPmoaccAplyInf(复杂查询)
func queryGldPmoaccAplyInfByPmoaccAplyIdRange(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	var err error
	var startPmoaccAplyId,endPmoaccAplyId string

	fmt.Println("- start query gldPmoaccAplyInf by pmoaccAplyId range")

	if len(args)<2{
		return shim.Error("Incorrect number of arguments. Range query expect startPmoaccAplyId and endPmoaccAplyId of the gldPmoaccAplyInf to query")
	}
	startPmoaccAplyId = strings.ToLower(args[0])
	endPmoaccAplyId = strings.ToLower(args[1])

	startKey := "BizGldPmoaccAplyInf" + startPmoaccAplyId
	endKey := "BizGldPmoaccAplyInf" + endPmoaccAplyId

	resultsIterator, err := stub.GetStateByRange(startKey, endKey)
	if err != nil {
		return shim.Error(err.Error())
	}
	defer resultsIterator.Close()                                                           // 迭代器用完需关闭！！

	if !resultsIterator.HasNext(){
		timestamp := time.Now().Format("2006-01-02 15:04:05")
		fmt.Printf("- can not find gldPmoaccAplyInf by pmoaccAplyId range ! %v \n",timestamp)
		return shim.Success([]byte("can not find gldPmoaccAplyInf by pmoaccAplyId range"))
	}
	var buffer bytes.Buffer
	buffer.WriteString("[")                                                              // 对查询结果进行封装，封装为对象数组json格式
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
	fmt.Printf("- query range by startPmoaccAplyId and endPmoaccAplyId successfully ! %v \n",timestamp)
//	fmt.Printf("- the query result is : \n")
//	fmt.Printf("  %v \n",buffer.String())

	return shim.Success(buffer.Bytes())                                                   // 最终，结果以[]byte形式返回

}

// 通过querystring查询gldPmoaccAplyInf(复杂查询)
func queryGldPmoaccAplyInfByQryStr(stub shim.ChaincodeStubInterface, args []string) pb.Response {

	var err error
	var queryString string

	//	fmt.Println("- start query gldAccvalRltv by queryString, the queryString is : \n %v \n",queryString )
	fmt.Println("- start query gldPmoaccAplyInf by queryString" )

	if len(args)!= 1{
		return shim.Error("Incorrect number of arguments. Expecting querystring of the gldPmoaccAplyInf to query")
	}
	queryString = args[0]

	if queryString == ""{                                                                       // 判定querystring是否为空
		return shim.Error("The queryString can not be null ")
	}
	qryStrBizGldPmoaccAplyInf := QryStrBizGldPmoaccAplyInf{}
	err = json.Unmarshal([]byte(queryString), &qryStrBizGldPmoaccAplyInf)
	if err != nil {
		return shim.Error(err.Error())
	}
	if isEmptyBizGldPmoaccAplyInf(qryStrBizGldPmoaccAplyInf.Selector){                          // 判断queryString中有没有查询条件
		return shim.Error("The selector of queryString can not be null ")
	}

	querystring := tranfBizGldPmoaccAplyInfToQryStr(qryStrBizGldPmoaccAplyInf.Selector)

	resultsIterator, err := stub.GetQueryResult(querystring)                        // 返回迭代器，querystring格式：{"selector":{"docType":"BizGldAccvalRltv","CreateUser":"jzk",...}}
	if err != nil {
		return shim.Error(err.Error())
	}
	defer resultsIterator.Close()

	if !resultsIterator.HasNext(){
		timestamp := time.Now().Format("2006-01-02 15:04:05")
		fmt.Printf("- can not find gldPmoaccAplyInf by querystring ! %v \n",timestamp)
		return shim.Success([]byte("can not find gldPmoaccAplyInf by querystring"))
	}
	var buffer bytes.Buffer
	buffer.WriteString("[")                                                                  // 封装成对象数组json串格式
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

// 通过condition查询gldPmoaccAplyInf(复杂查询)
func queryGldPmoaccAplyInfByCond(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	var err error

	fmt.Println("- start query gldPmoaccAplyInf by conditions" )

	if len(args) < 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}
	gldPmoaccAplyInf := args[0]                                                                     // 查询条件以json串传递，格式为：{"docType":"BizGldAccvalRltv","id":"id0001","createTime":"20171110",...}
	bizGldPmoaccAplyInf := BizGldPmoaccAplyInf{}
	err = json.Unmarshal([]byte(gldPmoaccAplyInf), &bizGldPmoaccAplyInf)
	if err != nil {
		return shim.Error(err.Error())
	}

	if isEmptyBizGldPmoaccAplyInf(bizGldPmoaccAplyInf){
		return shim.Error("The query conditions can not be null")
	}

	queryString := tranfBizGldPmoaccAplyInfToQryStr(bizGldPmoaccAplyInf)
	resultsIterator, err := stub.GetQueryResult(queryString)                                        // 返回迭代器
	if err != nil {
		return shim.Error(err.Error())
	}
	defer resultsIterator.Close()

	if !resultsIterator.HasNext(){
		timestamp := time.Now().Format("2006-01-02 15:04:05")
		fmt.Printf("- can not find gldPmoaccAplyInf by querystring ! %v \n",timestamp)
		return shim.Success([]byte("can not find gldPmoaccAplyInf by conditions"))
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

// 通过pmoaccAplyId查询gldPmoaccAplyInf的历史（复杂查询）
func queryGldPmoaccAplyInfHsyByPmoaccAplyId(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	var err error
	var pmoaccAplyId string

	fmt.Println("- start query gldPmoaccAplyInf history by pmoaccAplyId")

	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting pmoaccAplyId of the gldPmoaccAplyInf to query")
	}

	pmoaccAplyId = strings.ToLower(args[0])

	resultsIterator, err := stub.GetHistoryForKey("BizGldPmoaccAplyInf"+ pmoaccAplyId)                  // 返回迭代器
	if err != nil {
		return shim.Error(err.Error())
	}
	defer resultsIterator.Close()

	if !resultsIterator.HasNext(){
		timestamp := time.Now().Format("2006-01-02 15:04:05")
		fmt.Printf("- can not find gldPmoaccAplyInf history by pmoaccAplyId ! %v \n",timestamp)
		return shim.Success([]byte("can not find gldPmoaccAplyInf history by pmoaccAplyId"))
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
	fmt.Printf("- query history by pmoaccAplyId successfully ! %v \n",timestamp)
//	fmt.Printf("- the query result is : \n")
//	fmt.Printf("  %v \n",buffer.String())

	return shim.Success(buffer.Bytes())

}

