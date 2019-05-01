package ui

import (
	"github.com/civ0/rfc-reader/cache"
	"github.com/civ0/rfc-reader/index"
	"github.com/gdamore/tcell"
	"github.com/rivo/tview"
	"strings"
)

func getPages(rfc string) []string {
	const formFeed = string(12)
	return strings.Split(rfc, formFeed)
}

func GetRFCStatusCell(rfcStatus cache.RFCLocalCopyStatus) *tview.TableCell {
	const absentRune = rune('✕')
	const absentColor = tcell.ColorRed
	const presentRune = rune('✓')
	const presentColor = tcell.ColorGreen
	const unknownRune = rune('?')
	const unknownColor = tcell.ColorPink

	statusRune := absentRune
	statusColor := absentColor
	statusAttr := tcell.AttrBold

	switch rfcStatus {
	case cache.RFCLocalCopyStatusPresent:
		statusRune = presentRune
		statusColor = presentColor
	case cache.RFCLocalCopyStatusUnknown:
		statusRune = unknownRune
		statusColor = unknownColor
		statusAttr = statusAttr | tcell.AttrBlink
	}

	return tview.NewTableCell(string(statusRune)).
		SetTextColor(statusColor).
		SetAttributes(statusAttr)
}

func Run(index *index.RFCIndex) {
	app := tview.NewApplication()
	rfcTable, _ := RFCTable(index)

	statusText := tview.NewTextView()
	statusText.SetText("Status")

	flex := tview.NewFlex().SetDirection(tview.FlexRow)
	flex.SetTitle("RFC Reader")
	flex.SetBorder(true)

	abstractPreview := tview.NewTextView()
	abstractPreview.SetTitle("Abstract")
	abstractPreview.SetBorder(true)

	rfcTable.SetSelectionChangedFunc(func(row, column int) {
		abstractPreview.SetText(index.RFCEntries[row].Abstract)
	})

	rfcText := tview.NewTextView()
	rfcText.SetTitle("RFC")
	rfcText.SetBorder(true)

	rfcTable.SetSelectedFunc(func(row, column int) {
		rfc, err := cache.ReadRFC(index.RFCEntries[row].CanonicalName())
		if err != nil {
			rfcTable.SetCell(row, 1, GetRFCStatusCell(cache.RFCLocalCopyStatusUnknown))
		} else {
			rfcTable.SetCell(row, 1, GetRFCStatusCell(cache.RFCLocalCopyStatusPresent))
		}
		rfcText.SetText(rfc)
	})

	selectionFlex := tview.NewFlex().SetDirection(tview.FlexRow)
	selectionFlex.AddItem(rfcTable, 0, 3, true)
	selectionFlex.AddItem(abstractPreview, 0, 1, false)

	contentFlex := tview.NewFlex()
	contentFlex.AddItem(selectionFlex, 0, 1, true)
	contentFlex.AddItem(rfcText, 0, 3, false)

	flex.AddItem(contentFlex, 0, 1, true)
	flex.AddItem(statusText, 1, 0, false)

	if tviewErr := app.SetRoot(flex, true).Run(); tviewErr != nil {
		panic(tviewErr)
	}
}

func RFCTable(index *index.RFCIndex) (*tview.Table, error) {
	table := tview.NewTable()
	table.SetTitle("RFCs")
	table.SetBorder(true)
	table.SetSelectable(true, false)

	for i, entry := range index.RFCEntries {
		localCopyStatus, err := cache.RFCPresent(entry.CanonicalName())
		if err != nil {
			return nil, err
		}
		table.SetCell(i, 0,
			tview.NewTableCell(entry.DocID))

		table.SetCell(i, 1, GetRFCStatusCell(localCopyStatus))

		table.SetCell(i, 2,
			tview.NewTableCell(entry.Title))
	}

	// TODO: Error handling
	return table, nil
}
