package typescript

import (
	"errors"
	"os"
	"path/filepath"
	"strings"

	"github.com/iris-contrib/npm"
	"github.com/iris-contrib/plugin/editor"
	"github.com/kataras/iris"
	"github.com/kataras/iris/logger"
	"github.com/kataras/iris/utils"
)

/* Notes

The editor is working when the typescript plugin finds a typescript project (tsconfig.json),
also working only if one typescript project found (normaly is one for client-side).

*/

// Name the name of the plugin, is "TypescriptPlugin"
const Name = "TypescriptPlugin"

type (
	// Plugin the struct of the Typescript Plugin, holds all necessary fields & methods
	Plugin struct {
		config Config
		// taken from Activate
		pluginContainer iris.PluginContainer
		// taken at the PreListen
		logger *logger.Logger
	}
)

// Editor is just a shortcut for github.com/kataras/iris/plugin/editor.New()
// returns a new (Editor)Plugin, it's exists here because the typescript plugin has direct interest with the EditorPlugin
func Editor(username, password string) *editor.Plugin {
	editorCfg := editor.DefaultConfig()
	editorCfg.Username = username
	editorCfg.Password = password
	return editor.New(editorCfg)
}

// Plugin

// New creates & returns a new instnace typescript plugin
func New(cfg ...Config) *Plugin {
	c := DefaultConfig().Merge(cfg)

	if !strings.Contains(c.Ignore, nodeModules) {
		c.Ignore += "," + nodeModules
	}

	return &Plugin{config: c}
}

// implement the IPlugin & IPluginPreListen

// Activate ...
func (t *Plugin) Activate(container iris.PluginContainer) error {
	t.pluginContainer = container
	return nil
}

// GetName ...
func (t *Plugin) GetName() string {
	return Name + "[" + utils.RandomString(10) + "]" // this allows the specific plugin to be registed more than one time
}

// GetDescription TypescriptPlugin scans and compile typescript files with ease
func (t *Plugin) GetDescription() string {
	return Name + " scans and compile typescript files with ease. \n"
}

// PreListen ...
func (t *Plugin) PreListen(s *iris.Framework) {
	t.logger = s.Logger
	t.start()
}

//

// implementation

func (t *Plugin) start() {
	defaultCompilerArgs := t.config.Tsconfig.CompilerArgs() //these will be used if no .tsconfig found.
	if t.hasTypescriptFiles() {
		//Can't check if permission denied returns always exists = true....
		//typescriptModule := out + string(os.PathSeparator) + "typescript" + string(os.PathSeparator) + "bin"
		if !npm.Exists(t.config.Bin) {
			t.logger.Println("Installing typescript, please wait...")
			res := npm.Install("typescript")
			if res.Error != nil {
				t.logger.Print(res.Error.Error())
				return
			}
			t.logger.Print(res.Message)

		}

		projects := t.getTypescriptProjects()
		if len(projects) > 0 {
			watchedProjects := 0
			//typescript project (.tsconfig) found
			for _, project := range projects {
				cmd := utils.CommandBuilder("node", t.config.Bin, "-p", project[0:strings.LastIndex(project, utils.PathSeparator)]) //remove the /tsconfig.json)
				projectConfig := FromFile(project)

				if projectConfig.CompilerOptions.Watch {
					watchedProjects++
					// if has watch : true then we have to wrap the command to a goroutine (I don't want to use the .Start here)
					go func() {
						_, err := cmd.Output()
						if err != nil {
							t.logger.Println(err.Error())
							return
						}
					}()
				} else {

					_, err := cmd.Output()
					if err != nil {
						t.logger.Println(err.Error())
						return
					}

				}

			}
			t.logger.Printf("%d Typescript project(s) compiled ( %d monitored by a background file watcher ) ", len(projects), watchedProjects)
		} else {
			//search for standalone typescript (.ts) files and compile them
			files := t.getTypescriptFiles()

			if len(files) > 0 {
				watchedFiles := 0
				if t.config.Tsconfig.CompilerOptions.Watch {
					watchedFiles = len(files)
				}
				//it must be always > 0 if we came here, because of if hasTypescriptFiles == true.
				for _, file := range files {
					cmd := utils.CommandBuilder("node", t.config.Bin)
					cmd.AppendArguments(defaultCompilerArgs...)
					cmd.AppendArguments(file)
					_, err := cmd.Output()
					cmd.Args = cmd.Args[0 : len(cmd.Args)-1] //remove the last, which is the file
					if err != nil {
						t.logger.Println(err.Error())
						return
					}

				}
				t.logger.Printf("%d Typescript file(s) compiled ( %d monitored by a background file watcher )", len(files), watchedFiles)
			}

		}

		//editor activation
		if len(projects) == 1 && t.config.Editor != nil {
			dir := projects[0][0:strings.LastIndex(projects[0], utils.PathSeparator)]
			t.config.Editor.Dir(dir)
			t.pluginContainer.Add(t.config.Editor)
		}

	}
}

func (t *Plugin) hasTypescriptFiles() bool {
	root := t.config.Dir
	ignoreFolders := strings.Split(t.config.Ignore, ",")
	hasTs := false

	filepath.Walk(root, func(path string, fi os.FileInfo, err error) error {

		if fi.IsDir() {
			return nil
		}
		for i := range ignoreFolders {
			if strings.Contains(path, ignoreFolders[i]) {
				return nil
			}
		}
		if strings.HasSuffix(path, ".ts") {
			hasTs = true
			return errors.New("Typescript found, hope that will stop here")
		}

		return nil
	})
	return hasTs
}

func (t *Plugin) getTypescriptProjects() []string {
	var projects []string
	ignoreFolders := strings.Split(t.config.Ignore, ",")

	root := t.config.Dir
	//t.logger.Printf("\nSearching for typescript projects in %s", root)

	filepath.Walk(root, func(path string, fi os.FileInfo, err error) error {
		if fi.IsDir() {
			return nil
		}
		for i := range ignoreFolders {
			if strings.Contains(path, ignoreFolders[i]) {
				//t.logger.Println(path + " ignored")
				return nil
			}
		}

		if strings.HasSuffix(path, utils.PathSeparator+"tsconfig.json") {
			//t.logger.Printf("\nTypescript project found in %s", path)
			projects = append(projects, path)
		}

		return nil
	})
	return projects
}

// this is being called if getTypescriptProjects return 0 len, then we are searching for files using that:
func (t *Plugin) getTypescriptFiles() []string {
	var files []string
	ignoreFolders := strings.Split(t.config.Ignore, ",")

	root := t.config.Dir
	//t.logger.Printf("\nSearching for typescript files in %s", root)

	filepath.Walk(root, func(path string, fi os.FileInfo, err error) error {
		if fi.IsDir() {
			return nil
		}
		for i := range ignoreFolders {
			if strings.Contains(path, ignoreFolders[i]) {
				//t.logger.Println(path + " ignored")
				return nil
			}
		}

		if strings.HasSuffix(path, ".ts") {
			//t.logger.Printf("\nTypescript file found in %s", path)
			files = append(files, path)
		}

		return nil
	})
	return files
}

//
//
