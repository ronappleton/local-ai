package cmd

// This file defines CLI commands for managing user accounts. Currently it
// exposes a `create` subcommand allowing an administrator to create a new
// account directly from the terminal.

import (
	"codex/src/auth"
	"codex/src/memory"
	"fmt"
	"github.com/spf13/cobra"
)

// usersCmd is the parent for user related subcommands.
var usersCmd = &cobra.Command{
	Use:   "user",
	Short: "Manage users",
}

// createUserCmd creates a user in the local database.
var createUserCmd = &cobra.Command{
	Use:   "create [username] [email] [password]",
	Short: "Create a new user",
	Args:  cobra.ExactArgs(3),
	RunE: func(cmd *cobra.Command, args []string) error {
		username := args[0]
		email := args[1]
		password := args[2]
		db, err := memory.InitDB()
		if err != nil {
			return err
		}
		defer db.Close()
		return auth.CreateUser(db, username, email, password)
	},
}

// promoteUserCmd marks an existing user as an admin.
var promoteUserCmd = &cobra.Command{
	Use:   "promote [username]",
	Short: "Promote a user to admin",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		username := args[0]
		db, err := memory.InitDB()
		if err != nil {
			return err
		}
		defer db.Close()
		return auth.SetAdmin(db, username, true)
	},
}

// listUsersCmd prints all registered users.
var listUsersCmd = &cobra.Command{
	Use:   "list",
	Short: "List users",
	RunE: func(cmd *cobra.Command, args []string) error {
		db, err := memory.InitDB()
		if err != nil {
			return err
		}
		defer db.Close()
		list, err := auth.List(db)
		if err != nil {
			return err
		}
		for _, u := range list {
			fmt.Printf("%s\t%s\tadmin:%v verified:%v\n", u.Username, u.Email, u.Admin, u.Verified)
		}
		return nil
	},
}

// deleteAll indicates whether all users should be removed.
var deleteAll bool

// deleteUserCmd removes a specific user account or all if --all is set.
var deleteUserCmd = &cobra.Command{
	Use:   "delete [username]",
	Short: "Delete user accounts",
	Args:  cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		db, err := memory.InitDB()
		if err != nil {
			return err
		}
		defer db.Close()
		if deleteAll {
			return auth.DeleteAllUsers(db)
		}
		if len(args) != 1 {
			return fmt.Errorf("username required")
		}
		return auth.DeleteUser(db, args[0])
	},
}

func init() {
	usersCmd.AddCommand(createUserCmd)
	usersCmd.AddCommand(promoteUserCmd)
	usersCmd.AddCommand(listUsersCmd)
	deleteUserCmd.Flags().BoolVarP(&deleteAll, "all", "a", false, "delete all users")
	usersCmd.AddCommand(deleteUserCmd)
	rootCmd.AddCommand(usersCmd)
}
