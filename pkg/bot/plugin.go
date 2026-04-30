package bot

import (
	"context"
	"fmt"
	"sync"
)

// Plugin is the interface that all TeleRiHa plugins must implement.
type Plugin interface {
	// Name returns the plugin's name.
	Name() string

	// Init is called when the plugin is loaded.
	// The bot parameter provides access to the bot instance.
	Init(bot *Bot) error

	// Start is called when the bot starts.
	Start(ctx context.Context) error

	// Stop is called when the bot stops.
	// This should clean up any resources.
	Stop() error
}

// PluginFunc is a function that creates a Plugin.
// This is useful for plugin registration.
type PluginFunc func(*Bot) Plugin

// PluginRegistry manages plugin registration and lifecycle.
type PluginRegistry struct {
	// plugins is the map of plugin names to Plugin instances
	plugins map[string]Plugin

	// pluginFuncs is the map of plugin names to PluginFunc
	pluginFuncs map[string]PluginFunc

	// mu protects the plugins map
	mu sync.RWMutex
}

// NewPluginRegistry creates a new PluginRegistry.
func NewPluginRegistry() *PluginRegistry {
	return &PluginRegistry{
		plugins:    make(map[string]Plugin),
		pluginFuncs: make(map[string]PluginFunc),
	}
}

// Register registers a PluginFunc with the given name.
func (r *PluginRegistry) Register(name string, fn PluginFunc) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.pluginFuncs[name] = fn
}

// RegisterPlugin registers an already instantiated Plugin.
func (r *PluginRegistry) RegisterPlugin(plugin Plugin) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.plugins[plugin.Name()] = plugin
}

// Load loads a plugin by name and initializes it with the bot.
func (r *PluginRegistry) Load(name string, bot *Bot) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	fn, ok := r.pluginFuncs[name]
	if !ok {
		return fmt.Errorf("plugin %s not found", name)
	}

	plugin := fn(bot)
	if err := plugin.Init(bot); err != nil {
		return fmt.Errorf("failed to init plugin %s: %w", name, err)
	}

	r.plugins[name] = plugin
	return nil
}

// LoadAll loads all registered plugin functions.
func (r *PluginRegistry) LoadAll(bot *Bot) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	for name, fn := range r.pluginFuncs {
		plugin := fn(bot)
		if err := plugin.Init(bot); err != nil {
			return fmt.Errorf("failed to init plugin %s: %w", name, err)
		}
		r.plugins[name] = plugin
	}

	return nil
}

// StartAll starts all loaded plugins.
func (r *PluginRegistry) StartAll(ctx context.Context) error {
	r.mu.RLock()
	defer r.mu.RUnlock()

	for name, plugin := range r.plugins {
		if err := plugin.Start(ctx); err != nil {
			return fmt.Errorf("failed to start plugin %s: %w", name, err)
		}
	}

	return nil
}

// StopAll stops all loaded plugins.
func (r *PluginRegistry) StopAll() error {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var firstErr error
	for _, plugin := range r.plugins {
		if err := plugin.Stop(); err != nil {
			if firstErr == nil {
				firstErr = err
			}
			// Log the error but continue stopping other plugins
		}
	}

	return firstErr
}

// Get returns a loaded plugin by name.
func (r *PluginRegistry) Get(name string) (Plugin, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	plugin, ok := r.plugins[name]
	return plugin, ok
}

// List returns the names of all loaded plugins.
func (r *PluginRegistry) List() []string {
	r.mu.RLock()
	defer r.mu.RUnlock()

	names := make([]string, 0, len(r.plugins))
	for name := range r.plugins {
		names = append(names, name)
	}

	return names
}

// Unload unloads a plugin by name.
func (r *PluginRegistry) Unload(name string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	plugin, ok := r.plugins[name]
	if !ok {
		return fmt.Errorf("plugin %s not loaded", name)
	}

	if err := plugin.Stop(); err != nil {
		return fmt.Errorf("failed to stop plugin %s: %w", name, err)
	}

	delete(r.plugins, name)
	return nil
}

// BasePlugin provides a base implementation for plugins.
// Embed this in your plugin to get default implementations.
type BasePlugin struct {
	name   string
	bot    *Bot
	active bool
}

// NewBasePlugin creates a new BasePlugin.
func NewBasePlugin(name string) *BasePlugin {
	return &BasePlugin{
		name: name,
	}
}

// Name returns the plugin's name.
func (p *BasePlugin) Name() string {
	return p.name
}

// Init initializes the plugin with the bot.
func (p *BasePlugin) Init(bot *Bot) error {
	p.bot = bot
	p.active = true
	return nil
}

// Start starts the plugin.
func (p *BasePlugin) Start(ctx context.Context) error {
	p.active = true
	return nil
}

// Stop stops the plugin.
func (p *BasePlugin) Stop() error {
	p.active = false
	return nil
}

// Bot returns the bot instance.
func (p *BasePlugin) Bot() *Bot {
	return p.bot
}

// IsActive returns whether the plugin is active.
func (p *BasePlugin) IsActive() bool {
	return p.active
}

// PluginOption is a function that configures a plugin.
type PluginOption func(*Bot)

// WithPlugin adds a plugin to the bot.
func WithPlugin(plugin Plugin) BotOption {
	return func(b *Bot) {
		// Store the plugin in the bot's router or context
		// For now, we'll add it to the bot's plugins map
		// This is a placeholder; the actual implementation
		// would depend on how you want to manage plugins
		b.plugins[plugin.Name()] = plugin
	}
}
