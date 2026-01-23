package enum

//go:generate go run github.com/go-pkgz/enum@latest -type providerType -lower
type providerType int

const (
	ProviderTypeHetzner  providerType = iota // enum:alias=hetzner
	ProviderTypeAWS                          // enum:alias=aws
	ProviderTypeScaleway                     // enum:alias=scaleway
	ProviderTypeVsysHost                     // enum:alias=vsys_host
)

//go:generate go run github.com/go-pkgz/enum@latest -type accountType -lower
type accountType int

const (
	AccountTypeCloud accountType = iota // enum:alias=cloud
	AccountTypeRobot                    // enum:alias=robot
)

//go:generate go run github.com/go-pkgz/enum@latest -type serverStatus -lower
type serverStatus int

const (
	ServerStatusActive  serverStatus = iota // enum:alias=active
	ServerStatusPaused                      // enum:alias=paused
	ServerStatusDeleted                     // enum:alias=deleted
)

//go:generate go run github.com/go-pkgz/enum@latest -type serverType -lower
type serverType int

const (
	ServerTypeCloud serverType = iota // enum:alias=cloud
	ServerTypeRobot                   // enum:alias=robot
)

//go:generate go run github.com/go-pkgz/enum@latest -type logAction -lower
type logAction int

const (
	LogActionAdded   logAction = iota // enum:alias=added
	LogActionPaused                   // enum:alias=paused
	LogActionDeleted                  // enum:alias=deleted
	LogActionSynced                   // enum:alias=synced
	LogActionUpdated                  // enum:alias=updated
)

//go:generate go run github.com/go-pkgz/enum@latest -type theme -lower
type theme int

const (
	ThemeSystem theme = iota // enum:alias=
	ThemeLight
	ThemeDark
)

//go:generate go run github.com/go-pkgz/enum@latest -type viewMode -lower
type viewMode int

const (
	ViewModeTable viewMode = iota
	ViewModeCards
)
