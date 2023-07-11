package editor

import "github.com/AlecAivazis/survey/v2"

func OpenInEditor(body string, filePattern string) (string, error) {
	prompt := &survey.Editor{
		Default:  body,
		FileName: filePattern,
		AppendDefault: true,
		HideDefault: true,
	}
	content := ""

	err := survey.AskOne(prompt, &content)
	return content, err
}
