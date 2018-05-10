package main
import (
	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
	"encoding/json"
	"fmt"
	"strings"
	"reflect"
	"time"
	"bytes"
	"strconv"
)


/**************************************************************************************/
// 金票付款信息结构体
type BizGldPymInf struct{
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

	GldId             string   `json:"gldId"`    // 金票编号
	GldAmt            float64  `json:"gldAmt"`   // 金票金额
	GldBal            float64  `json:"gldBal"`   // 金票余额
	Pyr               string   `json:"pyr"`      // 付款方
	PymAmt            float64  `json:"pymAmt"`   // 付款金额
	PyTm              int64    `json:"pyTm"`     // 付款时间
	PymTp             string   `json:"pymTp"`    // 付款类型
	PymStCd           string   `json:"pymStCd"`  // 付款状态：未生效/有效/失效...
	PcsStCd           string   `json:"pcsStCd"`  // 流程状态
	AplyId            string   `json:"aplyId"`   // 申请编号

	// 2017.11.27 新增
	RspbPsnId          string    `json:"rspbPsnId"`     // 经办人编号
	HdlInstId          string    `json:"hdlInstId"`     // 经办机构编号
	HdlDt              int64     `json:"hdlDt"`         // 经办日期（营业日）

	// 2017.12.14  新增
	FreeZeNo           string     `json:"freeZeNo"`     // 冻结编号

}

// 用于解析queryString
type QryStrBizGldPymInf struct{
	Selector          BizGldPymInf         `json:"selector"`
}

// 判断BizGldPymInf变量是否为空
func isEmptyBizGldPymInf(arg BizGldPymInf) bool {
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

// BizGldPymInf变量变为queryString
func tranfBizGldPymInfToQryStr(arg BizGldPymInf)( string){
	queryString := "{\"selector\":{\"docType\":\"BizGldPymInf\""                  // 封装，头部

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
// 保存或更新gldPymInf
func saveOrUpdateGldPymInf(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	var err error
	fmt.Println("- start saveOrUpdate gldPymInf")

	if len(args) < 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}

	gldPymInf := args[0]
//	fmt.Println("- the received gldPymInf args is :",gldPymInf)
	bizGldPymInf1 := BizGldPymInf{}
	err = json.Unmarshal([]byte(gldPymInf), &bizGldPymInf1)
	if err != nil {
		return shim.Error(err.Error())
	}

	if bizGldPymInf1.GldId == "" {
		return shim.Error("gldId can't be null")
	} else {
		bizGldPymInfFromState, err := stub.GetState("BizGldPymInf" + bizGldPymInf1.GldId)
		if err != nil {
			return shim.Error("Failed to get bizGldPymInf:" + err.Error())
		} else if bizGldPymInfFromState == nil {
			bizGldPymInf1.ObjectType = "BizGldPymInf"
			bizGldPymInfToState, err := json.Marshal(bizGldPymInf1)
			err = stub.PutState("BizGldPymInf" + bizGldPymInf1.GldId, bizGldPymInfToState)
			if err != nil {
				return shim.Error(err.Error())
			}

			timestamp := time.Now().Format("2006-01-02 15:04:05")
			fmt.Printf("- save successfully ! %v \n",timestamp)
	//		fmt.Printf("- the key of record is : BizGldPymInf%v \n",bizGldPymInf1.GldId)
	//		fmt.Printf("- the value of record is ： %v \n",string(bizGldPymInfToState))

			return shim.Success(nil)
		} else {
			bizGldPymInf2 := BizGldPymInf{}
			err = json.Unmarshal([]byte(bizGldPymInfFromState), &bizGldPymInf2)
			value1 := reflect.ValueOf(&bizGldPymInf1).Elem()
			value2 := reflect.ValueOf(&bizGldPymInf2).Elem()
			for i:=0; i<value1.NumField(); i++{
				if !isEmpty( value1.Field(i).Interface() ){
					value2.Field(i).Set( value1.Field(i) )
				}
			}
			bizGldPymInfToState,err := json.Marshal(bizGldPymInf2)
			err = stub.PutState("BizGldPymInf" + bizGldPymInf1.GldId, bizGldPymInfToState)
			if err != nil {
				return shim.Error(err.Error())
			}

			timestamp := time.Now().Format("2006-01-02 15:04:05")
			fmt.Printf("- update successfully ! %v \n",timestamp)
	//		fmt.Printf("- the key of record is : BizGldPymInf%v \n",bizGldPymInf1.GldId)
	//		fmt.Printf("- the value of record is ： %v \n",string(bizGldPymInfToState))

		}
	}
	return shim.Success(nil)
}


// 删除gldPymInf
func deleteGldPymInf(stub shim.ChaincodeStubInterface, args []string) pb.Response {

	var err error
	var gldId, jsonResp string

	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting gldId of the gldPymInf to delete")
	}

	gldId = strings.ToLower(args[0])

	bizGldPymInfFromState, err := stub.GetState("BizGldPymInf" + gldId) //get the painting from chaincode state
	if err != nil {
		jsonResp = "{\"Error\":\"Failed to get state for " + gldId + "\"}"
		return shim.Error(jsonResp)
	} else if bizGldPymInfFromState == nil {
		jsonResp = "{\"Error\":\"bizGldPymInf does not exist: " + gldId + "\"}"
		return shim.Error(jsonResp)
	}
	bizGldPymInf := BizGldPymInf{}

	err = json.Unmarshal(bizGldPymInfFromState,&bizGldPymInf)
	bizGldPymInf.DelInd = "1"
	bizGldPymInfToState, err := json.Marshal(bizGldPymInf)
	err = stub.PutState("BizGldPymInf" + bizGldPymInf.GldId, bizGldPymInfToState)
	if err != nil {
		return shim.Error(err.Error())
	}

	timestamp := time.Now().Format("2006-01-02 15:04:05")
	fmt.Printf("- delete successfully ! %v \n",timestamp)
//	fmt.Printf("- the key of record is : BizGldPymInf%v \n",bizGldPymInf.GldId)
//	fmt.Printf("- the value of record is ： %v \n",string(bizGldPymInfToState))

	return shim.Success(nil)
}


// 通过gldId查询gldPymInf
func queryGldPymInfByGldId(stub shim.ChaincodeStubInterface, args []string) pb.Response {

	var err error
	var gldId, jsonResp string

	fmt.Println("- start query gldPymInf by gldId")


	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting gldId of the gldPymInf to query")
	}

	gldId = strings.ToLower(args[0])

	bizGldPymInfFromState, err := stub.GetState("BizGldPymInf" + gldId)

	if err != nil {
		jsonResp = "{\"Error\":\"Failed to get state for " + gldId + "\"}"
		return shim.Error(jsonResp)
	} else if bizGldPymInfFromState == nil {
		jsonResp = "{\"Error\":\"bizGldPymInf does not exist: " + gldId + "\"}"
		return shim.Error(jsonResp)
	}


	var buffer bytes.Buffer                                                  // 缓存，用来保存查询结果
	buffer.WriteString("{\"Key\":")
	buffer.WriteString("\"BizGldPymInf")
	buffer.WriteString(gldId)
	buffer.WriteString("\"")
	buffer.WriteString(", \"Value\":")
	buffer.WriteString(string(bizGldPymInfFromState))
	buffer.WriteString("}")

	timestamp := time.Now().Format("2006-01-02 15:04:05")
	fmt.Printf("- query by gldId successfully ! %v \n",timestamp)
//	fmt.Printf("- the query result is : \n")
//	fmt.Printf("  %v \n",buffer.String())

	return shim.Success(buffer.Bytes())
}


// 通过gldId范围查询gldPymInf(复杂查询)
func querGldPymInfByGldIdRange(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	var err error
	var startGldId,endGldId string

	fmt.Println("- start query gldPymInf by gldId range")

	if len(args)<2{
		return shim.Error("Incorrect number of arguments. Range query expect startGldId and endGldId of the gldPymInf to query")
	}
	startGldId = strings.ToLower(args[0])
	endGldId = strings.ToLower(args[1])

	startKey := "BizGldPymInf" + startGldId
	endKey := "BizGldPymInf" + endGldId

	resultsIterator, err := stub.GetStateByRange(startKey, endKey)
	if err != nil {
		return shim.Error(err.Error())
	}
	defer resultsIterator.Close()                                            // 迭代器用完需关闭！！

	if !resultsIterator.HasNext(){
		timestamp := time.Now().Format("2006-01-02 15:04:05")
		fmt.Printf("- can not find gldPymInf by gldId range ! %v \n",timestamp)
		return shim.Success([]byte("can not find gldPymInf by gldId range"))
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

// 通过querystring查询gldPymInf(复杂查询)
func queryGldPymInfByQryStr(stub shim.ChaincodeStubInterface, args []string) pb.Response {

	var err error
	var queryString string

	//	fmt.Println("- start query gldAccvalRltv by queryString, the queryString is : \n %v \n",queryString )
	fmt.Println("- start query gldPymInf by queryString" )

	if len(args)!= 1{
		return shim.Error("Incorrect number of arguments. Expecting querystring of the gldPymInf to query")
	}
	queryString = args[0]

	if queryString == ""{                                                     // 判定querystring是否为空
		return shim.Error("The queryString can not be null ")
	}
	qryStrBizGldPymInf := QryStrBizGldPymInf{}
	err = json.Unmarshal([]byte(queryString), &qryStrBizGldPymInf)
	if err != nil {
		return shim.Error(err.Error())
	}
	if isEmptyBizGldPymInf(qryStrBizGldPymInf.Selector){                     // 判断queryString中有没有查询条件
		return shim.Error("The selector of queryString can not be null ")
	}

	querystring := tranfBizGldPymInfToQryStr(qryStrBizGldPymInf.Selector)

	resultsIterator, err := stub.GetQueryResult(querystring)      // 返回迭代器，querystring格式：{"selector":{"docType":"BizGldAccvalRltv","CreateUser":"jzk",...}}
	if err != nil {
		return shim.Error(err.Error())
	}
	defer resultsIterator.Close()

	if !resultsIterator.HasNext(){
		timestamp := time.Now().Format("2006-01-02 15:04:05")
		fmt.Printf("- can not find gldPymInf by querystring ! %v \n",timestamp)
		return shim.Success([]byte("can not find gldPymInf by querystring"))
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

// 通过condition查询gldPymInf(复杂查询)
func queryGldPymInfByCond(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	var err error

	fmt.Println("- start query gldPymInf by conditions" )

	if len(args) < 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}
	gldPymInf := args[0]                                                         // 查询条件以json串传递，格式为：{"docType":"BizGldAccvalRltv","id":"id0001","createTime":"20171110",...}
	bizGldPymInf := BizGldPymInf{}
	err = json.Unmarshal([]byte(gldPymInf), &bizGldPymInf)
	if err != nil {
		return shim.Error(err.Error())
	}

	if isEmptyBizGldPymInf(bizGldPymInf){
		return shim.Error("The query conditions can not be null")
	}

	queryString := tranfBizGldPymInfToQryStr(bizGldPymInf)
	resultsIterator, err := stub.GetQueryResult(queryString)                      // 返回迭代器
	if err != nil {
		return shim.Error(err.Error())
	}
	defer resultsIterator.Close()

	if !resultsIterator.HasNext(){
		timestamp := time.Now().Format("2006-01-02 15:04:05")
		fmt.Printf("- can not find gldPymInf by querystring ! %v \n",timestamp)
		return shim.Success([]byte("can not find gldPymInf by conditions"))
	}
	var buffer bytes.Buffer
	buffer.WriteString("[")                                                   // 封装成对象数组json串格式
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

// 通过gldId查询gldPymInf的历史（复杂查询）
func queryGldPymInfHsyByGldId(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	var err error
	var gldId string

	fmt.Println("- start query gldPymInf history by gldId")

	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting gldId of the gldPymInf to query")
	}

	gldId = strings.ToLower(args[0])

	resultsIterator, err := stub.GetHistoryForKey("BizGldPymInf"+ gldId)      // 返回迭代器
	if err != nil {
		return shim.Error(err.Error())
	}
	defer resultsIterator.Close()

	if !resultsIterator.HasNext(){
		timestamp := time.Now().Format("2006-01-02 15:04:05")
		fmt.Printf("- can not find gldPymInf history by gldId ! %v \n",timestamp)
		return shim.Success([]byte("can not find gldPymInf history by gldId"))
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


