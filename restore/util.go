package restore

import (
	"flag"
	"fmt"
	cc "github.com/myENA/consul-backinator/common/config"
	ccns "github.com/myENA/consul-backinator/common/consul"
	"os"
	"strings"
)

// setupFlags initializes the instance configuration
func (c *Command) setupFlags(args []string) error {
	var cmdFlags *flag.FlagSet // instance flagset

	// init config if needed
	if config == nil {
		config = new(configStruct)
	}

	// init consul config if needed
	if config.consulConfig == nil {
		config.consulConfig = new(ccns.Config)
	}

	// init flagset
	cmdFlags = flag.NewFlagSet("restore", flag.ContinueOnError)
	cmdFlags.Usage = func() { fmt.Fprint(os.Stdout, c.Help()); os.Exit(0) }

	// declare flags
	cmdFlags.StringVar(&config.fileName, "file", "consul.bak",
		"Source")
	cmdFlags.StringVar(&config.cryptKey, "key", "password",
		"Passphrase for data encryption and signature validation")
	cmdFlags.BoolVar(&config.noKV, "nokv", false,
		"Do not attempt to restore kv data")
	cmdFlags.StringVar(&config.aclFileName, "acls", "",
		"Optional source filename for acl tokens")
	cmdFlags.StringVar(&config.queryFileName, "queries", "",
		"Optional source filename for query definitions")
	cmdFlags.StringVar(&config.pathTransform, "transform", "",
		"Optional path transformation")
	cmdFlags.BoolVar(&config.delTree, "delete", false,
		"Delete all keys under specified prefix")
	cmdFlags.StringVar(&config.consulPrefix, "prefix", "/",
		"Prefix for delete operation")

	// add shared flags
	cc.AddSharedConsulFlags(cmdFlags, config.consulConfig)

	// parse flags and ignore error
	if err := cmdFlags.Parse(args); err != nil {
		return nil
	}

	// populate potentially missing config items
	cc.AddEnvDefaults(config.consulConfig)

	// fixup prefix per upstream issue 2403
	// https://github.com/hashicorp/consul/issues/2403
	config.consulPrefix = strings.TrimPrefix(config.consulPrefix,
		ccns.Separator)

	// always okay
	return nil
}
