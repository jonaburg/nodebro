package main

import (
    "regexp"
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"sort"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"github.com/BurntSushi/toml"
    "nodebro/page"
    "nodebro/requests"
)

type Config struct {
    Pat   string     `toml:"pat"`
    Repos []RepoInfo `toml:"repos"`
}

type RepoInfo struct {
	Owner string `toml:"Owner"`
	Repo  string `toml:"Repo"`
}

const configFileRelative = ".config/nodebro/config"

func addRow(table *tview.Table, row int, owner, repo,  createdat, timediff, tag string, colorTag bool) {
    tagCell      := tview.NewTableCell(tag).SetAlign(tview.AlignCenter).SetBackgroundColor(tcell.ColorDefault)
    timeDiffCell := tview.NewTableCell(timediff).SetAlign(tview.AlignCenter).SetBackgroundColor(tcell.ColorDefault)
    table.SetCell(row, 0, tview.NewTableCell(owner).SetAlign(tview.AlignCenter).SetBackgroundColor(tcell.ColorDefault))
    table.SetCell(row, 1, tview.NewTableCell(repo).SetAlign(tview.AlignCenter).SetBackgroundColor(tcell.ColorDefault))
    table.SetCell(row, 2, tview.NewTableCell(createdat).SetAlign(tview.AlignCenter).SetBackgroundColor(tcell.ColorDefault))
    table.SetCell(row, 3, timeDiffCell)
    table.SetCell(row, 4, tagCell)
    tview.Styles.PrimitiveBackgroundColor = tcell.ColorDefault
    tview.Styles.ContrastBackgroundColor = tcell.ColorDefault
    table.SetSelectedStyle(tcell.Style{}.
        Background(tcell.ColorDefault).
        Foreground(tcell.ColorBlue). 
        Attributes(tcell.AttrNone))
}


func main() {
    // Load config
    homeDir, err := os.UserHomeDir()
    if err != nil {
		fmt.Printf("Error getting user home directory: %v", err)
	}
    	configPath := filepath.Join(homeDir, configFileRelative)
    	var config Config
    	if _, err := toml.DecodeFile(configPath, &config); err != nil {
    		fmt.Printf("Error decoding TOML file: %v", err)
    	}
        for _, repo := range config.Repos {
            fmt.Printf("Loading Config: Owner: %s, Repo: %s\n", repo.Owner, repo.Repo)
}
    authToken := config.Pat
    var allReleases []requests.Release
	var wg sync.WaitGroup
    var selectedRow int
	app := tview.NewApplication()
    patternDay := regexp.MustCompile(`\d+d\b`)
    patternHour := regexp.MustCompile(`\d+h\b`)
	table := tview.NewTable().
		SetBorders(true).
		SetBordersColor(tview.Styles.BorderColor).
		SetSeparator('|').
        SetSelectable(true,false)


        table.SetSelectedFunc(func(row, column int) {
         if row == 0  {
            return
        }
        selectedRow = row
        owner := table.GetCell(selectedRow, 0).Text
        repo := table.GetCell(selectedRow, 1).Text
        tag := table.GetCell(selectedRow, 4).Text
        var (
            releaseBody string
            releaseCreatedAt string
            releaseHourDiff string
            )
        // Find the selected release from allReleases
        for _, release := range allReleases {
            if release.Owner == owner && release.Repo == repo {
                releaseBody = release.Body
                releaseCreatedAt = release.CreatedAt.String()
                releaseHourDiff = release.TimeDifference
                break
          }
        }
        releasePage := page.ReleasePage(owner,repo,tag,releaseBody, releaseCreatedAt, releaseHourDiff, app)
        app.SetFocus(releasePage)
        app.SetRoot(releasePage, true)
    })

	addRow(table, 0, "Owner", "Repo", "Created At", "Time Ago", "Tag", false)
	//  placeholders for tag values while fetching
	for i := 1; i <= len(config.Repos); i++ {
		addRow(table, i, config.Repos[i-1].Owner, config.Repos[i-1].Repo, "Fetching...", "Calculating...", "Fetching...",  false)
	}

   header := func(text string) tview.Primitive {
		return tview.NewTextView().
        SetDynamicColors(true).
			SetText("[green]" + text + "\n")
	}

	MainGrid := tview.NewGrid().
        SetColumns(0,1).
        SetRows(1,0).
        SetBorders(true).
        AddItem(header("Github Releases Tracker  - Latest Tags"), 0, 0, 1, 2, 0, 0, false).
		AddItem(table, 1, 0, 1, 2, 0, 10, true)
        table.SetBackgroundColor(tcell.ColorDefault)

//	MainGrid := tview.NewGrid().
//        SetColumns(2,0,1).
//        SetRows(1,0,1).
//        SetBorders(true).
//        AddItem(header("Github Releases  - Latest Tags"), 0, 1, 1, 1, 0, 0, false).
//        AddItem(header("Github Releases  - Latest Tags"), 0, 1, 1, 1, 0, 0, false).
//		AddItem(table, 1, 1, 1, 1, 0, 0, true)
//        table.SetBackgroundColor(tcell.ColorDefault)

// Start a goroutine to fetch information for each repository
tagCh := make(chan requests.Release, len(config.Repos))
// Add a separate wait group for the goroutines reading from tagCh
var readWG sync.WaitGroup
readWG.Add(len(config.Repos))

for _, repo := range config.Repos {
    wg.Add(1)
    go func(repoInfo RepoInfo) {
        defer wg.Done()
        requests.FetchLatestRelease(requests.RepoInfo(repoInfo), tagCh, &readWG, authToken)
    }(repo)
}
go func() {
    // Wait for all writing goroutines to finish
    wg.Wait()
    close(tagCh) // Close the channel after all releases are received
}()
// Start a goroutine to read from tagCh and append to allReleases
go func() {
    // Wait for all reading goroutines to finish
    readWG.Wait()

    // Read from tagCh until it's closed
    for release := range tagCh {
        allReleases = append(allReleases, release)
    }

    // Update the table with fetched tag information
    app.QueueUpdateDraw(func() {
    // sort the releases based on creation time
    sort.Slice(allReleases, func(i, j int) bool {
        return allReleases[i].CreatedAt.After(allReleases[j].CreatedAt)
    })
        for i, release := range allReleases {
            addRow(table, i+1, release.Owner, release.Repo, release.CreatedAt.String(), release.TimeDifference, release.TagName, true)

            timeDiffCell := table.GetCell(i+1, 3)
            if timeDiffCell != nil {
                if patternHour.MatchString(release.TimeDifference) {
                    timeDiffCell.SetText(fmt.Sprintf("[green] %s", release.TimeDifference)).SetBackgroundColor(tcell.ColorDefault)
                }
                if patternDay.MatchString(release.TimeDifference) {
                    timeDiffCell.SetText(fmt.Sprintf("[yellow] %s", release.TimeDifference)).SetBackgroundColor(tcell.ColorDefault)
                }
            }
        }
        app.SetFocus(table)
    })
}()



app.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
    // Handle navigation based on the current mode
    var currentMode string = "table"
    switch currentMode {
    case "table":
        switch event.Key() {
        case tcell.KeyUp, 'k':
            // Move selection up
            selectedRow, _ := table.GetSelection()
            if selectedRow > 0 {
                table.Select(selectedRow-1, 0)
            }
            return nil
        case tcell.KeyDown, 'j':
            selectedRow, _ := table.GetSelection()
            if selectedRow < table.GetRowCount()-1 {
                table.Select(selectedRow+1, 0)
            }
            return nil
        }
    }
    if event.Rune() == 'x' {
        app.Stop()
        return nil
    }
    if event.Key() == tcell.KeyRune && event.Rune() == 'b' {
        app.SetRoot(MainGrid, true)
        app.SetFocus(table) 
        return nil
    }
    return event
})



	if err := app.SetRoot(MainGrid, true).Run(); err != nil {
		panic(err)
	}
}

