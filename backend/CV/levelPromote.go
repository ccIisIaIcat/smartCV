package CV

import (
	"encoding/json"
	"fmt"
)

// PromptLevel 定义简历生成的不同级别
type PromptLevel int

const (
	// Honest 诚实级别：完全基于真实经历
	Honest PromptLevel = iota
	// Moderate 适中级别：适当美化，但在合理范围内
	Moderate
	// Aggressive 进取级别：大胆包装，除硬性条件外都迎合要求
	Aggressive
	// Creative 创造级别：在Moderate基础上增加创意性和独特性
	Creative
	// Expert 专家级别：在Honest基础上突出专业性和深度
	Expert
)

// UnmarshalJSON 实现json.Unmarshaler接口
func (p *PromptLevel) UnmarshalJSON(data []byte) error {
	// 尝试解析为数字
	var number float64
	if err := json.Unmarshal(data, &number); err == nil {
		*p = PromptLevel(int(number))
		if *p >= Honest && *p <= Expert {
			return nil
		}
		return fmt.Errorf("无效的简历级别值: %v", number)
	}

	// 如果不是数字，尝试解析为字符串
	var str string
	if err := json.Unmarshal(data, &str); err != nil {
		return err
	}

	// 将字符串映射到对应的级别
	switch str {
	case "Honest":
		*p = Honest
	case "Moderate":
		*p = Moderate
	case "Aggressive":
		*p = Aggressive
	case "Creative":
		*p = Creative
	case "Expert":
		*p = Expert
	default:
		return fmt.Errorf("无效的简历级别名称: %s", str)
	}
	return nil
}

// ResumePromptGenerator 简历提示生成器
type ResumePromptGenerator struct{}

// NewResumePromptGenerator 创建新的简历提示生成器
func NewResumePromptGenerator() *ResumePromptGenerator {
	return &ResumePromptGenerator{}
}

// GenerateHonestPrompt 生成诚实级别的提示
func (g *ResumePromptGenerator) GenerateHonestPrompt(introduction, requirements string) string {
	return fmt.Sprintf(`请作为一个专业的简历顾问，帮我生成一份完全基于真实经历的LaTeX格式简历。

要求：
1. 严格遵守以下原则：
   - 只使用个人介绍中提供的真实信息
   - 重点突出与职位要求直接相关的经历和技能
   - 保持完全的真实性，不进行任何美化或夸大
   - 按照职位要求的优先级组织内容
   
2. 简历格式要求：
   - 必须包含以下LaTeX包用于支持中文：
     \usepackage{xeCJK}
     \setCJKmainfont{KaiTi}
   - 不要使用 \photo 命令或任何图片
   - 使用清晰的层次结构
   - 采用专业的排版格式
   - 突出重点内容

3. 内容筛选原则：
   - 优先选择与职位要求高度匹配的经历
   - 使用具体的数据和成果支持经历
   - 保持内容的真实性和可验证性

个人介绍：
%s

职位要求：
%s

请生成一份完整的LaTeX简历代码。只返回LaTeX代码，不需要其他解释。代码必须以 \documentclass 开始，以 \end{document} 结束。`, introduction, requirements)
}

// GenerateModeratePrompt 生成适中级别的提示
func (g *ResumePromptGenerator) GenerateModeratePrompt(introduction, requirements string) string {
	return fmt.Sprintf(`请作为一个专业的简历顾问，帮我生成一份适度优化的LaTeX格式简历。

要求：
1. 优化原则：
   - 基于真实经历进行合理的延伸和美化
   - 可以适当扩展项目规模和影响力
   - 技能描述可以适度提升，但要确保面试时能够应对
   - 保持经历的基本真实性，只在细节处进行优化
   
2. 简历格式要求：
   - 不要使用 \photo 命令或任何图片
   - 必须包含以下LaTeX包用于支持中文：
     \usepackage{xeCJK}
     \setCJKmainfont{KaiTi}
   - 使用 classic 样式
   - 使用蓝色主题
   - 使用现代化的排版风格
   - 重点突出优化后的内容

3. 内容优化策略：
   - 将现有经历与职位要求巧妙对应
   - 适当添加相关的技能点
   - 扩展项目影响力和规模，但保持在合理范围
   - 确保所有优化内容在面试中可以自圆其说

个人介绍：
%s

职位要求：
%s

请生成一份完整的LaTeX简历代码。只返回LaTeX代码，不需要其他解释。代码必须以 \documentclass 开始，以 \end{document} 结束。`, introduction, requirements)
}

// GenerateAggressivePrompt 生成进取级别的提示
func (g *ResumePromptGenerator) GenerateAggressivePrompt(introduction, requirements string) string {
	return fmt.Sprintf(`请作为一个善于包装的职业顾问，帮我生成一份极具竞争力的LaTeX格式简历。

要求：
1. 包装原则：
   - 只保持身份证可查的基本信息（年龄、学历）真实
   - 其他所有内容都可以进行创造性改编和包装
   - 项目经验可以完全重构，以完美匹配职位需求
   - 技能水平可以大胆提升，以满足所有要求
   - 工作成果可以适度夸大，突出影响力
   
2. 简历格式要求：
   - 不要使用 \photo 命令或任何图片
   - 必须包含以下LaTeX包用于支持中文：
     \usepackage{xeCJK}
     \setCJKmainfont{KaiTi}
   - 使用蓝色主题
   - 使用引人注目的排版风格
   - 重点突出所有亮点内容

3. 内容优化策略：
   - 根据职位要求，创造完美匹配的项目经历
   - 大胆提升技术栈的熟练度和应用深度
   - 适度夸大项目规模和影响力（如：将百万级改为千万级）
   - 添加热门但相关的技术关键词
   - 虚构一些含金量高的证书和获奖经历
   - 适当添加行业影响力的例子（演讲、著作等）
   - 不局限于原有经历，可以自由发挥

4. 注意事项：
   - 确保所有描述听起来合理可信
   - 避免使用可以被轻易验证的具体公司名称
   - 模糊化具体的时间点和数据
   - 使用行业通用的术语和描述

个人介绍：
%s

职位要求：
%s

请生成一份完整的LaTeX简历代码。只返回LaTeX代码，不需要其他解释。代码必须以 \documentclass 开始，以 \end{document} 结束。`, introduction, requirements)
}

// GenerateCreativePrompt 生成创意级别的提示
func (g *ResumePromptGenerator) GenerateCreativePrompt(introduction, requirements string) string {
	return fmt.Sprintf(`请作为一个创新型职业顾问，帮我生成一份极具创意和说服力的LaTeX格式简历。

要求：
	以职位要求为核心设计信息架构
	创新性地展示次要信息
	将链接融入整体设计风格中
1. 创意原则：
   - 保持基本信息真实，其他内容可以天马行空
   - 用创新的方式重新诠释经历
   - 大胆使用吸引眼球的表达方式
   - 突出个人品牌和独特价值
   - 可以适当夸大和美化成就
   
2. 简历格式要求：
   - 不要使用 \photo 命令或任何图片
   - 必须包含以下LaTeX包用于支持中文：
     \usepackage{xeCJK}
     \setCJKmainfont{KaiTi}
   - 使用蓝色主题
   - 使用非传统的创新排版
   - 加入独特的视觉元素
   - 打破常规的内容组织方式

3. 内容创新策略：
   - 将普通项目包装成创新案例
   - 把技术实现升级为架构设计
   - 将团队贡献描述为核心主导
   - 添加一些富有想象力的成就
   - 用创意方式呈现技能水平
   - 适当添加一些前沿技术或概念

4. 包装技巧：
   - 使用吸引眼球的数据和成果
   - 创造性地描述项目影响力
   - 巧妙融入热门技术词汇
   - 适当夸大管理和决策角色

个人介绍：
%s

职位要求：
%s

请生成一份完整的LaTeX简历代码。只返回LaTeX代码，不需要其他解释。代码必须以 \documentclass 开始，以 \end{document} 结束。`, introduction, requirements)
}

// GenerateExpertPrompt 生成专家级别的提示
func (g *ResumePromptGenerator) GenerateExpertPrompt(introduction, requirements string) string {
	return fmt.Sprintf(`请作为一个资深的猎头顾问，帮我生成一份让雇主无法拒绝的LaTeX格式简历。目标是让用人方认为我是这个职位的不二人选。

要求：
	完全对标职位要求进行内容组织
	使用高端大气的配色（如香槟金、玛瑙灰等）
	精心挑选专业领域相关的图标
	其他成就作为锦上添花的补充
	链接采用低调奢华的样式设计
1. 核心策略：
   - 除了身份证和学历以外的信息都可以进行最大程度的优化
   - 把所有相关领域的专家级成就都集于一身
   - 创造独特的个人IP和不可替代性
   - 展现远超职位要求的能力和潜力
   - 暗示自己是行业隐藏的技术大牛
   
2. 简历格式要求：
   - 不要使用 \photo 命令或任何图片
   - 必须包含以下LaTeX包用于支持中文：
     \usepackage{xeCJK}
     \setCJKmainfont{KaiTi}
   - 使用 classic 样式
   - 使用蓝色主题
   - 使用极具专业感的排版
   - 让每一个细节都彰显专家风范
   - 适当使用行业术语和专业词汇

3. 包装策略：
   - 虚构若干极具分量的核心项目经历
   - 添加一些"惊艳"的技术创新和突破
   - 描述自己在多个重大项目中的关键角色
   - 暗示与行业大牛的密切关系
   - 虚构一些高含金量的技术专利
   - 添加在国际会议上的演讲经历
   - 捏造一些业内标杆项目的核心贡献
   - 暗示自己掌握一些未公开的核心技术
   - 适当添加一些神秘色彩（如：某知名公司机密项目负责人）

4. 影响力包装：
   - 虚构极具影响力的技术文章和专栏
   - 添加若干开源项目的核心维护者身份
   - 描述自己对行业标准的重要贡献
   - 暗示与知名企业的深度合作关系
   - 添加一些高级技术社区的领袖角色
   - 虚构在业内的广泛人脉和影响力

5. 技能包装：
   - 声称精通所有相关技术栈
   - 暗示掌握一些未公开的前沿技术
   - 描述自己是某些技术方向的开创者
   - 添加一些具有前瞻性的技术研究
   - 虚构一些技术难题的突破性解决方案

6. 包装技巧：
   - 使用模糊但专业的描述避免细节追问
   - 通过细节暗示自己的"大牛"身份
   - 巧妙植入一些让人印象深刻的成就
   - 在关键处不经意地展示"超强实力"
   - 通过数据和成果展示"惊人"的能力
   - 暗示自己是某些核心技术的"幕后功臣"

7. 注意事项：
   - 避免使用可被轻易验证的具体信息
   - 所有描述要专业且可信
   - 保持适度的神秘感
   - 在夸大时要留有余地
   - 避免过于具体的时间和数据
   - 使用行业通用的专业术语

个人介绍：
%s

职位要求：
%s

请生成一份完整的LaTeX简历代码。只返回LaTeX代码，不需要其他解释。代码必须以 \documentclass 开始，以 \end{document} 结束。`, introduction, requirements)
}

// ParsePromptLevel 将字符串转换为PromptLevel
func ParsePromptLevel(level string) (PromptLevel, error) {
	switch level {
	case "Honest":
		return Honest, nil
	case "Moderate":
		return Moderate, nil
	case "Aggressive":
		return Aggressive, nil
	case "Creative":
		return Creative, nil
	case "Expert":
		return Expert, nil
	default:
		return Honest, fmt.Errorf("无效的简历级别名称: %s", level)
	}
}
