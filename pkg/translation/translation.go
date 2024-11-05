package translation

import (
	"context"
	"os"
	"path/filepath"

	"github.com/dathuynh1108/clean-arch-base/pkg/comjson"
	"github.com/dathuynh1108/clean-arch-base/pkg/config"
	"github.com/dathuynh1108/clean-arch-base/pkg/logger"
	"github.com/dathuynh1108/clean-arch-base/pkg/singleton"
	"github.com/nicksnyder/go-i18n/v2/i18n"

	"golang.org/x/text/language"
)

const (
	TranslationNotFoundMessage = "translation_not_found"
)

const (
	ContextKeyLocalizer string = "localizer" // Use string for set to echo context
)

var (
	bundleSingleton = singleton.NewSingleton(
		func() *i18n.Bundle {
			bundle := i18n.NewBundle(language.English)
			bundle.RegisterUnmarshalFunc("json", comjson.Unmarshal)
			return bundle
		},
		false,
	)
)

func GetBundle() *i18n.Bundle {
	return bundleSingleton.Get()
}

func Load() error {
	var (
		bundle = bundleSingleton.Get()
		folder = config.GetConfig().TranslationConfig.Folder
		logger = logger.GetLogger()
	)

	return filepath.Walk(folder, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !info.IsDir() && filepath.Ext(info.Name()) == ".json" {
			logger.Info("Loading translation file", path)
			bundle.MustLoadMessageFile(path)
		}

		return nil
	})
}

func TranslateKey(ctx context.Context, key string) (string, error) {
	localizer, ok := GetLocallizer(ctx)
	if !ok {
		return key, nil
	}

	message, err := localizer.Localize(&i18n.LocalizeConfig{
		MessageID: key,
	})
	if err != nil {
		switch err.(type) {
		case *i18n.MessageNotFoundErr:
			// return localizer.Localize(&i18n.LocalizeConfig{
			// 	MessageID: TranslationNotFoundMessage,
			// })
			return key, nil // Return key if translation not found
		default:
			return message, err
		}
	}

	return message, err
}

func TranslateKeyData(ctx context.Context, key string, data any) (string, error) {
	localizer, ok := GetLocallizer(ctx)
	if !ok {
		return key, nil
	}

	message, err := localizer.Localize(&i18n.LocalizeConfig{
		MessageID:    key,
		TemplateData: data,
	})
	if err != nil {
		switch err.(type) {
		case *i18n.MessageNotFoundErr:
			// return localizer.Localize(&i18n.LocalizeConfig{
			// 	MessageID: TranslationNotFoundMessage,
			// })
			return key, nil // Return key if translation not found
		default:
			return message, err
		}
	}

	return message, err
}

func GetLocallizer(ctx context.Context) (*i18n.Localizer, bool) {
	localizerObj := ctx.Value(ContextKeyLocalizer)
	if localizerObj == nil {
		return nil, false
	}

	localizer, ok := localizerObj.(*i18n.Localizer)
	return localizer, ok
}
