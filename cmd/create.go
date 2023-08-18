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
	"fmt"
	"os"
	"text/template"
	"time"

	"github.com/spf13/cobra"
)

type (
	createFlags struct {
		MigrationName string
	}
	createDatabaseFlags struct {
		TableName     string
		MigrationName string
	}
)

var (
	cFlags    = &createFlags{}
	createCmd = &cobra.Command{
		Use:     "create",
		Short:   "create an empty database migration file",
		Long:    ``,
		Example: "migrate create create-species",
		Run: func(cmd *cobra.Command, args []string) {
			generate("./cmd/templates/migration.yml", cFlags)
		},
	}

	dbFlags        = &createDatabaseFlags{}
	createTableCmd = &cobra.Command{
		Use:     "table",
		Short:   "create a migration file to create a new table",
		Long:    ``,
		Example: "migrate create table create-species",
		Run: func(cmd *cobra.Command, args []string) {
			dbFlags.MigrationName = fmt.Sprintf("create-%s", dbFlags.TableName)
			generate("./cmd/templates/create-table.yml", dbFlags)
		},
	}
)

func generate(tmplFile string, data any) {
	tmpl, err := os.ReadFile(tmplFile)
	cobra.CheckErr(err)

	parsed, err := template.New("migration").Parse(string(tmpl))
	cobra.CheckErr(err)

	currentTime := time.Now().Format("20060102150405")
	filename := fmt.Sprintf("%s-%s.yml", currentTime, getMigrationName(data))

	file, err := os.Create(fmt.Sprintf("%s/%s", "./migrations", filename))
	cobra.CheckErr(err)

	err = parsed.Execute(file, data)
	cobra.CheckErr(err)
}

func getMigrationName(data any) string {
	switch t := data.(type) {
	case *createFlags:
		return t.MigrationName
	case *createDatabaseFlags:
		return fmt.Sprintf("create-%s", t.TableName)
	default:
		return "empty-migration"
	}
}

func init() {
	createCmd.Flags().StringVarP(&cFlags.MigrationName, "name", "n", "", "Migration name")
	cobra.CheckErr(createCmd.MarkFlagRequired("name"))
	createCmd.AddCommand(createTableCmd)

	createTableCmd.Flags().StringVarP(&dbFlags.TableName, "table", "t", "", "Database table name")

	rootCmd.AddCommand(createCmd)
}
