package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/spf13/cobra"
)

var projectCmd = &cobra.Command{
	Use:   "project",
	Short: "Manage projects",
}

var projectListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all projects",
	RunE: func(cmd *cobra.Command, args []string) error {
		res, err := kbClient.Call("getAllProjects", nil)
		if err != nil {
			return err
		}

		if jsonOut {
			fmt.Println(string(res))
			return nil
		}

		var projects []map[string]interface{}
		if err := json.Unmarshal(res, &projects); err != nil {
			return err
		}

		w := tabwriter.NewWriter(os.Stdout, 0, 0, 3, ' ', 0)
		fmt.Fprintln(w, "ID\tNAME\tACTIVE")
		for _, p := range projects {
			fmt.Fprintf(w, "%v\t%v\t%v\n", p["id"], p["name"], p["is_active"])
		}
		w.Flush()
		return nil
	},
}

var projectCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a new project",
	RunE: func(cmd *cobra.Command, args []string) error {
		name, _ := cmd.Flags().GetString("name")
		desc, _ := cmd.Flags().GetString("description")
		
		params := map[string]interface{}{
			"name": name,
		}
		if desc != "" {
			params["description"] = desc
		}

		res, err := kbClient.Call("createProject", params)
		if err != nil {
			return err
		}

		if jsonOut {
			fmt.Println(string(res))
			return nil
		}

		fmt.Printf("Project created successfully with ID: %s\n", string(res))
		return nil
	},
}

var projectGetCmd = &cobra.Command{
	Use:   "get [id]",
	Short: "Get project by ID",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		params := map[string]interface{}{
			"project_id": args[0],
		}
		res, err := kbClient.Call("getProjectById", params)
		if err != nil {
			return err
		}

		if jsonOut {
			fmt.Println(string(res))
			return nil
		}
		
		var project map[string]interface{}
		if err := json.Unmarshal(res, &project); err != nil {
			return err
		}

		for k, v := range project {
			if v == nil {
				fmt.Printf("%s: None\n", k)
			} else {
				fmt.Printf("%s: %v\n", k, v)
			}
		}
		return nil
	},
}

var projectUpdateCmd = &cobra.Command{
	Use:   "update [id]",
	Short: "Update project",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		name, _ := cmd.Flags().GetString("name")
		params := map[string]interface{}{
			"id": args[0],
		}
		if name != "" {
			params["name"] = name
		}
		
		res, err := kbClient.Call("updateProject", params)
		if err != nil {
			return err
		}

		if jsonOut {
			fmt.Println(string(res))
			return nil
		}

		fmt.Println("Project updated successfully:", string(res))
		return nil
	},
}

var projectDeleteCmd = &cobra.Command{
	Use:   "delete [id]",
	Short: "Delete project by ID",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		params := map[string]interface{}{
			"project_id": args[0],
		}
		res, err := kbClient.Call("removeProject", params)
		if err != nil {
			return err
		}

		if jsonOut {
			fmt.Println(string(res))
			return nil
		}

		fmt.Println("Project deleted successfully:", string(res))
		return nil
	},
}

func init() {
	rootCmd.AddCommand(projectCmd)
	projectCmd.AddCommand(projectListCmd)
	
	projectCmd.AddCommand(projectCreateCmd)
	projectCreateCmd.Flags().StringP("name", "n", "", "Project Name (required)")
	projectCreateCmd.MarkFlagRequired("name")
	projectCreateCmd.Flags().StringP("description", "d", "", "Project Description")
	
	projectCmd.AddCommand(projectGetCmd)
	
	projectCmd.AddCommand(projectUpdateCmd)
	projectUpdateCmd.Flags().StringP("name", "n", "", "New Project Name")
	
	projectCmd.AddCommand(projectDeleteCmd)
}
