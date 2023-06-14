package command

var Commands = map[string]Command{
	"help":    NewHelpCommand(),
	"init":    NewInitCommand(),
	"dump":    NewDumpCommand(),
	"update":  NewUpdateCommand(),
	"ls":      NewLsCommand(),
	"get":     NewGetCommand(),
	"set":     NewSetCommand(),
	"rm":      NewRmCommand(),
	"version": NewVersionCommand(),
}
