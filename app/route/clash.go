package route

import (
	"encoding/base64"
	"fmt"
	"net/http"

	"github.com/microctar/licorice/app/config"
	"github.com/microctar/licorice/app/facade"
	"github.com/microctar/licorice/app/utils"
	"github.com/labstack/echo/v4"
)

// config :[]byte => yaml

func ExportClashConfig(ctx echo.Context) error {
	var rfpath string

	subscribtion_link_b64 := ctx.Param("link")
	rulefilename := ctx.Param("rulefile")

	subscribtion_link, b64err := base64.RawURLEncoding.DecodeString(subscribtion_link_b64)

	if b64err != nil {
		return b64err
	}

	if rulefilename != "" {
		rfpath = fmt.Sprintf("%s/%s", config.DefaultClashConfigPath, rulefilename)
	} else {
		rfpath = config.DefaultClashRule
	}

	enc_subscribtion, online_err := utils.GetOnlineContent(string(subscribtion_link))

	if online_err != nil {
		return online_err
	}

	clash := facade.ClashConfig{}

  collect_err := clash.Collect(enc_subscribtion, config.GetDefaultConfigDirectory(), rfpath)

  if collect_err != nil {
    return collect_err
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
