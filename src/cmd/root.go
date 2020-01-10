/*
Copyright © 2020 NAME HERE <EMAIL ADDRESS>

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
  "alfredman/man"
  "fmt"
  "log"
  "os"
  "time"

  aw "github.com/deanishe/awgo"
  "github.com/deanishe/awgo/update"
  "github.com/spf13/cobra"
)

var (
  cfgFile string

  // Command-line arguments
  doCheck    bool
  doDownload bool
  query      string

  cacheName   = "man_pages.json"  // Filename of cached repo list
  maxResults  = 200               // Number of results sent to Alfred
  minScore    = 10.0              // Minimum score for a result
  maxCacheAge = 180 * time.Minute // How long to cache repo list for

  // Icon to show if an update is available
  iconAvailable = &aw.Icon{Value: "update-available.png"}
  repo          = "theseceng/alfred-man" // GitHub repo
  wf            *aw.Workflow             // Our Workflow struct
)

// Name of the background job that checks for updates
const updateJobName = "checkForUpdate"

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
  Use:   "alfredman",
  Short: "Alfred workflow to interface with MAN pages",
  RunE:  alfredmanRun,
}

func alfredmanRun(cmd *cobra.Command, args []string) error {
  if cmd.Flags().Changed("update") {
    wf.Configure(aw.TextErrors(true))
    log.Println("Checking for updates...")
    if err := wf.CheckForUpdate(); err != nil {
      wf.FatalError(err)
    }
    return nil
  }

  if cmd.Flags().Changed("cache") {
    wf.Configure(aw.TextErrors(true))
    log.Printf("[main] retrieving man pages list...")
    man := man.NewManInterface()
    err := man.GetManDatabase()
    if err != nil {
      return err
    }

    manDatabase := man.Commands
    if err != nil {
      wf.FatalError(err)
    }
    if err := wf.Cache.StoreJSON(cacheName, manDatabase); err != nil {
      wf.FatalError(err)
    }
    log.Printf("[main] downloaded man pages database")
    return nil
  }
  return nil
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
  if err := rootCmd.Execute(); err != nil {
    fmt.Println(err)
    os.Exit(1)
  }
}

func init() {
  wf = aw.New(update.GitHub(repo))

  // Cobra also supports local flags, which will only run
  // when this action is called directly.
  rootCmd.Flags().BoolP("update", "u", false, "Check for workflow updates")
  rootCmd.Flags().BoolP("cache", "c", false, "Refresh alfred cache")
  rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
