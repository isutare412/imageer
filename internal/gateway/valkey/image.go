package valkey

func imageUploadDoneChannel(prefix string, imageID string) string {
	return prefix + imageID
}

func imageProcessDoneChannel(prefix string, imageID string) string {
	return prefix + imageID
}
