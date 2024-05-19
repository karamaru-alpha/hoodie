package di

import (
	"context"

	"cloud.google.com/go/spanner"
	"go.uber.org/fx"

	"github.com/karamaru-alpha/days/pkg/domain/database"
	dspanner "github.com/karamaru-alpha/days/pkg/infra/spanner"
)

type SpannerConfig interface {
	GetSpannerProjectID() string
	GetSpannerInstance() string
	GetSpannerDB() string
}

var (
	NewTransactionClient = func(ctx context.Context, cfg SpannerConfig) (*spanner.Client, error) {
		return dspanner.New(ctx, cfg.GetSpannerProjectID(), cfg.GetSpannerInstance(), cfg.GetSpannerDB())
	}
	NewDBTransactionTxManager = fx.Annotate(NewSpannerTxManager, fx.As(new(database.TransactionTxManager)))
)

func NewSpannerTxManager(client *spanner.Client) database.TxManager {
	return dspanner.NewTxManager(client)
}
