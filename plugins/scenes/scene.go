package scenes

import (
	"strconv"

	"github.com/Yashwanth-Kumar-26/teleriha/pkg/bot"
	"github.com/Yashwanth-Kumar-26/teleriha/plugins/session"
)

// SceneHandler handles a specific step in a scene.
type SceneHandler func(ctx *bot.Context) error

// Scene represents a state machine with multiple steps.
type Scene struct {
	Name  string
	Steps []SceneHandler
}

// NewScene creates a new scene with the given name and handlers.
func NewScene(name string, steps ...SceneHandler) *Scene {
	return &Scene{
		Name:  name,
		Steps: steps,
	}
}

// SceneManager manages multiple scenes and user progression.
type SceneManager struct {
	scenes map[string]*Scene
}

// NewSceneManager creates a new scene manager.
func NewSceneManager() *SceneManager {
	return &SceneManager{
		scenes: make(map[string]*Scene),
	}
}

// AddScene registers a scene.
func (sm *SceneManager) AddScene(s *Scene) {
	sm.scenes[s.Name] = s
}

const (
	sceneNameKey = "scene_name"
	sceneStepKey = "scene_step"
)

// Enter enters a user into a scene.
func (sm *SceneManager) Enter(ctx *bot.Context, sceneName string) error {
	sess := session.GetSession(ctx)
	if sess == nil {
		return nil
	}
	scene, ok := sm.scenes[sceneName]
	if !ok {
		return nil
	}
	_ = sess.Set(sceneNameKey, []byte(sceneName))
	_ = sess.Set(sceneStepKey, []byte("0"))
	
	if len(scene.Steps) > 0 {
		return scene.Steps[0](ctx)
	}
	return nil
}

// Leave forces the user to exit the current scene.
func (sm *SceneManager) Leave(ctx *bot.Context) error {
	sess := session.GetSession(ctx)
	if sess == nil {
		return nil
	}
	_ = sess.Set(sceneNameKey, []byte(""))
	_ = sess.Set(sceneStepKey, []byte("0"))
	return nil
}

// Middleware creates a bot middleware to intercept updates for active scenes.
func (sm *SceneManager) Middleware() bot.Middleware {
	return func(next bot.Handler) bot.Handler {
		return func(ctx *bot.Context) error {
			sess := session.GetSession(ctx)
			if sess == nil {
				return next(ctx)
			}
			
			nameBytes, _ := sess.Get(sceneNameKey)
			sceneName := string(nameBytes)
			if sceneName == "" {
				return next(ctx)
			}
			
			scene, ok := sm.scenes[sceneName]
			if !ok {
				return next(ctx)
			}

			stepBytes, _ := sess.Get(sceneStepKey)
			step := 0
			if len(stepBytes) > 0 {
				step, _ = strconv.Atoi(string(stepBytes))
			}

			if step >= len(scene.Steps) {
				_ = sm.Leave(ctx)
				return next(ctx)
			}

			err := scene.Steps[step](ctx)
			
			// Auto advance if still in same scene and didn't hit an error
			if err == nil {
				newNameBytes, _ := sess.Get(sceneNameKey)
				if string(newNameBytes) == sceneName {
					newStep := step + 1
					if newStep < len(scene.Steps) {
						_ = sess.Set(sceneStepKey, []byte(strconv.Itoa(newStep)))
					} else {
						_ = sm.Leave(ctx)
					}
				}
			}
			
			// Return err so the core router handles it
			return err
		}
	}
}
