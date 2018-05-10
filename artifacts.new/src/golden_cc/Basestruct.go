package main

import (
	"strconv"
)

/**************************************************************************************/
/*// 基础结构体
type Base struct {
	ObjectType string     `json:"docType"`             //docType is used to distinguish the various types of objects in state database
	Id         int64      `json:"id"`                  // 主键
	CreateTime string     `json:"createTime"`          // 创建时间
	UpdateTime string     `json:"updateTime"`          // 更新时间
	CreateUser string     `json:"createUser"`          // 创建人
	UpdateUser string     `json:"updateUser"`          // 更新人
	ExpdId     string     `json:"expdId"`              // 扩展ID
	DelInd     string     `json:"delInd"`              // 删除标志
	Version    int32      `json:"version"`             // 版本号
	TenantId   string     `json:"tenantId"`            // 租户ID
}
*/
/**************************************************************************************/
/*
// 空接口转int32
func getint32Fromintf(arg interface{}) int32{
	v,ok := arg.(int32)                           // 类型断言
	if ok{
		return v
	}else{
		return 0
	}
}

// 空接口转string
func getstringFromintf(arg interface{}) string{
	v,ok := arg.(string)
	if ok{
		return v
	}else{
		return ""
	}
}

// 空接口转float64
func getfloat64Fromintf(arg interface{}) float64{
	v,ok := arg.(float64)
	if ok{
		return v
	}else{
		return 0.0
	}
}

// 空接口转map
func getmapsliceFromintf(arg interface{}) []map[string]interface{}{
	v,ok := arg.([]map[string]interface{})
	if ok{
		return v
	}else{
		return nil
	}
}
*/

// 判断变量是否为初始零值
func isEmpty(arg interface{})(bool){
	switch v := arg.(type){
	case int32 :
		if v == 0{
			return true
		}else{
			return false
		}
	case int64 :
		if v ==0{
			return true
		}else{
			return false
		}
	case float64 :
		if v==0.0 {
			return true
		}else{
			return false
		}
	case string :
		if v==""{
			return true
		}else{
			return false
		}
	default :
		return true                 // 默认为空，不改变原值
	}
}


// 接口转换为string
func interfaceTostring(arg interface{})string {
	switch v := arg.(type){
	case int32 :
		val:= strconv.Itoa(int(v))
		return val
	case float64 :
		val := strconv.AppendFloat( []byte(""),v,'f',5,32 )
		return string(val)
	case string :
		return v
	default :
		return ""
	}
}

// 字符串首字母转小写
func strFirstToLower(str string) string {
	ans := ""
	for i,ch := range str{
		if i ==0 {
			if ch >= 65 && ch <=90{
				ans = ans + string(ch+32)        // 首字母，大写转小写
			}else{
				ans = ans + string(ch)
			}
		}else{
			ans = ans + string(ch)
		}
	}

	return ans
}

