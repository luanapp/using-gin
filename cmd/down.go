/*
Copyright © 2021 Luana Pimentel

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package cmd

import (
	"github.com/luanapp/gin-example/config/database"
	"github.com/spf13/cobra"
)

var (
	downFlags = &Flags{}
	downCmd   = &cobra.Command{
		Use:     "down",
		Short:   "undo changes already applied to the database",
		Long:    ``,
		Example: "migrate down",
		Run: func(cmd *cobra.Command, args []string) {
			if downFlags.filename != "" {
				err := database.DownFromFilename(downFlags.filename)
				cobra.CheckErr(err)
			} else {
				err := database.Down()
				cobra.CheckErr(err)
			}
		},
	}
)

func init() {
	downCmd.Flags().StringVarP(&downFlags.filename, "filename", "f", "", "Filename with the migration to run")
	cobra.CheckErr(downCmd.MarkFlagFilename("filename"))

	rootCmd.AddCommand(downCmd)
}