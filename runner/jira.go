package runner

import (
	"strings"
	"tix/creator"
	"tix/creator/jira"
	"tix/logger"
	"tix/md"
	"tix/runner/env"
	"tix/runner/field"
	"tix/settings"
	"tix/ticket"
)

type JiraRunner struct {
	envMap   map[string]string
	settings *settings.Settings
}

func NewJiraRunner(envMap map[string]string, settings *settings.Settings) TixRunner {
	return &JiraRunner{
		envMap: envMap,
		settings: settings,
	}
}

func (r JiraRunner) Run(markdownData []byte) error {
	if len(r.settings.Jira.Url) == 0 {
		return nil
	}

	err := env.CheckJiraEnvironment(r.envMap)
	if err != nil {
		return err
	}

	tickets, err := r.parseMarkdown(markdownData)

	if err != nil {
		return err
	}

	r.jiraCreator().CreateTickets(tickets)

	return nil
}

func (r JiraRunner) parseMarkdown(markdownData []byte) ([]*ticket.Ticket, error) {
	fieldState := field.JiraFieldState(*r.settings)
	markdownParser := md.NewParser(fieldState)
	return markdownParser.Parse(markdownData)
}

func (r JiraRunner) jiraCreator() creator.TicketCreator {
	api := r.createJiraApi()
	if r.settings.Jira.NoEpics {
		return jira.NewCreatorWithoutEpics(api)
	} else {
		return jira.NewCreatorWithEpics(api)
	}
}

func (r JiraRunner) createJiraApi() jira.Api {
	return jira.NewApi(r.envMap[env.JiraUsername], r.envMap[env.JiraApiToken], r.settings.Jira.Url)
}

func (r JiraRunner) DryRun(markdownData []byte) error {
	if len(r.settings.Jira.Url) == 0 {
		return nil
	}

	err := env.CheckJiraEnvironment(r.envMap)
	if err != nil {
		return err
	}

	tickets, err := r.parseMarkdown(markdownData)

	if err != nil {
		return err
	}

	var builder strings.Builder
	r.dryRunCreator(&builder).CreateTickets(tickets)
	logger.Message(builder.String())

	return nil
}

func (r JiraRunner) dryRunCreator(builder *strings.Builder) creator.TicketCreator {
	if r.settings.Jira.NoEpics {
		return jira.NewDryRunCreatorWithoutEpics(builder)
	} else {
		return jira.NewDryRunCreatorWithEpics(builder)
	}
}