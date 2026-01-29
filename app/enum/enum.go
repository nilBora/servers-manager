package enum

//go:generate go run github.com/go-pkgz/enum@latest -type serverStatus -lower
type serverStatus int

const (
	ServerStatusActive  serverStatus = iota // enum:alias=active
	ServerStatusPaused                      // enum:alias=paused
	ServerStatusDeleted                     // enum:alias=deleted
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
	ThemeSystem       theme = iota // enum:alias=
	ThemeLight                     // enum:alias=light
	ThemeDark                      // enum:alias=dark
	ThemeDarkElectric              // enum:alias=dark-electric
	ThemeDarkCyber                 // enum:alias=dark-cyber
)

//go:generate go run github.com/go-pkgz/enum@latest -type viewMode -lower
type viewMode int

const (
	ViewModeTable viewMode = iota
	ViewModeCards
)
