package config

// Task task-related config.
type Task struct {
	ForceDefaultAction bool `json:"defaultActionForce"`
}

// SetDefaults provides default settings.
func (t *Task) SetDefaults() {}
