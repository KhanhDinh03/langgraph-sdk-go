package http

func containsTextEventStream(contentType string) bool {
	return contentType == "text/event-stream" ||
		contentType == "text/event-stream; charset=utf-8"
}
