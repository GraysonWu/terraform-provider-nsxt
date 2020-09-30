/* Copyright © 2017 VMware, Inc. All Rights Reserved.
   SPDX-License-Identifier: MPL-2.0 */

package nsxt

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccDataSourceNsxtMacPool_basic(t *testing.T) {
	macPoolName := getMacPoolName()
	testResourceName := "data.nsxt_mac_pool.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testAccOnlyLocalManager(t); testAccTestMP(t); testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccNSXMacPoolReadTemplate(macPoolName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(testResourceName, "display_name", macPoolName),
					resource.TestCheckResourceAttrSet(testResourceName, "id"),
				),
			},
		},
	})
}

func testAccNSXMacPoolReadTemplate(name string) string {
	return fmt.Sprintf(`
data "nsxt_mac_pool" "test" {
  display_name = "%s"
}`, name)
}
