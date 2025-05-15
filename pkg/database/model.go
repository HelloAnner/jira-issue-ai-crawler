package database

import (
	"gorm.io/gorm"
)

type ConsumerIssue struct {
	ID                   uint      `gorm:"primarykey"`                                                                     // 主键
	JiraURL              string    `gorm:"column:jira_url;not null"`                                                       // 链接
	JiraKey              string    `gorm:"column:jira_key;type:varchar(191);uniqueIndex:idx_jira_key,length:191;not null"` // 工单号
	ModuleCategory       string    `gorm:"column:module_category;type:varchar(191);"`                                      // 模块分类
	IssueTitle           string    `gorm:"column:issue_title;type:text;type:varchar(256);"`                                // 工单标题
	IsOverdue            bool      `gorm:"column:is_overdue;type:tinyint(1);default:0"`                                    // 是否逾期
	ResponsibleDev       string    `gorm:"column:responsible_dev;type:varchar(191);"`                                      // 开发负责人
	ResponsibleQA        string    `gorm:"column:responsible_qa;type:varchar(191);"`                                       // 测试负责人
	SymptomCategory      string    `gorm:"column:symptom_category;type:varchar(191);"`                                     // 症状分类
	SymptomDescription   string    `gorm:"column:symptom_description;type:text"`                                           // 症状描述
	RootCauseCategory    string    `gorm:"column:root_cause_category;type:varchar(191);"`                                  // 根因分类
	RootCauseDescription string    `gorm:"column:root_cause_description;type:text"`                                        // 根因描述
	SolutionCategory     string    `gorm:"column:solution_category;type:varchar(191);"`                                    // 解决方案分类
	SolutionDescription  string    `gorm:"column:solution_description;type:text"`                                          // 解决方案描述
	IsClosed             bool      `gorm:"column:is_closed;type:tinyint(1);default:0"`                                     // 是否关闭
	IsFixed              bool      `gorm:"column:is_fixed;type:tinyint(1);default:0"`                                      // 是否修复
	DefectType           string    `gorm:"column:defect_type;type:varchar(191);"`                                          // 缺陷类型
	TechnicalDebtDesc    string    `gorm:"column:technical_debt_description;type:text"`                                    // 技术债务描述
	IndustrySolution     string    `gorm:"column:industry_solution;type:text"`                                             // 行业解决方案
	GapAnalysis          string    `gorm:"column:gap_analysis;type:text"`                                                  // 差距分析
	OriginalDescription  string    `gorm:"column:original_description;type:text"`                                          // 原始描述
	OriginalComments     string    `gorm:"column:original_comments;type:text"`                                            // 原始评论
	OriginalWorkLogs     string    `gorm:"column:original_work_logs;type:text"`                                            // 原始工作日志
}

type Repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) CreateOrUpdate(analysis *ConsumerIssue) error {
	return r.db.Save(analysis).Error
}

func (r *Repository) FindByJiraKey(jiraKey string) (*ConsumerIssue, error) {
	var analysis ConsumerIssue
	err := r.db.Where("jira_key = ?", jiraKey).First(&analysis).Error
	if err != nil {
		return nil, err
	}
	return &analysis, nil
}

func (r *Repository) ListAll() ([]ConsumerIssue, error) {
	var analyses []ConsumerIssue
	err := r.db.Find(&analyses).Error
	if err != nil {
		return nil, err
	}
	return analyses, nil
}
