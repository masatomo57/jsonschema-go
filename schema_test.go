package jsonschema

import (
	"encoding/json"
	"testing"
)

// verifyFunc はカスタム検証ロジック用の関数型
type verifyFunc func(t *testing.T, schema map[string]any)

// testCase は単一のテストケースを表す
type testCase struct {
	name    string
	input   any
	wantErr bool
	verify  verifyFunc
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
			wantErr: false,
			verify: func(t *testing.T, schema map[string]any) {
				if schema["type"] != "object" {
					t.Errorf("Expected type object, got %v", schema["type"])
				}

				props, ok := schema["properties"].(map[string]any)
				if !ok {
					t.Fatalf("Expected properties map")
				}

				// Nameフィールドの確認
				nameSchema, ok := props["Name"].(map[string]any)
				if !ok || nameSchema["type"] != "string" {
					t.Errorf("Expected Name to be string type")
				}

				// Ageフィールドの確認
				ageSchema, ok := props["Age"].(map[string]any)
				if !ok || ageSchema["type"] != "integer" {
					t.Errorf("Expected Age to be integer type")
				}

				// Scoreフィールドの確認
				scoreSchema, ok := props["Score"].(map[string]any)
				if !ok || scoreSchema["type"] != "number" {
					t.Errorf("Expected Score to be number type")
				}

				// Activeフィールドの確認
				activeSchema, ok := props["Active"].(map[string]any)
				if !ok || activeSchema["type"] != "boolean" {
					t.Errorf("Expected Active to be boolean type")
				}
			},
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
			wantErr: false,
			verify: func(t *testing.T, schema map[string]any) {
				props, ok := schema["properties"].(map[string]any)
				if !ok {
					t.Fatalf("Expected properties map")
				}

				// カスタムJSON名の確認
				if _, ok := props["id"]; !ok {
					t.Errorf("Expected 'id' property")
				}
				if _, ok := props["name"]; !ok {
					t.Errorf("Expected 'name' property")
				}
				if _, ok := props["email"]; !ok {
					t.Errorf("Expected 'email' property")
				}

				// Passwordが除外されていることを確認
				if _, ok := props["Password"]; ok {
					t.Errorf("Password should be excluded")
				}
				if _, ok := props["password"]; ok {
					t.Errorf("password should be excluded")
				}
			},
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
			wantErr: false,
			verify: func(t *testing.T, schema map[string]any) {
				props, ok := schema["properties"].(map[string]any)
				if !ok {
					t.Fatalf("Expected properties map")
				}

				addressSchema, ok := props["Address"].(map[string]any)
				if !ok {
					t.Fatalf("Expected Address property")
				}

				if addressSchema["type"] != "object" {
					t.Errorf("Expected Address to be object type")
				}

				addressProps, ok := addressSchema["properties"].(map[string]any)
				if !ok {
					t.Fatalf("Expected Address properties")
				}

				if _, ok := addressProps["Street"]; !ok {
					t.Errorf("Expected Street property in Address")
				}
				if _, ok := addressProps["City"]; !ok {
					t.Errorf("Expected City property in Address")
				}
			},
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
			wantErr: false,
			verify: func(t *testing.T, schema map[string]any) {
				props, ok := schema["properties"].(map[string]any)
				if !ok {
					t.Fatalf("Expected properties map")
				}

				// Tags配列の確認
				tagsSchema, ok := props["Tags"].(map[string]any)
				if !ok || tagsSchema["type"] != "array" {
					t.Errorf("Expected Tags to be array type")
				}

				tagsItems, ok := tagsSchema["items"].(map[string]any)
				if !ok || tagsItems["type"] != "string" {
					t.Errorf("Expected Tags items to be string type")
				}

				// Numbers配列の確認
				numbersSchema, ok := props["Numbers"].(map[string]any)
				if !ok || numbersSchema["type"] != "array" {
					t.Errorf("Expected Numbers to be array type")
				}

				numbersItems, ok := numbersSchema["items"].(map[string]any)
				if !ok || numbersItems["type"] != "integer" {
					t.Errorf("Expected Numbers items to be integer type")
				}
			},
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
			wantErr: false,
			verify: func(t *testing.T, schema map[string]any) {
				props, ok := schema["properties"].(map[string]any)
				if !ok {
					t.Fatalf("Expected properties map")
				}

				// Metadataマップの確認
				metadataSchema, ok := props["Metadata"].(map[string]any)
				if !ok || metadataSchema["type"] != "object" {
					t.Errorf("Expected Metadata to be object type")
				}

				metadataProps, ok := metadataSchema["additionalProperties"].(map[string]any)
				if !ok || metadataProps["type"] != "string" {
					t.Errorf("Expected Metadata additionalProperties to be string type")
				}

				// Configマップの確認
				configSchema, ok := props["Config"].(map[string]any)
				if !ok || configSchema["type"] != "object" {
					t.Errorf("Expected Config to be object type")
				}

				configProps, ok := configSchema["additionalProperties"].(map[string]any)
				if !ok || configProps["type"] != "integer" {
					t.Errorf("Expected Config additionalProperties to be integer type")
				}
			},
		},
		{
			name: "正常系: バリデーション制約",
			input: func() any {
				type ValidatedStruct struct {
					ID       int    `validate:"required,minimum=1,maximum=100"`
					Name     string `validate:"required,minLength=1,maxLength=100"`
					Email    string `validate:"pattern=^[a-z]+@[a-z]+\\.[a-z]+$"`
					Age      int    `validate:"minimum=0,maximum=150"`
					Optional string `json:"optional,omitempty"`
				}
				return ValidatedStruct{}
			}(),
			wantErr: false,
			verify: func(t *testing.T, schema map[string]any) {
				props, ok := schema["properties"].(map[string]any)
				if !ok {
					t.Fatalf("Expected properties map")
				}

				// ID制約の確認
				idSchema, ok := props["ID"].(map[string]any)
				if !ok {
					t.Fatalf("Expected ID property")
				}
				if idSchema["minimum"] != float64(1) {
					t.Errorf("Expected ID minimum to be 1, got %v", idSchema["minimum"])
				}
				if idSchema["maximum"] != float64(100) {
					t.Errorf("Expected ID maximum to be 100, got %v", idSchema["maximum"])
				}

				// Name制約の確認
				nameSchema, ok := props["Name"].(map[string]any)
				if !ok {
					t.Fatalf("Expected Name property")
				}
				if nameSchema["minLength"] != 1 {
					t.Errorf("Expected Name minLength to be 1, got %v", nameSchema["minLength"])
				}
				if nameSchema["maxLength"] != 100 {
					t.Errorf("Expected Name maxLength to be 100, got %v", nameSchema["maxLength"])
				}

				// Emailパターンの確認
				emailSchema, ok := props["Email"].(map[string]any)
				if !ok {
					t.Fatalf("Expected Email property")
				}
				expectedPattern := "^[a-z]+@[a-z]+\\.[a-z]+$"
				if emailSchema["pattern"] != expectedPattern {
					t.Errorf("Expected Email pattern to be %s, got %v", expectedPattern, emailSchema["pattern"])
				}

				// requiredフィールドの確認
				required, ok := schema["required"].([]string)
				if !ok {
					t.Fatalf("Expected required array")
				}

				requiredMap := make(map[string]bool)
				for _, r := range required {
					requiredMap[r] = true
				}

				if !requiredMap["ID"] {
					t.Errorf("Expected ID to be required")
				}
				if !requiredMap["Name"] {
					t.Errorf("Expected Name to be required")
				}
			},
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
					Name     string            `json:"name" validate:"required,minLength=1,maxLength=100"`
					Email    string            `json:"email" validate:"pattern=^[a-z]+@[a-z]+\\.[a-z]+$"`
					Age      int               `json:"age,omitempty" validate:"minimum=0,maximum=150"`
					Tags     []string          `json:"tags"`
					Metadata map[string]string `json:"metadata"`
					Address  Address           `json:"address"`
				}
				return User{}
			}(),
			wantErr: false,
			verify: func(t *testing.T, schema map[string]any) {
				props, ok := schema["properties"].(map[string]any)
				if !ok {
					t.Fatalf("Expected properties map")
				}

				// すべての期待されるプロパティが存在することを確認
				expectedProps := []string{"id", "name", "email", "age", "tags", "metadata", "address"}
				for _, prop := range expectedProps {
					if _, ok := props[prop]; !ok {
						t.Errorf("Expected property %s not found", prop)
					}
				}

				// ネストしたaddress構造体の確認
				addressSchema, ok := props["address"].(map[string]any)
				if !ok {
					t.Fatalf("Expected address property")
				}

				addressProps, ok := addressSchema["properties"].(map[string]any)
				if !ok {
					t.Fatalf("Expected address properties")
				}

				if _, ok := addressProps["street"]; !ok {
					t.Errorf("Expected street property in address")
				}
				if _, ok := addressProps["city"]; !ok {
					t.Errorf("Expected city property in address")
				}

				// デバッグ用にスキーマを出力
				jsonBytes, _ := json.MarshalIndent(schema, "", "  ")
				t.Logf("Generated schema:\n%s", string(jsonBytes))
			},
		},
		// エラーケース
		{
			name:    "異常系: nil値",
			input:   nil,
			wantErr: true,
			verify:  nil,
		},
		{
			name:    "異常系: 非構造体型",
			input:   "not a struct",
			wantErr: true,
			verify:  nil,
		},
		{
			name: "正常系: 構造体へのポインタ",
			input: func() any {
				type TestStruct struct {
					Name string
				}
				return &TestStruct{}
			}(),
			wantErr: false,
			verify:  nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			schema, err := Generate(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("Generate() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if err == nil && tt.verify != nil {
				tt.verify(t, schema)
			}
		})
	}
}
