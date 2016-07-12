package ryan

import "github.com/kataras/iris"

func initDefault() {
	iris.Default = iris.New()
	iris.Config = iris.Default.Config
	iris.Logger = iris.Default.Logger
	iris.Plugins = iris.Default.Plugins
	iris.Websocket = iris.Default.Websocket
	iris.Servers = iris.Default.Servers
	iris.Available = iris.Default.Available
}

/*
func TestRyan(t *testing.T) {
	iris.Get("/path1/?optional1/should/continue/even/?this/doesnt/exists", func(ctx *iris.Context) {
		ctx.Write(ctx.PathString())
	})

	e := iris.Tester(t)

	testPath1 := "/path1/optionalishere/should/continue/even/andthisishere/doesnt/exists"
	e.GET(testPath1).Expect().Status(StatusOK).Body().Equal(testPath1)
	testPath2 := "/path1/optionalishere/should/continue/even/doesnt/exists"
	e.GET(testPath2).Expect().Status(StatusOK).Body().Equal(testPath2)
	testPath3 := "/path1/should/continue/even/andthisishere/doesnt/exists"
	e.GET(testPath3).Expect().Status(StatusOK).Body().Equal(testPath3)
	testPath4 := "/path1/should/continue/even/doesnt/exists"
	e.GET(testPath4).Expect().Status(StatusOK).Body().Equal(testPath4)
	testPath5NotFound := "/path1/wrong/should/continue/even/doesnt/exists"
	e.GET(testPath5NotFound).Expect().Status(StatusNotFound)
}*/

func ExampleRyan_PreLookup() {
	initDefault()
	r := New().SetDebug(true)
	iris.Plugins.Add(r)

	iris.Get("/testpath", func(ctx *iris.Context) {})
	// Output:
	// Route with path: /testpath just registered

	iris.ListenVirtual()
}
