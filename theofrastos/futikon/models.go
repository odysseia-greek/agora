package futikon

type Role struct {
	Privileges []string `yaml:"privileges"`
	Cluster    []string `yaml:"cluster,omitempty"`
}

type RoleMapping struct {
	Indices []string `yaml:"indices"`
	Role    Role     `yaml:"role"`
}

type UserMapping struct {
	Role string `yaml:"role"`
}

type ILMRolloverPolicy struct {
	Name   string `yaml:"name"`
	MaxAge string `yaml:"max_age"`
}

type ILMPolicyPhase struct {
	Default  bool                `yaml:"default,omitempty"`
	Rollover []ILMRolloverPolicy `yaml:"rollover"`
}

type Config struct {
	Roles    map[string]RoleMapping    `yaml:"roles"`
	Users    map[string]UserMapping    `yaml:"users"`
	Policies map[string]ILMPolicyPhase `yaml:"policies"`
}
