package sesv2_test

import (
	"context"
	"errors"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/hashicorp/terraform-provider-aws/internal/acctest"
	"github.com/hashicorp/terraform-provider-aws/internal/conns"
	"github.com/hashicorp/terraform-provider-aws/internal/create"
	tfsesv2 "github.com/hashicorp/terraform-provider-aws/internal/service/sesv2"
	"github.com/hashicorp/terraform-provider-aws/names"
)

func TestAccSESV2EmailIdentityFeedbackAttributes_basic(t *testing.T) {
	rName := acctest.RandomEmailAddress(acctest.RandomDomainName())
	resourceName := "aws_sesv2_email_identity_feedback_attributes.test"
	emailIdentityName := "aws_sesv2_email_identity.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ErrorCheck:               acctest.ErrorCheck(t, names.SESV2EndpointID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckEmailIdentityDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccEmailIdentityFeedbackAttributesConfig_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckEmailIdentityFeedbackAttributesExist(emailIdentityName, false),
					resource.TestCheckResourceAttrPair(resourceName, "email_identity", emailIdentityName, "email_identity"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccSESV2EmailIdentityFeedbackAttributes_disappears(t *testing.T) {
	rName := acctest.RandomEmailAddress(acctest.RandomDomainName())
	resourceName := "aws_sesv2_email_identity_feedback_attributes.test"
	emailIdentityName := "aws_sesv2_email_identity.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ErrorCheck:               acctest.ErrorCheck(t, names.SESV2EndpointID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckEmailIdentityDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccEmailIdentityFeedbackAttributesConfig_emailForwardingEnabled(rName, true),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckEmailIdentityFeedbackAttributesExist(emailIdentityName, true),
					acctest.CheckResourceDisappears(acctest.Provider, tfsesv2.ResourceEmailIdentityFeedbackAttributes(), resourceName),
				),
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func TestAccSESV2EmailIdentityFeedbackAttributes_disappears_emailIdentity(t *testing.T) {
	rName := acctest.RandomEmailAddress(acctest.RandomDomainName())
	emailIdentityName := "aws_sesv2_email_identity.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ErrorCheck:               acctest.ErrorCheck(t, names.SESV2EndpointID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckEmailIdentityDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccEmailIdentityFeedbackAttributesConfig_emailForwardingEnabled(rName, true),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckEmailIdentityFeedbackAttributesExist(emailIdentityName, true),
					acctest.CheckResourceDisappears(acctest.Provider, tfsesv2.ResourceEmailIdentity(), emailIdentityName),
				),
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func TestAccSESV2EmailIdentityFeedbackAttributes_emailForwardingEnabled(t *testing.T) {
	rName := acctest.RandomEmailAddress(acctest.RandomDomainName())
	resourceName := "aws_sesv2_email_identity_feedback_attributes.test"
	emailIdentityName := "aws_sesv2_email_identity.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ErrorCheck:               acctest.ErrorCheck(t, names.SESV2EndpointID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckEmailIdentityDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccEmailIdentityFeedbackAttributesConfig_emailForwardingEnabled(rName, true),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckEmailIdentityFeedbackAttributesExist(emailIdentityName, true),
					resource.TestCheckResourceAttr(resourceName, "email_forwarding_enabled", "true"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccEmailIdentityFeedbackAttributesConfig_emailForwardingEnabled(rName, false),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckEmailIdentityFeedbackAttributesExist(emailIdentityName, false),
					resource.TestCheckResourceAttr(resourceName, "email_forwarding_enabled", "false"),
				),
			},
		},
	})
}

// testAccCheckEmailIdentityFeedbackAttributesExist verifies that both the email identity exists,
// and that the email forwarding enabled setting is correct
func testAccCheckEmailIdentityFeedbackAttributesExist(name string, emailForwardingEnabled bool) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[name]
		if !ok {
			return create.Error(names.SESV2, create.ErrActionCheckingExistence, tfsesv2.ResNameEmailIdentity, name, errors.New("not found"))
		}

		if rs.Primary.ID == "" {
			return create.Error(names.SESV2, create.ErrActionCheckingExistence, tfsesv2.ResNameEmailIdentity, name, errors.New("not set"))
		}

		conn := acctest.Provider.Meta().(*conns.AWSClient).SESV2Client

		out, err := tfsesv2.FindEmailIdentityByID(context.Background(), conn, rs.Primary.ID)
		if err != nil {
			return create.Error(names.SESV2, create.ErrActionCheckingExistence, tfsesv2.ResNameEmailIdentity, rs.Primary.ID, err)
		}
		if out == nil || out.FeedbackForwardingStatus != emailForwardingEnabled {
			return create.Error(names.SESV2, create.ErrActionCheckingExistence, tfsesv2.ResNameEmailIdentityFeedbackAttributes, rs.Primary.ID, errors.New("feedback attributes not set"))
		}

		return nil
	}
}

func testAccEmailIdentityFeedbackAttributesConfig_basic(rName string) string {
	return fmt.Sprintf(`
resource "aws_sesv2_email_identity" "test" {
  email_identity = %[1]q
}

resource "aws_sesv2_email_identity_feedback_attributes" "test" {
  email_identity = aws_sesv2_email_identity.test.email_identity
}
`, rName)
}

func testAccEmailIdentityFeedbackAttributesConfig_emailForwardingEnabled(rName string, emailForwardingEnabled bool) string {
	return fmt.Sprintf(`
resource "aws_sesv2_email_identity" "test" {
  email_identity = %[1]q
}

resource "aws_sesv2_email_identity_feedback_attributes" "test" {
  email_identity           = aws_sesv2_email_identity.test.email_identity
  email_forwarding_enabled = %[2]t
}
`, rName, emailForwardingEnabled)
}
