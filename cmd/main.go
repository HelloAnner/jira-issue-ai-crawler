package main

import (
	"context"
	"log"
	"time"

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

	// Start the sync loop
	ticker := time.NewTicker(time.Duration(cfg.Sync.Interval) * time.Minute)
	defer ticker.Stop()

	log.Printf("Starting Jira issue analysis service. Sync interval: %d minutes", cfg.Sync.Interval)

	// Run first sync immediately
	if err := syncIssues(jiraClient, analyzer, repo); err != nil {
		log.Printf("Initial sync failed: %v", err)
	}

	// Continue syncing periodically
	for range ticker.C {
		if err := syncIssues(jiraClient, analyzer, repo); err != nil {
			log.Printf("Sync failed: %v", err)
		}
	}
}

func syncIssues(jiraClient *jira.Client, analyzer *ai.Analyzer, repo *database.Repository) error {
	issues, err := jiraClient.GetIssues("project = YourProject AND updated >= -7d")
	if err != nil {
		return err
	}

	log.Printf("Found %d issues to analyze", len(issues))

	ctx := context.Background()
	for _, issue := range issues {
		existing, err := repo.FindByJiraKey(issue.Key)
		if err == nil && existing != nil {
			// Skip if we already have this issue and it hasn't been resolved yet
			if existing.ResolvedAt.IsZero() && issue.ResolvedAt.IsZero() {
				continue
			}
		}

		analysis, err := analyzer.AnalyzeIssue(ctx, &issue)
		if err != nil {
			log.Printf("Failed to analyze issue %s: %v", issue.Key, err)
			continue
		}

		dbAnalysis := &database.IssueAnalysis{
			JiraURL:              issue.URL,                     // 链接
			JiraKey:              issue.Key,                     // 工单号
			IssueTitle:           issue.Title,                   // 工单标题
			CreatedAt:            issue.CreatedAt,               // 创建时间
			ResolvedAt:           issue.ResolvedAt,              // 解决时间
			ResponsibleDev:       issue.Dev,                     // 开发负责人
			ResponsibleQA:        issue.QA,                      // 测试负责人
			ModuleCategory:       analysis.ModuleCategory,       // 模块分类
			OLALevel:             analysis.OLALevel,             // OLA等级
			IsOverdue:            analysis.IsOverdue,            // 是否逾期
			SymptomCategory:      analysis.SymptomCategory,      // 症状分类
			SymptomDescription:   analysis.SymptomDescription,   // 症状描述
			RootCauseCategory:    analysis.RootCauseCategory,    // 根因分类
			RootCauseDescription: analysis.RootCauseDescription, // 根因描述
			SolutionCategory:     analysis.SolutionCategory,     // 解决方案分类
			SolutionDescription:  analysis.SolutionDescription,  // 解决方案描述
			IsClosed:             analysis.IsClosed,             // 是否关闭
			IsFixed:              analysis.IsFixed,              // 是否修复
			DefectType:           analysis.DefectType,           // 缺陷类型
			TechnicalDebtDesc:    analysis.TechnicalDebtDesc,    // 技术债务描述
			IndustrySolution:     analysis.IndustrySolution,     // 行业解决方案
			GapAnalysis:          analysis.GapAnalysis,          // 差距分析
		}

		if err := repo.CreateOrUpdate(dbAnalysis); err != nil {
			log.Printf("Failed to save analysis for issue %s: %v", issue.Key, err)
			continue
		}

		log.Printf("Successfully analyzed and saved issue %s", issue.Key)
	}

	return nil
}
