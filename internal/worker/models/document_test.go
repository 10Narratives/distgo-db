package models_test

import (
	"testing"

	"githib.com/10Narratives/distgo-db/internal/worker/models"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestDocument_DeepCopy(t *testing.T) {
	t.Parallel()

	type fields struct {
		ID   uuid.UUID
		Data map[string]any
	}

	tests := []struct {
		name      string
		fields    fields
		wantValue func(t *testing.T, original, copied *models.Document)
	}{
		{
			name: "simple data",
			fields: fields{
				ID: uuid.New(),
				Data: map[string]any{
					"name": "test",
					"age":  25,
				},
			},
			wantValue: func(t *testing.T, original, copied *models.Document) {
				assert.Equal(t, original.ID, copied.ID)
				assert.Equal(t, original.Data, copied.Data)

				copied.Data["modified"] = true
				assert.NotEqual(t, original.Data, copied.Data)

				original.Data["nested"] = map[string]any{"key": "value"}
				assert.NotContains(t, copied.Data, "nested")
			},
		},
		{
			name: "nested structures",
			fields: fields{
				ID: uuid.New(),
				Data: map[string]any{
					"slice": []any{1, "two", false},
					"map": map[string]any{
						"inner": map[string]any{
							"value": 42,
						},
					},
				},
			},
			wantValue: func(t *testing.T, original, copied *models.Document) {
				origSlice := original.Data["slice"].([]any)
				copiedSlice := copied.Data["slice"].([]any)
				assert.Equal(t, origSlice, copiedSlice)

				copied.Data["slice"].([]any)[0] = "modified"
				assert.NotEqual(t, original.Data["slice"], copied.Data["slice"])
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			original := &models.Document{
				ID:   tt.fields.ID,
				Data: tt.fields.Data,
			}

			copied := original.DeepCopy()

			tt.wantValue(t, original, copied)
		})
	}
}
