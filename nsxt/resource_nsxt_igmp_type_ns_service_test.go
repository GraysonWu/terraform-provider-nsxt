/* Copyright © 2017 VMware, Inc. All Rights Reserved.
   SPDX-License-Identifier: MPL-2.0 */

package nsxt

import (
	"fmt"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"github.com/vmware/go-vmware-nsxt"
	"net/http"
	"testing"
)

func TestAccResourceNsxtIgmpTypeNsService_basic(t *testing.T) {
	serviceName := fmt.Sprintf("test-nsx-igmp-service")
	updateServiceName := fmt.Sprintf("%s-update", serviceName)
	testResourceName := "nsxt_igmp_type_ns_service.test"

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		CheckDestroy: func(state *terraform.State) error {
			return testAccNSXIgmpServiceCheckDestroy(state, serviceName)
		},
		Steps: []resource.TestStep{
			{
				Config: testAccNSXIgmpServiceCreateTemplate(serviceName),
				Check: resource.ComposeTestCheckFunc(
					testAccNSXIgmpServiceExists(serviceName, testResourceName),
					resource.TestCheckResourceAttr(testResourceName, "display_name", serviceName),
					resource.TestCheckResourceAttr(testResourceName, "description", "igmp service"),
					resource.TestCheckResourceAttr(testResourceName, "tag.#", "1"),
				),
			},
			{
				Config: testAccNSXIgmpServiceCreateTemplate(updateServiceName),
				Check: resource.ComposeTestCheckFunc(
					testAccNSXIgmpServiceExists(updateServiceName, testResourceName),
					resource.TestCheckResourceAttr(testResourceName, "display_name", updateServiceName),
					resource.TestCheckResourceAttr(testResourceName, "description", "igmp service"),
					resource.TestCheckResourceAttr(testResourceName, "tag.#", "1"),
				),
			},
		},
	})
}

func testAccNSXIgmpServiceExists(display_name string, resourceName string) resource.TestCheckFunc {
	return func(state *terraform.State) error {

		nsxClient := testAccProvider.Meta().(*nsxt.APIClient)

		rs, ok := state.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("NSX igmp service resource %s not found in resources", resourceName)
		}

		resourceID := rs.Primary.ID
		if resourceID == "" {
			return fmt.Errorf("NSX igmp service resource ID not set in resources ")
		}

		service, responseCode, err := nsxClient.GroupingObjectsApi.ReadIgmpTypeNSService(nsxClient.Context, resourceID)
		if err != nil {
			return fmt.Errorf("Error while retrieving igmp service ID %s. Error: %v", resourceID, err)
		}

		if responseCode.StatusCode != http.StatusOK {
			return fmt.Errorf("Error while checking if igmp service %s exists. HTTP return code was %d", resourceID, responseCode.StatusCode)
		}

		if display_name == service.DisplayName {
			return nil
		}
		return fmt.Errorf("NSX igmp ns service %s wasn't found", display_name)
	}
}

func testAccNSXIgmpServiceCheckDestroy(state *terraform.State, display_name string) error {
	nsxClient := testAccProvider.Meta().(*nsxt.APIClient)

	for _, rs := range state.RootModule().Resources {

		if rs.Type != "nsxt_igmp_set_ns_service" {
			continue
		}

		resourceID := rs.Primary.Attributes["id"]
		service, responseCode, err := nsxClient.GroupingObjectsApi.ReadIgmpTypeNSService(nsxClient.Context, resourceID)
		if err != nil {
			if responseCode.StatusCode != http.StatusOK {
				return nil
			}
			return fmt.Errorf("Error while retrieving L4 ns service ID %s. Error: %v", resourceID, err)
		}

		if display_name == service.DisplayName {
			return fmt.Errorf("NSX L4 ns service %s still exists", display_name)
		}
	}
	return nil
}

func testAccNSXIgmpServiceCreateTemplate(serviceName string) string {
	return fmt.Sprintf(`
resource "nsxt_igmp_type_ns_service" "test" {
    description = "igmp service"
    display_name = "%s"
    tag {
    	scope = "scope1"
        tag = "tag1"
    }
}`, serviceName)
}
