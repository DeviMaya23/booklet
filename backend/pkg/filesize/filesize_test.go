package filesize_test

import (
	"testing"

	"github.com/devi/booklet/pkg/filesize"
	"github.com/stretchr/testify/assert"
)

func TestMaxBytesForExtension_VideoReturns2GB(t *testing.T) {
	assert.Equal(t, int64(2*1024*1024*1024), filesize.MaxBytesForExtension("mp4"))
}

func TestMaxBytesForExtension_LayeredReturns800MB(t *testing.T) {
	assert.Equal(t, int64(800*1024*1024), filesize.MaxBytesForExtension("psd"))
}

func TestMaxBytesForExtension_UnknownReturns100MBFallback(t *testing.T) {
	assert.Equal(t, int64(100*1024*1024), filesize.MaxBytesForExtension("xyz"))
}
