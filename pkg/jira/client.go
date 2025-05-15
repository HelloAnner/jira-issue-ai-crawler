package jira

import (
	"time"

	"github.com/andygrunwald/go-jira"
)

type Issue struct {
	URL           string // 链接
	Key           string // 工单号
	Title         string // 工单标题
	CreatedAt     time.Time // 创建时间
	ResolvedAt    time.Time // 解决时间
	Dev           string // 开发负责人
	QA            string // 测试负责人
	Description   string // 工单描述	
	Comments      []string // 评论
	WorkLogs      []string // 工作日志
	CustomFields  map[string]interface{} // 自定义字段
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