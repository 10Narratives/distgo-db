// models/collection.go
package models_test

import (
	"testing"

	"githib.com/10Narratives/distgo-db/internal/worker/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCollection_Insert(t *testing.T) {
	t.Parallel()

	type fields struct {
		name string
	}

	type args struct {
		data map[string]any
	}

	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr require.ErrorAssertionFunc
	}{
		{
			name: "insertion",
			args: args{
				data: map[string]any{
					"name": "Test",
					"age":  21,
				},
			},
			fields: fields{
				name: "users",
			},
			wantErr: require.NoError,
		},
	}

	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			c := models.NewCollection(tt.fields.name)
			d := models.NewDocument(tt.args.data)

			err := c.Insert(d)
			tt.wantErr(t, err)
		})
	}

}

func TestCollection_FindByID(t *testing.T) {
	t.Parallel()

	type fields struct {
		name string
	}

	type args struct {
		data map[string]any
	}

	tests := []struct {
		name      string
		fields    fields
		args      args
		wantValue require.ValueAssertionFunc
		wantErr   require.ErrorAssertionFunc
	}{
		{
			name: "successfully found",
			fields: fields{
				name: "users",
			},
			args: args{
				data: map[string]any{
					"name": "User",
					"age":  100,
				},
			},
			wantValue: func(tt require.TestingT, got any, _ ...any) {
				val, ok := got.(*models.Document)
				require.True(tt, ok)

				name, contains := val.Data["name"]
				require.True(tt, contains)
				require.Equal(t, "User", name)

				age, contains := val.Data["age"]
				require.True(t, contains)
				require.Equal(t, 100, age)
			},
			wantErr: require.NoError,
		},
	}

	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			c := models.NewCollection(tt.fields.name)
			d := models.NewDocument(tt.args.data)

			err := c.Insert(d)
			require.NoError(t, err)

			found, err := c.FindByID(d.ID)
			tt.wantValue(t, found)
			tt.wantErr(t, err)
		})
	}
}

func TestCollection_Update(t *testing.T) {
	t.Parallel()

	type fields struct {
		name string
	}

	type args struct {
		original map[string]any
		toUpdate map[string]any
	}

	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr require.ErrorAssertionFunc
	}{
		{
			name: "successful update",
			fields: fields{
				name: "users",
			},
			args: args{
				original: map[string]any{
					"name": "User",
					"age":  100,
				},
				toUpdate: map[string]any{
					"name": "Username",
					"age":  99,
				},
			},
			wantErr: require.NoError,
		},
		{
			name: "nil data update",
			fields: fields{
				name: "users",
			},
			args: args{
				original: map[string]any{
					"name": "User",
					"age":  100,
				},
				toUpdate: nil,
			},
			wantErr: func(tt require.TestingT, err error, i ...interface{}) {
				assert.EqualError(t, err, "nil data", i...)
			},
		},
	}

	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			c := models.NewCollection(tt.fields.name)
			d := models.NewDocument(tt.args.original)

			err := c.Insert(d)
			require.NoError(t, err)

			err = c.Update(d.ID, tt.args.toUpdate)
			tt.wantErr(t, err)
		})
	}
}

func TestCollection_Delete(t *testing.T) {
	t.Parallel()

	type fields struct {
		name string
	}

	type args struct {
		data map[string]any
	}

	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr require.ErrorAssertionFunc
	}{
		{
			name: "successful deletion",
			fields: fields{
				name: "users",
			},
			args: args{
				data: map[string]any{
					"name": "User",
					"age":  100,
				},
			},
			wantErr: require.NoError,
		},
	}

	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			c := models.NewCollection(tt.name)
			d := models.NewDocument(tt.args.data)

			err := c.Insert(d)
			require.NoError(t, err)

			err = c.Delete(d.ID)
			tt.wantErr(t, err)
		})
	}
}

func TestCollection_Exists(t *testing.T) {
	t.Parallel()

	type fields struct {
		name string
	}

	type args struct {
		data map[string]any
	}

	tests := []struct {
		name      string
		fields    fields
		args      args
		wantValue require.BoolAssertionFunc
	}{
		{
			name: "successful deletion",
			fields: fields{
				name: "users",
			},
			args: args{
				data: map[string]any{
					"name": "User",
					"age":  100,
				},
			},
			wantValue: require.True,
		},
	}

	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			c := models.NewCollection(tt.name)
			d := models.NewDocument(tt.args.data)

			err := c.Insert(d)
			require.NoError(t, err)

			exists := c.Exists(d.ID)
			tt.wantValue(t, exists)
		})
	}
}

func TestCollection_GetChangesSince(t *testing.T) {
	// TODO: Implement testing of GetChangesSince
}

func TestCollection_TruncateChanges(t *testing.T) {
	// TODO: Implement testing of TruncateChanges
}
