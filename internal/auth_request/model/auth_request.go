package model

import (
	"time"

	"golang.org/x/text/language"

	"github.com/caos/zitadel/internal/iam/model"
)

type AuthRequest struct {
	ID            string
	AgentID       string
	CreationDate  time.Time
	ChangeDate    time.Time
	BrowserInfo   *BrowserInfo
	ApplicationID string
	CallbackURI   string
	TransferState string
	Prompt        Prompt
	PossibleLOAs  []LevelOfAssurance
	UiLocales     []string
	LoginHint     string
	MaxAuthAge    uint32
	Request       Request

	levelOfAssurance       LevelOfAssurance
	UserID                 string
	UserName               string
	LoginName              string
	DisplayName            string
	UserOrgID              string
	RequestedOrgID         string
	RequestedOrgName       string
	RequestedPrimaryDomain string
	SelectedIDPConfigID    string
	LinkingUsers           []*ExternalUser
	PossibleSteps          []NextStep
	PasswordVerified       bool
	MFAsVerified           []MFAType
	Audience               []string
	AuthTime               time.Time
	Code                   string
	LoginPolicy            *model.LoginPolicyView
	LabelPolicy            *model.LabelPolicyView
	AllowedExternalIDPs    []*model.IDPProviderView
}

type ExternalUser struct {
	IDPConfigID       string
	ExternalUserID    string
	DisplayName       string
	PreferredUsername string
	FirstName         string
	LastName          string
	NickName          string
	Email             string
	IsEmailVerified   bool
	PreferredLanguage language.Tag
	Phone             string
	IsPhoneVerified   bool
}

type Prompt int32

const (
	PromptUnspecified Prompt = iota
	PromptNone
	PromptLogin
	PromptConsent
	PromptSelectAccount
)

type LevelOfAssurance int

const (
	LevelOfAssuranceNone LevelOfAssurance = iota
)

func (a *AuthRequest) IsValid() bool {
	return a.ID != "" &&
		a.AgentID != "" &&
		a.BrowserInfo != nil && a.BrowserInfo.IsValid() &&
		a.ApplicationID != "" &&
		a.CallbackURI != "" &&
		a.Request != nil && a.Request.IsValid()
}
