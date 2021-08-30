package model

import (
	"net"
	"testing"
)

func TestAuthRequest_IsValid(t *testing.T) {
	type fields struct {
		ID            string
		AgentID       string
		BrowserInfo   *BrowserInfo
		ApplicationID string
		CallbackURI   string
		Request       Request
	}
	tests := []struct {
		name   string
		fields fields
		want   bool
	}{
		{
			"missing id, false",
			fields{},
			false,
		},
		{
			"missing agent id, false",
			fields{
				ID: "id",
			},
			false,
		},
		{
			"missing browser info, false",
			fields{
				ID:      "id",
				AgentID: "agentID",
			},
			false,
		},
		{
			"browser info invalid, false",
			fields{
				ID:          "id",
				AgentID:     "agentID",
				BrowserInfo: &BrowserInfo{},
			},
			false,
		},
		{
			"missing application id, false",
			fields{
				ID:      "id",
				AgentID: "agentID",
				BrowserInfo: &BrowserInfo{
					UserAgent:      "user agent",
					AcceptLanguage: "accept language",
					RemoteIP:       net.IPv4(29, 4, 20, 19),
				},
			},
			false,
		},
		{
			"missing callback uri, false",
			fields{
				ID:      "id",
				AgentID: "agentID",
				BrowserInfo: &BrowserInfo{
					UserAgent:      "user agent",
					AcceptLanguage: "accept language",
					RemoteIP:       net.IPv4(29, 4, 20, 19),
				},
				ApplicationID: "appID",
			},
			false,
		},
		{
			"missing request, false",
			fields{
				ID:      "id",
				AgentID: "agentID",
				BrowserInfo: &BrowserInfo{
					UserAgent:      "user agent",
					AcceptLanguage: "accept language",
					RemoteIP:       net.IPv4(29, 4, 20, 19),
				},
				ApplicationID: "appID",
				CallbackURI:   "schema://callback",
			},
			false,
		},
		{
			"request invalid, false",
			fields{
				ID:      "id",
				AgentID: "agentID",
				BrowserInfo: &BrowserInfo{
					UserAgent:      "user agent",
					AcceptLanguage: "accept language",
					RemoteIP:       net.IPv4(29, 4, 20, 19),
				},
				ApplicationID: "appID",
				CallbackURI:   "schema://callback",
				Request:       &AuthRequestOIDC{},
			},
			false,
		},
		{
			"valid auth request, true",
			fields{
				ID:      "id",
				AgentID: "agentID",
				BrowserInfo: &BrowserInfo{
					UserAgent:      "user agent",
					AcceptLanguage: "accept language",
					RemoteIP:       net.IPv4(29, 4, 20, 19),
				},
				ApplicationID: "appID",
				CallbackURI:   "schema://callback",
				Request: &AuthRequestOIDC{
					Scopes: []string{"openid"},
					CodeChallenge: &OIDCCodeChallenge{
						Challenge: "challenge",
						Method:    CodeChallengeMethodS256,
					},
				},
			},
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := &AuthRequest{
				ID:            tt.fields.ID,
				AgentID:       tt.fields.AgentID,
				BrowserInfo:   tt.fields.BrowserInfo,
				ApplicationID: tt.fields.ApplicationID,
				CallbackURI:   tt.fields.CallbackURI,
				Request:       tt.fields.Request,
			}
			if got := a.IsValid(); got != tt.want {
				t.Errorf("IsValid() = %v, want %v", got, tt.want)
			}
		})
	}
}
