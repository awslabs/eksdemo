package application

type Installer interface {
	Install(options Options) error
	Uninstall(options Options) error
	SetDryRun()
	Type() InstallerType
}

type InstallerType int

const (
	HelmInstaller InstallerType = iota
	ManifestInstaller
)
