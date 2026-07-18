package models

// ProjectExportPayload defines the JSON structure for exporting/importing a project
type ProjectExportPayload struct {
	Project  Project    `json:"project"`
	Workflow *Workflow  `json:"workflow,omitempty"`
	Todos    []Todo     `json:"todos"`
	Wikis    []WikiPage `json:"wikis"`
}
