/*
Copyright (c) 2015 VMware, Inc. All Rights Reserved.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

/*
This example program shows how the `finder` and `property` packages can
be used to navigate a vSphere inventory structure using govmomi.
*/

package main

import (
	"fmt"
	"net/url"
	"os"
	"text/tabwriter"

	"github.com/vmware/govmomi"
	"github.com/vmware/govmomi/find"
	"github.com/vmware/govmomi/property"
	"github.com/vmware/govmomi/units"
	"github.com/vmware/govmomi/vim25/mo"
	"github.com/vmware/govmomi/vim25/types"
	"golang.org/x/crypto/ssh/terminal"
	"golang.org/x/net/context"
)

func exit(err error) {
	fmt.Fprintf(os.Stderr, "Error: %s\n", err)
	os.Exit(1)
}

func askpasswords(args *Arguments) error {

	if args.AskPass {
		fmt.Print(">>> Password: ")
		password_b, err := terminal.ReadPassword(0)
		fmt.Println("")
		if err != nil {
			return err
		}

		password := string(password_b)
		args.Pass = password
	}

	return nil

}

func main() {

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	args := &Arguments{}
	args.Parse()
	err := args.Validate()
	if err != nil {
		exit(err)
	}

	askpasswords(args)

	//Form URL u here with entered flags.
	connection_url := fmt.Sprintf("https://%s:%s@%s/sdk", args.User, args.Pass, args.Host)

	//fmt.Println(connection_url)
	//os.Exit(1)

	// Parse URL from string
	u, err := url.Parse(connection_url)
	if err != nil {
		exit(err)
	}

	// Connect and log in to ESX or vCenter
	c, err := govmomi.NewClient(ctx, u, args.Insecure)
	if err != nil {
		exit(err)
	}

	f := find.NewFinder(c.Client, true)

	// Find one and only datacenter
	dc, err := f.DefaultDatacenter(ctx)
	if err != nil {
		exit(err)
	}

	// Make future calls local to this datacenter
	f.SetDatacenter(dc)

	// Find datastores in datacenter
	dss, err := f.DatastoreList(ctx, "*")
	if err != nil {
		exit(err)
	}

	pc := property.DefaultCollector(c.Client)

	// Convert datastores into list of references
	var refs []types.ManagedObjectReference
	for _, ds := range dss {
		refs = append(refs, ds.Reference())
	}

	// Retrieve summary property for all datastores
	var dst []mo.Datastore
	err = pc.Retrieve(ctx, refs, []string{"summary", "vm"}, &dst)
	if err != nil {
		exit(err)
	}

	// Print summary per datastore
	tw := tabwriter.NewWriter(os.Stdout, 2, 0, 2, ' ', 0)
	fmt.Fprintf(tw, "Name:\tType:\tCapacity:\tFree:\n")
	for _, ds := range dst {
		fmt.Fprintf(tw, "%s\t", ds.Summary.Name)
		fmt.Fprintf(tw, "%s\t", ds.Summary.Type)
		fmt.Fprintf(tw, "%s\t", units.ByteSize(ds.Summary.Capacity))
		fmt.Fprintf(tw, "%s\t", units.ByteSize(ds.Summary.FreeSpace))
		fmt.Fprintf(tw, "\n")

		if ds.Vm == nil {
			continue
		}

		var vms []mo.VirtualMachine
		err := pc.Retrieve(ctx, ds.Vm, []string{"name", "summary"}, &vms)
		if err != nil {
			fmt.Println(err)
			continue
		}
		fmt.Println("Virtual machines found:", len(vms))
		for _, vm := range vms {
			fmt.Fprintf(tw, "\t\t%s %s %s\n", vm.Name, vm.Summary.Guest.IpAddress, vm.Summary.Guest.HostName)
		}
		fmt.Fprintf(tw, "\n")
	}
	tw.Flush()

}
