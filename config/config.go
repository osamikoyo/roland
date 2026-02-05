package config

type Config struct {
	ModelFile  string `yaml:"model_file"`
	LLMLibPath string `yaml:"llm_lib_path"`
	ResponseLength int32 `yaml:"response_length"`
	CMDS map[string]map[string]string `yaml:"cmds"`
}
