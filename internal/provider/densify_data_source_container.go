package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/joelpereira/densify-api-client-go"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ datasource.DataSource              = &densifyDataSourceContainer{}
	_ datasource.DataSourceWithConfigure = &densifyDataSourceContainer{}
)

// NewDensifyDataSource is a helper function to simplify the provider implementation.
func NewDensifyDataSourceContainer() datasource.DataSource {
	return &densifyDataSourceContainer{}
}

// densifyDataSource is the data source implementation.
type densifyDataSourceContainer struct {
	// client  *hashicups.Client
	client *densify.Client
}

// densifyDataSourceModel maps the data source schema data.
// type densifyDataSourceModel struct {
// 	Recommendation []densifyRecoModel `tfsdk:"recommendation"`
// }

// densifyRecoModel maps coffees schema data.
// type densifyRecoModel struct {
type densifyDataSourcePodModel struct {
	// EntityId   types.Int64  `tfsdk:"entityId"`
	EntityId types.String `tfsdk:"entity_id"`
	Name     types.String `tfsdk:"name"`
	// OptimizationType types.String `tfsdk:"optimization_type"`
	AccountRef types.String `tfsdk:"account_ref"`

	// k8s variables
	Cluster        types.String `tfsdk:"cluster"`
	Namespace      types.String `tfsdk:"namespace"`
	ControllerType types.String `tfsdk:"controller_type"`
	PodName        types.String `tfsdk:"pod_name"`
	// ApprovalType   types.String `tfsdk:"approval_type"`
	ContainerCount types.Int64 `tfsdk:"container_count"`

	Containers map[string]densifyDataSourceContainerModel `tfsdk:"containers"`
}

type densifyDataSourceContainerModel struct {
	ContainerName    types.String `tfsdk:"container_name"`
	OptimizationType types.String `tfsdk:"optimization_type"`

	CurCPUReq types.String `tfsdk:"current_cpu_request"`
	CurCPULim types.String `tfsdk:"current_cpu_limit"`
	CurMemReq types.String `tfsdk:"current_mem_request"`
	CurMemLim types.String `tfsdk:"current_mem_limit"`

	RecCPUReq types.String `tfsdk:"recommended_cpu_request"`
	RecCPULim types.String `tfsdk:"recommended_cpu_limit"`
	RecMemReq types.String `tfsdk:"recommended_mem_request"`
	RecMemLim types.String `tfsdk:"recommended_mem_limit"`
}

// Metadata returns the data source type name.
func (d *densifyDataSourceContainer) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_container"
}

// Schema defines the schema for the data source.
func (d *densifyDataSourceContainer) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Fetches a Kubernetes (EKS/AKS/GKE) Container Recommendation from the Densify API.",
		Attributes: map[string]schema.Attribute{
			"entity_id": schema.StringAttribute{
				Computed:    true,
				Description: "Unique identifier for container resource.",
			},
			"name": schema.StringAttribute{
				Computed:    true,
				Description: "Container manifest name.",
			},
			"account_ref": schema.StringAttribute{
				Computed:    true,
				Description: "Account reference identifier.",
			},
			// "approval_type": schema.StringAttribute{
			// 	Computed:    true,
			// 	Description: "Approval type. If ITSM integration has been enabled, this field will define whether the recommendation has been reviewed & approved.",
			// },

			// Kubernetes variables
			"cluster": schema.StringAttribute{
				Computed:    true,
				Description: "The Kubernetes cluster name.",
			},
			"namespace": schema.StringAttribute{
				Computed:    true,
				Description: "The Kubernetes namespace.",
			},
			"controller_type": schema.StringAttribute{
				Computed:    true,
				Description: "The Kubernetes controller type. Ex. deployment, daemonset, statefulset, cronjob, job, pod.",
			},
			"pod_name": schema.StringAttribute{
				Computed:    true,
				Description: "The Kubernetes pod name.",
			},
			"container_count": schema.Int64Attribute{
				Computed:    true,
				Description: "The number of containers within the pod recommendation.",
			},

			// nested/multiple container recommendations
			// https://developer.hashicorp.com/terraform/plugin/framework/handling-data/attributes/map-nested
			"containers": schema.MapNestedAttribute{
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"container_name": schema.StringAttribute{
							Computed:    true,
							Description: "The Kubernetes container name.",
						},
						"optimization_type": schema.StringAttribute{
							Computed:    true,
							Description: "Type of optimization. Ex. Downsize, Upsize, Resize, Terminate, etc.",
						},
						"current_cpu_request": schema.StringAttribute{
							Computed:    true,
							Description: "The current CPU Request for resources (in millicores or m).",
						},
						"current_cpu_limit": schema.StringAttribute{
							Computed:    true,
							Description: "The current CPU Limit for resources (in millicores or m).",
						},
						"current_mem_request": schema.StringAttribute{
							Computed:    true,
							Description: "The current Memory Request for resources (in mebibytes or Mi).",
						},
						"current_mem_limit": schema.StringAttribute{
							Computed:    true,
							Description: "The current Memory Limit for resources (in mebibytes or Mi).",
						},

						"recommended_cpu_request": schema.StringAttribute{
							Computed:    true,
							Description: "The recommended CPU Request for resources (in millicores or m).",
						},
						"recommended_cpu_limit": schema.StringAttribute{
							Computed:    true,
							Description: "The recommended CPU Limit for resources (in millicores or m).",
						},
						"recommended_mem_request": schema.StringAttribute{
							Computed:    true,
							Description: "The recommended Memory Request for resources (in mebibytes or Mi).",
						},
						"recommended_mem_limit": schema.StringAttribute{
							Computed:    true,
							Description: "The recommended Memory Limit for resources (in mebibytes or Mi).",
						},
					},
				},
				Computed: true,
			},
		},
	}
}

// Configure adds the provider configured client to the data source.
func (d *densifyDataSourceContainer) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	tflog.Trace(ctx, "Configuring Densify API client")
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*densify.Client)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
			fmt.Sprintf("Expected *densify.Client, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}

	d.client = client
}

// Read refreshes the Terraform state with the latest data.
func (d *densifyDataSourceContainer) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	tflog.Trace(ctx, "Reading Densify API client")
	var state densifyDataSourcePodModel

	tflog.Debug(ctx, "Densify API client: calling GetAccountOrCluster")
	_, err := d.client.GetAccountOrCluster()
	if err != nil {
		if d.client.Query.SkipErrors {
			// skip the error message
			return
		}
		resp.Diagnostics.AddError(
			"Unable to Find Densify Account Number/Name",
			err.Error(),
		)
		return
	}
	tflog.Trace(ctx, "Densify API client: GetAccountOrCluster: success")
	tflog.Debug(ctx, "Densify API client: calling GetDensifyRecommendation")
	podReco, err := d.client.GetDensifyRecommendation()
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Find Densify Recommendation",
			err.Error(),
		)
		return
	}
	tflog.Trace(ctx, "Densify API client: GetDensifyRecommendation: success")

	// if we didn't get a recommendation, return an empty one (instead of nil)
	if podReco == nil {
		podReco = &densify.DensifyRecommendation{}
	}

	if podReco != nil {
		// Map response body to model
		state.EntityId = types.StringValue(podReco.EntityId)
		state.Name = types.StringValue(podReco.Name)
		state.AccountRef = types.StringValue(podReco.AccountIdRef)
		// state.OptimizationType = types.StringValue(podReco.RecommendationType)
		// state.ApprovalType = types.StringValue(podReco.ApprovalType)

		// Kubernetes variables
		state.Cluster = types.StringValue(podReco.Cluster)
		state.Namespace = types.StringValue(podReco.Namespace)
		state.ControllerType = types.StringValue(podReco.ControllerType)
		state.PodName = types.StringValue(podReco.PodService)

		cpuUnit := "m"  // millicores
		memUnit := "Mi" // mebibytes

		if state.Containers == nil {
			state.Containers = map[string]densifyDataSourceContainerModel{}
		}

		// containersMap := map[string]densifyDataSourceContainerModel{}
		tflog.Debug(ctx, fmt.Sprintf(`Num of Containers: %d`, len(podReco.Containers)))
		for i := 0; i < len(podReco.Containers); i++ {
			reco := podReco.Containers[i]
			c := densifyDataSourceContainerModel{}
			c.ContainerName = types.StringValue(reco.Container)
			if reco.Container != "" {
				state.Name = types.StringValue(reco.Container)
			}
			c.OptimizationType = types.StringValue(reco.RecommendationType)

			c.CurCPUReq = types.StringValue(fmt.Sprintf(`%d%s`, reco.CurrentCpuRequest, cpuUnit))
			c.CurCPULim = types.StringValue(fmt.Sprintf(`%d%s`, reco.CurrentCpuLimit, cpuUnit))
			c.CurMemReq = types.StringValue(fmt.Sprintf(`%d%s`, reco.CurrentMemRequest, memUnit))
			c.CurMemLim = types.StringValue(fmt.Sprintf(`%d%s`, reco.CurrentMemLimit, memUnit))

			if reco.RecommendedCpuRequest > 0 || reco.RecommendedMemRequest > 0 {
				c.RecCPUReq = types.StringValue(fmt.Sprintf(`%d%s`, reco.RecommendedCpuRequest, cpuUnit))
				c.RecCPULim = types.StringValue(fmt.Sprintf(`%d%s`, reco.RecommendedCpuLimit, cpuUnit))
				c.RecMemReq = types.StringValue(fmt.Sprintf(`%d%s`, reco.RecommendedMemRequest, memUnit))
				c.RecMemLim = types.StringValue(fmt.Sprintf(`%d%s`, reco.RecommendedMemLimit, memUnit))
			} else {
				// if there are no recommendations, take the fallback values and output them as recommended
				c.RecCPUReq = types.StringValue(reco.FallbackCpuRequest)
				c.RecCPULim = types.StringValue(reco.FallbackCpuLimit)
				c.RecMemReq = types.StringValue(reco.FallbackMemRequest)
				c.RecMemLim = types.StringValue(reco.FallbackMemLimit)
			}

			// add the container to the map of containers
			state.Containers[reco.Container] = c
		}
	}

	// now we can set the count of containers
	state.ContainerCount = types.Int64Value(int64(len(podReco.Containers)))

	// Set state
	diags := resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		tflog.Trace(ctx, fmt.Sprintf(`Errors: %s`, resp.Diagnostics.Errors()))
		return
	}
}
