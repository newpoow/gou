package gou

import (
	"fmt"
	"io/fs"
	"strings"

	"github.com/yaoapp/gou/helper"
	"github.com/yaoapp/kun/exception"
	"github.com/yaoapp/kun/utils"
)

// LoadAPI 载入数据接口
func LoadAPI(file fs.File) *API {
	defer file.Close()
	http := HTTP{}
	err := helper.UnmarshalFile(file, &http)
	if err != nil {
		panic(err)
	}

	return &API{
		File: file,
		Type: "http",
		HTTP: http,
	}
}

// Reload 重新载入API
func (api *API) Reload() {}

// Run 执行指令并返回结果 name = "models.user.Find", name = "plugins.user.Login"
func Run(name string, args ...interface{}) interface{} {
	typ, class, method := extraProcess(name)
	switch typ {
	case "models":
		return runModel(class, method, args...)
	case "plugins":
		return runPlugin(class, method, args...)
	}
	return nil
}

// runModel name = user, method = login, args = [1]
func runPlugin(name string, method string, args ...interface{}) interface{} {
	utils.Dump(name, method, args)
	return nil
}

// runModel name = user, method = find, args = [1]
func runModel(name string, method string, args ...interface{}) interface{} {
	mod := Select(name)
	switch method {
	case "find":
		validateArgs(name, method, args, 1)
		return mod.Find(args[0])
	}
	return nil
}

// validateArgs( args )
func validateArgs(name string, method string, args []interface{}, length int) {
	if len(args) < length {
		exception.New(
			fmt.Sprintf("Model:%s%s(args...); 参数错误", name, method),
			400,
		).Throw()
	}
}

// extraProcess 解析执行方法
func extraProcess(name string) (typ string, class string, method string) {
	namer := strings.Split(name, ".")
	last := len(namer) - 1
	if last < 2 {
		exception.New(
			fmt.Sprintf("Process:%s 格式错误", name),
			400,
		).Throw()
	}
	typ = strings.ToLower(namer[0])
	class = strings.ToLower(strings.Join(namer[1:last], "."))
	method = strings.ToLower(namer[last])
	return typ, class, method
}