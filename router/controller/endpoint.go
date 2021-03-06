package controller

import (
    "github.com/gin-gonic/gin"
    "net/http"
    "reflect"
    "strings"
)

type Endpoint interface {
    Get(c *gin.Context)
    Post(c *gin.Context)
    Put(c *gin.Context)
    Delete(c *gin.Context)
    Head(c *gin.Context)
    Options(c *gin.Context)
    Patch(c *gin.Context)
    Trace(c *gin.Context)

    GetPermission() string
    Register(r *gin.RouterGroup, decors ...func(handlerFunc gin.HandlerFunc) gin.HandlerFunc)
}

type EndpointBase struct {
    Path         string
}

func (ex *EndpointBase) Head(c *gin.Context) {
    c.JSON(http.StatusMethodNotAllowed, gin.H{
        "code": http.StatusMethodNotAllowed,
        "msg": "HEAD is not impl",
    })
}

func (ex *EndpointBase) Options(c *gin.Context) {
    c.JSON(http.StatusMethodNotAllowed, gin.H{
        "code": http.StatusMethodNotAllowed,
        "msg": "OPTIONS is not impl",
    })
}

func (ex *EndpointBase) Patch(c *gin.Context) {
    c.JSON(http.StatusMethodNotAllowed, gin.H{
        "code": http.StatusMethodNotAllowed,
        "msg": "PATCH is not impl",
    })
}

func (ex *EndpointBase) Trace(c *gin.Context) {
    c.JSON(http.StatusMethodNotAllowed, gin.H{
        "code": http.StatusMethodNotAllowed,
        "msg": "TRACE is not impl",
    })
}

func (ex *EndpointBase) Get(c *gin.Context) {
    c.JSON(http.StatusMethodNotAllowed, gin.H{
        "code": http.StatusMethodNotAllowed,
        "msg": "GET is not impl",
    })
}

func (ex *EndpointBase) Post(c *gin.Context) {
    c.JSON(http.StatusMethodNotAllowed, gin.H{
        "code": http.StatusMethodNotAllowed,
        "msg": "POST is not impl",
    })
}

func (ex *EndpointBase) Put(c *gin.Context) {
    c.JSON(http.StatusMethodNotAllowed, gin.H{
        "code": http.StatusMethodNotAllowed,
        "msg": "PUT is not impl",
    })
}

func (ex *EndpointBase) Delete(c *gin.Context) {
    c.JSON(http.StatusMethodNotAllowed, gin.H{
        "code": http.StatusMethodNotAllowed,
        "msg": "DELETE is not impl",
    })
}

func (ex *EndpointBase) GetPermission() string {
    return ""
}

func Dispatch(e interface{}) gin.HandlerFunc {
    funcMap := make(map[string]func(c *gin.Context))
    //ref := reflect.Indirect(reflect.ValueOf(e))
    ref := reflect.ValueOf(e)
    fillMap := func(method string) {
        b := make([]byte, len(method))
        copy(b, method)
        name := string(b)
        name = name[:1] + strings.ToLower(name[1:])

        m := ref.MethodByName(name)
        if m.IsValid() {
            i := m.Interface()
            if f, ok := i.(func(*gin.Context)); ok {
                funcMap[method] = f
            }
        }
    }

    am := []string{
        "GET",
        "POST",
        "PUT",
        "DELETE",
        "PATCH",
        "HEAD",
        "OPTIONS",
        "TRACE",
    }

    for i := range am {
        fillMap(am[i])
    }

    return func(c *gin.Context) {
        res, ok := funcMap[c.Request.Method]
        if !ok {
            c.JSON(http.StatusMethodNotAllowed, gin.H{
                "code": http.StatusMethodNotAllowed,
                "msg": "method not allow here",
            })
            return
        }
        res(c)
    }
}

func (ex *EndpointBase) Register(e interface{}, r *gin.RouterGroup, decors ...func(handlerFunc gin.HandlerFunc) gin.HandlerFunc) {
    r.Any(ex.Path, HandlerPipeline(Dispatch(e), decors...))
    res, ok := e.(Endpoint)
    if ok {
        res.GetPermission()
    }
}

func HandlerPipeline(target gin.HandlerFunc, decors ...func(handlerFunc gin.HandlerFunc) gin.HandlerFunc) gin.HandlerFunc {
    for i := range decors {
        d := decors[len(decors) - i - 1]
        target = d(target)
    }
    return target
}
