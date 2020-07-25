package bcs

import "C"
import (
	"fmt"
	rand2 "github.com/fabric-app/pkg/util/rand"
	clientMSP "github.com/hyperledger/fabric-sdk-go/pkg/client/msp"
	"github.com/hyperledger/fabric-sdk-go/pkg/core/cryptosuite"
	"github.com/hyperledger/fabric-sdk-go/pkg/msp"
	"log"
	"os"
)

func (c *Client) NewCAClient() {
	log.Print("Enroll registrar")
	ctxProvider := c.SDK.Context()
	mspClient, err := clientMSP.New(ctxProvider)
	//registrarEnrollID, registrarEnrollSecret := c.getRegistrarEnrollmentCredentials(ctxProvider)

	err = mspClient.Enroll("admin", clientMSP.WithSecret("1b48f5aae5d425142058fd2412e815251d26dd3e1175ee8e555f4dc4ed56b6fa"))
	if err != nil {
		log.Fatalf("enroll registrar failed: %v", err)
	}
	c.MC = mspClient
}

func (c *Client) removeUserData() {
	configBackend, err := c.SDK.Config()
	if err != nil {
		log.Fatal(err)
	}

	cryptoSuiteConfig := cryptosuite.ConfigFromBackend(configBackend)
	identityConfig, err := msp.ConfigFromBackend(configBackend)
	if err != nil {
		log.Fatal(err)
	}

	keyStorePath := cryptoSuiteConfig.KeyStorePath()
	credentialStorePath := identityConfig.CredentialStorePath()
	c.removePath(keyStorePath)
	c.removePath(credentialStorePath)
}

func (c *Client) removePath(storePath string) {
	err := os.RemoveAll(storePath)
	if err != nil {
		log.Fatalf("Cleaning up directory '%s' failed: %v", storePath, err)
	}
}

//func (c *Client) getRegistrarEnrollmentCredentials(ctxProvider context.ClientProvider) (string, string) {
//
//	ctx, err := ctxProvider()
//	if err != nil {
//		fmt.Printf("failed to get context: %v\n", err)
//	}
//
//	clientConfig := ctx.IdentityConfig().Client()
//	//if err != nil {
//	//	fmt.Printf("config.Client() failed: %v\n", err)
//	//}
//
//	myOrg := clientConfig.Organization
//
//	caConfig, ok := ctx.IdentityConfig().CAConfig(myOrg)
//	if ok {
//		fmt.Printf("CAConfig failed: %v\n", err)
//	}
//
//	return caConfig.Registrar.EnrollID, caConfig.Registrar.EnrollSecret
//}

// Register a new user
func (c *Client) RegisterUser(username, orgName, secret, identityTypeUser string) (string, bool) {
	// Register the new user
	log.Printf("User not found, registering new user: %v", username)
	testAttributes := []clientMSP.Attribute{
		{
			Name:  rand2.RandStringBytesMaskImprSrcUnsafe(10),
			Value: fmt.Sprintf("%s:ecert", rand2.RandStringBytesMaskImprSrcUnsafe(10)),
			ECert: true,
		},
		{
			Name:  rand2.RandStringBytesMaskImprSrcUnsafe(10),
			Value: fmt.Sprintf("%s:ecert", rand2.RandStringBytesMaskImprSrcUnsafe(10)),
			ECert: true,
		},
	}
	_, err := c.MC.Register(&clientMSP.RegistrationRequest{
		Name:        username,
		Type:        identityTypeUser,
		Attributes:  testAttributes,
		Affiliation: orgName,
		Secret:      secret, // Is ready to get hash?
	})
	if err != nil {
		return username + err.Error(), false
	}
	//signingIdentity, err := mspClient.GetSigningIdentity(username)
	//log.Printf("%s: %s", signingIdentity.Identifier().ID, string(signingIdentity.EnrollmentCertificate()[:]))
	return username + " register Successfully", true

}

// Enroll a user
func (c *Client) EnrollUser(username, orgName, secret, identityTypeUser string) (string, bool) {
	//ctxProvider := c.SDK.Context(fabsdk.WithOrg(orgName))
	//mspClient, err := clientMSP.New(ctxProvider)
	//if err != nil {
	//	log.Fatalf("Failed to create msp client: %s", err.Error())
	//	return username + " login error :" + err.Error(), false
	//}
	err := c.MC.Enroll(username, clientMSP.WithSecret(secret))
	if err != nil {
		log.Printf("enroll %s failed: %v", username, err)
		return err.Error(), false
	}
	//	signingIdentity, err := c.MC.GetSigningIdentity(username)
	//	log.Printf("%s: %s", signingIdentity.Identifier().ID, string(signingIdentity.EnrollmentCertificate()[:]))
	return username + " login success", true
}

// Reovoke a user
func (c *Client) RevokeUser(username, orgName, secret, identityTypeUser string) (string, bool) {
	request := clientMSP.RemoveIdentityRequest{
		ID:     username,
		Force:  true,
		CAName: "ca.org1.lzawt.com",
	}
	idr, err := c.MC.RemoveIdentity(&request)
	if err != nil {
		log.Printf("enroll %s failed: %v", username, err)
		return err.Error(), false
	}
	log.Println(idr)
	return username + " RevokeUser success", true
}

// getRegisteredUser get registered user. If user is not enrolled, enroll new user
func (c *Client) GetRegisteredUser(username, orgName, secret, identityTypeUser string) (string, bool) {
	//ctxProvider := c.SDK.Context(fabsdk.WithOrg(orgName))
	//mspClient, err := clientMSP.New(ctxProvider)
	//if err != nil {
	//	log.Fatalf("Failed to create msp client: %s", err.Error())
	//}

	signingIdentity, err := c.MC.GetSigningIdentity(username)
	if err != nil {
		log.Printf("Check if user %s is enrolled: %s", username, err.Error())
		testAttributes := []clientMSP.Attribute{
			{
				Name:  rand2.RandStringBytesMaskImprSrcUnsafe(10),
				Value: fmt.Sprintf("%s:ecert", rand2.RandStringBytesMaskImprSrcUnsafe(10)),
				ECert: true,
			},
			{
				Name:  rand2.RandStringBytesMaskImprSrcUnsafe(10),
				Value: fmt.Sprintf("%s:ecert", rand2.RandStringBytesMaskImprSrcUnsafe(10)),
				ECert: true,
			},
		}

		// Register the new user
		identity, err := c.MC.GetIdentity(username)
		if true {
			log.Printf("User %s does not exist, registering new user", username)
			_, err = c.MC.Register(&clientMSP.RegistrationRequest{
				Name:        username,
				Type:        identityTypeUser,
				Attributes:  testAttributes,
				Affiliation: orgName,
				Secret:      secret,
			})
		} else {
			log.Printf("Identity: %s", identity.Secret)
		}
		//enroll user
		err = c.MC.Enroll(username, clientMSP.WithSecret(secret))
		if err != nil {
			log.Printf("enroll %s failed: %v", username, err)
			return "failed " + err.Error(), false
		}

		return username + " enrolled Successfully", true
	}
	log.Printf("%s: %s", signingIdentity.Identifier().ID, string(signingIdentity.EnrollmentCertificate()[:]))
	return username + " already enrolled", true
}

func (c *Client) GetAllUsers() ([]string, error) {
	var users []string
	ids, err := c.MC.GetAllIdentities()
	if err != nil {
		return nil, err
	}
	for _, v := range ids {
		users = append(users, v.ID)
	}
	return users, nil
}

func (c *Client) ModifyUserSecret(user, passwd string) ([]byte, error) {
	res, err := c.MC.ModifyIdentity(&clientMSP.IdentityRequest{ID: user, Affiliation: "org2", Secret: passwd})
	if err != nil {
		return nil, err
	}
	return []byte(res.ID), nil
}
