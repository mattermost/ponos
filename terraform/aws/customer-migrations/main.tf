module "migration-bucket" {
  source                 = "git::https://github.com/mattermost/mattermost-cloud-monitoring.git//aws/migrate-to-cloud?ref=v1.1.0"
  region                 = var.region
  account_id             = var.account_id
  bucket                 = var.bucket
  customer_bucket_folder = var.customer_bucket_folder
  customer_policy_name   = var.customer_bucket_folder
  kms_key                = var.kms_key
}
