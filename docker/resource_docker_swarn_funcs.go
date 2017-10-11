package docker

import (
	"github.com/docker/docker/api/types/swarm"
	dc "github.com/fsouza/go-dockerclient"
	"github.com/hashicorp/terraform/helper/schema"
)

func resourceDockerSwarmExists(d *schema.ResourceData, meta interface{}) (bool, error) {
	client := meta.(*dc.Client)

	_, err := client.InspectSwarm(nil)
	if err != nil {
		return false, err
	}

	return true, nil
}

func resourceDockerSwarmCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*dc.Client)

	joinToken, hasJoinToken := d.GetOk("join_token")

	if !hasJoinToken {
		initOpts := dc.InitSwarmOptions{
			InitRequest: swarm.InitRequest{
				AdvertiseAddr:    "",
				AutoLockManagers: true,
				Availability:     swarm.NodeAvailabilityDrain,
				ForceNewCluster:  true,
				ListenAddr:       "",
			},
		}

		_, err := client.InitSwarm(initOpts)
		if err != nil {
			return err
		}
	} else {
		joinOpts := dc.JoinSwarmOptions{
			JoinRequest: swarm.JoinRequest{
				AdvertiseAddr: "",
				Availability:  swarm.NodeAvailabilityActive,
				JoinToken:     joinToken.(string),
				ListenAddr:    "",
				RemoteAddrs:   []string{""},
			},
		}

		err := client.JoinSwarm(joinOpts)
		if err != nil {
			return err
		}
	}

	return resourceDockerSwarmRead(d, meta)
}

func resourceDockerSwarmRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*dc.Client)

	swarm, err := client.InspectSwarm(nil)
	if err != nil {
		return err
	}

	d.SetId(swarm.ID)

	d.Set("manager_join_token", swarm.JoinTokens.Manager)
	d.Set("worker_join_token", swarm.JoinTokens.Worker)

	return nil
}

func resourceDockerSwarmDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*dc.Client)

	leaveOpts := dc.LeaveSwarmOptions{
		Force: true,
	}

	return client.LeaveSwarm(leaveOpts)
}
