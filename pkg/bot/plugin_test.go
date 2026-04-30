package bot

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

type mockPlugin struct {
	BasePlugin
	initCalled  bool
	startCalled bool
	stopCalled  bool
}

func (m *mockPlugin) Init(bot *Bot) error {
	m.initCalled = true
	return m.BasePlugin.Init(bot)
}

func (m *mockPlugin) Start(ctx context.Context) error {
	m.startCalled = true
	return m.BasePlugin.Start(ctx)
}

func (m *mockPlugin) Stop() error {
	m.stopCalled = true
	return m.BasePlugin.Stop()
}

func TestPluginRegistry(t *testing.T) {
	r := NewPluginRegistry()

	p := &mockPlugin{BasePlugin: *NewBasePlugin("test")}
	r.RegisterPlugin(p)

	assert.Equal(t, []string{"test"}, r.List())
	
	p2, ok := r.Get("test")
	assert.True(t, ok)
	assert.Equal(t, p, p2)

	err := r.StartAll(context.Background())
	assert.NoError(t, err)
	assert.True(t, p.startCalled)

	err = r.StopAll()
	assert.NoError(t, err)
	assert.True(t, p.stopCalled)

	err = r.Unload("test")
	assert.NoError(t, err)
	assert.Len(t, r.List(), 0)
}

func TestPluginRegistry_Load(t *testing.T) {
	r := NewPluginRegistry()
	b := New("token")

	r.Register("test", func(bot *Bot) Plugin {
		return &mockPlugin{BasePlugin: *NewBasePlugin("test")}
	})

	err := r.Load("test", b)
	assert.NoError(t, err)
	
	p, ok := r.Get("test")
	assert.True(t, ok)
	assert.True(t, p.(*mockPlugin).initCalled)
}

func TestPluginRegistry_LoadAll(t *testing.T) {
	r := NewPluginRegistry()
	b := New("token")

	r.Register("p1", func(bot *Bot) Plugin {
		return &mockPlugin{BasePlugin: *NewBasePlugin("p1")}
	})
	r.Register("p2", func(bot *Bot) Plugin {
		return &mockPlugin{BasePlugin: *NewBasePlugin("p2")}
	})

	err := r.LoadAll(b)
	assert.NoError(t, err)
	assert.ElementsMatch(t, []string{"p1", "p2"}, r.List())
}

func TestBasePlugin(t *testing.T) {
	p := NewBasePlugin("test")
	assert.Equal(t, "test", p.Name())
	assert.False(t, p.IsActive())

	b := New("token")
	p.Init(b)
	assert.Equal(t, b, p.Bot())
	assert.True(t, p.IsActive())

	p.Stop()
	assert.False(t, p.IsActive())
}

func TestWithPlugin(t *testing.T) {
	p := NewBasePlugin("test")
	b := New("token", WithPlugin(p))
	assert.Equal(t, p, b.plugins["test"])
}
