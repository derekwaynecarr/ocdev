package url

import (
	"fmt"

	"github.com/pkg/errors"
	"github.com/redhat-developer/ocdev/pkg/application"
	"github.com/redhat-developer/ocdev/pkg/component"
	"github.com/redhat-developer/ocdev/pkg/occlient"
	log "github.com/sirupsen/logrus"
)

type URL struct {
	Name string
	URL  string
}

// Delete deletes a URL
func Delete(client *occlient.Client, name string) error {
	return client.DeleteRoute(name)
}

// Create creates a URL
func Create(client *occlient.Client, cmp string) (*URL, error) {

	app, err := application.GetCurrentOrDefault(client)
	if err != nil {
		return nil, errors.Wrap(err, "unable to get current application")
	}

	labels, err := component.GetLabels(cmp, app, false)
	if err != nil {
		return nil, errors.Wrap(err, "unable to get labels")
	}

	route, err := client.CreateRoute(cmp, labels)
	if err != nil {
		return nil, errors.Wrap(err, "unable to create route")
	}

	return &URL{
		Name: route.Name,
		URL:  route.Spec.Host,
	}, nil
}

// List lists the URLs in an application. The results can further be narrowed
// down if a component name is provided, which will only list URLs for the
// given component
func List(client *occlient.Client, componentName string, applicationName string) ([]URL, error) {

	labelSelector := fmt.Sprintf("%v=%v", application.ApplicationLabel, applicationName)

	if componentName != "" {
		labelSelector = labelSelector + fmt.Sprintf(",%v=%v", component.ComponentLabel, componentName)
	}

	log.Debugf("Listing routes with label selector: %v", labelSelector)
	routes, err := client.ListRoutes(labelSelector)
	if err != nil {
		return nil, errors.Wrap(err, "unable to list route names")
	}

	var urls []URL
	for _, r := range routes {
		urls = append(urls, URL{
			Name: r.Name,
			URL:  r.Spec.Host,
		})
	}

	return urls, nil
}
