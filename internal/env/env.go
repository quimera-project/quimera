package env

var (
	Config = &config{}
)

type config struct {
	// Arch indicates the mutant arch
	Arch string

	// Version indicates the quimera version
	Version string

	// QuimeraDir indicates the quimera installed directory
	QuimeraDir string

	// Workshop indicates the quimera-workshop directory
	WorkshopDir string

	// PID identifies the current process ID
	PID int

	// Categories indicates the selected categories
	Categories []string

	// Tags indicates the selected tags
	Tags []string

	// OS indicates the selected operating systems
	OS []string

	// Theme indicates the selected theme
	Theme string

	// Silent indicates not to show terminal output
	Silent bool

	// SkipFailed indicates if the user wants to not show failed checks
	SkipFailed bool

	// Output indicates the output directory to store the results
	Output string

	// Markdown indicates to generate markdown output
	Markdown bool

	// Obsidian indicates to generate Obsidian markdown output
	Obsidian bool

	// Html indicates to generate html output
	Html bool

	// Json indicates to generate json output
	Json bool

	// Stdout indicates to generate stdout output
	Stdout bool

	// Raw indicates to generate raw documentation output
	Raw bool

	// Language indicates the selected language
	Language string

	// Language indicates not to ask for user input
	Batch bool

	// Uncolorized indicates not to show terminal output colors
	Uncolorized bool

	// Debug allows to show the checks JSON if an error occurred
	Debug bool
}
