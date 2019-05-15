package commands

import (
	"encoding/json"
	"os"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/SAP/cloud-mta/mta"
)

var _ = Describe("Resource", func() {

	AfterEach(func() {
		os.RemoveAll(getTestPath("result"))
	})

	It("Sanity", func() {
		os.MkdirAll(getTestPath("result"), os.ModePerm)
		addResourceMtaCmdPath = getTestPath("result", "mta.yaml")
		Ω(mta.CopyFile(getTestPath("mta.yaml"), addResourceMtaCmdPath, os.Create)).Should(Succeed())

		var err error

		hash, exists, err := mta.GetMtaHash(addResourceMtaCmdPath)
		addResourceCmdHashcode = hash
		Ω(err).Should(Succeed())
		Ω(exists).Should(BeTrue())
		oResource := mta.Resource{
			Name: "testResource",
			Type: "testType",
		}

		jsonData, err := json.Marshal(&oResource)
		addResourceCmdData = string(jsonData)
		Ω(addResourceCmd.RunE(nil, []string{})).Should(Succeed())
		oResource.Name = "test1"
		jsonData, err = json.Marshal(oResource)
		addResourceCmdData = string(jsonData)
		// hashcode of the mta.yaml is wrong now
		Ω(addResourceCmd.RunE(nil, []string{})).Should(HaveOccurred())
	})
})
