package github

import "time"

type PullRequestEventModel struct {
	Action      string `json:"action"`
	Number      int    `json:"number"`
	PullRequest struct {
		URL                string    `json:"url"`
		ID                 int64     `json:"id"`
		NodeID             string    `json:"node_id"`
		HTMLURL            string    `json:"html_url"`
		DiffURL            string    `json:"diff_url"`
		PatchURL           string    `json:"patch_url"`
		IssueURL           string    `json:"issue_url"`
		Number             int       `json:"number"`
		State              string    `json:"state"`
		Locked             bool      `json:"locked"`
		Title              string    `json:"title"`
		Body               any       `json:"body"`
		CreatedAt          time.Time `json:"created_at"`
		UpdatedAt          time.Time `json:"updated_at"`
		ClosedAt           any       `json:"closed_at"`
		MergedAt           any       `json:"merged_at"`
		MergeCommitSha     any       `json:"merge_commit_sha"`
		Assignee           any       `json:"assignee"`
		Assignees          []any     `json:"assignees"`
		RequestedReviewers []any     `json:"requested_reviewers"`
		RequestedTeams     []any     `json:"requested_teams"`
		Labels             []any     `json:"labels"`
		Milestone          any       `json:"milestone"`
		Draft              bool      `json:"draft"`
		CommitsURL         string    `json:"commits_url"`
		ReviewCommentsURL  string    `json:"review_comments_url"`
		ReviewCommentURL   string    `json:"review_comment_url"`
		CommentsURL        string    `json:"comments_url"`
		StatusesURL        string    `json:"statuses_url"`
		Head               struct {
			Label string `json:"label"`
			Ref   string `json:"ref"`
			Sha   string `json:"sha"`
		} `json:"head"`
		AuthorAssociation   string `json:"author_association"`
		AutoMerge           any    `json:"auto_merge"`
		ActiveLockReason    any    `json:"active_lock_reason"`
		Merged              bool   `json:"merged"`
		Mergeable           any    `json:"mergeable"`
		Rebaseable          any    `json:"rebaseable"`
		MergeableState      string `json:"mergeable_state"`
		MergedBy            any    `json:"merged_by"`
		Comments            int    `json:"comments"`
		ReviewComments      int    `json:"review_comments"`
		MaintainerCanModify bool   `json:"maintainer_can_modify"`
		Commits             int    `json:"commits"`
		Additions           int    `json:"additions"`
		Deletions           int    `json:"deletions"`
		ChangedFiles        int    `json:"changed_files"`
	} `json:"pull_request"`
	Repository struct {
		ID    int    `json:"id"`
		Name  string `json:"name"`
		Owner struct {
			Login string `json:"login"`
		} `json:"owner"`
	} `json:"repository"`
	Sender struct {
		Login string `json:"login"`
	} `json:"sender"`
}

type InstallationsResponseModel struct {
	ID int `json:"id"`
}

type InstallationTokenResponseModel struct {
	Token     string    `json:"token"`
	ExpiresAt time.Time `json:"expires_at"`
}

type PRMetadataResponseModel struct {
	URL    string `json:"url"`
	ID     int    `json:"id"`
	Number int    `json:"number"`
	Title  string `json:"title"`
	User   struct {
		Login string `json:"login"`
	} `json:"user"`
	Body string `json:"body"`
	Head struct {
		Label string `json:"label"`
		Ref   string `json:"ref"`
		Sha   string `json:"sha"`
	} `json:"head"`
	Base struct {
		Label string `json:"label"`
		Ref   string `json:"ref"`
		Sha   string `json:"sha"`
	} `json:"base"`
	Comments            int  `json:"comments"`
	ReviewComments      int  `json:"review_comments"`
	MaintainerCanModify bool `json:"maintainer_can_modify"`
	Commits             int  `json:"commits"`
	Additions           int  `json:"additions"`
	Deletions           int  `json:"deletions"`
	ChangedFiles        int  `json:"changed_files"`
}

type PRFileChangesResponseModel []struct {
	Sha         string `json:"sha"`
	Filename    string `json:"filename"`
	Status      string `json:"status"`
	Additions   int    `json:"additions"`
	Deletions   int    `json:"deletions"`
	Changes     int    `json:"changes"`
	BlobURL     string `json:"blob_url"`
	RawURL      string `json:"raw_url"`
	ContentsURL string `json:"contents_url"`
	Patch       string `json:"patch"`
}

type PRReviewRequestModel struct {
	CommitID string    `json:"commit_id"`
	Body     string    `json:"body"`
	Event    string    `json:"event"`
	Comments []Comment `json:"comments"`
}

type Comment struct {
	Path     string `json:"path"`
	Position int64  `json:"position"`
	Body     string `json:"body"`
}

type PRReviewResponseModel struct {
	ID                int64     `json:"id"`
	NodeID            string    `json:"node_id"`
	User              User      `json:"user"`
	Body              string    `json:"body"`
	State             string    `json:"state"`
	HTMLURL           string    `json:"html_url"`
	PullRequestURL    string    `json:"pull_request_url"`
	Links             Links     `json:"_links"`
	SubmittedAt       time.Time `json:"submitted_at"`
	CommitID          string    `json:"commit_id"`
	AuthorAssociation string    `json:"author_association"`
}

type Links struct {
	HTML        HTML `json:"html"`
	PullRequest HTML `json:"pull_request"`
}

type HTML struct {
	Href string `json:"href"`
}

type User struct {
	Login             string `json:"login"`
	ID                int64  `json:"id"`
	NodeID            string `json:"node_id"`
	AvatarURL         string `json:"avatar_url"`
	GravatarID        string `json:"gravatar_id"`
	URL               string `json:"url"`
	HTMLURL           string `json:"html_url"`
	FollowersURL      string `json:"followers_url"`
	FollowingURL      string `json:"following_url"`
	GistsURL          string `json:"gists_url"`
	StarredURL        string `json:"starred_url"`
	SubscriptionsURL  string `json:"subscriptions_url"`
	OrganizationsURL  string `json:"organizations_url"`
	ReposURL          string `json:"repos_url"`
	EventsURL         string `json:"events_url"`
	ReceivedEventsURL string `json:"received_events_url"`
	Type              string `json:"type"`
	SiteAdmin         bool   `json:"site_admin"`
}

type DiffChunk struct {
	FilePath      string
	CleanedCode   string
	OriginalDiff  string
	HunkHeader    string
	HunkStartLine int
	Position      int
}
