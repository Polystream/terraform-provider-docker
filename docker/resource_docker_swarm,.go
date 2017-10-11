package docker

import (
	"github.com/hashicorp/terraform/helper/schema"
)

func resourceDockerSwarm() *schema.Resource {
	return &schema.Resource{
		Create: resourceDockerSwarmCreate,
		Read:   resourceDockerSwarmRead,
		Delete: resourceDockerSwarmDelete,
		Exists: resourceDockerSwarmExists,

		Schema: map[string]*schema.Schema{},
	}
}
