//Package ryan is a Router which provides support for optional path parameters at any path part position
// use this plugin with your own risk.
// Note: This router is not fast as the original, I'm not putting too much efford here, it is used only for special cases only when you don't care so much about router's performanrce.
// It is used only on the routes you register with optional parameters, this will not reduce the performanrce of the rest of your routes.

// THIS IS NOT COMPLETED YET, not even started.

package ryan

import (
	"fmt"
	"strings"

	"github.com/kataras/iris"
)

// DefaultSymbol is the symbol for ryan router's optional parameters, which is '?'
const (
	DefaultSymbol = '?'
	wildcardName  = "_iris_optional_parameterized_ryan_router_"
)

// Ryan is the ryan router plugin
type Ryan struct {
	symbol byte
	debug  bool
}

// New creates and returns a Ryan router plugin
// it takes zero arguments, has only two configurable options, debug and symbol which can be changed via funcs.
func New() *Ryan {
	return &Ryan{symbol: DefaultSymbol}
}

// SetSymbol sets/changes the default symbol which is '?'
// returns itself
func (r *Ryan) SetSymbol(symbol byte) *Ryan {
	r.symbol = symbol
	return r
}

// SetDebug sets the debug option
// returns itself
func (r *Ryan) SetDebug(b bool) *Ryan {
	r.debug = b
	return r
}

// PreLookup catch the route before registers and change its path so it is compatible with the original router
func (r *Ryan) PreLookup(route iris.Route) {
	if r.debug {
		fmt.Println("Route with path: " + route.Path() + " just registered")
	}

	path := route.Path()
	if symbolIdx := strings.IndexByte(path, r.symbol); symbolIdx != -1 {
		///TODO: Implementation here.
		newPath := path[0:symbolIdx-1] + wildcardName
		parts := strings.Split(path, "/")
		for _, part := range parts {
			if part[0] == r.symbol {

			}
		}

		// change the path
		route.SetPath(newPath)

		// prepends the middleware for this route, yes each route will have its own 'router' first middleware because we want to use the ryan only when it is nessecary,
		// we don't use iris.UseGlobal feature for performanrce reasons
		route.SetMiddleware(append(iris.Middleware{iris.HandlerFunc(func(ctx *iris.Context) {

		})}, route.Middleware()...))
	}

}
