package embedfs

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/things-go/embedfs/testdata"
)

func TestDir_Open(t *testing.T) {
	file := Dir{FS: testdata.Staticfs, Dir: "static/css"}
	_, err := file.Open("3.png")
	require.NoError(t, err)
}
