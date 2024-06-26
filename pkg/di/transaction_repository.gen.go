// Code generated by protoc-gen-days (generator/spanner) . DO NOT EDIT.
// This file is meant to be re-generated in place and/or deleted at any time.

package di

import (
	"go.uber.org/fx"

	"github.com/karamaru-alpha/days/pkg/infra/spanner/repository"
)

var TransactionRepositorySet = fx.Provide(
	repository.NewUserRepository,
	repository.NewUserExampleRepository,
)
