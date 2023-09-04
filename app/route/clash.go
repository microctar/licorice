package route

import (
	"encoding/base64"
	"fmt"
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/microctar/licorice/app/config"
	"github.com/microctar/licorice/app/facade"
	"github.com/microctar/licorice/app/utils"
	"github.com/patrickmn/go-cache"
)

// config :[]byte => yaml

func ExportClashConfig(cache *cache.Cache, acldir string, defaultrulefile string) echo.HandlerFunc {
	return func(ctx echo.Context) error {
		var b64Links strings.Builder

		subsLinkB64 := ctx.Param("link")
		rulefile := ctx.Param("rulefile")

		if subsLink, b64err := base64.RawURLEncoding.DecodeString(subsLinkB64); b64err != nil {
			return b64err
		} else {
			b64Links.Write(subsLink)
		}

		rfpath := fmt.Sprintf("%s/%s", config.DefaultClashRulePath, rulefile)

		if rulefile == "" {
			rfpath = defaultrulefile
		}

		encSubscription, onlineErr := utils.GetOnlineContent(b64Links.String())

		if onlineErr != nil {
			return onlineErr
		}

		clash := facade.NewCachedGenerator("clash", cache)

		if collectErr := clash.Collect(encSubscription, acldir, rfpath); collectErr != nil {
			return collectErr
		}

		data, epterr := clash.Export()

		if epterr != nil {
			return epterr
		}

		// Content-Disposition reference => https://developer.mozilla.org/en-US/docs/Web/HTTP/Headers/Content-Disposition
		ctx.Response().Header().Set("Content-Disposition", `attachment; filename="config.yaml"`)

		ctx.Response().WriteHeader(http.StatusOK)
		ctx.Response().Writer.Header().Add("Content-Type", "application/x-yaml; charset=UTF-8")

		_, err := ctx.Response().Writer.Write(data)

		return err
	}
}
