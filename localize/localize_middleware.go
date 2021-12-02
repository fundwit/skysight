package localize

import (
	"github.com/gin-gonic/gin"
	"golang.org/x/text/language"
)

func LocalizeMiddleware(i18nPath string) gin.HandlerFunc {
	bundleOption := WithBundle(&BundleCfg{
		RootPath:        i18nPath,
		AcceptLanguage:  []language.Tag{language.Chinese, language.English},
		DefaultLanguage: language.English,
	})

	return Localize(bundleOption)
}
