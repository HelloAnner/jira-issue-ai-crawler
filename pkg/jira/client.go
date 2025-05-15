package jira

import (
	"time"

	"github.com/andygrunwald/go-jira"
)

type Issue struct {
	URL           string
	Key           string
	Title         string
	CreatedAt     time.Time
	ResolvedAt    time.Time
	Dev           string
	QA            string
	Description   string
	Comments      []string
	WorkLogs      []string
	CustomFields  map[string]interface{}
}

type Client struct {
	client *jira.Client
}

func NewClient(url, username, password string) (*Client, error) {
	tp := jira.BasicAuthTransport{
		Username: username,
		Password: password,
	}

	client, err := jira.NewClient(tp.Client(), url)
	if err != nil {
		return nil, err
	}

	return &Client{client: client}, nil
}

func (c *Client) GetIssues(jql string) ([]Issue, error) {
	var issues []Issue
	
	options := &jira.SearchOptions{
		MaxResults: 100,
		Fields: []string{
			"summary",
			"description",
			"created",
			"resolutiondate",
			"assignee",
			"customfield_*",
			"comment",
			"worklog",
		},
	}

	chunk, resp, err := c.client.Issue.Search(jql, options)
	if err != nil {
		return nil, err
	}

	for _, item := range chunk {
		issue := Issue{
			URL:          item.Self,
			Key:          item.Key,
			Title:        item.Fields.Summary,
			Description:  item.Fields.Description,
			CustomFields: make(map[string]interface{}),
		}

		// Parse created time
		createdTime, err := time.Parse("2006-01-02T15:04:05.999-0700", item.Fields.Created)
		if err == nil {
			issue.CreatedAt = createdTime
		}

		// Parse resolution time
		if resolutionDate := item.Fields.Resolutiondate.String(); resolutionDate != "" {
			if resolvedTime, err := time.Parse("2006-01-02T15:04:05.999-0700", resolutionDate); err == nil {
				issue.ResolvedAt = &resolvedTime
			}
		}

		if item.Fields.Assignee != nil {
			issue.Dev = item.Fields.Assignee.DisplayName
		}

		// Extract comments
		if item.Fields.Comments != nil {
			for _, comment := range item.Fields.Comments.Comments {
				issue.Comments = append(issue.Comments, comment.Body)
			}
		}

		// Extract worklogs
		if item.Fields.Worklog != nil {
			for _, worklog := range item.Fields.Worklog.Worklogs {
				issue.WorkLogs = append(issue.WorkLogs, worklog.Comment)
			}
		}

		// Extract custom fields
		for key, value := range item.Fields.Unknowns {
			issue.CustomFields[key] = value
		}

		issues = append(issues, issue)
	}

	for resp.Total > len(issues) {
		options.StartAt = len(issues)
		chunk, resp, err = c.client.Issue.Search(jql, options)
		if err != nil {
			return nil, err
		}

		for _, item := range chunk {
			issue := Issue{
				URL:          item.Self,
				Key:          item.Key,
				Title:        item.Fields.Summary,
				Description:  item.Fields.Description,
				CustomFields: make(map[string]interface{}),
			}

			// Parse created time
			createdTime, err := time.Parse("2006-01-02T15:04:05.999-0700", item.Fields.Created)
			if err == nil {
				issue.CreatedAt = createdTime
			}

			// Parse resolution time
			if resolutionDate := item.Fields.Resolutiondate; resolutionDate != "" {
				if resolvedTime, err := time.Parse("2006-01-02T15:04:05.999-0700", resolutionDate); err == nil {
					issue.ResolvedAt = &resolvedTime
				}
			}

			if item.Fields.Assignee != nil {
				issue.Dev = item.Fields.Assignee.DisplayName
			}

			// Extract comments
			if item.Fields.Comments != nil {
				for _, comment := range item.Fields.Comments.Comments {
					issue.Comments = append(issue.Comments, comment.Body)
				}
			}

			// Extract worklogs
			if item.Fields.Worklog != nil {
				for _, worklog := range item.Fields.Worklog.Worklogs {
					issue.WorkLogs = append(issue.WorkLogs, worklog.Comment)
				}
			}

			// Extract custom fields
			for key, value := range item.Fields.Unknowns {
				issue.CustomFields[key] = value
			}

			issues = append(issues, issue)
		}
	}

	return issues, nil
} 