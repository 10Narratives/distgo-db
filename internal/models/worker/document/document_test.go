package documentmodels_test

import (
	"testing"
	"time"

	documentmodels "github.com/10Narratives/distgo-db/internal/models/worker/document"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDocument_Copy(t *testing.T) {
	t.Parallel()

	var (
		documentID uuid.UUID      = uuid.New()
		content    map[string]any = map[string]any{
			"string_field": "some_string_value",
			"int_field":    42,
			"float_field":  3.14159,
			"bool_field":   true,
			"nil_field":    nil,
			"slice_field":  []interface{}{"a", 1, false, nil},
			"struct_field": map[string]interface{}{
				"nested_key": "nested_value",
				"nested_slice": []interface{}{
					map[string]interface{}{"id": 1, "active": true},
					map[string]interface{}{"id": 2, "active": false},
				},
			},
			"deep_nested": map[string]interface{}{
				"level1": map[string]interface{}{
					"level2": map[string]interface{}{
						"level3": "deep_value",
					},
				},
			},
		}
		createdAt time.Time = time.Now()
		updatedAt time.Time = time.Now()
	)

	type fields struct {
		ID        uuid.UUID
		Content   map[string]any
		CreatedAt time.Time
		UpdatedAt time.Time
	}
	tests := []struct {
		name    string
		fields  fields
		wantVal require.ValueAssertionFunc
	}{
		{
			name: "successful copy",
			fields: fields{
				ID:        documentID,
				Content:   content,
				CreatedAt: createdAt,
				UpdatedAt: updatedAt,
			},
			wantVal: func(tt require.TestingT, got interface{}, i2 ...interface{}) {
				document, ok := got.(documentmodels.Document)
				require.True(t, ok)

				assert.Equal(t, documentID, document.ID)
				assert.Equal(t, content, document.Content)
				assert.Equal(t, createdAt, document.CreatedAt)
				assert.Equal(t, updatedAt, document.UpdatedAt)
			},
		},
	}
	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			src := documentmodels.Document{
				ID:        documentID,
				Content:   content,
				CreatedAt: createdAt,
				UpdatedAt: updatedAt,
			}

			copy := src.Copy()

			tt.wantVal(t, copy)
		})
	}
}
