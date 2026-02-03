//go:build windows

package app

import (
	"fmt"
	"syscall"

	"github.com/lxn/walk"
	"github.com/lxn/walk/declarative"
	"github.com/lxn/win"

	"shutdown-alert/internal/config"
)

// Win32 message constants
const (
	WM_QUERYENDSESSION = 0x0011
)

// ShellExecute constants
const (
	SW_SHOWNORMAL = 1
)

// appInstance holds the current App instance for WndProc callback.
// This is necessary because Windows callbacks cannot capture Go closures.
var appInstance *App

// origWndProc holds the original window procedure.
var origWndProc uintptr

// App represents the main application.
type App struct {
	mw        *walk.MainWindow
	ni        *walk.NotifyIcon
	urlToOpen string
}

// NewApp creates a new application instance.
func NewApp() *App {
	return &App{
		urlToOpen: config.GetTargetURL(),
	}
}

// Run initializes and runs the application.
func (a *App) Run() error {
	// Store the app instance for WndProc callback
	appInstance = a

	err := a.createMainWindow()
	if err != nil {
		return fmt.Errorf("failed to create main window: %w", err)
	}

	err = a.initNotifyIcon()
	if err != nil {
		return fmt.Errorf("failed to initialize notify icon: %w", err)
	}

	a.installWndProcHook()

	// Start the message loop (blocking call)
	a.mw.Run()

	// Cleanup after the app exits
	if a.ni != nil {
		_ = a.ni.Dispose()
	}

	// Clear the global instance
	appInstance = nil

	return nil
}

// createMainWindow creates a hidden main window.
func (a *App) createMainWindow() error {
	return declarative.MainWindow{
		AssignTo: &a.mw,
		Title:    "Shutdown Alert",
		Visible:  false,
		Layout:   declarative.VBox{},
	}.Create()
}

// wndProcCallback is the custom window procedure to intercept WM_QUERYENDSESSION.
func wndProcCallback(hwnd win.HWND, msg uint32, wParam, lParam uintptr) uintptr {
	if msg == WM_QUERYENDSESSION {
		if appInstance != nil {
			appInstance.showConfirmationDialog()
		}
		// Return TRUE to allow the session to end
		return 1
	}

	// Call the original window procedure
	return win.CallWindowProc(origWndProc, hwnd, msg, wParam, lParam)
}

// installWndProcHook installs the custom window procedure.
func (a *App) installWndProcHook() {
	hwnd := a.mw.Handle()
	origWndProc = win.SetWindowLongPtr(hwnd, win.GWLP_WNDPROC, syscall.NewCallback(wndProcCallback))
}

// initNotifyIcon imperatively creates and configures the notification icon.
func (a *App) initNotifyIcon() error {
	icon, err := walk.NewIconFromFile("internal/icon/icon.ico")
	if err != nil {
		return fmt.Errorf("failed to load icon: %w", err)
	}

	a.ni, err = walk.NewNotifyIcon(a.mw)
	if err != nil {
		return fmt.Errorf("failed to create notify icon: %w", err)
	}

	if err := a.ni.SetIcon(icon); err != nil {
		return fmt.Errorf("failed to set icon: %w", err)
	}
	if err := a.ni.SetToolTip("Shutdown Alert is running."); err != nil {
		return fmt.Errorf("failed to set tooltip: %w", err)
	}

	// Create the exit action
	exitAction := walk.NewAction()
	if err := exitAction.SetText("E&xit"); err != nil {
		return fmt.Errorf("failed to set exit text: %w", err)
	}
	exitAction.Triggered().Attach(func() {
		walk.App().Exit(0)
	})

	// Add exit action to the context menu
	if err := a.ni.ContextMenu().Actions().Add(exitAction); err != nil {
		return fmt.Errorf("failed to add exit action: %w", err)
	}

	return a.ni.SetVisible(true)
}

// showConfirmationDialog displays the shutdown confirmation message.
func (a *App) showConfirmationDialog() {
	msg := fmt.Sprintf("You are about to sign out.\nDo you want to open %s first?", a.urlToOpen)

	result := walk.MsgBox(a.mw, "Shutdown Confirmation", msg, walk.MsgBoxYesNo|walk.MsgBoxIconQuestion)

	if result == walk.DlgCmdYes {
		a.openURL()
	}
}

// openURL opens the target URL using ShellExecute.
func (a *App) openURL() {
	verb, _ := syscall.UTF16PtrFromString("open")
	file, _ := syscall.UTF16PtrFromString(a.urlToOpen)

	win.ShellExecute(a.mw.Handle(), verb, file, nil, nil, SW_SHOWNORMAL)
}
