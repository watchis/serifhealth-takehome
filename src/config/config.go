package config

type Config struct {
	FilePath string `yaml:"file_path"`

	StateName         string `yaml:"state_name"`
	StateAbbreviation string `yaml:"state_abbreviation"`
	IncludeHighmark   bool   `yaml:"include_highmark"`
}
