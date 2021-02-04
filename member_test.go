// Copyright (c) 2021, ZeroTier, Inc.
// All rights reserved.
//
// Redistribution and use in source and binary forms, with or without
// modification, are permitted provided that the following conditions are met:
//
// 1. Redistributions of source code must retain the above copyright notice, this
//    list of conditions and the following disclaimer.
//
// 2. Redistributions in binary form must reproduce the above copyright notice,
//    this list of conditions and the following disclaimer in the documentation
//    and/or other materials provided with the distribution.
//
// 3. Neither the name of the copyright holder nor the names of its
//    contributors may be used to endorse or promote products derived from
//    this software without specific prior written permission.
//
// THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS "AS IS"
// AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE
// IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE ARE
// DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT HOLDER OR CONTRIBUTORS BE LIABLE
// FOR ANY DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR CONSEQUENTIAL
// DAMAGES (INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR
// SERVICES; LOSS OF USE, DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER
// CAUSED AND ON ANY THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY,
// OR TORT (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE
// OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.

package ztcentral

import (
	"context"
	"testing"
	"time"

	"github.com/zerotier/go-ztcentral/pkg/testutil"
	"github.com/zerotier/go-ztidentity"
)

func TestGetMember(t *testing.T) {
	testutil.NeedsToken(t)

	c := NewClient(testutil.InitTokenFromEnv())

	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	net, err := c.NewNetwork(ctx, "get-member-network")
	if err != nil {
		t.Fatal(err)
	}
	defer c.DeleteNetwork(ctx, net.Config.ID)

	if _, err := c.GetMember(ctx, net.Config.ID, "123456789"); err == nil {
		t.Fatal("Tried to get a fake member and succeeded")
	}

	aliceID := ztidentity.NewZeroTierIdentity()

	alice, err := c.CreateAuthorizedMember(ctx, net.Config.ID, aliceID.IDString(), "alice")
	if err != nil {
		t.Fatal(err)
	}

	res, err := c.GetMember(ctx, net.Config.ID, alice.NodeID)
	if err != nil {
		t.Fatal(err)
	}

	if res.NetworkID != net.Config.ID {
		t.Fatal("network ID of member was not equivalent")
	}

	if res.NodeID != alice.NodeID {
		t.Fatal("member IDs were not equivalent")
	}
}

func TestGetMembers(t *testing.T) {
	testutil.NeedsToken(t)

	c := NewClient(testutil.InitTokenFromEnv())

	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	net, err := c.NewNetwork(ctx, "get-members-network")
	if err != nil {
		t.Fatal(err)
	}
	defer c.DeleteNetwork(ctx, net.Config.ID)

	users := map[string]ztidentity.ZeroTierIdentity{
		"awk":    ztidentity.NewZeroTierIdentity(),
		"bash":   ztidentity.NewZeroTierIdentity(),
		"cc":     ztidentity.NewZeroTierIdentity(),
		"dpkg":   ztidentity.NewZeroTierIdentity(),
		"edlin":  ztidentity.NewZeroTierIdentity(),
		"finger": ztidentity.NewZeroTierIdentity(),
		"gopher": ztidentity.NewZeroTierIdentity(),
	}

	for name, id := range users {
		_, err := c.CreateAuthorizedMember(ctx, net.Config.ID, id.IDString(), name)
		if err != nil {
			t.Fatal(err)
		}
	}

	members, err := c.GetMembers(ctx, net.Config.ID)
	if err != nil {
		t.Fatal(err)
	}

	for _, member := range members {
		id, ok := users[member.Name]
		if !ok {
			t.Fatal("could not find member in pre-populated table")
		}

		if id.IDString() != member.Config.ID {
			t.Fatalf("IDs were not equal for member %q", member.Name)
		}
	}
}
