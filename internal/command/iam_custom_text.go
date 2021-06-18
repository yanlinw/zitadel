package command

import (
	"context"

	"github.com/caos/zitadel/internal/domain"
	caos_errs "github.com/caos/zitadel/internal/errors"
	iam_repo "github.com/caos/zitadel/internal/repository/iam"
	"github.com/caos/zitadel/internal/telemetry/tracing"
)

func (c *Commands) AddDefaultCustomTextFile(ctx context.Context, storageKey string) (*domain.ObjectDetails, error) {
	if storageKey == "" {
		return nil, caos_errs.ThrowInvalidArgument(nil, "IAM-4n9fs", "Errors.Assets.EmptyKey")
	}
	iamAgg := iam_repo.NewAggregate()
	_, err := c.eventstore.PushEvents(ctx, iam_repo.NewCustomTextFileUploadedEvent(ctx, &iamAgg.Aggregate, storageKey))
	if err != nil {
		return nil, err
	}
	return nil, nil
	//err = AppendAndReduce(existingPolicy, pushedEvents...)
	//if err != nil {
	//	return nil, err
	//}
	//return writeModelToObjectDetails(&existingPolicy.LabelPolicyWriteModel.WriteModel), nil
}

func (c *Commands) RemoveDefaultCustomTextFile(ctx context.Context) (*domain.ObjectDetails, error) {
	existingPolicy, err := c.defaultLabelPolicyWriteModelByID(ctx)
	if err != nil {
		return nil, err
	}

	if existingPolicy.State == domain.PolicyStateUnspecified || existingPolicy.State == domain.PolicyStateRemoved {
		return nil, caos_errs.ThrowNotFound(nil, "IAM-Xc8Kf", "Errors.IAM.LabelPolicy.NotFound")
	}

	err = c.RemoveAsset(ctx, domain.IAMID, existingPolicy.LogoKey)
	if err != nil {
		return nil, err
	}
	iamAgg := IAMAggregateFromWriteModel(&existingPolicy.LabelPolicyWriteModel.WriteModel)
	pushedEvents, err := c.eventstore.PushEvents(ctx, iam_repo.NewLabelPolicyLogoRemovedEvent(ctx, iamAgg, existingPolicy.LogoKey))
	if err != nil {
		return nil, err
	}
	err = AppendAndReduce(existingPolicy, pushedEvents...)
	if err != nil {
		return nil, err
	}
	return writeModelToObjectDetails(&existingPolicy.LabelPolicyWriteModel.WriteModel), nil
}

func (c *Commands) defaultLabelPolicyWriteModelByID(ctx context.Context) (policy *IAMLabelPolicyWriteModel, err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()

	writeModel := NewIAMLabelPolicyWriteModel()
	err = c.eventstore.FilterToQueryReducer(ctx, writeModel)
	if err != nil {
		return nil, err
	}
	return writeModel, nil
}

func (c *Commands) getDefaultLabelPolicy(ctx context.Context) (*domain.LabelPolicy, error) {
	policyWriteModel, err := c.defaultLabelPolicyWriteModelByID(ctx)
	if err != nil {
		return nil, err
	}
	policy := writeModelToLabelPolicy(&policyWriteModel.LabelPolicyWriteModel)
	policy.Default = true
	return policy, nil
}
