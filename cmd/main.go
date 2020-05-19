package main

import (
	"flag"
	"fmt"
	"github.com/kpfaulkner/azureps/pkg"
)

func main() {
	fmt.Printf("so it begins....\n")

	subscriptionID := flag.String("subID", "", "subscription ID")
	tenantID := flag.String("tenID", "", "tenant ID")
	clientID := flag.String("clientID", "", "client ID")
	clientSecret := flag.String("clientSecret", "", "client Secret")

	flag.Parse()

  ps := pkg.NewAzurePowerShellRemote(*subscriptionID, *tenantID, *clientID, *clientSecret)
  ps.RunPowerShell("myrg", "myvm", "echo \"testing\" > c:/temp/ken2")

}
