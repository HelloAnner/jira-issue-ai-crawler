package database

import (
	"time"

	"gorm.io/gorm"
)

type IssueAnalysis struct {
	ID                   uint       `gorm:"primarykey"`                                  // 主键
	JiraURL              string     `gorm:"column:jira_url;not null"`                    // 链接
	JiraKey              string     `gorm:"column:jira_key;uniqueIndex;not null"`        // 工单号
	ModuleCategory       string     `gorm:"column:module_category"`                      // 模块分类
	OLALevel             string     `gorm:"column:ola_level"`                            // OLA等级
	IssueTitle           string     `gorm:"column:issue_title;type:text"`                // 工单标题
	CreatedAt            time.Time  `gorm:"column:created_at"`                           // 创建时间
	ResolvedAt           time.Time `gorm:"column:resolved_at"`                          // 解决时间
	IsOverdue            bool       `gorm:"column:is_overdue"`                           // 是否逾期
	ResponsibleDev       string     `gorm:"column:responsible_dev"`                      // 开发负责人
	ResponsibleQA        string     `gorm:"column:responsible_qa"`                       // 测试负责人
	SymptomCategory      string     `gorm:"column:symptom_category"`                     // 症状分类
	SymptomDescription   string     `gorm:"column:symptom_description;type:text"`        // 症状描述
	RootCauseCategory    string     `gorm:"column:root_cause_category"`                  // 根因分类
	RootCauseDescription string     `gorm:"column:root_cause_description;type:text"`     // 根因描述
	SolutionCategory     string     `gorm:"column:solution_category"`                    // 解决方案分类
	SolutionDescription  string     `gorm:"column:solution_description;type:text"`       // 解决方案描述
	IsClosed             bool       `gorm:"column:is_closed"`                            // 是否关闭
	IsFixed              bool       `gorm:"column:is_fixed"`                             // 是否修复
	DefectType           string     `gorm:"column:defect_type"`                          // 缺陷类型
	TechnicalDebtDesc    string     `gorm:"column:technical_debt_description;type:text"` // 技术债务描述
	IndustrySolution     string     `gorm:"column:industry_solution;type:text"`          // 行业解决方案
	GapAnalysis          string     `gorm:"column:gap_analysis;type:text"`               // 差距分析
	CreatedTime          time.Time  `gorm:"column:created_time;autoCreateTime"`          // 创建时间
	UpdatedTime          time.Time  `gorm:"column:updated_time;autoUpdateTime"`          // 更新时间
}

type Repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) CreateOrUpdate(analysis *IssueAnalysis) error {
	return r.db.Save(analysis).Error
}

func (r *Repository) FindByJiraKey(jiraKey string) (*IssueAnalysis, error) {
	var analysis IssueAnalysis
	err := r.db.Where("jira_key = ?", jiraKey).First(&analysis).Error
	if err != nil {
		return nil, err
	}
	return &analysis, nil
}

func (r *Repository) ListAll() ([]IssueAnalysis, error) {
	var analyses []IssueAnalysis
	err := r.db.Find(&analyses).Error
	if err != nil {
		return nil, err
	}
	return analyses, nil
}
