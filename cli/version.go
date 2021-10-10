// SPDX-License-Identifier: Apache-2.0

package cli

import "github.com/spf13/cobra"

// NewVersionCmd returns version command.
func NewVersionCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "version",
		Short: "Pandas services version",
		Long:  `Pandas services version: get version of Pandas Things Service`,
		Run: func(cmd *cobra.Command, args []string) {
			v, err := sdk.Version()
			if err != nil {
				logError(err)
				return
			}

			logJSON(v)
		},
	}
}
