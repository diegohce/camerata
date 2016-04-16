package vmware

import (
	"errors"
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

func (be *cloudVmWare) uses_bastion(ip string) bool {

	args := be.args

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

func (be *cloudVmWare) askpasswords() error {

	args := be.args

	if args.Password == "" {
		fmt.Fprint(os.Stderr, "Password: ")
		password_b, err := terminal.ReadPassword(0)
		fmt.Fprintln(os.Stderr, "")
		if err != nil {
			return err
		}

		password := string(password_b)
		args.Password = password
	}

	return nil

}

func (be *cloudVmWare) vmwareInventory() error {

	args := be.args

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	be.askpasswords()

	//Form URL u here with entered flags.
	connection_url := fmt.Sprintf("https://%s:%s@%s/sdk",
		args.User, args.Password, args.Host)

	// Parse URL from string
	u, err := url.Parse(connection_url)
	if err != nil {
		return err
	}

	// Connect and log in to ESX or vCenter
	c, err := govmomi.NewClient(ctx, u, args.Insecure)
	if err != nil {
		return err
	}

	f := find.NewFinder(c.Client, true)

	// Find one and only datacenter
	dc, err := f.DefaultDatacenter(ctx)
	if err != nil {
		return err
	}

	// Make future calls local to this datacenter
	f.SetDatacenter(dc)

	// Find datastores in datacenter
	dss, err := f.DatastoreList(ctx, "*")
	if err != nil {
		return err
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
		return err
	}

	switch args.OutputFormat {
	case "toml":
		return be.tomlOutput(ctx, dst, pc)

	case "csv":
		return be.csvOutput(ctx, dst, pc)

	default:
		be.tomlOutput(ctx, dst, pc)
		return errors.New("Invalid output format: " + args.OutputFormat)

	}

	return nil

	//	if args.Bastion != "" {
	//		fmt.Println("[bastion]")
	//		fmt.Printf("host=\"%s\"\n", args.Bastion)
	//		fmt.Println("user=\"MUST_SET_USER\"")
	//		fmt.Println("password=\"\"")
	//		fmt.Println("")
	//	}

	//	fmt.Println("[servers]")

	//	// Print VMs per datastore
	//	for _, ds := range dst {

	//		if ds.Vm == nil {
	//			continue
	//		}

	//		var vms []mo.VirtualMachine
	//		err := pc.Retrieve(ctx, ds.Vm, []string{"name", "summary"}, &vms)
	//		if err != nil {
	//			fmt.Println(err)
	//			continue
	//		}
	//		//		fmt.Println("Virtual machines found:", len(vms))
	//		for _, vm := range vms {
	//			if vm.Summary.Guest.IpAddress == "" {
	//				fmt.Printf("#\t[servers.%s]\n", strings.Replace(vm.Name, " ", "_", -1))
	//				//			fmt.Printf("\t\t%s %s %s\n", vm.Name, vm.Summary.Guest.IpAddress, vm.Summary.Guest.HostName)
	//				fmt.Printf("#\thost=\"%s\"\n\n", vm.Summary.Guest.IpAddress)
	//			} else {
	//				fmt.Printf("\t[servers.%s]\n", strings.Replace(vm.Name, " ", "_", -1))
	//				fmt.Printf("\thost=\"%s\"\n", vm.Summary.Guest.IpAddress)
	//				if args.Bastion != "" && be.uses_bastion(vm.Summary.Guest.IpAddress) {
	//					fmt.Println("\tuse_bastion=true")
	//				}
	//				fmt.Println("")
	//			}
	//		}
	//	}
	//
	//	return nil
}

func (be *cloudVmWare) tomlOutput(ctx context.Context, dst []mo.Datastore, pc *property.Collector) error {

	args := be.args

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
		for _, vm := range vms {
			if vm.Summary.Guest.IpAddress == "" {
				fmt.Printf("#\t[servers.%s]\n", strings.Replace(vm.Name, " ", "_", -1))
				fmt.Printf("#\thost=\"%s\"\n\n", vm.Summary.Guest.IpAddress)
			} else {
				fmt.Printf("\t[servers.%s]\n", strings.Replace(vm.Name, " ", "_", -1))
				fmt.Printf("\thost=\"%s\"\n", vm.Summary.Guest.IpAddress)
				if args.Bastion != "" && be.uses_bastion(vm.Summary.Guest.IpAddress) {
					fmt.Println("\tuse_bastion=true")
				}
			}
			fmt.Println("")
		}
	}
	return nil
}

func (be *cloudVmWare) csvOutput(ctx context.Context, dst []mo.Datastore, pc *property.Collector) error {

	args := be.args

	fmt.Println("Host, ip address, uses bastion")

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
		for _, vm := range vms {
			fmt.Print(vm.Name, ",", vm.Summary.Guest.IpAddress)
			if args.Bastion != "" && be.uses_bastion(vm.Summary.Guest.IpAddress) {
				fmt.Print(", true")
			}
			fmt.Println("")
		}
	}
	return nil
}
