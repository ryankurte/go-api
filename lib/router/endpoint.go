package router

// Endpoint for internal use
type endpoint struct {
	path string
	// Input object
	i interface{}
	// Output object
	o interface{}
	// Base function
	f interface{}
	// Wrapped function
	w interface{}
}
