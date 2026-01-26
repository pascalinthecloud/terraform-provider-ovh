package ovh

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccCloudProjectRegionStorage_basic(t *testing.T) {
	bucketName := acctest.RandomWithPrefix(test_prefix)
	serviceName := os.Getenv("OVH_CLOUD_PROJECT_SERVICE_TEST")
	config := fmt.Sprintf(`
	resource "ovh_cloud_project_storage" "storage" {
		service_name = "%s"
		region_name = "GRA"
		name = "%s"
		versioning = {
			status = "enabled"
		}
	}
	`, serviceName, bucketName)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheckCloud(t); testAccCheckCloudProjectExists(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("ovh_cloud_project_storage.storage", "name", bucketName),
					resource.TestCheckResourceAttr("ovh_cloud_project_storage.storage", "region", "GRA"),
					resource.TestCheckResourceAttr("ovh_cloud_project_storage.storage", "versioning.status", "enabled"),
					resource.TestCheckResourceAttr("ovh_cloud_project_storage.storage", "encryption.sse_algorithm", "plaintext"),
					resource.TestCheckResourceAttrSet("ovh_cloud_project_storage.storage", "virtual_host"),
					// Verify ID is populated with composite format: service_name/region_name/name
					resource.TestCheckResourceAttr("ovh_cloud_project_storage.storage", "id", fmt.Sprintf("%s/GRA/%s", serviceName, bucketName)),
				),
			},
			{
				ImportState:                          true,
				ImportStateVerify:                    true,
				ImportStateVerifyIdentifierAttribute: "name",
				ResourceName:                         "ovh_cloud_project_storage.storage",
				ImportStateId:                        fmt.Sprintf("%s/GRA/%s", os.Getenv("OVH_CLOUD_PROJECT_SERVICE_TEST"), bucketName),
				ImportStateVerifyIgnore:              []string{"created_at"}, // Ignore created_at since its value is invalid in response of the POST.
			},
		},
	})
}

func TestAccCloudProjectRegionStorage_withReplication(t *testing.T) {
	bucketName := acctest.RandomWithPrefix(test_prefix)
	replicaBucketName := acctest.RandomWithPrefix(test_prefix)
	serviceName := os.Getenv("OVH_CLOUD_PROJECT_SERVICE_TEST")

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheckCloud(t); testAccCheckCloudProjectExists(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
					resource "ovh_cloud_project_storage" "storage" {
						service_name = "%s"
						region_name = "GRA"
						name = "%s"

						versioning = {
							status = "enabled"
						}

						replication = {
							rules = [
								{
									id          = "test"
									priority    = 1
									status      = "enabled"
									destination = {
										name   = "%s"
										region = "GRA"
										remove_on_main_bucket_deletion = true
									}
									filter = {
										"prefix" = "test"
										"tags"   = {
											"key": "test"
										}
									}
									delete_marker_replication = "disabled"
								}
							]
						}
					}`, serviceName, bucketName, replicaBucketName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("ovh_cloud_project_storage.storage", "name", bucketName),
					resource.TestCheckResourceAttr("ovh_cloud_project_storage.storage", "region", "GRA"),
					resource.TestCheckResourceAttr("ovh_cloud_project_storage.storage", "versioning.status", "enabled"),
					resource.TestCheckResourceAttr("ovh_cloud_project_storage.storage", "encryption.sse_algorithm", "plaintext"),
					resource.TestCheckResourceAttr("ovh_cloud_project_storage.storage", "replication.rules.#", "1"),
					resource.TestCheckResourceAttr("ovh_cloud_project_storage.storage", "replication.rules.0.id", "test"),
					resource.TestCheckResourceAttr("ovh_cloud_project_storage.storage", "replication.rules.0.priority", "1"),
					resource.TestCheckResourceAttr("ovh_cloud_project_storage.storage", "replication.rules.0.status", "enabled"),
					resource.TestCheckResourceAttr("ovh_cloud_project_storage.storage", "replication.rules.0.destination.name", replicaBucketName),
					resource.TestCheckResourceAttr("ovh_cloud_project_storage.storage", "replication.rules.0.destination.region", "GRA"),
					resource.TestCheckResourceAttr("ovh_cloud_project_storage.storage", "replication.rules.0.filter.prefix", "test"),
					resource.TestCheckResourceAttrSet("ovh_cloud_project_storage.storage", "virtual_host"),
					// Verify ID is populated with composite format: service_name/region_name/name
					resource.TestCheckResourceAttr("ovh_cloud_project_storage.storage", "id", fmt.Sprintf("%s/GRA/%s", serviceName, bucketName)),
				),
			},
			{
				Config: fmt.Sprintf(`
					resource "ovh_cloud_project_storage" "storage" {
						service_name = "%s"
						region_name = "GRA"
						name = "%s"

						versioning = {
							status = "enabled"
						}

						replication = {
							rules = [
								{
									id          = "test"
									priority    = 1
									status      = "enabled"
									destination = {
										name   = "%s"
										region = "GRA"
										remove_on_main_bucket_deletion = true
									}
									filter = {
										"prefix" = "test-updated"
										"tags"   = {
											"key": "test-updated"
										}
									}
									delete_marker_replication = "disabled"
								}
							]
						}
					} `, serviceName, bucketName, replicaBucketName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("ovh_cloud_project_storage.storage", "name", bucketName),
					resource.TestCheckResourceAttr("ovh_cloud_project_storage.storage", "region", "GRA"),
					resource.TestCheckResourceAttr("ovh_cloud_project_storage.storage", "versioning.status", "enabled"),
					resource.TestCheckResourceAttr("ovh_cloud_project_storage.storage", "encryption.sse_algorithm", "plaintext"),
					resource.TestCheckResourceAttr("ovh_cloud_project_storage.storage", "replication.rules.#", "1"),
					resource.TestCheckResourceAttr("ovh_cloud_project_storage.storage", "replication.rules.0.id", "test"),
					resource.TestCheckResourceAttr("ovh_cloud_project_storage.storage", "replication.rules.0.priority", "1"),
					resource.TestCheckResourceAttr("ovh_cloud_project_storage.storage", "replication.rules.0.status", "enabled"),
					resource.TestCheckResourceAttr("ovh_cloud_project_storage.storage", "replication.rules.0.destination.name", replicaBucketName),
					resource.TestCheckResourceAttr("ovh_cloud_project_storage.storage", "replication.rules.0.destination.region", "GRA"),
					resource.TestCheckResourceAttr("ovh_cloud_project_storage.storage", "replication.rules.0.filter.prefix", "test-updated"),
					resource.TestCheckResourceAttrSet("ovh_cloud_project_storage.storage", "virtual_host"),
					// Verify ID is populated with composite format after update
					resource.TestCheckResourceAttr("ovh_cloud_project_storage.storage", "id", fmt.Sprintf("%s/GRA/%s", serviceName, bucketName)),
				),
			},
			{
				ImportState:                          true,
				ImportStateVerify:                    true,
				ImportStateVerifyIdentifierAttribute: "name",
				ResourceName:                         "ovh_cloud_project_storage.storage",
				ImportStateId:                        fmt.Sprintf("%s/GRA/%s", os.Getenv("OVH_CLOUD_PROJECT_SERVICE_TEST"), bucketName),
				// Ignore created_at since its value is invalid in response of the POST.
				// Also ignore remove_on_main_bucket_deletion since its computed value is not returned by the API.
				ImportStateVerifyIgnore: []string{"created_at", "replication.rules.0.destination.remove_on_main_bucket_deletion"},
			},
		},
	})
}

func TestAccCloudProjectRegionStorage_withObjectLock(t *testing.T) {
	bucketName := acctest.RandomWithPrefix(test_prefix)
	serviceName := os.Getenv("OVH_CLOUD_PROJECT_SERVICE_TEST")

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheckCloud(t); testAccCheckCloudProjectExists(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Case 1: P7D (7 days) -> Expect API to return P1W, provider should handle it.
			{
				Config: fmt.Sprintf(`
					resource "ovh_cloud_project_storage" "storage" {
						service_name = "%s"
						region_name = "GRA"
						name = "%s"
						object_lock = {
							status = "enabled"
							rule = {
								mode = "compliance"
								period = "P7D"
							}
						}
					}`, serviceName, bucketName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("ovh_cloud_project_storage.storage", "name", bucketName),
					resource.TestCheckResourceAttr("ovh_cloud_project_storage.storage", "object_lock.status", "enabled"),
					resource.TestCheckResourceAttr("ovh_cloud_project_storage.storage", "object_lock.rule.mode", "compliance"),
					resource.TestCheckResourceAttr("ovh_cloud_project_storage.storage", "object_lock.rule.period", "P7D"),
				),
			},
			// Case 2: Update to P2W (2 weeks) -> Standard weeks format.
			{
				Config: fmt.Sprintf(`
					resource "ovh_cloud_project_storage" "storage" {
						service_name = "%s"
						region_name = "GRA"
						name = "%s"
						object_lock = {
							status = "enabled"
							rule = {
								mode = "compliance"
								period = "P2W"
							}
						}
					}`, serviceName, bucketName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("ovh_cloud_project_storage.storage", "object_lock.rule.period", "P2W"),
				),
			},
			// Case 3: Update to P14D (14 days) -> Expect API to return P2W, provider should handle it.
			{
				Config: fmt.Sprintf(`
					resource "ovh_cloud_project_storage" "storage" {
						service_name = "%s"
						region_name = "GRA"
						name = "%s"
						object_lock = {
							status = "enabled"
							rule = {
								mode = "compliance"
								period = "P14D"
							}
						}
					}`, serviceName, bucketName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("ovh_cloud_project_storage.storage", "object_lock.rule.period", "P14D"),
				),
			},
			// Case 4: Update to P1D (1 Day) -> Smallest unit.
			{
				Config: fmt.Sprintf(`
					resource "ovh_cloud_project_storage" "storage" {
						service_name = "%s"
						region_name = "GRA"
						name = "%s"
						object_lock = {
							status = "enabled"
							rule = {
								mode = "compliance"
								period = "P1D"
							}
						}
					}`, serviceName, bucketName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("ovh_cloud_project_storage.storage", "object_lock.rule.period", "P1D"),
				),
			},
			// Case 5: Update to P8D (8 Days) -> Not divisible by 7. Expect API to keep P8D or handle accordingly.
			{
				Config: fmt.Sprintf(`
					resource "ovh_cloud_project_storage" "storage" {
						service_name = "%s"
						region_name = "GRA"
						name = "%s"
						object_lock = {
							status = "enabled"
							rule = {
								mode = "compliance"
								period = "P8D"
							}
						}
					}`, serviceName, bucketName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("ovh_cloud_project_storage.storage", "object_lock.rule.period", "P8D"),
				),
			},
			// Case 6: Update to P364D (364 Days) -> 52 Weeks. Expect API to likely return P52W.
			{
				Config: fmt.Sprintf(`
					resource "ovh_cloud_project_storage" "storage" {
						service_name = "%s"
						region_name = "GRA"
						name = "%s"
						object_lock = {
							status = "enabled"
							rule = {
								mode = "compliance"
								period = "P364D"
							}
						}
					}`, serviceName, bucketName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("ovh_cloud_project_storage.storage", "object_lock.rule.period", "P364D"),
				),
			},
			// Case 7: P1Y (1 Year) -> Test year-based duration
			{
				Config: fmt.Sprintf(`
					resource "ovh_cloud_project_storage" "storage" {
						service_name = "%s"
						region_name = "GRA"
						name = "%s"
						object_lock = {
							status = "enabled"
							rule = {
								mode = "compliance"
								period = "P1Y"
							}
						}
					}`, serviceName, bucketName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("ovh_cloud_project_storage.storage", "object_lock.rule.period", "P1Y"),
				),
			}},
	})
}
