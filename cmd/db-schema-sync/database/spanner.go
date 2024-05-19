package database

import (
	"context"
	"fmt"
	"os/exec"

	spanner "cloud.google.com/go/spanner/admin/database/apiv1"
)

func NewSpannerAdminClient(ctx context.Context) (*spanner.DatabaseAdminClient, error) {
	cli, err := spanner.NewDatabaseAdminClient(ctx)
	if err != nil {
		return nil, fmt.Errorf("fail to connect to spanner. err: %w", err)
	}
	return cli, nil
}

func DiffSpanner(ctx context.Context, projectID, instance, db, ddlFilePath string) (string, error) {
	url := fmt.Sprintf("spanner://projects/%s/instances/%s/databases/%s?x-clean-statements=true", projectID, instance, db)
	output, err := exec.CommandContext(ctx, "hammer", "diff", url, ddlFilePath).CombinedOutput()
	if err != nil {
		// outputに標準エラー出力の値が入る
		return "", fmt.Errorf("fail to get schema diff. output = %s: %w", string(output), err)
	}
	return string(output), nil
}
