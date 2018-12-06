package arm

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

const mysqlTemplate = `{
    "$schema": "http://schema.management.azure.com/schemas/2014-04-01-preview/deploymentTemplate.json#",
    "contentVersion": "1.0.0.0",
    "parameters": {
        "administratorLogin": {
            "type": "string"
        },
        "administratorLoginPassword": {
            "type": "securestring"
        },
        "location": {
            "type": "string"
        },
        "serverName": {
            "type": "string"
        },
        "version": {
            "type": "string"
        },
        "sslEnforcement": {
            "type": "string"
        },
		"tags": {
			"type": "object"
		}
    },
    "resources": [
        {
            "apiVersion": "2017-12-01-preview",
            "kind": "",
            "location": "[parameters('location')]",
            "name": "[parameters('serverName')]",
			"tags": "[parameters('tags')]",
            "properties": {
                "version": "[parameters('version')]",
                "administratorLogin": "[parameters('administratorLogin')]",
                "administratorLoginPassword": "[parameters('administratorLoginPassword')]",
				"sslEnforcement": "[parameters('sslEnforcement')]",
                "storageProfile": {
                    "storageMB": "102400",
                    "backupRetentionDays": 7,
                    "geoRedundantBackup": "Disabled"
                }
            },
            "sku": {
                "name": "GP_Gen5_4",
                "tier": "GeneralPurpose",
                "capacity": 4,
                "size": 102400,
                "family": "Gen5"
            },
            "type": "Microsoft.DBforMySQL/servers"
        }
    ],
	"outputs": {
		"MYSQL_HOST": {
			"type": "string",
			"value": "[reference(parameters('serverName')).fullyQualifiedDomainName]"
		}
	}
}`

func TestMySQLTemplateExists(t *testing.T) {
	tpl, err := GetTemplate("mysql")
	require.NoError(t, err)
	assert.NotEmpty(t, tpl)
	assert.Equal(t, mysqlTemplate, string(tpl))
}
