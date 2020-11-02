package reader

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestReadConfig_Files(t *testing.T) {
	a := assert.New(t)
	rc := ReadConfig{
		Dir:       "../../example/pkg/",
		Recursive: true,
	}

	files, err := rc.Files()
	a.NoError(err)
	a.NotEmpty(files)
}
