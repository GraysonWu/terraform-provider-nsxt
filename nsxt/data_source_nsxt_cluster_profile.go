/* Copyright © 2017 VMware, Inc. All Rights Reserved.
   SPDX-License-Identifier: MPL-2.0 */

package nsxt

import (
	"fmt"
	"strings"

	"net/http"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/vmware/go-vmware-nsxt/manager"
)

func dataSourceNsxtClusterProfile() *schema.Resource {
	return &schema.Resource{
		Read:               dataSourceNsxtClusterProfileRead,
		DeprecationMessage: mpObjectDataSourceDeprecationMessage,
		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeString,
				Description: "Unique ID of this resource",
				Optional:    true,
				Computed:    true,
			},
			"display_name": {
				Type:        schema.TypeString,
				Description: "The display name of this resource",
				Optional:    true,
				Computed:    true,
			},
			"description": {
				Type:        schema.TypeString,
				Description: "Description of this resource",
				Optional:    true,
				Computed:    true,
			},
			"resource_type": {
				Type:        schema.TypeString,
				Description: "Supported cluster profiles(Enum: EdgeHighAvailabilityProfile, BridgeHighAvailabilityClusterProfile)",
				Required:    true,
			},
		},
	}
}

func dataSourceNsxtClusterProfileRead(d *schema.ResourceData, m interface{}) error {
	// Read an cluster profile by name or id
	nsxClient := m.(nsxtClients).NsxtClient
	if nsxClient == nil {
		return dataSourceNotSupportedError()
	}

	objID := d.Get("id").(string)
	objName := d.Get("display_name").(string)
	var obj manager.ClusterProfile
	if objID != "" {
		// Get by id
		objGet, resp, err := nsxClient.NetworkTransportApi.GetClusterProfile(nsxClient.Context, objID)

		if resp != nil && resp.StatusCode == http.StatusNotFound {
			return fmt.Errorf("Cluster profile %s was not found", objID)
		}
		if err != nil {
			return fmt.Errorf("Error while reading cluster profile %s: %v", objID, err)
		}
		obj = objGet

	} else if objName == "" {
		return fmt.Errorf("Error obtaining cluster profile ID or name during read")
	} else {
		// Get by full name/prefix
		// TODO use 2nd parameter localVarOptionals for paging
		objList, _, err := nsxClient.NetworkTransportApi.ListClusterProfiles(nsxClient.Context, nil)
		if err != nil {
			return fmt.Errorf("Error while reading cluster profile: %v", err)
		}
		// go over the list to find the correct one (prefer a perfect match. If not - prefix match)
		var perfectMatch []manager.ClusterProfile
		var prefixMatch []manager.ClusterProfile
		for _, objInList := range objList.Results {
			if strings.HasPrefix(objInList.DisplayName, objName) {
				prefixMatch = append(prefixMatch, objInList)
			}
			if objInList.DisplayName == objName {
				perfectMatch = append(perfectMatch, objInList)
			}
		}
		if len(perfectMatch) > 0 {
			if len(perfectMatch) > 1 {
				return fmt.Errorf("Found multiple cluster profiles with name '%s'", objName)
			}
			obj = perfectMatch[0]
		} else if len(prefixMatch) > 0 {
			if len(prefixMatch) > 1 {
				return fmt.Errorf("Found multiple cluster profiles with name starting with '%s'", objName)
			}
			obj = prefixMatch[0]
		} else {
			return fmt.Errorf("Cluster profile with name '%s' was not found", objName)
		}
	}

	d.SetId(obj.Id)
	d.Set("display_name", obj.DisplayName)
	d.Set("description", obj.Description)
	d.Set("resource_type", obj.ResourceType)

	return nil
}
