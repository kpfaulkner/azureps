package pkg

import (
	"fmt"
	"github.com/kpfaulkner/azureauth"
	"io/ioutil"
	"net/http"
	"strings"
)

type AzurePowerShellRemote struct {
	azureAuth *azureauth.AzureAuth
	subscriptionID string
	tenantID string
	clientID string
}

func NewAzurePowerShellRemote( subscriptionID string, tenantID string, clientID string, clientSecret string ) AzurePowerShellRemote {

	apsr := AzurePowerShellRemote{}
	apsr.azureAuth = azureauth.NewAzureAuth(subscriptionID, tenantID, clientID, clientSecret)
	apsr.subscriptionID = subscriptionID
	apsr.clientID = clientID
	apsr.tenantID = tenantID
	return apsr
}

// RunPowerShell run powershell script. This absolutely needs some script validation for malicious content
// but currently is only used for running pre-defined scripts on target VM.
func (ps *AzurePowerShellRemote) RunPowerShell(resourceGroup string, vmName string, script string) error {

	// always start with a refreshed token.
	err := ps.azureAuth.RefreshToken()
	if err != nil {
		fmt.Printf("Unable to refresh token : %s\n", err.Error())
		return err
	}

	template := "https://management.azure.com/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Compute/virtualMachines/%s/runCommand?api-version=2018-04-01"
	url := fmt.Sprintf(template, ps.subscriptionID, resourceGroup, vmName)

	bodyTemplate := "{\"commandId\":\"RunPowerShellScript\",\"script\":['%s']}"

	// need a way to validate this script... or at least look for malicious stuff! TODO(kpfaulkner)
	body := fmt.Sprintf(bodyTemplate, script)
	req, err := http.NewRequest("POST", url, strings.NewReader(body))
	req.Header.Add("Authorization", "Bearer " + ps.azureAuth.CurrentToken().AccessToken)
	req.Header.Add("Content-type", "application/json")

	client := http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("error on post %s\n", err.Error())
		panic(err)
	}

	if resp.StatusCode != 202 {
		b, _ := ioutil.ReadAll(resp.Body)
		fmt.Printf("body is %s!!\n", string(b))
	}
	return nil
}
