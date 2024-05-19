package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"runtime"
	"strings"

	"cloud.google.com/go/spanner/admin/database/apiv1/databasepb"
	"golang.org/x/sync/errgroup"

	"github.com/karamaru-alpha/days/cmd/db-schema-sync/database"
)

type code int

const (
	success code = 0
	failure code = 1

	transactionDDLFilePath = "db/ddl/transaction.gen.sql"
)

func main() {
	os.Exit(int(cmd()))
}

func cmd() code {
	// Transaction Spanner
	transactionProjectID := os.Getenv("PROJECT_ID")
	transactionInstance := os.Getenv("SPANNER_INSTANCE")
	transactionDB := os.Getenv("SPANNER_DB")

	ctx := context.Background()
	eg, ctx := errgroup.WithContext(ctx)
	eg.Go(func() (err error) {
		defer func() {
			if r := recover(); r != nil {
				err = recoverError(r)
			}
		}()
		return syncSpannerSchema(ctx, transactionProjectID, transactionInstance, transactionDB, transactionDDLFilePath, "transaction")
	})

	if err := eg.Wait(); err != nil {
		slog.Error(fmt.Sprintf("fail to sync schema err = %+v", err))
		return failure
	}

	return success
}

func syncSpannerSchema(ctx context.Context, projectID, instance, db, ddlFilePath, target string) error {
	slog.Info(fmt.Sprintf("[Spanner (%s)] start to sync.", target))

	diff, err := database.DiffSpanner(ctx, projectID, instance, db, ddlFilePath)
	if err != nil {
		return err
	}
	if diff == "" {
		slog.Info(fmt.Sprintf("[Spanner (%s)] no changes.", target))
		return nil
	}

	diffStatements := strings.Split(diff, ";")
	validStatements := make([]string, 0, len(diffStatements))
	for _, statement := range diffStatements {
		statement = strings.TrimSpace(statement)
		if statement == "" {
			continue
		}
		validStatements = append(validStatements, statement)
	}
	if len(validStatements) == 0 {
		slog.Info(fmt.Sprintf("[Spanner (%s)] no changes.", target))
		return nil
	}

	spannerAdminClient, err := database.NewSpannerAdminClient(ctx)
	if err != nil {
		return err
	}
	defer func() {
		if err := spannerAdminClient.Close(); err != nil {
			slog.Error(fmt.Sprintf("fail to close spanner connection. err = %+v", err))
		}
	}()

	slog.Info(fmt.Sprintf("[Spanner (%s)] execute sync. ddl = %v", target, validStatements))
	if _, err := spannerAdminClient.UpdateDatabaseDdl(ctx, &databasepb.UpdateDatabaseDdlRequest{
		Database:   fmt.Sprintf("projects/%s/instances/%s/databases/%s", projectID, instance, db),
		Statements: validStatements,
	}); err != nil {
		return err
	}

	slog.Info(fmt.Sprintf("[Spanner (%s)] finished to sync.", target))
	return nil
}

func recoverError(r any) error {
	var stacktrace string
	for depth := 0; ; depth++ {
		_, file, line, ok := runtime.Caller(depth)
		if !ok {
			break
		}
		stacktrace += fmt.Sprintf("        %v:%d\n", file, line)
	}
	return fmt.Errorf("panic occurred. recovered = %v\nstacktrace = %s", r, stacktrace)
}
