package main

import (
	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
	"encoding/json"
	"fmt"
	"reflect"
	"time"
	"strings"
	"bytes"
	"strconv"
)
/**************************************************************************************/
// 金票账款关联结构体
type BizGldAccvalRltv struct{
	// 结构体变量通过json.marshal转换成json串后，若有标签`json:"tagname"` ，则tagname转换为键，若无标签，则成员名转换为键，成员的值转换为值
	// 	json串通过json.unmarshl转换为结构体变量，按json串中每个key依序查询结构体中字段是否匹配，若匹配则将value值赋予目标字段的值
	//    顺序：1）一个包含key标签的字段 2）一个名为key的字段
	// 若用Base，则BizGldAccvalRltv变量的第0个成员就为对象了，嵌套关系而非替代
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

	GldId      string     `json:"gldId"`               // 金票编号
	RcvbId     string     `json:"rcvbId"`              // 应收账款编号
	EstbAmt    float64    `json:"estbAmt"`             // 开立金额

	// 2017.12.14 新增
	EstbBal    float64    `json:"estbBal"`             // 开立余额
}

// 用于解析queryString
type QryStrBizGldAccvalRltv struct{
	Selector          BizGldAccvalRltv         `json:"selector"`
}

// 判断BizGldAccvalRltv变量是否为空
func isEmptyBizGldAccvalRltv(arg BizGldAccvalRltv) bool {

	value := reflect.ValueOf(arg)
	num := 0
	for i:=0;i<value.NumField();i++{
		if isEmpty( value.Field(i).Interface() ){
			num++                                   // 计数结构体中为空成员个数
		}
	}

	if num == value.NumField(){                     // 所有成员为空，结构体才为空
		return true
	}else{
		return false
	}

}

// BizGldAccvalRltv变量变为queryString
func tranfBizGldAccvalRltvToQryStr(arg BizGldAccvalRltv)( string){
	queryString := "{\"selector\":{\"docType\":\"BizGldAccvalRltv\""                  // 封装，头部

	value := reflect.ValueOf(arg)
	typ := reflect.TypeOf(arg)
	for i:=1;i<value.NumField();i++{                                                  // Field(0)跳过
		if !isEmpty( value.Field(i).Interface() ){
			keyname := typ.Field(i).Name                                              // 取成员名称
			keystring := strFirstToLower(keyname)                                     // 首字母小写，生成json串中字段名称
			valuestring := interfaceTostring(value.Field(i).Interface())              // 取成员的值，需转换为字符串
			partstring := fmt.Sprintf(",\"%v\":\"%v\"",keystring,valuestring)
			queryString = queryString + partstring
		}
	}
	queryString = queryString + "}}"                                                   // 封装，尾部
	// {"selector":{"docType":"BizGldAccvalRltv","CreateUser":"jzk"},....}

	return queryString
}


/**************************************************************************************/
// 保存或更新gldAccvalRltv
func saveOrUpdateGldAccvalRltv(stub shim.ChaincodeStubInterface, args []string) pb.Response{
	var err error
	fmt.Println("- start saveOrUpdate gldAccvalRltv")                                  // print都是输出到container的log上

	/*1、首先判断参数个数是否正确*/
	if len(args) < 1 {                                                                    // 数据以json串的形式传递进来，只有一个string变量
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}

	/*2、对接收到的数据进行解析*/
	gldAccvalRltv := args[0]
//	fmt.Println("- the received gldAccvalRltv args is :",gldAccvalRltv)                   // 将接收到的参数打印出来    println既可以打印字符串也可以打印变量，printf只能打印字符串
	bizGldAccvalRltv1 := BizGldAccvalRltv{}                                               // 创建变量，分配内存；json串需与结构体类型保持一致
	err = json.Unmarshal([]byte(gldAccvalRltv), &bizGldAccvalRltv1)                       // json串转换；peer调用方法，由其完成json串与相应数据结构的对应，并非被动盲目接收json串
	if err != nil {
		return shim.Error(err.Error())
	}

	/*3、判断数据是否正确：金票ID是否存在*/
	if bizGldAccvalRltv1.GldId == "" {
		return shim.Error("gldId can't be null")
	} else {
	/*4、判断是save还是update*/
		bizGldAccvalRltvFromState, err := stub.GetState("BizGldAccvalRltv"+               // 以*.ObjectType+GldId作为key，一张金票可能对应多个结构体，防止数据覆盖
			                                            bizGldAccvalRltv1.GldId)          // 不能 bizGldAccvalRltv1.ObjectType + bizGldAccvalRltv1.GldId
		if err != nil {
			return shim.Error("Failed to get bizGldAccvalRltv:" + err.Error())           // 读取过程出错
		} else if bizGldAccvalRltvFromState == nil {                                     // 读取成功，worldstate中没有该条记录key-value
	/*5、数据save*/
			bizGldAccvalRltv1.ObjectType = "BizGldAccvalRltv"
			bizGldAccvalRltvToState, err := json.Marshal(bizGldAccvalRltv1)
			err = stub.PutState("BizGldAccvalRltv" + bizGldAccvalRltv1.GldId,             // 函数收到的参数是json串，json,Marshal 输出类型为[]byte,json.Unmarshal 输入类型为[]byte
													 bizGldAccvalRltvToState)             // stub.PutState输入中，value是[]byte类型；stub.GutState输出是[]byte类型
			if err != nil {
				return shim.Error(err.Error())
			}

			timestamp := time.Now().Format("2006-01-02 15:04:05")                  // 时间戳标准化
			fmt.Printf("- save successfully ! %v\n",timestamp)                    // 输出数据，用于验证逻辑是否正确
			fmt.Printf("- the key of record is : BizGldAccvalRltv%v \n",bizGldAccvalRltv1.GldId)
			fmt.Printf("- the value of record is ： %v \n",string(bizGldAccvalRltvToState))

			return shim.Success(nil)
		} else {
	/*6、数据update*/
			bizGldAccvalRltv2 := BizGldAccvalRltv{}
			err = json.Unmarshal([]byte(bizGldAccvalRltvFromState), &bizGldAccvalRltv2)   // 本身就是[]byte，没必要再进行转换
			value1 := reflect.ValueOf(&bizGldAccvalRltv1).Elem()                          // bizGldAccvalRltv1：保存通过形参传入、解析后的数据；bizGldAccvalRltv2：保存从worldstate中读取的原有数据
			value2 := reflect.ValueOf(&bizGldAccvalRltv2).Elem()                          // & 必须有；.Elem()必须有
			for i:=0; i<value1.NumField(); i++{                                           // 遍历结构体成员
				if !isEmpty( value1.Field(i).Interface() ){                               // 判断成员是否为空，value1.Field(i)得到的是Value类型变量，Interface()必须有
					value2.Field(i).Set( value1.Field(i) )                                // 设置结构体成员值
				}
			}
			bizGldAccvalRltvToState,err := json.Marshal(bizGldAccvalRltv2)                // 局部变量作用范围就在该{}内
			err = stub.PutState("BizGldAccvalRltv" + bizGldAccvalRltv1.GldId, bizGldAccvalRltvToState)
			if err != nil {
				return shim.Error(err.Error())
			}

			timestamp := time.Now().Format("2006-01-02 15:04:05")
			fmt.Printf("- update successfully ! %v \n",timestamp)
			fmt.Printf("- the key of record is : BizGldAccvalRltv%v \n",bizGldAccvalRltv1.GldId)
			fmt.Printf("- the value of record is ： %v \n",string(bizGldAccvalRltvToState))

		}
	}
	return shim.Success(nil)
}


// 通过gldId查询gldAccvalRltv
func queryGldAccvalRltvByGldId(stub shim.ChaincodeStubInterface, args []string) pb.Response {

	var err error
	var gldId, jsonResp string

	fmt.Println("- start query gldAccvalRltv by gldId")

	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting gldId of the gldAccvalRltv to query")
	}

	gldId = strings.ToLower(args[0])                                              // 通过GldId查询gldAccvalRltv

	bizGldAccvalRltvFromState, err := stub.GetState("BizGldAccvalRltv" + gldId)

	if err != nil {                                                               // 查询过程出现错误
		jsonResp = "{\"Error\":\"Failed to get state for " + gldId + "\"}"
		return shim.Error(jsonResp)
	} else if bizGldAccvalRltvFromState == nil {                                 // 没有查到相应记录
		jsonResp = "{\"Error\":\"bizgldAccvalRltv does not exist: " + gldId + "\"}"
		return shim.Error(jsonResp)
	}

	var buffer bytes.Buffer                                                      // 缓存，用来保存查询结果
	buffer.WriteString("{\"Key\":")
	buffer.WriteString("\"BizGldAccvalRltv")
	buffer.WriteString(gldId)
	buffer.WriteString("\"")
	buffer.WriteString(", \"Value\":")
	buffer.WriteString(string(bizGldAccvalRltvFromState))
	buffer.WriteString("}")
	// {"Key":"BizGldAccvalRltvgld0001", "Value":{"docType":"BizGldAccvalRltv","id":"id0001","createTime":"20171110",...}}

	timestamp := time.Now().Format("2006-01-02 15:04:05")
	fmt.Printf("- query by gldId successfully ! %v \n",timestamp)
	fmt.Printf("- the query result is : \n")
	fmt.Printf("  %v \n",buffer.String())

	return shim.Success(buffer.Bytes())
}

// 通过gldId范围查询gldAccvalRltv(复杂查询)
func queryGldAccvalRltvByGldIdRange(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	var err error
	var startGldId,endGldId string

	fmt.Println("- start query gldAccvalRltv by gldId range")

	if len(args)<2{
		return shim.Error("Incorrect number of arguments. Range query expect startGldId and endGldId of the gldAccvalRltv to query")
	}
	startGldId = strings.ToLower(args[0])
	endGldId = strings.ToLower(args[1])

	startKey := "BizGldAccvalRltv" + startGldId                              // 封装为key的形式
	endKey := "BizGldAccvalRltv" + endGldId

	resultsIterator, err := stub.GetStateByRange(startKey, endKey)           // 返回结果为迭代器，StateQueryIterator 结构体类型
	if err != nil {                                                          // 范围查询过程出国
		return shim.Error(err.Error())
	}
	defer resultsIterator.Close()                                           // 迭代器用完需关闭！！

	if !resultsIterator.HasNext(){
		timestamp := time.Now().Format("2006-01-02 15:04:05")
		fmt.Printf("- can not find gldAccvalRltv by gldId range ! %v \n",timestamp)
		return shim.Success([]byte("can not find gldAccvalRltv by gldId range"))
	}
	var buffer bytes.Buffer                                                  // 缓存，用来保存查询结果
	buffer.WriteString("[")                                               // 对查询结果进行封装，封装为对象数组json格式

	bArrayMemberAlreadyWritten := false
	for resultsIterator.HasNext() {                                          // 遍历迭代器中的 key-value
		queryResponse, err := resultsIterator.Next()                         // queryResponse 是 queryresult.KV 结构体类型，内有三个成员{Namespace，Key，Value}

		if err != nil {
			return shim.Error(err.Error())                                   // 迭代过程出错
		}
		if bArrayMemberAlreadyWritten == true {                               // json串中，第一个元素前不需要加','
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
	buffer.WriteString("]")                                                // 封装后的格式为 [{"key":" queryResponse.Key","Record":queryResponse.Value},......]
	// [{"Key":"BizGldAccvalRltvgld0001", "Value":{"docType":"BizGldAccvalRltv","id":"id0001","createTime":"20171110",...}}
	//  ,{"Key":"BizGldAccvalRltvgld0002", "Value":{"docType":"BizGldAccvalRltv","id":"id0001","createTime":"20171110",...}}
	//   ...]

	timestamp := time.Now().Format("2006-01-02 15:04:05")
	fmt.Printf("- query range by satrtGldId and endGldId successfully ! %v \n",timestamp)
	fmt.Printf("- the query result is : \n")
	fmt.Printf("  %v \n",buffer.String())

	return shim.Success(buffer.Bytes())                                      // 最终，结果以[]byte形式返回

}

// 通过querystring查询gldAccvalRltv(复杂查询)
func queryGldAccvalRltvByQryStr(stub shim.ChaincodeStubInterface, args []string) pb.Response {

	var err error
	var queryString string

//	fmt.Println("- start query gldAccvalRltv by queryString, the queryString is : \n %v \n",queryString )
	fmt.Println("- start query gldAccvalRltv by queryString" )                            //  querystring格式：{"selector":{"docType":"BizGldAccvalRltv","CreateUser":"jzk",...}}

	if len(args)!= 1{
		return shim.Error("Incorrect number of arguments. Expecting querystring of the gldAccvalRltv to query")
	}
	queryString = args[0]                                                                    // 这里不能转换为小写

	// 对收到的querystring进行判断
	if queryString == ""{                                                                    // 判定querystring是否为空
		return shim.Error("The queryString can not be null ")
	}
	qryStrBizGldAccvalRltv := QryStrBizGldAccvalRltv{}
	err = json.Unmarshal([]byte(queryString), &qryStrBizGldAccvalRltv)
	if err != nil {
		return shim.Error(err.Error())
	}
	if isEmptyBizGldAccvalRltv(qryStrBizGldAccvalRltv.Selector){                             // 判断queryString中有没有查询条件
		return shim.Error("The selector of queryString can not be null ")
	}

	querystring := tranfBizGldAccvalRltvToQryStr(qryStrBizGldAccvalRltv.Selector)            // 生成querystring不能用json.marshal，会将空字段也包含进去
	// queryString：{"selector":{"docType":"BizGldAccvalRltv","createUser":"jzk",...}}

	resultsIterator, err := stub.GetQueryResult(querystring)                                 // 返回迭代器，querystring格式：{"selector":{"docType":"BizGldAccvalRltv","CreateUser":"jzk",...}}
	if err != nil {
		return shim.Error(err.Error())
	}
	defer resultsIterator.Close()

	if !resultsIterator.HasNext(){
		timestamp := time.Now().Format("2006-01-02 15:04:05")
		fmt.Printf("- can not find gldAccvalRltv by querystring ! %v \n",timestamp)
		return shim.Success([]byte("can not find gldAccvalRltv by querystring"))
	}
	var buffer bytes.Buffer                                                   // 缓存，用来保存查询结果
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
	// [{"Key":"BizGldAccvalRltvgld0001", "Value":{"docType":"BizGldAccvalRltv","id":"id0001","createTime":"20171110",...}}
	//  ,{"Key":"BizGldAccvalRltvgld0002", "Value":{"docType":"BizGldAccvalRltv","id":"id0001","createTime":"20171110",...}}
	//   ...]

	timestamp := time.Now().Format("2006-01-02 15:04:05")
	fmt.Printf("- query by queryString successfully ! %v \n",timestamp)
	fmt.Printf("- the query result is : \n")
	fmt.Printf("  %v \n",buffer.String())

	return shim.Success(buffer.Bytes())

}

// 通过condition查询gldAccvalRltv(复杂查询)
func queryGldAccvalRltvByCond(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	var err error

	fmt.Println("- start query gldAccvalRltv by conditions" )

	if len(args) < 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}
	gldAccvalRltv := args[0]                                                 // 查询条件以json串传递，格式为：{"docType":"BizGldAccvalRltv","id":"id0001","createTime":"20171110",...}
	bizGldAccvalRltv := BizGldAccvalRltv{}
	err = json.Unmarshal([]byte(gldAccvalRltv), &bizGldAccvalRltv)
	if err != nil {
		return shim.Error(err.Error())
	}

	if isEmptyBizGldAccvalRltv(bizGldAccvalRltv){
		return shim.Error("The query conditions can not be null")
	}

	queryString := tranfBizGldAccvalRltvToQryStr(bizGldAccvalRltv)            // 生成queryString不能用json.marshal，会将空字段也包含进去
	// queryString：{"selector":{"docType":"BizGldAccvalRltv","createUser":"jzk",...}}

	resultsIterator, err := stub.GetQueryResult(queryString)                  // 返回迭代器
	if err != nil {
		return shim.Error(err.Error())
	}
	defer resultsIterator.Close()

	if !resultsIterator.HasNext(){
		timestamp := time.Now().Format("2006-01-02 15:04:05")
		fmt.Printf("- can not find gldAccvalRltv by querystring ! %v \n",timestamp)
		return shim.Success([]byte("can not find gldAccvalRltv by conditions"))
	}
	var buffer bytes.Buffer                                                   // 缓存，用来保存查询结果
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
	// [{"Key":"BizGldAccvalRltvgld0001", "Value":{"docType":"BizGldAccvalRltv","id":"id0001","createTime":"20171110",...}}
	//  ,{"Key":"BizGldAccvalRltvgld0002", "Value":{"docType":"BizGldAccvalRltv","id":"id0001","createTime":"20171110",...}}
	//   ...]

	timestamp := time.Now().Format("2006-01-02 15:04:05")
	fmt.Printf("- query by conditions successfully ! %v \n",timestamp)
	fmt.Printf("- the query result is : \n")
	fmt.Printf("  %v \n",buffer.String())

	return shim.Success(buffer.Bytes())

}

// 通过gldId查询gldAccvalRltv的历史（复杂查询）
func queryGldAccvalRltvHsyByGldId(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	var err error
	var gldId string

	fmt.Println("- start query gldAccvalRltv history by gldId")

	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting gldId of the gldAccvalRltv to query")
	}

	gldId = strings.ToLower(args[0])

	resultsIterator, err := stub.GetHistoryForKey("BizGldAccvalRltv"+ gldId)      // 返回迭代器
	if err != nil {
		return shim.Error(err.Error())
	}
	defer resultsIterator.Close()

	if !resultsIterator.HasNext(){
		timestamp := time.Now().Format("2006-01-02 15:04:05")
		fmt.Printf("- can not find gldAccvalRltv history by gldId ! %v \n",timestamp)
		return shim.Success([]byte("can not find gldAccvalRltv history by gldId"))
	}
	var buffer bytes.Buffer
	buffer.WriteString("[")

	bArrayMemberAlreadyWritten := false
	for resultsIterator.HasNext() {
		response, err := resultsIterator.Next()              // response为queryresult.KeyModification 结构体类型，内有四个成员{TxId、Value、Timestamp、IsDelete}
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
		if response.IsDelete {                              // 若已删除则value为null
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
	// [{""TxId":"response.TxId", "Value":string(response.Value), "Timestamp":"time.Unix()String()", "IsDelete":"response.IsDelete"}
	//  ,{""TxId":"response.TxId", "Value":string(response.Value), "Timestamp":"time.Unix()String()", "IsDelete":"response.IsDelete"}
	//  ,...]

	timestamp := time.Now().Format("2006-01-02 15:04:05")
	fmt.Printf("- query history by gldId successfully ! %v \n",timestamp)
	fmt.Printf("- the query result is : \n")
	fmt.Printf("  %v \n",buffer.String())

	return shim.Success(buffer.Bytes())

}



/*  小结：
    1、worldstate中，数据是以key-value的形式存储的，key一般为gldId也可为其他Id，
       value可为任意类型，一般以json的形式存储的
	2、worldstate中，若一条数据已存在，重新写入会导致原记录被覆盖
	3、worldstate中，删除数据并非真的删除，只是将其中的删除标志位置位
	4、chaincode中，每个函数都有两个形参，stub用来调用fabric提供的接口，args用来接收
	   peer传来的数据，需保存的数据以json串的形式传递，json串unmarshal直接解析成结构体
	5、chaindode中，每一个函数/方法都对应一种结构体，函数是由外部peer调用的，
	   由peer保证调用的函数与传入的参数的类型保持一致
*/
