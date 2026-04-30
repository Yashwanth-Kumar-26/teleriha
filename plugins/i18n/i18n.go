package i18n

import (
	"fmt"

	"github.com/Yashwanth-Kumar-26/teleriha/pkg/bot"
)

// I18nPlugin implements the bot.Plugin interface.
type I18nPlugin struct {
	bot.BasePlugin
	defaultLang  string
	translations map[string]map[string]string
}

// New creates a new I18nPlugin.
func New(defaultLang string, translations map[string]map[string]string) *I18nPlugin {
	return &I18nPlugin{
		BasePlugin:   *bot.NewBasePlugin("i18n"),
		defaultLang:  defaultLang,
		translations: translations,
	}
}

// Init initializes the plugin.
func (p *I18nPlugin) Init(b *bot.Bot) error {
	return p.BasePlugin.Init(b)
}

// Translate returns the translated string for the given key and context.
func (p *I18nPlugin) Translate(ctx *bot.Context, key string, args ...interface{}) string {
	lang := p.defaultLang
	if ctx.Sender != nil && ctx.Sender.LanguageCode != "" {
		if _, ok := p.translations[ctx.Sender.LanguageCode]; ok {
			lang = ctx.Sender.LanguageCode
		}
	}

	langTrans, ok := p.translations[lang]
	if !ok {
		return key
	}

	format, ok := langTrans[key]
	if !ok {
		// Fallback to default lang if current is not default
		if lang != p.defaultLang {
			if defTrans, ok := p.translations[p.defaultLang]; ok {
				if f, ok := defTrans[key]; ok {
					format = f
				} else {
					return key
				}
			} else {
				return key
			}
		} else {
			return key
		}
	}

	return fmt.Sprintf(format, args...)
}
