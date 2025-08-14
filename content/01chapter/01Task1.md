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

    ![cloudshell1](../images/cloudshell-01.png)
    ![cloudshell2](../images/cloudshell-02.png)

* Select **Yes** when asked if you would like to stay signed in

    ![cloudshell3](../images/cloudshell-03.png)

* If you are presented with a "Welcome to Microsoft Azure" screen, click **Cancel**
    
    ![cloudshell4](../images/cloudshell-04.jpg
    )

* Click on Cloud Shell icon on the Top Right side of the portal

    ![cloudshell5](../images/cloudshell-05.png)

* Select **Bash**

    ![cloudshell6](../images/cloudshell-06.png)

* Next, you will see a "Getting started" page.
    * Select **Mount Storage Account**
    * Choose **Internal-Training** as the Storage account subscription
    * Click Apply

    ![cloudshell7](../images/cloudshell-07.png)

* On the Mount storage account  screen 
  * click **Select existing storage account**
  * click **Next**

    ![cloudshell8](../images/cloudshell-08.png)

* On the Select storage account screen (values in drop down)
  * choose **Internal-Training** as description
  * resource group will be **"student number"-appsec-102-workshop** 
  * storage account name will be "student number" followed by some random numbers and letters
  * File share will be **cloudshellshare**
  * Click **Select**

    ![cloudShell10](cloudshell-10.png)

* Your Cloud shell is now configured.