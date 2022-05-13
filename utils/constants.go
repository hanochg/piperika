package utils

const (
	ArtifactoryPipelinesSourceId         = 6
	HttpClientCtxKey                     = "httpClient"
	BaseUiUrl                            = "baseUiUrl"
	BranchName                           = "branchName"
	ConfigCtxKey                         = "configurations"
	ForceFlag                            = "forceFlag"
	InfraReportLinksStep                 = "release_process_links"
	InfraReportEnvVersionSuffix          = "_VERSION"
	InfraReportEnvAdHocReleaseBranchName = "AD_HOC_RELEASE_BRANCH_NAME"
	InfraReportReleasePipeSuffix         = "_release"
	InfraReportBuildPipeSuffix           = "_build"
	InfraReportPostReleasePipeSuffix     = "_post_release"
)

var InfraReportServices = [...]string{"access", "metadata", "event", "router", "jfconnect", "mission_control", "integration"}
