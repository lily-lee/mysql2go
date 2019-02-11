// Copyright © 2018 NAME HERE <EMAIL ADDRESS>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cmd

import (
	"fmt"
	"os"

	"github.com/lily-lee/mysql2go/convert"
	"github.com/spf13/cobra"
)

var cfgFile string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "mysql2go sqlfile [gofile]",
	Short: "mysql2go",
	Long: `Hi~, Welcome to mysql2go.
	
mysql2go is used to convert mysql table structure to go struct.`,
	Run: func(cmd *cobra.Command, args []string) {
		if infile == "" && len(args) > 0 {
			infile = args[0]
		}

		if outfile == "" && len(args) > 1 {
			outfile = args[1]
		}

		convert.Convert(infile, outfile)
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

var infile, outfile string

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.Flags().StringVarP(&infile, "infile", "i", "", "input file path，eg: your sql file.")
	rootCmd.Flags().StringVarP(&outfile, "outfile", "o", "", "output file path. go file.")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
}
