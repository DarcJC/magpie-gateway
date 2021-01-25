package rest

import (
	"github.com/gin-gonic/gin"
	uuid "github.com/satori/go.uuid"
	"log"
	"magpie-gateway/service"
	"net/http/httputil"
	"net/url"
)

type Service struct {
	service.Base
	BaseLocation string     // request source
	Endpoints    []Endpoint // endpoints belong to this service
	proxy *httputil.ReverseProxy // reference of proxy instance
}

/*
 New return a new instance of service
 param base must have a scheme:
 http://xxxx.xxx:1234/magpie-gateway
 <http|https>://<domain|IP>:<port(number)>[/][path]
 */
func New(uuid uuid.UUID, base string) *Service {
	return &Service{
		Base: service.Base{
			ID:   uuid,
			Type: service.TYPE_REST,
		},
		BaseLocation: base,
		Endpoints:    nil,
		proxy: nil,
	}
}

/*
 initProxy initialize reserve proxy of this service
 return true if ok
 */
func (s *Service) initProxy() bool {
	u, err := url.Parse(s.BaseLocation)
	if err != nil {
		return false
	}
	s.proxy = httputil.NewSingleHostReverseProxy(u)
	return true
}

func (s *Service) invoke(c *gin.Context) {
	s.proxy.ServeHTTP(c.Writer, c.Request)
}

/*
 return a handler for gin
 */
func (s *Service) Handler() gin.HandlerFunc {
	if s.proxy == nil {
		if ok := s.initProxy(); !ok {
			log.Printf("[WARN] Could not init proxy to rest service (Service ID: %s)", s.ID)
			return nil
		}
	}

	return func(c *gin.Context) {
		s.invoke(c)
	}
}
