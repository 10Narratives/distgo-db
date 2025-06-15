package utils_test

import (
	"testing"

	"github.com/10Narratives/distgo-db/internal/lib/grpc/utils"
	"github.com/stretchr/testify/require"
)

func TestParseName(t *testing.T) {
	type args struct {
		name string
	}
	tests := []struct {
		name string
		args args
		want utils.ParsedName
	}{
		{
			name: "Full path with /documents/",
			args: args{
				name: "databases/mdb/collections/pillars/documents/host_1012",
			},
			want: utils.ParsedName{
				DatabaseID:   "mdb",
				CollectionID: "pillars",
				DocumentID:   "host_1012",
			},
		},
		{
			name: "Without document",
			args: args{
				name: "databases/mdb/collections/pillars",
			},
			want: utils.ParsedName{
				DatabaseID:   "mdb",
				CollectionID: "pillars",
				DocumentID:   "",
			},
		},
		{
			name: "Db only",
			args: args{
				name: "databases/mdb",
			},
			want: utils.ParsedName{
				DatabaseID:   "mdb",
				CollectionID: "",
				DocumentID:   "",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := utils.ParseName(tt.args.name)
			require.Equalf(t, tt.want, got,
				"ParseName(%q) = %+v, want %+v", tt.args.name, got, tt.want)
		})
	}
}
