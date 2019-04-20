package ui

import (
	"github.com/civ0/rfc-reader/cache"
	"github.com/civ0/rfc-reader/index"
	"github.com/rivo/tview"
	"strings"
)

func getPages(rfc string) []string {
	const formFeed = string(12)
	return strings.Split(rfc, formFeed)
}

func Run(index index.RFCIndex) {
	app := tview.NewApplication()
	rfcTable, _ := RFCTable(index)

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
		rfc, _ := cache.ReadRFC(index.RFCEntries[row].CanonicalName())
		rfcText.SetText(rfc)
	})

	selectionFlex := tview.NewFlex().SetDirection(tview.FlexRow)
	selectionFlex.AddItem(rfcTable, 0, 3, true)
	selectionFlex.AddItem(abstractPreview, 0, 1, false)

	flex := tview.NewFlex()
	flex.SetTitle("RFC Reader")
	flex.SetBorder(true)
	flex.AddItem(selectionFlex, 0, 1, true)
	flex.AddItem(rfcText, 0, 3, false)

	if tviewErr := app.SetRoot(flex, true).Run(); tviewErr != nil {
		panic(tviewErr)
	}
}

func RFCTable(index index.RFCIndex) (*tview.Table, error) {
	table := tview.NewTable()
	table.SetTitle("RFCs")
	table.SetBorder(true)
	table.SetSelectable(true, false)

	for i, entry := range index.RFCEntries {
		table.SetCell(i, 0,
			tview.NewTableCell(entry.DocID))
		table.SetCell(i, 1,
			tview.NewTableCell(entry.Title))
	}

	// TODO: Error handling
	return table, nil
}
