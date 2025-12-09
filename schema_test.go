package jsonschema

import (
	"encoding/json"
	"reflect"
	"testing"
)

// testCase は単一のテストケースを表す
type testCase struct {
	name    string
	input   any
	want    map[string]any
	wantErr bool
}

func TestGenerate(t *testing.T) {
	tests := []testCase{
		{
			name: "正常系: 基本型",
			input: func() any {
				type BasicStruct struct {
					Name   string
					Age    int
					Score  float64
					Active bool
				}
				return BasicStruct{}
			}(),
			want: map[string]any{
				"type": "object",
				"properties": map[string]any{
					"Name": map[string]any{
						"type": "string",
					},
					"Age": map[string]any{
						"type": "integer",
					},
					"Score": map[string]any{
						"type": "number",
					},
					"Active": map[string]any{
						"type": "boolean",
					},
				},
				"required":             []string{"Name", "Age", "Score", "Active"},
				"additionalProperties": false,
			},
			wantErr: false,
		},
		{
			name: "正常系: JSONタグ",
			input: func() any {
				type TaggedStruct struct {
					ID       int    `json:"id"`
					Name     string `json:"name"`
					Email    string `json:"email,omitempty"`
					Password string `json:"-"`
				}
				return TaggedStruct{}
			}(),
			want: map[string]any{
				"type": "object",
				"properties": map[string]any{
					"id": map[string]any{
						"type": "integer",
					},
					"name": map[string]any{
						"type": "string",
					},
					"email": map[string]any{
						"type": "string",
					},
				},
				"required":             []string{"id", "name"},
				"additionalProperties": false,
			},
			wantErr: false,
		},
		{
			name: "正常系: ネストした構造体",
			input: func() any {
				type Address struct {
					Street string
					City   string
				}

				type User struct {
					Name    string
					Address Address
				}
				return User{}
			}(),
			want: map[string]any{
				"type": "object",
				"properties": map[string]any{
					"Name": map[string]any{
						"type": "string",
					},
					"Address": map[string]any{
						"type": "object",
						"properties": map[string]any{
							"Street": map[string]any{
								"type": "string",
							},
							"City": map[string]any{
								"type": "string",
							},
						},
						"required":             []string{"Street", "City"},
						"additionalProperties": false,
					},
				},
				"required":             []string{"Name", "Address"},
				"additionalProperties": false,
			},
			wantErr: false,
		},
		{
			name: "正常系: 配列",
			input: func() any {
				type ArrayStruct struct {
					Tags    []string
					Numbers []int
					Scores  []float64
				}
				return ArrayStruct{}
			}(),
			want: map[string]any{
				"type": "object",
				"properties": map[string]any{
					"Tags": map[string]any{
						"type": "array",
						"items": map[string]any{
							"type": "string",
						},
					},
					"Numbers": map[string]any{
						"type": "array",
						"items": map[string]any{
							"type": "integer",
						},
					},
					"Scores": map[string]any{
						"type": "array",
						"items": map[string]any{
							"type": "number",
						},
					},
				},
				"required":             []string{"Tags", "Numbers", "Scores"},
				"additionalProperties": false,
			},
			wantErr: false,
		},
		{
			name: "正常系: マップ",
			input: func() any {
				type MapStruct struct {
					Metadata map[string]string
					Config   map[string]int
				}
				return MapStruct{}
			}(),
			want: map[string]any{
				"type": "object",
				"properties": map[string]any{
					"Metadata": map[string]any{
						"type": "object",
						"additionalProperties": map[string]any{
							"type": "string",
						},
					},
					"Config": map[string]any{
						"type": "object",
						"additionalProperties": map[string]any{
							"type": "integer",
						},
					},
				},
				"required":             []string{"Metadata", "Config"},
				"additionalProperties": false,
			},
			wantErr: false,
		},
		{
			name: "正常系: バリデーション制約",
			input: func() any {
				type ValidatedStruct struct {
					ID       int      `validate:"required,minimum=1,maximum=100"`
					Name     string   `validate:"required,pattern=^[a-zA-Z]+$"`
					Email    string   `validate:"pattern=^[a-z]+@[a-z]+\\.[a-z]+$,format=email"`
					Age      int      `validate:"minimum=0,maximum=150"`
					Score    float64  `validate:"multipleOf=0.5,exclusiveMinimum=0,exclusiveMaximum=100"`
					Tags     []string `validate:"minItems=1,maxItems=10"`
					Optional string   `json:"optional,omitempty"`
				}
				return ValidatedStruct{}
			}(),
			want: map[string]any{
				"type": "object",
				"properties": map[string]any{
					"ID": map[string]any{
						"type":    "integer",
						"minimum": float64(1),
						"maximum": float64(100),
					},
					"Name": map[string]any{
						"type":    "string",
						"pattern": "^[a-zA-Z]+$",
					},
					"Email": map[string]any{
						"type":    "string",
						"pattern": "^[a-z]+@[a-z]+\\.[a-z]+$",
						"format":  "email",
					},
					"Age": map[string]any{
						"type":    "integer",
						"minimum": float64(0),
						"maximum": float64(150),
					},
					"Score": map[string]any{
						"type":             "number",
						"multipleOf":       float64(0.5),
						"exclusiveMinimum": float64(0),
						"exclusiveMaximum": float64(100),
					},
					"Tags": map[string]any{
						"type": "array",
						"items": map[string]any{
							"type": "string",
						},
						"minItems": 1,
						"maxItems": 10,
					},
					"optional": map[string]any{
						"type": "string",
					},
				},
				"required":             []string{"ID", "Name", "Email", "Age", "Score", "Tags"},
				"additionalProperties": false,
			},
			wantErr: false,
		},
		{
			name: "正常系: 複雑な例",
			input: func() any {
				type Address struct {
					Street string `json:"street" validate:"required"`
					City   string `json:"city" validate:"required"`
				}

				type User struct {
					ID       int               `json:"id" validate:"required"`
					Name     string            `json:"name" validate:"required,pattern=^[a-zA-Z ]+$"`
					Email    string            `json:"email" validate:"pattern=^[a-z]+@[a-z]+\\.[a-z]+$,format=email"`
					Age      int               `json:"age,omitempty" validate:"minimum=0,maximum=150"`
					Tags     []string          `json:"tags" validate:"minItems=0,maxItems=20"`
					Metadata map[string]string `json:"metadata"`
					Address  Address           `json:"address"`
				}
				return User{}
			}(),
			want: map[string]any{
				"type": "object",
				"properties": map[string]any{
					"id": map[string]any{
						"type": "integer",
					},
					"name": map[string]any{
						"type":    "string",
						"pattern": "^[a-zA-Z ]+$",
					},
					"email": map[string]any{
						"type":    "string",
						"pattern": "^[a-z]+@[a-z]+\\.[a-z]+$",
						"format":  "email",
					},
					"age": map[string]any{
						"type":    "integer",
						"minimum": float64(0),
						"maximum": float64(150),
					},
					"tags": map[string]any{
						"type": "array",
						"items": map[string]any{
							"type": "string",
						},
						"minItems": 0,
						"maxItems": 20,
					},
					"metadata": map[string]any{
						"type": "object",
						"additionalProperties": map[string]any{
							"type": "string",
						},
					},
					"address": map[string]any{
						"type": "object",
						"properties": map[string]any{
							"street": map[string]any{
								"type": "string",
							},
							"city": map[string]any{
								"type": "string",
							},
						},
						"required":             []string{"street", "city"},
						"additionalProperties": false,
					},
				},
				"required":             []string{"id", "name", "email", "tags", "metadata", "address"},
				"additionalProperties": false,
			},
			wantErr: false,
		},
		// エラーケース
		{
			name:    "異常系: nil値",
			input:   nil,
			want:    nil,
			wantErr: true,
		},
		{
			name:    "異常系: 非構造体型",
			input:   "not a struct",
			want:    nil,
			wantErr: true,
		},
		{
			name: "正常系: 構造体へのポインタ",
			input: func() any {
				type TestStruct struct {
					Name string
				}
				return &TestStruct{}
			}(),
			want: map[string]any{
				"type": "object",
				"properties": map[string]any{
					"Name": map[string]any{
						"type": "string",
					},
				},
				"required":             []string{"Name"},
				"additionalProperties": false,
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			schema, err := Generate(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("Generate() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if err == nil && tt.want != nil {
				if !reflect.DeepEqual(schema, tt.want) {
					gotJSON, _ := json.MarshalIndent(schema, "", "  ")
					wantJSON, _ := json.MarshalIndent(tt.want, "", "  ")
					t.Errorf("Generate() = \n%s\nwant \n%s", string(gotJSON), string(wantJSON))
				}
			}
		})
	}
}
