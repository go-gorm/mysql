package mysql

import (
	"testing"

	"gorm.io/gorm/schema"
)

func TestDataTypeOfGeneratedColumn(t *testing.T) {
	dialector := Dialector{Config: &Config{}}
	tests := []struct {
		name  string
		field *schema.Field
		want  string
	}{
		{
			name:  "computed column renders a STORED generated column",
			field: &schema.Field{DataType: "decimal(12,2)", TagSettings: map[string]string{"GENERATED": "price * quantity"}},
			want:  "decimal(12,2) GENERATED ALWAYS AS (price * quantity) STORED",
		},
		{
			name:  "computed expression keeps commas",
			field: &schema.Field{DataType: "varchar(64)", TagSettings: map[string]string{"GENERATED": "concat(first_name, last_name)"}},
			want:  "varchar(64) GENERATED ALWAYS AS (concat(first_name, last_name)) STORED",
		},
		{
			// `identity` is reserved for identity columns, which MySQL renders
			// through its native AUTO_INCREMENT rather than a computed column.
			name:  "identity keyword is not treated as a computed column",
			field: &schema.Field{DataType: schema.Int, Size: 64, AutoIncrement: true, TagSettings: map[string]string{"GENERATED": "identity"}},
			want:  "bigint AUTO_INCREMENT",
		},
		{
			name:  "identity with an explicit mode is also reserved",
			field: &schema.Field{DataType: schema.Int, Size: 64, AutoIncrement: true, TagSettings: map[string]string{"GENERATED": "identity always"}},
			want:  "bigint AUTO_INCREMENT",
		},
		{
			name:  "a bare generated tag is ignored",
			field: &schema.Field{DataType: "decimal(10,2)", TagSettings: map[string]string{"GENERATED": "GENERATED"}},
			want:  "decimal(10,2)",
		},
		{
			name:  "a lowercase generated expression is not mistaken for a bare tag",
			field: &schema.Field{DataType: "decimal(10,2)", TagSettings: map[string]string{"GENERATED": "generated"}},
			want:  "decimal(10,2) GENERATED ALWAYS AS (generated) STORED",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := dialector.DataTypeOf(tt.field); got != tt.want {
				t.Errorf("DataTypeOf() = %q, want %q", got, tt.want)
			}
		})
	}
}
