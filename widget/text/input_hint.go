package text

import gkey "gioui.org/io/key"

// InputHint 定义了屏幕键盘的类型，也就是给出了用户可能要输入的数据类型的提示。
type InputHint = gkey.InputHint

const (
	// HintAny 表示期待任何类型的输入。
	HintAny InputHint = iota
	// HintText 表示期待文本输入。它可能会激活自动修正和建议。
	HintText
	// HintNumeric 表示期待数字输入。它可能会激活0-9、"."和","的快捷方式。
	HintNumeric
	// HintEmail 表示期待电子邮件输入。它可能会激活常用的电子邮件字符的快捷方式，如"@"和".com"。
	HintEmail
	// HintURL 表示期待URL输入。它可能会激活常见的URL片段的快捷方式，如"/"和".com"。
	HintURL
	// HintTelephone 表示期待电话号码输入。它可能会激活0-9、"#"和"*"的快捷方式。
	HintTelephone
	// HintPassword 表示期待密码输入。它可能会禁用自动修正并启用密码自动填充。
	HintPassword
)
