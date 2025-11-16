package types

type StackType int

const (
	StackTypeDockerCompose StackType = 1
	StackTypeDockerSwarm   StackType = 2
)

type Stack struct {
	ID              int              `json:"Id"`
	Name            string           `json:"Name"`
	Type            StackType        `json:"Type"`
	EndpointID      int              `json:"EndpointId"`
	SwarmID         string           `json:"SwarmId,omitempty"`
	EntryPoint      string           `json:"EntryPoint,omitempty"`
	Env             []EnvVar         `json:"Env,omitempty"`
	ResourceControl *ResourceControl `json:"ResourceControl,omitempty"`
	Status          int              `json:"Status"`
	ProjectPath     string           `json:"ProjectPath,omitempty"`
	CreationDate    int64            `json:"CreationDate"`
	CreatedBy       string           `json:"CreatedBy"`
	UpdateDate      int64            `json:"UpdateDate"`
	UpdatedBy       string           `json:"UpdatedBy"`
	AdditionalFiles []string         `json:"AdditionalFiles,omitempty"`
	AutoUpdate      interface{}      `json:"AutoUpdate,omitempty"`
	Option          interface{}      `json:"Option,omitempty"`
	GitConfig       *GitConfig       `json:"GitConfig,omitempty"`
	FromAppTemplate bool             `json:"FromAppTemplate"`
	Namespace       string           `json:"Namespace,omitempty"`
}

type EnvVar struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

type ResourceControl struct {
	ID                 int          `json:"Id"`
	ResourceID         string       `json:"ResourceId"`
	SubResourceIds     []string     `json:"SubResourceIds"`
	Type               int          `json:"Type"`
	UserAccesses       []UserAccess `json:"UserAccesses"`
	TeamAccesses       []TeamAccess `json:"TeamAccesses"`
	Public             bool         `json:"Public"`
	AdministratorsOnly bool         `json:"AdministratorsOnly"`
	System             bool         `json:"System"`
}

type UserAccess struct {
	UserID      int `json:"UserId"`
	AccessLevel int `json:"AccessLevel"`
}

type TeamAccess struct {
	TeamID      int `json:"TeamId"`
	AccessLevel int `json:"AccessLevel"`
}

type GitConfig struct {
	URL            string   `json:"URL,omitempty"`
	ReferenceName  string   `json:"ReferenceName,omitempty"`
	ConfigFilePath string   `json:"ConfigFilePath,omitempty"`
	Authentication *GitAuth `json:"Authentication,omitempty"`
	ConfigHash     string   `json:"ConfigHash,omitempty"`
	TLSSkipVerify  bool     `json:"TLSSkipVerify"`
}

type GitAuth struct {
	Username        string `json:"Username,omitempty"`
	Password        string `json:"Password,omitempty"`
	GitCredentialID int    `json:"GitCredentialID"`
}

type StackFilters struct {
	EndpointID int    `json:"EndpointId,omitempty"`
	SwarmID    string `json:"SwarmId,omitempty"`
}

func (st StackType) String() string {
	switch st {
	case StackTypeDockerCompose:
		return "compose"
	case StackTypeDockerSwarm:
		return "swarm"
	default:
		return "unknown"
	}
}

func (s Stack) StatusString() string {
	switch s.Status {
	case 1:
		return "running"
	case 2:
		return "stopped"
	case 3:
		return "failed"
	default:
		return "unknown"
	}
}

type Pair struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

type AutoUpdateSettings struct {
	Interval       string `json:"interval,omitempty"`
	Webhook        string `json:"webhook,omitempty"`
	ForcePullImage bool   `json:"forcePullImage,omitempty"`
	ForceUpdate    bool   `json:"forceUpdate,omitempty"`
	JobID          string `json:"jobID,omitempty"`
}

type StackCreateSwarmGitPayload struct {
	Name                     string              `json:"name"`
	RepositoryURL            string              `json:"repositoryURL"`
	SwarmID                  string              `json:"swarmID"`
	ComposeFile              string              `json:"composeFile,omitempty"`
	RepositoryReferenceName  string              `json:"repositoryReferenceName,omitempty"`
	RepositoryAuthentication bool                `json:"repositoryAuthentication,omitempty"`
	RepositoryUsername       string              `json:"repositoryUsername,omitempty"`
	RepositoryPassword       string              `json:"repositoryPassword,omitempty"`
	Env                      []Pair              `json:"env,omitempty"`
	AdditionalFiles          []string            `json:"additionalFiles,omitempty"`
	AutoUpdate               *AutoUpdateSettings `json:"autoUpdate,omitempty"`
	TLSSkipVerify            bool                `json:"tlsskipVerify,omitempty"`
	FromAppTemplate          bool                `json:"fromAppTemplate,omitempty"`
}
