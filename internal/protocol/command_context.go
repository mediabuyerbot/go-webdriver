package protocol

type CommandContext struct {
	id     string
	client Client
}

func NewCommandContext(cli Client, id string) *CommandContext {
	return &CommandContext{
		id:     id,
		client: cli,
	}
}
