package domain

import (
	"github.com/caos/zitadel/internal/eventstore/v1/models"
)

type Project struct {
	models.ObjectRoot

	State                          ProjectState
	Name                           string
	ProjectRoleAssertion           bool
	ProjectRoleCheck               bool
	HasProjectCheck                bool
	PrivateLabelingSetting         PrivateLabelingSetting
	LoginPolicySetting             LoginPolicySetting
	RegisterOnProjectResourceOwner bool
}

type ProjectState int32

const (
	ProjectStateUnspecified ProjectState = iota
	ProjectStateActive
	ProjectStateInactive
	ProjectStateRemoved
)

type PrivateLabelingSetting int32

const (
	PrivateLabelingSettingUnspecified PrivateLabelingSetting = iota
	PrivateLabelingSettingEnforceProjectResourceOwnerPolicy
	PrivateLabelingSettingAllowLoginUserResourceOwnerPolicy
)

type LoginPolicySetting int32

const (
	LoginPolicySettingUnspecified LoginPolicySetting = iota
	LoginPolicySettingEnforceProjectResourceOwnerPolicy
	LoginPolicySettingAllowLoginUserResourceOwnerPolicy
)

func (o *Project) IsValid() bool {
	return o.Name != ""
}
