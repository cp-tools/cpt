module github.com/cp-tools/cpt

go 1.14

require (
	github.com/AlecAivazis/survey/v2 v2.1.1
	github.com/cp-tools/cpt-lib v1.5.1
	github.com/fatih/color v1.9.0
	github.com/gosuri/uilive v0.0.4
	github.com/gosuri/uitable v0.0.4
	github.com/inconshreveable/go-update v0.0.0-20160112193335-8152e7eb6ccf
	github.com/infixint943/cookiejar v0.1.0
	github.com/kballard/go-shellquote v0.0.0-20180428030007-95032a82bc51
	github.com/knadh/koanf v0.13.0
	github.com/mitchellh/go-homedir v1.1.0
	github.com/oleiade/serrure v0.0.0-20160812094227-28794589ac9b
	github.com/spf13/cobra v1.0.1-0.20200719220246-c6fe2d4df810
	github.com/spf13/pflag v1.0.5
	github.com/spf13/viper v1.7.1
	github.com/tidwall/gjson v1.6.1
)

replace github.com/spf13/pflag => github.com/infixint943/pflag v1.0.6-0.20200801213445-50927f2f338d

replace github.com/AlecAivazis/survey/v2 => github.com/infixint943/survey/v2 v2.1.2-0.20201001232057-ad85ff5097a6

replace github.com/mitchellh/go-homedir => github.com/infixint943/go-homedir v1.1.1-0.20200627072908-00f1ec2bf896

replace github.com/knadh/koanf => github.com/infixint943/koanf master