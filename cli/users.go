// SPDX-License-Identifier: Apache-2.0

package cli

import (
	"encoding/json"

	mfxsdk "github.com/cloustone/pandas/sdk/go"
	"github.com/spf13/cobra"
)

var cmdUsers = []cobra.Command{
	cobra.Command{
		Use:   "create",
		Short: "create <username> <password>",
		Long:  `Creates new user`,
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) != 2 {
				logUsage(cmd.Short)
				return
			}

			user := mfxsdk.User{
				Email:    args[0],
				Password: args[1],
			}
			if err := sdk.CreateUser(user); err != nil {
				logError(err)
				return
			}

			logOK()
		},
	},
	cobra.Command{
		Use:   "get",
		Short: "get <user_auth_token>",
		Long:  `Returns user object`,
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) != 1 {
				logUsage(cmd.Short)
				return
			}

			u, err := sdk.User(args[0])
			if err != nil {
				logError(err)
				return
			}

			logJSON(u)
		},
	},
	cobra.Command{
		Use:   "token",
		Short: "token <username> <password>",
		Long:  `Creates new token`,
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) != 2 {
				logUsage(cmd.Short)
				return
			}

			user := mfxsdk.User{
				Email:    args[0],
				Password: args[1],
			}
			token, err := sdk.CreateToken(user)
			if err != nil {
				logError(err)
				return
			}

			logCreated(token)
		},
	},
	cobra.Command{
		Use:   "update",
		Short: "update <JSON_string> <user_auth_token>",
		Long:  `Update user metadata`,
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) != 2 {
				logUsage(cmd.Short)
				return
			}

			var user mfxsdk.User
			if err := json.Unmarshal([]byte(args[0]), &user.Metadata); err != nil {
				logError(err)
				return
			}

			if err := sdk.UpdateUser(user, args[1]); err != nil {
				logError(err)
				return
			}

			logOK()
		},
	},
	cobra.Command{
		Use:   "password",
		Short: "password <old_password> <password> <user_auth_token>",
		Long:  `Update user password`,
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) != 3 {
				logUsage(cmd.Short)
				return
			}

			if err := sdk.UpdatePassword(args[0], args[1], args[2]); err != nil {
				logError(err)
				return
			}

			logOK()
		},
	},
}

// NewUsersCmd returns users command.
func NewUsersCmd() *cobra.Command {
	cmd := cobra.Command{
		Use:   "users",
		Short: "Users management",
		Long:  `Users management: create accounts and tokens"`,
		Run: func(cmd *cobra.Command, args []string) {
			logUsage("Usage: users [create | get | update | token | password]")
		},
	}

	for i := range cmdUsers {
		cmd.AddCommand(&cmdUsers[i])
	}

	return &cmd
}
