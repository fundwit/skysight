package localize

import (
	"context"
	"errors"
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"golang.org/x/text/language"
	"gopkg.in/yaml.v2"
)

const (
	defaultRootPath = "./i18n"
)

var (
	defaultLanguage       = language.English
	defaultAcceptLanguage = []language.Tag{
		defaultLanguage,
		language.Chinese,
	}

	defaultBundleConfig = &BundleCfg{
		RootPath:        defaultRootPath,
		AcceptLanguage:  defaultAcceptLanguage,
		DefaultLanguage: defaultLanguage,
	}
)

// defaultGetLngHandler ...
func defaultGetLngHandler(context *gin.Context, defaultLng string) string {
	if context == nil || context.Request == nil {
		return defaultLng
	}

	lng := context.Query("lang")
	if lng != "" {
		return lng
	}

	lng = context.GetHeader("Accept-Language")
	if lng != "" {
		return lng
	}

	return defaultLng
}

type GinI18n struct {
	bundleConfig *BundleCfg

	bundle          *i18n.Bundle
	currentContext  *gin.Context
	localizerByLng  map[string]*i18n.Localizer
	defaultLanguage language.Tag
	getLngHandler   GetLngHandler
}

func (i *GinI18n) setBundleConfig(cfg *BundleCfg) {
	i.bundleConfig = cfg
}

// getMessage get localize message by lng and messageID
func (i *GinI18n) getMessage(param interface{}) (string, error) {
	defaultLang := i.defaultLanguage.String()
	lng := i.getLngHandler(i.currentContext, defaultLang)
	localizer := i.getLocalizerByLng(lng)

	localizeConfig, err := i.getLocalizeConfig(param)
	if err != nil {
		return fmt.Sprint(param), err
	}

	message, err := localizer.Localize(localizeConfig)
	if err != nil {
		return fmt.Sprint(param), err
	}

	return message, nil
}

// mustGetMessage ...
func (i *GinI18n) mustGetMessage(param interface{}) string {
	message, _ := i.getMessage(param)
	return message
}

func (i *GinI18n) setCurrentContext(ctx context.Context) {
	i.currentContext = ctx.(*gin.Context)
}

func (i *GinI18n) setBundle(cfg *BundleCfg) {
	bundle := i18n.NewBundle(cfg.DefaultLanguage)
	bundle.RegisterUnmarshalFunc("yaml", yaml.Unmarshal)

	i.bundle = bundle
	i.defaultLanguage = cfg.DefaultLanguage

	i.loadMessageFiles(cfg)
	i.setLocalizerByLng(cfg.AcceptLanguage)
}

func (i *GinI18n) setGetLngHandler(handler GetLngHandler) {
	i.getLngHandler = handler
}

// loadMessageFiles load all file localize to bundle
func (i *GinI18n) loadMessageFiles(config *BundleCfg) {
	for _, lng := range config.AcceptLanguage {
		path := config.RootPath + "/" + lng.String() + ".yaml"
		i.bundle.MustLoadMessageFile(path)
	}
}

// setLocalizerByLng set localizer by language
func (i *GinI18n) setLocalizerByLng(acceptLanguage []language.Tag) {
	i.localizerByLng = map[string]*i18n.Localizer{}
	for _, lng := range acceptLanguage {
		lngStr := lng.String()
		i.localizerByLng[lngStr] = i.newLocalizer(lngStr)
	}

	// set defaultLanguage if it isn't exist
	defaultLng := i.defaultLanguage.String()
	if _, hasDefaultLng := i.localizerByLng[defaultLng]; !hasDefaultLng {
		i.localizerByLng[defaultLng] = i.newLocalizer(defaultLng)
	}
}

// newLocalizer create a localizer by language
func (i *GinI18n) newLocalizer(lng string) *i18n.Localizer {
	lngDefault := i.defaultLanguage.String()
	langs := []string{
		lng,
	}

	if lng != lngDefault {
		langs = append(langs, lngDefault)
	}

	localizer := i18n.NewLocalizer(
		i.bundle,
		langs...,
	)
	return localizer
}

// getLocalizerByLng get localizer by language
func (i *GinI18n) getLocalizerByLng(lng string) *i18n.Localizer {
	acceptLangs := ParseAcceptLanguage(lng)

	for _, al := range acceptLangs {
		localizer, hasValue := i.localizerByLng[al.Lang]
		if hasValue {
			return localizer
		}
	}

	return i.localizerByLng[i.defaultLanguage.String()]
}

func (i *GinI18n) getLocalizeConfig(param interface{}) (*i18n.LocalizeConfig, error) {
	switch paramValue := param.(type) {
	case string:
		localizeConfig := &i18n.LocalizeConfig{
			MessageID: paramValue,
		}
		return localizeConfig, nil
	case *i18n.LocalizeConfig:
		return paramValue, nil
	}

	msg := fmt.Sprintf("un supported localize param: %v", param)
	return nil, errors.New(msg)
}
