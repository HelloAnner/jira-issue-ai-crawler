# Jira Issue AI Analysis

This project automatically analyzes Jira issues using AI to extract structured insights and stores them in a database for further analysis. It periodically syncs with Jira, processes new and updated issues, and maintains a structured database of issue analyses.

## Features

- Periodic synchronization with Jira
- AI-powered analysis of issue content
- Structured data extraction for:
  - Module categorization
  - OLA level assessment
  - Issue symptoms and root causes
  - Solution analysis
  - Technical debt evaluation
  - Industry comparison
- MySQL database storage for analysis results

## Configuration

Create a `config.yaml` file with the following structure:

```yaml
jira:
  url: "https://your-jira-instance.com"
  username: "your-username"
  password: "your-password"

ai:
  api_key: "your-openai-api-key"
  model: "gpt-3.5-turbo"
  temperature: 0.5

db:
  host: "localhost"
  port: 3306
  username: "root"
  password: "your-db-password"
  database: "jira"

sync:
  interval: 10 # minutes
```

## Database Schema

The application creates a table called `issue_analysis` with the following structure:

- `id`: Primary key
- `jira_url`: URL to the Jira issue
- `jira_key`: Unique Jira issue key
- `module_category`: Module classification
- `ola_level`: Service level classification
- `issue_title`: Issue title
- `created_at`: Issue creation time
- `resolved_at`: Issue resolution time
- `is_overdue`: Whether the issue is overdue
- `responsible_dev`: Responsible developer
- `responsible_qa`: Responsible QA
- `symptom_category`: Issue symptom category
- `symptom_description`: Detailed symptom description
- `root_cause_category`: Root cause category
- `root_cause_description`: Detailed root cause description
- `solution_category`: Solution approach category
- `solution_description`: Detailed solution description
- `is_closed`: Whether the issue is closed
- `is_fixed`: Whether the issue is fixed
- `defect_type`: Type of defect
- `technical_debt_description`: Technical debt analysis
- `industry_solution`: Industry standard solutions
- `gap_analysis`: Gap analysis with industry leading products

## Setup

1. Install Go 1.21 or later
2. Clone the repository
3. Copy `config.yaml.example` to `config.yaml` and update with your settings
4. Create the MySQL database
5. Run the application:

```bash
go run cmd/main.go
```

## Project Structure

```
.
├── cmd/
│   └── main.go           # Main application entry point
├── pkg/
│   ├── ai/              # AI analysis package
│   ├── config/          # Configuration handling
│   ├── database/        # Database operations
│   └── jira/            # Jira API client
├── config.yaml          # Configuration file
├── go.mod              # Go module file
└── README.md           # This file
```

## Dependencies

- github.com/andygrunwald/go-jira: Jira API client
- github.com/sashabaranov/go-openai: OpenAI API client
- gopkg.in/yaml.v2: YAML configuration parser
- gorm.io/gorm: ORM library
- gorm.io/driver/mysql: MySQL driver for GORM
