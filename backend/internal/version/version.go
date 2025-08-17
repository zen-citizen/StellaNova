package version

var (
	buildTime = "unknown"
	commitSHA = "unknown"
)

func GetBuildTime() string {
	return buildTime
}

func GetCommitSHA() string {
	return commitSHA
}
