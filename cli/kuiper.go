package cli

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/cloustone/pandas/kuiper"
	mfxsdk "github.com/cloustone/pandas/sdk/go"
	"github.com/spf13/cobra"
)

var createCommand cobra.Command = cobra.Command{
	Use:   "create",
	Short: "create kuiper stream and rule",
	Run:   func(cmd *cobra.Command, args []string) {},
}

var createStreamCommand cobra.Command = cobra.Command{
	Use:   "stream",
	Short: "create stream <stream> [-f stream_def_file] <user_auth_token>",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 3 {
			sfile := args[1]
			stream, err := readDef(sfile, "stream")
			if err != nil {
				logError(err)
				return
			}
			if reply, err := sdk.CreateKuiperStream(string(stream), args[2]); err != nil {
				logError(err)
			} else {
				logJSON(reply)
			}
		} else if len(args) == 2 {
			if reply, err := sdk.CreateKuiperStream(args[0], args[1]); err != nil {
				logError(err)
			} else {
				logJSON(reply)
			}
		} else {
			logUsage("create stream <stream> [-f stream_def_file] <user_auth_token>\n")
		}
	},
}

var createRuleCommand cobra.Command = cobra.Command{
	Use:   "rule",
	Short: "create rule <rule_name> [rule_json | -f rule_def_file] <user_auth_token>",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 4 { // using rule definition file
			sfile := args[2]
			rule, err := readDef(sfile, "rule")
			if err != nil {
				logError(err)
				return
			}
			rname := args[0]
			desc := kuiper.Rule{Name: rname, SQL: string(rule)}
			if reply, err := sdk.CreateKuiperRule(desc, args[3]); err != nil {
				logError(err)
			} else {
				logJSON(reply)
			}
		} else if len(args) == 3 {
			rname := args[0]
			rjson := args[1]
			desc := kuiper.Rule{Name: rname, SQL: rjson}
			if reply, err := sdk.CreateKuiperRule(desc, args[3]); err != nil {
				logError(err)
			} else {
				logJSON(reply)
			}
		} else {
			logUsage("create rule <rule_name> [rule_json | -f rule_def_file] <user_auth_token>\n")
		}
	},
}

var describeCommand = cobra.Command{
	Use:   "describe",
	Short: "describe stream $stream_name | describe rule $rule_name | describe plugin $plugin_type $plugin_name",
	Run:   func(cmd *cobra.Command, args []string) {},
}

var describeStreamCommand cobra.Command = cobra.Command{
	Use:   "stream",
	Short: "describe stream <stream_name> <user_auth_token>",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) != 2 {
			logUsage("describe stream <stream_name> <user_auth_token>.\n")
			return
		}
		if reply, err := sdk.KuiperStream(args[0], args[1]); err != nil {
			logError(err)
		} else {
			logJSON(reply)
		}

	},
}

var describeRuleCommand cobra.Command = cobra.Command{
	Use:   "rule",
	Short: "describe rule <rule_name> <user_auth_token>",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) != 2 {
			logUsage("describe rule <rule_name> <user_auth_token>.\n")
			return
		}
		if reply, err := sdk.KuiperRule(args[0], args[1]); err != nil {
			logError(err)
		} else {
			logJSON(reply)
		}
	},
}

var describePluginCommand cobra.Command = cobra.Command{
	Use:   "plugin",
	Short: "describe plugin <plugin_type> <plugin_name>",
	Run:   func(cmd *cobra.Command, args []string) {},
}

var dropCommand cobra.Command = cobra.Command{
	Use:   "drop",
	Short: "drop stream $stream_name | drop rule $rule_name | drop plugin $plugin_type $plugin_name -r $stop",
	Run:   func(cmd *cobra.Command, args []string) {},
}

var dropStreamCommand cobra.Command = cobra.Command{
	Use:   "stream",
	Short: "drop stream <stream_name> <user_auth_token>",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) != 2 {
			logUsage("drop stream <stream_name> <user_auth_token>")
			return
		}
		if err := sdk.DeleteKuiperStream(args[0], args[1]); err != nil {
			logError(err)
		}
	},
}

var dropRuleCommand cobra.Command = cobra.Command{
	Use:   "rule",
	Short: "drop rule <rule_name> <user_auth_token>",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) != 2 {
			logUsage("drop rule <rule_name> <user_auth_token>")
			return
		}
		if err := sdk.DeleteKuiperRule(args[0], args[1]); err != nil {
			logError(err)
		}
	},
}

var dropPluginCommand cobra.Command = cobra.Command{
	Use:   "plugin",
	Short: "drop plugin <plugin_type> <plugin_name> -s stop",
	Run:   func(cmd *cobra.Command, args []string) {},
}

var showStreamsCommand cobra.Command = cobra.Command{
	Use:   "streams",
	Short: "show streams <user_auth_token>",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) != 1 {
			logUsage("streams <user_auth_token>")
			return
		}
		if reply, err := sdk.KuiperStreams(args[0]); err != nil {
			logError(err)
		} else {
			logJSON(reply)
		}
	},
}

var showRulesCommand cobra.Command = cobra.Command{
	Use:   "rules",
	Short: "rules < user_auth_token>",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) != 1 {
			logUsage("rules <user_auth_token>")
			return
		}
		if reply, err := sdk.KuiperRules(args[0]); err != nil {
			logError(err)
		} else {
			logJSON(reply)
		}
	},
}

var showPluginsCommand cobra.Command = cobra.Command{
	Use:   "plugins",
	Short: "plugins <plugin_type> <user_auth_token>",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) != 2 {
			logUsage("plugins <plugin_type> <user_auth_token>")
			return
		}
		ptype, err := getPluginType(args[0])
		if err != nil {
			logError(err)
			return
		}
		if reply, err := sdk.KuiperPlugins(ptype, args[1]); err != nil {
			logError(err)
		} else {
			logJSON(reply)
		}
	},
}

var statusCommand cobra.Command = cobra.Command{
	Use:   "status",
	Short: "list status rule <rule_id>",
	Run:   func(cmd *cobra.Command, args []string) {},
}

var statusRuleCommand cobra.Command = cobra.Command{
	Use:   "rule",
	Short: "status rule <rule_id> <user_auth_token>",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) != 3 {
			logUsage("status rule <rule_id> <user_auth_token>")
			return
		}
		if reply, err := sdk.KuiperRuleStatus(args[1], args[2]); err != nil {
			logError(err)
		} else {
			logJSON(reply)
		}
	},
}

var startCommand cobra.Command = cobra.Command{
	Use:   "start",
	Short: "start rule <rule_id> <user_auth_token>",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) != 3 {
			logUsage("start rule <rule_id> <user_auth_token>")
			return
		}
		if err := sdk.StartKuiperRule(args[0], args[1]); err != nil {
			logError(err)
		}
	},
}

var stopCommand cobra.Command = cobra.Command{
	Use:   "stop",
	Short: "stop rule <rule_name> <user_auth_token>",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) != 3 {
			logUsage("stop rule <rule_name> <user_auth_token>")
			return
		}
		if err := sdk.StopKuiperRule(args[1], args[2]); err != nil {
			logError(err)
		}
	},
}

var restartCommand cobra.Command = cobra.Command{
	Use:   "restart",
	Short: "restart rule <rule_name> <user_auth_token>",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) != 3 {
			logUsage("restart rule <rule_name> <user_auth_token>")
			return
		}
		if err := sdk.RestartKuiperRule(args[1], args[2]); err != nil {
			logError(err)
		}
	},
}

func getPluginType(arg string) (ptype mfxsdk.KuiperPluginType, err error) {
	switch arg {
	case "source":
		ptype = mfxsdk.KuiperPluginSource
	case "sink":
		ptype = mfxsdk.KuiperPluginSink
	default:
		err = fmt.Errorf("Invalid plugin type %s, should be \"source\", \"sink\" or \"function\".\n", arg)
	}
	return
}

func readDef(sfile string, t string) ([]byte, error) {
	if _, err := os.Stat(sfile); os.IsNotExist(err) {
		return nil, fmt.Errorf("The specified %s defenition file %s is not existed.\n", t, sfile)
	}
	fmt.Printf("Creating a new %s from file %s.\n", t, sfile)
	if rule, err := ioutil.ReadFile(sfile); err != nil {
		return nil, fmt.Errorf("Failed to read from %s definition file %s.\n", t, sfile)
	} else {
		return rule, nil
	}
}

// NewKuiperCmd return rule engin kuiper command.
func NewKuiperCmd() *cobra.Command {
	cmd := cobra.Command{
		Use:   "kuiper",
		Short: "kuiper management",
		Long:  `kuiper management: create, get, update or delete Rules`,
		Run: func(cmd *cobra.Command, args []string) {
			logUsage("kuiper[create | describe| drop| show]")
		},
	}

	// Create
	createCommand.AddCommand(&createStreamCommand)
	createCommand.AddCommand(&createRuleCommand)

	// Describe
	describeCommand.AddCommand(&describeStreamCommand)
	describeCommand.AddCommand(&describeRuleCommand)
	describeCommand.AddCommand(&describePluginCommand)
	cmd.AddCommand(&describeCommand)

	// Drop
	dropCommand.AddCommand(&dropStreamCommand)
	dropCommand.AddCommand(&dropRuleCommand)
	dropCommand.AddCommand(&dropPluginCommand)
	cmd.AddCommand(&dropCommand)

	// Show
	cmd.AddCommand(&showStreamsCommand)
	cmd.AddCommand(&showRulesCommand)
	cmd.AddCommand(&showPluginsCommand)

	// Status
	statusCommand.AddCommand(&statusRuleCommand)
	cmd.AddCommand(&statusCommand)

	return &cmd
}
