package internal

import (
	"gopkg.in/yaml.v3"
	"os"
	"path/filepath"
)

type StackConfig struct {
	UrlPath string `yaml:"urlPath"`
	Port    string `yaml:"port"`
}

type StackConfigServiceImpl struct {
	stackConfigs map[string]StackConfig
}

func (s *StackConfigServiceImpl) GetStackConfig(stackName string) StackConfig {
	if stackConfig, found := s.stackConfigs[stackName]; found {
		return stackConfig
	}
	Logger.Error("error: StackConfig not found for '%s'", stackName)
	return StackConfig{"/", "80"}
}

func ProvideStackConfigService(stackDir string) StackConfigService {
	stackConfigs := make(map[string]StackConfig)

	files, err := os.ReadDir(stackDir)
	if err != nil {
		Logger.Fatal("error when reading directory %s: %w", stackDir, err)
	}

	for _, file := range files {
		if !file.IsDir() {
			continue
		}
		stackConfigFilePath := filepath.Join(stackDir, file.Name(), "app.yml")
		stackConfigs[file.Name()] = loadConfig(stackConfigFilePath)
	}
	return &StackConfigServiceImpl{stackConfigs: stackConfigs}
}

func loadConfig(configPath string) StackConfig {
	config := StackConfig{UrlPath: "/", Port: "80"}
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		Logger.Debug("file %s does not exist, providing default config instead", configPath)
		return config
	}

	fileContent, err := os.ReadFile(configPath)
	if err != nil {
		Logger.Fatal("error when reading file %s: %w", configPath, err)
	}
	if err := yaml.Unmarshal(fileContent, &config); err != nil {
		Logger.Fatal("error when unmarshalling %s: %w", configPath, err)
	}
	return config
}
