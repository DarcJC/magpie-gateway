package rest

import (
	"fmt"
	"github.com/gin-gonic/gin"
	uuid "github.com/satori/go.uuid"
	"log"
	"magpie-gateway/service"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
)

type Service struct {
	service.Base
	BaseLocation string     // request source
	Endpoints    []Endpoint // endpoints belong to this service
	directorCache func(r *http.Request)
}

/*
 New return a new instance of service
 param base must have a scheme:
 http://xxxx.xxx:1234/magpie-gateway
 <http|https>://<domain|IP>:<port(number)>[/][path]
 */
func New(uuid uuid.UUID, base string, path string) *Service {
	return &Service{
		Base: service.Base{
			ID:   uuid,
			Type: service.TypeRest,
			Path: path,
		},
		BaseLocation: base,
		Endpoints:    nil,
	}
}

func (s *Service) director() func(r *http.Request) {
	if s.directorCache == nil {
		u, err := url.Parse(s.BaseLocation)
		if err != nil {
			log.Printf("[WARN] Could not parse base location of service %s", s.ID)
			return nil
		}
		s.directorCache = func (req *http.Request) {
			req.URL.Scheme = u.Scheme
			req.URL.Host = u.Host
			path := strings.Replace(req.URL.Path, fmt.Sprintf("services/%s", s.Path), "", 1)
			req.URL.Path = path
		}
	}
	return s.directorCache
}

func (s *Service) invoke(c *gin.Context) {
	proxy := &httputil.ReverseProxy{Director: s.director()}
	proxy.ServeHTTP(c.Writer, c.Request)
}

/*
 return a handler for gin
 */
func (s *Service) Handler() gin.HandlerFunc {
	return func(c *gin.Context) {
		s.invoke(c)
	}
}
