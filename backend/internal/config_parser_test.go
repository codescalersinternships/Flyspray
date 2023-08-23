package internal

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

var validConfig = `
{
  "server": {
    "host": "localhost",
    "port": 8080
  },
  "mailSender": {
    "email": "email",
    "sendgrid_key": "key",
    "timeout": 30
  },
  "db": {
    "file": "filepath"
  },
  "jwt": {
    "secret": "key",
    "timeout": 5
  }
}
`

func TestReadConfigFile(t *testing.T) {
	t.Run("valid", func(t *testing.T) {
		dir := t.TempDir()
		configPath := filepath.Join(dir, "/config.json")

		err := os.WriteFile(configPath, []byte(validConfig), 0644)
		assert.Nil(t, err)

		data, err := ReadConfigFile(configPath)
		assert.Nil(t, err)
		assert.NotEmpty(t, data)
	})

	t.Run("have not access to open", func(t *testing.T) {
		dir := t.TempDir()
		configPath := filepath.Join(dir, "/config.json")

		err := os.WriteFile(configPath, []byte(validConfig), 0000)
		assert.Nil(t, err)

		_, err = ReadConfigFile(configPath)
		assert.NotNil(t, err)
	})

	InvalidTests := []struct {
		name  string
		input string
	}{
		{
			name:  "not json",
			input: ``,
		}, {
			name: "empty host",
			input: `{
				"server": {
					"host": "",
					"port": 8080
				},
				"mailSender": {
					"email": "email",
					"sendgrid_key": "key",
					"timeout": 30
				},
				"db": {
					"file": "filepath"
				},
				"jwt": {
					"secret": "key",
					"timeout": 5
				}
			}`,
		}, {
			name: "missed port",
			input: `{
				"server": {
					"host": "localhost",
				},
				"mailSender": {
					"email": "email",
					"sendgrid_key": "key",
					"timeout": 30
				},
				"db": {
					"file": "filepath"
				},
				"jwt": {
					"secret": "key",
					"timeout": 5
				}
			}`,
		}, {
			name: "empty mail",
			input: `{
				"server": {
					"host": "localhost",
					"port": 8080
				},
				"mailSender": {
					"email": "",
					"sendgrid_key": "key",
					"timeout": 30
				},
				"db": {
					"file": "filepath"
				},
				"jwt": {
					"secret": "key",
					"timeout": 5
				}
			}`,
		}, {
			name: "empty sendgrid_key",
			input: `{
				"server": {
					"host": "localhost",
					"port": 8080
				},
				"mailSender": {
					"email": "email",
					"sendgrid_key": "",
					"timeout": 30
				},
				"db": {
					"file": "filepath"
				},
				"jwt": {
					"secret": "key",
					"timeout": 5
				}
			}`,
		}, {
			name: "missed mail sender timeout",
			input: `{
				"server": {
					"host": "localhost",
					"port": 8080
				},
				"mailSender": {
					"email": "email",
					"sendgrid_key": "key",
				},
				"db": {
					"file": "filepath"
				},
				"jwt": {
					"secret": "key",
					"timeout": 5
				}
			}`,
		}, {
			name: "empty db file path",
			input: `{
				"server": {
					"host": "localhost",
					"port": 8080
				},
				"mailSender": {
					"email": "email",
					"sendgrid_key": "key",
					"timeout": 30
				},
				"db": {
					"file": ""
				},
				"jwt": {
					"secret": "key",
					"timeout": 5
				}
			}`,
		}, {
			name: "empty jwt secret",
			input: `{
				"server": {
					"host": "localhost",
					"port": 8080
				},
				"mailSender": {
					"email": "email",
					"sendgrid_key": "key",
					"timeout": 30
				},
				"db": {
					"file": "filepath"
				},
				"jwt": {
					"secret": "",
					"timeout": 5
				}
			}`,
		}, {
			name: "missed jwt timeout",
			input: `{
				"server": {
					"host": "localhost",
					"port": 8080
				},
				"mailSender": {
					"email": "email",
					"sendgrid_key": "key",
					"timeout": 30
				},
				"db": {
					"file": "filepath"
				},
				"jwt": {
					"secret": "key",
				}
			}`,
		},
		{
			name: "mail sender timeout less than minimum",
			input: `{
				"server": {
					"host": "localhost",
					"port": 8080
				},
				"mailSender": {
					"email": "email",
					"sendgrid_key": "key",
					"timeout": 29
				},
				"db": {
					"file": "filepath"
				},
				"jwt": {
					"secret": "key",
					"timeout": 5
				}
			}`,
		}, {
			name: "jwt timeout less than minimum",
			input: `{
				"server": {
					"host": "localhost",
					"port": 8080
				},
				"mailSender": {
					"email": "email",
					"sendgrid_key": "key",
					"timeout": 30
				},
				"db": {
					"file": "filepath"
				},
				"jwt": {
					"secret": "key",
					"timeout": 4
				}
			}`,
		},
	}

	for _, tc := range InvalidTests {
		t.Run(tc.name, func(t *testing.T) {
			dir := t.TempDir()
			configPath := filepath.Join(dir, "/config.json")

			err := os.WriteFile(configPath, []byte(tc.input), 0644)
			assert.Nil(t, err)

			_, err = ReadConfigFile(configPath)
			assert.NotNil(t, err)
		})
	}
}
