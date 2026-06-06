package models

import "time"

// Enums

type Provider string

const (
	ProviderLocal    Provider = "LOCAL"
	ProviderGitHub   Provider = "GITHUB"
	ProviderGoogle   Provider = "GOOGLE"
	ProviderFarcaster Provider = "FARCASTER"
	ProviderWallet   Provider = "WALLET"
	ProviderGeneric  Provider = "GENERIC"
)

type Role string

const (
	RoleSuperAdmin Role = "SUPERADMIN"
	RoleAdmin      Role = "ADMIN"
	RoleUser       Role = "USER"
)

type State string

const (
	StateQueue     State = "QUEUE"
	StatePublished State = "PUBLISHED"
	StateError     State = "ERROR"
	StateDraft     State = "DRAFT"
)

type SubscriptionTier string

const (
	TierStandard SubscriptionTier = "STANDARD"
	TierPro      SubscriptionTier = "PRO"
	TierTeam     SubscriptionTier = "TEAM"
	TierUltimate SubscriptionTier = "ULTIMATE"
)

type Period string

const (
	PeriodMonthly Period = "MONTHLY"
	PeriodYearly  Period = "YEARLY"
)

type ShortLinkPreference string

const (
	ShortLinkAsk ShortLinkPreference = "ASK"
	ShortLinkYes ShortLinkPreference = "YES"
	ShortLinkNo  ShortLinkPreference = "NO"
)

type OrderStatus string

const (
	OrderPending   OrderStatus = "PENDING"
	OrderAccepted  OrderStatus = "ACCEPTED"
	OrderCanceled  OrderStatus = "CANCELED"
	OrderCompleted OrderStatus = "COMPLETED"
)

type CreationMethod string

const (
	CreationUnknown  CreationMethod = "UNKNOWN"
	CreationWeb      CreationMethod = "WEB"
	CreationMCP      CreationMethod = "MCP"
	CreationAPI      CreationMethod = "API"
	CreationAutoPost CreationMethod = "AUTOPOST"
	CreationCLI      CreationMethod = "CLI"
)

type AnnouncementColor string

const (
	ColorInfo    AnnouncementColor = "INFO"
	ColorWarning AnnouncementColor = "WARNING"
	ColorError   AnnouncementColor = "ERROR"
)

type ApprovedSubmitForOrder string

const (
	ApprovedNo                  ApprovedSubmitForOrder = "NO"
	ApprovedWaitingConfirmation ApprovedSubmitForOrder = "WAITING_CONFIRMATION"
	ApprovedYes                 ApprovedSubmitForOrder = "YES"
)

// Models

type Organization struct {
	ID          string              `json:"id" db:"id"`
	Name        string              `json:"name" db:"name"`
	Description *string             `json:"description,omitempty" db:"description"`
	APIKey      *string             `json:"apiKey,omitempty" db:"api_key"`
	PaymentID   *string             `json:"paymentId,omitempty" db:"payment_id"`
	StreakSince  *time.Time          `json:"streakSince,omitempty" db:"streak_since"`
	ShortLink   ShortLinkPreference `json:"shortlink" db:"shortlink"`
	AllowTrial  bool                `json:"allowTrial" db:"allow_trial"`
	IsTrialing  bool                `json:"isTrailing" db:"is_trailing"`
	CreatedAt   time.Time           `json:"createdAt" db:"created_at"`
	UpdatedAt   time.Time           `json:"updatedAt" db:"updated_at"`
}

type User struct {
	ID                    string    `json:"id" db:"id"`
	Email                 string    `json:"email" db:"email"`
	Password              *string   `json:"-" db:"password"`
	ProviderName          Provider  `json:"providerName" db:"provider_name"`
	Name                  *string   `json:"name,omitempty" db:"name"`
	LastName              *string   `json:"lastName,omitempty" db:"last_name"`
	IsSuperAdmin          bool      `json:"isSuperAdmin" db:"is_super_admin"`
	Bio                   *string   `json:"bio,omitempty" db:"bio"`
	Audience              int       `json:"audience" db:"audience"`
	PictureID             *string   `json:"pictureId,omitempty" db:"picture_id"`
	ProviderID            *string   `json:"providerId,omitempty" db:"provider_id"`
	Timezone              int       `json:"timezone" db:"timezone"`
	LastReadNotifications time.Time `json:"lastReadNotifications" db:"last_read_notifications"`
	InviteID              *string   `json:"inviteId,omitempty" db:"invite_id"`
	Activated             bool      `json:"activated" db:"activated"`
	Account               *string   `json:"account,omitempty" db:"account"`
	ConnectedAccount      bool      `json:"connectedAccount" db:"connected_account"`
	LastOnline            time.Time `json:"lastOnline" db:"last_online"`
	SendSuccessEmails     bool      `json:"sendSuccessEmails" db:"send_success_emails"`
	SendFailureEmails     bool      `json:"sendFailureEmails" db:"send_failure_emails"`
	SendStreakEmails       bool      `json:"sendStreakEmails" db:"send_streak_emails"`
	CreatedAt             time.Time `json:"createdAt" db:"created_at"`
	UpdatedAt             time.Time `json:"updatedAt" db:"updated_at"`
}

type UserOrganization struct {
	ID             string    `json:"id" db:"id"`
	UserID         string    `json:"userId" db:"user_id"`
	OrganizationID string    `json:"organizationId" db:"organization_id"`
	Disabled       bool      `json:"disabled" db:"disabled"`
	Role           Role      `json:"role" db:"role"`
	CreatedAt      time.Time `json:"createdAt" db:"created_at"`
	UpdatedAt      time.Time `json:"updatedAt" db:"updated_at"`
}

type Post struct {
	ID                         string                 `json:"id" db:"id"`
	State                      State                  `json:"state" db:"state"`
	PublishDate                time.Time              `json:"publishDate" db:"publish_date"`
	OrganizationID             string                 `json:"organizationId" db:"organization_id"`
	IntegrationID              string                 `json:"integrationId" db:"integration_id"`
	Content                    string                 `json:"content" db:"content"`
	Delay                      int                    `json:"delay" db:"delay"`
	Group                      string                 `json:"group" db:"group"`
	Title                      *string                `json:"title,omitempty" db:"title"`
	Description                *string                `json:"description,omitempty" db:"description"`
	ParentPostID               *string                `json:"parentPostId,omitempty" db:"parent_post_id"`
	ReleaseID                  *string                `json:"releaseId,omitempty" db:"release_id"`
	ReleaseURL                 *string                `json:"releaseUrl,omitempty" db:"release_url"`
	Settings                   *string                `json:"settings,omitempty" db:"settings"`
	Image                      *string                `json:"image,omitempty" db:"image"`
	SubmittedForOrderID        *string                `json:"submittedForOrderId,omitempty" db:"submitted_for_order_id"`
	SubmittedForOrganizationID *string                `json:"submittedForOrganizationId,omitempty" db:"submitted_for_organization_id"`
	ApprovedSubmitForOrder     ApprovedSubmitForOrder `json:"approvedSubmitForOrder" db:"approved_submit_for_order"`
	CreationMethod             CreationMethod         `json:"creationMethod" db:"creation_method"`
	LastMessageID              *string                `json:"lastMessageId,omitempty" db:"last_message_id"`
	IntervalInDays             *int                   `json:"intervalInDays,omitempty" db:"interval_in_days"`
	Error                      *string                `json:"error,omitempty" db:"error"`
	DeletedAt                  *time.Time             `json:"deletedAt,omitempty" db:"deleted_at"`
	CreatedAt                  time.Time              `json:"createdAt" db:"created_at"`
	UpdatedAt                  time.Time              `json:"updatedAt" db:"updated_at"`
}

type Integration struct {
	ID                   string     `json:"id" db:"id"`
	InternalID           string     `json:"internalId" db:"internal_id"`
	OrganizationID       string     `json:"organizationId" db:"organization_id"`
	Name                 string     `json:"name" db:"name"`
	Picture              *string    `json:"picture,omitempty" db:"picture"`
	ProviderIdentifier   string     `json:"providerIdentifier" db:"provider_identifier"`
	Type                 string     `json:"type" db:"type"`
	Token                string     `json:"-" db:"token"`
	Disabled             bool       `json:"disabled" db:"disabled"`
	TokenExpiration      *time.Time `json:"tokenExpiration,omitempty" db:"token_expiration"`
	RefreshToken         *string    `json:"-" db:"refresh_token"`
	Profile              *string    `json:"profile,omitempty" db:"profile"`
	DeletedAt            *time.Time `json:"deletedAt,omitempty" db:"deleted_at"`
	InBetweenSteps       bool       `json:"inBetweenSteps" db:"in_between_steps"`
	RefreshNeeded        bool       `json:"refreshNeeded" db:"refresh_needed"`
	PostingTimes         string     `json:"postingTimes" db:"posting_times"`
	CustomInstanceDetails *string   `json:"customInstanceDetails,omitempty" db:"custom_instance_details"`
	CustomerID           *string    `json:"customerId,omitempty" db:"customer_id"`
	RootInternalID       *string    `json:"rootInternalId,omitempty" db:"root_internal_id"`
	AdditionalSettings   *string    `json:"additionalSettings,omitempty" db:"additional_settings"`
	CreatedAt            time.Time  `json:"createdAt" db:"created_at"`
	UpdatedAt            *time.Time `json:"updatedAt,omitempty" db:"updated_at"`
}

type Subscription struct {
	ID             string           `json:"id" db:"id"`
	OrganizationID string           `json:"organizationId" db:"organization_id"`
	Tier           SubscriptionTier `json:"subscriptionTier" db:"subscription_tier"`
	Identifier     *string          `json:"identifier,omitempty" db:"identifier"`
	CancelAt       *time.Time       `json:"cancelAt,omitempty" db:"cancel_at"`
	Period         Period           `json:"period" db:"period"`
	TotalChannels  int              `json:"totalChannels" db:"total_channels"`
	IsLifetime     bool             `json:"isLifetime" db:"is_lifetime"`
	DeletedAt      *time.Time       `json:"deletedAt,omitempty" db:"deleted_at"`
	CreatedAt      time.Time        `json:"createdAt" db:"created_at"`
	UpdatedAt      time.Time        `json:"updatedAt" db:"updated_at"`
}

type Media struct {
	ID                string     `json:"id" db:"id"`
	Name              string     `json:"name" db:"name"`
	OriginalName      *string    `json:"originalName,omitempty" db:"original_name"`
	Path              string     `json:"path" db:"path"`
	OrganizationID    string     `json:"organizationId" db:"organization_id"`
	FileSize          int        `json:"fileSize" db:"file_size"`
	Type              string     `json:"type" db:"type"`
	Thumbnail         *string    `json:"thumbnail,omitempty" db:"thumbnail"`
	Alt               *string    `json:"alt,omitempty" db:"alt"`
	ThumbnailTimestamp *int      `json:"thumbnailTimestamp,omitempty" db:"thumbnail_timestamp"`
	DeletedAt         *time.Time `json:"deletedAt,omitempty" db:"deleted_at"`
	CreatedAt         time.Time  `json:"createdAt" db:"created_at"`
	UpdatedAt         time.Time  `json:"updatedAt" db:"updated_at"`
}

type Tag struct {
	ID             string     `json:"id" db:"id"`
	Name           string     `json:"name" db:"name"`
	Color          string     `json:"color" db:"color"`
	OrgID          string     `json:"orgId" db:"org_id"`
	DeletedAt      *time.Time `json:"deletedAt,omitempty" db:"deleted_at"`
	CreatedAt      time.Time  `json:"createdAt" db:"created_at"`
	UpdatedAt      time.Time  `json:"updatedAt" db:"updated_at"`
}

type Notification struct {
	ID             string     `json:"id" db:"id"`
	OrganizationID string     `json:"organizationId" db:"organization_id"`
	Content        string     `json:"content" db:"content"`
	Link           *string    `json:"link,omitempty" db:"link"`
	DeletedAt      *time.Time `json:"deletedAt,omitempty" db:"deleted_at"`
	CreatedAt      time.Time  `json:"createdAt" db:"created_at"`
	UpdatedAt      time.Time  `json:"updatedAt" db:"updated_at"`
}

type Webhook struct {
	ID             string     `json:"id" db:"id"`
	Name           string     `json:"name" db:"name"`
	OrganizationID string     `json:"organizationId" db:"organization_id"`
	URL            string     `json:"url" db:"url"`
	DeletedAt      *time.Time `json:"deletedAt,omitempty" db:"deleted_at"`
	CreatedAt      time.Time  `json:"createdAt" db:"created_at"`
	UpdatedAt      time.Time  `json:"updatedAt" db:"updated_at"`
}

type Set struct {
	ID             string    `json:"id" db:"id"`
	OrganizationID string    `json:"organizationId" db:"organization_id"`
	Name           string    `json:"name" db:"name"`
	Content        string    `json:"content" db:"content"`
	CreatedAt      time.Time `json:"createdAt" db:"created_at"`
	UpdatedAt      time.Time `json:"updatedAt" db:"updated_at"`
}

type Signature struct {
	ID             string     `json:"id" db:"id"`
	OrganizationID string     `json:"organizationId" db:"organization_id"`
	Content        string     `json:"content" db:"content"`
	AutoAdd        bool       `json:"autoAdd" db:"auto_add"`
	DeletedAt      *time.Time `json:"deletedAt,omitempty" db:"deleted_at"`
	CreatedAt      time.Time  `json:"createdAt" db:"created_at"`
	UpdatedAt      time.Time  `json:"updatedAt" db:"updated_at"`
}

type AutoPost struct {
	ID              string     `json:"id" db:"id"`
	OrganizationID  string     `json:"organizationId" db:"organization_id"`
	Title           string     `json:"title" db:"title"`
	Content         *string    `json:"content,omitempty" db:"content"`
	OnSlot          bool       `json:"onSlot" db:"on_slot"`
	SyncLast        bool       `json:"syncLast" db:"sync_last"`
	URL             string     `json:"url" db:"url"`
	LastURL         string     `json:"lastUrl" db:"last_url"`
	Active          bool       `json:"active" db:"active"`
	AddPicture      bool       `json:"addPicture" db:"add_picture"`
	GenerateContent bool       `json:"generateContent" db:"generate_content"`
	Integrations    string     `json:"integrations" db:"integrations"`
	DeletedAt       *time.Time `json:"deletedAt,omitempty" db:"deleted_at"`
	CreatedAt       time.Time  `json:"createdAt" db:"created_at"`
	UpdatedAt       time.Time  `json:"updatedAt" db:"updated_at"`
}

type ThirdParty struct {
	ID             string     `json:"id" db:"id"`
	OrganizationID string     `json:"organizationId" db:"organization_id"`
	Identifier     string     `json:"identifier" db:"identifier"`
	Name           string     `json:"name" db:"name"`
	InternalID     string     `json:"internalId" db:"internal_id"`
	APIKey         string     `json:"apiKey" db:"api_key"`
	DeletedAt      *time.Time `json:"deletedAt,omitempty" db:"deleted_at"`
	CreatedAt      time.Time  `json:"createdAt" db:"created_at"`
	UpdatedAt      time.Time  `json:"updatedAt" db:"updated_at"`
}

type OAuthApp struct {
	ID             string     `json:"id" db:"id"`
	OrganizationID string     `json:"organizationId" db:"organization_id"`
	Name           string     `json:"name" db:"name"`
	Description    *string    `json:"description,omitempty" db:"description"`
	PictureID      *string    `json:"pictureId,omitempty" db:"picture_id"`
	RedirectURL    string     `json:"redirectUrl" db:"redirect_url"`
	ClientID       string     `json:"clientId" db:"client_id"`
	ClientSecret   string     `json:"-" db:"client_secret"`
	DeletedAt      *time.Time `json:"deletedAt,omitempty" db:"deleted_at"`
	CreatedAt      time.Time  `json:"createdAt" db:"created_at"`
	UpdatedAt      time.Time  `json:"updatedAt" db:"updated_at"`
}

type Announcement struct {
	ID          string            `json:"id" db:"id"`
	Title       string            `json:"title" db:"title"`
	Description string            `json:"description" db:"description"`
	Color       AnnouncementColor `json:"color" db:"color"`
	CreatedAt   time.Time         `json:"createdAt" db:"created_at"`
}

type Comment struct {
	ID             string     `json:"id" db:"id"`
	Content        string     `json:"content" db:"content"`
	OrganizationID string     `json:"organizationId" db:"organization_id"`
	PostID         string     `json:"postId" db:"post_id"`
	UserID         string     `json:"userId" db:"user_id"`
	DeletedAt      *time.Time `json:"deletedAt,omitempty" db:"deleted_at"`
	CreatedAt      time.Time  `json:"createdAt" db:"created_at"`
	UpdatedAt      time.Time  `json:"updatedAt" db:"updated_at"`
}

type Credits struct {
	ID             string    `json:"id" db:"id"`
	OrganizationID string    `json:"organizationId" db:"organization_id"`
	Credits        int       `json:"credits" db:"credits"`
	Type           string    `json:"type" db:"type"`
	CreatedAt      time.Time `json:"createdAt" db:"created_at"`
	UpdatedAt      time.Time `json:"updatedAt" db:"updated_at"`
}

type GitHub struct {
	ID             string     `json:"id" db:"id"`
	Login          *string    `json:"login,omitempty" db:"login"`
	Name           *string    `json:"name,omitempty" db:"name"`
	Token          string     `json:"-" db:"token"`
	JobID          *string    `json:"jobId,omitempty" db:"job_id"`
	OrganizationID string     `json:"organizationId" db:"organization_id"`
	CreatedAt      time.Time  `json:"createdAt" db:"created_at"`
	UpdatedAt      time.Time  `json:"updatedAt" db:"updated_at"`
}

type Customer struct {
	ID        string     `json:"id" db:"id"`
	Name      string     `json:"name" db:"name"`
	OrgID     string     `json:"orgId" db:"org_id"`
	DeletedAt *time.Time `json:"deletedAt,omitempty" db:"deleted_at"`
	CreatedAt time.Time  `json:"createdAt" db:"created_at"`
	UpdatedAt time.Time  `json:"updatedAt" db:"updated_at"`
}
