package main

import (
	"fmt"
	"log"
	"os"
	"sort"
	"text/tabwriter"
	"time"

	transmission "github.com/metalmatze/transmission-exporter"
)

func main() {
	client := transmission.New("http://localhost:9091", nil)

	for {
		torrents, err := client.GetTorrents()
		if err != nil {
			log.Fatal(err)
		}

		sort.Sort(sort.Reverse(transmission.ByRatio(torrents)))
		//sort.Sort(sort.Reverse(transmission.ByDate(torrents)))

		print("\033[H\033[2J")
		w := tabwriter.NewWriter(os.Stdout, 0, 0, 1, ' ', 0)

		fmt.Fprintln(w, "Name\tRatio\tUp\tDown\tPercent\tDate")
		for _, t := range torrents {
			fmt.Fprintf(w,
				"%s\t%f\t%d\t%d\t%f\t%v\n",
				t.Name,
				t.UploadRatio,
				t.RateUpload,
				t.RateDownload,
				t.PercentDone,
				time.Unix(int64(t.Date), 0),
			)
		}
		w.Flush()

		time.Sleep(time.Second)
	}
}
