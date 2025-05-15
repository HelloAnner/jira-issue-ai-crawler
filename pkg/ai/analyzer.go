package ai

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/anner/jira-issue-ai-crawler/pkg/jira"
	"github.com/tmc/langchaingo/llms"
	"github.com/tmc/langchaingo/llms/openai"
)

type Analysis struct {
	ModuleCategory       string `json:"module_category"`            // 模块分类
	OLALevel             string `json:"ola_level"`                  // 确定OLA级别
	IsOverdue            bool   `json:"is_overdue"`                 // 是否逾期
	SymptomCategory      string `json:"symptom_category"`           // 症状分类
	SymptomDescription   string `json:"symptom_description"`        // 症状描述
	RootCauseCategory    string `json:"root_cause_category"`        // 根因分类
	RootCauseDescription string `json:"root_cause_description"`     // 根因描述
	SolutionCategory     string `json:"solution_category"`          // 解决方案分类
	SolutionDescription  string `json:"solution_description"`       // 解决方案描述
	IsClosed             bool   `json:"is_closed"`                  // 是否关闭
	IsFixed              bool   `json:"is_fixed"`                   // 是否修复
	DefectType           string `json:"defect_type"`                // 缺陷类型
	TechnicalDebtDesc    string `json:"technical_debt_description"` // 技术债务描述
	IndustrySolution     string `json:"industry_solution"`          // 行业解决方案
	GapAnalysis          string `json:"gap_analysis"`               // 差距分析
}

type Analyzer struct {
	llm llms.LLM
}

func NewAnalyzer(apiKey, model, baseURL string) (*Analyzer, error) {
	llm, err := openai.New(
		openai.WithToken(apiKey),
		openai.WithModel(model),
		openai.WithBaseURL(baseURL),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create OpenAI client: %v", err)
	}

	return &Analyzer{
		llm: llm,
	}, nil
}

func (a *Analyzer) AnalyzeIssue(ctx context.Context, issue *jira.Issue) (*Analysis, error) {
	content := []llms.MessageContent{
		llms.TextParts(
			llms.ChatMessageTypeSystem,
			`你是经验丰富的系统分析师。分析Jira问题详情并提供结构化分析。
你的响应应该符合以下JSON结构：
{
    "module_category": "string - 模块分类",
    "ola_level": "string - 确定OLA级别",
    "is_overdue": "boolean - 是否逾期",
    "symptom_category": "string - 症状分类",
    "symptom_description": "string - 症状描述",
    "root_cause_category": "string - 根因分类",
    "root_cause_description": "string - 根因描述",
    "solution_category": "string - 解决方案分类",
    "solution_description": "string - 解决方案描述",
    "is_closed": "boolean - 是否关闭",
    "is_fixed": "boolean - 是否修复",
    "defect_type": "string - 缺陷类型",
    "technical_debt_description": "string - 技术债务描述",
    "industry_solution": "string - 行业解决方案",
    "gap_analysis": "string - 行业领先产品分析"
}`),
		llms.TextParts(llms.ChatMessageTypeHuman, buildPrompt(issue)),
	}

	result, err := a.llm.GenerateContent(ctx, content, llms.WithMaxTokens(2048))
	if err != nil {
		return nil, fmt.Errorf("OpenAI API error: %v", err)
	}

	var analysis Analysis
	err = json.Unmarshal([]byte(result.Choices[0].Content), &analysis)
	if err != nil {
		return nil, fmt.Errorf("failed to parse AI response: %v", err)
	}

	return &analysis, nil
}

func buildPrompt(issue *jira.Issue) string {
	prompt := fmt.Sprintf(`分析以下Jira问题:
标题: %s
工单号: %s
描述: %s
创建时间: %s
任务解决过程中的评论: %s
任务解决过程中的工作日志: %s
`, issue.Title,
		issue.Key,
		issue.Description,
		issue.CreatedAt.Format("2006-01-02"),
		strings.Join(issue.Comments, "\n"),
		strings.Join(issue.WorkLogs, "\n"),
	)

	if !issue.ResolvedAt.IsZero() {
		prompt += fmt.Sprintf("解决时间: %s\n", issue.ResolvedAt.Format("2006-01-02"))
	}

	if len(issue.Comments) > 0 {
		prompt += "\n评论:\n"
		for _, comment := range issue.Comments {
			prompt += fmt.Sprintf("- %s\n", comment)
		}
	}

	if len(issue.WorkLogs) > 0 {
		prompt += "\n工作日志:\n"
		for _, worklog := range issue.WorkLogs {
			prompt += fmt.Sprintf("- %s\n", worklog)
		}
	}

	return prompt
}
