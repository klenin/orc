package router

import (
    "github.com/klenin/orc/mvc/controllers"
    "net/http"
    "reflect"
    "strings"
)

type FastCGIServer struct{}

func (this FastCGIServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
    url := r.URL
    parts := strings.Split(url.Path, "/")
    controllerName := "indexcontroller"
    methodName := "index"

    if len(parts) < 2 {
        // index
    } else if len(parts) < 3 {
        if parts[1] != "" {
            controllerName = parts[1]
        }
    } else {
        controllerName = parts[1]
        if parts[2] != "" {
            methodName = parts[2]
        }
    }
    controller := FindController(controllerName)
    if controller != nil {
        controller.Elem().FieldByName("Request").Set(reflect.ValueOf(r))
        controller.Elem().FieldByName("Response").Set(reflect.ValueOf(w))
        cType := controller.Type()
        cMethod := FindMethod(cType, methodName)
        if cMethod != nil {
            params := PopulateParams(*cMethod, parts)
            allParams := make([]reflect.Value, 0)
            cMethod.Func.Call(append(append(allParams, *controller), params...))
        } else {
            http.Error(w, "Unable to locate index method in controller.", http.StatusMethodNotAllowed)
        }
    } else {
        http.Error(w, "Unable to locate default controller.", http.StatusMethodNotAllowed)
    }
}

func FindController(controllerName string) *reflect.Value {
    baseController := new(controllers.BaseController)
    cmt := reflect.TypeOf(baseController)
    count := cmt.NumMethod()
    for i := 0; i < count; i++ {
        cmt_method := cmt.Method(i)
        if strings.ToLower(cmt_method.Name) == strings.ToLower(controllerName) {
            params := make([]reflect.Value, 1)
            params[0] = reflect.ValueOf(baseController)
            result := cmt_method.Func.Call(params)
            return &result[0]
        }
    }
    return nil
}

func FindMethod(cType reflect.Type, methodName string) *reflect.Method {
    count := cType.NumMethod()
    for i := 0; i < count; i++ {
        method := cType.Method(i)
        if strings.ToLower(method.Name) == strings.ToLower(methodName) {
            return &method
        }
    }
    return nil
}

func PopulateParams(method reflect.Method, parts []string) []reflect.Value {
    numParams := method.Type.NumIn() - 1
    params := make([]reflect.Value, numParams)
    for x := 0; x < numParams; x++ {
        if len(parts) > (x + 3) {
            params[x] = reflect.ValueOf(parts[x+3])
        } else {
            params[x] = reflect.ValueOf("")
        }
    }
    return params
}
