package router

import (
    "github.com/klenin/orc/mvc/controllers"
    "net/http"
    "reflect"
    "strings"
)

type FastCGIServer struct{}

func (this FastCGIServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
    parts := strings.Split(r.URL.Path, "/")
    controllerName := "indexcontroller"
    methodName := "index"

    if len(parts) >= 2 && parts[1] != "" {
        controllerName = parts[1]
    }
    if len(parts) >= 3 && parts[2] != "" {
        methodName = parts[2]
    }

    if controller := findController(controllerName); controller != nil {
        controller.Elem().FieldByName("Request").Set(reflect.ValueOf(r))
        controller.Elem().FieldByName("Response").Set(reflect.ValueOf(w))
        cType := controller.Type()
        if cMethod := findMethod(cType, methodName); cMethod != nil {
            parts = append(parts, make([]string, cMethod.Type.NumIn())...)
            cMethod.Func.Call(append(
                []reflect.Value{*controller},
                stringToValueSlice(parts[3:3 + cMethod.Type.NumIn() - 1])...,
            ))
        } else {
            http.Error(w, "Unable to locate method in controller.", http.StatusMethodNotAllowed)
        }
    } else {
        http.Error(w, "Unable to locate controller.", http.StatusMethodNotAllowed)
    }
}

func findController(controllerName string) *reflect.Value {
    baseController := new(controllers.BaseController)
    cmt := reflect.TypeOf(baseController)
    if cmt_method := findMethod(cmt, controllerName); cmt_method != nil {
        params := []reflect.Value{reflect.ValueOf(baseController)}
        result := cmt_method.Func.Call(params)
        return &result[0]
    }
    return nil
}

func findMethod(cType reflect.Type, methodName string) *reflect.Method {
    for i := 0; i < cType.NumMethod(); i++ {
        method := cType.Method(i)
        if strings.ToLower(method.Name) == strings.ToLower(methodName) {
            return &method
        }
    }
    return nil
}

func stringToValueSlice(a []string) (r []reflect.Value) {
    r = make([]reflect.Value, len(a))
    for k, v := range a {
        r[k] = reflect.ValueOf(v)
    }
    return r
}
