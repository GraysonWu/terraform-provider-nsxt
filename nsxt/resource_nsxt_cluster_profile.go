/* Copyright © 2018 VMware, Inc. All Rights Reserved.
   SPDX-License-Identifier: MPL-2.0 */

package nsxt

import (
	"fmt"
	"log"
	"net/http"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/vmware/go-vmware-nsxt/manager"
)

func resourceNsxtClusterProfile() *schema.Resource {
	return &schema.Resource{
		Create: resourceNsxtClusterProfileCreate,
		Read:   resourceNsxtClusterProfileRead,
		Update: resourceNsxtClusterProfileUpdate,
		Delete: resourceNsxtClusterProfileDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		DeprecationMessage: mpObjectResourceDeprecationMessage,
		Schema: map[string]*schema.Schema{
			"revision": getRevisionSchema(),
			"description": {
				Type:        schema.TypeString,
				Description: "Description of this resource",
				Optional:    true,
			},
			"display_name": {
				Type:        schema.TypeString,
				Description: "The display name of this resource. Defaults to ID if not set",
				Optional:    true,
				Computed:    true,
			},
			"tag": getTagsSchema(),
			"resource_type": {
				Type:        schema.TypeString,
				Description: "Supported cluster profiles(Enum: EdgeHighAvailabilityProfile, BridgeHighAvailabilityClusterProfile)",
				Required:    true,
			},
		},
	}
}

func resourceNsxtClusterProfileCreate(d *schema.ResourceData, m interface{}) error {
	nsxClient := m.(nsxtClients).NsxtClient
	if nsxClient == nil {
		return resourceNotSupportedError()
	}

	description := d.Get("description").(string)
	displayName := d.Get("display_name").(string)
	tags := getTagsFromSchema(d)
	resourceType := d.Get("resource_type").(string)
	clusterProfile := manager.ClusterProfile{
		DisplayName:  displayName,
		Description:  description,
		Tags:         tags,
		ResourceType: resourceType,
	}

	clusterProfile, resp, err := nsxClient.NetworkTransportApi.CreateClusterProfile(nsxClient.Context, clusterProfile)

	if err != nil {
		return fmt.Errorf("Error during ClusterProfile create: %v", err)
	}

	if resp.StatusCode != http.StatusCreated {
		return fmt.Errorf("Unexpected status returned during ClusterProfile create: %v", resp.StatusCode)
	}
	d.SetId(clusterProfile.Id)

	return resourceNsxtClusterProfileRead(d, m)
}

func resourceNsxtClusterProfileRead(d *schema.ResourceData, m interface{}) error {
	nsxClient := m.(nsxtClients).NsxtClient
	if nsxClient == nil {
		return resourceNotSupportedError()
	}

	id := d.Id()
	if id == "" {
		return fmt.Errorf("Error obtaining logical object id")
	}

	clusterProfile, resp, err := nsxClient.NetworkTransportApi.GetClusterProfile(nsxClient.Context, id)
	if resp != nil && resp.StatusCode == http.StatusNotFound {
		log.Printf("[DEBUG] ClusterProfile %s not found", id)
		d.SetId("")
		return nil
	}
	if err != nil {
		return fmt.Errorf("Error during ClusterProfile read: %v", err)
	}

	d.Set("revision", clusterProfile.Revision)
	d.Set("description", clusterProfile.Description)
	d.Set("display_name", clusterProfile.DisplayName)
	setTagsInSchema(d, clusterProfile.Tags)
	d.Set("resource_type", clusterProfile.ResourceType)

	return nil
}

func resourceNsxtClusterProfileUpdate(d *schema.ResourceData, m interface{}) error {
	nsxClient := m.(nsxtClients).NsxtClient
	if nsxClient == nil {
		return resourceNotSupportedError()
	}

	id := d.Id()
	if id == "" {
		return fmt.Errorf("Error obtaining logical object id")
	}

	displayName := d.Get("display_name").(string)
	description := d.Get("description").(string)
	tags := getTagsFromSchema(d)
	revision := int64(d.Get("revision").(int))
	resourceType := d.Get("resource_type").(string)
	clusterProfile := manager.ClusterProfile{
		DisplayName:  displayName,
		Description:  description,
		Tags:         tags,
		Revision:     revision,
		ResourceType: resourceType,
	}

	_, resp, err := nsxClient.NetworkTransportApi.UpdateClusterProfile(nsxClient.Context, id, clusterProfile)

	if err != nil || resp.StatusCode == http.StatusNotFound {
		return fmt.Errorf("Error during ClusterProfile update: %v", err)
	}

	return resourceNsxtClusterProfileRead(d, m)
}

func resourceNsxtClusterProfileDelete(d *schema.ResourceData, m interface{}) error {
	nsxClient := m.(nsxtClients).NsxtClient
	if nsxClient == nil {
		return resourceNotSupportedError()
	}

	id := d.Id()
	if id == "" {
		return fmt.Errorf("Error obtaining logical object id")
	}

	resp, err := nsxClient.NetworkTransportApi.DeleteClusterProfile(nsxClient.Context, id)
	if err != nil {
		return fmt.Errorf("Error during ClusterProfile delete: %v", err)
	}

	if resp.StatusCode == http.StatusNotFound {
		log.Printf("[DEBUG] ClusterProfile %s not found", id)
		d.SetId("")
	}
	return nil
}
