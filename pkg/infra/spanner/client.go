package spanner

import (
	"context"
	"fmt"

	"cloud.google.com/go/spanner"

	"github.com/karamaru-alpha/days/pkg/derrors"
)

func New(ctx context.Context, projectID, instance, db string) (*spanner.Client, error) {
	database := fmt.Sprintf("projects/%s/instances/%s/databases/%s", projectID, instance, db)
	client, err := spanner.NewClient(ctx, database)
	if err != nil {
		return nil, derrors.Wrap(err, derrors.Internal, err.Error())
	}

	return client, nil
}
