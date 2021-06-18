package assets

import (
	"context"
	"strings"

	"github.com/caos/zitadel/internal/api/authz"
	"github.com/caos/zitadel/internal/command"
	"github.com/caos/zitadel/internal/domain"
	"github.com/caos/zitadel/internal/iam/model"
	"github.com/caos/zitadel/internal/management/repository"
)

func (h *Handler) UploadDefaultTranslationFile() Uploader {
	return &translationFileUploader{true, []string{"image/"}, 1 << 19}
}

func (h *Handler) UploadTranslationFile() Uploader {
	return &translationFileUploader{false, []string{"image/"}, 1 << 19}
}

type translationFileUploader struct {
	defaultTranslations bool
	contentTypes        []string
	maxSize             int64
}

func (l *translationFileUploader) ContentTypeAllowed(contentType string) bool {
	for _, ct := range l.contentTypes {
		if strings.HasPrefix(contentType, ct) {
			return true
		}
	}
	return false
}

func (l *translationFileUploader) MaxFileSize() int64 {
	return l.maxSize
}

func (l *translationFileUploader) ObjectName(_ authz.CtxData) (string, error) {
	return domain.TranslationFilePath, nil
}

func (l *translationFileUploader) BucketName(ctxData authz.CtxData) string {
	if l.defaultTranslations {
		return domain.IAMID
	}
	return ctxData.OrgID
}

func (l *translationFileUploader) Callback(ctx context.Context, info *domain.AssetInfo, orgID string, commands *command.Commands) error {
	_, err := commands.AddLogoLabelPolicy(ctx, orgID, info.Key)
	return err
}

func (h *Handler) GetDefaultTranslationFile() Downloader {
	return &translationFileDownloader{org: h.orgRepo, defaultTranslations: true}
}

func (h *Handler) GetTranslationFile() Downloader {
	return &translationFileDownloader{org: h.orgRepo, defaultTranslations: true}
}

type translationFileDownloader struct {
	org                 repository.OrgRepository
	defaultTranslations bool
}

func (l *translationFileDownloader) ObjectName(ctx context.Context, path string) (string, error) {
	policy, err := getTranslationFile(ctx, l.defaultTranslations, l.org)
	if err != nil {
		return "", nil
	}
	return policy.LogoURL, nil
}

func (l *translationFileDownloader) BucketName(ctx context.Context, id string) string {
	return getTranslationFileBucketName(ctx, l.defaultTranslations, l.org)
}

func getTranslationFile(ctx context.Context, defaultTranslations bool, orgRepo repository.OrgRepository) (*model.LabelPolicyView, error) {
	if defaultTranslations {
		return orgRepo.GetDefaultLabelPolicy(ctx)
	}
	return orgRepo.GetLabelPolicy(ctx)
}

func getTranslationFileBucketName(ctx context.Context, defaultTranslations bool, org repository.OrgRepository) string {
	if defaultTranslations {
		return domain.IAMID
	}
	policy, err := getLabelPolicy(ctx, defaultTranslations, org)
	if err != nil {
		return ""
	}
	if policy.Default {
		return domain.IAMID
	}
	return authz.GetCtxData(ctx).OrgID
}
