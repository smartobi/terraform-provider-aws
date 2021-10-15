package backup

import (
	"fmt"
	"log"
	"regexp"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/backup"
	"github.com/hashicorp/aws-sdk-go-base/tfawserr"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/hashicorp/terraform-provider-aws/internal/conns"
)

func ResourceVaultLockConfiguration() *schema.Resource {
	return &schema.Resource{
		Create: resourceVaultLockConfigurationCreate,
		Delete: resourceVaultLockConfigurationDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"backup_vault_name": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringMatch(regexp.MustCompile(`^[a-zA-Z0-9\-\_\.]{1,50}$`), "must consist of lowercase letters, numbers, and hyphens."),
			},
			"changeable_for_days": {
				Type:         schema.TypeInt,
				Optional:     true,
				ValidateFunc: validation.IntAtLeast(3),
			},
			"max_retention_days": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"min_retention_days": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"backup_vault_arn": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceVaultLockConfigurationCreate(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*conns.AWSClient).BackupConn

	input := &backup.PutBackupVaultLockConfigurationInput{
		BackupVaultName: aws.String(d.Get("backup_vault_name").(string)),
	}

	if v, ok := d.GetOk("changeable_for_days"); ok {
		input.changeable_for_days = aws.Int64(int64(v.(int)))
	}

	if v, ok := d.GetOk("max_retention_days"); ok {
		input.max_retention_days = aws.Int64(int64(v.(int)))
	}

	if v, ok := d.GetOk("min_retention_days"); ok {
		input.min_retention_days = aws.Int64(int64(v.(int)))
	}

	_, err := conn.PutBackupVaultLockConfiguration(input)
	if err != nil {
		return fmt.Errorf("error creating Backup Vault Lock Configuration (%s): %w", d.Id(), err)
	}

	d.SetId(d.Get("backup_vault_name").(string))

	// note: BackupVaultLockConfiguration currently does not have a GetBackupVaultLockConfiguration
	// Reference: https://pkg.go.dev/github.com/aws/aws-sdk-go-v2/service/backup
	return nil
}

func resourceVaultLockConfigurationDelete(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*conns.AWSClient).BackupConn

	input := &backup.DeleteBackupVaultLockConfigurationInput{
		BackupVaultName: aws.String(d.Id()),
	}

	_, err := conn.DeleteBackupVaultLockConfiguration(input)
	if err != nil {
		if tfawserr.ErrMessageContains(err, backup.ErrCodeResourceNotFoundException, "") {
			return nil
		}
		return fmt.Errorf("error deleting Backup Vault Lock Configuration (%s): %w", d.Id(), err)
	}

	return nil
}
