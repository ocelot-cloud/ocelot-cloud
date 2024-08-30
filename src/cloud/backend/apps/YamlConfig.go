package apps

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
	logger.Error("error: StackConfig not found for '%s'", stackName)
	return StackConfig{"/", "80"}
}

func provideStackConfigService(stackDir string) StackConfigService {
	stackConfigs := make(map[string]StackConfig)

	files, err := os.ReadDir(stackDir)
	if err != nil {
		logger.Fatal("error when reading directory %s: %v", stackDir, err)
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
	newConfig := StackConfig{UrlPath: "/", Port: "80"}
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		logger.Debug("file %s does not exist, providing default config instead", configPath)
		return newConfig
	}

	fileContent, err := os.ReadFile(configPath)
	if err != nil {
		logger.Fatal("error when reading file %s: %w", configPath, err)
	}
	if err = yaml.Unmarshal(fileContent, &newConfig); err != nil {
		logger.Fatal("error when unmarshalling %s: %w", configPath, err)
	}
	return newConfig
}
