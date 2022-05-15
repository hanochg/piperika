package utils

const (
	ArtifactoryPipelinesSourceId         = 6
	HttpClientCtxKey                     = "httpClient"
	BaseUiUrl                            = "baseUiUrl"
	BranchName                           = "branchName"
	ConfigCtxKey                         = "configurations"
	ForceFlag                            = "forceFlag"
	DefaultProject                       = "default"
	DefaultProjectId                     = 1
	InfraReportLinksStep                 = "release_process_links"
	InfraReportEnvVersionSuffix          = "_VERSION"
	InfraReportEnvAdHocReleaseBranchName = "AD_HOC_RELEASE_BRANCH_NAME"
	InfraReportReleasePipeSuffix         = "_release"
	InfraReportBuildPipeSuffix           = "_build"
	InfraReportPostReleasePipeSuffix     = "_post_release"
)

var InfraReportServiceProject = map[string]string{
	"access":          "Access",
	"metadata":        DefaultProject,
	"event":           DefaultProject,
	"router":          DefaultProject,
	"jfconnect":       DefaultProject,
	"mission_control": DefaultProject,
	"integration":     DefaultProject,
}
