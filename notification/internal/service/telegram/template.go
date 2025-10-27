package telegram

type TemplateRenderer interface {
	Render(templateName string, data interface{}) (string, error)
}
