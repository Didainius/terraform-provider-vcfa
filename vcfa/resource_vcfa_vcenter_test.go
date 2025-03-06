//go:build tm || ALL || functional

package vcfa

import (
	"fmt"
	"regexp"
	"sync"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

var doOnceFirstTests sync.Once

// type firstTests struct {
// 	// runOnce   sync.Once
// 	// setupOnce sync.Once
// 	t     *testing.T
// 	tests []testEntry
// }

// type testEntry struct {
// 	Name     string
// 	testFunc func(*testing.T)
// }

var priorityTestCleanupFunc func() error

func testAccPriority(t *testing.T) {
	doOnceFirstTests.Do(func() {
		firstTestAcc(t)
	})
}

func firstTestAcc(t *testing.T) {
	tests := []func(*testing.T){
		TestAccVcfaNsxManager,
		TestAccVcfaVcenter,
		TestAccVcfaVcenterInvalid,
	}

	for index, test := range tests {
		t.Run(fmt.Sprintf("%d ", index), test)
	}

	// doOnceTestAccVcfaVcenter.Do(func() {
	// 	t.Run("TestAccVcfaVcenter", testAccVcfaVcenter)
	// })

	// setup shared things for other tests

	if vcfaTestVerbose {
		fmt.Println("# Will setup vCenter and NSX Manager")
	}
	cleanup, err := setupVcAndNsx()
	if err != nil {
		fmt.Printf("error setting up shared VC and NSX: %s", err)
	}

	priorityTestCleanupFunc = cleanup

	// priorityTestCleanupFunc = func(t *testing.T) {
	// 	err = cleanupFunc()
	// 	if err != nil {
	// 		t.Errorf("error cleaning up: %s", err)
	// 	}
	// }

	if vcfaTestVerbose {
		fmt.Println("# Done")
	}

}

var doOnceTestAccVcfaVcenter sync.Once

func TestAccVcfaVcenter(t *testing.T) {
	doOnceTestAccVcfaVcenter.Do(func() {
		t.Run("TestAccVcfaVcenter", testAccVcfaVcenter)
	})
}

func testAccVcfaVcenter(t *testing.T) {
	preTestChecks(t)
	skipIfNotSysAdmin(t)

	if !testConfig.Tm.CreateVcenter {
		t.Skipf("Skipping vCenter creation")
	}

	var params = StringMap{
		"Org":             testConfig.Tm.Org,
		"VcenterUsername": testConfig.Tm.VcenterUsername,
		"VcenterPassword": testConfig.Tm.VcenterPassword,
		"VcenterUrl":      testConfig.Tm.VcenterUrl,
		"NsxUsername":     testConfig.Tm.NsxManagerUsername,
		"NsxPassword":     testConfig.Tm.NsxManagerPassword,
		"NsxUrl":          testConfig.Tm.NsxManagerUrl,

		"Testname": t.Name(),

		"Tags": "tm",
	}
	testParamsNotEmpty(t, params)

	configText1 := templateFill(testAccVcfaVcenterStep1, params)
	params["FuncName"] = t.Name() + "-step2"
	configText2 := templateFill(testAccVcfaVcenterStep2, params)

	params["FuncName"] = t.Name() + "-step3"
	configText3 := templateFill(testAccVcfaVcenterStep3, params)

	params["FuncName"] = t.Name() + "-step4"
	configText4 := templateFill(testAccVcfaVcenterStep4DS, params)

	debugPrintf("#[DEBUG] CONFIGURATION step1: %s\n", configText1)
	debugPrintf("#[DEBUG] CONFIGURATION step2: %s\n", configText2)
	debugPrintf("#[DEBUG] CONFIGURATION step3: %s\n", configText3)
	debugPrintf("#[DEBUG] CONFIGURATION step4: %s\n", configText4)
	if vcfaShortTest {
		t.Skip(acceptanceTestsSkipped)
		return
	}

	resource.Test(t, resource.TestCase{
		ProviderFactories: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: configText1,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("vcfa_vcenter.test", "id"),
					resource.TestCheckResourceAttr("vcfa_vcenter.test", "name", t.Name()),
					resource.TestCheckResourceAttr("vcfa_vcenter.test", "description", ""),
					resource.TestCheckResourceAttr("vcfa_vcenter.test", "is_enabled", "true"),
					resource.TestCheckResourceAttr("vcfa_vcenter.test", "has_proxy", "false"),
					resource.TestCheckResourceAttrSet("vcfa_vcenter.test", "cluster_health_status"),
					resource.TestCheckResourceAttrSet("vcfa_vcenter.test", "is_connected"),
					resource.TestCheckResourceAttrSet("vcfa_vcenter.test", "connection_status"),
					resource.TestCheckResourceAttrSet("vcfa_vcenter.test", "mode"),
					resource.TestCheckResourceAttrSet("vcfa_vcenter.test", "uuid"),
					resource.TestCheckResourceAttrSet("vcfa_vcenter.test", "vcenter_version"),
				),
			},
			{
				Config: configText2,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("vcfa_vcenter.test", "id"),
					resource.TestCheckResourceAttr("vcfa_vcenter.test", "name", t.Name()+"-rename"),
					resource.TestCheckResourceAttr("vcfa_vcenter.test", "is_enabled", "false"),
					resource.TestCheckResourceAttr("vcfa_vcenter.test", "description", "description from Terraform"),
					resource.TestCheckResourceAttr("vcfa_vcenter.test", "has_proxy", "false"),
					resource.TestCheckResourceAttrSet("vcfa_vcenter.test", "cluster_health_status"),
					resource.TestCheckResourceAttrSet("vcfa_vcenter.test", "is_connected"),
					resource.TestCheckResourceAttrSet("vcfa_vcenter.test", "connection_status"),
					resource.TestCheckResourceAttrSet("vcfa_vcenter.test", "mode"),
					resource.TestCheckResourceAttrSet("vcfa_vcenter.test", "uuid"),
					resource.TestCheckResourceAttrSet("vcfa_vcenter.test", "vcenter_version"),
					resource.TestCheckResourceAttr("vcfa_vcenter.test", "status", "READY"),
				),
			},
			{
				Config: configText3,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("vcfa_vcenter.test", "id"),
					resource.TestCheckResourceAttr("vcfa_vcenter.test", "name", t.Name()),
					resource.TestCheckResourceAttr("vcfa_vcenter.test", "description", ""),
					resource.TestCheckResourceAttr("vcfa_vcenter.test", "is_enabled", "true"),
					resource.TestCheckResourceAttr("vcfa_vcenter.test", "has_proxy", "false"),
					resource.TestCheckResourceAttrSet("vcfa_vcenter.test", "cluster_health_status"),
					resource.TestCheckResourceAttrSet("vcfa_vcenter.test", "is_connected"),
					resource.TestCheckResourceAttrSet("vcfa_vcenter.test", "connection_status"),
					resource.TestCheckResourceAttrSet("vcfa_vcenter.test", "mode"),
					resource.TestCheckResourceAttrSet("vcfa_vcenter.test", "uuid"),
					resource.TestCheckResourceAttrSet("vcfa_vcenter.test", "vcenter_version"),
				),
			},
			{
				ResourceName:            "vcfa_vcenter.test",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateId:           params["Testname"].(string),
				ImportStateVerifyIgnore: []string{"password", "auto_trust_certificate", "refresh_vcenter_on_read", "refresh_policies_on_read", "refresh_vcenter_on_create", "refresh_policies_on_create", "nsx_manager_id"},
			},
			{
				Config: configText4,
				Check: resource.ComposeTestCheckFunc(
					resourceFieldsEqual("data.vcfa_vcenter.test", "vcfa_vcenter.test", []string{"%", "nsx_manager_id"}),
				),
			},
		},
	})

	postTestChecks(t)
}

const testAccVcfaVcenterPrerequisites = `
resource "vcfa_nsx_manager" "test" {
  name                   = "{{.Testname}}"
  description            = "terraform test"
  username               = "{{.NsxUsername}}"
  password               = "{{.NsxPassword}}"
  url                    = "{{.NsxUrl}}"
  auto_trust_certificate = true
}
`

const testAccVcfaVcenterStep1 = testAccVcfaVcenterPrerequisites + `
resource "vcfa_vcenter" "test" {
  name                     = "{{.Testname}}"
  url                      = "{{.VcenterUrl}}"
  auto_trust_certificate   = true
  refresh_vcenter_on_read  = true
  refresh_policies_on_read = true
  username                 = "{{.VcenterUsername}}"
  password                 = "{{.VcenterPassword}}"
  is_enabled               = true
  nsx_manager_id           = vcfa_nsx_manager.test.id
}
`

const testAccVcfaVcenterStep2 = testAccVcfaVcenterPrerequisites + `
resource "vcfa_vcenter" "test" {
  name                   = "{{.Testname}}-rename"
  description            = "description from Terraform"
  auto_trust_certificate = true
  url                    = "{{.VcenterUrl}}"
  username               = "{{.VcenterUsername}}"
  password               = "{{.VcenterPassword}}"
  is_enabled             = false
  nsx_manager_id         = vcfa_nsx_manager.test.id
}
`

const testAccVcfaVcenterStep3 = testAccVcfaVcenterPrerequisites + `
resource "vcfa_vcenter" "test" {
  name                   = "{{.Testname}}"
  url                    = "{{.VcenterUrl}}"
  auto_trust_certificate = true
  username               = "{{.VcenterUsername}}"
  password               = "{{.VcenterPassword}}"
  is_enabled             = true
  nsx_manager_id         = vcfa_nsx_manager.test.id
}
`

const testAccVcfaVcenterStep4DS = testAccVcfaVcenterStep3 + `
data "vcfa_vcenter" "test" {
  name = vcfa_vcenter.test.name
}
`

var doOnceTestAccVcfaVcenterInvalid sync.Once

func TestAccVcfaVcenterInvalid(t *testing.T) {
	doOnceTestAccVcfaVcenter.Do(func() {
		t.Run("TestAccVcfaVcenterInvalid", testAccVcfaVcenterInvalid)
	})
}

func testAccVcfaVcenterInvalid(t *testing.T) {
	preTestChecks(t)
	skipIfNotSysAdmin(t)

	// test fails on purpose
	if vcfaShortTest {
		t.Skip(acceptanceTestsSkipped)
		return
	}

	var params = StringMap{
		"Org":             testConfig.Tm.Org,
		"VcenterUsername": testConfig.Tm.VcenterUsername,
		"VcenterPassword": "invalid",
		"VcenterUrl":      testConfig.Tm.VcenterUrl,
		"NsxUsername":     testConfig.Tm.NsxManagerUsername,
		"NsxPassword":     testConfig.Tm.NsxManagerPassword,
		"NsxUrl":          testConfig.Tm.NsxManagerUrl,

		"Testname": t.Name(),

		"Tags": "tm",
	}
	testParamsNotEmpty(t, params)

	configText1 := templateFill(testAccVcfaVcenterStep1, params)
	params["FuncName"] = t.Name() + "-step2"
	params["VcenterPassword"] = testConfig.Tm.VcenterPassword
	configText2 := templateFill(testAccVcfaVcenterStep1, params)

	debugPrintf("#[DEBUG] CONFIGURATION step1: %s\n", configText1)
	debugPrintf("#[DEBUG] CONFIGURATION step2: %s\n", configText2)

	resource.Test(t, resource.TestCase{
		ProviderFactories: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config:      configText1,
				ExpectError: regexp.MustCompile(`Failed to connect to the vCenter`),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("vcfa_vcenter.test", "id"),
				),
			},
			{
				Config: configText2,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("vcfa_vcenter.test", "id"),
					resource.TestCheckResourceAttr("vcfa_vcenter.test", "name", t.Name()),
					resource.TestCheckResourceAttr("vcfa_vcenter.test", "is_enabled", "true"),
					resource.TestCheckResourceAttrSet("vcfa_vcenter.test", "cluster_health_status"),
					resource.TestCheckResourceAttrSet("vcfa_vcenter.test", "is_connected"),
					resource.TestCheckResourceAttrSet("vcfa_vcenter.test", "connection_status"),
					resource.TestCheckResourceAttrSet("vcfa_vcenter.test", "mode"),
					resource.TestCheckResourceAttrSet("vcfa_vcenter.test", "uuid"),
					resource.TestCheckResourceAttrSet("vcfa_vcenter.test", "vcenter_version"),
					resource.TestCheckResourceAttr("vcfa_vcenter.test", "status", "READY"),
				),
			},
		},
	})

	postTestChecks(t)
}
