// Package injection provides the injection engine
//
// The injection engine is used to manage the dependencies of the application
package injection

func New() interface{} {
	return &engine{}
}
