CREATE TABLE IF NOT EXISTS issue_analysis (
    id BIGINT AUTO_INCREMENT PRIMARY KEY, -- 主键
    jira_url VARCHAR(255) NOT NULL, -- 链接
    jira_key VARCHAR(50) NOT NULL, -- 工单号
    module_category VARCHAR(100), -- 模块分类
    ola_level VARCHAR(50), -- OLA等级
    issue_title TEXT, -- 工单标题
    created_at DATETIME, -- 创建时间
    resolved_at DATETIME, -- 解决时间
    is_overdue BOOLEAN, -- 是否逾期
    responsible_dev VARCHAR(100), -- 开发负责人
    responsible_qa VARCHAR(100), -- 测试负责人
    symptom_category VARCHAR(100), -- 症状分类
    symptom_description TEXT, -- 症状描述
    root_cause_category VARCHAR(100), -- 根因分类
    root_cause_description TEXT, -- 根因描述
    solution_category VARCHAR(3000), -- 解决方案分类
    solution_description TEXT, -- 解决方案描述
    is_closed BOOLEAN, -- 是否关闭
    is_fixed BOOLEAN, -- 是否修复
    defect_type VARCHAR(100), -- 缺陷类型
    technical_debt_description TEXT, -- 技术债务描述
    industry_solution TEXT, -- 行业解决方案
    gap_analysis TEXT, -- 差距分析
    created_time DATETIME DEFAULT CURRENT_TIMESTAMP, -- 创建时间
    updated_time DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP, -- 更新时间
    UNIQUE KEY unique_jira_key (jira_key) -- 工单号唯一索引
); 