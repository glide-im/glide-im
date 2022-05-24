package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/glide-im/glideim/im/api/comm"
	"github.com/glide-im/glideim/im/api/router"
	"github.com/glide-im/glideim/im/auth"
	"github.com/glide-im/glideim/im/message"
	"github.com/glide-im/glideim/pkg/logger"
	"net/http"
	"reflect"
)

type CommonParam struct {
	Uid    int64
	Device int64
	Data   String
}

type String string

func (s *String) UnmarshalJSON(bytes []byte) error {
	*s = String(bytes)
	return nil
}

type CommonResponse struct {
	Code int
	Msg  string
	Data interface{}
}

type Validatable interface {
	Validate() error
}

var g *gin.Engine
var rt gin.IRoutes
var typeRequestInfo = reflect.TypeOf((*route.Context)(nil))
var typeError = reflect.TypeOf((*error)(nil)).Elem()

//run http server
func run(addr string, port int) error {

	g = gin.Default()
	rt = g.Use(crosMiddleware())
	initRoute()

	ad := fmt.Sprintf("%s:%d", addr, port)
	return g.Run(ad)
}

func onParamValidateFailed(ctx *gin.Context, err error) {
	logger.D("validate request param failed %v", err)
	_ = ctx.BindJSON(CommonResponse{
		Code: 300,
		Msg:  "invalid parameter",
		Data: nil,
	})
}

func onParamError(ctx *gin.Context, err error) {
	logger.E("resolve api param error %v", err)
	_ = ctx.BindJSON(CommonResponse{
		Code: 300,
		Msg:  "parameter parse error",
		Data: nil,
	})
}

func onHandlerFuncErr(ctx *gin.Context, err error, handlerParam []reflect.Value) {
	errBiz, ok := err.(*comm.ErrApiBiz)
	if ok {
		ctx.JSON(http.StatusOK, CommonResponse{
			Code: errBiz.Code,
			Msg:  errBiz.Error(),
			Data: nil,
		})
		return
	}

	errUnexpected, ok := err.(*comm.ErrUnexpected)
	if ok {
		logger.D("api error, path:%s\n\t%s", ctx.FullPath(), errUnexpected.Line)
		context, ok := handlerParam[0].Interface().(*route.Context)
		if ok {
			logger.D("uid:%d, device:%d", context.Uid, context.Device)
		}
		if len(handlerParam) == 2 {
			marshal, _ := json.Marshal(handlerParam[1].Interface())
			logger.D("param: %v", string(marshal))
		}
		logger.E("msg=%s, origin=%v", errUnexpected.Msg, errUnexpected.Origin)
		ctx.JSON(http.StatusOK, CommonResponse{
			Code: errUnexpected.Code,
			Msg:  errUnexpected.Error(),
			Data: nil,
		})
		return
	}

	ctx.JSON(http.StatusOK, CommonResponse{
		Code: 500,
		Msg:  err.Error(),
		Data: nil,
	})
}

func getContext(ctx *gin.Context) *route.Context {
	info := &route.Context{
		Uid:    0,
		Device: 0,
		R: func(msg *message.Message) {
			response := CommonResponse{
				Code: 100,
				Msg:  "success",
				Data: msg.GetData(),
			}
			ctx.JSON(http.StatusOK, &response)
		},
	}
	a, exists := ctx.Get(CtxKeyAuthInfo)
	if exists {
		authInfo, ok := a.(*auth.AuthInfo)
		if ok {
			info.Uid = authInfo.Uid
			info.Device = authInfo.Device
		} else {
			logger.E("cast request context auth info (%s) failed, the value is: %v", CtxKeyAuthInfo, a)
		}
	}
	return info
}

func getHandler(path string, fn interface{}) func(ctx *gin.Context) {
	handleFunc, paramType, hasParam, validate := reflectHandleFunc(path, fn)
	return func(context *gin.Context) {
		ctx := getContext(context)
		if ctx == nil {
			onParamValidateFailed(context, errors.New("authentication failed"))
			return
		}
		var handlerParam []reflect.Value
		if hasParam {
			param := reflect.New(paramType).Interface()
			err := context.BindJSON(&param)
			if err != nil {
				onParamError(context, errors.New("invalid parameter"))
				return
			}
			if validate {
				err = param.(Validatable).Validate()
				if err != nil {
					onParamValidateFailed(context, err)
					return
				}
			}
			handlerParam = valOf(ctx, param)
		} else {
			handlerParam = valOf(ctx)
		}
		errV := handleFunc.Call(handlerParam)[0].Interface()

		if errV != nil {
			err := errV.(error)
			onHandlerFuncErr(context, err, handlerParam)
		}
	}
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

	if typeHandleFunc.NumOut() != 1 || !typeHandleFunc.Out(0).Implements(typeError) {
		panic("route handler must return an error param, path: " + path)
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
