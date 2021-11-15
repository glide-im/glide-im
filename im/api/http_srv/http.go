package http_srv

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"go_im/im/api/auth"
	"go_im/im/api/router"
	"go_im/im/message"
	"go_im/pkg/logger"
	"net/http"
	"reflect"
)

type CommonParam struct {
	Uid    int64
	Device int64
	Data   string
}

type Validatable interface {
	Validate(data string) error
}

var g *gin.Engine
var typeRequestInfo = reflect.TypeOf((*route.Context)(nil))

func Run(addr string, port int) error {

	g = gin.Default()

	initRoute()

	ad := fmt.Sprintf("%s:%d", addr, port)
	return g.Run(ad)
}

func initRoute() {
	authApi := auth.AuthApi{}
	// TODO 2021-11-15 完成其他 api 的 http 服务
	post("/api/auth/register", authApi.Register)
}

func onParamValidateFailed(ctx *gin.Context, err error) {
	logger.E("validate request param failed %v", err)
}

func onParamError(ctx *gin.Context, err error) {
	logger.E("resolve api param error %v", err)
}

func requestParam(ctx *gin.Context) (*route.Context, string) {
	commonP := &CommonParam{}
	e := ctx.ShouldBindJSON(commonP)
	if e != nil {
		onParamError(ctx, e)
		return nil, ""
	}
	info := &route.Context{
		Uid:    commonP.Uid,
		Device: commonP.Device,
		R: func(message *message.Message) {
			ctx.JSON(http.StatusOK, message)
		},
	}
	return info, commonP.Data
}

func deserialize(data string, i interface{}) error {
	return json.Unmarshal([]byte(data), i)
}

func post(path string, fn interface{}) {

	handleFunc, paramType, hasParam, validate := reflectHandleFunc(path, fn)

	g.POST(path, func(context *gin.Context) {
		if hasParam {
			reqInfo, data := requestParam(context)
			if reqInfo == nil {
				return
			}
			param := reflect.New(paramType).Interface()
			if validate {
				v := param.(Validatable)
				err := v.Validate(data)
				if err != nil {
					onParamValidateFailed(context, err)
					return
				}
			} else {
				if hasParam {
					err := deserialize(data, param)
					if err != nil {
						onParamError(context, err)
						return
					}
				} else {
					handleFunc.Call(valOf(reqInfo))
					return
				}
			}
			p := reflect.ValueOf(param).Interface()
			handleFunc.Call(valOf(reqInfo, p))
		}
	})
}

func valOf(i ...interface{}) []reflect.Value {
	var rt []reflect.Value
	for _, i2 := range i {
		rt = append(rt, reflect.ValueOf(i2))
	}
	return rt
}

func reflectHandleFunc(path string, handleFunc interface{}) (reflect.Value, reflect.Type, bool, bool) {
	typeHandleFunc := reflect.TypeOf(handleFunc)

	if typeHandleFunc.Kind() != reflect.Func {
		panic("the route handleFunc must be a function, path: " + path)
	}

	argNum := typeHandleFunc.NumIn()

	if argNum == 0 || argNum > 2 {
		panic("route handleFunc bad arguments, path: " + path)
	}

	shouldValidate := false
	var typeParam reflect.Type
	// reflect request param
	if argNum == 2 {
		if typeHandleFunc.In(1).Kind() != reflect.Ptr {
			panic("route handleFunc param must be pointer, route: " + path)
		}

		typeParam = typeHandleFunc.In(1).Elem()
		if typeParam.Kind() != reflect.Struct {
			panic("the second arg of handleFunc must struct")
		}
		_, shouldValidate = reflect.New(typeParam).Interface().(route.Validatable)
	}

	// reflect first param
	if !typeHandleFunc.In(0).AssignableTo(typeRequestInfo) {
		panic("route handleFunc bad arguments, route: " + path)
	}
	valueHandleFunc := reflect.ValueOf(handleFunc)
	return valueHandleFunc, typeParam, argNum == 2, shouldValidate
}
