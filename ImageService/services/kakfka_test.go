package services

import "testing"

func TestMsgUnarshal(t *testing.T) {
	tests := []struct {
		name          string
		kMsg          []byte
		wantUserID    int
		wantProductID int
	}{
		{
			name:          "Pure JSON",
			kMsg:          []byte(`{"u":1,"p":1}`),
			wantUserID:    1,
			wantProductID: 1,
		},
		{
			name:          "Not a JSON",
			kMsg:          []byte(`{"u":2"p":2}`),
			wantUserID:    -1,
			wantProductID: -1,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotUserID, gotProductID := MsgUnarshal(tt.kMsg)
			if gotUserID != tt.wantUserID {
				t.Errorf("MsgUnarshal() gotUserID = %v, want %v", gotUserID, tt.wantUserID)
			}
			if gotProductID != tt.wantProductID {
				t.Errorf("MsgUnarshal() gotProductID = %v, want %v", gotProductID, tt.wantProductID)
			}
		})
	}
}
