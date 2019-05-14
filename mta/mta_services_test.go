package mta

import (
	"encoding/json"
	"errors"
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"os"
	"path/filepath"
	"reflect"

	ghodss "github.com/ghodss/yaml"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/gofrs/flock"
)

var _ = Describe("MtaServices", func() {

	schemaVersion := "1.1"
	oMtaInput := MTA{
		ID:            "test",
		Version:       "1.2",
		SchemaVersion: &schemaVersion,
		Description:   "test mta creation",
	}

	AfterEach(func() {
		err := os.RemoveAll(getTestPath("result"))
		Ω(err).Should(Succeed())
	})
	var _ = Describe("CreateMta", func() {
		It("Create MTA", func() {
			jsonData, err := json.Marshal(oMtaInput)
			Ω(err).Should(Succeed())
			mtaPath := getTestPath("result", "temp.mta.yaml")
			Ω(CreateMta(mtaPath, string(jsonData), os.MkdirAll)).Should(Succeed())
			Ω(mtaPath).Should(BeAnExistingFile())
			yamlData, err := ioutil.ReadFile(mtaPath)
			Ω(err).Should(Succeed())
			oMtaOutput, err := Unmarshal(yamlData)
			Ω(err).Should(Succeed())
			Ω(reflect.DeepEqual(oMtaInput, *oMtaOutput)).Should(BeTrue())
		})

		It("Create MTA with wrong json format", func() {
			wrongJSON := "{Name:fff"
			mtaPath := getTestPath("result", "temp.mta.yaml")
			Ω(CreateMta(mtaPath, wrongJSON, os.MkdirAll)).Should(HaveOccurred())
		})

		It("Create MTA fail to create file", func() {
			jsonData, err := json.Marshal(oMtaInput)
			Ω(err).Should(Succeed())
			mtaPath := getTestPath("result", "temp.mta.yaml")
			Ω(CreateMta(mtaPath, string(jsonData), mkDirsErr)).Should(HaveOccurred())
		})
	})

	var _ = Describe("CopyFile", func() {
		It("Copy file content", func() {
			jsonData, err := json.Marshal(oMtaInput)
			Ω(err).Should(Succeed())
			sourceFilePath := getTestPath("result", "temp.mta.yaml")
			targetFilePath := getTestPath("result", "temp2.mta.yaml")
			Ω(CreateMta(sourceFilePath, string(jsonData), os.MkdirAll)).Should(Succeed())
			Ω(CopyFile(sourceFilePath, targetFilePath, os.Create)).Should(Succeed())
			Ω(targetFilePath).Should(BeAnExistingFile())
			yamlData, err := ioutil.ReadFile(targetFilePath)
			Ω(err).Should(Succeed())
			oOutput, err := Unmarshal(yamlData)
			Ω(err).Should(Succeed())
			Ω(reflect.DeepEqual(oMtaInput, *oOutput)).Should(BeTrue())
		})

		It("Copy file with non existing path", func() {
			sourceFilePath := "c:/temp/test1"
			targetFilePath := "c:/temp/test2"
			Ω(CopyFile(sourceFilePath, targetFilePath, os.Create)).Should(HaveOccurred())
		})

		It("Copy file fail to create destination file", func() {
			jsonData, err := json.Marshal(oMtaInput)
			Ω(err).Should(Succeed())
			sourceFilePath := getTestPath("result", "temp.mta.yaml")
			targetFilePath := getTestPath("result", "temp2.mta.yaml")
			Ω(CreateMta(sourceFilePath, string(jsonData), os.MkdirAll)).Should(Succeed())
			Ω(CopyFile(sourceFilePath, targetFilePath, createErr)).Should(HaveOccurred())
			Ω(targetFilePath).ShouldNot(BeAnExistingFile())
		})
	})

	var _ = Describe("deleteFile", func() {
		It("Delete file", func() {
			jsonData, err := json.Marshal(oMtaInput)
			Ω(err).Should(Succeed())
			mtaPath := getTestPath("result", "temp.mta.yaml")
			Ω(CreateMta(mtaPath, string(jsonData), os.MkdirAll)).Should(Succeed())
			Ω(mtaPath).Should(BeAnExistingFile())
			Ω(DeleteFile(mtaPath)).Should(Succeed())
			Ω(mtaPath).ShouldNot(BeAnExistingFile())
		})
	})

	var _ = Describe("addModule", func() {
		It("Add module", func() {
			oModule := Module{
				Name: "testModule",
				Type: "testType",
				Path: "test",
			}

			mtaPath := getTestPath("result", "temp.mta.yaml")

			jsonRootData, err := json.Marshal(oMtaInput)
			Ω(err).Should(Succeed())
			Ω(CreateMta(mtaPath, string(jsonRootData), os.MkdirAll)).Should(Succeed())

			jsonModuleData, err := json.Marshal(oModule)
			Ω(err).Should(Succeed())
			Ω(AddModule(mtaPath, string(jsonModuleData), ghodss.Marshal)).Should(Succeed())

			oMtaInput.Modules = append(oMtaInput.Modules, &oModule)
			Ω(mtaPath).Should(BeAnExistingFile())
			yamlData, err := ioutil.ReadFile(mtaPath)
			Ω(err).Should(Succeed())
			oMtaOutput, err := Unmarshal(yamlData)
			Ω(err).Should(Succeed())
			Ω(reflect.DeepEqual(oMtaInput, *oMtaOutput)).Should(BeTrue())
		})

		It("Add module to non existing mta.yaml file", func() {
			json := "{name:fff}"
			mtaPath := getTestPath("result", "mta.yaml")
			Ω(AddModule(mtaPath, json, ghodss.Marshal)).Should(HaveOccurred())
		})

		It("Add module to wrong mta.yaml format", func() {
			wrongJSON := "{TEST:fff}"
			oModule := Module{
				Name: "testModule",
				Type: "testType",
				Path: "test",
			}

			mtaPath := getTestPath("result", "mta.yaml")
			Ω(CreateMta(mtaPath, wrongJSON, os.MkdirAll)).Should(Succeed())

			jsonModuleData, err := json.Marshal(oModule)
			Ω(err).Should(Succeed())
			Ω(AddModule(mtaPath, string(jsonModuleData), ghodss.Marshal)).Should(HaveOccurred())
		})

		It("Add module with wrong json format", func() {
			wrongJSON := "{name:fff"

			mtaPath := getTestPath("result", "temp.mta.yaml")
			jsonRootData, err := json.Marshal(oMtaInput)
			Ω(err).Should(Succeed())
			Ω(CreateMta(mtaPath, string(jsonRootData), os.MkdirAll)).Should(Succeed())

			Ω(AddModule(mtaPath, wrongJSON, ghodss.Marshal)).Should(HaveOccurred())
		})

		It("Add module fails to marshal", func() {
			oModule := Module{
				Name: "testModule",
				Type: "testType",
				Path: "test",
			}

			mtaPath := getTestPath("result", "temp.mta.yaml")

			jsonRootData, err := json.Marshal(oMtaInput)
			Ω(err).Should(Succeed())
			Ω(CreateMta(mtaPath, string(jsonRootData), os.MkdirAll)).Should(Succeed())

			jsonModuleData, err := json.Marshal(oModule)
			Ω(err).Should(Succeed())
			Ω(AddModule(mtaPath, string(jsonModuleData), marshalErr)).Should(HaveOccurred())
		})
	})

	var _ = Describe("addResource", func() {
		It("Add resource", func() {
			oResource := Resource{
				Name: "testResource",
				Type: "testType",
			}

			mtaPath := getTestPath("result", "temp.mta.yaml")

			jsonRootData, err := json.Marshal(oMtaInput)
			Ω(err).Should(Succeed())
			Ω(CreateMta(mtaPath, string(jsonRootData), os.MkdirAll)).Should(Succeed())

			jsonResourceData, err := json.Marshal(oResource)
			Ω(err).Should(Succeed())
			Ω(AddResource(mtaPath, string(jsonResourceData), ghodss.Marshal)).Should(Succeed())

			oMtaInput.Resources = append(oMtaInput.Resources, &oResource)
			Ω(mtaPath).Should(BeAnExistingFile())
			yamlData, err := ioutil.ReadFile(mtaPath)
			Ω(err).Should(Succeed())
			oMtaOutput, err := Unmarshal(yamlData)
			Ω(err).Should(Succeed())
			Ω(reflect.DeepEqual(oMtaInput, *oMtaOutput)).Should(BeTrue())
		})

		It("Add resource to non existing mta.yaml file", func() {
			json := "{name:fff}"
			mtaPath := getTestPath("result", "mta.yaml")
			Ω(AddResource(mtaPath, json, ghodss.Marshal)).Should(HaveOccurred())
		})

		It("Add resource to wrong mta.yaml format", func() {
			wrongJSON := "{TEST:fff}"
			oResource := Resource{
				Name: "testResource",
				Type: "testType",
			}

			mtaPath := getTestPath("result", "mta.yaml")
			Ω(CreateMta(mtaPath, wrongJSON, os.MkdirAll)).Should(Succeed())

			jsonResourceData, err := json.Marshal(oResource)
			Ω(err).Should(Succeed())
			Ω(AddResource(mtaPath, string(jsonResourceData), ghodss.Marshal)).Should(HaveOccurred())
		})

		It("Add resource with wrong json format", func() {
			wrongJSON := "{name:fff"

			mtaPath := getTestPath("result", "temp.mta.yaml")
			jsonRootData, err := json.Marshal(oMtaInput)
			Ω(err).Should(Succeed())
			Ω(CreateMta(mtaPath, string(jsonRootData), os.MkdirAll)).Should(Succeed())

			Ω(AddResource(mtaPath, wrongJSON, ghodss.Marshal)).Should(HaveOccurred())
		})

		It("Add resource fails to marshal", func() {
			oResource := Resource{
				Name: "testResource",
				Type: "testType",
			}

			mtaPath := getTestPath("result", "temp.mta.yaml")

			jsonRootData, err := json.Marshal(oMtaInput)
			Ω(err).Should(Succeed())
			Ω(CreateMta(mtaPath, string(jsonRootData), os.MkdirAll)).Should(Succeed())

			jsonResourceData, err := json.Marshal(oResource)
			Ω(err).Should(Succeed())
			Ω(AddResource(mtaPath, string(jsonResourceData), marshalErr)).Should(HaveOccurred())
		})
	})
})

var _ = Describe("Module", func() {
	AfterEach(func() {
		err := os.RemoveAll(getTestPath("result"))
		Ω(err).Should(Succeed())
	})

	It("Sanity", func() {
		os.MkdirAll(getTestPath("result"), os.ModePerm)
		mtaPath := getTestPath("result", "mta.yaml")
		Ω(CopyFile(getTestPath("mta.yaml"), mtaPath, os.Create)).Should(Succeed())

		var err error
		mtaHashCode, err := GetMtaHash(mtaPath)
		Ω(err).Should(Succeed())
		oModule := Module{
			Name: "testModule",
			Type: "testType",
			Path: "test",
		}

		jsonData, err := json.Marshal(oModule)
		moduleJSON := string(jsonData)
		err = ModifyMta(mtaPath, func() error {
			return AddModule(mtaPath, moduleJSON, yaml.Marshal)
		}, mtaHashCode)
		Ω(err).Should(Succeed())
		oModule.Name = "test1"
		jsonData, err = json.Marshal(oModule)
		moduleJSON = string(jsonData)
		// hashcode of the mta.yaml is wrong now
		err = ModifyMta(mtaPath, func() error {
			return AddModule(mtaPath, moduleJSON, ghodss.Marshal)
		}, mtaHashCode)
		Ω(err).Should(HaveOccurred())
	})
	It("Locking fails", func() {
		os.MkdirAll(getTestPath("result"), os.ModePerm)
		mtaPath := getTestPath("result", "mta.yaml")
		lockFilePath := filepath.Join(filepath.Dir(mtaPath), "mta-lock.lock")
		lock := flock.New(lockFilePath)
		Ω(lock.TryLock()).Should(BeTrue())
		Ω(CopyFile(getTestPath("mta.yaml"), mtaPath, os.Create)).Should(Succeed())

		var err error
		mtaHashCode, err := GetMtaHash(mtaPath)
		Ω(err).Should(Succeed())
		oModule := Module{
			Name: "testModule",
			Type: "testType",
			Path: "test",
		}

		jsonData, err := json.Marshal(oModule)
		moduleJSON := string(jsonData)
		err = ModifyMta(mtaPath, func() error {
			return AddModule(mtaPath, moduleJSON, yaml.Marshal)
		}, mtaHashCode)
		Ω(err).Should(HaveOccurred())
		Ω(err.Error()).Should(Equal(`failed to lock the "C:\Users\i019379\go\src\github.com\SAP\cloud-mta\mta\testdata\result\mta.yaml" file for modification`))
		err = lock.Close()
		Ω(err).Should(Succeed())
		err = os.Remove(lockFilePath)
		//Ω(err).Should(Succeed())
		err = ModifyMta(mtaPath, func() error {
			return AddModule(mtaPath, moduleJSON, yaml.Marshal)
		}, mtaHashCode)
		Ω(err).Should(Succeed())
	})
})

func mkDirsErr(path string, perm os.FileMode) error {
	return errors.New("err")
}

func createErr(path string) (*os.File, error) {
	return nil, errors.New("err")
}

func marshalErr(o interface{}) ([]byte, error) {
	return nil, errors.New("err")
}
