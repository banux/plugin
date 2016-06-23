package oauth

import (
	"github.com/iris-contrib/gothic"
	"github.com/kataras/iris"
	"github.com/markbates/goth"
)

// Plugin is a plugin which helps you to use OAuth/OAuth2 apis from famous websites
// See more at https://github.com/iris-contrib/gothic
type Plugin struct {
	Config          Config
	successHandlers []iris.HandlerFunc
	failHandler     iris.HandlerFunc
}

// New returns a new OAuth plugin
// receives one parameter of type 'Config'
func New(cfg Config) *Plugin {
	c := DefaultConfig().MergeSingle(cfg)
	return &Plugin{Config: c}
}

// OnSuccess registers handler(s) which fires when the user logged in successfully
func (p *Plugin) OnSuccess(handlersFn ...iris.HandlerFunc) {
	p.successHandlers = append(p.successHandlers, handlersFn...)
}

// OnFail registers handler which fires when the user failed to logged in
// underhood it justs registers an error handler to the StatusUnauthorized(400 status code), same as 'iris.OnError(400,handler)'
func (p *Plugin) OnFail(handler iris.HandlerFunc) {
	p.failHandler = handler
}

// User returns the user for the particular client
// if user is not validated  or not found it returns nil
// same as 'ctx.Get(config's ContextKey field).(goth.User)'
func (p *Plugin) User(ctx *iris.Context) (u goth.User) {
	return ctx.Get(p.Config.ContextKey).(goth.User)
}

// PreListen init the providers and the routes before server's listens
func (p *Plugin) PreListen(s *iris.Framework) {
	oauthProviders := p.Config.GenerateProviders(s.HTTPServer.FullHost())
	if len(oauthProviders) > 0 {
		goth.UseProviders(oauthProviders...)
		// set the mux path to handle these providers
		s.Get(p.Config.Path+"/:provider", func(ctx *iris.Context) {
			err := gothic.BeginAuthHandler(ctx)
			if err != nil {
				s.Logger.Warningf("\n[IRIS OAUTH MIDDLEWARE] Error:" + err.Error())
			}
		})

		authMiddleware := func(ctx *iris.Context) {
			user, err := gothic.CompleteUserAuth(ctx)
			if err != nil {
				ctx.EmitError(iris.StatusUnauthorized)
				ctx.Log(err.Error())
				return
			}
			ctx.Set(p.Config.ContextKey, user)
			ctx.Next()
		}

		p.successHandlers = append([]iris.HandlerFunc{authMiddleware}, p.successHandlers...)

		s.Get(p.Config.Path+"/:provider/callback", p.successHandlers...)(p.Config.RouteName)

		// register the error handler
		if p.failHandler != nil {
			s.OnError(iris.StatusUnauthorized, p.failHandler)
		}
	}
}
