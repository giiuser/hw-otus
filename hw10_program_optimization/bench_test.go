package hw10programoptimization

import (
	"archive/zip"
	"testing"

	"github.com/stretchr/testify/require"
)

func BenchmarkStat(b *testing.B) {
	uz, err := zip.OpenReader("testdata/users.dat.zip")
	if err == nil {
		defer uz.Close()
	}

	data, err := uz.File[0].Open()
	require.NoError(b, err)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		GetDomainStat(data, "com")
	}
}
