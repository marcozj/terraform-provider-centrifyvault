package centrify

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	logger "github.com/marcozj/golang-sdk/logging"
	vault "github.com/marcozj/golang-sdk/platform"
	"github.com/marcozj/golang-sdk/restapi"
)

func dataSourceOidcWebApp() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceOidcWebAppRead,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Name of the Web App",
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Description of the Web App",
			},
			"application_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Application ID. Specify the name or 'target' that the mobile application uses to find this application.",
			},
			"oauth_profile": getOidcProfileSchema(),
		},
	}
}

func dataSourceOidcWebAppRead(d *schema.ResourceData, m interface{}) error {
	logger.Infof("Finding Oidc webapp")
	client := m.(*restapi.RestClient)
	object := vault.NewOidcWebApp(client)
	object.Name = d.Get("name").(string)
	object.ApplicationID = d.Get("application_id").(string)

	// We can't use simple Query method because it doesn't return all attributes
	err := object.GetByName()
	if err != nil {
		return fmt.Errorf("error retrieving Oauth webapp with name '%s': %s", object.Name, err)
	}
	d.SetId(object.ID)

	schemamap, err := vault.GenerateSchemaMap(object)
	if err != nil {
		return err
	}
	//logger.Debugf("Generated Map for resourceOidcWebAppRead(): %+v", schemamap)
	for k, v := range schemamap {
		switch k {
		case "oauth_profile":
			d.Set(k, []interface{}{v})
		case "challenge_rule":
			d.Set(k, v.(map[string]interface{})["rule"])
		case "workflow_settings":
			if v.(string) != "" {
				wfschema, err := convertWorkflowSchema(v.(string))
				if err != nil {
					return err
				}
				d.Set("workflow_approver", wfschema)
				d.Set(k, v)
			}
		default:
			d.Set(k, v)
		}
	}

	return nil
}