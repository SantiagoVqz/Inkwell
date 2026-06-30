package main

import (
	"database/sql"
	"fmt"

	"github.com/santivqzv/inkwell/internal/config"
	"github.com/santivqzv/inkwell/internal/store"
	"github.com/spf13/cobra"
)

// app holds the dependencies shared by every command — currently just the
// resolved database path. It's threaded through command constructors instead
// of living in a package global, so the command tree stays testable (each test
// builds a fresh tree pointed at its own temp DB).
type app struct {
	dbPath string
}

// queries opens the store (which also runs migrations) and wraps it in the
// sqlc-generated Queries. Callers own the returned *sql.DB and must Close it.
func (a *app) queries() (*sql.DB, *store.Queries, error) {
	db, err := store.Open(a.dbPath)
	if err != nil {
		return nil, nil, err
	}
	return db, store.New(db), nil
}

func newRootCmd() *cobra.Command {
	a := &app{}

	root := &cobra.Command{
		Use:   "inkwell",
		Short: "A source-agnostic research engine that ingests feeds into a vault",
		// Don't dump usage text on a runtime error (e.g. a DB failure); usage is
		// only helpful for actual misuse, which cobra still reports on its own.
		SilenceUsage: true,
	}

	root.PersistentFlags().StringVar(&a.dbPath, "db", config.Default().DatabasePath,
		"path to the SQLite database")

	root.AddCommand(
		newVersionCmd(),
		newMigrateCmd(a),
		newFeedsCmd(a),
	)
	return root
}

func newVersionCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "version",
		Short: "Print the version",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, _ []string) error {
			fmt.Fprintln(cmd.OutOrStdout(), "inkwell", version)
			return nil
		},
	}
}

func newMigrateCmd(a *app) *cobra.Command {
	return &cobra.Command{
		Use:   "migrate",
		Short: "Create or update the database schema",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, _ []string) error {
			// store.Open applies every pending migration as part of opening.
			db, err := store.Open(a.dbPath)
			if err != nil {
				return err
			}
			defer db.Close()
			fmt.Fprintf(cmd.OutOrStdout(), "schema up to date at %s\n", a.dbPath)
			return nil
		},
	}
}
