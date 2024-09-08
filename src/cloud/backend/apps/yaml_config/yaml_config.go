package yaml_config

import (
	"gopkg.in/yaml.v3"
	"ocelot/backend/tools"
	"os"
	"path/filepath"
)

type ConfigServiceType interface {
	GetAppConfig(appName string) appConfig
}

var logger = tools.Logger

type appConfig struct {
	UrlPath string `yaml:"urlPath"`
	Port    string `yaml:"port"`
}

type configServiceImpl struct {
	stackConfigs map[string]appConfig
}

func (s *configServiceImpl) GetAppConfig(stackName string) appConfig {
	if stackConfig, found := s.stackConfigs[stackName]; found {
		return stackConfig
	}
	logger.Error("error: StackConfig not found for '%s'", stackName)
	return appConfig{"/", "80"}
}

func ProvideAppConfigService(stackDir string) ConfigServiceType {
	stackConfigs := make(map[string]appConfig)

	files, err := os.ReadDir(stackDir)
	if err != nil {
		a, _ := os.Getwd()
		logger.Debug("current dir is: %s", a)
		logger.Fatal("error when reading directory %s: %v", stackDir, err)
	}

	for _, file := range files {
		if !file.IsDir() {
			continue
		}
		stackConfigFilePath := filepath.Join(stackDir, file.Name(), "app.yml")
		stackConfigs[file.Name()] = loadConfig(stackConfigFilePath)
	}
	return &configServiceImpl{stackConfigs: stackConfigs}
}

func loadConfig(configPath string) appConfig {
	newConfig := appConfig{UrlPath: "/", Port: "80"}
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
