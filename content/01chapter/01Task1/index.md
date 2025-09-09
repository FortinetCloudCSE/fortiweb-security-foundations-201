---
title: "Task 1: Setup Azure Cloudshell"
linkTitle: "Task 1: Setup Azure Cloudshell"
weight: 10
---

|                            |    |  
|----------------------------| ----
| **Goal**                   | Login to Azure
| **Task**                   | Create an Azure account and login.
| **Verify task completion** | You will receive an email


{{% notice info %}} The below document references "student number" for a couple of the steps.  This is the first portion of the Username you received in the initial email with your Azure credentials.  For example if your username is **web10@fortinetcloud.onmicrosoft.com** then your student number would be **web10**  {{% /notice %}}


#### **Setup your AzureCloud Shell**

* Login to Azure Cloud Portal [https://portal.azure.com/](https://portal.azure.com/) with the provided login/password

    {{< figure src="cloudshell-01.png" alt="cloudshell1" >}}
    {{< figure src="cloudshell-02.png" alt="cloudshell2" >}}

* Select **Yes** when asked if you would like to stay signed in

    {{< figure src="cloudshell-03.png" alt="cloudshell3" >}}

* If you are presented with a "Welcome to Microsoft Azure" screen, click **Cancel**

    {{< figure src="cloudshell-04.jpg" alt="cloudshell4" >}}

* Click on Cloud Shell icon on the Top Right side of the portal

    {{< figure src="cloudshell-05.png" alt="cloudshell5" >}}

* Select **Bash**

    {{< figure src="cloudshell-06.png" alt="cloudshell6" >}}

* Next, you will see a "Getting started" page.
    * Select **Mount Storage Account**
    * Choose **Internal-Training** as the Storage account subscription
    * Click Apply

    {{< figure src="cloudshell-07.png" alt="cloudshell7" >}}

* On the Mount storage account  screen 
  * click **Select existing storage account**
  * click **Next**

    {{< figure src="cloudshell-08.png" alt="cloudshell8" >}}

* On the Select storage account screen (values in drop down)
  * choose **Internal-Training** as description
  * resource group will be **"student number"-appsec-102-workshop** 
  * storage account name will be "student number" followed by some random numbers and letters
  * File share will be **cloudshellshare**
  * Click **Select**

    {{< figure src="cloudshell-10.png" alt="cloudShell10" >}}

* Your Cloud shell is now configured.