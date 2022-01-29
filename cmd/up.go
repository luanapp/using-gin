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
	"github.com/spf13/cobra"

	"github.com/luanapp/gin-example/config/database"
)

type (
	Flags struct {
		filename string
	}
)

var (
	upFlags = &Flags{}
	upCmd   = &cobra.Command{
		Use:     "up",
		Short:   "apply changes not yet applied to the database",
		Long:    ``,
		Example: "migrate up",
		Run: func(cmd *cobra.Command, args []string) {
			database.InitializeDB()
			if upFlags.filename != "" {
				err := database.UpFromFilename(upFlags.filename)
				cobra.CheckErr(err)
			} else {
				err := database.Up()
				cobra.CheckErr(err)
			}
		},
	}
)

func init() {
	upCmd.Flags().StringVarP(&upFlags.filename, "filename", "f", "", "Filename with the migration to run")
	cobra.CheckErr(upCmd.MarkFlagFilename("filename"))

	rootCmd.AddCommand(upCmd)
}
