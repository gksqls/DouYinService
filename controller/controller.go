package controller

var (
	result = make(map[string]interface{})
)

func success() map[string]interface{} {
	result["msg"] = "SUCCESS"
	result["code"] = 0
	return result
}
