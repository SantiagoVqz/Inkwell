package main

import (
	"fmt"
	"os"
	"strconv"
	"text/tabwriter"
	"time"

	"github.com/santivqzv/inkwell/internal/opml"
	"github.com/santivqzv/inkwell/internal/store"
	"github.com/spf13/cobra"
)

func newFeedsCmd(a *app) *cobra.Command {
	feeds := &cobra.Command{
		Use:   "feeds",
		Short: "Manage feed subscriptions",
	}
	feeds.AddCommand(
		newFeedsAddCmd(a),
		newFeedsImportCmd(a),
		newFeedsListCmd(a),
		newFeedsRemoveCmd(a),
		newFeedsSetActiveCmd(a, true),  // activate
		newFeedsSetActiveCmd(a, false), // deactivate
	)
	return feeds
}

func newFeedsImportCmd(a *app) *cobra.Command {
	return &cobra.Command{
		Use:   "import <file.opml>",
		Short: "Import feeds from an OPML file (e.g. a Feedly export)",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			f, err := os.Open(args[0])
			if err != nil {
				return err
			}
			defer f.Close()

			parsed, err := opml.Parse(f)
			if err != nil {
				return err
			}
			if len(parsed) == 0 {
				fmt.Fprintf(cmd.OutOrStdout(), "no feeds found in %s\n", args[0])
				return nil
			}

			db, q, err := a.queries()
			if err != nil {
				return err
			}
			defer db.Close()

			now := time.Now().UTC().Format(time.RFC3339)
			var added, skipped int
			for _, pf := range parsed {
				rows, err := q.CreateFeedIfNew(cmd.Context(), store.CreateFeedIfNewParams{
					Url:       pf.URL,
					Title:     pf.Title,
					CreatedAt: now,
				})
				if err != nil {
					return fmt.Errorf("import %q: %w", pf.URL, err)
				}
				if rows == 1 {
					added++
				} else {
					skipped++ // url already present — ON CONFLICT DO NOTHING
				}
			}
			fmt.Fprintf(cmd.OutOrStdout(),
				"imported %s: %d new, %d already present\n", args[0], added, skipped)
			return nil
		},
	}
}

func newFeedsAddCmd(a *app) *cobra.Command {
	var title string
	cmd := &cobra.Command{
		Use:   "add <url>",
		Short: "Subscribe to a feed",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			db, q, err := a.queries()
			if err != nil {
				return err
			}
			defer db.Close()

			feed, err := q.CreateFeed(cmd.Context(), store.CreateFeedParams{
				Url:       args[0],
				Title:     title,
				CreatedAt: time.Now().UTC().Format(time.RFC3339), // app owns the RFC3339 format
			})
			if err != nil {
				return fmt.Errorf("add feed %q: %w", args[0], err)
			}
			fmt.Fprintf(cmd.OutOrStdout(), "added feed %d: %s\n", feed.ID, feed.Url)
			return nil
		},
	}
	cmd.Flags().StringVar(&title, "title", "", "optional human-readable title")
	return cmd
}

func newFeedsListCmd(a *app) *cobra.Command {
	return &cobra.Command{
		Use:   "list",
		Short: "List subscribed feeds",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, _ []string) error {
			db, q, err := a.queries()
			if err != nil {
				return err
			}
			defer db.Close()

			feeds, err := q.ListFeeds(cmd.Context())
			if err != nil {
				return err
			}
			if len(feeds) == 0 {
				fmt.Fprintln(cmd.OutOrStdout(), "no feeds yet — add one with `inkwell feeds add <url>`")
				return nil
			}

			// tabwriter aligns columns by padding to the widest cell per column.
			w := tabwriter.NewWriter(cmd.OutOrStdout(), 0, 0, 2, ' ', 0)
			fmt.Fprintln(w, "ID\tACTIVE\tURL\tTITLE")
			for _, f := range feeds {
				fmt.Fprintf(w, "%d\t%t\t%s\t%s\n", f.ID, f.Active, f.Url, f.Title)
			}
			return w.Flush()
		},
	}
}

func newFeedsRemoveCmd(a *app) *cobra.Command {
	return &cobra.Command{
		Use:   "remove <id>",
		Short: "Unsubscribe from a feed (also deletes its entries)",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			id, err := parseID(args[0])
			if err != nil {
				return err
			}
			db, q, err := a.queries()
			if err != nil {
				return err
			}
			defer db.Close()

			if err := q.DeleteFeed(cmd.Context(), id); err != nil {
				return err
			}
			fmt.Fprintf(cmd.OutOrStdout(), "removed feed %d\n", id)
			return nil
		},
	}
}

// newFeedsSetActiveCmd builds either `activate` or `deactivate` from one body,
// since they differ only by the boolean they write and the word they print.
func newFeedsSetActiveCmd(a *app, active bool) *cobra.Command {
	use, past := "activate", "activated"
	if !active {
		use, past = "deactivate", "deactivated"
	}
	return &cobra.Command{
		Use:   use + " <id>",
		Short: "Enable or disable fetching for a feed",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			id, err := parseID(args[0])
			if err != nil {
				return err
			}
			db, q, err := a.queries()
			if err != nil {
				return err
			}
			defer db.Close()

			if err := q.SetFeedActive(cmd.Context(), store.SetFeedActiveParams{
				Active: active,
				ID:     id,
			}); err != nil {
				return err
			}
			fmt.Fprintf(cmd.OutOrStdout(), "%s feed %d\n", past, id)
			return nil
		},
	}
}

func parseID(s string) (int64, error) {
	id, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		return 0, fmt.Errorf("invalid feed id %q: must be a number", s)
	}
	return id, nil
}
