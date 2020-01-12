package cmd

import (
	"alfredman/man"
	"log"
	"os"
	"os/exec"

	aw "github.com/deanishe/awgo"
	"github.com/spf13/cobra"
)

var searchCmd = &cobra.Command{
	Use:   "search",
	Short: "A brief description of your application",
	RunE:  searchRun,
}

func searchRun(cmd *cobra.Command, args []string) error {
	query, err := cmd.Flags().GetString("term")
	if err != nil {
		return err
	}

	if wf.UpdateCheckDue() && !wf.IsRunning(updateJobName) {
		log.Println("Running update check in background...")

		cmd := exec.Command(os.Args[0], "-check")
		if err := wf.RunInBackground(updateJobName, cmd); err != nil {
			log.Printf("Error starting update check: %s", err)
		}
	}

	// Only show update status if query is empty.
	if query == "" && wf.UpdateAvailable() {
		wf.Configure(aw.SuppressUIDs(true))

		wf.NewItem("Update available!").
			Subtitle("↩ to install").
			Autocomplete("workflow:update").
			Valid(false).
			Icon(iconAvailable)
	}

	commands := man.NewManInterface()
	if wf.Cache.Exists(cacheName) {
		if err := wf.Cache.LoadJSON(cacheName, &commands.Commands); err != nil {
			wf.FatalError(err)
		}
	}

	if wf.Cache.Expired(cacheName, maxCacheAge) {
		wf.Rerun(0.3)
		if !wf.IsRunning("cache") {
			cmd := exec.Command(os.Args[0], "--cache")
			if err := wf.RunInBackground("cache", cmd); err != nil {
				wf.FatalError(err)
			}
		} else {
			log.Printf("cache job already running.")
		}
		// Cache is also "expired" if it doesn't exist. So if there are no
		// cached data, show a corresponding message and exit.
		if len(commands.Commands) == 0 {
			wf.NewItem("Downloading man pagees").
				Icon(aw.IconInfo)
			wf.SendFeedback()
			return nil
		}

	}

	// Add results for cached repos
	for _, r := range commands.Commands {
		wf.NewItem(r.Name).
			Subtitle(r.Description).
			Arg(r.ManArg).
			UID(r.ManURI).
			Valid(true)
	}

	wf.NewItem("Reset update status").
		Autocomplete("workflow:delcache").
		Icon(aw.IconTrash).
		Valid(false)

	// Filter results against query if user entered one
	if query != "" {
		res := wf.Filter(query)
		log.Printf("[main] %d/%d man pages match \"%s\"", len(res), len(commands.Commands), query)
	}

	wf.WarnEmpty("No man pages found", "Try a different query?")

	// Send results/warning message to Alfred
	wf.SendFeedback()

	return nil
}

func init() {
	searchCmd.Flags().StringP("term", "t", "", "Term to search man pages for")
	rootCmd.AddCommand(searchCmd)
}
