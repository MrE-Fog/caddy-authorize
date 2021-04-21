// Copyright 2020 Paul Greenberg greenpau@outlook.com
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package testutils

import (
	"context"
	"fmt"
	"github.com/greenpau/caddy-auth-jwt/pkg/acl"
	"github.com/greenpau/caddy-auth-jwt/pkg/kms"
	"github.com/greenpau/caddy-auth-jwt/pkg/user"
	"github.com/greenpau/caddy-auth-jwt/pkg/utils"
	"net/http"
	"time"
)

// InjectedTestToken is an instance of injected token.
type InjectedTestToken struct {
	Name string
	// The locations to inject a token in this test.
	Location string
	// The basic user claims.
	User *user.User
}

// NewInjectedTestToken returns an instance of injected token.
func NewInjectedTestToken(name, location, cfg string) *InjectedTestToken {
	cfg = `{
        "exp": ` + fmt.Sprintf("%d", time.Now().Add(10*time.Minute).Unix()) + `,
        "iat": ` + fmt.Sprintf("%d", time.Now().Add(10*time.Minute*-1).Unix()) + `,
        "nbf": ` + fmt.Sprintf("%d", time.Date(2015, 10, 10, 12, 0, 0, 0, time.UTC).Unix()) + `,
        ` + cfg + `
        "email":  "smithj@outlook.com",
        "origin": "localhost",
        "sub":    "smithj@outlook.com",
        "roles": "anonymous guest"
    }`
	usr, err := user.NewUser(cfg)
	if err != nil {
		panic(err)
	}
	tkn := &InjectedTestToken{
		Name:     name,
		Location: location,
		User:     usr,
	}
	return tkn
}

// NewTestUser returns test User with claims.
func NewTestUser() *user.User {
	cfg := `{
        "exp": ` + fmt.Sprintf("%d", time.Now().Add(10*time.Minute).Unix()) + `,
        "iat": ` + fmt.Sprintf("%d", time.Now().Add(10*time.Minute*-1).Unix()) + `,
        "nbf": ` + fmt.Sprintf("%d", time.Date(2015, 10, 10, 12, 0, 0, 0, time.UTC).Unix()) + `,
        "name":   "Smith, John",
        "email":  "smithj@outlook.com",
        "origin": "localhost",
        "sub":    "smithj@outlook.com",
        "roles": "anonymous guest"
    }`
	usr, err := user.NewUser(cfg)
	if err != nil {
		panic(err)
	}
	return usr
}

// NewTestGuestAccessList return ACL with guest access.
func NewTestGuestAccessList() *acl.AccessList {
	ctx := context.Background()
	rules := []*acl.RuleConfiguration{
		{
			Comment: "guest access list",
			Conditions: []string{
				"exact match roles anonymous guest",
			},
			Action: `allow`,
		},
	}
	accessList := acl.NewAccessList()
	if err := accessList.AddRules(ctx, rules); err != nil {
		panic(err)
	}
	return accessList
}

// NewTestGuestAccessListWithLogger return ACL with guest access and logger.
func NewTestGuestAccessListWithLogger() *acl.AccessList {
	ctx := context.Background()
	logger := utils.NewLogger()
	rules := []*acl.RuleConfiguration{
		{
			Comment: "guest access list",
			Conditions: []string{
				"exact match roles anonymous guest",
			},
			Action: `allow log`,
		},
	}
	accessList := acl.NewAccessList()
	accessList.SetLogger(logger)
	if err := accessList.AddRules(ctx, rules); err != nil {
		panic(err)
	}
	return accessList
}

// NewTestKeyManagers returns an instance of key manager
func NewTestKeyManagers(method string, secret interface{}) []*kms.KeyManager {
	tokenConfig, err := kms.NewTokenConfig(method, secret)
	if err != nil {
		panic(err)
	}
	keyManager, err := kms.NewKeyManager(tokenConfig)
	if err != nil {
		panic(err)
	}
	return []*kms.KeyManager{keyManager}
}

// NewTestKeyManager returns an instance of key manager.
func NewTestKeyManager(cfg string) *kms.KeyManager {
	tokenConfig, err := kms.NewTokenConfig(cfg)
	if err != nil {
		panic(err)
	}
	keyManager, err := kms.NewKeyManager(tokenConfig)
	if err != nil {
		panic(err)
	}
	return keyManager
}

// GetSharedKey returns shared key for HS algorithms.
func GetSharedKey() string {
	return "8b53b66e-7071-4f7c-ab9a-3ec9dd891704"
}

// GetCookie returns http cookie.
func GetCookie(name, value string, ttl int) *http.Cookie {
	return &http.Cookie{
		Name:    name,
		Value:   value,
		Expires: time.Now().Add(30 * time.Duration(ttl)),
	}
}

// NewTestSigningKey returns signing key.
func NewTestSigningKey() *kms.Key {
	method := "HS512"
	secret := GetSharedKey()
	kms := NewTestKeyManagers(method, secret)
	_, keys := kms[0].GetKeys()
	for _, k := range keys {
		return k
	}
	return nil
}
