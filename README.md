## Repository information

This Repository contains all plugins for the [Iris web framework](https://github.com/kataras/iris).

You can contribute also, just make a pull request, try to keep conversion, configuration file: './myplugin/config.go' & plugin: './myplugin/myplugin.go'.


## How can I register a plugin?

```go
iris.Plugins.Add(theplugin)
``` 