package plugins

// RegisterHandler interface for handling route registration
type RegisterHandler interface {
	Register(route string, method string, input interface{}, output interface{})
}
