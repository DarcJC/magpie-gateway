package rest

type Endpoint struct {
	Path string
	Permission string // zero value "" means dont need permission to access
}
