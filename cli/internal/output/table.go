package output

import (
	"fmt"
	"os"
	"strings"
	"text/tabwriter"
)

func FormatFileSize(size int64) string {
	const (
		KB = 1024
		MB = KB * 1024
		GB = MB * 1024
		TB = GB * 1024
	)
	switch {
	case size >= TB:
		return fmt.Sprintf("%.2f TB", float64(size)/float64(TB))
	case size >= GB:
		return fmt.Sprintf("%.2f GB", float64(size)/float64(GB))
	case size >= MB:
		return fmt.Sprintf("%.2f MB", float64(size)/float64(MB))
	case size >= KB:
		return fmt.Sprintf("%.2f KB", float64(size)/float64(KB))
	default:
		return fmt.Sprintf("%d B", size)
	}
}

func (p *Printer) PrintFileTable(path string, files []JSONFile) {
	if p.JSON {
		p.PrintSuccess(map[string]interface{}{
			"path":  path,
			"files": files,
		})
		return
	}

	fmt.Printf("Path: %s\n\n", path)

	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	fmt.Fprintln(w, "NAME\tSIZE\tTYPE\tMODIFIED")
	fmt.Fprintln(w, "----\t----\t----\t--------")

	for _, f := range files {
		name := f.Name
		typ := "file"
		if f.IsDir {
			typ = "dir"
			name += "/"
		}
		fmt.Fprintf(w, "%s\t%s\t%s\t%s\n", name, FormatFileSize(f.Size), typ, f.UpdateTime)
	}
	w.Flush()
}

func (p *Printer) PrintFileList(path string, files []JSONFile) {
	if p.JSON {
		p.PrintSuccess(map[string]interface{}{
			"path":  path,
			"files": files,
		})
		return
	}

	for _, f := range files {
		if f.IsDir {
			fmt.Printf("%s/\n", f.Name)
		} else {
			fmt.Println(f.Name)
		}
	}
}

func (p *Printer) PrintStatTable(stat JSONStat) {
	if p.JSON {
		p.PrintSuccess(stat)
		return
	}

	fmt.Printf("  Name:     %s\n", stat.Name)
	fmt.Printf("  Type:     %s\n", map[bool]string{true: "Directory", false: "File"}[stat.IsDir])
	fmt.Printf("  Size:     %s\n", FormatFileSize(stat.Size))
	if stat.Sha1 != "" {
		fmt.Printf("  SHA1:     %s\n", stat.Sha1)
	}
	if stat.PickCode != "" {
		fmt.Printf("  Pick:     %s\n", stat.PickCode)
	}
	fmt.Printf("  ID:       %s\n", stat.FileID)
	fmt.Printf("  Created:  %s\n", stat.CreateTime)
	fmt.Printf("  Modified: %s\n", stat.UpdateTime)
	if stat.IsDir {
		fmt.Printf("  Files:    %d\n", stat.FileCount)
		fmt.Printf("  Dirs:     %d\n", stat.DirCount)
	}
	if len(stat.Parents) > 0 {
		var parts []string
		for _, pp := range stat.Parents {
			parts = append(parts, pp.Name)
		}
		fmt.Printf("  Path:     /%s\n", strings.Join(parts, "/"))
	}
}

func (p *Printer) PrintOfflineTable(tasks []map[string]interface{}) {
	if p.JSON {
		return
	}

	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	fmt.Fprintln(w, "NAME\tSTATUS\tPROGRESS\tSIZE")
	fmt.Fprintln(w, "----\t------\t--------\t----")

	for _, t := range tasks {
		fmt.Fprintf(w, "%s\t%s\t%.1f%%\t%s\n",
			t["name"], t["status"], t["percent"], FormatFileSize(t["size"].(int64)))
	}
	w.Flush()
}
