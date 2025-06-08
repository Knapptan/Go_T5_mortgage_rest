// Package config_test содержит модульные тесты для пакета config.
package config_test

import (
	"fmt"
	"os"
	"testing"

	. "github.com/Knapptan/Go_T5_mortgage_rest/config"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestLoad_Success проверяет успешную загрузку конфигурации с заданным портом.
func TestLoad_Success(t *testing.T) {
	content := []byte(`port: 9090`)
	tmpfile, err := os.CreateTemp("", "*.yaml")
	require.NoError(t, err)
	defer os.Remove(tmpfile.Name())

	_, err = tmpfile.Write(content)
	require.NoError(t, err)
	require.NoError(t, tmpfile.Close())

	cfg, err := Load(tmpfile.Name())
	require.NoError(t, err)
	assert.Equal(t, 9090, cfg.Port)
}

// TestLoad_DefaultPort проверяет, что по умолчанию устанавливается порт 8080, если он не задан.
func TestLoad_DefaultPort(t *testing.T) {
	content := []byte(``)
	tmpfile, err := os.CreateTemp("", "*.yaml")
	require.NoError(t, err)
	defer os.Remove(tmpfile.Name())

	_, err = tmpfile.Write(content)
	require.NoError(t, err)
	require.NoError(t, tmpfile.Close())

	cfg, err := Load(tmpfile.Name())
	require.NoError(t, err)
	assert.Equal(t, 8080, cfg.Port)
}

// TestLoad_FileNotExist проверяет поведение при отсутствии конфигурационного файла.
func TestLoad_FileNotExist(t *testing.T) {
	_, err := Load("non_existent_file.yaml")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "error reading config")
}

// TestLoad_InvalidYAML проверяет, что возвращается ошибка при неверном формате YAML.
func TestLoad_InvalidYAML(t *testing.T) {
	content := []byte(`port: "should_be_number"`)
	tmpfile, err := os.CreateTemp("", "*.yaml")
	require.NoError(t, err)
	defer os.Remove(tmpfile.Name())

	_, err = tmpfile.Write(content)
	require.NoError(t, err)
	require.NoError(t, tmpfile.Close())

	_, err = Load(tmpfile.Name())
	require.Error(t, err)
	assert.Contains(t, err.Error(), "error parsing config")
}

// TestLoad_EmptyConfig проверяет поведение при пустом объекте конфигурации.
func TestLoad_EmptyConfig(t *testing.T) {
	content := []byte(`{}`)
	tmpfile, err := os.CreateTemp("", "*.yaml")
	require.NoError(t, err)
	defer os.Remove(tmpfile.Name())

	_, err = tmpfile.Write(content)
	require.NoError(t, err)
	require.NoError(t, tmpfile.Close())

	cfg, err := Load(tmpfile.Name())
	require.NoError(t, err)
	assert.Equal(t, 8080, cfg.Port)
}

// TestLoad_PortBoundaries проверяет корректность обработки крайних значений порта.
func TestLoad_PortBoundaries(t *testing.T) {
	tests := []struct {
		name     string
		port     int
		expected int
	}{
		{"min port", 1, 1},
		{"max port", 65535, 65535},
		{"zero port", 0, 8080},
		{"negative port", -1, 8080},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			content := []byte(fmt.Sprintf("port: %d", tt.port))
			tmpfile, err := os.CreateTemp("", "*.yaml")
			require.NoError(t, err)
			defer os.Remove(tmpfile.Name())

			_, err = tmpfile.Write(content)
			require.NoError(t, err)
			require.NoError(t, tmpfile.Close())

			cfg, err := Load(tmpfile.Name())
			require.NoError(t, err)
			assert.Equal(t, tt.expected, cfg.Port)
		})
	}
}
