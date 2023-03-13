package cmd

type CmdRun struct {
	// Command to run
	Command string `json:"command" binding:"required"`

	// List of environment variables
	Envs []string `json:"envs,omitempty"`

	// Specifies the current working directory for command execution
	Cwd string `json:"cwd,omitempty"`
}

type RunCommandResponse struct {
	Output string
}
