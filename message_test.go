package samhook

import (
	"reflect"
	"strings"
	"testing"

	"github.com/bytedance/sonic"
)

func TestMessage_JSONSerialization(t *testing.T) {
	tests := []struct {
		name    string
		message Message
		want    string
	}{
		{
			name: "基本訊息",
			message: Message{
				Text:     "Hello",
				Username: "bot",
			},
			want: `{"text":"Hello","username":"bot"}`,
		},
		{
			name: "完整訊息",
			message: Message{
				Text:      "Hello",
				Username:  "bot",
				IconEmoji: ":robot_face:",
				Channel:   "#general",
			},
			want: `{"text":"Hello","username":"bot","icon_emoji":":robot_face:","channel":"#general"}`,
		},
		{
			name:    "空訊息",
			message: Message{},
			want:    `{}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := sonic.Marshal(tt.message)
			if err != nil {
				t.Fatalf("sonic.Marshal() error = %v", err)
			}

			// 解析 JSON 進行比較（忽略順序）
			var gotMap, wantMap map[string]interface{}
			if err := sonic.Unmarshal(got, &gotMap); err != nil {
				t.Fatalf("failed to unmarshal got: %v", err)
			}
			if err := sonic.Unmarshal([]byte(tt.want), &wantMap); err != nil {
				t.Fatalf("failed to unmarshal want: %v", err)
			}

			if !reflect.DeepEqual(gotMap, wantMap) {
				t.Errorf("sonic.Marshal() = %s, want %s", got, tt.want)
			}
		})
	}
}

func TestMessage_OmitEmpty(t *testing.T) {
	msg := Message{
		Text: "Hello",
		// 其他欄位為空
	}

	data, err := sonic.Marshal(msg)
	if err != nil {
		t.Fatalf("sonic.Marshal() error = %v", err)
	}

	jsonStr := string(data)

	// 驗證空欄位被省略
	if strings.Contains(jsonStr, "username") {
		t.Error("empty username should be omitted")
	}
	if strings.Contains(jsonStr, "icon_url") {
		t.Error("empty icon_url should be omitted")
	}
	if strings.Contains(jsonStr, "attachments") {
		t.Error("empty attachments should be omitted")
	}

	// 驗證非空欄位被包含
	if !strings.Contains(jsonStr, "text") {
		t.Error("non-empty text should be included")
	}
}

func TestMessage_JSONTags(t *testing.T) {
	msg := Message{
		Text:      "test",
		Username:  "bot",
		IconURL:   "https://example.com/icon.png",
		IconEmoji: ":emoji:",
		Channel:   "#channel",
	}

	data, err := sonic.Marshal(msg)
	if err != nil {
		t.Fatalf("sonic.Marshal() error = %v", err)
	}

	jsonStr := string(data)

	// 驗證 JSON 標籤與 Slack API 規範一致
	expectedTags := []string{
		`"text"`,
		`"username"`,
		`"icon_url"`,   // 不是 iconUrl
		`"icon_emoji"`, // 不是 iconEmoji
		`"channel"`,
	}

	for _, tag := range expectedTags {
		if !strings.Contains(jsonStr, tag) {
			t.Errorf("expected JSON tag %s not found in %s", tag, jsonStr)
		}
	}
}

func TestMessage_RoundTrip(t *testing.T) {
	original := Message{
		Text:      "Test message",
		Username:  "test-bot",
		IconEmoji: ":robot_face:",
		Channel:   "#test",
		Attachments: []Attachment{
			{
				Color: Good,
				Title: "Test",
			},
		},
	}

	// 序列化
	data, err := sonic.Marshal(original)
	if err != nil {
		t.Fatalf("sonic.Marshal() error = %v", err)
	}

	// 反序列化
	var unmarshaled Message
	if err := sonic.Unmarshal(data, &unmarshaled); err != nil {
		t.Fatalf("sonic.Unmarshal() error = %v", err)
	}

	// 比較
	if unmarshaled.Text != original.Text {
		t.Errorf("Text: got %q, want %q", unmarshaled.Text, original.Text)
	}
	if unmarshaled.Username != original.Username {
		t.Errorf("Username: got %q, want %q", unmarshaled.Username, original.Username)
	}
	if len(unmarshaled.Attachments) != len(original.Attachments) {
		t.Errorf("Attachments length: got %d, want %d",
			len(unmarshaled.Attachments), len(original.Attachments))
	}
}

func TestAttachment_JSONSerialization(t *testing.T) {
	attachment := Attachment{
		Color:      Good,
		Title:      "Test Title",
		TitleLink:  "https://example.com",
		Text:       "Test text",
		Fallback:   "Fallback text",
		AuthorName: "Test Author",
		Fields: []Field{
			{Title: "Field 1", Value: "Value 1", Short: true},
			{Title: "Field 2", Value: "Value 2", Short: false},
		},
	}

	data, err := sonic.Marshal(attachment)
	if err != nil {
		t.Fatalf("sonic.Marshal() error = %v", err)
	}

	// 驗證關鍵欄位存在
	jsonStr := string(data)
	expectedFields := []string{
		`"color"`,
		`"title"`,
		`"text"`,
		`"fallback"`,
		`"fields"`,
	}

	for _, field := range expectedFields {
		if !strings.Contains(jsonStr, field) {
			t.Errorf("expected field %s not found", field)
		}
	}

	// 驗證 JSON 標籤正確性（使用下劃線）
	if !strings.Contains(jsonStr, `"author_name"`) {
		t.Error("expected author_name tag (with underscore)")
	}
	if !strings.Contains(jsonStr, `"title_link"`) {
		t.Error("expected title_link tag (with underscore)")
	}
}

func TestAttachment_OmitEmpty(t *testing.T) {
	attachment := Attachment{
		Color: Good,
		Title: "Test",
		// 其他欄位為空
	}

	data, err := sonic.Marshal(attachment)
	if err != nil {
		t.Fatalf("sonic.Marshal() error = %v", err)
	}

	jsonStr := string(data)

	// 空欄位應該被省略
	if strings.Contains(jsonStr, "pretext") {
		t.Error("empty pretext should be omitted")
	}
	if strings.Contains(jsonStr, "fields") {
		t.Error("empty fields should be omitted")
	}
}

func TestField_JSONSerialization(t *testing.T) {
	tests := []struct {
		name  string
		field Field
		want  string
	}{
		{
			name: "短欄位",
			field: Field{
				Title: "Title",
				Value: "Value",
				Short: true,
			},
			want: `{"title":"Title","value":"Value","short":true}`,
		},
		{
			name: "長欄位",
			field: Field{
				Title: "Title",
				Value: "Value",
				Short: false,
			},
			// 注意：false 值在 omitempty 時會被省略，所以期望結果不包含 short
			want: `{"title":"Title","value":"Value"}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := sonic.Marshal(tt.field)
			if err != nil {
				t.Fatalf("sonic.Marshal() error = %v", err)
			}

			var gotMap, wantMap map[string]interface{}
			sonic.Unmarshal(got, &gotMap)
			sonic.Unmarshal([]byte(tt.want), &wantMap)

			if !reflect.DeepEqual(gotMap, wantMap) {
				t.Errorf("sonic.Marshal() = %s, want %s", got, tt.want)
			}
		})
	}
}

func TestField_BooleanSerialization(t *testing.T) {
	field := Field{
		Title: "Test",
		Value: "Value",
		Short: true,
	}

	data, err := sonic.Marshal(field)
	if err != nil {
		t.Fatalf("sonic.Marshal() error = %v", err)
	}

	// 驗證布林值正確序列化為 true/false（不是 1/0）
	jsonStr := string(data)
	if !strings.Contains(jsonStr, `"short":true`) {
		t.Error("boolean should be serialized as true, not 1")
	}
}

func TestMessage_WithAttachments(t *testing.T) {
	msg := Message{
		Text: "Test",
		Attachments: []Attachment{
			{
				Color: Good,
				Title: "Attachment 1",
				Fields: []Field{
					{Title: "Field 1", Value: "Value 1", Short: true},
				},
			},
			{
				Color: Warning,
				Title: "Attachment 2",
			},
		},
	}

	data, err := sonic.Marshal(msg)
	if err != nil {
		t.Fatalf("sonic.Marshal() error = %v", err)
	}

	// 驗證嵌套結構正確序列化
	jsonStr := string(data)
	if !strings.Contains(jsonStr, `"attachments"`) {
		t.Error("attachments should be included")
	}
	if !strings.Contains(jsonStr, `"fields"`) {
		t.Error("nested fields should be included")
	}

	// 驗證可以反序列化
	var unmarshaled Message
	if err := sonic.Unmarshal(data, &unmarshaled); err != nil {
		t.Fatalf("sonic.Unmarshal() error = %v", err)
	}

	if len(unmarshaled.Attachments) != 2 {
		t.Errorf("expected 2 attachments, got %d", len(unmarshaled.Attachments))
	}
}

func TestColorConstants(t *testing.T) {
	if Warning != "#FFBB00" {
		t.Errorf("Warning color incorrect: expected #FFBB00, got %s", Warning)
	}
	if Danger != "#FF0000" {
		t.Errorf("Danger color incorrect: expected #FF0000, got %s", Danger)
	}
	if Good != "#00FF00" {
		t.Errorf("Good color incorrect: expected #00FF00, got %s", Good)
	}
}
