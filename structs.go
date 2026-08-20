package main

type Ollama struct {
	Model  string `json:"model"`
	System string `json:"system"`
	Prompt string `json:"prompt"`
	Stream bool   `json:"stream"`
	Think  bool   `json:"bool"`
}

type OllamaGenerateResponse struct {
	Response string `json:"response"`
}
