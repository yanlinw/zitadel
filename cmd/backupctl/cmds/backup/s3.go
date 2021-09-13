package backup

import (
	"context"
	"github.com/caos/orbos/mntr"
	"github.com/caos/zitadel/cmd/backupctl/cmds/helpers"
	"github.com/caos/zitadel/pkg/backup"
	"github.com/spf13/cobra"
	"io/ioutil"
)

func S3Command(ctx context.Context, monitor mntr.Monitor) *cobra.Command {
	var (
		backupName    string
		backupNameEnv string
		assetEndpoint string
		assetAKID     string
		assetSAK      string
		assetPrefix   string
		assetRegion   string
		destEndpoint  string
		destBucket    string
		destAKID      string
		destSAK       string
		destRegion    string
		configPath    string
		certsDir      string
		host          string
		port          string
		cmd           = &cobra.Command{
			Use:   "s3",
			Short: "Backup to S3 storage",
			Long:  "Backup to S3 storage",
		}
	)

	flags := cmd.Flags()
	flags.StringVar(&backupName, "backupname", "", "Backupname used in destination file path")
	flags.StringVar(&backupNameEnv, "backupnameenv", "", "Backupnameenv used in destination file path")
	flags.StringVar(&assetEndpoint, "asset-endpoint", "", "Endpoint for the asset S3 storage")
	flags.StringVar(&assetAKID, "asset-akid", "", "AccessKeyID for the asset S3 storage")
	flags.StringVar(&assetSAK, "asset-sak", "", "SecretAccessKey for the asset S3 storage")
	flags.StringVar(&assetPrefix, "asset-prefix", "", "Bucket-Prefix in the asset S3 storage")
	flags.StringVar(&assetRegion, "asset-region", "", "Bucket-Region for the asset S3 storage")
	flags.StringVar(&destEndpoint, "destination-endpoint", "", "Endpoint for the destination S3 storage")
	flags.StringVar(&destAKID, "destination-akid", "", "AccessKeyID for the destination S3 storage")
	flags.StringVar(&destSAK, "destination-sak", "", "SecretAccessKey for the destination S3 storage")
	flags.StringVar(&destBucket, "destination-bucket", "", "Bucketname in the destination S3 storage")
	flags.StringVar(&destRegion, "destination-region", "", "Region for the destination S3 storage")
	flags.StringVar(&configPath, "configpath", "", "Path used to save rsync configuration")
	flags.StringVar(&certsDir, "certs-dir", "", "Folder with certificates used to connect to cockroachdb")
	flags.StringVar(&host, "host", "", "Host used to connect to cockroachdb")
	flags.StringVar(&port, "port", "", "Port used to connect to cockroachdb")

	cmd.RunE = func(cmd *cobra.Command, args []string) (err error) {

		if err := helpers.ValidateBackupFlags(
			backupName,
			backupNameEnv,
		); err != nil {
			return err
		}

		if err := helpers.ValidateDestinationS3Flags(
			destEndpoint,
			destAKID,
			destSAK,
			destBucket,
			destRegion,
		); err != nil {
			return err
		}

		if err := helpers.ValidateSourceS3Flags(
			assetEndpoint,
			assetAKID,
			assetSAK,
			assetPrefix,
			assetRegion,
		); err != nil {
			return err
		}

		if err := helpers.ValidateCockroachFlags(
			certsDir,
			host,
			port,
		); err != nil {
			return err
		}

		destinationAKIDData, err := ioutil.ReadFile(destAKID)
		if err != nil {
			return err
		}

		destinationSAKData, err := ioutil.ReadFile(destSAK)
		if err != nil {
			return err
		}

		if err := backup.CreateS3BucketIfNotExists(destEndpoint, string(destinationAKIDData), string(destinationSAKData), destBucket, destRegion); err != nil {
			return err
		}

		if err := backup.RsyncBackupS3ToS3(
			ctx,
			backupName,
			backupNameEnv,
			"destination",
			destEndpoint,
			destAKID,
			destSAK,
			destBucket,
			destRegion,
			"source",
			assetEndpoint,
			assetAKID,
			assetSAK,
			assetPrefix,
			assetRegion,
			configPath,
		); err != nil {
			return err
		}

		if err := backup.CockroachBackupToS3(
			ctx,
			certsDir,
			destBucket,
			backupName,
			backupNameEnv,
			host,
			port,
			destAKID,
			destSAK,
			destEndpoint,
			destRegion,
		); err != nil {
			return err
		}

		return nil
	}
	return cmd
}
