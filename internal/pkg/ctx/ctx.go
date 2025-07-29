package ctx

import (
	"context"

	"github.com/Eizeed/2025-07-29/internal/pkg/config"
	"github.com/Eizeed/2025-07-29/pkg/assert"
)

func GetAppConfigFromConfig(ctx context.Context) *config.AppConfig {
	cfg, ok := ctx.Value(AppConfigKey{}).(*config.AppConfig)
	assert.Assert(ok, "Assertion failed. AppConfig from Context is expected to be there")

	return cfg
}
