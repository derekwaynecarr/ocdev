package cmd

import (
	"fmt"
	"os"

	"github.com/redhat-developer/ocdev/pkg/application"
	"github.com/redhat-developer/ocdev/pkg/component"
	"github.com/redhat-developer/ocdev/pkg/url"
	"github.com/spf13/cobra"
)

var (
	urlListComponent   string
	urlListApplication string
)

var urlCmd = &cobra.Command{
	Use:   "url",
	Short: "Expose component to the outside world",
	Long: `Expose component to the outside world.
The URLs that are generated using this command, can be used to access the
deployed components from outside the cluster.
`,
}

var urlCreateCmd = &cobra.Command{
	Use:   "create [component name]",
	Short: "Create a URL for a component",
	Long: `Create a URL for a component.
The created URL can be used to access the specified component from outside the
OpenShift cluster.
`,
	Example: `# Create a URL for the current component.
ocdev url create

# Create a URL for a specific component
ocdev url create <component name>
`,
	Args: cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		client := getOcClient()

		var cmp string
		switch len(args) {
		case 0:
			var err error
			cmp, err = component.GetCurrent(client)
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
		case 1:
			cmp = args[0]
		default:
			fmt.Println("unable to get component")
			os.Exit(1)
		}

		fmt.Printf("Adding URL to component: %v\n", cmp)
		u, err := url.Create(client, cmp)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		fmt.Printf("URL created for component: %v\n\n"+
			"%v - %v\n", cmp, u.Name, u.URL)
	},
}

var urlDeleteCmd = &cobra.Command{
	Use:   "delete <URL>",
	Short: "Delete a URL",
	Long: `Delete a URL.
Deleted the given URL, hence making the service inaccessible.
`,
	Example: `# Delete a URL to a component
ocdev url delete <URL>
`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		client := getOcClient()
		u := args[0]
		if err := url.Delete(client, u); err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		fmt.Printf("Deleted URL: %v\n", u)
	},
}

var urlListCmd = &cobra.Command{
	Use:   "list",
	Short: "List URLs",
	Long: `List URLs.
Lists all the available URLs which can be used to access the components.`,
	Example: ` # List the available URLs
ocdev url list
`,
	Args: cobra.ExactArgs(0),
	Run: func(cmd *cobra.Command, args []string) {
		client := getOcClient()

		var app string
		if urlListApplication == "" {
			var err error
			app, err = application.GetCurrent(client)
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
		} else {
			app = urlListApplication
		}

		cmp := urlListComponent
		urls, err := url.List(client, cmp, app)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		if len(urls) == 0 {
			fmt.Printf("No URLs found for component %v in application %v\n", cmp, app)
		} else {
			fmt.Printf("Found the following URLs for component %v in application %v:\n", cmp, app)
			for _, u := range urls {
				fmt.Printf("%v - %v\n", u.Name, u.URL)
			}
		}
	},
}

func init() {
	urlListCmd.Flags().StringVarP(&urlListApplication, "application", "a", "", "list URLs for application")
	urlListCmd.Flags().StringVarP(&urlListComponent, "component", "c", "", "list URLs for component")

	urlCmd.AddCommand(urlListCmd)
	urlCmd.AddCommand(urlDeleteCmd)
	urlCmd.AddCommand(urlCreateCmd)
	rootCmd.AddCommand(urlCmd)
}
