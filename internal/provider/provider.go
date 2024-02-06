package provider

import (
	"context"
	"os"
	"strconv"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/joelpereira/densify-api-client-go"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ provider.Provider = &densifyProvider{}
)

// New is a helper function to simplify provider server and testing implementation.
func New(version string) func() provider.Provider {
	return func() provider.Provider {
		return &densifyProvider{
			version: version,
		}
	}
}

// Densify API query configuration settings.
type DensifySettings struct {
	instance     string
	username     string
	password     string
	timeout      int
	techPlatform string
	// cloud.
	accountName          string
	accountNumber        string
	systemName           string
	fallbackInstanceType string
	continueIfError      bool
	// containers/k8s.
	cluster        string
	namespace      string
	controllerType string
	podName        string
	containerName  string
	fallbackCPUReq string
	fallbackCPULim string
	fallbackMemReq string
	fallbackMemLim string
}

// densifyProvider is the provider implementation.
type densifyProvider struct {
	// version is set to the provider version on release, "dev" when the
	// provider is built and ran locally, and "test" when running acceptance
	// testing.
	version string
}

// densifyProviderModel maps provider schema data to a Go type.
type densifyProviderModel struct {
	DensifyInstance      types.String `tfsdk:"densify_instance"`
	Username             types.String `tfsdk:"username"`
	Password             types.String `tfsdk:"password"`
	ApiTimeout           types.Int64  `tfsdk:"api_timeout"`
	TechPlatform         types.String `tfsdk:"tech_platform"`
	AccountNumber        types.String `tfsdk:"account_number"`
	AccountName          types.String `tfsdk:"account_name"`
	SystemName           types.String `tfsdk:"system_name"`
	FallbackInstanceType types.String `tfsdk:"fallback_instance_type"`
	ContinueIfError      types.Bool   `tfsdk:"continue_if_error"`

	K8sCluster        types.String `tfsdk:"cluster"`
	K8sNamespace      types.String `tfsdk:"namespace"`
	K8sControllerType types.String `tfsdk:"controller_type"`
	K8sPodName        types.String `tfsdk:"pod_name"`
	K8sContainerName  types.String `tfsdk:"container_name"`

	K8sFallbackCPUReq types.String `tfsdk:"fallback_cpu_req"`
	K8sFallbackCPULim types.String `tfsdk:"fallback_cpu_lim"`
	K8sFallbackMemReq types.String `tfsdk:"fallback_mem_req"`
	K8sFallbackMemLim types.String `tfsdk:"fallback_mem_lim"`
}

// Metadata returns the provider type name.
func (p *densifyProvider) Metadata(_ context.Context, _ provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "densify"
	resp.Version = p.version
}

// Schema defines the provider-level schema for configuration data.
func (p *densifyProvider) Schema(_ context.Context, _ provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Retrieve right-sizing optimizations directly from the Densify API for cloud and/or container resources.",
		Attributes: map[string]schema.Attribute{
			"densify_instance": schema.StringAttribute{
				Optional:    true,
				Description: "URI for your Densify instance. May also be provided via DENSIFY_INSTANCE environment variable. Ex. https://instance.densify.com:8443",
			},
			"username": schema.StringAttribute{
				Optional:    true,
				Description: "Username to authenticate to Densify API. May also be provided via DENSIFY_USERNAME environment variable. Contact your Account Manager to request a service account details.",
			},
			"password": schema.StringAttribute{
				Optional:    true,
				Sensitive:   true,
				Description: "Password to authenticate to Densify API. May also be provided via DENSIFY_PASSWORD environment variable. Contact your Account Manager to request a service account details.",
			},
			"api_timeout": schema.Int64Attribute{
				Optional:    true,
				Description: "The Densify API timeout. The default value is 30 seconds but this can be adjusted via the DENSIFY_API_TIMEOUT environment variable.",
			},
			"tech_platform": schema.StringAttribute{
				Optional:    true,
				Description: "Which Cloud Service Provider (CSP) / technology platform to use for the Densify API. May also be provided via DENSIFY_TECH_PLATFORM environment variable. Accepted values are: aws, azure, gcp, k8s, kubernetes.",
			},

			// cloud parameters.
			"account_number": schema.StringAttribute{
				Optional:    true,
				Description: "The CSP (Cloud Service Provider) account number to check for a recommendation.",
			},
			"account_name": schema.StringAttribute{
				Optional:    true,
				Description: "The CSP (Cloud Service Provider) account name to check for a recommendation.",
			},
			"system_name": schema.StringAttribute{
				Optional:    true,
				Description: "The system name to check for a recommendation.",
			},
			"fallback_instance_type": schema.StringAttribute{
				Optional:    true,
				Description: "The fallback / default instance type to use. You may use the approved_type output value which will use this fallback instance by default, until a recommendation is generated by Densify and approved (manually or with full ITSM integration).",
			},
			"continue_if_error": schema.BoolAttribute{
				Optional:    true,
				Description: "Prevent errors from interupting the terraform deployment.",
			},

			// k8s parameters.
			"cluster": schema.StringAttribute{
				Optional:    true,
				Description: "Kubernetes namespace to look for a recommendation in Densify.",
			},
			"namespace": schema.StringAttribute{
				Optional:    true,
				Description: "Kubernetes namespace to look for a recommendation in Densify.",
			},
			"controller_type": schema.StringAttribute{
				Optional:    true,
				Description: "Kubernetes controller type to look for a recommendation in Densify. Accepted values are: deployment, replicaset, statefulset, daemonset, cronjob, job, pod.",
			},
			"pod_name": schema.StringAttribute{
				Optional:    true,
				Description: "Kubernetes pod name to look for a recommendation in Densify.",
			},
			"container_name": schema.StringAttribute{
				Optional:    true,
				Description: "Kubernetes container name to look for a recommendation in Densify.",
			},
			"fallback_cpu_req": schema.StringAttribute{
				Optional:    true,
				Description: "Fallback CPU request values, in millicores (m).",
			},
			"fallback_cpu_lim": schema.StringAttribute{
				Optional:    true,
				Description: "Fallback CPU limit values, in millicores (m).",
			},
			"fallback_mem_req": schema.StringAttribute{
				Optional:    true,
				Description: "Fallback Memory request values, in mebibytes (Mi).",
			},
			"fallback_mem_lim": schema.StringAttribute{
				Optional:    true,
				Description: "Fallback Memory limit values, in mebibytes (Mi).",
			},
		},
	}
}

// Configure prepares a Densify API client for data sources and resources.
func (p *densifyProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	tflog.Debug(ctx, "Configuring Densify client")
	// Retrieve provider data from configuration
	var config densifyProviderModel
	diags := req.Config.Get(ctx, &config)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// validate all the config parameters for the provider.
	config.ValidateProviderParameters(resp)

	if resp.Diagnostics.HasError() {
		return
	}

	// Default values to environment variables, but override with Terraform configuration value if set.
	densifysettings := DensifySettings{}
	tflog.Debug(ctx, "Loading Densify Settings from Environment Variables")
	densifysettings.LoadEnvironmentVariablesSettings(config)
	tflog.Debug(ctx, "Loading Densify Settings from Provider config")
	densifysettings.LoadConfigSettings(config)
	tflog.Debug(ctx, "Validating Densify Settings have all the required values")
	densifysettings.ValidateSettings(resp)

	if resp.Diagnostics.HasError() {
		return
	}

	// Send variables to log.
	ctx = tflog.SetField(ctx, "densify_instance", densifysettings.instance)
	ctx = tflog.SetField(ctx, "densify_username", densifysettings.username)
	ctx = tflog.SetField(ctx, "densify_password", densifysettings.password)
	ctx = tflog.MaskFieldValuesWithFieldKeys(ctx, "densify_password")
	ctx = tflog.SetField(ctx, "densify_api_timeout", densifysettings.timeout)
	ctx = tflog.SetField(ctx, "densify_tech_platform", densifysettings.techPlatform)
	ctx = tflog.SetField(ctx, "densify_account_name", densifysettings.accountName)
	ctx = tflog.SetField(ctx, "densify_account_number", densifysettings.accountNumber)
	ctx = tflog.SetField(ctx, "densify_system_name", densifysettings.systemName)
	ctx = tflog.SetField(ctx, "densify_fallback_instance_type", densifysettings.fallbackInstanceType)
	ctx = tflog.SetField(ctx, "densify_continue_if_error", densifysettings.continueIfError)
	ctx = tflog.SetField(ctx, "densify_cluster", densifysettings.cluster)
	ctx = tflog.SetField(ctx, "densify_namespace", densifysettings.namespace)
	ctx = tflog.SetField(ctx, "densify_controller_type", densifysettings.controllerType)
	ctx = tflog.SetField(ctx, "densify_pod_name", densifysettings.podName)
	ctx = tflog.SetField(ctx, "densify_container_name", densifysettings.containerName)
	ctx = tflog.SetField(ctx, "densify_fallback_cpu_req", densifysettings.fallbackCPUReq)
	ctx = tflog.SetField(ctx, "densify_fallback_cpu_lim", densifysettings.fallbackCPULim)
	ctx = tflog.SetField(ctx, "densify_fallback_mem_req", densifysettings.fallbackMemReq)
	ctx = tflog.SetField(ctx, "densify_fallback_mem_lim", densifysettings.fallbackMemLim)
	tflog.Debug(ctx, "Creating Densify API client")

	// Create a new Densify client using the configuration values.
	client, err := densify.NewDensifyClient(&densifysettings.instance, &densifysettings.username, &densifysettings.password, densifysettings.timeout)
	if err != nil && !densifysettings.continueIfError {
		resp.Diagnostics.AddError(
			"Unable to Create Densify API Client",
			"An unexpected error occurred when creating the Densify API client. "+
				"If the error is not clear, please contact the provider developers.\n\n"+
				"Densify Client Error: "+err.Error(),
		)
		return
	}

	// set configuration for Densify API Client.
	densifyAPIQuery := densify.DensifyAPIQuery{
		AnalysisTechnology: densifysettings.techPlatform,
		AccountName:        densifysettings.accountName,
		AccountNumber:      densifysettings.accountNumber,
		SystemName:         densifysettings.systemName,
		FallbackInstance:   densifysettings.fallbackInstanceType,
		SkipErrors:         densifysettings.continueIfError,

		K8sCluster:        densifysettings.cluster,
		K8sNamespace:      densifysettings.namespace,
		K8sControllerType: densifysettings.controllerType,
		K8sPodName:        densifysettings.podName,
		K8sContainerName:  densifysettings.containerName,

		FallbackCPURequest: densifysettings.fallbackCPUReq,
		FallbackCPULimit:   densifysettings.fallbackCPULim,
		FallbackMemRequest: densifysettings.fallbackMemReq,
		FallbackMemLimit:   densifysettings.fallbackMemLim,
	}

	tflog.Debug(ctx, "Validating Densify client query")
	err = client.ConfigureQuery(&densifyAPIQuery)
	if err != nil && !densifysettings.continueIfError {
		resp.Diagnostics.AddError(
			"Unable to create Densify query",
			"An unexpected error occurred when creating the Densify API query. "+
				"If the error is not clear, please contact the provider developers.\n\n"+
				"Densify Client Query Error: "+err.Error(),
		)
		return
	}
	tflog.Debug(ctx, "Validated Densify client query")

	// Make the Densify client available during DataSource and Resource type Configure methods.
	resp.DataSourceData = client
	resp.ResourceData = client

	tflog.Debug(ctx, "Configured Densify client", map[string]any{"success": true})
}

// DataSources defines the data sources implemented in the provider.
func (p *densifyProvider) DataSources(ctx context.Context) []func() datasource.DataSource {
	tflog.Trace(ctx, "Densify client DataSources")
	return []func() datasource.DataSource{
		NewDensifyDataSourceCloud,
		NewDensifyDataSourceContainer,
	}
}

// Resources defines the resources implemented in the provider.
func (p *densifyProvider) Resources(ctx context.Context) []func() resource.Resource {
	tflog.Trace(ctx, "Densify client Resources")
	return []func() resource.Resource{}
}

func (config *densifyProviderModel) ValidateProviderParameters(resp *provider.ConfigureResponse) {
	// If practitioner provided a configuration value for any of the
	// attributes, it must be a known value.
	if config.DensifyInstance.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("densify_instance"),
			"Unknown Densify API Instance",
			"The provider cannot create the Densify API client as there is an unknown configuration value for the Densify API instance name. "+
				"Either target apply the source of the value first, set the value statically in the configuration, or use the DENSIFY_INSTANCE environment variable.",
		)
	}

	if config.Username.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("username"),
			"Unknown Densify API Username",
			"The provider cannot create the Densify API client as there is an unknown configuration value for the Densify API username. "+
				"Either target apply the source of the value first, set the value statically in the configuration, or use the DENSIFY_USERNAME environment variable.",
		)
	}

	if config.Password.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("password"),
			"Unknown Densify API Password",
			"The provider cannot create the Densify API client as there is an unknown configuration value for the Densify API password. "+
				"Either target apply the source of the value first, set the value statically in the configuration, or use the DENSIFY_PASSWORD environment variable.",
		)
	}

	if config.TechPlatform.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("tech_platform"),
			"Unknown Densify API Technology Platform",
			"The provider cannot create the Densify API client as there is an unknown configuration value for the Densify API Technology Platform. "+
				"Either target apply the source of the value first, set the value statically in the configuration, or use the DENSIFY_TECH_PLATFORM environment variable.",
		)
	}

	// KUBERNETES/CONTAINERS.
	if strings.ToLower(config.TechPlatform.ValueString()) == "k8s" || strings.ToLower(config.TechPlatform.ValueString()) == "kubernetes" {
		if config.K8sCluster.IsUnknown() {
			resp.Diagnostics.AddAttributeError(
				path.Root("cluster"),
				"Unknown Kubernetes Cluster",
				"The provider cannot create the Densify API client as there is an unknown configuration value for the Kubernetes Cluster. "+
					"Either target apply the source of the value first, set the value statically in the configuration, or use the DENSIFY_CLUSTER environment variable.",
			)
		}
		if config.K8sNamespace.IsUnknown() {
			resp.Diagnostics.AddAttributeError(
				path.Root("namespace"),
				"Unknown Kubernetes Namespace",
				"The provider cannot create the Densify API client as there is an unknown configuration value for the Kubernetes Namespace. "+
					"Either target apply the source of the value first, set the value statically in the configuration, or use the DENSIFY_NAMESPACE environment variable.",
			)
		}
		if config.K8sControllerType.IsUnknown() {
			resp.Diagnostics.AddAttributeError(
				path.Root("controller_type"),
				"Unknown Kubernetes Controller Type",
				"The provider cannot create the Densify API client as there is an unknown configuration value for the Kubernetes Controller Type. "+
					"Either target apply the source of the value first, set the value statically in the configuration, or use the DENSIFY_CONTROLLER_TYPE environment variable.",
			)
		}
		if config.K8sPodName.IsUnknown() {
			resp.Diagnostics.AddAttributeError(
				path.Root("pod_name"),
				"Unknown Kubernetes Pod Name",
				"The provider cannot create the Densify API client as there is an unknown configuration value for the Kubernetes Pod Name. "+
					"Either target apply the source of the value first, set the value statically in the configuration, or use the DENSIFY_POD_NAME environment variable.",
			)
		}
	} else { // CLOUD.
		if config.AccountName.IsUnknown() {
			resp.Diagnostics.AddAttributeError(
				path.Root("account_name"),
				"Unknown Densify API Account Name",
				"The provider cannot create the Densify API client as there is an unknown configuration value for the Densify API Account Name. "+
					"Either target apply the source of the value first, set the value statically in the configuration, or use the DENSIFY_ACCOUNT_NAME environment variable.",
			)
		}
		if config.SystemName.IsUnknown() {
			resp.Diagnostics.AddAttributeError(
				path.Root("system_name"),
				"Unknown Densify System Name",
				"The provider cannot create the Densify API client as there is an unknown configuration value for the Densify System Name. "+
					"Either target apply the source of the value first, set the value statically in the configuration, or use the DENSIFY_SYSTEM_NAME environment variable.",
			)
		}
	}
}

// Load Densify settings from Environment Variables.
func (densifysettings *DensifySettings) LoadEnvironmentVariablesSettings(config densifyProviderModel) {
	// set default timeout (seconds);
	tout := 45
	// gracefully handle if the timeout is not a valid int.
	if val, err := strconv.Atoi(os.Getenv("DENSIFY_API_TIMEOUT")); err == nil {
		// make sure the timeout (seconds) is between 1-300 (5 mins) seconds.
		if val >= 1 && val <= 300 {
			tout = val
		}
	}
	densifysettings.timeout = tout
	densifysettings.instance = os.Getenv("DENSIFY_INSTANCE")
	densifysettings.username = os.Getenv("DENSIFY_USERNAME")
	densifysettings.password = os.Getenv("DENSIFY_PASSWORD")
	densifysettings.techPlatform = os.Getenv("DENSIFY_TECH_PLATFORM")
	densifysettings.accountName = os.Getenv("DENSIFY_ACCOUNT_NAME")
	densifysettings.accountNumber = os.Getenv("DENSIFY_ACCOUNT_NUMBER")
	densifysettings.systemName = os.Getenv("DENSIFY_SYSTEM_NAME")
	densifysettings.fallbackInstanceType = os.Getenv("DENSIFY_FALLBACK_INSTANCE_TYPE")
	densifysettings.continueIfError = false // default is false
	if strings.ToLower(os.Getenv("DENSIFY_CONTINUE_IF_ERROR")) == "true" {
		densifysettings.continueIfError = true
	}
	densifysettings.cluster = os.Getenv("DENSIFY_CLUSTER")
	densifysettings.namespace = os.Getenv("DENSIFY_NAMESPACE")
	densifysettings.controllerType = os.Getenv("DENSIFY_CONTROLLER_TYPE")
	densifysettings.podName = os.Getenv("DENSIFY_POD_NAME")
	densifysettings.containerName = os.Getenv("DENSIFY_CONTAINER_NAME")
}

// Load Densify settings from Config provided by the user for the Terraform Provider.
func (densifysettings *DensifySettings) LoadConfigSettings(config densifyProviderModel) {
	if !config.DensifyInstance.IsNull() {
		densifysettings.instance = config.DensifyInstance.ValueString()
	}
	if !config.Username.IsNull() {
		densifysettings.username = config.Username.ValueString()
	}
	if !config.Password.IsNull() {
		densifysettings.password = config.Password.ValueString()
	}
	if !config.ApiTimeout.IsNull() {
		densifysettings.timeout = int(config.ApiTimeout.ValueInt64())
	}
	if !config.TechPlatform.IsNull() {
		densifysettings.techPlatform = config.TechPlatform.ValueString()
	}
	if !config.AccountNumber.IsNull() {
		densifysettings.accountNumber = config.AccountNumber.ValueString()
	}
	if !config.AccountName.IsNull() {
		densifysettings.accountName = config.AccountName.ValueString()
	}
	if !config.SystemName.IsNull() {
		densifysettings.systemName = config.SystemName.ValueString()
	}
	if !config.FallbackInstanceType.IsNull() {
		densifysettings.fallbackInstanceType = config.FallbackInstanceType.ValueString()
	}
	if !config.ContinueIfError.IsNull() {
		densifysettings.continueIfError = config.ContinueIfError.ValueBool()
	}

	if !config.K8sCluster.IsNull() {
		densifysettings.cluster = config.K8sCluster.ValueString()
	}
	if !config.K8sNamespace.IsNull() {
		densifysettings.namespace = config.K8sNamespace.ValueString()
	}
	if !config.K8sControllerType.IsNull() {
		densifysettings.controllerType = config.K8sControllerType.ValueString()
	}
	if !config.K8sPodName.IsNull() {
		densifysettings.podName = config.K8sPodName.ValueString()
	}
	if !config.K8sContainerName.IsNull() {
		densifysettings.containerName = config.K8sContainerName.ValueString()
	}

	if !config.K8sFallbackCPUReq.IsNull() {
		densifysettings.fallbackCPUReq = config.K8sFallbackCPUReq.ValueString()
	}
	if !config.K8sFallbackCPULim.IsNull() {
		densifysettings.fallbackCPULim = config.K8sFallbackCPULim.ValueString()
	}
	if !config.K8sFallbackMemReq.IsNull() {
		densifysettings.fallbackMemReq = config.K8sFallbackMemReq.ValueString()
	}
	if !config.K8sFallbackMemLim.IsNull() {
		densifysettings.fallbackMemLim = config.K8sFallbackMemLim.ValueString()
	}
}

// Load Densify settings from Config provided by the user for the Terraform Provider.
func (densifysettings *DensifySettings) ValidateSettings(resp *provider.ConfigureResponse) {
	// If any of the expected configurations are missing, return errors with provider-specific guidance.

	if densifysettings.instance == "" {
		resp.Diagnostics.AddAttributeError(
			path.Root("densify_instance"),
			"Missing Densify API Instance Name",
			"The provider cannot create the Densify API client as there is a missing or empty value for the Densify API Instance Name. "+
				"Set the instance value in the configuration or use the DENSIFY_INSTANCE environment variable. "+
				"If either is already set, ensure the value is not empty.",
		)
	}

	if densifysettings.username == "" {
		resp.Diagnostics.AddAttributeError(
			path.Root("username"),
			"Missing Densify API Username",
			"The provider cannot create the Densify API client as there is a missing or empty value for the Densify API username. "+
				"Set the username value in the configuration or use the DENSIFY_USERNAME environment variable. "+
				"If either is already set, ensure the value is not empty.",
		)
	}

	if densifysettings.password == "" {
		resp.Diagnostics.AddAttributeError(
			path.Root("password"),
			"Missing Densify API Password",
			"The provider cannot create the Densify API client as there is a missing or empty value for the Densify API password. "+
				"Set the password value in the configuration or use the DENSIFY_PASSWORD environment variable. "+
				"If either is already set, ensure the value is not empty.",
		)
	}

	if densifysettings.techPlatform == "" {
		resp.Diagnostics.AddAttributeError(
			path.Root("tech_platform"),
			"Missing Densify API Technology Platform",
			"The provider cannot create the Densify API client as there is a missing or empty value for the Densify API Technology Platform. "+
				"Set the tech_platform value in the configuration or use the DENSIFY_TECH_PLATFORM environment variable. "+
				"If either is already set, ensure the value is not empty.",
		)
	}

	// KUBERNETES/CONTAINERS.
	if strings.ToLower(densifysettings.techPlatform) == "k8s" || strings.ToLower(densifysettings.techPlatform) == "kubernetes" {
		if densifysettings.cluster == "" {
			resp.Diagnostics.AddAttributeError(
				path.Root("cluster"),
				"Missing Kubernetes Cluster Name",
				"The provider cannot create the Densify API client as there is a missing or empty value for the Cluster Name. "+
					"Set the cluster value in the configuration or use the DENSIFY_CLUSTER environment variable. "+
					"If either is already set, ensure the value is not empty.",
			)
		}
		if densifysettings.namespace == "" {
			resp.Diagnostics.AddAttributeError(
				path.Root("namespace"),
				"Missing Kubernetes Namespace",
				"The provider cannot create the Densify API client as there is a missing or empty value for the Namespace. "+
					"Set the namespace value in the configuration or use the DENSIFY_NAMESPACE environment variable. "+
					"If either is already set, ensure the value is not empty.",
			)
		}
		if densifysettings.controllerType == "" {
			resp.Diagnostics.AddAttributeError(
				path.Root("controller_type"),
				"Missing Kubernetes Cluster Name",
				"The provider cannot create the Densify API client as there is a missing or empty value for the Controller Type. "+
					"Set the controller_type value in the configuration. "+
					"If it is already set, ensure the value is not empty.",
			)
		}
		if densifysettings.podName == "" {
			resp.Diagnostics.AddAttributeError(
				path.Root("pod_name"),
				"Missing Kubernetes Pod Name",
				"The provider cannot create the Densify API client as there is a missing or empty value for the Pod Name. "+
					"Set the pod_name value in the configuration. "+
					"If it is already set, ensure the value is not empty.",
			)
		}

	} else {
		if densifysettings.accountName == "" && densifysettings.accountNumber == "" {
			resp.Diagnostics.AddAttributeError(
				path.Root("account_number"),
				"Missing Densify API Account Name/Number",
				"The provider cannot create the Densify API client as there is a missing or empty value for the Densify API Account Number or Account Name. "+
					"Set the account_number value in the configuration or use the DENSIFY_ACCOUNT_NUMBER environment variable. "+
					"If either is already set, ensure the value is not empty.",
			)
		}
		if densifysettings.systemName == "" {
			resp.Diagnostics.AddAttributeError(
				path.Root("system_name"),
				"Missing Densify System Name",
				"The provider cannot create the Densify API client as there is a missing or empty value for the Densify System Name. "+
					"Set the system_name value in the configuration or use the DENSIFY_SYSTEM_NAME environment variable. "+
					"If either is already set, ensure the value is not empty.",
			)
		}
	}
}
