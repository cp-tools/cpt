module github.com/cp-tools/cpt

go 1.14

require (
	github.com/AlecAivazis/survey/v2 v2.2.5
	github.com/cp-tools/cpt-lib/v2 v2.1.3
	github.com/fatih/color v1.10.0
	github.com/gosuri/uilive v0.0.4
	github.com/gosuri/uitable v0.0.4
	github.com/inconshreveable/go-update v0.0.0-20160112193335-8152e7eb6ccf
	github.com/kballard/go-shellquote v0.0.0-20180428030007-95032a82bc51
	github.com/knadh/koanf v0.14.0
	github.com/mattn/go-runewidth v0.0.9 // indirect
	github.com/mitchellh/go-homedir v1.1.0
	github.com/mitchellh/mapstructure v1.4.0
	github.com/spf13/cobra v1.1.1
	github.com/stretchr/testify v1.6.1
	golang.org/x/mod v0.1.0
	gopkg.in/yaml.v2 v2.4.0
)

replace github.com/spf13/pflag => github.com/cornfeedhobo/pflag v1.0.2

replace github.com/mitchellh/go-homedir => github.com/infinitepr0/go-homedir v1.2.0

replace github.com/spf13/cobra => github.com/infinitepr0/cobra v1.0.1-0.20201026004338-69df80ec4c29
