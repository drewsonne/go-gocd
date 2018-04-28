package gocd

import (
	"context"
	"fmt"
	"github.com/stretchr/testify/assert"
	"regexp"
	"testing"
	"net/http"
	"io/ioutil"
)

func TestRole(t *testing.T) {
	t.Run("Create/GoCD", testRoleCreateGoCD)
	t.Run("List", testRoleList)
}

func testRoleCreateGoCD(t *testing.T) {
	if runIntegrationTest(t) {

		ctx := context.Background()

		roles := []*Role{
			{
				Name: "spacetiger",
				Type: "gocd",
				Attributes: &RoleAttributesGoCD{
					Users: []string{"alice", "bob", "robin"},
				},
			},
			{
				Name: "my-mock-gocd-role",
				Type: "gocd",
				Attributes: &RoleAttributesGoCD{
					Users: []string{"user-one", "user-two"},
				},
			},
			// Currently there's no fixtures to test the plugin roles,
			// so until there is a way, we can not test plugin role types.
			//{
			//	Name: "blackbird",
			//	Type: "plugin",
			//	Attributes: &RoleAttributesGoCD{
			//		AuthConfigId: String("ldap"),
			//		Properties: []*RoleAttributeProperties{
			//			{
			//				Key:   "UserGroupMembershipAttribute",
			//				Value: "memberOf",
			//			},
			//			{
			//				Key:   "GroupIdentifiers",
			//				Value: "ou=admins,ou=groups,ou=system,dc=example,dc=com",
			//			},
			//		},
			//	},
			//},
		}

		// Test role creation
		for _, role := range roles {
			roleResponse, _, err := intClient.Roles.Create(ctx, role)
			assert.NoError(t, err)

			assert.Regexp(t, regexp.MustCompile("^[a-f0-9]{32}--gzip$"), roleResponse.Version)
			role.Version = roleResponse.Version
			role.Links = roleResponse.Links

			assert.Equal(t, role, roleResponse)
		}

		// Test role listing
		rolesResponses, _, err := intClient.Roles.List(ctx)
		assert.NoError(t, err)

		for i, roleResponse := range rolesResponses {
			assert.Regexp(t, regexp.MustCompile("^[a-f0-9]{32}--gzip$"), roles[i].Version)
			roleResponse.Version = roles[i].Version

			roles[i].Links = roleResponse.Links
			assert.Equal(t, roles[i], roleResponse)
		}

		// Test role update
		roles[0].Attributes.Users = []string{"new-admin"}
		roleUpdateResponse, _, err := intClient.Roles.Update(ctx, roles[0].Name, roles[0])
		assert.NoError(t, err)
		updatedRole, _, err := intClient.Roles.Get(ctx, roleUpdateResponse.Name)
		assert.NoError(t, err)
		assert.Regexp(t, regexp.MustCompile("^[a-f0-9]{32}--gzip$"), updatedRole.Version)
		roles[0].Version = updatedRole.Version
		roles[0].Links = updatedRole.Links
		assert.Equal(t, updatedRole, roles[0])

		// Test role delete
		for _, role := range roles {
			result, _, err := intClient.Roles.Delete(ctx, role.Name)
			assert.Equal(t, fmt.Sprintf("The role '%s' was deleted successfully.", role.Name), result)
			assert.NoError(t, err)
		}
		roleResponse, _, err := intClient.Roles.List(ctx)
		assert.NoError(t, err)
		assert.Empty(t, roleResponse)

	}
}

func testRoleList(t *testing.T) {
	setup()
	defer teardown()
	mux.HandleFunc("/api/admin/security/roles", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, r.Method, "GET", "Unexpected HTTP method")
		assert.Equal(t, apiV1, r.Header.Get("Accept"))

		j, _ := ioutil.ReadFile("test/resources/role.2.json")

		fmt.Fprint(w, string(j))
	})

	r, _, err := client.Roles.List(context.Background())

	assert.NoError(t, err)

	assert.Equal(t, []*Role{
		{
			Name: "spacetiger",
			Type: "gocd",
			Attributes: &RoleAttributesGoCD{
				Users: []string{"alice", "bob", "robin"},
			},
		},
		{
			Name: "blackbird",
			Type: "plugin",
			Attributes: &RoleAttributesGoCD{
				AuthConfigID: String("ldap"),
				Properties: []*RoleAttributeProperties{
					{
						Key:   "UserGroupMembershipAttribute",
						Value: "memberOf",
					},
					{
						Key:   "GroupIdentifiers",
						Value: "ou=admins,ou=groups,ou=system,dc=example,dc=com",
					},
				},
			},
		},
	}, r)
}
