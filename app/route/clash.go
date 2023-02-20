package route

import (
	"encoding/base64"
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/microctar/licorice/app/config"
	"github.com/microctar/licorice/app/facade"
	"github.com/microctar/licorice/app/utils"
	"github.com/patrickmn/go-cache"
)

// config :[]byte => yaml

func ExportClashConfig(cache *cache.Cache, acldir string, defaultrulefile string) echo.HandlerFunc {
	return func(ctx echo.Context) error {
		var rfpath string

		subsLinkB64 := ctx.Param("link")
		rulefile := ctx.Param("rulefile")

		subsLink, b64err := base64.RawURLEncoding.DecodeString(subsLinkB64)

		if b64err != nil {
			return b64err
		}

		if rulefile == "" {
			rfpath = defaultrulefile
		} else {
			rfpath = fmt.Sprintf("%s/%s", config.DefaultClashRulePath, rulefile)
		}

		encSubscription, onlineErr := utils.GetOnlineContent(string(subsLink))

		if onlineErr != nil {
			return onlineErr
		}

		clash := facade.NewCachedGenerator("clash", cache)

		collectErr := clash.Collect(encSubscription, acldir, rfpath)

		if collectErr != nil {
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
