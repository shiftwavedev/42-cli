package cmd

// Common API structures for 42 Intranet API

// User represents a 42 user with basic information
type User struct {
	ID    int    `json:"id"`
	Login string `json:"login"`
	URL   string `json:"url"`
}

// UserProfile represents a complete user profile from the API
type UserProfile struct {
	Login           string `json:"login"`
	FirstName       string `json:"first_name"`
	LastName        string `json:"last_name"`
	Email           string `json:"email"`
	Phone           string `json:"phone"`
	DisplayName     string `json:"displayname"`
	Staff           bool   `json:"staff?"`
	CorrectionPoint int    `json:"correction_point"`
	PoolMonth       string `json:"pool_month"`
	PoolYear        string `json:"pool_year"`
	Location        string `json:"location"`
	Wallet          int    `json:"wallet"`
	CreatedAt       string `json:"created_at"`
	UpdatedAt       string `json:"updated_at"`
	Alumni          bool   `json:"alumni?"`
	Active          bool   `json:"active?"`
	Campus          []struct {
		Name    string `json:"name"`
		Country string `json:"country"`
	} `json:"campus"`
	CursusUsers []struct {
		Grade  string  `json:"grade"`
		Level  float64 `json:"level"`
		Cursus struct {
			Name string `json:"name"`
		} `json:"cursus"`
	} `json:"cursus_users"`
}

// Project represents a project structure
type Project struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
	Slug string `json:"slug"`
}

// ProjectUser represents a user's project
type ProjectUser struct {
	ID        int     `json:"id"`
	Name      string  `json:"name"`
	Slug      string  `json:"slug"`
	Status    string  `json:"status"`
	Validated bool    `json:"validated?"`
	FinalMark int     `json:"final_mark"`
	Project   Project `json:"project"`
	CursusIds []int   `json:"cursus_ids"`
	MarkedAt  string  `json:"marked_at"`
	Marked    bool    `json:"marked"`
	Retriable bool    `json:"retriable"`
	Teams     []Team  `json:"teams"`
}

// Team represents a project team
type Team struct {
	ID      int    `json:"id"`
	Name    string `json:"name"`
	RepoURL string `json:"repo_url"`
	Users   []struct {
		ID    int    `json:"id"`
		Login string `json:"login"`
	} `json:"users"`
}

// TeamUser represents a user in a team
type TeamUser struct {
	ID             int    `json:"id"`
	Login          string `json:"login"`
	URL            string `json:"url"`
	Leader         bool   `json:"leader"`
	Occurrence     int    `json:"occurrence"`
	Validated      bool   `json:"validated"`
	ProjectsUserID int    `json:"projects_user_id"`
}

// TeamInfo represents detailed team information
type TeamInfo struct {
	ID                int        `json:"id"`
	Name              string     `json:"name"`
	URL               string     `json:"url"`
	FinalMark         *int       `json:"final_mark"`
	ProjectID         int        `json:"project_id"`
	CreatedAt         string     `json:"created_at"`
	UpdatedAt         string     `json:"updated_at"`
	Status            string     `json:"status"`
	TerminatingAt     *string    `json:"terminating_at"`
	Users             []TeamUser `json:"users"`
	Locked            bool       `json:"locked?"`
	Validated         *bool      `json:"validated?"`
	Closed            bool       `json:"closed?"`
	RepoURL           string     `json:"repo_url"`
	RepoUUID          string     `json:"repo_uuid"`
	LockedAt          string     `json:"locked_at"`
	ClosedAt          string     `json:"closed_at"`
	ProjectSessionID  int        `json:"project_session_id"`
	ProjectGitlabPath string     `json:"project_gitlab_path"`
	Project           Project    `json:"project"`
}

// Flag represents an evaluation flag
type Flag struct {
	ID        int    `json:"id"`
	Name      string `json:"name"`
	Positive  bool   `json:"positive"`
	Icon      string `json:"icon"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

// Language represents a language
type Language struct {
	ID         int    `json:"id"`
	Name       string `json:"name"`
	Identifier string `json:"identifier"`
	CreatedAt  string `json:"created_at"`
	UpdatedAt  string `json:"updated_at"`
}

// Scale represents an evaluation scale
type Scale struct {
	ID                 int        `json:"id"`
	EvaluationID       int        `json:"evaluation_id"`
	Name               string     `json:"name"`
	IsPrimary          bool       `json:"is_primary"`
	Comment            string     `json:"comment"`
	IntroductionMd     string     `json:"introduction_md"`
	DisclaimerMd       string     `json:"disclaimer_md"`
	GuidelinesMd       string     `json:"guidelines_md"`
	CreatedAt          string     `json:"created_at"`
	CorrectionNumber   int        `json:"correction_number"`
	Duration           int        `json:"duration"`
	ManualSubscription bool       `json:"manual_subscription"`
	Languages          []Language `json:"languages"`
	Flags              []Flag     `json:"flags"`
	Free               bool       `json:"free"`
}

// ScaleTeam represents a scale team evaluation
type ScaleTeam struct {
	ID                   int      `json:"id"`
	ScaleID              int      `json:"scale_id"`
	Comment              *string  `json:"comment"`
	CreatedAt            string   `json:"created_at"`
	UpdatedAt            string   `json:"updated_at"`
	Feedback             *string  `json:"feedback"`
	FinalMark            *int     `json:"final_mark"`
	Flag                 Flag     `json:"flag"`
	BeginAt              string   `json:"begin_at"`
	Correcteds           []User   `json:"correcteds"`
	Corrector            User     `json:"corrector"`
	Truant               struct{} `json:"truant"`
	FilledAt             *string  `json:"filled_at"`
	QuestionsWithAnswers []any    `json:"questions_with_answers"`
	Scale                Scale    `json:"scale"`
	Team                 TeamInfo `json:"team"`
	Feedbacks            []any    `json:"feedbacks"`
}

// Location represents user location information
type Location struct {
	ID       int     `json:"id"`
	BeginAt  string  `json:"begin_at"`
	EndAt    *string `json:"end_at"`
	Primary  bool    `json:"primary"`
	Host     string  `json:"host"`
	CampusID int     `json:"campus_id"`
	User     User    `json:"user"`
}

// TokenResponse represents OAuth token response
type TokenResponse struct {
	AccessToken  string `json:"access_token"`
	TokenType    string `json:"token_type"`
	ExpiresIn    int    `json:"expires_in"`
	Scope        string `json:"scope"`
	CreatedAt    int64  `json:"created_at"`
	RefreshToken string `json:"refresh_token"`
}

// CallbackResponse represents OAuth callback response
type CallbackResponse struct {
	Code  string
	Error string
}
