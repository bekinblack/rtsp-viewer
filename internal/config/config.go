package config

import (
	"fmt"
	"os"
	"rtsp-viewer/internal/model"

	"gopkg.in/yaml.v3"
)

const configFile = "config.yaml"

func Load() (model.Form, error) {
	empty := model.Form{}
	data, err := os.ReadFile(configFile)
	if err != nil {
		return empty, err
	}

	var form model.Form
	if err := yaml.Unmarshal(data, &form); err != nil {
		return empty, err
	}

	form.Password, err = Decode(form.Password)
	if err != nil {
		return empty, err
	}

	return form, nil
}

func Save(f model.Form) error {
	f.Password = Encode(f.Password)
	data, err := yaml.Marshal(&f)
	if err != nil {
		return fmt.Errorf("ошибка сериализации YAML: %v", err)
	}

	if err := os.WriteFile(configFile, data, 0644); err != nil {
		return fmt.Errorf("ошибка записи файла: %v", err)
	}

	return nil
}
