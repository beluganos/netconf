// -*- coding: utf-8 -*-

// Copyright (C) 2018 Nippon Telegraph and Telephone Corporation.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//    http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or
// implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package lxdlib

import (
	lxd "github.com/lxc/lxd/client"
	api "github.com/lxc/lxd/shared/api"
	log "github.com/sirupsen/logrus"
)

type Client struct {
	Base   string
	Server lxd.ContainerServer
}

func NewClient(base string) *Client {
	return &Client{
		Base:   base,
		Server: nil,
	}
}

func (c *Client) Connect() error {
	conn, err := Connect()
	if err != nil {
		return err
	}

	c.Server = conn
	return nil
}

func (c *Client) GetProfile(name string) (*api.Profile, string, error) {
	return c.Server.GetProfile(name)
}

func (c *Client) UpdateProfile(name string, profile api.ProfilePut) error {
	_, etag, err := c.Server.GetProfile(name)
	if err != nil {
		return nil
	}

	return c.Server.UpdateProfile(name, profile, etag)
}

func (c *Client) CreateOrUpdateProfile(name string, profile *api.ProfilePut) error {
	if _, etag, err := c.Server.GetProfile(name); err == nil {
		return c.Server.UpdateProfile(name, *profile, etag)
	}

	return c.Server.CreateProfile(api.ProfilesPost{
		ProfilePut: *profile,
		Name:       name,
	})
}

func (c *Client) InitializeProfile(name string, logdir, mntIf, bridgeIf string) (*api.Profile, string, error) {
	profile := NewDefaultProfile(mntIf, bridgeIf)
	SetMontToProfle(profile, "/var/log", logdir)

	if err := c.CreateOrUpdateProfile(name, profile); err != nil {
		return nil, "", err
	}

	return c.GetProfile(name)
}

func (c *Client) DeleteProfile(name string) error {
	log.Debugf("Profile(%s) deleted", name)
	return c.Server.DeleteProfile(name)
}

func (c *Client) GetContainer(name string) (*api.Container, string, error) {
	return c.Server.GetContainer(name)
}

func (c *Client) CreateContainer(name string) (*api.Container, string, error) {

	err := func() error {
		container := NewDefaultContainer(name)
		if _, etag, err := c.Server.GetContainer(name); err == nil {
			op, err := c.Server.UpdateContainer(name, container.ContainerPut, etag)
			if err != nil {
				return err
			}
			log.Debugf("Container(%s) initailized.", name)
			return op.Wait()

		} else {
			op, err := CreateContainer(c.Server, c.Base, container)
			if err != nil {
				return err
			}
			log.Debugf("Container(%s) created.", name)
			return op.Wait()
		}
	}()

	if err != nil {
		return nil, "", err
	}

	return c.GetContainer(name)
}

func (c *Client) DeleteContainer(name string) error {
	op, err := c.Server.DeleteContainer(name)
	if err != nil {
		log.Errorf("DeleteContainer error. %s", err)
		return err
	}

	log.Debugf("Container(%s) deleted.", name)
	return op.Wait()
}

func (c *Client) UpdateContainerState(name string, action string) error {
	state, _, err := c.Server.GetContainerState(name)
	if err != nil {
		log.Errorf("GetContainerState error. %s", err)
		return err
	}

	log.Debugf("GetContainerState: current state %s", state.Status)

	var same_state bool = false
	switch action {
	case "start":
		same_state = (state.Status == api.Running.String())
	case "stop":
		same_state = (state.Status == api.Stopped.String())
	}

	if same_state {
		log.Debugf("Container state no need to change.")
		return nil
	}

	statePut := api.ContainerStatePut{
		Action:  action,
		Timeout: -1,
		Force:   true,
		//Stateful: true,
	}
	op, err := c.Server.UpdateContainerState(name, statePut, "")
	if err != nil {
		log.Errorf("UpdateContainerState error. %s", err)
		return err
	}

	log.Debugf("UpdateContainerState ok. %s %s", name, action)
	return op.Wait()
}
