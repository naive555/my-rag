package rag

import (
	"strconv"
	"strings"
)

func BuildPrompt(chunks []string, q, lang string, maxSentences int) string {
	var sb strings.Builder

	sb.WriteString("SYSTEM:\n")
	sb.WriteString("You are an assistant for a business knowledge system.\n")
	sb.WriteString("Answer ONLY in " + lang + ". ")
	sb.WriteString("DO NOT use any other language.\n")
	sb.WriteString("Be concise. Limit your answer to ")
	sb.WriteString(strconv.Itoa(maxSentences))
	sb.WriteString(" sentences.\n")

	if len(chunks) > 0 {
		sb.WriteString("CONTEXT:\n")
		sb.WriteString("Use ONLY the provided context. ")
		sb.WriteString("If the answer is not in the context, say \"I don't know\".\n")
		sb.WriteString(strings.Join(chunks, "\n---\n"))
		sb.WriteString("\n")
	} else {
		sb.WriteString("If you are unsure, say \"I don't know\".\n")
	}

	sb.WriteString("USER:\n")
	sb.WriteString(q)

	return sb.String()
}

func BuildPromptWithHistory(
	chunks []string,
	history []string,
	q string,
	lang string,
	maxSentences int,
) string {
	var sb strings.Builder

	sb.WriteString("SYSTEM:\n")
	sb.WriteString("You are an assistant for a business knowledge system.\n")
	sb.WriteString("Answer ONLY in " + lang + ". ")
	sb.WriteString("DO NOT use any other language.\n")
	sb.WriteString("Be concise. Limit your answer to ")
	sb.WriteString(strconv.Itoa(maxSentences))
	sb.WriteString(" sentences.\n")

	if len(chunks) > 0 {
		sb.WriteString("CONTEXT:\n")
		sb.WriteString("Use ONLY the provided context. ")
		sb.WriteString("If the answer is not in the context, say \"I don't know\".\n")
		sb.WriteString(strings.Join(chunks, "\n---\n"))
		sb.WriteString("\n")
	} else {
		sb.WriteString("If you are unsure, say \"I don't know\".\n")
	}

	if len(history) > 0 {
		sb.WriteString("CHAT HISTORY:\n")
		sb.WriteString(strings.Join(history, "\n"))
		sb.WriteString("\n\n")
	}

	sb.WriteString("USER:\n")
	sb.WriteString(q)

	return sb.String()
}

// Simple Prompt
func BuildSimplePrompt(chunks []string, q, lang string) string {
	var sb strings.Builder

	sb.WriteString("SYSTEM:\n")
	sb.WriteString("Answer in " + lang + ". ")
	sb.WriteString("Be concise. Max 3 sentences. \n")

	if len(chunks) > 0 {
		sb.WriteString("SYSTEM:\nYou answer only from the context. If missing, say you don't know.\n\n" +
			"CONTEXT:\n" + strings.Join(chunks, "\n---\n") +
			"\n\nUSER:\n" + q)
	} else {
		sb.WriteString("You are a helpful assistant.\n\nUser:\n" + q)
	}

	return sb.String()
}
