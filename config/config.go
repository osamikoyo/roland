package config

type LLMs struct{
	Parser string `yaml:"parser"`
	Listener string `yaml:"listener"`
}

type Config struct {
	LLMs LLMs `yaml:"llms"`
	LLamaLibPath     string                       `yaml:"llm_lib_path"`

	ResponseLength int32                        `yaml:"response_length"`
	CMDS           map[string]map[string]string `yaml:"cmds"`
	SaveCMDOutput  bool                         `yaml:"save_cmd_output"`
}
