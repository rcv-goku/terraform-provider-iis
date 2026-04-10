package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/maxjoehnk/terraform-provider-iis/iis"
)

const ApplicationKey = "application"

func resourceVirtualDirectory() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceVirtualDirectoryCreate,
		ReadContext:   resourceVirtualDirectoryRead,
		UpdateContext: resourceVirtualDirectoryUpdate,
		DeleteContext: resourceVirtualDirectoryDelete,

		Schema: map[string]*schema.Schema{
			PathKey: {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			PhysicalPathKey: {
				Type:     schema.TypeString,
				Required: true,
			},
			WebsiteKey: {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			ApplicationKey: {
				Type:     schema.TypeString,
				Optional: true,
			},
			"location": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceVirtualDirectoryCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*iis.Client)
	req := buildCreateVirtualDirectoryRequest(d)
	tflog.Debug(ctx, "Creating virtual directory: "+toJSON(req))
	vdir, err := client.CreateVirtualDirectory(ctx, req)
	if err != nil {
		return diag.FromErr(err)
	}
	tflog.Debug(ctx, "Created virtual directory: "+toJSON(vdir))
	d.SetId(vdir.ID)
	return resourceVirtualDirectoryRead(ctx, d, m)
}

func resourceVirtualDirectoryRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*iis.Client)
	vdir, err := client.ReadVirtualDirectory(ctx, d.Id())
	if err != nil {
		if iis.IsNotFoundError(err) {
			tflog.Warn(ctx, "Virtual directory not found, removing from state: "+d.Id())
			d.SetId("")
			return nil
		}
		return diag.FromErr(err)
	}
	tflog.Debug(ctx, "Read virtual directory: "+toJSON(vdir))
	if err = d.Set(PathKey, vdir.Path); err != nil {
		return diag.FromErr(err)
	}
	if err = d.Set(PhysicalPathKey, vdir.PhysicalPath); err != nil {
		return diag.FromErr(err)
	}
	if err = d.Set(WebsiteKey, vdir.Website.ID); err != nil {
		return diag.FromErr(err)
	}
	if vdir.WebApp.ID != "" {
		if err = d.Set(ApplicationKey, vdir.WebApp.ID); err != nil {
			return diag.FromErr(err)
		}
	}
	if err = d.Set("location", vdir.Location); err != nil {
		return diag.FromErr(err)
	}
	return nil
}

func resourceVirtualDirectoryUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*iis.Client)
	id := d.Id()

	if d.HasChange(PhysicalPathKey) {
		req := iis.UpdateVirtualDirectoryRequest{
			PhysicalPath: d.Get(PhysicalPathKey).(string),
		}
		tflog.Debug(ctx, "Updating virtual directory: "+toJSON(id))
		vdir, err := client.UpdateVirtualDirectory(ctx, id, req)
		if err != nil {
			return diag.FromErr(err)
		}
		tflog.Debug(ctx, "Updated virtual directory: "+toJSON(vdir))
	}

	return resourceVirtualDirectoryRead(ctx, d, m)
}

func resourceVirtualDirectoryDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*iis.Client)
	id := d.Id()
	tflog.Debug(ctx, "Deleting virtual directory: "+toJSON(id))
	if err := client.DeleteVirtualDirectory(ctx, id); err != nil {
		return diag.FromErr(err)
	}
	tflog.Debug(ctx, "Deleted virtual directory: "+toJSON(id))
	return nil
}

func buildCreateVirtualDirectoryRequest(d *schema.ResourceData) iis.CreateVirtualDirectoryRequest {
	req := iis.CreateVirtualDirectoryRequest{
		Path:         d.Get(PathKey).(string),
		PhysicalPath: d.Get(PhysicalPathKey).(string),
		Website:      iis.Reference{ID: d.Get(WebsiteKey).(string)},
	}
	if appID, ok := d.GetOk(ApplicationKey); ok {
		req.WebApp = iis.Reference{ID: appID.(string)}
	}
	return req
}
