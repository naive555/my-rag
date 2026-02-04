package rag

import "strings"

func BuildPrompt(chunks []string, q string) string {
	var sb strings.Builder

	sb.WriteString("SYSTEM:\n")
	sb.WriteString("Answer using only the context. ")
	sb.WriteString("Be concise. Max 3 sentences. ")
	sb.WriteString("Do not explain your reasoning. ")
	sb.WriteString("If not found, say \"I don't know\".\n\n")

	sb.WriteString("CONTEXT:\n")
	if len(chunks) == 0 {
		sb.WriteString("No context.\n")
	} else {
		sb.WriteString(strings.Join(chunks, "\n"))
	}
	sb.WriteString("\n\n")

	sb.WriteString("USER:\n")
	sb.WriteString(q)

	return sb.String()
}

// func BuildPrompt(chunks []string, q string) string {
// 	var prompt string

// 	if len(chunks) > 0 {
// 		prompt = "SYSTEM:\nYou answer only from the context. If missing, say you don't know.\n\n" +
// 			"CONTEXT:\n" + strings.Join(chunks, "\n---\n") +
// 			"\n\nUSER:\n" + q
// 	} else {
// 		prompt = "You are a helpful assistant.\n\nUser:\n" + q
// 	}

// 	return prompt
// }
