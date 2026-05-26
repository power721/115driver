// Copyright 2025 The 115driver Authors. All rights reserved.
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/power721/115driver/mcp/server"
	"github.com/power721/115driver/pkg/driver"
)

var (
	cookie = flag.String("cookie", "", "115 driver cookie for authentication")
	help   = flag.Bool("help", false, "display help information")
)

func main() {
	flag.Parse()

	if *help {
		fmt.Fprintf(os.Stderr, "Usage: %s [OPTIONS]\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "115 Driver MCP Server - Provides access to 115 cloud storage via MCP protocol\n\n")
		fmt.Fprintf(os.Stderr, "Options:\n")
		flag.PrintDefaults()
		fmt.Fprintf(os.Stderr, "\nExamples:\n")
		fmt.Fprintf(os.Stderr, "  Run server with cookie: %s --cookie=\"UID=xxx;CID=xxx;SEID=xxx\"\n", os.Args[0])
		os.Exit(1)
	}

	if *cookie == "" {
		fmt.Fprintf(os.Stderr, "Error: Cookie is required\n")
		fmt.Fprintf(os.Stderr, "Usage: %s --cookie=\"UID=xxx;CID=xxx;SEID=xxx\"\n", os.Args[0])
		os.Exit(1)
	}

	cr := &driver.Credential{}
	cr.FromCookie(*cookie)
	// Create 115 driver client and authenticate
	client := driver.New(driver.UA(driver.UA115Browser)).ImportCredential(cr)

	// Check login status
	if err := client.CookieCheck(); err != nil {
		log.Fatalf("Authentication failed: %v", err)
	}

	// Create and start the MCP server
	s := server.NewServer().WithClient(client)
	if err := s.Start(context.Background()); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}
