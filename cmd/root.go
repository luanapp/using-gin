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
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "migrate",
	Short: "Apply db migrations into database",
	Long: `This is a basic migration command which has to commands, as follows:

up: applies the changes in the field "up" in the migration file
down: undo applied changes in the database, using the field "down" in the migration file

The migration file must be a yaml with the following fields:
name: migration name. This is meant for you! Organize your shit, we won't use this field anywhere
description: same as above
up: sql to use when running the migration
down: sql that undo the changes done in the "up" field`,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	cobra.CheckErr(rootCmd.Execute())
}
