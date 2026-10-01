package ollama

import "strings"

const (
	_thinkOpen  = "<think>"
	_thinkClose = "</think>"
)

func ParseThinkingAndReply(msg string) (thinking, reply string) {
	startIdx := strings.Index(msg, _thinkOpen)
	if startIdx == -1 {
		reply = msg
		return
	}
	endIdx := strings.Index(msg, _thinkClose) + len(_thinkClose)
	if endIdx == -1 {
		reply = msg
		return
	}
	thinking = msg[startIdx:endIdx]
	if len(msg) > endIdx {
		reply = strings.TrimSpace(msg[endIdx:])
	}
	return
}
