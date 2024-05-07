// Code generated by internal/generate/tagresource/main.go; DO NOT EDIT.

package dynamodb

import (
	"context"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-aws/internal/conns"
	"github.com/hashicorp/terraform-provider-aws/internal/errs/sdkdiag"
	tftags "github.com/hashicorp/terraform-provider-aws/internal/tags"
	"github.com/hashicorp/terraform-provider-aws/internal/tfresource"
	"github.com/hashicorp/terraform-provider-aws/names"
)

// @SDKResource("aws_dynamodb_tag", name="DynamoDB Resource Tag")
func resourceTag() *schema.Resource {
	return &schema.Resource{
		CreateWithoutTimeout: resourceTagCreate,
		ReadWithoutTimeout:   resourceTagRead,
		UpdateWithoutTimeout: resourceTagUpdate,
		DeleteWithoutTimeout: resourceTagDelete,

		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"resource_arn": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			names.AttrKey: {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			names.AttrValue: {
				Type:     schema.TypeString,
				Required: true,
			},
		},
	}
}

func resourceTagCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics { // nosemgrep:ci.semgrep.tags.calling-UpdateTags-in-resource-create
	var diags diag.Diagnostics
	conn := meta.(*conns.AWSClient).DynamoDBClient(ctx)

	identifier := d.Get("resource_arn").(string)
	key := d.Get(names.AttrKey).(string)
	value := d.Get(names.AttrValue).(string)

	if err := updateTagsV2(ctx, conn, identifier, nil, map[string]string{key: value}); err != nil {
		return sdkdiag.AppendErrorf(diags, "creating %s resource (%s) tag (%s): %s", names.DynamoDB, identifier, key, err)
	}

	d.SetId(tftags.SetResourceID(identifier, key))

	return append(diags, resourceTagRead(ctx, d, meta)...)
}

func resourceTagRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	conn := meta.(*conns.AWSClient).DynamoDBClient(ctx)

	identifier, key, err := tftags.GetResourceID(d.Id())
	if err != nil {
		return sdkdiag.AppendFromErr(diags, err)
	}

	value, err := GetTag(ctx, conn, identifier, key)

	if !d.IsNewResource() && tfresource.NotFound(err) {
		log.Printf("[WARN] %s resource (%s) tag (%s) not found, removing from state", names.DynamoDB, identifier, key)
		d.SetId("")
		return diags
	}

	if err != nil {
		return sdkdiag.AppendErrorf(diags, "reading %s resource (%s) tag (%s): %s", names.DynamoDB, identifier, key, err)
	}

	d.Set("resource_arn", identifier)
	d.Set(names.AttrKey, key)
	d.Set(names.AttrValue, value)

	return diags
}

func resourceTagUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	conn := meta.(*conns.AWSClient).DynamoDBClient(ctx)

	identifier, key, err := tftags.GetResourceID(d.Id())
	if err != nil {
		return sdkdiag.AppendFromErr(diags, err)
	}

	if err := updateTagsV2(ctx, conn, identifier, nil, map[string]string{key: d.Get(names.AttrValue).(string)}); err != nil {
		return sdkdiag.AppendErrorf(diags, "updating %s resource (%s) tag (%s): %s", names.DynamoDB, identifier, key, err)
	}

	return append(diags, resourceTagRead(ctx, d, meta)...)
}

func resourceTagDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	conn := meta.(*conns.AWSClient).DynamoDBClient(ctx)

	identifier, key, err := tftags.GetResourceID(d.Id())
	if err != nil {
		return sdkdiag.AppendFromErr(diags, err)
	}

	if err := updateTagsV2(ctx, conn, identifier, map[string]string{key: d.Get(names.AttrValue).(string)}, nil); err != nil {
		return sdkdiag.AppendErrorf(diags, "deleting %s resource (%s) tag (%s): %s", names.DynamoDB, identifier, key, err)
	}

	return diags
}
