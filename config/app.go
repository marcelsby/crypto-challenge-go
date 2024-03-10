package config

import (
	"encoding/hex"
	"fmt"
	"strings"

	"github.com/cristalhq/aconfig"
	"github.com/cristalhq/aconfig/aconfigdotenv"
)

// Criar uma struct com as propriedades requisitadas pela aplicação
type AppConfig struct {
	Database struct {
		User     string
		Password string
		Host     string `default:"localhost"`
		Port     int    `default:"3306"`
		DbName   string `env:"NAME"`
	}

	Cryptography struct {
		SecretKey string
	}
}

func GetAppConfig(configFilePath string) *AppConfig {
	var cfg AppConfig

	loader := aconfig.LoaderFor(&cfg, aconfig.Config{
		SkipEnv:            true,
		SkipFlags:          true,
		FailOnFileNotFound: true,
		FileDecoders: map[string]aconfig.FileDecoder{
			".env": aconfigdotenv.New(),
		},
		Files: []string{configFilePath},
	})

	if err := loader.Load(); err != nil {
		panic(err)
	}

	validationResult := validateAppConfig(&cfg)

	if len(validationResult) > 0 {
		for k, v := range validationResult {
			fmt.Printf("%s:\n", k)
			for _, v := range *v {
				fmt.Printf("\t- %s\n", v)
			}
		}
		panic(".env validation error")
	}

	return &cfg
}

func validateAppConfig(cfg *AppConfig) map[string]*[]string {
	validationErrors := make(map[string]*[]string)

	fieldsToMakeBlankValidation := map[string]string{
		"Database.User":          cfg.Database.User,
		"Database.Password":      cfg.Database.Password,
		"Database.DbName":        cfg.Database.DbName,
		"Cryptography.SecretKey": cfg.Cryptography.SecretKey,
	}

	blankFieldValidationResults := validateBlankFields(fieldsToMakeBlankValidation)

	if len(fieldsToMakeBlankValidation) > 0 {
		for k, v := range blankFieldValidationResults {
			if _, ok := validationErrors[k]; !ok {
				validationErrors[k] = &[]string{}
			}

			*validationErrors[k] = append(*validationErrors[k], v)
		}
	}

	secretKeyValidationErrors := []string{}

	if decodedSecretKey, err := hex.DecodeString(cfg.Cryptography.SecretKey); err != nil {
		secretKeyValidationErrors = append(secretKeyValidationErrors, err.Error())
	} else if len(decodedSecretKey) != 32 {
		secretKeyValidationErrors = append(secretKeyValidationErrors, "Must represent exactly 32 bytes (64 hex-characters).")
	}

	if len(secretKeyValidationErrors) > 0 {
		key := "Cryptography.SecretKey"

		if _, ok := validationErrors[key]; !ok {
			validationErrors[key] = &secretKeyValidationErrors
		} else {
			concatenedErrors := append(*validationErrors[key], secretKeyValidationErrors...)
			validationErrors[key] = &concatenedErrors
		}
	}

	return validationErrors
}

func validateBlankFields(fieldNameAndValue map[string]string) map[string]string {
	validationResult := make(map[string]string)

	for k, v := range fieldNameAndValue {
		if len(strings.TrimSpace(v)) == 0 {
			validationResult[k] = "Must be a non-blank string."
		}
	}

	return validationResult
}
