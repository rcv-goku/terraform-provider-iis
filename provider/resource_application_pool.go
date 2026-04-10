package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/maxjoehnk/terraform-provider-iis/iis"
)

const NameKey = "name"
const StatusKey = "status"

func resourceApplicationPool() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceApplicationPoolCreate,
		ReadContext:   resourceApplicationPoolRead,
		UpdateContext: resourceApplicationPoolUpdate,
		DeleteContext: resourceApplicationPoolDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			NameKey: {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			StatusKey: {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "started",
			},
			"managed_runtime_version": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: ".NET CLR version for the app pool (e.g., v4.0, v2.0)",
			},
			"auto_start": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
			},
			"pipeline_mode": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "Integrated",
				Description: "Managed pipeline mode: Integrated or Classic",
			},
			"enable_32bit_win64": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"queue_length": {
				Type:     schema.TypeInt,
				Optional: true,
				Default:  1000,
			},

			"cpu": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"limit": {
							Type:     schema.TypeInt,
							Optional: true,
							Computed: true,
						},
						"limit_interval": {
							Type:     schema.TypeInt,
							Optional: true,
							Computed: true,
						},
						"action": {
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
							Description: "CPU action: NoAction, KillW3wp, Throttle, ThrottleUnderLoad",
						},
						"processor_affinity_enabled": {
							Type:     schema.TypeBool,
							Optional: true,
							Computed: true,
						},
						"processor_affinity_mask32": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"processor_affinity_mask64": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
					},
				},
			},

			"process_model": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"idle_timeout": {
							Type:     schema.TypeInt,
							Optional: true,
							Computed: true,
						},
						"max_processes": {
							Type:     schema.TypeInt,
							Optional: true,
							Computed: true,
						},
						"pinging_enabled": {
							Type:     schema.TypeBool,
							Optional: true,
							Computed: true,
						},
						"ping_interval": {
							Type:     schema.TypeInt,
							Optional: true,
							Computed: true,
						},
						"ping_response_time": {
							Type:     schema.TypeInt,
							Optional: true,
							Computed: true,
						},
						"shutdown_time_limit": {
							Type:     schema.TypeInt,
							Optional: true,
							Computed: true,
						},
						"startup_time_limit": {
							Type:     schema.TypeInt,
							Optional: true,
							Computed: true,
						},
						"idle_timeout_action": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
					},
				},
			},

			"identity": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"identity_type": {
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
							Description: "ApplicationPoolIdentity, LocalSystem, LocalService, NetworkService, SpecificUser",
						},
						"username": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"load_user_profile": {
							Type:     schema.TypeBool,
							Optional: true,
							Computed: true,
						},
					},
				},
			},

			"recycling": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"disable_overlapped_recycle": {
							Type:     schema.TypeBool,
							Optional: true,
							Computed: true,
						},
						"disable_recycle_on_config_change": {
							Type:     schema.TypeBool,
							Optional: true,
							Computed: true,
						},
						"log_events": {
							Type:     schema.TypeList,
							Optional: true,
							Computed: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"time": {
										Type:     schema.TypeBool,
										Optional: true,
										Computed: true,
									},
									"requests": {
										Type:     schema.TypeBool,
										Optional: true,
										Computed: true,
									},
									"schedule": {
										Type:     schema.TypeBool,
										Optional: true,
										Computed: true,
									},
									"memory": {
										Type:     schema.TypeBool,
										Optional: true,
										Computed: true,
									},
									"isapi_unhealthy": {
										Type:     schema.TypeBool,
										Optional: true,
										Computed: true,
									},
									"on_demand": {
										Type:     schema.TypeBool,
										Optional: true,
										Computed: true,
									},
									"config_change": {
										Type:     schema.TypeBool,
										Optional: true,
										Computed: true,
									},
									"private_memory": {
										Type:     schema.TypeBool,
										Optional: true,
										Computed: true,
									},
								},
							},
						},
						"periodic_restart": {
							Type:     schema.TypeList,
							Optional: true,
							Computed: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"time_interval": {
										Type:     schema.TypeInt,
										Optional: true,
										Computed: true,
									},
									"private_memory": {
										Type:     schema.TypeInt,
										Optional: true,
										Computed: true,
									},
									"request_limit": {
										Type:     schema.TypeInt,
										Optional: true,
										Computed: true,
									},
									"virtual_memory": {
										Type:     schema.TypeInt,
										Optional: true,
										Computed: true,
									},
								},
							},
						},
					},
				},
			},

			"rapid_fail_protection": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"enabled": {
							Type:     schema.TypeBool,
							Optional: true,
							Computed: true,
						},
						"load_balancer_capabilities": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"interval": {
							Type:     schema.TypeInt,
							Optional: true,
							Computed: true,
						},
						"max_crashes": {
							Type:     schema.TypeInt,
							Optional: true,
							Computed: true,
						},
						"auto_shutdown_exe": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"auto_shutdown_params": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
					},
				},
			},

			"process_orphaning": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"enabled": {
							Type:     schema.TypeBool,
							Optional: true,
							Computed: true,
						},
						"orphan_action_exe": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"orphan_action_params": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

// expandAppPool builds an iis.ApplicationPool from Terraform resource data.
// The ID and Status fields are left for the caller to populate as needed.
func expandAppPool(d *schema.ResourceData) iis.ApplicationPool {
	pool := iis.ApplicationPool{
		Name:                  d.Get(NameKey).(string),
		ManagedRuntimeVersion: d.Get("managed_runtime_version").(string),
		AutoStart:             d.Get("auto_start").(bool),
		PipelineMode:          d.Get("pipeline_mode").(string),
		Enable32BitWin64:      d.Get("enable_32bit_win64").(bool),
		QueueLength:           int64(d.Get("queue_length").(int)),
	}

	if v, ok := d.GetOk("cpu"); ok {
		cpuList := v.([]interface{})
		if len(cpuList) == 1 {
			m := cpuList[0].(map[string]interface{})
			pool.CPU = iis.CPU{
				Limit:                    int64(m["limit"].(int)),
				LimitInterval:            int64(m["limit_interval"].(int)),
				Action:                   m["action"].(string),
				ProcessorAffinityEnabled: m["processor_affinity_enabled"].(bool),
				ProcessorAffinityMask32:  m["processor_affinity_mask32"].(string),
				ProcessorAffinityMask64:  m["processor_affinity_mask64"].(string),
			}
		}
	}

	if v, ok := d.GetOk("process_model"); ok {
		pmList := v.([]interface{})
		if len(pmList) == 1 {
			m := pmList[0].(map[string]interface{})
			pool.ProcessModel = iis.ProcessModel{
				IdleTimeout:       int64(m["idle_timeout"].(int)),
				MaxProcesses:      int64(m["max_processes"].(int)),
				PingingEnabled:    m["pinging_enabled"].(bool),
				PingInterval:      int64(m["ping_interval"].(int)),
				PingResponseTime:  int64(m["ping_response_time"].(int)),
				ShutdownTimeLimit: int64(m["shutdown_time_limit"].(int)),
				StartupTimeLimit:  int64(m["startup_time_limit"].(int)),
				IdleTimeoutAction: m["idle_timeout_action"].(string),
			}
		}
	}

	if v, ok := d.GetOk("identity"); ok {
		idList := v.([]interface{})
		if len(idList) == 1 {
			m := idList[0].(map[string]interface{})
			pool.Identity = iis.Identity{
				IdentityType:    m["identity_type"].(string),
				Username:        m["username"].(string),
				LoadUserProfile: m["load_user_profile"].(bool),
			}
		}
	}

	if v, ok := d.GetOk("recycling"); ok {
		rList := v.([]interface{})
		if len(rList) == 1 {
			m := rList[0].(map[string]interface{})
			recycling := iis.Recycling{
				DisableOverlappedRecycle:     m["disable_overlapped_recycle"].(bool),
				DisableRecycleOnConfigChange: m["disable_recycle_on_config_change"].(bool),
			}

			if leList, ok := m["log_events"].([]interface{}); ok && len(leList) == 1 {
				le := leList[0].(map[string]interface{})
				recycling.LogEvents = iis.LogEvents{
					Time:           le["time"].(bool),
					Requests:       le["requests"].(bool),
					Schedule:       le["schedule"].(bool),
					Memory:         le["memory"].(bool),
					IsapiUnhealthy: le["isapi_unhealthy"].(bool),
					OnDemand:       le["on_demand"].(bool),
					ConfigChange:   le["config_change"].(bool),
					PrivateMemory:  le["private_memory"].(bool),
				}
			}

			if prList, ok := m["periodic_restart"].([]interface{}); ok && len(prList) == 1 {
				pr := prList[0].(map[string]interface{})
				recycling.PeriodicRestart = iis.PeriodicRestart{
					TimeInterval:  int64(pr["time_interval"].(int)),
					PrivateMemory: int64(pr["private_memory"].(int)),
					RequestLimit:  int64(pr["request_limit"].(int)),
					VirtualMemory: int64(pr["virtual_memory"].(int)),
				}
			}

			pool.Recycling = recycling
		}
	}

	if v, ok := d.GetOk("rapid_fail_protection"); ok {
		rfpList := v.([]interface{})
		if len(rfpList) == 1 {
			m := rfpList[0].(map[string]interface{})
			pool.RapidFailProtection = iis.RapidFailProtection{
				Enabled:                  m["enabled"].(bool),
				LoadBalancerCapabilities: m["load_balancer_capabilities"].(string),
				Interval:                 int64(m["interval"].(int)),
				MaxCrashes:               int64(m["max_crashes"].(int)),
				AutoShutdownExe:          m["auto_shutdown_exe"].(string),
				AutoShutdownParams:       m["auto_shutdown_params"].(string),
			}
		}
	}

	if v, ok := d.GetOk("process_orphaning"); ok {
		poList := v.([]interface{})
		if len(poList) == 1 {
			m := poList[0].(map[string]interface{})
			pool.ProcessOrphaning = iis.ProcessOrphaning{
				Enabled:            m["enabled"].(bool),
				OrphanActionExe:    m["orphan_action_exe"].(string),
				OrphanActionParams: m["orphan_action_params"].(string),
			}
		}
	}

	return pool
}

// flattenAppPool writes all fields from an iis.ApplicationPool into Terraform state.
func flattenAppPool(d *schema.ResourceData, pool *iis.ApplicationPool) error {
	if err := d.Set(NameKey, pool.Name); err != nil {
		return err
	}
	if err := d.Set(StatusKey, pool.Status); err != nil {
		return err
	}
	if err := d.Set("managed_runtime_version", pool.ManagedRuntimeVersion); err != nil {
		return err
	}
	if err := d.Set("auto_start", pool.AutoStart); err != nil {
		return err
	}
	if err := d.Set("pipeline_mode", pool.PipelineMode); err != nil {
		return err
	}
	if err := d.Set("enable_32bit_win64", pool.Enable32BitWin64); err != nil {
		return err
	}
	if err := d.Set("queue_length", int(pool.QueueLength)); err != nil {
		return err
	}

	cpu := []map[string]interface{}{{
		"limit":                      int(pool.CPU.Limit),
		"limit_interval":             int(pool.CPU.LimitInterval),
		"action":                     pool.CPU.Action,
		"processor_affinity_enabled": pool.CPU.ProcessorAffinityEnabled,
		"processor_affinity_mask32":  pool.CPU.ProcessorAffinityMask32,
		"processor_affinity_mask64":  pool.CPU.ProcessorAffinityMask64,
	}}
	if err := d.Set("cpu", cpu); err != nil {
		return err
	}

	pm := []map[string]interface{}{{
		"idle_timeout":        int(pool.ProcessModel.IdleTimeout),
		"max_processes":       int(pool.ProcessModel.MaxProcesses),
		"pinging_enabled":     pool.ProcessModel.PingingEnabled,
		"ping_interval":       int(pool.ProcessModel.PingInterval),
		"ping_response_time":  int(pool.ProcessModel.PingResponseTime),
		"shutdown_time_limit": int(pool.ProcessModel.ShutdownTimeLimit),
		"startup_time_limit":  int(pool.ProcessModel.StartupTimeLimit),
		"idle_timeout_action": pool.ProcessModel.IdleTimeoutAction,
	}}
	if err := d.Set("process_model", pm); err != nil {
		return err
	}

	identity := []map[string]interface{}{{
		"identity_type":     pool.Identity.IdentityType,
		"username":          pool.Identity.Username,
		"load_user_profile": pool.Identity.LoadUserProfile,
	}}
	if err := d.Set("identity", identity); err != nil {
		return err
	}

	logEvents := []map[string]interface{}{{
		"time":           pool.Recycling.LogEvents.Time,
		"requests":       pool.Recycling.LogEvents.Requests,
		"schedule":       pool.Recycling.LogEvents.Schedule,
		"memory":         pool.Recycling.LogEvents.Memory,
		"isapi_unhealthy": pool.Recycling.LogEvents.IsapiUnhealthy,
		"on_demand":      pool.Recycling.LogEvents.OnDemand,
		"config_change":  pool.Recycling.LogEvents.ConfigChange,
		"private_memory": pool.Recycling.LogEvents.PrivateMemory,
	}}
	periodicRestart := []map[string]interface{}{{
		"time_interval":  int(pool.Recycling.PeriodicRestart.TimeInterval),
		"private_memory": int(pool.Recycling.PeriodicRestart.PrivateMemory),
		"request_limit":  int(pool.Recycling.PeriodicRestart.RequestLimit),
		"virtual_memory": int(pool.Recycling.PeriodicRestart.VirtualMemory),
	}}
	recycling := []map[string]interface{}{{
		"disable_overlapped_recycle":      pool.Recycling.DisableOverlappedRecycle,
		"disable_recycle_on_config_change": pool.Recycling.DisableRecycleOnConfigChange,
		"log_events":                      logEvents,
		"periodic_restart":                periodicRestart,
	}}
	if err := d.Set("recycling", recycling); err != nil {
		return err
	}

	rfp := []map[string]interface{}{{
		"enabled":                    pool.RapidFailProtection.Enabled,
		"load_balancer_capabilities": pool.RapidFailProtection.LoadBalancerCapabilities,
		"interval":                   int(pool.RapidFailProtection.Interval),
		"max_crashes":                int(pool.RapidFailProtection.MaxCrashes),
		"auto_shutdown_exe":          pool.RapidFailProtection.AutoShutdownExe,
		"auto_shutdown_params":       pool.RapidFailProtection.AutoShutdownParams,
	}}
	if err := d.Set("rapid_fail_protection", rfp); err != nil {
		return err
	}

	po := []map[string]interface{}{{
		"enabled":              pool.ProcessOrphaning.Enabled,
		"orphan_action_exe":    pool.ProcessOrphaning.OrphanActionExe,
		"orphan_action_params": pool.ProcessOrphaning.OrphanActionParams,
	}}
	if err := d.Set("process_orphaning", po); err != nil {
		return err
	}

	return nil
}

func resourceApplicationPoolCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*iis.Client)
	pool := expandAppPool(d)
	tflog.Debug(ctx, "Creating application pool: "+toJSON(pool))
	created, err := client.CreateAppPool(ctx, pool)
	if err != nil {
		return diag.FromErr(err)
	}
	tflog.Debug(ctx, "Created application pool: "+toJSON(created))
	d.SetId(created.ID)

	// IIS Admin API only accepts name on POST — PATCH the full config after creation
	pool.Status = d.Get(StatusKey).(string)
	tflog.Debug(ctx, "Applying full configuration via PATCH")
	_, err = client.UpdateAppPool(ctx, created.ID, pool)
	if err != nil {
		return diag.FromErr(err)
	}

	return resourceApplicationPoolRead(ctx, d, m)
}

func resourceApplicationPoolRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*iis.Client)
	id := d.Id()
	appPool, err := client.ReadAppPool(ctx, id)
	if err != nil {
		if iis.IsNotFoundError(err) {
			tflog.Warn(ctx, "Application pool not found, removing from state: "+id)
			d.SetId("")
			return nil
		}
		return diag.FromErr(err)
	}
	tflog.Debug(ctx, "Read application pool: "+toJSON(appPool))
	if err := flattenAppPool(d, appPool); err != nil {
		return diag.FromErr(err)
	}
	return nil
}

func resourceApplicationPoolUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*iis.Client)
	id := d.Id()

	pool := expandAppPool(d)
	pool.Status = d.Get(StatusKey).(string)
	tflog.Debug(ctx, "Updating application pool: "+toJSON(id)+", payload: "+toJSON(pool))
	updated, err := client.UpdateAppPool(ctx, id, pool)
	if err != nil {
		return diag.FromErr(err)
	}
	tflog.Debug(ctx, "Updated application pool: "+toJSON(updated))
	return resourceApplicationPoolRead(ctx, d, m)
}

func resourceApplicationPoolDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*iis.Client)
	id := d.Id()
	tflog.Debug(ctx, "Deleting application pool: "+toJSON(id))
	err := client.DeleteAppPool(ctx, id)
	if err != nil {
		return diag.FromErr(err)
	}
	tflog.Debug(ctx, "Deleted application pool: "+toJSON(id))
	return nil
}
