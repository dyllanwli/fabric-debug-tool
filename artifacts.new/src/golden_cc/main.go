/*
	更新：2017.12.14  jzk
    功能：金票相关记录的保存、更新、删除、查询
	状态：调试、验证通过
	改动：根据测试人员测试情况进行以下修改
	      1）部分结构体增加了新的成员
		  2）所有结构体中ID、与实践和日期有关的成员，数据类型全部改为int64
*/
package main

import (
	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
//	"strconv"
	"fmt"
	"time"
//	"reflect"
//	"bytes"
	"net/http"
	//"log"
)

// chaincode操作类型
type ChainCodeImpl struct{
}

// 初始化方法
func (c *ChainCodeImpl)Init(stub shim.ChaincodeStubInterface) pb.Response{

	timestamp := time.Now().Format("2006-01-02 15:04:05")
	fmt.Println("init complete !",timestamp)

	fmt.Println("==================================================")
	fmt.Println("-now is in init function")
	fmt.Println("==================================================")
	fmt.Println("")
	fmt.Println("")
	fmt.Println("")
	_, err := http.Get("http://www.baidu.com")
	if err != nil {
		// handle error
		fmt.Println("==================================================")
		fmt.Println("-http can not get www.baidu.com")
		fmt.Println("-the error is :",err)
		fmt.Println("==================================================")
		fmt.Println("")
		fmt.Println("")
		fmt.Println("")
	//	return
	}else{
		fmt.Println("==================================================")
		fmt.Println("-http get www.baidu.com success!!!")
		fmt.Println("==================================================")
		fmt.Println("")
		fmt.Println("")
		fmt.Println("")
/*		defer resp.Body.Close()

		headers := resp.Header

		for k, v := range headers {
			fmt.Printf("k=%v, v=%v\n", k, v)
		}

		fmt.Printf("resp status %s,statusCode %d\n", resp.Status, resp.StatusCode)

		fmt.Printf("resp Proto %s\n", resp.Proto)

		fmt.Printf("resp content length %d\n", resp.ContentLength)

		fmt.Printf("resp transfer encoding %v\n", resp.TransferEncoding)

		fmt.Printf("resp Uncompressed %t\n", resp.Uncompressed)

		fmt.Println(reflect.TypeOf(resp.Body)) // *http.gzipReader

		buf := bytes.NewBuffer(make([]byte, 0, 512))

		length, _ := buf.ReadFrom(resp.Body)

		fmt.Println(len(buf.Bytes()))
		fmt.Println(length)
		fmt.Println(string(buf.Bytes()))
		*/
	}


	return shim.Success(nil)
}

// ************* 测试，http www.baidu.com  ***********************
func callExtSer(stub shim.ChaincodeStubInterface, args []string) pb.Response{
	var err error
	fmt.Println("- start callExtSer")
	if len(args) < 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}

	url := args[0]
	fmt.Println("- the external service is ：",url)
	resp, err := http.Get("http://www.baidu.com")
	if err != nil {
		// handle error
		fmt.Println("==================================================")
		fmt.Println("-now is in callExtSer")
		fmt.Println("-can not call the external service :",url)
		fmt.Println("-the error is :", err)
		fmt.Println("==================================================")
		fmt.Println("")
		fmt.Println("")
		fmt.Println("")
	}else{
		fmt.Println("==================================================")
		fmt.Println("-now is in callExtSer")
		fmt.Println("-call the external service successfully :",url)
		fmt.Println("-the response is :")
		fmt.Println(resp)
		fmt.Println("==================================================")
	}

	return shim.Success(nil)

}

// 外部调用统一入口
func (c *ChainCodeImpl)Invoke(stub shim.ChaincodeStubInterface) pb.Response {
	function, args := stub.GetFunctionAndParameters()
	timestamp := time.Now().Format("2006-01-02 15:04:05")
	fmt.Println("invoke is running " + function,timestamp)

	switch function{
	// 结构体1：BizGldAccvalRltv
	case "saveOrUpdateGldAccvalRltv":
		return saveOrUpdateGldAccvalRltv(stub,args)
	case "queryGldAccvalRltvByGldId":
		return queryGldAccvalRltvByGldId(stub,args)
	case "queryGldAccvalRltvByGldIdRange":
		return queryGldAccvalRltvByGldIdRange(stub,args)
	case "queryGldAccvalRltvByQryStr":
		return queryGldAccvalRltvByQryStr(stub,args)
	case "queryGldAccvalRltvByCond":
		return queryGldAccvalRltvByCond(stub,args)
	case "queryGldAccvalRltvHsyByGldId":
		return queryGldAccvalRltvHsyByGldId(stub,args)

	// 结构体2：BizGldCfm
	case "saveOrUpdateGldCfm":
		return saveOrUpdateGldCfm(stub,args)
	case "deleteGldCfm":
		return deleteGldCfm(stub,args)
	case "queryGldCfmByCfmAplyId":
		return queryGldCfmByCfmAplyId(stub,args)
	case "queryGldCfmByCfmAplyIdRange":
		return queryGldCfmByCfmAplyIdRange(stub,args)
	case "queryGldCfmByQryStr":
		return queryGldCfmByQryStr(stub,args)
	case "queryGldCfmByCond":
		return queryGldCfmByCond(stub,args)
	case "queryGldCfmHsyByGldId":
		return queryGldCfmHsyByGldId(stub,args)

	// 结构体3：BizGldCfmDtl
	case "saveOrUpdateGldCfmDtl":
		return saveOrUpdateGldCfmDtl(stub,args)
	case "deleteGldCfmDtl":
		return deleteGldCfmDtl(stub,args)
	case "queryGldCfmDtlByGldId":
		return queryGldCfmDtlByGldId(stub,args)
	case "queryGldCfmDtlByGldIdRange":
		return queryGldCfmDtlByGldIdRange(stub,args)
	case "queryGldCfmDtlByQryStr":
		return queryGldCfmDtlByQryStr(stub,args)
	case "queryGldCfmDtlByCond":
		return queryGldCfmDtlByCond(stub,args)
	case "queryGldCfmDtlHsyByGldId":
		return queryGldCfmDtlHsyByGldId(stub,args)

	// 结构体4：BizGldCntClsfDtl
	case "saveGldCntClsfDtl":
		return saveGldCntClsfDtl(stub,args)
	case "queryGldCntClsfDtlByGldId":
		return queryGldCntClsfDtlByGldId(stub,args)
	case "queryGldCntClsfDtlByGldIdRange":
		return queryGldCntClsfDtlByGldIdRange(stub,args)
	case "queryGldCntClsfDtlByQryStr":
		return queryGldCntClsfDtlByQryStr(stub,args)
	case "queryGldCntClsfDtlByCond":
		return queryGldCntClsfDtlByCond(stub,args)
	case "queryGldCntClsfDtlHsyByGldId":
		return queryGldCntClsfDtlHsyByGldId(stub,args)

	// 结构体5：BizGldDocRltv
	case "saveGldDocRltv":
		return saveGldDocRltv(stub,args)
	case "queryGldDocRltvByGldId":
		return queryGldDocRltvByGldId(stub,args)
	case "queryGldDocRltvByGldIdRange":
		return queryGldDocRltvByGldIdRange(stub,args)
	case "queryGldDocRltvByQryStr":
		return queryGldDocRltvByQryStr(stub,args)
	case "queryGldDocRltvByCond":
		return queryGldDocRltvByCond(stub,args)
	case "queryGldDocRltvHsyByGldId":
		return queryGldDocRltvHsyByGldId(stub,args)

	// 结构体6：BizGldInf
	case "saveOrUpdateGldInf":
		return saveOrUpdateGldInf(stub,args)
	case "queryGldInfByGldId":
		return queryGldInfByGldId(stub,args)
	case "queryGldInfByGldIdRange":
		return queryGldInfByGldIdRange(stub,args)
	case "queryGldInfByQryStr":
		return queryGldInfByQryStr(stub,args)
	case "queryGldInfByCond":
		return queryGldInfByCond(stub,args)
	case "queryGldInfHsyByGldId":
		return queryGldInfHsyByGldId(stub,args)

	// 结构体7：BizGldPmoaccAply
	case "saveOrUpdateGldPmoaccAply":
		return saveOrUpdateGldPmoaccAply(stub,args)
	case "deleteGldPmoaccAply":
		return deleteGldPmoaccAply(stub,args)
	case "queryGldPmoaccAplyByGldId":
		return queryGldPmoaccAplyByGldId(stub,args)
	case "queryGldPmoaccAplyByGldIdRange":
		return queryGldPmoaccAplyByGldIdRange(stub,args)
	case "queryGldPmoaccAplyByQryStr":
		return queryGldPmoaccAplyByQryStr(stub,args)
	case "queryGldPmoaccAplyByCond":
		return queryGldPmoaccAplyByCond(stub,args)
	case "queryGldPmoaccAplyHsyByGldId":
		return queryGldPmoaccAplyHsyByGldId(stub,args)

	// 结构体8：BizGldPmoaccAplyDtl
	case "saveOrUpdateGldPmoaccAplyDtl":
		return saveOrUpdateGldPmoaccAplyDtl(stub,args)
	case "deleteGldPmoaccAplyDtl":
		return deleteGldPmoaccAplyDtl(stub,args)
	case "queryGldPmoaccAplyDtlByGldId":
		return queryGldPmoaccAplyDtlByGldId(stub,args)
	case "queryGldPmoaccAplyDtlByGldIdRange":
		return queryGldPmoaccAplyDtlByGldIdRange(stub,args)
	case "queryGldPmoaccAplyDtlByQryStr":
		return queryGldPmoaccAplyDtlByQryStr(stub,args)
	case "queryGldPmoaccAplyDtlByCond":
		return queryGldPmoaccAplyDtlByCond(stub,args)
	case "queryGldPmoaccAplyDtlHsyByGldId":
		return queryGldPmoaccAplyDtlHsyByGldId(stub,args)

	// 结构体9：BizGldPmoaccAplyInf
	case "saveOrUpdateGldPmoaccAplyInf":
		return saveOrUpdateGldPmoaccAplyInf(stub,args)
	case "deleteGldPmoaccAplyInf":
		return deleteGldPmoaccAplyInf(stub,args)
	case "queryGldPmoaccAplyInfByPmoaccAplyId":
		return queryGldPmoaccAplyInfByPmoaccAplyId(stub,args)
	case "queryGldPmoaccAplyInfByPmoaccAplyIdRange":
		return queryGldPmoaccAplyInfByPmoaccAplyIdRange(stub,args)
	case "queryGldPmoaccAplyInfByQryStr":
		return queryGldPmoaccAplyInfByQryStr(stub,args)
	case "queryGldPmoaccAplyInfByCond":
		return queryGldPmoaccAplyInfByCond(stub,args)
	case "queryGldPmoaccAplyInfHsyByPmoaccAplyId":
		return queryGldPmoaccAplyInfHsyByPmoaccAplyId(stub,args)

	// 结构体10：BizGldPymInf
	case "saveOrUpdateGldPymInf":
		return saveOrUpdateGldPymInf(stub,args)
	case "deleteGldPymInf":
		return deleteGldPymInf(stub,args)
	case "queryGldPymInfByGldId":
		return queryGldPymInfByGldId(stub,args)
	case "querGldPymInfByGldIdRange":
		return querGldPymInfByGldIdRange(stub,args)
	case "queryGldPymInfByQryStr":
		return queryGldPymInfByQryStr(stub,args)
	case "queryGldPymInfByCond":
		return queryGldPymInfByCond(stub,args)
	case "queryGldPymInfHsyByGldId":
		return queryGldPymInfHsyByGldId(stub,args)

	// 结构体11：BizGldPyRcrd
	case "saveOrUpdateGldPyRcrd":
		return saveOrUpdateGldPyRcrd(stub,args)
	case "deleteGldPyRcrd":
		return deleteGldPyRcrd(stub,args)
	case "queryGldPyRcrdByOriNewGldId":
		return queryGldPyRcrdByOriNewGldId(stub,args)
	case "queryGldPyRcrdByGldIdRange":
		return queryGldPyRcrdByGldIdRange(stub,args)
	case "queryGldPyRcrdByQryStr":
		return queryGldPyRcrdByQryStr(stub,args)
	case "queryGldPyRcrdByCond":
		return queryGldPyRcrdByCond(stub,args)
	case "querGldPyRcrdHsyByGldId":
		return querGldPyRcrdHsyByGldId(stub,args)

	// 结构体12：BizGldSigntoacptRcrd
	case "saveOrUpdateGldSigntoacptRcrd":
		return saveOrUpdateGldSigntoacptRcrd(stub,args)
	case "queryGldSigntoacptRcrdByGldId":
		return queryGldSigntoacptRcrdByGldId(stub,args)
	case "queryGldSigntoacptRcrdByGldIdRange":
		return queryGldSigntoacptRcrdByGldIdRange(stub,args)
	case "queryGldSigntoacptRcrdByQryStr":
		return queryGldSigntoacptRcrdByQryStr(stub,args)
	case "queryGldSigntoacptRcrdByCond":
		return queryGldSigntoacptRcrdByCond(stub,args)
	case "queryGldSigntoacptRcrdHsyByGldId":
		return queryGldSigntoacptRcrdHsyByGldId(stub,args)

	// 测试 http
	case "callExtSer":
		return callExtSer(stub,args)

	// 所调用的方法未定义
	default :
		fmt.Println("invoke did not find func: " + function)
		return shim.Error("Received unknown function invocation")
	}
}

// chaincode入口
func main(){
	timestamp := time.Now().Format("2006-01-02 15:04:05")
	err := shim.Start(new(ChainCodeImpl))
	if err != nil {
		fmt.Printf("Error starting chaincode: %s  %s", err,timestamp)
	}

	fmt.Println("==================================================")
	fmt.Println("-now is in main function")
	fmt.Println("==================================================")
	fmt.Println("")
	fmt.Println("")
	fmt.Println("")
	_, err = http.Get("http://www.baidu.com")
	if err != nil {
		// handle error
		fmt.Println("==================================================")
		fmt.Println("-http can not get www.baidu.com")
		fmt.Println("-the error is :",err)
		fmt.Println("==================================================")
	//	return
	}else{
		fmt.Println("==================================================")
		fmt.Println("-http get www.baidu.com success!!!")
		fmt.Println("==================================================")

/*		defer resp.Body.Close()

		headers := resp.Header

		for k, v := range headers {
			fmt.Printf("k=%v, v=%v\n", k, v)
		}

		fmt.Printf("resp status %s,statusCode %d\n", resp.Status, resp.StatusCode)

		fmt.Printf("resp Proto %s\n", resp.Proto)

		fmt.Printf("resp content length %d\n", resp.ContentLength)

		fmt.Printf("resp transfer encoding %v\n", resp.TransferEncoding)

		fmt.Printf("resp Uncompressed %t\n", resp.Uncompressed)

		fmt.Println(reflect.TypeOf(resp.Body)) // *http.gzipReader

		buf := bytes.NewBuffer(make([]byte, 0, 512))

		length, _ := buf.ReadFrom(resp.Body)

		fmt.Println(len(buf.Bytes()))
		fmt.Println(length)
		fmt.Println(string(buf.Bytes()))
		*/
	}
}

