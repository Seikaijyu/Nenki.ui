package text

import gtext "github.com/Seikaijyu/gio/text"

// WrapPolicy 就像是一本故事书，告诉我们如何正确地把一个长长的故事分成一页一页的小故事。
type WrapPolicy = gtext.WrapPolicy

const (
	// WrapHeuristically 就像是我们在读故事的时候，尽量让每个单词都完整的出现在一页上。
	// 只有当一个单词太长，无法完整的出现在一页上时，我们才会把它分开。
	// 而且，如果一页的最后一个单词被截断了，我们会尽量让更多的字母保留在这一页。
	WrapHeuristically WrapPolicy = iota
	// WrapWords 就是我们坚决不允许一个单词被分开出现在两页上。就算这样做会让某个单词超出页面。
	WrapWords
	// WrapGraphemes 就是我们为了让每一页都尽量多的包含文字，不惜牺牲阅读的流畅性，
	// 在任何地方都可能把一个单词分开。
	WrapGraphemes
)
package text

import gtext "gioui.org/text"

// WrapPolicy 就像是一本故事书，告诉我们如何正确地把一个长长的故事分成一页一页的小故事。
type WrapPolicy = gtext.WrapPolicy

const (
	// WrapHeuristically 就像是我们在读故事的时候，尽量让每个单词都完整的出现在一页上。
	// 只有当一个单词太长，无法完整的出现在一页上时，我们才会把它分开。
	// 而且，如果一页的最后一个单词被截断了，我们会尽量让更多的字母保留在这一页。
	WrapHeuristically WrapPolicy = iota
	// WrapWords 就是我们坚决不允许一个单词被分开出现在两页上。就算这样做会让某个单词超出页面。
	WrapWords
	// WrapGraphemes 就是我们为了让每一页都尽量多的包含文字，不惜牺牲阅读的流畅性，
	// 在任何地方都可能把一个单词分开。
	WrapGraphemes
)
