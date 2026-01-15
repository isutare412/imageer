package valkey

func imageProcessDoneChannel(prefix string, imageID string) string {
	return prefix + imageID
}
