// Moodwave CLI — main entry point.
//
// Command dispatcher for all moodwave subcommands.
// Uses stdlib flag parsing only — no external CLI framework.
package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"

	"github.com/moodwave/moodwave/internal/config"
	"github.com/moodwave/moodwave/internal/platform"
)

const usage = `
MOODWAVE — terminal mood music companion

USAGE
  moodwave <command> [flags]

COMMANDS
  init      Initialize Moodwave config and directories
  scan      Scan the repository and detect mood
  mood      Show the current mood profile
  play      Start playback matched to detected mood
  search    Search YouTube for tracks and play chosen selection
  pause     Pause playback
  stop      Stop playback
  next      Skip to next recommendation
  queue     Show the current music queue
  status    Show current CLI and playback status
  config    View or edit configuration
  theme     Switch visual theme
  visual    Switch visual mode
  source    List or switch music sources
  doctor    Run diagnostics on all subsystems
  update    Update the CLI binary in-place to the latest release

FLAGS
  --help, -h      Show this help message
  --version, -v   Show version information
  --debug         Enable debug logging
  --no-color      Disable ANSI color output
  --no-animation  Disable terminal animations
  --path <dir>    Path to scan (default: current directory)

EXAMPLES
  moodwave init
  moodwave scan
  moodwave play
  moodwave doctor
  moodwave visual spectrum
  moodwave update

DOCS
  https://github.com/moodwave/moodwave/tree/main/docs
`

func main() {
	// Graceful panic recovery to prevent raw stack trace outputs to the user
	defer func() {
		if r := recover(); r != nil {
			fmt.Fprintf(os.Stderr, "\n\033[1;31mError: An unexpected internal error occurred: %v\033[0m\n", r)
			fmt.Fprintf(os.Stderr, "Please report this issue at https://github.com/Boredooms/Moodwave-CLI/issues\n")
			// Restore cursor and screen just in case
			fmt.Print("\033[?25h")
			os.Exit(1)
		}
	}()

	// Dispatch subcommand.
	subcommand := "auto"
	var remainingArgs []string
	var globalArgs []string

	if len(os.Args) >= 2 {
		subcommand = os.Args[1]
		switch subcommand {
		case "--help", "-h", "help":
			fmt.Print(usage)
			return
		case "--version", "-v", "version":
			fmt.Printf("moodwave %s (built %s)\n", config.Version, config.BuildTime)
			return
		}
		remainingArgs = filterFlags(os.Args[2:])
		globalArgs = os.Args[2:]
	}

	// Set up signal handling for clean exits.
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-sigs
		fmt.Fprintf(os.Stderr, "\n\033[?25h") // restore cursor on Ctrl-C
		cancel()
	}()

	// Detect terminal capabilities.
	caps := platform.Detect()

	// Determine the project root for config loading.
	projectRoot, _ := os.Getwd()

	// Parse global flags.
	globalFlags := parseGlobalFlags(globalArgs)

	// Load configuration.
	cfg, err := config.Load(projectRoot)
	if err != nil {
		fatalf("config error: %v", err)
	}

	// Apply global flag overrides to config.
	if globalFlags.noColor || !caps.HasColor {
		cfg.Visual.NoColor = true
	}
	if globalFlags.noAnimation || caps.ReducedMotion {
		cfg.Visual.NoAnimation = true
	}
	if globalFlags.debug {
		cfg.Debug = true
	}
	if globalFlags.path != "" {
		projectRoot = globalFlags.path
	}

	// Ensure directories exist.
	if err := config.EnsureDirectories(cfg); err != nil {
		if cfg.Debug {
			fmt.Fprintf(os.Stderr, "warn: %v\n", err)
		}
	}

	app := &App{
		ctx:         ctx,
		cfg:         cfg,
		caps:        caps,
		projectRoot: projectRoot,
		debug:       globalFlags.debug,
	}

	var cmdErr error
	switch subcommand {
	case "auto":
		if app.caps.IsTTY {
			cmdErr = app.cmdWelcome()
		} else {
			fmt.Println("🌊 MOODWAVE — Autonomous Developer Music companion")
			fmt.Println("  [1/2] Scanning codebase...")
			if err := app.cmdScan(nil); err != nil {
				fatalf("scan failed: %v", err)
			}
			fmt.Println("\n  [2/2] Searching matching tracks...")
			cmdErr = app.cmdPlay(nil)
		}
	case "init":
		cmdErr = app.cmdInit()
	case "scan":
		cmdErr = app.cmdScan(remainingArgs)
	case "mood":
		cmdErr = app.cmdMood()
	case "play":
		cmdErr = app.cmdPlay(remainingArgs)
	case "search":
		cmdErr = app.cmdSearch(remainingArgs)
	case "pause":
		cmdErr = app.cmdPause()
	case "stop":
		cmdErr = app.cmdStop()
	case "next":
		cmdErr = app.cmdNext()
	case "queue":
		cmdErr = app.cmdQueue()
	case "status":
		cmdErr = app.cmdStatus()
	case "config":
		cmdErr = app.cmdConfig(remainingArgs)
	case "theme":
		cmdErr = app.cmdTheme(remainingArgs)
	case "visual":
		cmdErr = app.cmdVisual(remainingArgs)
	case "source":
		cmdErr = app.cmdSource(remainingArgs)
	case "doctor":
		cmdErr = app.cmdDoctor()
	case "update", "upgrade":
		cmdErr = app.cmdSelfUpdate()
	default:
		fmt.Fprintf(os.Stderr, "moodwave: unknown command %q\n\nRun 'moodwave --help' for usage.\n", subcommand)
		os.Exit(1)
	}

	if cmdErr != nil {
		fatalf("%v", cmdErr)
	}
}

// globalFlags holds parsed top-level flags.
type globalFlags struct {
	noColor     bool
	noAnimation bool
	debug       bool
	path        string
}

// parseGlobalFlags scans the remaining args for known global flags.
func parseGlobalFlags(args []string) globalFlags {
	var f globalFlags
	for i := 0; i < len(args); i++ {
		switch args[i] {
		case "--no-color":
			f.noColor = true
		case "--no-animation":
			f.noAnimation = true
		case "--debug":
			f.debug = true
		case "--path":
			if i+1 < len(args) {
				f.path = args[i+1]
				i++
			}
		}
	}
	return f
}

// filterFlags removes flag arguments (--flag or --flag value), leaving only positional args.
func filterFlags(args []string) []string {
	var out []string
	skip := false
	for _, a := range args {
		if skip {
			skip = false
			continue
		}
		if len(a) >= 2 && a[:2] == "--" {
			if a == "--path" {
				skip = true // next arg is the value
			}
			continue
		}
		out = append(out, a)
	}
	return out
}

// fatalf prints an error and exits.
func fatalf(format string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, "\033[31merror:\033[0m "+format+"\n", args...)
	os.Exit(1)
}

// runDefaultScreen shows a brief idle screen when no command is given.
func runDefaultScreen() {
	caps := platform.Detect()

	if caps.HasColor {
		fmt.Print("\033[1m\033[97m")
	}
	fmt.Print(`
  ███╗   ███╗ ██████╗  ██████╗ ██████╗ ██╗    ██╗ █████╗ ██╗   ██╗███████╗
  ████╗ ████║██╔═══██╗██╔═══██╗██╔══██╗██║    ██║██╔══██╗██║   ██║██╔════╝
  ██╔████╔██║██║   ██║██║   ██║██║  ██║██║ █╗ ██║███████║██║   ██║█████╗  
  ██║╚██╔╝██║██║   ██║██║   ██║██║  ██║██║███╗██║██╔══██║╚██╗ ██╔╝██╔══╝  
  ██║ ╚═╝ ██║╚██████╔╝╚██████╔╝██████╔╝╚███╔███╔╝██║  ██║ ╚████╔╝ ███████╗
  ╚═╝     ╚═╝ ╚═════╝  ╚═════╝ ╚═════╝  ╚══╝╚══╝ ╚═╝  ╚═╝  ╚═══╝  ╚══════╝
`)
	if caps.HasColor {
		fmt.Print("\033[0m")
	}

	fmt.Println("  terminal mood music companion")
	fmt.Println()
	fmt.Println("  Run 'moodwave --help' for usage.")
	fmt.Printf("  Version: %s\n\n", config.Version)
}

// ──────────────────────────────────────────────────────────────────────────────
// App is the command runner — holds shared state across subcommands.
// ──────────────────────────────────────────────────────────────────────────────

// App holds shared CLI state.
type App struct {
	ctx         context.Context
	cfg         *config.Config
	caps        platform.Capabilities
	projectRoot string
	debug       bool
}

// sessionFile returns the path to the session state file.
func (a *App) sessionFile() string {
	return filepath.Join(a.cfg.Paths.CacheDir, "session.json")
}

func (a *App) debugf(format string, args ...interface{}) {
	if a.debug {
		fmt.Fprintf(os.Stderr, "[debug] "+format+"\n", args...)
	}
}
