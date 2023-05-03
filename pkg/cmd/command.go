package cmd

type Command struct {
	Parent      string
	Name        string
	Description string
	Aliases     []string
	Args        []string
	Hidden      bool
}
