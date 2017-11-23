/* 
hypercloud_test.go

A "comprehensive" set of tests for the hypercloud go library.

These tests require the following environment variables to be set:

    - HC_BASE_URL
    - HC_ACCESS_KEY
    - HC_SECRET_KEY
    - HC_SSH_PUB

for Authentication

Have a duck:
       ..---.. 
     .'  _    `. 
 __..'  (o)    : 
`..__          ; 
     `.       / 
       ;      `..---...___ 
     .'                   `~-. .-') 
    .                         ' _.' 
   :                           : 
   \                           ' 
    +                         J 
     `._                   _.' 
        `~--....___...---~' mh 

*/

package hypercloud

import (

    "fmt"
    "os"
    "strings"
    "testing"

    hypercloud "bitbucket.org/mistarhee/hypercloud" //Replace with the official repo before making a PR
)

func TestHypercloud(t *testing.T){
    base_url := os.Getenv("HC_BASE_URL")
    access_key := os.Getenv("HC_ACCESS_KEY")
    secret_key := os.Getenv("HC_SECRET_KEY")

    hc := hypercloud.NewHypercloud(base_url, access_key, secret_key)

    /* Lets start by grabbing all the data we need to create an instance */

    //Get all regions, find the Sydney (SY3) region.
    regions, errs := hc.RegionList()
    if errs != nil {
        t.Fatalf("Error occurred in getting RegionList: \n%v", errs)
    }

    var mRegion string
    mRegion = nil //uninit for "safety"
    for _, region := range regions {
        //Find SY3
        reg := region.(map[string]interface{})
        if reg["code"].(string) == "SY3"{
            mRegion = reg["id"].(string)
            break
        }
    }

    if mRegion == nil {
        t.Fatalf("Failed to get the region id for SY3")
    }

    // Lets grab the Standard performance tier for disks and instances in the SY3 region
    var mInstanceTier string
    var mDiskTier string

    instanceTiers, err := hc.PerformanceTierListInstances()
    if err != nil {
        t.Fatalf("Error occurred in getting PerformanceTierListInstances: \n%v", errs)
    }

    for _, its := range instanceTiers.([]interface{}) {
        if its.(map[string]interface{})["region"].(map[string]interface{})["id"].(string) == mRegion && its.(map[string]interface{})["name"] == "Standard" {
            mInstanceTier = its.(map[string]interface{})["id"].(string)
            break
        }
    }

    if mInstanceTier == nil {
        t.Fatalf("Failed to get the standard instance tier id for SY3")
    }

    diskTiers, err := hc.PerformanceTierListDisks()
    if err != nil {
        t.Fatalf("Error occurred in getting PerformanceTierListDisks: \n%v", errs)
    }

    for _, dts := range diskTiers.([]interface{}) {
        if dts.(map[string]interface{})["region"].(map[string]interface{})["id"].(string) = mRegion && dts.(map[string]interface{})["name"] == "Standard" {
            mDiskTier = dts.(map[string]interface{})["id"].string()
            break
        }
    }

    if mDiskTier == nil {
        t.Fatalf("Failed to get the standard disk tier id for SY3")
    }

    // Make a blank 10G disk of specified performance tier
    diskMap = make(map[string]interface{})
    diskMap["name"] = "hypercloud-test-disk"
    diskMap["performance_tier"] = mDiskTier
    diskMap["region"] = mRegion
    diskMap["size"] = 10
    newDisk, err = hc.DiskCreate(diskMap)
    if err != nil {
        t.Fatalf("Failed to create new disk: \n%v", errs)
    }

    defer hc.DiskDelete(newDisk.(map[string]interface{})["id"].(string))


    //Actually, lets resize it to say 20 G
    diskMap = make(map[string]interface{})
    diskMap["size"] = 20
    updatedDisk, err := hc.DiskResize(newDisk["id"].(string), diskMap)
    if err != nil {
        t.Fatalf("Failed to resize the new disk: \n%v", err)
    }


    //Lets make a generic new instance in SY3
    instanceMap = make(map[string]interface{})
    instanceMap["name"] = "hypercloud-test-instance"
    instanceMap["performance_tier"] = mInstanceTier
    instanceMap["region"] = mRegion
    instanceMap["memory"] = 2048

    newInstance, err := hc.InstanceCreate(instanceMap)
    if err != nil {
        t.Fatalf("Failed to create the new instance: \n%v", errs)
    }

    defer hc.InstanceDelete(newInstance.(map[string]interface{})["id"].(string))

    //Now we need a boot disk for this instance
    //Search all templates for a Ubuntu 16.10
    var mTemplateId string
    templates, err := hc.TemplateList()
    if err != nil {
        t.Fatalf("Failed to list all templates: \n%v", errs)
    }

    for _, t := range templates.([]interface{}) {
        if t.(map[string]interface{})["name"] == "Ubuntu 16.10" && t.(map[string]interface{})["region"].(map[string]interface{})["id"] == mRegion {
            mTemplateId = t.(map[string]interface{})["id"]
            break
        }
    }

    if mTemplateId == nil {
        t.Fatalf("Failed to get template id for Ubuntu 16.10 in SY3")
    }

    diskMap = make(map[string]interface{})
    diskMap["name"] = "hypercloud-test-boot-disk"

    bootDisk, err := hc.DiskClone(mTemplateId, diskMap)
    if err != nil {
        t.Fatalf("Unable to clone the template: \n%v", err)
    }

    defer hc.DiskDelete(bootDisk.(map[string]interface{})["id"])

    //Now lets make a public/private IP for this guy

    //Public IP
    netMap = make(map[string]interface{})
    netMap["region"] = mRegion
    mPubIp, err = hc.IPAddressCreate(netMap)
    if err != nil {
        t.Fatalf("Unable to allocate new public IP in SY3: \n%v", err)
    }

    defer hc.IPAddressDelete(mPubIp.(map[string]interface{})["id"])

    //Private IP
    //Make a network adapter for this test
    netMap = make(map[string]interface{})
    netMap["name"] = "hypercloud-test-network-adapter"
    netMap["region"] = mRegion
    netMap["specification"] = "10.6.9.0/24" //gonna delete it anyway

    mNetAdapter, err = hc.NetworkCreate(netMap)
    if err != nil {
        t.Fatalf("Unable to create private network: \n%v", err)
    }

    defer hc.NetworkDelete(mNetAdapter.(map[string]interface{})["id"])

    netMap = make(map[string]interface{})
}

