package service

import (
    "github.com/gin-gonic/gin"
    uuid "github.com/satori/go.uuid"
    "log"
    "net/http"
    "net/http/httputil"
    "net/url"
)

type Service struct {
	Base
	Source        string     // request source
	directorCache func(r *http.Request)
}

/*
 New return a new instance of service
 param base must have a scheme:
 http://xxxx.xxx:1234/magpie-gateway
 <http|https>://<domain|IP>:<port(number)>[/][path]
 */
func New(uuid uuid.UUID, base string) *Service {
	return &Service{
		Base: Base{
			ID:   uuid,
			Type: TypeRest,
            Endpoints: nil,
		},
		Source:    base,
	}
}

func (s *Service) AddPermission(name, desc, key string) error {
    // TODO 
    return nil
}

func (s *Service) director() func(r *http.Request) {
	if s.directorCache == nil {
		u, err := url.Parse(s.Source)
		if err != nil {
			log.Printf("[WARN] Could not parse base location of service %s", s.ID)
			return nil
		}
		s.directorCache = func (req *http.Request) {
			req.URL.Scheme = u.Scheme
			req.URL.Host = u.Host
			// path := strings.Replace(req.URL.Path, fmt.Sprintf("services/%s", s.Path), "", 1)
			// TODO finish reserve proxy
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
		path := c.GetString("path")
		if path == "" {
			c.JSON(http.StatusBadGateway, gin.H{
				"code": http.StatusBadGateway,
				"msg": "Server configuration error #1",
			})
			return
		}
		s.invoke(c)
	}
}
