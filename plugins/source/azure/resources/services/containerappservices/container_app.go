package containerappservices

import (
	"context"

	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/appcontainers/armappcontainers/v3"
	"github.com/cloudquery/cloudquery/plugins/source/azure/client"
	"github.com/cloudquery/plugin-sdk/v4/schema"
	"github.com/cloudquery/plugin-sdk/v4/transformers"
)

func ContainerApp() *schema.Table {
	return &schema.Table{
		Name:                 "azure_container_app",
		Resolver:             fetchContainerApps,
		PostResourceResolver: client.LowercaseIDResolver,
		Description:          "https://learn.microsoft.com/en-us/rest/api/containerapps/container-apps/list-by-subscription?view=rest-containerapps-2023-05-01&tabs=HTTP#containerapp",
		Multiplex:            client.SubscriptionMultiplexRegisteredNamespace("azure_container_app", client.Namespacemicrosoft_app),
		Transform:            transformers.TransformWithStruct(&armappcontainers.ContainerApp{}, transformers.WithPrimaryKeys("ID")),
		Columns:              schema.ColumnList{client.SubscriptionID},
	}
}

func fetchContainerApps(ctx context.Context, meta schema.ClientMeta, parent *schema.Resource, res chan<- any) error {
	cl := meta.(*client.Client)
	svc, err := armappcontainers.NewContainerAppsClient(cl.SubscriptionId, cl.Creds, cl.Options)
	if err != nil {
		return err
	}
	pager := svc.NewListBySubscriptionPager(nil)
	for pager.More() {
		p, err := pager.NextPage(ctx)
		if err != nil {
			return err
		}
		res <- p.Value
	}
	return nil
}
