package main

import (
	"context"
	"fmt"
	"log"
	"strings"
	"sync"

	"github.com/anner/jira-issue-ai-crawler/pkg/ai"
	"github.com/anner/jira-issue-ai-crawler/pkg/config"
	"github.com/anner/jira-issue-ai-crawler/pkg/database"
	"github.com/anner/jira-issue-ai-crawler/pkg/jira"
)

func main() {
	// Load configuration
	cfg, err := config.LoadConfig("config.yaml")
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Initialize database connection
	db, err := database.NewConnection(
		cfg.DB.Host,
		cfg.DB.Port,
		cfg.DB.Username,
		cfg.DB.Password,
		cfg.DB.Database,
	)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// Initialize repositories and services
	repo := database.NewRepository(db)
	jiraClient, err := jira.NewClient(cfg.Jira.URL, cfg.Jira.Username, cfg.Jira.Password)
	if err != nil {
		log.Fatalf("Failed to create Jira client: %v", err)
	}
	analyzer, err := ai.NewAnalyzer(cfg.AI.APIKey, cfg.AI.Model, cfg.AI.BaseURL)
	if err != nil {
		log.Fatalf("Failed to create AI analyzer: %v", err)
	}

	log.Printf("Starting Jira issue analysis service. Sync interval: %d minutes", cfg.Sync.Interval)

	// Run first sync immediately
	if err := syncIssues(jiraClient, analyzer, repo, cfg); err != nil {
		log.Printf("Initial sync failed: %v", err)
	}
}

func syncIssues(jiraClient *jira.Client, analyzer *ai.Analyzer, repo *database.Repository, cfg *config.Config) error {
	issues, err := jiraClient.GetIssues(cfg.Jira.JQL)
	if err != nil {
		return err
	}

	log.Printf("Found %d issues to analyze", len(issues))

	ctx := context.Background()

	// 创建任务通道
	taskChan := make(chan jira.Issue, len(issues))
	// 创建错误通道
	errChan := make(chan error, len(issues))
	// 创建等待组
	var wg sync.WaitGroup

	// 启动10个worker
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func(workerID int) {
			defer wg.Done()
			for issue := range taskChan {
				// 检查issue是否已存在
				existing, err := repo.FindByJiraKey(issue.Key)
				if err == nil && existing != nil {
					log.Printf("Worker %d: issue %s already exists", workerID, issue.Key)
					continue
				}

				log.Printf("Worker %d: Analyzing issue %s", workerID, issue.Key)

				analysis, err := analyzer.AnalyzeIssue(ctx, &issue)
				if err != nil {
					errChan <- fmt.Errorf("Worker %d failed to analyze issue %s: %v", workerID, issue.Key, err)
					continue
				}

				dbAnalysis := &database.ConsumerIssue{
					JiraURL:              cfg.Jira.URL + "/browse/" + issue.Key,
					JiraKey:              issue.Key,
					IssueTitle:           issue.Title,
					ResponsibleDev:       issue.Dev,
					ResponsibleQA:        issue.QA,
					ModuleCategory:       analysis.ModuleCategory,
					SymptomCategory:      analysis.SymptomCategory,
					SymptomDescription:   analysis.SymptomDescription,
					RootCauseCategory:    analysis.RootCauseCategory,
					RootCauseDescription: analysis.RootCauseDescription,
					SolutionCategory:     analysis.SolutionCategory,
					SolutionDescription:  analysis.SolutionDescription,
					IsClosed:             analysis.IsClosed,
					IsFixed:              analysis.IsFixed,
					DefectType:           analysis.DefectType,
					TechnicalDebtDesc:    analysis.TechnicalDebtDesc,
					IndustrySolution:     analysis.IndustrySolution,
					GapAnalysis:          analysis.GapAnalysis,
					OriginalDescription:  issue.Description,
					OriginalComments:     strings.Join(issue.Comments, "\n"),
					OriginalWorkLogs:     strings.Join(issue.WorkLogs, "\n"),
				}

				if err := repo.CreateOrUpdate(dbAnalysis); err != nil {
					errChan <- fmt.Errorf("Worker %d failed to save analysis for issue %s: %v", workerID, issue.Key, err)
					continue
				}

				log.Printf("Worker %d: Successfully analyzed and saved issue %s", workerID, issue.Key)
			}
		}(i)
	}

	// 发送任务到通道
	for _, issue := range issues {
		taskChan <- issue
	}
	close(taskChan)

	// 等待所有worker完成
	wg.Wait()
	close(errChan)

	// 收集并处理错误
	var errors []error
	for err := range errChan {
		errors = append(errors, err)
	}

	// 如果有错误，返回第一个错误
	if len(errors) > 0 {
		return fmt.Errorf("encountered %d errors during processing, first error: %v", len(errors), errors[0])
	}

	return nil
}
