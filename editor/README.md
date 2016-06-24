## Package information

Editor Plugin is just a bridge between Iris and [alm-tools](http://alm.tools).


[alm-tools](http://alm.tools) is a typescript online IDE/Editor, made by [@basarat](https://twitter.com/basarat) one of the top contributors of the [Typescript](http://www.typescriptlang.org).

Iris gives you the opportunity to edit your client-side using the alm-tools editor, via the editor plugin.


This plugin starts it's own server, if Iris server is using TLS then the editor will use the same key and cert.


## Install

```sh
$ go get -u github.com/iris-contrib/plugin/editor
```

## How to use

```go

package main

import (
	"github.com/kataras/iris"
	"github.com/iris-contrib/plugin/editor"
)

func main(){
	editorConfig:= editor.Config {
		Host: "127.0.0.1",
		Port: 4444,
		WorkingDir: "/path/to/the/client/side/directory",
		Username: "myusername",
		Password: "mypassword"

	}

	iris.Plugins.Add(editor.New(editorConfig))

	iris.Get("/", func (ctx *iris.Context){})

	iris.Listen(":8080")
}


```

> Note for username, password: The Authorization specifies the authentication mechanism (in this case Basic) followed by the username and password.
Although, the string aHR0cHdhdGNoOmY= may look encrypted it is simply a base64 encoded version of <username>:<password>.
Would be readily available to anyone who could intercept the HTTP request. [Read more.](https://www.httpwatch.com/httpgallery/authentication/)

> The editor can't work if the directory doesn't contains a [tsconfig.json](http://www.typescriptlang.org/docs/handbook/tsconfig.json.html).

> If you are using the [typescript plugin](https://github.com/iris-contrib/plugin/tree/master/typescript) you don't have to call the .Dir(...)
