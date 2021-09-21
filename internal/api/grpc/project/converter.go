package project

import (
	"github.com/caos/zitadel/internal/api/grpc/object"
	"github.com/caos/zitadel/internal/domain"
	"github.com/caos/zitadel/internal/errors"
	proj_model "github.com/caos/zitadel/internal/project/model"
	"github.com/caos/zitadel/internal/query"
	proj_pb "github.com/caos/zitadel/pkg/grpc/project"
)

func ProjectViewsToPb(projects []*query.Project) []*proj_pb.Project {
	o := make([]*proj_pb.Project, len(projects))
	for i, org := range projects {
		o[i] = ProjectViewToPb(org)
	}
	return o
}

func ProjectViewToPb(project *query.Project) *proj_pb.Project {
	return &proj_pb.Project{
		Id:    project.ID,
		State: projectStateToPb(project.State),
		Name:  project.Name,
		Details: object.ToViewDetailsPb(
			project.Sequence,
			project.CreationDate,
			project.ChangeDate,
			project.ResourceOwner,
		),
	}
}

func GrantedProjectViewsToPb(projects []*query.ProjectGrant) []*proj_pb.GrantedProject {
	p := make([]*proj_pb.GrantedProject, len(projects))
	for i, project := range projects {
		p[i] = GrantedProjectViewToPb(project)
	}
	return p
}

func GrantedProjectViewToPb(project *query.ProjectGrant) *proj_pb.GrantedProject {
	return &proj_pb.GrantedProject{
		ProjectId:        project.ProjectID,
		GrantId:          project.GrantID,
		Details:          object.ToViewDetailsPb(project.Sequence, project.CreationDate, project.ChangeDate, project.ResourceOwner),
		ProjectName:      project.ProjectName,
		State:            projectGrantStateToPb(project.State),
		ProjectOwnerId:   project.ResourceOwner,
		ProjectOwnerName: project.ResourceOwnerName,
		GrantedOrgId:     project.GrantedOrgID,
		GrantedOrgName:   project.OrgName,
		GrantedRoleKeys:  project.GrantedRoleKeys,
	}
}
func ProjectQueriesToModel(queries []*proj_pb.ProjectQuery) (_ []query.SearchQuery, err error) {
	q := make([]query.SearchQuery, len(queries))
	for i, query := range queries {
		q[i], err = ProjectQueryToModel(query)
		if err != nil {
			return nil, err
		}
	}
	return q, nil
}

func ProjectQueryToModel(apiQuery *proj_pb.ProjectQuery) (query.SearchQuery, error) {
	switch q := apiQuery.Query.(type) {
	case *proj_pb.ProjectQuery_NameQuery:
		return query.NewProjectNameSearchQuery(object.TextMethodToQuery(q.NameQuery.Method), q.NameQuery.Name)
	default:
		return nil, errors.ThrowInvalidArgument(nil, "ORG-vR9nC", "List.Query.Invalid")
	}
}

func projectStateToPb(state domain.ProjectState) proj_pb.ProjectState {
	switch state {
	case domain.ProjectStateActive:
		return proj_pb.ProjectState_PROJECT_STATE_ACTIVE
	case domain.ProjectStateInactive:
		return proj_pb.ProjectState_PROJECT_STATE_INACTIVE
	default:
		return proj_pb.ProjectState_PROJECT_STATE_UNSPECIFIED
	}
}

func projectGrantStateToPb(state domain.ProjectGrantState) proj_pb.ProjectGrantState {
	switch state {
	case domain.ProjectGrantStateActive:
		return proj_pb.ProjectGrantState_PROJECT_GRANT_STATE_ACTIVE
	case domain.ProjectGrantStateInactive:
		return proj_pb.ProjectGrantState_PROJECT_GRANT_STATE_INACTIVE
	default:
		return proj_pb.ProjectGrantState_PROJECT_GRANT_STATE_UNSPECIFIED
	}
}

func privateLabelingSettingToPb(setting domain.PrivateLabelingSetting) proj_pb.PrivateLabelingSetting {
	switch setting {
	case domain.PrivateLabelingSettingAllowLoginUserResourceOwnerPolicy:
		return proj_pb.PrivateLabelingSetting_PRIVATE_LABELING_SETTING_ALLOW_LOGIN_USER_RESOURCE_OWNER_POLICY
	case domain.PrivateLabelingSettingEnforceProjectResourceOwnerPolicy:
		return proj_pb.PrivateLabelingSetting_PRIVATE_LABELING_SETTING_ENFORCE_PROJECT_RESOURCE_OWNER_POLICY
	default:
		return proj_pb.PrivateLabelingSetting_PRIVATE_LABELING_SETTING_UNSPECIFIED
	}
}

func grantedProjectStateToPb(state proj_model.ProjectState) proj_pb.ProjectGrantState {
	switch state {
	case proj_model.ProjectStateActive:
		return proj_pb.ProjectGrantState_PROJECT_GRANT_STATE_ACTIVE
	case proj_model.ProjectStateInactive:
		return proj_pb.ProjectGrantState_PROJECT_GRANT_STATE_INACTIVE
	default:
		return proj_pb.ProjectGrantState_PROJECT_GRANT_STATE_UNSPECIFIED
	}
}

func GrantedProjectQueriesToModel(queries []*proj_pb.ProjectQuery) (_ []*proj_model.ProjectGrantViewSearchQuery, err error) {
	q := make([]*proj_model.ProjectGrantViewSearchQuery, len(queries))
	for i, query := range queries {
		q[i], err = GrantedProjectQueryToModel(query)
		if err != nil {
			return nil, err
		}
	}
	return q, nil
}

func GrantedProjectQueryToModel(query *proj_pb.ProjectQuery) (*proj_model.ProjectGrantViewSearchQuery, error) {
	switch q := query.Query.(type) {
	case *proj_pb.ProjectQuery_NameQuery:
		return GrantedProjectQueryNameToModel(q.NameQuery), nil
	default:
		return nil, errors.ThrowInvalidArgument(nil, "ORG-Ags42", "List.Query.Invalid")
	}
}

func GrantedProjectQueryNameToModel(query *proj_pb.ProjectNameQuery) *proj_model.ProjectGrantViewSearchQuery {
	return &proj_model.ProjectGrantViewSearchQuery{
		Key:    proj_model.GrantedProjectSearchKeyName,
		Method: object.TextMethodToModel(query.Method),
		Value:  query.Name,
	}
}

func RoleQueriesToModel(queries []*proj_pb.RoleQuery) (_ []*proj_model.ProjectRoleSearchQuery, err error) {
	q := make([]*proj_model.ProjectRoleSearchQuery, len(queries))
	for i, query := range queries {
		q[i], err = RoleQueryToModel(query)
		if err != nil {
			return nil, err
		}
	}
	return q, nil
}

func RoleQueryToModel(query *proj_pb.RoleQuery) (*proj_model.ProjectRoleSearchQuery, error) {
	switch q := query.Query.(type) {
	case *proj_pb.RoleQuery_KeyQuery:
		return RoleQueryKeyToModel(q.KeyQuery), nil
	case *proj_pb.RoleQuery_DisplayNameQuery:
		return RoleQueryDisplayNameToModel(q.DisplayNameQuery), nil
	default:
		return nil, errors.ThrowInvalidArgument(nil, "ORG-Ags42", "List.Query.Invalid")
	}
}

func RoleQueryKeyToModel(query *proj_pb.RoleKeyQuery) *proj_model.ProjectRoleSearchQuery {
	return &proj_model.ProjectRoleSearchQuery{
		Key:    proj_model.ProjectRoleSearchKeyKey,
		Method: object.TextMethodToModel(query.Method),
		Value:  query.Key,
	}
}

func RoleQueryDisplayNameToModel(query *proj_pb.RoleDisplayNameQuery) *proj_model.ProjectRoleSearchQuery {
	return &proj_model.ProjectRoleSearchQuery{
		Key:    proj_model.ProjectRoleSearchKeyDisplayName,
		Method: object.TextMethodToModel(query.Method),
		Value:  query.DisplayName,
	}
}

func RolesToPb(roles []*proj_model.ProjectRoleView) []*proj_pb.Role {
	r := make([]*proj_pb.Role, len(roles))
	for i, role := range roles {
		r[i] = RoleToPb(role)
	}
	return r
}

func RoleToPb(role *proj_model.ProjectRoleView) *proj_pb.Role {
	return &proj_pb.Role{
		Key:         role.Key,
		Details:     object.ToViewDetailsPb(role.Sequence, role.CreationDate, role.ChangeDate, role.ResourceOwner),
		DisplayName: role.DisplayName,
		Group:       role.Group,
	}
}
