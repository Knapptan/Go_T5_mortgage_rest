// Тесты пакета config
package config

import (
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLoad_Success(t *testing.T) {
	// Создаем временный YAML-файл
	content := []byte(`
port: 9090
`)
	tmpfile, err := os.CreateTemp("", "config.*.yaml")
	require.NoError(t, err)
	defer os.Remove(tmpfile.Name())

	_, err = tmpfile.Write(content)
	require.NoError(t, err)
	require.NoError(t, tmpfile.Close())

	// Вызываем функцию
	cfg, err := Load(tmpfile.Name())
	require.NoError(t, err)

	// Проверяем результаты
	assert.Equal(t, 9090, cfg.Port)
}

func TestLoad_DefaultPort(t *testing.T) {
	// Создаем конфиг без порта
	content := []byte(``)
	tmpfile, err := os.CreateTemp("", "config.*.yaml")
	require.NoError(t, err)
	defer os.Remove(tmpfile.Name())

	_, err = tmpfile.Write(content)
	require.NoError(t, err)
	require.NoError(t, tmpfile.Close())

	cfg, err := Load(tmpfile.Name())
	require.NoError(t, err)
	assert.Equal(t, 8080, cfg.Port) // проверяем значение по умолчанию
}

func TestLoad_FileNotExist(t *testing.T) {
	_, err := Load("non_existent_file.yaml")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "error reading config")
}

func TestLoad_InvalidYAML(t *testing.T) {
	// Создаем битый YAML
	content := []byte(`port: "should_be_number"`)
	tmpfile, err := os.CreateTemp("", "config.*.yaml")
	require.NoError(t, err)
	defer os.Remove(tmpfile.Name())

	_, err = tmpfile.Write(content)
	require.NoError(t, err)
	require.NoError(t, tmpfile.Close())

	_, err = Load(tmpfile.Name())
	require.Error(t, err)
	assert.Contains(t, err.Error(), "error parsing config")
}

func TestLoad_EmptyConfig(t *testing.T) {
	content := []byte(`{}`)
	tmpfile, err := os.CreateTemp("", "config.*.yaml")
	require.NoError(t, err)
	defer os.Remove(tmpfile.Name())

	_, err = tmpfile.Write(content)
	require.NoError(t, err)
	require.NoError(t, tmpfile.Close())

	cfg, err := Load(tmpfile.Name())
	require.NoError(t, err)
	assert.Equal(t, 8080, cfg.Port) // должно подставиться значение по умолчанию
}

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
			tmpfile, err := os.CreateTemp("", "config.*.yaml")
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
