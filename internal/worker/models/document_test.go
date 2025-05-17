package models_test

import (
	"testing"

	"githib.com/10Narratives/distgo-db/internal/worker/models"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func TestDocument_Get(t *testing.T) {
	t.Parallel()

	type fields struct {
		id   uuid.UUID
		data map[string]any
	}

	type args struct {
		path string
	}

	tests := []struct {
		name      string
		fields    fields
		args      args
		wantValue require.ValueAssertionFunc
		wantError require.ErrorAssertionFunc
	}{
		{
			name: "get string",
			fields: fields{
				id: uuid.New(),
				data: map[string]any{
					"nickname": "User",
					"email":    "awesome_user@mail.com",
				},
			},
			args: args{path: "nickname"},
			wantValue: func(tt require.TestingT, got any, _ ...any) {
				nickname, ok := got.(string)
				require.True(t, ok)
				require.Equal(t, "User", nickname)
			},
			wantError: require.NoError,
		},
		{
			name: "get map",
			fields: fields{
				id: uuid.New(),
				data: map[string]any{
					"nickname": "User",
					"email":    "awesome_user@mail.com",
					"address": map[string]any{
						"city":   "Moscow",
						"street": "Some Moscow street",
						"house":  "007",
					},
				},
			},
			args: args{path: "address"},
			wantValue: func(tt require.TestingT, got any, _ ...any) {
				address, ok := got.(map[string]any)
				require.True(t, ok)

				require.Contains(t, address, "city")
				require.Contains(t, address, "street")
				require.Contains(t, address, "house")
			},
			wantError: require.NoError,
		},
		{
			name: "empty path",
			fields: fields{
				id: uuid.New(),
				data: map[string]any{
					"nickname": "User",
					"email":    "awesome_user@mail.com",
					"address": map[string]any{
						"city":   "Moscow",
						"street": "Some Moscow street",
						"house":  "007",
					},
				},
			},
			args:      args{path: ""},
			wantValue: require.Empty,
			wantError: func(tt require.TestingT, err error, _ ...any) {
				require.EqualError(t, err, models.ErrIncorrectPath.Error())
			},
		},
		{
			name: "first level is empty",
			fields: fields{
				id: uuid.New(),
				data: map[string]any{
					"nickname": "User",
					"email":    "awesome_user@mail.com",
					"address": map[string]any{
						"city":   "Moscow",
						"street": "Some Moscow street",
						"house":  "007",
					},
				},
			},
			args:      args{path: ".city"},
			wantValue: require.Empty,
			wantError: func(tt require.TestingT, err error, _ ...any) {
				require.EqualError(t, err, models.ErrIncorrectPath.Error())
			},
		},
		{
			name: "level not found",
			fields: fields{
				id: uuid.New(),
				data: map[string]any{
					"nickname": "User",
					"email":    "awesome_user@mail.com",
					"address": map[string]any{
						"city":   "Moscow",
						"street": "Some Moscow street",
						"house":  "007",
					},
				},
			},
			args:      args{path: "age"},
			wantValue: require.Empty,
			wantError: func(tt require.TestingT, err error, _ ...any) {
				require.EqualError(t, err, models.ErrLevelNotFound.Error())
			},
		},
		{
			name: "intermediate is not a map",
			fields: fields{
				id: uuid.New(),
				data: map[string]any{
					"nickname": "User",
					"email":    "awesome_user@mail.com",
					"address": map[string]any{
						"city":   "Moscow",
						"street": "Some Moscow street",
						"house":  "007",
					},
				},
			},
			args:      args{path: "address.city.dummy"},
			wantValue: require.Empty,
			wantError: func(tt require.TestingT, err error, _ ...any) {
				require.EqualError(t, err, models.ErrIntermediateNotMap.Error())
			},
		},
	}

	for _, tc := range tests {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			d := models.Document{
				ID:   tc.fields.id,
				Data: tc.fields.data,
			}
			val, err := d.Get(tc.args.path)

			tc.wantValue(t, val)
			tc.wantError(t, err)
		})
	}
}

func TestDocument_Set(t *testing.T) {
	t.Parallel()

	type fields struct {
		ID   uuid.UUID
		Data map[string]any
	}

	type args struct {
		path  string
		value any
	}

	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr require.ErrorAssertionFunc
	}{
		{
			name: "successful set",
			fields: fields{
				ID: uuid.New(),
				Data: map[string]any{
					"nickname": "User",
					"email":    "awesome_user@mail.com",
				},
			},
			args: args{
				path:  "address.city",
				value: "Moscow",
			},
			wantErr: require.NoError,
		},
		{
			name: "empty path",
			fields: fields{
				ID: uuid.New(),
				Data: map[string]any{
					"nickname": "User",
					"email":    "awesome_user@mail.com",
				},
			},
			args: args{
				path:  "",
				value: "",
			},
			wantErr: func(tt require.TestingT, err error, i ...any) {
				require.EqualError(t, err, models.ErrIncorrectPath.Error())
			},
		},
		{
			name: "first level is empty",
			fields: fields{
				ID: uuid.New(),
				Data: map[string]any{
					"nickname": "User",
					"email":    "awesome_user@mail.com",
				},
			},
			args: args{
				path:  ".city",
				value: "",
			},
			wantErr: func(tt require.TestingT, err error, i ...any) {
				require.EqualError(t, err, models.ErrIncorrectPath.Error())
			},
		},
	}

	for _, tc := range tests {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			d := models.Document{
				ID:   tc.fields.ID,
				Data: tc.fields.Data,
			}
			err := d.Set(tc.args.path, tc.args.value)

			tc.wantErr(t, err)
		})
	}

}
