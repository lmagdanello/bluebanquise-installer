package system

var PythonRequirements = []string{
	"ansible",
	"ansible-core",
	"netaddr",
	"clustershell",
	"jmespath",
	"jinja2",
	"pymysql",
}

type PackageDefinition struct {
	OSID     string
	Version  string
	Packages []string
	PostHook func() error
}

var DependenciePackages = []PackageDefinition{
	{
		OSID:    "ubuntu",
		Version: "24.04",
		Packages: []string{
			"python3.12", "python3.12-pip", "python3.12-venv",
			"ssh", "curl", "git",
		},
	},
	{
		OSID:    "ubuntu",
		Version: "22.04",
		Packages: []string{
			"python3.12", "python3.12-pip", "python3.12-venv",
			"ssh", "curl", "git",
		},
	},
	{
		OSID:    "ubuntu",
		Version: "20.04",
		Packages: []string{
			"build-essential", "zlib1g-dev", "libncurses5-dev", "libgdbm-dev",
			"libnss3-dev", "libssl-dev", "libreadline-dev", "libffi-dev",
			"libsqlite3-dev", "wget", "libbz2-dev", "pkg-config", "ssh",
			"curl", "git",
		},
		PostHook: BuildPython311FromSource,
	},
	{
		OSID:    "rhel",
		Version: "7",
		Packages: []string{
			"epel-release", "openssh",
			"centos-release-scl-rh", "centos-release-scl", "rh-python38",
		},
	},
	{
		OSID:    "rhel",
		Version: "8",
		Packages: []string{
			"git", "python39", "python3-pip", "python3-policycoreutils", "openssh-clients", "python39-setuptools",
		},
	},
	{
		OSID:    "rhel",
		Version: "9",
		Packages: []string{
			"git", "python3.12", "python3.12-pip", "python3-policycoreutils", "openssh-clients", "python3.12-setuptools",
		},
	},
	{
		OSID:    "debian",
		Version: "11",
		Packages: []string{
			"python3", "python3-pip", "python3-venv", "git", "ssh", "curl",
		},
	},
	{
		OSID:    "debian",
		Version: "12",
		Packages: []string{
			"python3.12", "python3.12-pip", "python3.12-venv", "git", "ssh", "curl",
		},
	},
	{
		OSID:    "opensuse-leap",
		Version: "15.5",
		Packages: []string{
			"python3", "python3-pip", "python311", "python311-pip", "git", "openssh", "curl",
		},
		PostHook: LinkPython311AsDefault,
	},
	{
		OSID:    "opensuse-leap",
		Version: "15.6",
		Packages: []string{
			"python3", "python3-pip", "python311", "python311-pip", "git", "openssh", "curl",
		},
		PostHook: LinkPython311AsDefault,
	},
}
