package main
import (
	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
	"encoding/json"
	"fmt"
	"reflect"
	"time"
	"strings"
	"strconv"
	"bytes"
)


/**************************************************************************************/
// 金票信息结构体
type BizGldInf struct{
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

	GldId          string      `json:"gldId"`           // 金票编号
	OriGldId       string      `json:"oriGldId"`        // 原金票编号
	CtrId          string      `json:"ctrId"`           // 合同编号
	Estb           string      `json:"estb"`            // 开立方：交易买方
	Pypt           string      `json:"pypt"`            // 支付方：交易买方
	RcPty          string      `json:"rcPty"`           // 持票人/接收方：交易卖方
	Fctr           string      `json:"fctr"`            // 保理商
	GldAmt         float64     `json:"gldAmt"`          // 金票金额：金票最初金额，也就是合同额
	GldBal         float64     `json:"gldBal"`          // 金票余额：目前金票剩余额度
	PyAmt          float64     `json:"pyAmt"`           // 支付金额：买房已支付的金额，为金票金额与金票余额之差
	ToPayAmt       float64     `json:"toPayAmt"`        // 待支付金额：其实就是金票余额         开立时：金票待付金额 = 金票金额  付款时：金票待付金额 = 金票待付金额 - 上次支付金额
	FncAmt         float64     `json:"fncAmt"`          // 融资金额：已融资的额度
	FncBal         float64     `json:"fncBal"`          // 融资余额：已融资额度减去买方的还款额
	PmoaccTamt     float64     `json:"pmoaccTamt"`      // 垫付总金额：保理商需垫付的总金额，相当于融资金额
	PmoaccTotBal   float64     `json:"pmoaccTotBal"`    // 垫付总余额：保理商需垫付的总余额，相当于融资余额
	PymTamt        float64     `json:"pymTamt"`         // 付款总金额：买方需付款的总金额
	EstbDay        int64       `json:"estbDay"`         // 开立日/支付日：金票开立日期
	ExDay          int64       `json:"exDay"`           // 到期日：金票到期还款日期
	SignToacptDay  int64       `json:"signToacptDay"`   // 签收日：卖方签收的日期，需保理商确认
	GldStCd        string      `json:"gldStCd"`         // 金票状态：0:未签收、1：已签收、2：部分垫付、3：已垫付、4：已核销、5：未确认、6：已失效
	PcsStCd        string      `json:"pcsStCd"`         // 流程状态代码
	PrnGldId       string      `json:"prnGldId"`        // 父金票编号集合，以逗号隔开
	FncPct         float64     `json:"fncPct"`          // 融资比例：应收装款转让率
	LockInd        string      `json:"lockInd"`         // 锁标志

	// 2017.11.27 新增
	RspbPsnId          string    `json:"rspbPsnId"`     // 经办人编号
	HdlInstId          string    `json:"hdlInstId"`     // 经办机构编号
	HdlDt              int64     `json:"hdlDt"`         // 经办日期（营业日）
	Opin               string    `json:"opin"`          // 意见

	// 2017.12.14  新增
	ToClsAmt           float64   `json:"toClsAmt"`      // 待清分金额

}


// 用于解析queryString
type QryStrBizGldInf struct{
	Selector          BizGldInf         `json:"selector"`
}

// 判断BizGldAccvalRltv变量是否为空
func isEmptyBizGldInf(arg BizGldInf) bool {
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

// BizGldInf变量变为queryString
func tranfBizGldInfToQryStr(arg BizGldInf)( string){
	queryString := "{\"selector\":{\"docType\":\"BizGldInf\""                          // 封装，头部

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
// 保存或更新gldInf
func saveOrUpdateGldInf(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	var err error
	fmt.Println("- start saveOrUpdate gldInf")

	if len(args) < 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}

	gldInf := args[0]
//	fmt.Println("- the received gldInf args is :",gldInf)
	bizGldInf1 := BizGldInf{}
	err = json.Unmarshal([]byte(gldInf), &bizGldInf1)
	if err != nil {
		return shim.Error(err.Error())
	}

	if bizGldInf1.GldId == "" {
		return shim.Error("gldId can't be null")
	} else {
		bizGldInfFromState, err := stub.GetState("BizGldInf" + bizGldInf1.GldId)
		if err != nil {
			return shim.Error("Failed to get bizGldInf:" + err.Error())
		} else if bizGldInfFromState == nil {
			bizGldInf1.ObjectType = "BizGldInf"
			bizGldInf1.DelInd = "0"                                                      // 删除标志，初始值为0
			bizGldInfToState, err := json.Marshal(bizGldInf1)
			err = stub.PutState("BizGldInf" +bizGldInf1.GldId, bizGldInfToState)
			if err != nil {
				return shim.Error(err.Error())
			}

			timestamp := time.Now().Format("2006-01-02 15:04:05")
			fmt.Printf("- save successfully ! %v \n",timestamp)
	//		fmt.Printf("- the key of record is : BizGldInf%v \n",bizGldInf1.GldId)
	//		fmt.Printf("- the value of record is ： %v \n",string(bizGldInfToState))

			return shim.Success(nil)
		} else {
			bizGldInf2 := BizGldInf{}
			err = json.Unmarshal([]byte(bizGldInfFromState), &bizGldInf2)
			value1 := reflect.ValueOf(&bizGldInf1).Elem()
			value2 := reflect.ValueOf(&bizGldInf2).Elem()
			for i:=0; i<value1.NumField(); i++{
				if !isEmpty( value1.Field(i).Interface() ){
					value2.Field(i).Set( value1.Field(i) )
				}
			}
			bizGldInfToState,err := json.Marshal(bizGldInf2)
			err = stub.PutState("BizGldInf" + bizGldInf1.GldId, bizGldInfToState)
			if err != nil {
				return shim.Error(err.Error())
			}

			timestamp := time.Now().Format("2006-01-02 15:04:05")
			fmt.Printf("- update successfully ! %v \n",timestamp)
//			fmt.Printf("- the key of record is : BizGldInf%v \n",bizGldInf1.GldId)
//			fmt.Printf("- the value of record is ： %v \n",string(bizGldInfToState))

		}
	}
	return shim.Success(nil)
}



// 通过gldId查询gldInf
func queryGldInfByGldId(stub shim.ChaincodeStubInterface, args []string) pb.Response {

	var err error
	var gldId, jsonResp string

	fmt.Println("- start query gldInf by gldId")

	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting gldId of the gldInf to query")
	}

	gldId = strings.ToLower(args[0])

	bizGldInfFromState, err := stub.GetState("BizGldInf" + gldId)

	if err != nil {
		jsonResp = "{\"Error\":\"Failed to get state for " + gldId + "\"}"
		return shim.Error(jsonResp)
	} else if bizGldInfFromState == nil {
		jsonResp = "{\"Error\":\"bizGldInf does not exist: " + gldId + "\"}"
		return shim.Error(jsonResp)
	}


	var buffer bytes.Buffer                                                  // 缓存，用来保存查询结果
	buffer.WriteString("{\"Key\":")
	buffer.WriteString("\"BizGldInf")
	buffer.WriteString(gldId)
	buffer.WriteString("\"")
	buffer.WriteString(", \"Value\":")
	buffer.WriteString(string(bizGldInfFromState))
	buffer.WriteString("}")

	timestamp := time.Now().Format("2006-01-02 15:04:05")
	fmt.Printf("- query by gldId successfully ! %v \n",timestamp)
//	fmt.Printf("- the query result is : \n")
//	fmt.Printf("  %v \n",buffer.String())

	return shim.Success(buffer.Bytes())
}


// 通过gldId范围查询gldInf(复杂查询)
func queryGldInfByGldIdRange(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	var err error
	var startGldId,endGldId string

	fmt.Println("- start query gldInf by gldId range")

	if len(args)<2{
		return shim.Error("Incorrect number of arguments. Range query expect startGldId and endGldId of the gldInf to query")
	}
	startGldId = strings.ToLower(args[0])
	endGldId = strings.ToLower(args[1])

	startKey := "BizGldInf" + startGldId
	endKey := "BizGldInf" + endGldId

	resultsIterator, err := stub.GetStateByRange(startKey, endKey)
	if err != nil {
		return shim.Error(err.Error())
	}
	defer resultsIterator.Close()                                            // 迭代器用完需关闭！！

	if !resultsIterator.HasNext(){
		timestamp := time.Now().Format("2006-01-02 15:04:05")
		fmt.Printf("- can not find gldInf by gldId range ! %v \n",timestamp)
		return shim.Success([]byte("can not find gldInf by gldId range"))
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

// 通过querystring查询gldInf(复杂查询)
func queryGldInfByQryStr(stub shim.ChaincodeStubInterface, args []string) pb.Response {

	var err error
	var queryString string

	//	fmt.Println("- start query gldAccvalRltv by queryString, the queryString is : \n %v \n",queryString )
	fmt.Println("- start query gldInf by queryString" )

	if len(args)!= 1{
		return shim.Error("Incorrect number of arguments. Expecting querystring of the gldInf to query")
	}
	queryString = args[0]


	if queryString == ""{                                                     // 判定querystring是否为空
		return shim.Error("The queryString can not be null ")
	}
	qryStrBizGldInf := QryStrBizGldInf{}
	err = json.Unmarshal([]byte(queryString), &qryStrBizGldInf)
	if err != nil {
		return shim.Error(err.Error())
	}
	if isEmptyBizGldInf(qryStrBizGldInf.Selector){                            // 判断queryString中有没有查询条件
		return shim.Error("The selector of queryString can not be null ")
	}

	querystring := tranfBizGldInfToQryStr(qryStrBizGldInf.Selector)

	resultsIterator, err := stub.GetQueryResult(querystring)      // 返回迭代器，querystring格式：{"selector":{"docType":"BizGldAccvalRltv","CreateUser":"jzk",...}}
	if err != nil {
		return shim.Error(err.Error())
	}
	defer resultsIterator.Close()

	if !resultsIterator.HasNext(){
		timestamp := time.Now().Format("2006-01-02 15:04:05")
		fmt.Printf("- can not find gldInf by querystring ! %v \n",timestamp)
		return shim.Success([]byte("can not find gldInf by querystring"))
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

// 通过condition查询gldInf(复杂查询)
func queryGldInfByCond(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	var err error

	fmt.Println("- start query gldInf by conditions" )

	if len(args) < 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}
	gldInf := args[0]                                                              // 查询条件以json串传递，格式为：{"docType":"BizGldAccvalRltv","id":"id0001","createTime":"20171110",...}
	bizGldInf := BizGldInf{}
	err = json.Unmarshal([]byte(gldInf), &bizGldInf)
	if err != nil {
		return shim.Error(err.Error())
	}

	if isEmptyBizGldInf(bizGldInf){
		return shim.Error("The query conditions can not be null")
	}

	queryString := tranfBizGldInfToQryStr(bizGldInf)
	resultsIterator, err := stub.GetQueryResult(queryString)                       // 返回迭代器
	if err != nil {
		return shim.Error(err.Error())
	}
	defer resultsIterator.Close()

	if !resultsIterator.HasNext(){
		timestamp := time.Now().Format("2006-01-02 15:04:05")
		fmt.Printf("- can not find gldInf by querystring ! %v \n",timestamp)
		return shim.Success([]byte("can not find gldInf by conditions"))
	}
	var buffer bytes.Buffer
	buffer.WriteString("[")                                                     // 封装成对象数组json串格式
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

// 通过gldId查询gldInf的历史（复杂查询）
func queryGldInfHsyByGldId(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	var err error
	var gldId string

	fmt.Println("- start query gldInf history by gldId")

	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting gldId of the gldInf to query")
	}

	gldId = strings.ToLower(args[0])

	resultsIterator, err := stub.GetHistoryForKey("BizGldInf"+ gldId)                 // 返回迭代器
	if err != nil {
		return shim.Error(err.Error())
	}
	defer resultsIterator.Close()

	if !resultsIterator.HasNext(){
		timestamp := time.Now().Format("2006-01-02 15:04:05")
		fmt.Printf("- can not find gldInf history by gldId ! %v \n",timestamp)
		return shim.Success([]byte("can not find gldInf history by gldId"))
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
		if response.IsDelete {                                                       // 若已删除则value为null
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

