package cli

import (
	"fmt"
	"github.com/getevo/evo/v2/lib/generic"
	"os"
	"os/exec"
	"runtime"
	"strings"
)

var commands []Command

type Params []Param

func (params Params) Get(s string) generic.Value {
	for _, p := range params {
		if p.Switch == s {
			return p.Value
		}
	}
	return generic.Parse(nil)
}

type Context struct {
	Params Params
}

func (c *Context) Param(s string) generic.Value {
	return c.Params.Get(s)
}

func (c *Context) Log(s string, params ...interface{}) {
	fmt.Printf("[LOG] "+s+"\n", params...)
}

func (c *Context) Debug(s string, params ...interface{}) {
	fmt.Printf("[DEBUG] "+s+"\n", params...)
}

func (c *Context) Error(s string, params ...interface{}) {
	fmt.Printf("[ERR] "+s, params...)
}

func (c *Context) Fatal(s string, params ...interface{}) {
	fmt.Printf("[FATAL] "+s+"\n", params...)
}

func (c *Context) Print(s string, params ...interface{}) {
	fmt.Printf(s+"\n", params...)
}

func (c *Context) ClearScreen() {
	ClearTerminal()
}

type Param struct {
	Switch       string
	Usage        string
	Required     bool
	Value        generic.Value
	DefaultValue interface{}
	set          bool
}

type Command struct {
	Switch   string
	Usage    string
	Params   Params
	Action   func(params *Context)
	Commands []Command
	Default  bool
}

func (c Command) getHelp() string {
	var builder = ""
	builder += fmt.Sprintf("  %s \t %s\n", c.Switch, c.Usage)
	for _, param := range c.Params {
		remarks := ""
		if param.Required {
			remarks = " (required)"
		}
		if param.DefaultValue != nil {
			remarks += fmt.Sprintf(" (default: %v)", param.DefaultValue)
		}
		builder += fmt.Sprintf("    -%s \t %s%s\n", param.Switch, param.Usage, remarks)
	}
	return builder
}

func Register(c ...Command) {
	commands = append(commands, c...)
}

func Run() {
	parseFlags(1, commands)
}

func parseFlags(pos int, commands []Command) {
	var found = false
	if len(os.Args) < 2 {
		return
	}
	if os.Args[pos] == "help" || os.Args[pos] == "-h" {
		fmt.Println("Usage: ")
		for _, command := range commands {
			fmt.Println(command.getHelp())
		}
		os.Exit(0)
	}
	for idx, _ := range commands {
		var command = commands[idx]

		if command.Switch == os.Args[pos] {
			var ctx = Context{}

			if len(os.Args) > pos+1 && os.Args[pos+1] == "help" || os.Args[pos+1] == "-h" {
				fmt.Println("Usage: ")
				for _, cmd := range command.Commands {
					fmt.Println(cmd.getHelp())
				}
				os.Exit(0)
			}

			for j := 0; j < len(command.Params); j++ {
				command.Params[j].Value = generic.Parse(command.Params[j].DefaultValue)
			}
			pos = parseParams(pos+1, &command)
			ctx.Params = command.Params
			command.Action(&ctx)
			found = true
			if len(command.Commands) > 0 && len(os.Args) > pos {
				parseFlags(pos, command.Commands)
			}
		}
	}

	if !found {
		fmt.Printf("Unknown command: %s\n", os.Args[pos])
		os.Exit(1)
	}

}

func parseParams(pos int, command *Command) int {
	defer func() {
		for k := 0; k < len(command.Params); k++ {
			if !command.Params[k].set && command.Params[k].Required {
				fmt.Printf("Missing required parameter: -%s\n", command.Params[k].Switch)
				os.Exit(1)
			}
		}
	}()
	var p = pos
	for p < len(os.Args) {
		if os.Args[p][0] != '-' {
			return pos
		}
		var found bool
		for k := 0; k < len(command.Params); k++ {
			var v = strings.SplitN(os.Args[p], "=", 2)
			if "-"+command.Params[k].Switch == v[0] {
				found = true
				pos++
				if len(v) == 2 {
					command.Params[k].Value = generic.Parse(v[1])
					command.Params[k].set = true
				}
				break
			}
		}
		p++
		if !found {
			fmt.Printf("Unknown parameter: %s\n", os.Args[pos-1])
			os.Exit(1)
		}
	}
	return pos
}

func ClearTerminal() {
	switch runtime.GOOS {
	case "darwin":
		runCmd("clear")
	case "linux":
		runCmd("clear")
	case "windows":
		runCmd("cmd", "/c", "cls")
	default:
		runCmd("clear")
	}
}

func runCmd(name string, arg ...string) {
	cmd := exec.Command(name, arg...)
	cmd.Stdout = os.Stdout
	_ = cmd.Run()
}
