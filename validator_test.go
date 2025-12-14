package samhook

import "testing"

func TestValidateWebhookURL(t *testing.T) {
	tests := []struct {
		name    string
		url     string
		wantErr bool
	}{
		{
			name:    "有效的 HTTPS URL",
			url:     "https://example.com/webhook",
			wantErr: false,
		},
		{
			name:    "有效的 HTTP URL",
			url:     "http://example.com/webhook",
			wantErr: false,
		},
		{
			name:    "空 URL",
			url:     "",
			wantErr: true,
		},
		{
			name:    "無效的 URL 格式",
			url:     "not-a-valid-url",
			wantErr: true,
		},
		{
			name:    "無協議的 URL",
			url:     "example.com/webhook",
			wantErr: true,
		},
		{
			name:    "FTP 協議（不允許）",
			url:     "ftp://example.com/webhook",
			wantErr: true,
		},
		{
			name:    "無主機的 URL",
			url:     "https://",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateWebhookURL(tt.url)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateWebhookURL() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
