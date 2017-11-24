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
        t.Logf("Error occurred in getting RegionList: \n%v", errs)
        t.FailNow()
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
        t.Logf("Failed to get the region id for SY3")
        t.FailNow()
    }

    // Lets grab the Standard performance tier for disks and instances in the SY3 region
    var mInstanceTier string
    var mDiskTier string

    instanceTiers, err := hc.PerformanceTierListInstances()
    if err != nil {
        t.Logf("Error occurred in getting PerformanceTierListInstances: \n%v", errs)
        t.FailNow()
    }

    for _, its := range instanceTiers.([]interface{}) {
        if its.(map[string]interface{})["region"].(map[string]interface{})["id"].(string) == mRegion && its.(map[string]interface{})["name"].(string) == "Standard" {
            mInstanceTier = its.(map[string]interface{})["id"].(string)
            break
        }
    }

    if mInstanceTier == nil {
        t.Logf("Failed to get the standard instance tier id for SY3")
        t.FailNow()
    }

    diskTiers, err := hc.PerformanceTierListDisks()
    if err != nil {
        t.Logf("Error occurred in getting PerformanceTierListDisks: \n%v", errs)
        t.FailNow()
    }

    for _, dts := range diskTiers.([]interface{}) {
        if dts.(map[string]interface{})["region"].(map[string]interface{})["id"].(string) = mRegion && dts.(map[string]interface{})["name"].(string) == "Standard" {
            mDiskTier = dts.(map[string]interface{})["id"].(string)
            break
        }
    }

    if mDiskTier == nil {
        t.Logf("Failed to get the standard disk tier id for SY3")
        t.FailNow()
    }

    // Make a blank 10G disk of specified performance tier
    var mDisk string
    diskMap = make(map[string]interface{})
    diskMap["name"] = "hypercloud-test-disk"
    diskMap["performance_tier"] = mDiskTier
    diskMap["region"] = mRegion
    diskMap["size"] = 10
    newDisk, err := hc.DiskCreate(diskMap)
    if err != nil {
        t.Logf("Failed to create new disk: \n%v", errs)
        t.FailNow()
    }

    mDisk = newDisk.(map[string]interface{})["id"].(string)
    defer hc.DiskDelete(mDisk)

    //Actually, lets resize it to say 20 G
    diskMap = make(map[string]interface{})
    diskMap["size"] = 20
    updatedDisk, err := hc.DiskResize(mDisk, diskMap)
    if err != nil {
        t.Logf("Failed to resize the new disk: \n%v", err)
        t.FailNow()
    }

    //Now we need a boot disk for this instance
    //Search all templates for a Ubuntu 16.10
    var mTemplateId string
    templates, err := hc.TemplateList()
    if err != nil {
        t.Logf("Failed to list all templates: \n%v", errs)
        t.FailNow()
    }

    for _, t := range templates.([]interface{}) {
        if t.(map[string]interface{})["name"].(string) == "Ubuntu 16.10" && t.(map[string]interface{})["region"].(map[string]interface{})["id"].(string) == mRegion {
            mTemplateId = t.(map[string]interface{})["id"].(string)
            break
        }
    }

    if mTemplateId == nil {
        t.Logf("Failed to get template id for Ubuntu 16.10 in SY3")
        t.FailNow()
    }

    diskMap = make(map[string]interface{})
    diskMap["name"] = "hypercloud-test-boot-disk"

    var mBootDisk string
    bootDisk, err := hc.DiskClone(mTemplateId, diskMap)
    if err != nil {
        t.Logf("Unable to clone the template: \n%v", err)
        t.FailNow()
    }

    defer hc.DiskDelete(bootDisk.(map[string]interface{})["id"].(string))

    //Now lets make a public/private IP for this guy

    //Public IP
    var mPubIp string
    netMap = make(map[string]interface{})
    netMap["region"] = mRegion
    pubIp, err := hc.IPAddressCreate(netMap)
    if err != nil {
        t.Logf("Unable to allocate new public IP in SY3: \n%v", err)
        t.FailNow()
    }

    mPubIp = pubIp.(map[string]interface{})["id"].(string)
    defer hc.IPAddressDelete(mPubIp)

    //Private IP
    //Make a network adapter for this test
    var mPrivIp string
    var mNetAdapter string
    netMap = make(map[string]interface{})
    netMap["name"] = "hypercloud-test-network-adapter"
    netMap["region"] = mRegion
    netMap["specification"] = "10.6.9.0/24" //gonna delete it anyway

    netAdapter, err := hc.NetworkCreate(netMap)
    if err != nil {
        t.Logf("Unable to create private network: \n%v", err)
        t.FailNow()
    }
    mNetAdapter = netAdapter.(map[string]interface{}).["id"].(string)
    defer hc.NetworkDelete(mNetAdapter)

    // Make a private IP
    netMap = make(map[string]interface{})
    netMap["name"] = "hypercloud-test-private-ip"
    netMap["network"] = mNetAdapter.(map[string]interface{})["id"].(string)

    privIp := hc.IPAddressCreate(netMap)
    if err != nil {
        t.Logf("Unable to create private ip: \n%v", err)
        t.FailNow()
    }
    mPrivIp = privIp.(map[string]interface{})["id"].(string)
    defer hc.IPAddressDelete(mPrivIp)

    //Lets make a generic new instance in SY3
    var mInstance string

    instanceMap := make(map[string]interface{})
    instanceMap["name"] = "hypercloud-test-instance"
    instanceMap["performance_tier"] = mInstanceTier
    instanceMap["region"] = mRegion
    instanceMap["memory"] = 2048

    newInstance, err := hc.InstanceCreate(instanceMap)
    if err != nil {
        t.Logf("Failed to create the new instance: \n%v", errs)
        t.FailNow()
    }

    mInstance = newInstance.(map[string]interface{})["id"].(string)

    defer hc.InstanceDelete(mInstance)

    //Attach disks/IP addresses to the guy
    updateInstance = make(map[string]interface{})

    updateInstance["disks"] = []string{mBootDisk, mDisk}
    updateInstance["network_adapters"] = []interface{}
    // Add the private IP
    privateNetwork = make(map[string]interface{})
    privateNetwork["network"] = privIp.(map[string]interface{})["network_id"].(string)
    privateNetwork["ip_addresses"] = []string{mPrivIp}

    updateInstance["network_adapters"] = append(updateInstance["network_adapters"], privateNetwork)

    //Add the public IP
    publicNetwork = make(map[string]interface{})
    publicNetwork["network"] = pubIp.(map[string]interface{})["network_id"].(string)
    publicNetwork["ip_addresses"] = []string{mPubIp}


    //Call the update function
    updateInstance["network_adapters"] = append(updateInstance["network_adapters"], publicNetwork)

    newInstance, err = hc.InstanceUpdate(mInstance, updateInstance)
    if err != nil {
        t.Logf("Unable to update instance: \n%v", err)
        t.FailNow()
    }
}

