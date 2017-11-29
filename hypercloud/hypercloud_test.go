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
    "os"
    "testing"
    "time"
)

func TestHypercloud(t *testing.T){
    base_url := os.Getenv("HC_BASE_URL")
    access_key := os.Getenv("HC_ACCESS_KEY")
    secret_key := os.Getenv("HC_SECRET_KEY")

    hc, err := NewHypercloud(base_url, access_key, secret_key)
    if err != nil {
        t.Logf("Failed to create initial hypercloud object: \n%v", err)
        t.FailNow()
    }

    /* Lets start by grabbing all the data we need to create an instance */

    //Get all regions, find the Sydney (SY3) region.
    regions, errs := hc.RegionList()
    if errs != nil {
        t.Logf("Error occurred in getting RegionList: \n%v", errs)
        t.FailNow()
    }

    var mRegion string
    for _, region := range regions.([]interface{}) {
        //Find SY3
        reg := region.(map[string]interface{})
        if reg["code"].(string) == "SY3"{
            mRegion = reg["id"].(string)
            break
        }
    }

    if mRegion == "" {
        t.Logf("Failed to get the region id for SY3")
        t.FailNow()
    }

    // Lets grab the Standard performance tier for disks and instances in the SY3 region
    var mInstanceTier string
    var mDiskTier string

    instanceTiers, err := hc.PerformanceTierListInstance()
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

    if mInstanceTier == "" {
        t.Logf("Failed to get the standard instance tier id for SY3")
        t.FailNow()
    }

    diskTiers, err := hc.PerformanceTierListDisk()
    if err != nil {
        t.Logf("Error occurred in getting PerformanceTierListDisks: \n%v", errs)
        t.FailNow()
    }

    for _, dts := range diskTiers.([]interface{}) {
        if dts.(map[string]interface{})["region"].(map[string]interface{})["id"].(string) == mRegion && dts.(map[string]interface{})["name"].(string) == "Standard" {
            mDiskTier = dts.(map[string]interface{})["id"].(string)
            break
        }
    }

    if mDiskTier == "" {
        t.Logf("Failed to get the standard disk tier id for SY3")
        t.FailNow()
    }

    // Make a blank 10G disk of specified performance tier
    var mDisk string
    diskMap := make(map[string]interface{})
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
    // Wait for resources to be up 
    end := time.Now().Add(time.Duration(30) * time.Second)
    for end.After(time.Now()){
        diskInfo, err := hc.DiskInfo(mDisk)
        if err != nil {
            t.Logf("Failed to grab disk info: \n%v", err)
            t.FailNow()
        }
        if diskInfo.(map[string]interface{})["state"] == "unattached" {
            break
        }
    }
    diskInfo, err := hc.DiskInfo(mDisk)
    if err != nil {
        t.Logf("Failed to grab disk info: \n%v", err)
        t.FailNow()
    }
    if diskInfo.(map[string]interface{})["state"] != "unattached" {
        t.Logf("Failed to create the new disk: \n(timeout)")
        t.FailNow()
    }

    //Actually, lets resize it to say 20 G
    diskMap = make(map[string]interface{})
    diskMap["size"] = 20
    _, err = hc.DiskResize(mDisk, diskMap)
    if err != nil {
        t.Logf("Failed to resize the new disk: \n%v", err)
        t.FailNow()
    }
    // Wait for resources to be up 
    end = time.Now().Add(time.Duration(30) * time.Second)
    for end.After(time.Now()){
        diskInfo, err := hc.DiskInfo(mDisk)
        if err != nil {
            t.Logf("Failed to grab disk info: \n%v", err)
            t.FailNow()
        }
        if diskInfo.(map[string]interface{})["state"] == "unattached" {
            break
        }
    }
    diskInfo, err = hc.DiskInfo(mDisk)
    if err != nil {
        t.Logf("Failed to grab disk info: \n%v", err)
        t.FailNow()
    }
    if diskInfo.(map[string]interface{})["state"] != "unattached" {
        t.Logf("Failed to create the new disk: \n(timeout)")
        t.FailNow()
    }

    //Now we need a boot disk for this instance
    //Search all templates for a Ubuntu 16.04
    var mTemplateId string
    templates, err := hc.TemplateList()
    if err != nil {
        t.Logf("Failed to list all templates: \n%v", errs)
        t.FailNow()
    }

    for _, t := range templates.([]interface{}) {
        if t.(map[string]interface{})["slug"].(string) == "ubuntu-16-04" && t.(map[string]interface{})["region"].(map[string]interface{})["id"].(string) == mRegion {
            mTemplateId = t.(map[string]interface{})["id"].(string)
            break
        }
    }

    if mTemplateId == "" {
        t.Logf("Failed to get template id for Ubuntu 16.04 in SY3")
        t.FailNow()
    }

    diskMap = make(map[string]interface{})
    diskMap["name"] = "hypercloud-test-boot-disk"
    diskMap["performance_tier"] = mDiskTier
    diskMap["region"] = mRegion
    diskMap["size"] = 10
    diskMap["template"] = mTemplateId

    var mBootDisk string
    bootDisk, err := hc.DiskCreate(diskMap)
    if err != nil {
        t.Logf("Unable to create the boot disk: \n%v", err)
        t.FailNow()
    }
    mBootDisk = bootDisk.(map[string]interface{})["id"].(string)
    // Wait for resources to be up 
    end = time.Now().Add(time.Duration(30) * time.Second)
    for end.After(time.Now()){
        diskInfo, err := hc.DiskInfo(mBootDisk)
        if err != nil {
            t.Logf("Failed to grab disk info: \n%v", err)
            t.FailNow()
        }
        if diskInfo.(map[string]interface{})["state"] == "unattached" {
            break
        }
    }
    diskInfo, err = hc.DiskInfo(mBootDisk)
    if err != nil {
        t.Logf("Failed to grab disk info: \n%v", err)
        t.FailNow()
    }
    if diskInfo.(map[string]interface{})["state"] != "unattached" {
        t.Logf("Failed to create the new disk: \n(timeout)")
        t.FailNow()
    }

    defer hc.DiskDelete(mBootDisk)

    //Now lets make a public/private IP for this guy

    //Public IP
    var mPubIp string
    netMap := make(map[string]interface{})
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
    mNetAdapter = netAdapter.(map[string]interface{})["id"].(string)

    end = time.Now().Add(time.Duration(30) * time.Second)
    for end.After(time.Now()){
        adapInfo, err := hc.NetworkInfo(mNetAdapter)
        if err != nil {
            t.Logf("Failed to grab network adapter info: \n %v", err)
            t.FailNow()
        }
        if adapInfo.(map[string]interface{})["state"] == "ready" {
            break
        }
    }
    adapInfo, err := hc.NetworkInfo(mNetAdapter)
    if err != nil {
        t.Logf("Failed to grab network info: \n%v")
        t.FailNow()
    }
    if adapInfo.(map[string]interface{})["state"] != "ready" {
        t.Logf("Failed to create new network adapter: \n(timeout)")
        t.FailNow()
    }

    defer hc.NetworkDelete(mNetAdapter)

    // Make a private IP
    netMap = make(map[string]interface{})
    netMap["name"] = "hypercloud-test-private-ip"
    netMap["network"] = mNetAdapter

    privIp, err := hc.IPAddressCreate(netMap)
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

    newInstance, err := hc.InstanceAssemble(instanceMap)
    if err != nil {
        t.Logf("Failed to create the new instance: \n%v", errs)
        t.FailNow()
    }

    mInstance = newInstance.(map[string]interface{})["id"].(string)

    // Wait for resources to be up 
    end = time.Now().Add(time.Duration(30) * time.Second)
    for end.After(time.Now()){
        instanceInfo, err := hc.InstanceInfo(mInstance)
        if err != nil {
            t.Logf("Failed to retrieve instance info: \n%v", err)
            t.FailNow()
        }
        if instanceInfo.(map[string]interface{})["state"] == "stopped" {
            break
        }
    }
    instanceInfo, err := hc.InstanceInfo(mInstance)
    if err != nil {
        t.Logf("Failed to retrieve instance info: \n%v", err)
        t.FailNow()
    }
    if instanceInfo.(map[string]interface{})["state"] != "stopped" {
        t.Logf("Failed to create the new instance: \n(Timeout)")
        t.FailNow()
    }

    defer hc.InstanceDelete(mInstance)

    //Attach disks/IP addresses to the guy
    updateInstance := make(map[string]interface{})

    updateInstance["disks"] = []string{mBootDisk, mDisk}

    var updateInstanceNA []interface{}
    // Add the private IP
    privateNetwork := make(map[string]interface{})
    privateNetwork["network"] = privIp.(map[string]interface{})["network_id"].(string)
    privateNetwork["ip_addresses"] = []string{mPrivIp}

    updateInstanceNA = append(updateInstanceNA, privateNetwork)

    //Add the public IP
    publicNetwork := make(map[string]interface{})
    publicNetwork["network"] = pubIp.(map[string]interface{})["network_id"].(string)
    publicNetwork["ip_addresses"] = []string{mPubIp}
    updateInstanceNA = append(updateInstanceNA, publicNetwork)

    //Call the update function
    updateInstance["network_adapters"] = updateInstanceNA

    newInstance, err = hc.InstanceUpdate(mInstance, updateInstance)
    if err != nil {
        t.Logf("Unable to update instance: \n%v", err)
        t.FailNow()
    }
}

