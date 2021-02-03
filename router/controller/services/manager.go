package services

import (
    "github.com/gin-gonic/gin"
    uuid "github.com/satori/go.uuid"
    "magpie-gateway/router"
    "magpie-gateway/router/controller"
    "magpie-gateway/service"
    "net/http"
)

type ServiceManagerEndpoint struct {
    controller.EndpointBase
}

type SMEPutData struct {
    ID string `json:"id" binding:"required"`
    Name string `json:"name" binding:"required"`
    Desc string `json:"desc" binding:"required"`
    Source string `json:"source" binding:"required"`
}

type SMEDeleteData struct {
    ID string `json:"id" binding:"required"`
}

func (s *ServiceManagerEndpoint) Put(ctx *gin.Context) {
    var data SMEPutData

    if err := ctx.BindJSON(&data); err != nil {
        ctx.JSON(http.StatusBadRequest, gin.H{
            "code": http.StatusBadRequest,
            "msg": "magpie could not bind data from this request",
        })
        return
    }

    u, err := uuid.FromString(data.ID)
    if err != nil || u == uuid.Nil {
        ctx.JSON(http.StatusBadRequest, gin.H{
            "code": http.StatusBadRequest,
            "msg": "bad field id",
        })
        return
    }

    e := service.GetServiceEngine()
    if err := e.Manager.CreateService(u.String(), data.Name, data.Desc, data.Source); err != nil {
        ctx.JSON(http.StatusConflict, gin.H{
            "code": http.StatusConflict,
            "msg": err.Error(),
        })
        return
    }

    ctx.JSON(http.StatusCreated, gin.H{
        "code": http.StatusCreated,
        "msg": "success",
    })

}

func (s *ServiceManagerEndpoint) Delete(ctx *gin.Context) {
    var data SMEDeleteData

    if err := ctx.BindJSON(&data); err != nil {
        ctx.JSON(http.StatusBadRequest, gin.H{
            "code": http.StatusBadRequest,
            "msg": "magpie could not bind data from this request",
        })
        return
    }

    u, err := uuid.FromString(data.ID)
    if err != nil {
        ctx.JSON(http.StatusBadRequest, gin.H{
            "code": http.StatusBadRequest,
            "msg": "bad field id",
        })
        return
    }

    e := service.GetServiceEngine()
    ser := e.Manager.GetService(u.String())

    if ser == nil {
        ctx.JSON(http.StatusNotFound, gin.H{
            "code": http.StatusNotFound,
            "msg": "service not found",
        })
        return
    }

    if err := ser.Deactivate(); err != nil {
        ctx.JSON(http.StatusInternalServerError, gin.H{
            "code": http.StatusInternalServerError,
            "msg": "could not set activated flag to this service",
        })
        return
    }

    ctx.JSON(http.StatusOK, gin.H{
        "code": http.StatusOK,
        "msg": "success",
    })
}

func init() {
    endpoint := &ServiceManagerEndpoint{
        EndpointBase: controller.EndpointBase{
            Path: "manage",
        },
    }
    endpoint.Register(endpoint, router.Router.Group("/service"))
}
