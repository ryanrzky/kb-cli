package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/spf13/cobra"
)

var taskCmd = &cobra.Command{
	Use:   "task",
	Short: "Manage tasks",
}

var taskListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all tasks for a project",
	RunE: func(cmd *cobra.Command, args []string) error {
		projectID, _ := cmd.Flags().GetString("project-id")
		
		params := map[string]interface{}{
			"project_id": projectID,
			"status_id":  1, // Active tasks
		}
		
		res, err := kbClient.Call("getAllTasks", params)
		if err != nil {
			return err
		}

		if jsonOut {
			fmt.Println(string(res))
			return nil
		}

		var tasks []map[string]interface{}
		if err := json.Unmarshal(res, &tasks); err != nil {
			return err
		}

		// Fetch columns for mapping
		colRes, err := kbClient.Call("getColumns", map[string]interface{}{"project_id": projectID})
		colMap := make(map[string]string)
		if err == nil {
			var cols []map[string]interface{}
			if json.Unmarshal(colRes, &cols) == nil {
				for _, c := range cols {
					colMap[fmt.Sprintf("%v", c["id"])] = fmt.Sprintf("%v", c["title"])
				}
			}
		}

		// Fetch users for mapping
		userRes, err := kbClient.Call("getProjectUsers", map[string]interface{}{"project_id": projectID})
		userMap := make(map[string]string)
		if err == nil {
			var users map[string]interface{}
			if json.Unmarshal(userRes, &users) == nil {
				for k, v := range users {
					userMap[fmt.Sprintf("%v", k)] = fmt.Sprintf("%v", v)
				}
			}
		}

		w := tabwriter.NewWriter(os.Stdout, 0, 0, 3, ' ', 0)
		fmt.Fprintln(w, "ID\tTITLE\tCOLUMN\tASSIGNEE")
		for _, t := range tasks {
			colID := fmt.Sprintf("%v", t["column_id"])
			colName := colMap[colID]
			if colName == "" {
				colName = colID
			}

			ownerID := fmt.Sprintf("%v", t["owner_id"])
			assigneeName := userMap[ownerID]
			if assigneeName == "" || ownerID == "0" {
				assigneeName = "Unassigned"
			}

			fmt.Fprintf(w, "%v\t%v\t%v\t%v\n", t["id"], t["title"], colName, assigneeName)
		}
		w.Flush()
		return nil
	},
}

var taskCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a new task",
	RunE: func(cmd *cobra.Command, args []string) error {
		projectID, _ := cmd.Flags().GetString("project-id")
		title, _ := cmd.Flags().GetString("title")
		desc, _ := cmd.Flags().GetString("description")
		colID, _ := cmd.Flags().GetString("column-id")
		color, _ := cmd.Flags().GetString("color")
		
		params := map[string]interface{}{
			"project_id": projectID,
			"title":      title,
		}
		if desc != "" {
			params["description"] = desc
		}
		if colID != "" {
			params["column_id"] = colID
		}
		if color != "" {
			params["color_id"] = color
		}

		res, err := kbClient.Call("createTask", params)
		if err != nil {
			return err
		}

		if jsonOut {
			fmt.Println(string(res))
			return nil
		}

		fmt.Printf("Task created successfully with ID: %s\n", string(res))
		return nil
	},
}

var taskGetCmd = &cobra.Command{
	Use:   "get [id]",
	Short: "Get task by ID",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		params := map[string]interface{}{
			"task_id": args[0],
		}
		res, err := kbClient.Call("getTask", params)
		if err != nil {
			return err
		}

		if jsonOut {
			fmt.Println(string(res))
			return nil
		}
		
		var task map[string]interface{}
		if err := json.Unmarshal(res, &task); err != nil {
			return err
		}

		for k, v := range task {
			if v == nil {
				fmt.Printf("%s: None\n", k)
			} else {
				fmt.Printf("%s: %v\n", k, v)
			}
		}
		return nil
	},
}

var taskUpdateCmd = &cobra.Command{
	Use:   "update [id]",
	Short: "Update task",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		title, _ := cmd.Flags().GetString("title")
		color, _ := cmd.Flags().GetString("color")
		params := map[string]interface{}{
			"id": args[0],
		}
		if title != "" {
			params["title"] = title
		}
		if color != "" {
			params["color_id"] = color
		}
		
		res, err := kbClient.Call("updateTask", params)
		if err != nil {
			return err
		}

		if jsonOut {
			fmt.Println(string(res))
			return nil
		}

		fmt.Println("Task updated successfully:", string(res))
		return nil
	},
}

var taskDeleteCmd = &cobra.Command{
	Use:   "delete [id]",
	Short: "Delete task by ID",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		params := map[string]interface{}{
			"task_id": args[0],
		}
		res, err := kbClient.Call("removeTask", params)
		if err != nil {
			return err
		}

		if jsonOut {
			fmt.Println(string(res))
			return nil
		}

		fmt.Println("Task deleted successfully:", string(res))
		return nil
	},
}

func init() {
	rootCmd.AddCommand(taskCmd)
	
	taskCmd.AddCommand(taskListCmd)
	taskListCmd.Flags().StringP("project-id", "p", "", "Project ID (required)")
	taskListCmd.MarkFlagRequired("project-id")
	
	taskCmd.AddCommand(taskCreateCmd)
	taskCreateCmd.Flags().StringP("project-id", "p", "", "Project ID (required)")
	taskCreateCmd.MarkFlagRequired("project-id")
	taskCreateCmd.Flags().StringP("title", "t", "", "Task Title (required)")
	taskCreateCmd.MarkFlagRequired("title")
	taskCreateCmd.Flags().StringP("description", "d", "", "Task Description")
	taskCreateCmd.Flags().StringP("column-id", "c", "", "Column ID")
	taskCreateCmd.Flags().String("color", "", "Color ID (e.g. yellow, blue, red, green)")
	
	taskCmd.AddCommand(taskGetCmd)
	
	taskCmd.AddCommand(taskUpdateCmd)
	taskUpdateCmd.Flags().StringP("title", "t", "", "New Task Title")
	taskUpdateCmd.Flags().String("color", "", "New Color ID (e.g. yellow, blue, red, green)")
	
	taskCmd.AddCommand(taskDeleteCmd)
}
