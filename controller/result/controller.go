package result

var (
	result = make(map[string]interface{})
)

// 返回默认成功JSON
func Success() map[string]interface{} {
	result["msg"] = "SUCCESS"
	result["code"] = 0
	return result
}
