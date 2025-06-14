package app

import (
	"fmt"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/lipgloss"

	"github.com/luca2040/clipse/config"
)

var (
	style                 = lipgloss.NewStyle()
	titleStyle, descStyle string
	appStyle              = style.Padding(1, 2)
	statusMessageStyle    = style.Foreground(
		lipgloss.AdaptiveColor{
			Light: defaultMsgColor,
			Dark:  defaultMsgColor,
		},
	).Render

	previewTitleStyle = func() lipgloss.Style {
		b := lipgloss.RoundedBorder()
		b.Right = borderRightChar
		return style.BorderStyle(b) //.MarginTop(1)
	}()

	previewInfoStyle = func() lipgloss.Style {
		b := lipgloss.RoundedBorder()
		b.Left = borderLeftChar
		return style.BorderStyle(b) //MarginBottom(1)
	}()
)

func (d itemDelegate) itemFilterStyle(i item) string {
	titleStyle := style.
		Foreground(lipgloss.Color(d.theme.DimmedTitle)).
		PaddingLeft(2).
		Render(i.titleBase)

	descStyle := style.
		Foreground(lipgloss.Color(d.theme.DimmedDesc)).
		PaddingLeft(2).
		Render(i.descriptionBase)

	return fmt.Sprintf("%s\n%s", titleStyle, descStyle)
}

func (d itemDelegate) itemChosenStyle(i item) string {
	titleStyle = style.
		Foreground(lipgloss.Color(d.theme.SelectedTitle)).
		PaddingLeft(1).
		BorderLeft(true).BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color(d.theme.SelectedDescBorder)).
		Render(i.titleBase)

	descStyle = style.
		Foreground(lipgloss.Color(d.theme.SelectedDesc)).
		PaddingLeft(1).
		BorderLeft(true).BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color(d.theme.SelectedDescBorder)).
		Render(i.descriptionBase)

	if i.pinned {
		descStyle += styledPin(d.theme)
	}

	return fmt.Sprintf("%s\n%s", titleStyle, descStyle)
}

func (d itemDelegate) itemSelectedStyle(i item) string {

	titleStyle = style.
		Foreground(lipgloss.Color(d.theme.SelectedTitle)).
		PaddingLeft(2).
		Render(i.titleBase)

	descStyle = style.
		Foreground(lipgloss.Color(d.theme.SelectedDesc)).
		PaddingLeft(2).
		Render(i.descriptionBase)

	if i.pinned {
		descStyle += styledPin(d.theme)
	}

	return fmt.Sprintf("%s\n%s", titleStyle, descStyle)
}

func (d itemDelegate) itemNormalStyle(i item) string {
	titleStyle = style.
		Foreground(lipgloss.Color(d.theme.NormalTitle)).
		PaddingLeft(2).
		Render(i.titleBase)

	descStyle = style.
		Foreground(lipgloss.Color(d.theme.NormalDesc)).
		PaddingLeft(2).
		Render(i.descriptionBase)

	if i.pinned {
		descStyle += styledPin(d.theme)
	}

	return fmt.Sprintf("%s\n%s", titleStyle, descStyle)
}

func styledList(clipboardList list.Model, ct config.CustomTheme) list.Model {
	clipboardList.FilterInput.PromptStyle = style.
		Foreground(lipgloss.Color(ct.FilterPrompt)).
		PaddingTop(1)
	clipboardList.FilterInput.TextStyle = style.Foreground(lipgloss.Color(ct.FilterText))
	clipboardList.Styles.StatusBarFilterCount = style.Foreground(lipgloss.Color(ct.FilterInfo))
	clipboardList.FilterInput.Cursor.Style = style.Foreground(lipgloss.Color(ct.FilterCursor))
	clipboardList.Styles.StatusEmpty = style.Foreground(lipgloss.Color(ct.FilterInfo))
	clipboardList.Help.Styles.ShortKey = style.Foreground(lipgloss.Color(ct.HelpKey))
	clipboardList.Help.Styles.ShortDesc = style.Foreground(lipgloss.Color(ct.HelpDesc))
	clipboardList.Help.Styles.FullKey = style.Foreground(lipgloss.Color(ct.HelpKey))
	clipboardList.Help.Styles.FullDesc = style.Foreground(lipgloss.Color(ct.HelpDesc))
	clipboardList.Paginator.ActiveDot = style.
		Foreground(lipgloss.Color(ct.PageActiveDot)).
		Render("•")
	clipboardList.Paginator.InactiveDot = style.
		Foreground(lipgloss.Color(ct.PageInactiveDot)).
		Render("•")
	clipboardList.Styles.StatusBar = style.
		Foreground(lipgloss.Color(ct.TitleInfo)).
		MarginBottom(1).
		MarginLeft(2)
	clipboardList.Styles.Title = style.
		Foreground(lipgloss.Color(ct.TitleFore)).
		Background(lipgloss.Color(ct.TitleBack)).
		MarginTop(1).
		Align(lipgloss.Position(1))
	clipboardList.Styles.DividerDot = style.
		Foreground(lipgloss.Color(ct.DividerDot)).
		SetString("•").
		PaddingLeft(1).
		PaddingRight(1)
	clipboardList.Help.FullSeparator = style.
		Foreground(lipgloss.Color(ct.DividerDot)).
		PaddingLeft(1).
		PaddingRight(1).
		Render("•")
	clipboardList.Help.ShortSeparator = style.
		Foreground(lipgloss.Color(ct.DividerDot)).
		PaddingLeft(1).
		PaddingRight(1).
		Render("•")
	clipboardList.Styles.NoItems = style.
		Foreground(lipgloss.Color(ct.TitleInfo)).
		PaddingBottom(1).
		PaddingLeft(2)

	return clipboardList
}

func styledHelp(help help.Model, ct config.CustomTheme) help.Model {
	help.Styles.ShortKey = style.Foreground(lipgloss.Color(ct.HelpKey))
	help.Styles.ShortDesc = style.Foreground(lipgloss.Color(ct.HelpDesc))
	help.Styles.FullKey = style.Foreground(lipgloss.Color(ct.HelpKey))
	help.Styles.FullDesc = style.Foreground(lipgloss.Color(ct.HelpDesc))
	help.FullSeparator = style.Foreground(lipgloss.Color(ct.DividerDot)).
		PaddingLeft(1).
		PaddingRight(1).
		Render("•")
	help.ShortSeparator = style.
		Foreground(lipgloss.Color(ct.DividerDot)).
		PaddingLeft(1).
		PaddingRight(1).
		Render("•")
	return help
}

func styledStatusMessage(ct config.CustomTheme) func(strs ...string) string {
	return style.
		Foreground(lipgloss.AdaptiveColor{Light: ct.StatusMsg, Dark: ct.StatusMsg}).
		Render
}

func styledPin(theme config.CustomTheme) string {
	return style.
		Foreground(lipgloss.Color(theme.PinIndicatorColor)).
		Render(pinChar)
}

func (m *Model) styledPreviewHeader(str string) string {
	return style.
		Foreground(lipgloss.Color(m.theme.PreviewBorder)).
		MarginTop(2).
		Render(str)
}

func (m *Model) styledPreviewFooter(str string) string {
	return style.
		Foreground(lipgloss.Color(m.theme.PreviewBorder)).
		Render(str)
}

func (m *Model) styledPreviewContent(content string) string {
	return style.
		Foreground(lipgloss.Color(m.theme.PreviewedText)).
		Render(content)
}
