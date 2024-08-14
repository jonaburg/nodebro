package page

import (
//	"time"
    "fmt"
    "github.com/gdamore/tcell/v2"
    "github.com/rivo/tview"
)

var sectionDisplayed bool
var currentMode string = "table"

func scrollUpLineByLine(textView *tview.TextView) {
   offsetY, _ := textView.GetScrollOffset()
   textView.ScrollTo(offsetY-1, 0)
}
func scrollDownLineByLine(textView *tview.TextView) {
   offsetY, _ := textView.GetScrollOffset()
   textView.ScrollTo(offsetY+1, 0)
}

func ansiblePane(grid *tview.Grid) {
   header := func(text string) tview.Primitive {
		return tview.NewTextView().
        SetDynamicColors(true).
			SetText("[blue]" + text + "\n")
	}
    ansiblePane := tview.NewTextView().
        SetText("Here is the new selection").
        SetTextAlign(tview.AlignCenter).
        SetDynamicColors(true)
    grid.AddItem(header("Ansible Man"), 0, 0, 1, 3, 0, 0, false)
    grid.AddItem(ansiblePane, 1, 0, 1, 3, 0, 100, false)
    grid.SetColumns(1)
    grid.SetRows(1, 1, 0, 3)
}

func tagDropdown(app *tview.Application, grid *tview.Grid) {
   dropdown := tview.NewDropDown().
   		SetLabel("Select an option (hit Enter): ").
   		SetOptions([]string{"8.0.1", "8.1.0", "v8.3.2", "polkadot-v.90", "Fifth"}, nil).
        SetFieldBackgroundColor(tcell.ColorDefault).
        SetPrefixTextColor(tcell.ColorYellow).
       SetFieldWidth(30).
       SetFinishedFunc(func (tcell.Key) {
           fmt.Printf("thats it")
       })


       //       	SetFormAttributes(labelWidth int, labelColor, bgColor, fieldTextColor, fieldBgColor tcell.Color) FormItem
       
       grid.AddItem(dropdown, 1, 0, 1, 3, 0, 10, false)
       grid.SetColumns(2)
       grid.SetRows(1, 1, 0, 3)
       app.SetRoot(grid, true)
       app.SetFocus(dropdown)
}

func ReleasePage(owner, repo, tag, body, createdat, timediff string, app *tview.Application) tview.Primitive {

    createdAtStr := createdat


    textView := tview.NewTextView().
         SetText("[red]Owner: " + owner + "\n" +
            "[white]Repo: " + repo + "\n" +
            "[yellow]Tag:[white] " + tag + "\n" +
            "[yellow]Created At:[white] " + createdAtStr + "\n" + 
            body).
        SetDynamicColors(true).
        SetScrollable(true) // Enable scrolling

   header := func(text string) tview.Primitive {
		return tview.NewTextView().
        SetDynamicColors(true).
			SetText("[blue]" + text + "\n")
	}

    grid := tview.NewGrid().
        SetColumns(1).
        SetRows(1, 0, 2).
        SetBorders(true).
        AddItem(header("Release Notes"), 0, 0, 1, 3, 0, 0, false).
        AddItem(textView, 1, 0, 1, 3, 0, 10, false).
        AddItem(header("(b) to go back. (t) to view releases. (c) to compile"), 2, 0, 1, 3, 0, 0, false)


    grid.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
        textView.InputHandler()(event, nil)
        switch event.Rune() {
        case 'k':
//            textView.ScrollToBeginning()
            scrollUpLineByLine(textView)
            return nil
        case 'j':
//           textView.ScrollToEnd()
            scrollDownLineByLine(textView)
            return nil
        case 'G':
            textView.ScrollToEnd()
		    return nil
        case 'l':
		    return nil
        case 'h':
		    return nil
        case 'c':
		    ansiblePane(grid)
        case 't':
		    tagDropdown(app, grid)
        }
	
        return event
    })


    return grid
}

