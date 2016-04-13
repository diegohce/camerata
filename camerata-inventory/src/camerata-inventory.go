package main

import (
	"fmt"
	"net/url"
	"os"
	"strings"

	"github.com/vmware/govmomi"
	"github.com/vmware/govmomi/find"
	"github.com/vmware/govmomi/property"
	//	"github.com/vmware/govmomi/units"
	"github.com/vmware/govmomi/vim25/mo"
	"github.com/vmware/govmomi/vim25/types"
	"golang.org/x/crypto/ssh/terminal"
	"golang.org/x/net/context"
)

func exit(err error) {
	fmt.Fprintf(os.Stderr, "Error: %s\n", err)
	os.Exit(1)
}

func uses_bastion(args *Arguments, ip string) bool {

	if args.Bastion != "" {
		if args.BastionNets != "" {
			for _, prefix := range strings.Split(args.BastionNets, ",") {
				if strings.HasPrefix(ip, strings.TrimSpace(prefix)) {
					return true
				}
			}
		}
	}
	return false
}

func askpasswords(args *Arguments) error {

	if args.AskPass {
		fmt.Fprint(os.Stderr, "Password: ")
		password_b, err := terminal.ReadPassword(0)
		fmt.Fprintln(os.Stderr, "")
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

	// Retrieve summary && VMs properties for all datastores
	var dst []mo.Datastore
	err = pc.Retrieve(ctx, refs, []string{"summary", "vm"}, &dst)
	if err != nil {
		exit(err)
	}

	if args.Bastion != "" {
		fmt.Println("[bastion]")
		fmt.Printf("host=\"%s\"\n", args.Bastion)
		fmt.Println("user=\"MUST_SET_USER\"")
		fmt.Println("password=\"\"")
		fmt.Println("")
	}

	fmt.Println("[servers]")

	// Print VMs per datastore
	for _, ds := range dst {

		if ds.Vm == nil {
			continue
		}

		var vms []mo.VirtualMachine
		err := pc.Retrieve(ctx, ds.Vm, []string{"name", "summary"}, &vms)
		if err != nil {
			fmt.Println(err)
			continue
		}
		//		fmt.Println("Virtual machines found:", len(vms))
		for _, vm := range vms {
			if vm.Summary.Guest.IpAddress == "" {
				fmt.Printf("#\t[servers.%s]\n", strings.Replace(vm.Name, " ", "_", -1))
				//			fmt.Printf("\t\t%s %s %s\n", vm.Name, vm.Summary.Guest.IpAddress, vm.Summary.Guest.HostName)
				fmt.Printf("#\thost=\"%s\"\n\n", vm.Summary.Guest.IpAddress)
			} else {
				fmt.Printf("\t[servers.%s]\n", strings.Replace(vm.Name, " ", "_", -1))
				fmt.Printf("\thost=\"%s\"\n", vm.Summary.Guest.IpAddress)
				if args.Bastion != "" && uses_bastion(args, vm.Summary.Guest.IpAddress) {
					fmt.Println("\tuse_bastion=true")
				}
				fmt.Println("")
			}
		}
	}

}