package commands

type Param struct {
	Name     string `yaml:"name"`
	Type     string `yaml:"type"`
	Required bool   `yaml:"required"`
	Default  string `yaml:"default"`
	Prompt   string `yaml:"prompt"`
}

type Command struct {
	Name        string  `yaml:"name"`
	Description string  `yaml:"description"`
	Params      []Param `yaml:"params"`
	Command     string  `yaml:"command"`
	Sudo        bool    `yaml:"sudo"`
}

type CommandCategory struct {
	Category string    `yaml:"category"`
	Icon     string    `yaml:"icon"`
	Commands []Command `yaml:"commands"`
}
